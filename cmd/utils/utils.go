package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/LorezV/url-shorter.git/cmd/config"
	"io"
	"net/http"
)

type ContextKey string

// GenerateID returns random string with len = 6
func GenerateID() (string, error) {
	b, err := GenerateRandom(6)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// GenerateRandom returns random string.
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// EncodeUserID encode string with user id by secret key from config.
func EncodeUserID(id string) []byte {
	h := hmac.New(sha256.New, []byte(config.AppConfig.SecretKey))
	h.Write([]byte(id))
	return h.Sum(nil)
}
