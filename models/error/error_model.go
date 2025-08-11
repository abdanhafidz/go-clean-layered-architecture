package http_error

import "errors"

var (
	BAD_REQUEST_ERROR     = errors.New("Invalid Request Format !")
	INTERNAL_SERVER_ERROR = errors.New("Internal Server Error!")
	UNAUTHORIZED          = errors.New("Unauthorized, you don't have permission to access this service!")
	DATA_NOT_FOUND        = errors.New("There is not data with given credential / given parameter!")
	TIMEOUT               = errors.New("Server took to long respond!")
	EXISTING_ACCOUNT      = errors.New("There is existing account!")
	INVALID_TOKEN         = errors.New("Invalid Authentication Payload!")
	DUPLICATE_DATA        = errors.New("Duplicate data !")
	ACCOUNT_NOT_FOUND     = errors.New("There is no account with given credential!")
	WRONG_PASSWORD        = errors.New("Your password is wrong for given account credential, please recheck!")
)
