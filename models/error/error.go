package http_error

import "errors"

var (
	BAD_REQUEST_ERROR           = errors.New("Invalid Request Format !")
	INTERNAL_SERVER_ERROR       = errors.New("Internal Server Error!")
	UNAUTHORIZED                = errors.New("Unauthorized, you don't have permission to access this service!")
	DATA_NOT_FOUND              = errors.New("There is not data with given credential / given parameter!")
	TIMEOUT                     = errors.New("Server took to long respond!")
	EXISTING_ACCOUNT            = errors.New("There is existing account!")
	INVALID_TOKEN               = errors.New("Invalid Authentication Payload!")
	DUPLICATE_DATA              = errors.New("Duplicate data !")
	ACCOUNT_NOT_FOUND           = errors.New("There is no account with given credential!")
	WRONG_PASSWORD              = errors.New("Your password is wrong for given account credential, please recheck!")
	INVALID_ACCOUNT_DIGITS      = errors.New("Your account 3 digits is not found in account number data")
	EXPIRED_TOKEN               = errors.New("Token expired")
	ALREADY_REGISTERED_TO_EVENT = errors.New("Account already registered to this event")
	EMAIL_ALREADY_EXISTS        = errors.New("Email already registered")
	NOT_REGISTERED_TO_EVENT     = errors.New("Account is not registered to this event")
	INVALID_OTP                 = errors.New("Invalid OTP Code")
	ERR_PROBLEM_SET_NOT_FOUND   = errors.New("problem set not found")
	ERR_QUESTION_NOT_FOUND      = errors.New("question not found")
	EVENT_FINISHED              = errors.New("The event has ended, you were disallowed to do the exam!")
	EVENT_NOT_STARTED           = errors.New("Take it easy, event hasn't starting yet! you cannot do the exam!")
	EXAMS_SUBMITTED             = errors.New("You've submitted the exam, you were diasallowed to answer the question!")
)
