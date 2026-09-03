package auth

import (
	"errors"
	"strconv"
	"time"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	UID uint64 `json:"uid"`
	TID uint64 `json:"tid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateTokenPair(userID, tenantID uint64, role string) (string, string, error) {
	cfg := config.AppCfg
	accessClaims := &Claims{
		UID: userID,
		TID: tenantID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatUint(userID, 10),
		},
	}
	refreshClaims := &Claims{
		UID: userID,
		TID: tenantID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatUint(userID, 10),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	access, err := accessToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	refresh, err := refreshToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.AppCfg
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func GetUserFromDB(userID uint64) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetTenantMembership(tenantID, userID uint64) (*model.TenantMember, error) {
	var member model.TenantMember
	if err := database.DB.Where("tenant_id = ? AND user_id = ? AND status = 'active'", tenantID, userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}
