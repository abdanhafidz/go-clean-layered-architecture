package provider

import (
    "context"
    "io"
)

type StorageProvider interface {
    UploadFile(ctx context.Context, file io.Reader, destinationPath string, contentType string) (string, error)
}
