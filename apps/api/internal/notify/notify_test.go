package notify_test

import (
	"context"
	"testing"
	"orbit/apps/api/internal/notify"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationService(t *testing.T) {
	svc := notify.NewNotificationService(nil)
	assert.NotNil(t, svc)
}

func TestSendNotification(t *testing.T) {
	svc := notify.NewNotificationService(nil)
	err := svc.Send(context.Background(), 1, "info", "Test", "Body", "")
	assert.Error(t, err)
}

func TestBroadcast(t *testing.T) {
	svc := notify.NewNotificationService(nil)
	svc.Broadcast(context.Background(), 1, "info", "Test", "Body")
}

func TestGetUnreadCount(t *testing.T) {
	count, err := notify.GetUnreadCount(999999)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
