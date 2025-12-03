package http_error

import "errors"

var (
	// ================= GENERAL =================
	BAD_REQUEST_ERROR     = errors.New("Invalid request format")
	INTERNAL_SERVER_ERROR = errors.New("Internal server error")
	TIMEOUT               = errors.New("Server took too long to respond")
	NOT_FOUND_ERROR       = errors.New("Resource not found")
	DUPLICATE_DATA        = errors.New("Duplicate data")
	INVALID_DATA_PAYLOAD  = errors.New("Invalid data payload provided")

	// ================= AUTH & ACCOUNT =================
	UNAUTHORIZED           = errors.New("Unauthorized, you don't have permission to access this service")
	EXISTING_ACCOUNT       = errors.New("Account already exists")
	INVALID_TOKEN          = errors.New("Invalid authentication payload")
	ACCOUNT_NOT_FOUND      = errors.New("There is no account with the given credentials")
	WRONG_PASSWORD         = errors.New("Invalid password, please check your credentials")
	INVALID_ACCOUNT_DIGITS = errors.New("Your account 3 digits is not found in account number data")
	EXPIRED_TOKEN          = errors.New("Token expired")
	INVALID_OTP            = errors.New("Invalid OTP code")
	EMAIL_ALREADY_EXISTS   = errors.New("Email already registered")

	// ================= EVENT & EXAM =================
	ALREADY_REGISTERED_TO_EVENT = errors.New("Account already registered to this event")
	EMAIL_ALREADY_EXISTS        = errors.New("Email already registered")
	NOT_REGISTERED_TO_EVENT     = errors.New("Account is not registered to this event")
	ERR_PROBLEM_SET_NOT_FOUND   = errors.New("Problem set not found")
	ERR_QUESTION_NOT_FOUND      = errors.New("Question not found")
	EVENT_FINISHED              = errors.New("The event has ended, you are disallowed to take the exam")
	EVENT_NOT_STARTED           = errors.New("Take it easy, event hasn't started yet! You cannot take the exam")
	EXAMS_SUBMITTED             = errors.New("You have submitted the exam, you are disallowed to answer the question")
	IMAGE_REQUIRED              = errors.New("Image is required")
	DESCRIPTION_REQUIRED        = errors.New("Description is required")
	CODE_REQUIRED               = errors.New("Code is required")

	// ================= FILE UPLOAD =================
	FILE_TOO_LARGE               = errors.New("File size exceeds the maximum limit")
	INVALID_FILE_TYPE            = errors.New("File type is not permitted for the selected context")
	UPLOAD_FAILED                = errors.New("Failed to upload file to storage provider")
	PARTIAL_UPLOAD_FAILURE       = errors.New("Some files failed validation or upload")
	INVALID_UPLOAD_CONTEXT_ERROR = errors.New("Invalid upload context")

	// ================= ACADEMY =================
	TITLE_REQUIRED       = errors.New("Title cannot be empty")
	SLUG_REQUIRED        = errors.New("Slug cannot be empty")
	ACADEMY_ID_REQUIRED  = errors.New("Academy ID is required")
	MATERIAL_ID_REQUIRED = errors.New("Material ID is required")

	ACADEMY_NOT_FOUND     = errors.New("Academy not found")
	MATERIAL_NOT_FOUND    = errors.New("Material not found")
	CONTENT_NOT_FOUND     = errors.New("Content not found")
	ACADEMY_HAS_MATERIALS = errors.New("Cannot delete academy because it still has materials")
	MATERIAL_HAS_CONTENTS = errors.New("Cannot delete material because it still has contents")
)
