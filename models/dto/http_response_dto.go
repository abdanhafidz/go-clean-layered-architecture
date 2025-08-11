package dto

type SuccessResponse[TResponse any] struct {
	Data     TResponse `json:"data"`
	Message  any       `json:"message"`
	MetaData any       `json:"meta_data"`
}

type ErrorResponse struct {
	Error    error `json:"errors"`
	Message  any   `json:"message"`
	MetaData any   `json:"meta_data"`
}
