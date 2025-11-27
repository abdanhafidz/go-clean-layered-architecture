package services

import (
    "context"
    "errors"
    "fmt"
    "io"
    "mime/multipart"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "github.com/google/uuid"

    entity "abdanhafidz.com/go-boilerplate/models/entity"
    http_error "abdanhafidz.com/go-boilerplate/models/error"
)

const MB = 1024 * 1024

type uploadConfig struct {
	MaxBytes    int64
	AllowedExts map[string]bool
	PathPrefix  string
	MaxCount    int
}

type StorageProvider interface {
	UploadFile(ctx context.Context, file io.Reader, destinationPath string, contentType string) (string, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.File, error)
}

type UploadService struct {
	storageProvider StorageProvider
	fileRepo        FileRepository
}

func NewUploadService(storage StorageProvider, repo FileRepository) *UploadService {
	return &UploadService{
		storageProvider: storage,
		fileRepo:        repo,
	}
}

func (s *UploadService) UploadFiles(ctx context.Context, files []*multipart.FileHeader, uploadContext string, accountID uuid.UUID) ([]entity.File, error) {
	config, err := s.getUploadConfig(uploadContext)
	if err != nil {
		return nil, err
	}

	if len(files) > config.MaxCount {
		return nil, http_error.INVALID_DATA_PAYLOAD
	}

	var uploadedFiles []entity.File
	var failedCount int
	var lastErr error

	for _, fileHeader := range files {
		fileEntity, err := s.processSingleFile(ctx, fileHeader, config, uploadContext, accountID)
		if err != nil {
			failedCount++
			lastErr = err
			continue
		}
		uploadedFiles = append(uploadedFiles, *fileEntity)
	}

	if failedCount > 0 {
		if len(uploadedFiles) > 0 {
			return uploadedFiles, http_error.PARTIAL_UPLOAD_FAILURE
		}
		return nil, lastErr
	}

	return uploadedFiles, nil
}

func (s *UploadService) GetFileByID(ctx context.Context, fileID uuid.UUID, accountID uuid.UUID) (*entity.File, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, http_error.NOT_FOUND_ERROR) {
			return nil, http_error.NOT_FOUND_ERROR
		}
		return nil, http_error.INTERNAL_SERVER_ERROR
	}

	if file == nil || file.AccountId != accountID {
		return nil, http_error.NOT_FOUND_ERROR
	}

	return file, nil
}

func (s *UploadService) processSingleFile(ctx context.Context, fileHeader *multipart.FileHeader, config uploadConfig, uploadContext string, accountID uuid.UUID) (*entity.File, error) {
	if err := s.validateFile(fileHeader, config); err != nil {
		return nil, err
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(fileHeader.Filename)))
	storedFilename := s.generateStoredFilename(fileHeader.Filename, ext)
	storagePath := s.generateStoragePath(config.PathPrefix, uploadContext, storedFilename, accountID)

	src, err := fileHeader.Open()
	if err != nil {
		return nil, http_error.INTERNAL_SERVER_ERROR
	}
	defer src.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	publicURL, err := s.storageProvider.UploadFile(ctx, src, storagePath, contentType)
	if err != nil {
		return nil, http_error.UPLOAD_FAILED
	}

	fileEntity := &entity.File{
		Id:           uuid.New(),
		OriginalName: fileHeader.Filename,
		StoredName:   storedFilename,
		MimeType:     contentType,
		Size:         fileHeader.Size,
		Path:         publicURL,
		Context:      uploadContext,
		AccountId:    accountID,
		CreatedAt:    time.Now(),
	}

	if err := s.fileRepo.Create(ctx, fileEntity); err != nil {
		return nil, http_error.INTERNAL_SERVER_ERROR
	}

	return fileEntity, nil
}

func (s *UploadService) validateFile(file *multipart.FileHeader, config uploadConfig) error {
	if file.Size == 0 || file.Size > config.MaxBytes {
		return http_error.FILE_TOO_LARGE
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(file.Filename)))
	if !config.AllowedExts[ext] {
		return http_error.INVALID_FILE_TYPE
	}

	blockedExts := map[string]bool{".exe": true, ".sh": true, ".bat": true, ".php": true}
	if blockedExts[ext] {
		return http_error.INVALID_FILE_TYPE
	}

	return nil
}

func (s *UploadService) generateStoredFilename(originalName string, ext string) string {
	originalNameWithoutExt := strings.TrimSuffix(originalName, ext)
	reg := regexp.MustCompile("[^a-zA-Z0-9._-]+")
	sanitizedRawName := reg.ReplaceAllString(originalNameWithoutExt, "_")

	uniqueID := uuid.New()
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d-%s%s", uniqueID.String(), timestamp, sanitizedRawName, ext)
}

func (s *UploadService) generateStoragePath(prefix, contextType, filename string, accountID uuid.UUID) string {
	switch contextType {
	case "submission":
		now := time.Now()
		return fmt.Sprintf("%s/%d/%02d/%s", prefix, now.Year(), now.Month(), filename)
	case "material":
		return fmt.Sprintf("%s/%s/%s", prefix, accountID.String(), filename)
	default:
		return fmt.Sprintf("%s/%s", prefix, filename)
	}
}

func (s *UploadService) getUploadConfig(contextType string) (uploadConfig, error) {
	codeExts := map[string]bool{".cpp": true, ".c": true, ".py": true, ".java": true, ".go": true, ".js": true, ".txt": true}
	imgExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	docExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true}

	allExts := make(map[string]bool)
	for k, v := range codeExts {
		allExts[k] = v
	}
	for k, v := range imgExts {
		allExts[k] = v
	}
	for k, v := range docExts {
		allExts[k] = v
	}

	switch contextType {
	case "image":
		return uploadConfig{
			MaxBytes:    10 * MB,
			AllowedExts: imgExts,
			PathPrefix:  "images",
			MaxCount:    5,
		}, nil
	case "material":
		return uploadConfig{
			MaxBytes:    10 * MB,
			AllowedExts: docExts,
			PathPrefix:  "materials",
			MaxCount:    1,
		}, nil
	case "submission":
		return uploadConfig{
			MaxBytes:    1 * MB,
			AllowedExts: codeExts,
			PathPrefix:  "submissions",
			MaxCount:    1,
		}, nil
	case "general":
		return uploadConfig{
			MaxBytes:    5 * MB,
			AllowedExts: allExts,
			PathPrefix:  "temp",
			MaxCount:    5,
		}, nil
	default:
		return uploadConfig{}, http_error.INVALID_DATA_PAYLOAD
	}
}

func (s *UploadService) UploadRawFile(ctx context.Context, reader io.Reader, originalName string, contentType string, uploadContext string, accountID uuid.UUID) (*entity.File, error) {
    config, err := s.getUploadConfig(uploadContext)
    if err != nil {
        return nil, http_error.BAD_REQUEST_ERROR
    }

    ext := strings.ToLower(strings.TrimSpace(filepath.Ext(originalName)))
    if ext == "" {
        if strings.Contains(contentType, "pdf") {
            ext = ".pdf"
        } else if strings.Contains(contentType, "png") {
            ext = ".png"
        } else if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
            ext = ".jpg"
        } else {
            ext = ".bin"
        }
    }

    storedFilename := s.generateStoredFilename(originalName, ext)
    storagePath := s.generateStoragePath(config.PathPrefix, uploadContext, storedFilename, accountID)

    publicURL, err := s.storageProvider.UploadFile(ctx, reader, storagePath, contentType)
    if err != nil {
        return nil, err
    }

    fileEntity := &entity.File{
        Id:           uuid.New(),
        OriginalName: originalName,
        StoredName:   storedFilename,
        MimeType:     contentType,
        Size:         0,
        Path:         publicURL,
        Context:      uploadContext,
        AccountId:    accountID,
        CreatedAt:    time.Now(),
    }

    if err := s.fileRepo.Create(ctx, fileEntity); err != nil {
        return nil, err
    }

    return fileEntity, nil
}