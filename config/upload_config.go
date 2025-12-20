package config

import (
	models "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
)

type UploadRule struct {
    MaxBytes    int64
    AllowedExts map[string]bool
    PathPrefix  string
    MaxCount    int
}

type UploadConfig interface {
    Get(contextType string) (UploadRule, error)
}

type uploadConfig struct{}

func NewUploadConfig() UploadConfig { return &uploadConfig{} }

func (c *uploadConfig) Get(contextType string) (UploadRule, error) {
    codeExts := map[string]bool{".cpp": true, ".c": true, ".py": true, ".java": true, ".go": true, ".js": true, ".txt": true}
    imgExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
    docExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true}

    allExts := make(map[string]bool)
    for k, v := range codeExts { allExts[k] = v }
    for k, v := range imgExts { allExts[k] = v }
    for k, v := range docExts { allExts[k] = v }

    switch contextType {
    case "image":
        return UploadRule{ MaxBytes: 10 * models.MB, AllowedExts: imgExts, PathPrefix: "images", MaxCount: 5 }, nil
    case "material":
        return UploadRule{ MaxBytes: 10 * models.MB, AllowedExts: docExts, PathPrefix: "materials", MaxCount: 1 }, nil
    case "submission":
        return UploadRule{ MaxBytes: 1 * models.MB, AllowedExts: codeExts, PathPrefix: "submissions", MaxCount: 1 }, nil
    case "general":
        return UploadRule{ MaxBytes: 5 * models.MB, AllowedExts: allExts, PathPrefix: "temp", MaxCount: 5 }, nil
    default:
        return UploadRule{}, http_error.INVALID_UPLOAD_CONTEXT_ERROR
    }
}

