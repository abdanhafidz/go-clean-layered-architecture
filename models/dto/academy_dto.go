package dto

import "github.com/google/uuid"

type CreateAcademyRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type UpdateAcademyRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type CreateMaterialRequest struct {
	AcademyId   uuid.UUID `json:"academy_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

type CreateContentRequest struct {
	MaterialId uuid.UUID `json:"material_id" binding:"required"`
	Title      string    `json:"title" binding:"required"`
	Contents   string    `json:"contents"`
}

// ================= RESPONSE DTOs =================

type AcademyProgressResponse struct {
	Id                      uuid.UUID `json:"id"`
	AccountId               uuid.UUID `json:"account_id"`
	AcademyId               uuid.UUID `json:"academy_id"`
	Status                  string    `json:"status"`
	Progress                float64   `json:"progress_percentage"`
	TotalCompletedMaterials uint      `json:"total_completed_materials"`
	CompletedAt             *string   `json:"completed_at"`
}

type AcademyContentResponse struct {
	Id     uuid.UUID `json:"id"`
	Order  uint      `json:"order"`
	Title  string    `json:"title"`
	Slug   string    `json:"slug"`
	Status string    `json:"status"`
}

type AcademyMaterialResponse struct {
	Id                     uuid.UUID                `json:"id"`
	Order                  uint                     `json:"order"`
	Title                  string                   `json:"title"`
	Slug                   string                   `json:"slug"`
	Status                 string                   `json:"status"`
	Progress               float64                  `json:"progress"`
	TotalCompletedContents uint                     `json:"total_completed_contents"`
	ContentsCount          int64                    `json:"contents_count"`
	Contents               []AcademyContentResponse `json:"contents"`
}

type AcademyDetailResponse struct {
	Id             uuid.UUID                 `json:"id"`
	Title          string                    `json:"title"`
	Slug           string                    `json:"slug"`
	Code           string                    `json:"code"`
	Description    string                    `json:"description"`
	ImageUrl       string                    `json:"image_url"`
	MaterialsCount int64                     `json:"materials_count"`
	UserProgress   *AcademyProgressResponse  `json:"user_progress"`
	Materials      []AcademyMaterialResponse `json:"materials"`
}

type MaterialProgressResponse struct {
	Id                     uuid.UUID `json:"id"`
	AccountId              uuid.UUID `json:"account_id"`
	AcademyId              uuid.UUID `json:"academy_id"`
	MaterialId             uuid.UUID `json:"material_id"`
	Progress               float64   `json:"progress_percentage"`
	TotalCompletedContents uint      `json:"total_completed_contents"`
	Status                 string    `json:"status"`
	CompletedAt            *string   `json:"completed_at"`
}

type ContentDetailResponse struct {
	Id     uuid.UUID `json:"id"`
	Order  uint      `json:"order"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
}

type MaterialDetailResponse struct {
	Id            uuid.UUID                 `json:"id"`
	AcademyId     uuid.UUID                 `json:"academy_id"`
	Title         string                    `json:"title"`
	Slug          string                    `json:"slug"`
	Description   string                    `json:"description"`
	Order         uint                      `json:"order"`
	ContentsCount int64                     `json:"contents_count"`
	Progress      *MaterialProgressResponse `json:"progress"`
	Contents      []ContentDetailResponse   `json:"contents"`
	Meta          map[string]string         `json:"meta"`
}
