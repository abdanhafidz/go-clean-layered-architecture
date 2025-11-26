package provider

import (
	"context"
	"mime/multipart"
)

type StorageProvider interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	GetFileURL(path string) (string, error)
	DeleteFile(path string) error
}