package dto

import (
	// Gunakan path yang SAMA PERSIS dengan yang ada di OptionsRequest
	entity "abdanhafidz.com/go-boilerplate/models/entity"

	"github.com/google/uuid"
	"time"
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

func FormatFileResponse(f *entity.File, baseURL string) FileResponse {
	fullURL := baseURL + f.Path

	return FileResponse{
		Id:           f.Id,
		OriginalName: f.OriginalName,
		URL:          fullURL,
		MimeType:     f.MimeType,
		Size:         f.Size,
		CreatedAt:    f.CreatedAt,
	}
}