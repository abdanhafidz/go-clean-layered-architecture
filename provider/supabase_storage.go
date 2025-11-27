package provider

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	storage_go "github.com/supabase-community/storage-go"
)

type StorageProvider interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	GetFileURL(path string) (string, error)
	DeleteFile(path string) error
}

type SupabaseStorage struct {
	client     *storage_go.Client
	bucketName string
	url        string
}

func NewSupabaseStorage(url string, key string, bucketName string) *SupabaseStorage {
	client := storage_go.NewClient(url+"/storage/v1", key, nil)
	return &SupabaseStorage{
		client:     client,
		bucketName: bucketName,
		url:        url,
	}
}

func (s *SupabaseStorage) UploadFile(ctx context.Context, file io.Reader, destinationPath string, contentType string) (string, error) {
	_, err := s.client.UploadFile(s.bucketName, destinationPath, file, storage_go.FileOptions{
		ContentType: &contentType, 
		Upsert:      new(bool), 
	})

	if err != nil {
		return "", err
	}
	publicURL := s.client.GetPublicUrl(s.bucketName, destinationPath).SignedURL
	if publicURL == "" {
		publicURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, s.bucketName, destinationPath)
	}

	return publicURL, nil
}
