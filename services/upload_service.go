package services

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	// Import package error kamu (sesuaikan path-nya)
	http_error "abdanhafidz.com/go-boilerplate/models/error"
)

type StorageProvider interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	GetFileURL(path string) (string, error)
}

type UploadService struct {
	storageProvider StorageProvider
}

func NewUploadService(storage StorageProvider) *UploadService {
	return &UploadService{storageProvider: storage}
}

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".pdf": true,
}

const MaxFileSize = 5 * 1024 * 1024 // 5 MB

func (s *UploadService) UploadFiles(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	var uploadedURLs []string

	for _, fileHeader := range files {
		// 1. Validasi Ukuran
		if fileHeader.Size > MaxFileSize {
			// Return error standar dari http_error
			return nil, http_error.FILE_TOO_LARGE
		}

		// 2. Validasi Ekstensi
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedExtensions[ext] {
			// Return error standar dari http_error
			return nil, http_error.INVALID_FILE_TYPE
		}

		src, err := fileHeader.Open()
		if err != nil {
			return nil, http_error.BAD_REQUEST_ERROR
		}

		// 3. Upload
		path, err := s.storageProvider.UploadFile(ctx, src, fileHeader)
		src.Close() // Tutup file
		
		if err != nil {
			// Jika error dari provider, kita bisa return error asli 
			// atau bungkus jadi INTERNAL_SERVER_ERROR / UPLOAD_FAILED
			return nil, http_error.UPLOAD_FAILED
		}

		// 4. Get URL
		url, err := s.storageProvider.GetFileURL(path)
		if err != nil {
			return nil, http_error.INTERNAL_SERVER_ERROR
		}

		uploadedURLs = append(uploadedURLs, url)
	}

	return uploadedURLs, nil
}