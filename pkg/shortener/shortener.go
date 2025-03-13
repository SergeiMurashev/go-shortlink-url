package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)


func GenerateShortLink() (string, error) {
	b := make([]byte, 6) 
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	shortLink := base64.RawURLEncoding.EncodeToString(b)
	return strings.TrimRight(shortLink, "="), nil
}