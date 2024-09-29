package utils

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

// noinspection SpellCheckingInspection
const (
	randomCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GetRandomString(length int) string {
	maxNumber := big.NewInt(int64(len(randomCharset)))
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		randomValue, err := rand.Int(rand.Reader, maxNumber)
		if err != nil {
			return ""
		}
		result[i] = randomCharset[randomValue.Int64()]
	}
	return string(result)
}

func IsValidUrl(value string) bool {
	if _, err := url.ParseRequestURI(value); err != nil {
		return false
	}
	parsedUrl, err := url.Parse(value)
	if err != nil || parsedUrl.Scheme == "" {
		return false
	}
	return true
}
