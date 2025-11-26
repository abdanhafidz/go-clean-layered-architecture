package dto

import "github.com/google/uuid"

type CreateAcademyRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type UpdateAcademyRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CreateMaterialRequest struct {
	AcademyId   uuid.UUID `json:"academy_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

type CreateContentRequest struct {
	AcademyMaterialId uuid.UUID `json:"academy_material_id" binding:"required"`
	Title             string    `json:"title" binding:"required"`
	Contents          string    `json:"contents"`
	Order             uint      `json:"order"`
}
