package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type AcademyService interface {
	// Academy
	CreateAcademy(ctx context.Context, req dto.CreateAcademyRequest) (entity.Academy, error)
	GetAcademy(ctx context.Context, accountId string, slug string) (entity.Academy, error)
	ListAcademies(ctx context.Context,accountId string) ([]entity.Academy, error)
	GetAcademyDetail(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	UpdateAcademy(ctx context.Context, id uuid.UUID, req dto.UpdateAcademyRequest) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error

	// Material
	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (entity.AcademyMaterial, error)
	GetMaterial(ctx context.Context, academySlug string, materialSlug string) (entity.AcademyMaterial, error)

	// Content
	CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error)
	GetContent(ctx context.Context, academySlug string, materialSlug string, order uint) (entity.AcademyContent, error)
	UpdateContentProgress(ctx context.Context, accountId string, academySlug string, materialSlug string, order uint) (entity.AcademyContentProgress, entity.AcademyMaterialProgress, entity.AcademyProgress, error)
	UpdateMaterialProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy, material entity.AcademyMaterial) (entity.AcademyMaterialProgress, error)
	UpdateAcademyProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy) (entity.AcademyProgress, error)
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
		Id:             uuid.New(),
		Title:          req.Title,
		Slug:           slugVal,
		Description:    req.Description,
		ImageUrl:       req.ImageUrl,
		MaterialsCount: 0,
	}

	return s.repo.CreateAcademy(ctx, a)
}

func (s *academyService) GetAcademy(ctx context.Context, accountId string, slug string) (entity.Academy, error) {
	return s.repo.GetAcademyWithProgress(ctx, accountId, slug)
}

func (s *academyService) ListAcademies(ctx context.Context,accountId string) ([]entity.Academy, error) {
	return s.repo.ListAcademy(ctx,accountId)
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

func (s *academyService) GetMaterial(ctx context.Context, academySlug string, materialSlug string) (entity.AcademyMaterial, error) {
	if strings.TrimSpace(academySlug) == "" || strings.TrimSpace(materialSlug) == "" {
		return entity.AcademyMaterial{}, errors.New("slug required")
	}
	academy, err := s.repo.GetAcademyBySlug(ctx, academySlug)
	if err != nil {
		return entity.AcademyMaterial{}, errors.New("academy not found: " + err.Error())
	}

	academyId := academy.Id
	return s.repo.GetMaterialBySlug(ctx, academyId, materialSlug)
}

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

	orderCount, _ := s.repo.CountMaterialsByAcademyID(ctx, req.AcademyId)
	order := uint(orderCount + 1)

	m := entity.AcademyMaterial{
		Id:            uuid.New(),
		AcademyId:     req.AcademyId,
		Title:         req.Title,
		Slug:          slugVal,
		Description:   req.Description,
		Order:         order,
		ContentsCount: 0,
	}

	// Update total materials in academy
	a, _ := s.repo.GetAcademyByID(ctx, req.AcademyId)
	a.MaterialsCount = a.MaterialsCount + 1
	s.repo.UpdateAcademy(ctx, a)

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
	count, _ := s.repo.CountContentsByMaterialID(ctx, req.AcademyMaterialId)
	order := uint(count + 1)

	c := entity.AcademyContent{
		Id:                uuid.New(),
		AcademyMaterialId: req.AcademyMaterialId,
		Title:             req.Title,
		Contents:          req.Contents,
		Order:             order,
	}

	// Update total progress in material
	m, _ := s.repo.GetMaterialByID(ctx, req.AcademyMaterialId)
	m.ContentsCount = m.ContentsCount + 1
	s.repo.UpdateMaterial(ctx, m)

	return s.repo.CreateContent(ctx, c)
}

func (s *academyService) GetContent(ctx context.Context, academySlug string, materialSlug string, order uint) (entity.AcademyContent, error) {
	material, err := s.GetMaterial(ctx, academySlug, materialSlug)
	if err != nil {
		return entity.AcademyContent{}, errors.New("material not found")
	}
	materialId := material.Id

	return s.repo.GetContentBySlug(ctx, materialId, order)
}

// Progress

func (s *academyService) UpdateContentProgress(ctx context.Context, accountId string, academySlug string, materialSlug string, order uint) (entity.AcademyContentProgress, entity.AcademyMaterialProgress, entity.AcademyProgress, error) {
	accountID, err := uuid.Parse(accountId)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("invalid account_id format")
	}

	academy, err := s.repo.GetAcademyBySlug(ctx, academySlug)

	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("academy not found")
	}

	material, err := s.repo.GetMaterialBySlug(ctx, academy.Id, materialSlug)

	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("material not found")
	}

	content, err := s.repo.GetContentBySlug(ctx, material.Id, order)

	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("content not found")
	}

	acp := entity.AcademyContentProgress{
		Id:                uuid.New(),  
		AccountId:         accountID,   
		AcademyId:         academy.Id,  
		AcademyMaterialId: material.Id, 
		ContentId:         content.Id,  
		Status:            "COMPLETED",
		CompletedAt:       time.Now(),
	}

	_, err = s.repo.UpsertContentProgress(ctx, acp)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("failed to upsert content progress: " + err.Error())
	}
	amp, err := s.UpdateMaterialProgress(ctx, accountID, academy, material)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("failed to update material progress: " + err.Error())
	}
	ap, err := s.UpdateAcademyProgress(ctx, accountID, academy)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, errors.New("failed to update academy progress: " + err.Error())
	}

	return acp, amp, ap, nil
}

func (s *academyService) UpdateMaterialProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy, material entity.AcademyMaterial) (entity.AcademyMaterialProgress, error) {
	//Count total completed contents for material progress update
	totalContentsCompleted, _ := s.repo.CountCompletedContentsByMaterialAndAccount(ctx, accountId, material.Id)
	m, err := s.repo.GetMaterialByID(ctx, material.Id)
	if err != nil {
		return entity.AcademyMaterialProgress{}, errors.New("material not found")
	}
	status := "IN_PROGRESS"

	if totalContentsCompleted == m.ContentsCount {
		status = "COMPLETED"
	}

	amp := entity.AcademyMaterialProgress{
		Id:                     uuid.New(),
		AccountId:              accountId,
		AcademyId:              academy.Id,
		AcademyMaterialId:      material.Id,
		Progress:               float64((float64(totalContentsCompleted) / float64(m.ContentsCount)) * 100),
		TotalCompletedContents: uint(totalContentsCompleted),
		Status:                 status,
	}

	_, err = s.repo.UpsertMaterialProgress(ctx, amp)
	return amp, err
}

func (s *academyService) UpdateAcademyProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy) (entity.AcademyProgress, error) {
	//Count total completed materials for academy progress update
	totalMaterialsCompleted, _ := s.repo.CountCompletedMaterialsByAcademyAndAccount(ctx, accountId, academy.Id)
	status := "IN_PROGRESS"

	if totalMaterialsCompleted == academy.MaterialsCount {
		status = "COMPLETED"
	}
	ap := entity.AcademyProgress{
		Id:                      uuid.New(),
		AccountId:               accountId,
		AcademyId:               academy.Id,
		Progress:                float64((float64(totalMaterialsCompleted) / float64(academy.MaterialsCount) * 100)),
		TotalCompletedMaterials: uint(totalMaterialsCompleted),
		Status:                  status,
	}
	_, err := s.repo.UpsertAcademyProgress(ctx, ap)
	return ap, err
}
