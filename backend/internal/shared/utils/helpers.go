package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID returns a new UUID v4 string.
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateCode generates a prefixed code like "MBR-20260316-ABCD".
func GenerateCode(prefix string) string {
	date := time.Now().Format("20060102")
	suffix := randomHex(4)
	return fmt.Sprintf("%s-%s-%s", strings.ToUpper(prefix), date, strings.ToUpper(suffix))
}

// GenerateNumber generates a numeric string of given length, e.g. "0001".
func GenerateNumber(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	num := ""
	for _, v := range b {
		num += fmt.Sprintf("%d", v%10)
	}
	return num[:length]
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
