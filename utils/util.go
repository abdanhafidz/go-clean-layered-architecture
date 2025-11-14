package utils

import (
	"time"

	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
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

	// kalau belum mulai (startTime > now), remaining = full duration
	if startTime.After(now) {
		return int(dueTime.Sub(startTime).Seconds())
	}

	remaining := int(dueTime.Sub(now).Seconds())

	if remaining < 0 {
		return 0
	}
	return remaining / 60
}
