package utils

import (
	"crypto/rand"
	"math/big"
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
