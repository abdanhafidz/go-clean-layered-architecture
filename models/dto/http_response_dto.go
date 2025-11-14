package dto

type SuccessResponse[TResponse any] struct {
	Status   string    `json:"status"`
	Data     TResponse `json:"data"`
	Message  any       `json:"message"`
	MetaData any       `json:"meta_data"`
}

type ErrorResponse struct {
	Status   string `json:"status"`
	Error    error  `json:"errors"`
	Message  any    `json:"message"`
	MetaData any    `json:"meta_data"`
}
