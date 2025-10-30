package utils

import (
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
