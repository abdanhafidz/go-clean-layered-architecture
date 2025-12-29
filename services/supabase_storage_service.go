package services

import (
	"context"
	"fmt"
	"io"
	"strings"

	http_error "abdanhafidz.com/go-boilerplate/models/error"
	storage_go "github.com/supabase-community/storage-go"
)

type StorageService interface {
	UploadFile(ctx context.Context, file io.Reader, destinationPath string, contentType string) (string, error)
}

type supabaseStorageService struct {
	client     *storage_go.Client
	bucketName string
	url        string
}

func NewSupabaseStorageService(url string, key string, bucketName string) StorageService {

	if url == "" || key == "" || bucketName == "" {
		fmt.Errorf(" supabase storage config is empty (url, key, and bucket are required)")
		return nil
	}
	if !strings.HasPrefix(url, "https://") || !strings.Contains(url, ".supabase.co") {
		fmt.Errorf("supabase storage url is invalid")
		return nil
	}
	if strings.Count(key, ".") != 2 {
		fmt.Errorf("supabase service key is not a valid compact JWS")
		return nil
	}

	client := storage_go.NewClient(url+"/storage/v1", key, nil)
	return &supabaseStorageService{client: client, bucketName: bucketName, url: url}
}

func (s *supabaseStorageService) UploadFile(ctx context.Context, file io.Reader, destinationPath string, contentType string) (string, error) {
	_, err := s.client.UploadFile(s.bucketName, destinationPath, file, storage_go.FileOptions{ContentType: &contentType, Upsert: new(bool)})
	if err != nil {
		return "", fmt.Errorf("%w: %v", http_error.UPLOAD_FAILED, err)
	}

	publicURL := s.client.GetPublicUrl(s.bucketName, destinationPath).SignedURL
	if publicURL == "" {
		publicURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, s.bucketName, destinationPath)
	}
	return publicURL, nil
}
