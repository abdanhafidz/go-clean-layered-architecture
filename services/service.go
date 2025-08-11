package services

import (
	"errors"

	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/utils"
	"gorm.io/gorm"
)

func RepoError(repo_err error) (err error) {
	if repo_err != nil {
		if errors.Is(repo_err, gorm.ErrDuplicatedKey) {
			return http_error.DUPLICATE_DATA

		} else if errors.Is(repo_err, gorm.ErrRecordNotFound) {
			return http_error.DATA_NOT_FOUND
		} else {
			utils.InternalErrorLog(repo_err)
			return http_error.INTERNAL_SERVER_ERROR
		}
	}
	return nil
}
