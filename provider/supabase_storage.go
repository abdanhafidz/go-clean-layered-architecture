package provider

import (
	"context"
	"fmt"
	"io"

	// Pastikan import library supabase client yang Anda pakai benar
	storage_go "github.com/supabase-community/storage-go" 
)

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
		Upsert:      new(bool), // Use new(bool) to create a pointer to false
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
