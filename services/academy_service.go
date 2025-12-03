package services

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils"
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
	GetContent(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	UpdateContentProgress(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContentProgress, entity.AcademyMaterialProgress, entity.AcademyProgress, error)

	GetAcademyResponse(ctx context.Context, accountId uuid.UUID, slug string) (*dto.AcademyDetailResponse, error)
	GetMaterialResponse(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string) (*dto.MaterialDetailResponse, error)
}
type academyService struct {
	repo repositories.AcademyRepository
}

func NewAcademyService(repo repositories.AcademyRepository) AcademyService {
	return &academyService{repo: repo}
}

func (s *academyService) GetAcademy(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error) {
	return s.academyRepo.GetAcademyWithProgress(ctx, accountId, slug)
}

func (s *academyService) GetAcademyDetail(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	a, _, err := s.academyRepo.GetAcademyWithMaterials(ctx, id)
	return a, err
}

func (s *academyService) CreateAcademy(ctx context.Context, req dto.CreateAcademyRequest) (entity.Academy, error) {
	if strings.TrimSpace(req.Title) == "" {
		return entity.Academy{}, http_error.TITLE_REQUIRED
	}

	if strings.TrimSpace(req.Code) == "" {
		return entity.Academy{}, http_error.CODE_REQUIRED
	}

	if strings.TrimSpace(req.Description) == "" {
		return entity.Academy{}, http_error.DESCRIPTION_REQUIRED
	}

	if strings.TrimSpace(req.ImageUrl) == "" {
		return entity.Academy{}, http_error.IMAGE_REQUIRED
	}

	slugVal := req.Slug
	if slugVal == "" {
		slugVal = slug.Make(req.Title)
	}

	if _, err := s.academyRepo.GetAcademyBySlug(ctx, slugVal); err == nil {
		return entity.Academy{}, http_error.DUPLICATE_DATA
	}

	a := entity.Academy{
		Id:             uuid.New(),
		Title:          req.Title,
		Slug:           slugVal,
		Code:           req.Code,
		Description:    req.Description,
		ImageUrl:       req.ImageUrl,
		MaterialsCount: 0,
	}

	return s.academyRepo.CreateAcademy(ctx, a)
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
	if req.ImageUrl != "" {
		existing.ImageUrl = req.ImageUrl
	}
	return s.academyRepo.UpdateAcademy(ctx, existing)
}

func (s *academyService) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	_, mats, err := s.repo.GetAcademyWithMaterials(ctx, id)
	if err != nil {
		return errors.New("academy not found")
	}
	if len(mats) > 0 {
		return http_error.ACADEMY_HAS_MATERIALS
	}
	return s.academyRepo.DeleteAcademy(ctx, id)
}

func (s *academyService) ListAcademies(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error) {
	return s.academyRepo.ListAcademy(ctx, accountId)
}

func (s *academyService) GetMaterial(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string) (entity.AcademyMaterial, error) {
	if strings.TrimSpace(academySlug) == "" || strings.TrimSpace(materialSlug) == "" {
		return entity.AcademyMaterial{}, http_error.SLUG_REQUIRED
	}
	academy, err := s.academyRepo.GetAcademyBySlug(ctx, academySlug)
	if err != nil {
		return entity.AcademyMaterial{}, http_error.ACADEMY_NOT_FOUND
	}
	return s.academyRepo.GetMaterialWithProgress(ctx, accountId, academy.Id, materialSlug)
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

	var createdMaterial entity.AcademyMaterial

	err := s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		orderCount, _ := txRepo.CountMaterialsByAcademyID(ctx, req.AcademyId)

		m := entity.AcademyMaterial{
			Id:            uuid.New(),
			AcademyId:     req.AcademyId,
			Title:         req.Title,
			Slug:          slugVal,
			Description:   req.Description,
			Order:         uint(orderCount + 1),
			ContentsCount: 0,
		}

		res, err := txRepo.CreateMaterial(ctx, m)
		if err != nil {
			return err
		}
		createdMaterial = res

		realCount, err := txRepo.CountMaterialsByAcademyID(ctx, req.AcademyId)
		if err != nil {
			return err
		}
		academy, _ := txRepo.GetAcademyByID(ctx, req.AcademyId)
		academy.MaterialsCount = int64(realCount)
		if _, err := txRepo.UpdateAcademy(ctx, academy); err != nil {
			return err
		}

		if err := txRepo.BatchRecalculateAcademyProgress(ctx, req.AcademyId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.AcademyMaterial{}, err
	}

	return createdMaterial, nil
}

func (s *academyService) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	m, err := s.academyRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return http_error.MATERIAL_NOT_FOUND
	}

	return s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		if err := txRepo.DeleteContentProgressByMaterialID(ctx, id); err != nil {
			return err
		}
		if err := txRepo.DeleteMaterialProgressByMaterialID(ctx, id); err != nil {
			return err
		}

		if err := txRepo.DeleteMaterial(ctx, id); err != nil {
			return err
		}

		if err := txRepo.DecrementMaterialOrdersGreaterThan(ctx, m.AcademyId, m.Order); err != nil {
			return err
		}

		realCount, err := txRepo.CountMaterialsByAcademyID(ctx, m.AcademyId)
		if err != nil {
			return err
		}

		academy, _ := txRepo.GetAcademyByID(ctx, m.AcademyId)
		academy.MaterialsCount = int64(realCount)
		if _, err := txRepo.UpdateAcademy(ctx, academy); err != nil {
			return err
		}

		if err := txRepo.BatchRecalculateAcademyProgress(ctx, m.AcademyId); err != nil {
			return err
		}

		return nil
	})
}

func (s *academyService) CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error) {
	if req.MaterialId == uuid.Nil {
		return entity.AcademyContent{}, http_error.MATERIAL_ID_REQUIRED
	}

	var createdContent entity.AcademyContent

	err := s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		m, err := txRepo.GetMaterialByID(ctx, req.MaterialId)
		if err != nil {
			return http_error.MATERIAL_NOT_FOUND
		}

		count, _ := txRepo.CountContentsByMaterialID(ctx, req.MaterialId)

		c := entity.AcademyContent{
			Id:         uuid.New(),
			MaterialId: req.MaterialId,
			Title:      req.Title,
			Contents:   req.Contents,
			Order:      uint(count + 1),
		}
		res, err := txRepo.CreateContent(ctx, c)
		if err != nil {
			return err
		}
		createdContent = res

		realCount, err := txRepo.CountContentsByMaterialID(ctx, req.MaterialId)
		if err != nil {
			return err
		}

		m.ContentsCount = realCount
		if _, err := txRepo.UpdateMaterial(ctx, m); err != nil {
			return err
		}

		if err := txRepo.BatchRecalculateMaterialProgress(ctx, req.MaterialId); err != nil {
			return err
		}

		if err := txRepo.BatchRecalculateAcademyProgress(ctx, m.AcademyId); err != nil {
			return err
		}

		return nil
	})

	return createdContent, err
}

func (s *academyService) GetContent(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContent, error) {
	material, err := s.GetMaterial(ctx, accountId, academySlug, materialSlug)
	if err != nil {
		return entity.AcademyContent{}, err
	}
	return s.academyRepo.GetContentWithProgress(ctx, accountId, material.AcademyId, material.Id, order)
}

func (s *academyService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	c, err := s.academyRepo.GetContentByID(ctx, id)
	if err != nil {
		return http_error.CONTENT_NOT_FOUND
	}

	return s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		if err := txRepo.DeleteContentProgressByContentID(ctx, id); err != nil {
			return err
		}
		if err := txRepo.DeleteContent(ctx, id); err != nil {
			return err
		}
		if err := txRepo.DecrementContentOrdersGreaterThan(ctx, c.MaterialId, c.Order); err != nil {
			return err
		}
		realCount, err := txRepo.CountContentsByMaterialID(ctx, c.MaterialId)
		if err != nil {
			return err
		}

		material, _ := txRepo.GetMaterialByID(ctx, c.MaterialId)
		material.ContentsCount = realCount
		if _, err := txRepo.UpdateMaterial(ctx, material); err != nil {
			return err
		}

		if err := txRepo.BatchRecalculateMaterialProgress(ctx, c.MaterialId); err != nil {
			return err
		}
		if err := txRepo.BatchRecalculateAcademyProgress(ctx, material.AcademyId); err != nil {
			return err
		}

		return nil
	})
}

func (s *academyService) UpdateContentProgress(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContentProgress, entity.AcademyMaterialProgress, entity.AcademyProgress, error) {
	academy, err := s.academyRepo.GetAcademyBySlug(ctx, academySlug)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, http_error.ACADEMY_NOT_FOUND
	}
	material, err := s.academyRepo.GetMaterialBySlug(ctx, academy.Id, materialSlug)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, http_error.MATERIAL_NOT_FOUND
	}
	content, err := s.academyRepo.GetContentBySlug(ctx, material.Id, order)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, http_error.CONTENT_NOT_FOUND
	}

	var acp entity.AcademyContentProgress
	var amp entity.AcademyMaterialProgress
	var ap entity.AcademyProgress

	err = s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		existingACP, _ := txRepo.GetContentProgress(ctx, accountId, academy.Id, material.Id, content.Id)
		acpID := existingACP.Id
		if acpID == uuid.Nil {
			acpID = uuid.New()
		}

		acp = entity.AcademyContentProgress{
			Id:          acpID,
			AccountId:   accountId,
			AcademyId:   academy.Id,
			MaterialId:  material.Id,
			ContentId:   content.Id,
			Status:      entity.StatusCompleted,
			CompletedAt: utils.Ptr(time.Now()),
		}
		if _, err := txRepo.UpsertContentProgress(ctx, acp); err != nil {
			return err
		}

		totalContentsCompleted, _ := txRepo.CountCompletedContentsByMaterialAndAccount(ctx, accountId, material.Id)
		m, _ := txRepo.GetMaterialByID(ctx, material.Id)

		matStatus := entity.StatusInProgress
		var matCompletedAt *time.Time
		progressPct := 0.0

		if m.ContentsCount > 0 {
			progressPct = (float64(totalContentsCompleted) / float64(m.ContentsCount)) * 100
			progressPct = math.Round(progressPct*100) / 100
			if totalContentsCompleted >= m.ContentsCount {
				matStatus = entity.StatusCompleted
				matCompletedAt = utils.Ptr(time.Now())
				progressPct = 100
			}
		} else {
			matStatus = entity.StatusCompleted
			progressPct = 100
		}

		existingAMP, _ := txRepo.GetMaterialProgress(ctx, accountId, academy.Id, material.Id)
		ampID := existingAMP.Id
		if ampID == uuid.Nil {
			ampID = uuid.New()
		}

		amp = entity.AcademyMaterialProgress{
			Id:                     ampID,
			AccountId:              accountId,
			AcademyId:              academy.Id,
			MaterialId:             material.Id,
			Progress:               progressPct,
			TotalCompletedContents: uint(totalContentsCompleted),
			Status:                 matStatus,
			CompletedAt:            matCompletedAt,
		}
		if _, err := txRepo.UpsertMaterialProgress(ctx, amp); err != nil {
			return err
		}

		accumulatedProgress, _ := txRepo.GetAccumulatedMaterialProgress(ctx, accountId, academy.Id)
		a, _ := txRepo.GetAcademyByID(ctx, academy.Id)

		acadStatus := entity.StatusNotStarted
		var acadCompletedAt *time.Time
		acadProgressPct := 0.0

		if a.MaterialsCount > 0 {
			acadProgressPct = accumulatedProgress / float64(a.MaterialsCount)
			acadProgressPct = math.Round(acadProgressPct*100) / 100

			if acadProgressPct >= 100 {
				acadStatus = entity.StatusCompleted
				acadCompletedAt = utils.Ptr(time.Now())
				acadProgressPct = 100
			} else if acadProgressPct > 0 {
				acadStatus = entity.StatusInProgress
			}
		}

		totalMaterialsCompleted, _ := txRepo.CountCompletedMaterialsByAcademyAndAccount(ctx, accountId, academy.Id)

		existingAP, _ := txRepo.GetAcademyProgress(ctx, accountId, academy.Id)
		apID := existingAP.Id
		if apID == uuid.Nil {
			apID = uuid.New()
		}

		ap = entity.AcademyProgress{
			Id:                      apID,
			AccountId:               accountId,
			AcademyId:               academy.Id,
			Progress:                acadProgressPct,
			TotalCompletedMaterials: uint(totalMaterialsCompleted),
			Status:                  acadStatus,
			CompletedAt:             acadCompletedAt,
		}
		if _, err := txRepo.UpsertAcademyProgress(ctx, ap); err != nil {
			return err
		}

		return nil
	})

	return acp, amp, ap, err
}

func (s *academyService) GetAcademyResponse(ctx context.Context, accountId uuid.UUID, slug string) (*dto.AcademyDetailResponse, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, http_error.SLUG_REQUIRED
	}

	academy, err := s.academyRepo.GetAcademyBySlug(ctx, slug)
	if err != nil {
		return nil, http_error.ACADEMY_NOT_FOUND
	}

	academyProgress, err := s.academyRepo.GetAcademyProgress(ctx, accountId, academy.Id)
	if err != nil {
		academyProgress = entity.AcademyProgress{Status: entity.StatusNotStarted}
	}

	materials, err := s.academyRepo.GetMaterialsWithContents(ctx, academy.Id)
	if err != nil {
		materials = []entity.AcademyMaterial{}
	}

	materialProgressMap, err := s.academyRepo.GetBatchMaterialProgress(ctx, accountId, academy.Id)
	if err != nil {
		materialProgressMap = make(map[uuid.UUID]entity.AcademyMaterialProgress)
	}

	var allContentIds []uuid.UUID
	for _, m := range materials {
		for _, c := range m.Contents {
			allContentIds = append(allContentIds, c.Id)
		}
	}
	contentProgressMap, err := s.academyRepo.GetContentProgressBatch(ctx, accountId, academy.Id, allContentIds)
	if err != nil {
		contentProgressMap = make(map[uuid.UUID]entity.AcademyContentProgress)
	}

	resp := &dto.AcademyDetailResponse{
		Id:             academy.Id,
		Title:          academy.Title,
		Slug:           academy.Slug,
		Code:           academy.Code,
		Description:    academy.Description,
		ImageUrl:       academy.ImageUrl,
		MaterialsCount: academy.MaterialsCount,
		UserProgress: &dto.AcademyProgressResponse{
			Id:                      academyProgress.Id,
			AccountId:               academyProgress.AccountId,
			AcademyId:               academyProgress.AcademyId,
			Status:                  academyProgress.Status,
			Progress:                academyProgress.Progress,
			TotalCompletedMaterials: academyProgress.TotalCompletedMaterials,
			CompletedAt:             utils.TimePtrToString(academyProgress.CompletedAt),
		},
	}

	dtMaterials := make([]dto.AcademyMaterialResponse, len(materials))
	for i, m := range materials {
		var matProg entity.AcademyMaterialProgress
		if p, ok := materialProgressMap[m.Id]; ok {
			matProg = p
		} else {
			matProg = entity.AcademyMaterialProgress{Status: entity.StatusNotStarted}
		}

		dtContents := make([]dto.AcademyContentResponse, len(m.Contents))
		for j, c := range m.Contents {
			cStatus := entity.StatusNotStarted
			if cp, ok := contentProgressMap[c.Id]; ok {
				cStatus = cp.Status
			}
			dtContents[j] = dto.AcademyContentResponse{
				Id:     c.Id,
				Order:  c.Order,
				Title:  c.Title,
				Status: cStatus,
			}
		}

		dtMaterials[i] = dto.AcademyMaterialResponse{
			Id:                     m.Id,
			Order:                  m.Order,
			Title:                  m.Title,
			Slug:                   m.Slug,
			Status:                 matProg.Status,
			Progress:               matProg.Progress,
			TotalCompletedContents: matProg.TotalCompletedContents,
			ContentsCount:          m.ContentsCount,
			Contents:               dtContents,
		}
	}
	resp.Materials = dtMaterials

	return resp, nil
}

func (s *academyService) GetMaterialResponse(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string) (*dto.MaterialDetailResponse, error) {
	if strings.TrimSpace(academySlug) == "" || strings.TrimSpace(materialSlug) == "" {
		return nil, http_error.SLUG_REQUIRED
	}

	academy, err := s.academyRepo.GetAcademyBySlug(ctx, academySlug)
	if err != nil {
		return nil, http_error.ACADEMY_NOT_FOUND
	}
	material, err := s.academyRepo.GetMaterialBySlug(ctx, academy.Id, materialSlug)
	if err != nil {
		return nil, http_error.MATERIAL_NOT_FOUND
	}

	materialProgress, _ := s.academyRepo.GetMaterialProgress(ctx, accountId, academy.Id, material.Id)

	_, contents, err := s.academyRepo.GetMaterialWithContents(ctx, material.Id)
	if err != nil {
		contents = []entity.AcademyContent{}
	}

	contentIds := make([]uuid.UUID, len(contents))
	for i, c := range contents {
		contentIds[i] = c.Id
	}
	contentProgressMap, err := s.academyRepo.GetContentProgressBatch(ctx, accountId, academy.Id, contentIds)
	if err != nil {
		contentProgressMap = make(map[uuid.UUID]entity.AcademyContentProgress)
	}

	resp := &dto.MaterialDetailResponse{
		Id:            material.Id,
		AcademyId:     material.AcademyId,
		Title:         material.Title,
		Slug:          material.Slug,
		Description:   material.Description,
		Order:         material.Order,
		ContentsCount: material.ContentsCount,
		Meta: map[string]string{
			"academy_slug":  academySlug,
			"material_slug": materialSlug,
		},
		Progress: &dto.MaterialProgressResponse{
			Id:                     materialProgress.Id,
			AccountId:              materialProgress.AccountId,
			AcademyId:              materialProgress.AcademyId,
			MaterialId:             materialProgress.MaterialId,
			Progress:               materialProgress.Progress,
			TotalCompletedContents: materialProgress.TotalCompletedContents,
			Status:                 materialProgress.Status,
			CompletedAt:            utils.TimePtrToString(materialProgress.CompletedAt),
		},
	}

	dtContents := make([]dto.ContentDetailResponse, len(contents))
	for i, c := range contents {
		cStatus := entity.StatusNotStarted
		if cp, ok := contentProgressMap[c.Id]; ok {
			cStatus = cp.Status
		}

		dtContents[i] = dto.ContentDetailResponse{
			Id:     c.Id,
			Order:  c.Order,
			Title:  c.Title,
			Status: cStatus,
		}
	}
	resp.Contents = dtContents

	return resp, nil
}
