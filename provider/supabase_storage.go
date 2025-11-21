package provider

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorage struct {
	client      *storage_go.Client
	bucketName  string
	supabaseURL string
}

func NewSupabaseStorage(url, secretKey, bucketName string) *SupabaseStorage {
	client := storage_go.NewClient(url+"/storage/v1", secretKey, nil)

	return &SupabaseStorage{
		client:      client,
		bucketName:  bucketName,
		supabaseURL: url,
	}
}

func (s *SupabaseStorage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)

	contentType := header.Header.Get("Content-Type")

	_, err := s.client.UploadFile(s.bucketName, filename, file, storage_go.FileOptions{
		ContentType: &contentType,
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to upload to supabase: %v", err)
	}

	return filename, nil
}

func (s *SupabaseStorage) GetFileURL(path string) (string, error) {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.supabaseURL, s.bucketName, path), nil
}

func (s *SupabaseStorage) DeleteFile(path string) error {
	_, err := s.client.RemoveFile(s.bucketName, []string{path})
	return err
}