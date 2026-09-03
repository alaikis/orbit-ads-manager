package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"orbit/apps/api/config"
)

func Encrypt(plaintext string) (string, error) {
	cfg := config.AppCfg
	key := []byte(cfg.Encryption.MasterKey)
	if len(key) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return fmt.Sprintf("v1:%s", base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func Decrypt(ciphertext string) (string, error) {
	cfg := config.AppCfg
	key := []byte(cfg.Encryption.MasterKey)
	if len(key) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	if len(ciphertext) < 4 || ciphertext[:3] != "v1:" {
		return "", fmt.Errorf("invalid ciphertext format")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext[3:])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
