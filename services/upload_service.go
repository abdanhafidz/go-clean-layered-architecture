package services

import (
	"context"
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

type StorageProvider interface {
	UploadFile(ctx context.Context, file io.Reader, destinationPath string, contentType string) (string, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
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

type contextConfig struct {
	MaxBytes    int64
	AllowedExts map[string]bool
	PathPrefix  string
}

func (s *UploadService) getConfig(uploadContext string) (contextConfig, error) {
	switch uploadContext {
	case "avatar":
		return contextConfig{
			MaxBytes:    5 * 1024 * 1024,
			AllowedExts: map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true},
			PathPrefix:  "avatars",
		}, nil
	case "material":
		// Documents (.pdf) Max Size 10 MB
		return contextConfig{
			MaxBytes:    10 * 1024 * 1024,
			AllowedExts: map[string]bool{".pdf": true},
			PathPrefix:  "materials",
		}, nil
	case "submission":
		// Source Code Max Size 1 MB
		return contextConfig{
			MaxBytes:    1 * 1024 * 1024,
			AllowedExts: map[string]bool{".cpp": true, ".c": true, ".py": true, ".java": true, ".go": true, ".js": true, ".txt": true},
			PathPrefix:  "submissions",
		}, nil
	case "general":
		return contextConfig{
			MaxBytes: 5 * 1024 * 1024,
			AllowedExts: map[string]bool{
				".jpg": true, ".png": true, ".pdf": true, ".c": true,
				".cpp": true, ".py": true, ".go": true, ".js": true,
				".txt": true,
			},
			PathPrefix: "temp",
		}, nil
	default:
		return contextConfig{}, fmt.Errorf("invalid context")
	}
}

func (s *UploadService) UploadFiles(ctx context.Context, files []*multipart.FileHeader, uploadContext string, accountID uuid.UUID) ([]entity.File, error) {
	var uploadedFiles []entity.File
	var failedCount int // Counter untuk kegagalan (Partial Failure)

	config, err := s.getConfig(uploadContext)
	if err != nil {
		// Konteks tidak valid adalah kesalahan permintaan yang fatal (400)
		return nil, http_error.BAD_REQUEST_ERROR
	}

	// Check Max Files/Req. Submission/Material = 1, Images/General = 5
	if uploadContext == "submission" || uploadContext == "material" {
		if len(files) > 1 {
			return nil, http_error.BAD_REQUEST_ERROR // Melanggar Max Files/Req
		}
	}
	// Asumsi: Max Files/Req Images/General adalah 5, jika > 5 ini akan dianggap BAD_REQUEST_ERROR
	if (uploadContext == "avatar" || uploadContext == "general") && len(files) > 5 {
		return nil, http_error.BAD_REQUEST_ERROR
	}

	for _, fileHeader := range files {
		// Logika kegagalan validasi file tunggal (non-atomic)

		// 1. Check Empty File (0-byte file)
		if fileHeader.Size == 0 {
			failedCount++
			// Di sini, idealnya kita akan mencatat error detail, tapi untuk menjaga signature kita hanya menghitung.
			continue
		}

		// 2. Check Size Limit
		if fileHeader.Size > config.MaxBytes {
			failedCount++
			continue
		}

		// 3. Check Extension Validity (Allowed Exts)
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !config.AllowedExts[ext] {
			failedCount++
			continue
		}

		// 4. Blocklist check (Explicitly reject executables)
		// Note: Blocklist untuk .svg tidak diterapkan di sini karena dapat diperbolehkan jika disanitasi.
		if ext == ".exe" || ext == ".sh" || ext == ".bat" || ext == ".php" {
			failedCount++
			continue
		}

		// Filename Sanitization & Naming Strategy
		originalNameWithoutExt := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
		reg := regexp.MustCompile("[^a-zA-Z0-9._-]+")
		sanitizedRawName := reg.ReplaceAllString(originalNameWithoutExt, "_")

		// Storage Naming Strategy: {UUID}-{Timestamp}-{SanitizedFilename}.{Ext}
		uniqueID := uuid.New()
		timestamp := time.Now().Unix()
		storedFilename := fmt.Sprintf("%s-%d-%s%s", uniqueID.String(), timestamp, sanitizedRawName, ext)

		// 🎯 PERBAIKAN LOGIKA PENENTUAN storagePath
		var storagePath string
		now := time.Now()

		switch uploadContext {
		case "submission":
			// Structure: /submissions/{year}/{month}/{uuid}-{filename}.cpp
			storagePath = fmt.Sprintf("%s/%d/%02d/%s", config.PathPrefix, now.Year(), now.Month(), storedFilename)
		case "material":
			// Structure: /materials/{material_id}/{uuid}-{filename}.pdf
			// ASUMSI: Menggunakan accountID sebagai {material_id}
			storagePath = fmt.Sprintf("%s/%s/%s", config.PathPrefix, accountID.String(), storedFilename)
		case "avatar":
			// Structure: /avatars/{uuid}-{filename}.jpg
			storagePath = fmt.Sprintf("%s/%s", config.PathPrefix, storedFilename)
		case "general":
			// Structure: /temp/{filename}
			storagePath = fmt.Sprintf("%s/%s", config.PathPrefix, storedFilename)
		default:
			// Fallback (Seharusnya tidak tercapai karena sudah divalidasi oleh getConfig)
			storagePath = fmt.Sprintf("unknown/%s", storedFilename)
		}

		// Open file stream
		src, err := fileHeader.Open()
		if err != nil {
			failedCount++
			continue // Lanjut ke file berikutnya
		}

		// Upload to Storage
		contentType := fileHeader.Header.Get("Content-Type")
		publicURL, err := s.storageProvider.UploadFile(ctx, src, storagePath, contentType)
		src.Close() // Pastikan stream ditutup setelah upload

		if err != nil {
			failedCount++
			continue // Lanjut ke file berikutnya
		}

		// Create Entity
		fileEntity := entity.File{
			Id:           uniqueID,
			OriginalName: fileHeader.Filename,
			StoredName:   storedFilename,
			MimeType:     contentType,
			Size:         fileHeader.Size,
			Path:         publicURL,
			Context:      uploadContext,
			AccountId:    accountID,
			CreatedAt:    time.Now(),
		}

		// Save to Database
		if err := s.fileRepo.Create(ctx, &fileEntity); err != nil {

			return nil, http_error.INTERNAL_SERVER_ERROR
		}

		uploadedFiles = append(uploadedFiles, fileEntity)
	}

	if failedCount > 0 {
		if len(uploadedFiles) > 0 {
			return uploadedFiles, http_error.PARTIAL_UPLOAD_FAILURE
		} else {
			return nil, http_error.INVALID_FILE_TYPE
		}
	}

	return uploadedFiles, nil
}
