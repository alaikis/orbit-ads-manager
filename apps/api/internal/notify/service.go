package notify

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type NotificationService struct {
	config *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{config: cfg}
}

func (s *NotificationService) Send(ctx context.Context, userID uint64, notifType, title, body, link string) error {
	n := model.Notification{
		UserID: userID,
		Type:   notifType,
		Title:  title,
		Body:   body,
		Link:   link,
	}
	if err := database.DB.Create(&n).Error; err != nil {
		return err
	}

	go s.sendEmailIfNeeded(userID, notifType, title, body)
	return nil
}

func (s *NotificationService) sendEmailIfNeeded(userID uint64, notifType, title, body string) {
	if notifType != "critical" {
		return
	}

	cfg := s.config
	if cfg.SMTP.Host == "" {
		return
	}

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return
	}

	auth := smtp.PlainAuth(cfg.SMTP.User, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.Host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", user.Email, title, body))

	if err := smtp.SendMail(fmt.Sprintf("%s:%s", cfg.SMTP.Host, cfg.SMTP.Port), auth, cfg.SMTP.User, []string{user.Email}, msg); err != nil {
		log.Printf("failed to send email: %v", err)
		return
	}

	logEntry := model.EmailLog{
		ToEmail:  user.Email,
		Subject:  title,
		Template: "default",
		Status:   "sent",
	}
	database.DB.Create(&logEntry)
}

func (s *NotificationService) Broadcast(ctx context.Context, tenantID uint64, notifType, title, body string) {
	var members []model.TenantMember
	database.DB.Where("tenant_id = ? AND status = 'active'", tenantID).Find(&members)
	for _, m := range members {
		_ = s.Send(ctx, m.UserID, notifType, title, body, "")
	}
}

func (s *NotificationService) MarkRead(userID uint64, notificationID uint64) error {
	return database.DB.Model(&model.Notification{}).Where("id = ? AND user_id = ?", notificationID, userID).Update("read_at", "now()").Error
}

func GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	return count, database.DB.Model(&model.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&count).Error
}
