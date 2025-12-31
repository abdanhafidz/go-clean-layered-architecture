package utils

import (
	"os"
	"regexp"
	"strings"
	"time"

	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"github.com/google/uuid"
)

func ToUUID(s any) (uuid.UUID, error) {
	sStr, ok := s.(string)
	if !ok {
		return uuid.UUID{}, http_error.INTERNAL_SERVER_ERROR
	}

	res, err := uuid.Parse(sStr)
	if err != nil {
		return uuid.UUID{}, http_error.INTERNAL_SERVER_ERROR
	}

	return res, nil
}
func CalculateRemainingTime(startTime, dueTime time.Time) int {
	now := time.Now()
	if startTime.After(now) {
		return int(dueTime.Sub(startTime).Seconds())
	}
	remaining := int(dueTime.Sub(now).Seconds())
	if remaining < 0 {
		return 0
	}
	return remaining / 60
}

func Ptr[T any](v T) *T {
	return &v
}

func TimePtrToString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func ValidateCode(code string) error {
	var CodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{6,12}$`)
	if !CodeRegex.MatchString(code) {
		return http_error.INVALID_CODE
	}
	return nil
}

func GetEnv(key string) string {
	// 1. Normal env
	if val := os.Getenv(key); val != "" {
		return val
	}

	// 2. File-based secret
	if file := os.Getenv(key + "_FILE"); file != "" {
		data, err := os.ReadFile(file)
		if err == nil {
			return strings.TrimSpace(string(data))
		} else {
			panic(err)
		}
	}

	return ""
}
