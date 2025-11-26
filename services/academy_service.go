package services

import (
	"context"
	"errors"
	"strings"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type AcademyService interface {
	// Academy
	CreateAcademy(ctx context.Context, req dto.CreateAcademyRequest) (entity.Academy, error)
	GetAcademy(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	ListAcademies(ctx context.Context) ([]entity.Academy, error)
	GetAcademyDetail(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	UpdateAcademy(ctx context.Context, id uuid.UUID, req dto.UpdateAcademyRequest) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error

	// Material
	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (entity.AcademyMaterial, error)

	// Content
	CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error)
}
type academyService struct {
	repo repositories.AcademyRepository
}

func NewAcademyService(repo repositories.AcademyRepository) AcademyService {
	return &academyService{repo: repo}
}

//
// ===== Academy =====
//

func (s *academyService) CreateAcademy(ctx context.Context, req dto.CreateAcademyRequest) (entity.Academy, error) {
	if strings.TrimSpace(req.Title) == "" {
		return entity.Academy{}, errors.New("title required")
	}

	slugVal := req.Slug
	if slugVal == "" {
		slugVal = slug.Make(req.Title)
	}

	if _, err := s.repo.GetAcademyBySlug(ctx, slugVal); err == nil {
		return entity.Academy{}, errors.New("slug already exists")
	}

	a := entity.Academy{
		Title:       req.Title,
		Slug:        slugVal,
		Description: req.Description,
	}

	return s.repo.CreateAcademy(ctx, a)
}

func (s *academyService) GetAcademy(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	return s.repo.GetAcademyByID(ctx, id)
}

func (s *academyService) ListAcademies(ctx context.Context) ([]entity.Academy, error) {
	return s.repo.ListAcademy(ctx)
}

func (s *academyService) GetAcademyDetail(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	a, _, err := s.repo.GetAcademyWithMaterials(ctx, id)
	return a, err
}

func (s *academyService) UpdateAcademy(ctx context.Context, id uuid.UUID, req dto.UpdateAcademyRequest) (entity.Academy, error) {
	existing, err := s.repo.GetAcademyByID(ctx, id)
	if err != nil {
		return entity.Academy{}, errors.New("academy not found")
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}

	if req.Slug != "" {
		existing.Slug = req.Slug
	} else {
		existing.Slug = slug.Make(existing.Title)
	}

	return s.repo.UpdateAcademy(ctx, existing)
}

func (s *academyService) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	_, mats, err := s.repo.GetAcademyWithMaterials(ctx, id)
	if err != nil {
		return errors.New("academy not found")
	}
	if len(mats) > 0 {
		return errors.New("cannot delete academy with materials")
	}
	return s.repo.DeleteAcademy(ctx, id)
}

//
// ===== Material =====
//

func (s *academyService) CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (entity.AcademyMaterial, error) {
	if req.AcademyId == uuid.Nil {
		return entity.AcademyMaterial{}, errors.New("academy_id required")
	}
	if _, err := s.repo.GetAcademyByID(ctx, req.AcademyId); err != nil {
		return entity.AcademyMaterial{}, errors.New("academy not found")
	}

	slugVal := req.Slug
	if slugVal == "" {
		slugVal = slug.Make(req.Title)
	}

	m := entity.AcademyMaterial{
		AcademyId:   req.AcademyId,
		Title:       req.Title,
		Slug:        slugVal,
		Description: req.Description,
	}

	return s.repo.CreateMaterial(ctx, m)
}

//
// ===== Content =====
//

func (s *academyService) CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error) {
	if req.AcademyMaterialId == uuid.Nil {
		return entity.AcademyContent{}, errors.New("academy_material_id required")
	}

	if _, err := s.repo.GetMaterialByID(ctx, req.AcademyMaterialId); err != nil {
		return entity.AcademyContent{}, errors.New("material not found")
	}

	// auto order last++
	list, _ := s.repo.ListContentsByMaterialID(ctx, req.AcademyMaterialId)
	order := req.Order
	if order == 0 {
		order = uint(len(list) + 1)
	}

	c := entity.AcademyContent{
		AcademyMaterialId: req.AcademyMaterialId,
		Title:             req.Title,
		Contents:          req.Contents,
		Order:             order,
	}

	return s.repo.CreateContent(ctx, c)
}
