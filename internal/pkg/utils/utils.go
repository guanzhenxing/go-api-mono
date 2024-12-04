package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

// GenerateRequestID 生成请求ID
func GenerateRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b)
}
