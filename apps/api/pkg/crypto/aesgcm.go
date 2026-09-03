package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"orbit/apps/api/config"
)

var (
	ErrInvalidKeyLength = errors.New("crypto: master key must be 32 bytes")
	ErrCiphertextTooShort = errors.New("crypto: ciphertext too short")
)

func gcm() (cipher.AEAD, error) {
	key := config.AppCfg.Encryption.MasterKey
	if len(key) != 32 {
		return nil, ErrInvalidKeyLength
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func Encrypt(plaintext []byte) (ciphertext []byte, nonce []byte, err error) {
	g, err := gcm()
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, g.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext = g.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func Decrypt(ciphertext []byte, nonce []byte) ([]byte, error) {
	g, err := gcm()
	if err != nil {
		return nil, err
	}
	if len(ciphertext) == 0 || len(nonce) != g.NonceSize() {
		return nil, ErrCiphertextTooShort
	}
	return g.Open(nil, nonce, ciphertext, nil)
}

func EncryptString(plaintext string) (string, string, error) {
	ct, nonce, err := Encrypt([]byte(plaintext))
	if err != nil {
		return "", "", err
	}
	return string(ct), string(nonce), nil
}

func DecryptString(ciphertext string, nonce string) (string, error) {
	pt, err := Decrypt([]byte(ciphertext), []byte(nonce))
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

func Last4(s string) string {
	if len(s) <= 4 {
		return s
	}
	return s[len(s)-4:]
}
