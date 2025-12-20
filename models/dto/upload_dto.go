package dto

import (
	"time"

	"github.com/google/uuid"
)

type FileResponse struct {
	Id           uuid.UUID `json:"id"`
	OriginalName string    `json:"original_name"`
	URL          string    `json:"url"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

type FileUploadResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []FileResponse `json:"data"`
}

type FileResponseSingle struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    FileResponse `json:"data"`
}
