package services

import (
	"context"
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
	GetAcademy(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error)
	GetAcademyDetail(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	CreateAcademy(ctx context.Context, req dto.CreateAcademyRequest) (entity.Academy, error)
	UpdateAcademy(ctx context.Context, id uuid.UUID, req dto.UpdateAcademyRequest) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error
	ListAcademies(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error)

	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (entity.AcademyMaterial, error)
	GetMaterial(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string) (entity.AcademyMaterial, error)
	DeleteMaterial(ctx context.Context, id uuid.UUID) error

	CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error)
	GetContent(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	UpdateContentProgress(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContentProgress, entity.AcademyMaterialProgress, entity.AcademyProgress, error)
	UpdateMaterialProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy, material entity.AcademyMaterial) (entity.AcademyMaterialProgress, error)
	UpdateAcademyProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy) (entity.AcademyProgress, error)

	sanitizeAcademyProgress(a entity.Academy, startedMaterials int64, accumulatedProgress float64) entity.Academy
	sanitizeMaterialProgress(m entity.AcademyMaterial) entity.AcademyMaterial
}

type academyService struct {
	repo repositories.AcademyRepository
}

func NewAcademyService(repo repositories.AcademyRepository) AcademyService {
	return &academyService{repo: repo}
}

// ================= ACADEMY =================

func (s *academyService) GetAcademy(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error) {
	a, err := s.repo.GetAcademyWithProgress(ctx, accountId, slug)
	if err != nil {
		return a, err
	}
	// Ambil data real-time untuk sanitasi
	startedCount, _ := s.repo.CountStartedMaterialsByAcademyAndAccount(ctx, accountId, a.Id)
	accumulatedProgress, _ := s.repo.GetAccumulatedMaterialProgress(ctx, accountId, a.Id)

	return s.sanitizeAcademyProgress(a, startedCount, accumulatedProgress), nil
}

func (s *academyService) GetAcademyDetail(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	a, _, err := s.repo.GetAcademyWithMaterials(ctx, id)
	return a, err
}

func (s *academyService) CreateAcademy(ctx context.Context, req dto.CreateAcademyRequest) (entity.Academy, error) {
	if strings.TrimSpace(req.Title) == "" {
		return entity.Academy{}, http_error.TITLE_REQUIRED
	}
	slugVal := req.Slug
	if slugVal == "" {
		slugVal = slug.Make(req.Title)
	}
	if _, err := s.repo.GetAcademyBySlug(ctx, slugVal); err == nil {
		return entity.Academy{}, http_error.DUPLICATE_DATA
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

func (s *academyService) UpdateAcademy(ctx context.Context, id uuid.UUID, req dto.UpdateAcademyRequest) (entity.Academy, error) {
	existing, err := s.repo.GetAcademyByID(ctx, id)
	if err != nil {
		return entity.Academy{}, http_error.ACADEMY_NOT_FOUND
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Slug != "" {
		existing.Slug = req.Slug
	}
	return s.repo.UpdateAcademy(ctx, existing)
}

func (s *academyService) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	_, mats, err := s.repo.GetAcademyWithMaterials(ctx, id)
	if err != nil {
		return http_error.ACADEMY_NOT_FOUND
	}
	if len(mats) > 0 {
		return http_error.ACADEMY_HAS_MATERIALS
	}
	return s.repo.DeleteAcademy(ctx, id)
}

func (s *academyService) ListAcademies(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error) {
	list, err := s.repo.ListAcademy(ctx, accountId)
	if err != nil {
		return nil, err
	}

	for i := range list {
		// Note: Untuk list, kita bisa pass 0 jika ingin performa cepat,
		// atau ambil data real jika ingin akurat seperti detail.
		// Disini kita ambil real data (hati-hati N+1 jika data ribuan).
		startedCount, _ := s.repo.CountStartedMaterialsByAcademyAndAccount(ctx, accountId, list[i].Id)
		accumulated, _ := s.repo.GetAccumulatedMaterialProgress(ctx, accountId, list[i].Id)

		list[i] = s.sanitizeAcademyProgress(list[i], startedCount, accumulated)
	}

	return list, nil
}

// ================= MATERIAL =================

func (s *academyService) GetMaterial(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string) (entity.AcademyMaterial, error) {
	if strings.TrimSpace(academySlug) == "" || strings.TrimSpace(materialSlug) == "" {
		return entity.AcademyMaterial{}, http_error.SLUG_REQUIRED
	}
	academy, err := s.repo.GetAcademyBySlug(ctx, academySlug)
	if err != nil {
		return entity.AcademyMaterial{}, http_error.ACADEMY_NOT_FOUND
	}

	m, err := s.repo.GetMaterialWithProgress(ctx, accountId, academy.Id, materialSlug)
	if err != nil {
		return m, err
	}

	return s.sanitizeMaterialProgress(m), nil
}

func (s *academyService) CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (entity.AcademyMaterial, error) {
	if req.AcademyId == uuid.Nil {
		return entity.AcademyMaterial{}, http_error.ACADEMY_ID_REQUIRED
	}
	if _, err := s.repo.GetAcademyByID(ctx, req.AcademyId); err != nil {
		return entity.AcademyMaterial{}, http_error.ACADEMY_NOT_FOUND
	}
	slugVal := req.Slug
	if slugVal == "" {
		slugVal = slug.Make(req.Title)
	}
	orderCount, _ := s.repo.CountMaterialsByAcademyID(ctx, req.AcademyId)
	m := entity.AcademyMaterial{
		Id:            uuid.New(),
		AcademyId:     req.AcademyId,
		Title:         req.Title,
		Slug:          slugVal,
		Description:   req.Description,
		Order:         uint(orderCount + 1),
		ContentsCount: 0,
	}
	var createdMaterial entity.AcademyMaterial
	err := s.repo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
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
		return nil
	})

	if err != nil {
		return entity.AcademyMaterial{}, err
	}

	return s.sanitizeMaterialProgress(createdMaterial), nil
}

func (s *academyService) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	m, err := s.repo.GetMaterialByID(ctx, id)
	if err != nil {
		return http_error.MATERIAL_NOT_FOUND
	}

	// 1. Cascade delete progress
	if err := s.repo.DeleteContentProgressByMaterialID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteMaterialProgressByMaterialID(ctx, id); err != nil {
		return err
	}

	// 2. Delete the Material
	if err := s.repo.DeleteMaterial(ctx, id); err != nil {
		return err
	}

	// 3. Reorder Logic
	if err := s.repo.DecrementMaterialOrdersGreaterThan(ctx, m.AcademyId, m.Order); err != nil {
		return err
	}

	// 4. Recalculate Counter in Academy (Parent)
	realCount, err := s.repo.CountMaterialsByAcademyID(ctx, m.AcademyId)
	if err != nil {
		return err
	}

	academy, err := s.repo.GetAcademyByID(ctx, m.AcademyId)
	if err == nil {
		academy.MaterialsCount = int64(realCount)
		s.repo.UpdateAcademy(ctx, academy)
	}

	return nil
}

// ================= CONTENT =================

func (s *academyService) CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error) {
	if req.MaterialId == uuid.Nil {
		return entity.AcademyContent{}, http_error.MATERIAL_ID_REQUIRED
	}
	m, err := s.repo.GetMaterialByID(ctx, req.MaterialId)
	if err != nil {
		return entity.AcademyContent{}, http_error.MATERIAL_NOT_FOUND
	}

	var createdContent entity.AcademyContent

	err = s.repo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		count, _ := txRepo.CountContentsByMaterialID(ctx, req.MaterialId)
		order := uint(count + 1)

		c := entity.AcademyContent{
			Id:         uuid.New(),
			MaterialId: req.MaterialId,
			Title:      req.Title,
			Contents:   req.Contents,
			Order:      order,
			AcademyContentProgress: entity.AcademyContentProgress{
				Status: entity.StatusNotStarted,
			},
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

		if err := txRepo.BatchResetMaterialProgressStatus(ctx, req.MaterialId); err != nil {
			return err
		}

		if err := txRepo.BatchResetAcademyProgressStatus(ctx, m.AcademyId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.AcademyContent{}, err
	}

	return createdContent, nil
}

func (s *academyService) GetContent(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContent, error) {
	material, err := s.GetMaterial(ctx, accountId, academySlug, materialSlug)
	if err != nil {
		return entity.AcademyContent{}, err
	}
	return s.repo.GetContentWithProgress(ctx, accountId, material.AcademyId, material.Id, order)
}

func (s *academyService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetContentByID(ctx, id)
	if err != nil {
		return http_error.CONTENT_NOT_FOUND
	}

	if err := s.repo.DeleteContentProgressByContentID(ctx, id); err != nil {
		return err
	}

	if err := s.repo.DeleteContent(ctx, id); err != nil {
		return err
	}

	if err := s.repo.DecrementContentOrdersGreaterThan(ctx, c.MaterialId, c.Order); err != nil {
		return err
	}

	realCount, err := s.repo.CountContentsByMaterialID(ctx, c.MaterialId)
	if err != nil {
		return err
	}

	material, err := s.repo.GetMaterialByID(ctx, c.MaterialId)
	if err == nil {
		material.ContentsCount = realCount
		s.repo.UpdateMaterial(ctx, material)
	}

	return nil
}

// ================= PROGRESS (TRANSACTIONAL) =================

func (s *academyService) UpdateContentProgress(ctx context.Context, accountId uuid.UUID, academySlug string, materialSlug string, order uint) (entity.AcademyContentProgress, entity.AcademyMaterialProgress, entity.AcademyProgress, error) {
	academy, err := s.repo.GetAcademyBySlug(ctx, academySlug)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, http_error.ACADEMY_NOT_FOUND
	}
	material, err := s.repo.GetMaterialBySlug(ctx, academy.Id, materialSlug)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, http_error.MATERIAL_NOT_FOUND
	}
	content, err := s.repo.GetContentBySlug(ctx, material.Id, order)
	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, http_error.CONTENT_NOT_FOUND
	}

	var acp entity.AcademyContentProgress
	var amp entity.AcademyMaterialProgress
	var ap entity.AcademyProgress

	err = s.repo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {

		// --- STEP A: UPDATE CONTENT PROGRESS ---
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

		// --- STEP B: UPDATE MATERIAL PROGRESS ---
		totalContents, err := txRepo.CountCompletedContentsByMaterialAndAccount(ctx, accountId, material.Id)
		if err != nil {
			return err
		}

		m, _ := txRepo.GetMaterialByID(ctx, material.Id)

		matStatus := entity.StatusInProgress
		var matCompletedAt *time.Time
		progressPct := 0.0

		if m.ContentsCount > 0 {
			progressPct = (float64(totalContents) / float64(m.ContentsCount)) * 100
			progressPct = math.Round(progressPct*100) / 100

			if totalContents >= m.ContentsCount {
				matStatus = entity.StatusCompleted
				matCompletedAt = utils.Ptr(time.Now())
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
			TotalCompletedContents: uint(totalContents),
			Status:                 matStatus,
			CompletedAt:            matCompletedAt,
		}
		if _, err := txRepo.UpsertMaterialProgress(ctx, amp); err != nil {
			return err
		}

		// --- STEP C: UPDATE ACADEMY PROGRESS ---
		totalMaterials, err := txRepo.CountCompletedMaterialsByAcademyAndAccount(ctx, accountId, academy.Id)
		if err != nil {
			return err
		}

		accumulatedProgress, _ := txRepo.GetAccumulatedMaterialProgress(ctx, accountId, academy.Id)

		acadStatus := entity.StatusNotStarted
		var acadCompletedAt *time.Time
		acadProgressPct := 0.0

		if academy.MaterialsCount > 0 {
			acadProgressPct = accumulatedProgress / float64(academy.MaterialsCount)
			acadProgressPct = math.Round(acadProgressPct*100) / 100

			if totalMaterials >= int64(academy.MaterialsCount) {
				acadStatus = entity.StatusCompleted
				acadCompletedAt = utils.Ptr(time.Now())
			} else if accumulatedProgress > 0 {
				acadStatus = entity.StatusInProgress
			}
		} else {
			acadStatus = entity.StatusCompleted
			acadProgressPct = 100
		}

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
			TotalCompletedMaterials: uint(totalMaterials),
			Status:                  acadStatus,
			CompletedAt:             acadCompletedAt,
		}
		if _, err := txRepo.UpsertAcademyProgress(ctx, ap); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.AcademyContentProgress{}, entity.AcademyMaterialProgress{}, entity.AcademyProgress{}, err
	}

	return acp, amp, ap, nil
}

func (s *academyService) UpdateAcademyProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy) (entity.AcademyProgress, error) {
	totalCompleted, _ := s.repo.CountCompletedMaterialsByAcademyAndAccount(ctx, accountId, academy.Id)
	accumulatedProgress, _ := s.repo.GetAccumulatedMaterialProgress(ctx, accountId, academy.Id)

	status := entity.StatusNotStarted
	var completedAt *time.Time
	pct := 0.0

	if academy.MaterialsCount > 0 {
		pct = accumulatedProgress / float64(academy.MaterialsCount)
		pct = math.Round(pct*100) / 100

		if totalCompleted >= int64(academy.MaterialsCount) {
			status = entity.StatusCompleted
			completedAt = utils.Ptr(time.Now())
		} else if accumulatedProgress > 0 {
			status = entity.StatusInProgress
		}
	}

	existingAP, _ := s.repo.GetAcademyProgress(ctx, accountId, academy.Id)
	apID := existingAP.Id
	if apID == uuid.Nil {
		apID = uuid.New()
	}

	ap := entity.AcademyProgress{
		Id:                      apID,
		AccountId:               accountId,
		AcademyId:               academy.Id,
		Progress:                pct,
		TotalCompletedMaterials: uint(totalCompleted),
		Status:                  status,
		CompletedAt:             completedAt,
	}
	_, err := s.repo.UpsertAcademyProgress(ctx, ap)
	return ap, err
}

func (s *academyService) UpdateMaterialProgress(ctx context.Context, accountId uuid.UUID, academy entity.Academy, material entity.AcademyMaterial) (entity.AcademyMaterialProgress, error) {
	total, _ := s.repo.CountCompletedContentsByMaterialAndAccount(ctx, accountId, material.Id)
	m, _ := s.repo.GetMaterialByID(ctx, material.Id)

	status := entity.StatusInProgress
	var completedAt *time.Time
	pct := 0.0
	if m.ContentsCount > 0 {
		pct = (float64(total) / float64(m.ContentsCount)) * 100
		pct = math.Round(pct*100) / 100

		if total >= m.ContentsCount {
			status = entity.StatusCompleted
			completedAt = utils.Ptr(time.Now())
		}
	}

	existingAMP, _ := s.repo.GetMaterialProgress(ctx, accountId, academy.Id, material.Id)
	ampID := existingAMP.Id
	if ampID == uuid.Nil {
		ampID = uuid.New()
	}

	amp := entity.AcademyMaterialProgress{
		Id:                     ampID,
		AccountId:              accountId,
		AcademyId:              academy.Id,
		MaterialId:             material.Id,
		Progress:               pct,
		TotalCompletedContents: uint(total),
		Status:                 status,
		CompletedAt:            completedAt,
	}
	_, err := s.repo.UpsertMaterialProgress(ctx, amp)
	return amp, err
}

// HELPER

func (s *academyService) sanitizeAcademyProgress(a entity.Academy, startedMaterials int64, accumulatedProgress float64) entity.Academy {
	if a.MaterialsCount > 0 {
		completed := a.AcademyProgresss.TotalCompletedMaterials

		realProgress := accumulatedProgress / float64(a.MaterialsCount)
		realProgress = math.Round(realProgress*100) / 100

		a.AcademyProgresss.Progress = realProgress

		// Logika Status
		if completed == 0 {
			if startedMaterials > 0 {
				a.AcademyProgresss.Status = entity.StatusInProgress
			} else {
				a.AcademyProgresss.Status = entity.StatusNotStarted
			}
			a.AcademyProgresss.CompletedAt = nil
		} else if int64(completed) < int64(a.MaterialsCount) {
			a.AcademyProgresss.Status = entity.StatusInProgress
			a.AcademyProgresss.CompletedAt = nil
		} else if int64(completed) >= int64(a.MaterialsCount) {
			a.AcademyProgresss.Status = entity.StatusCompleted
			a.AcademyProgresss.Progress = 100
		}
	}
	return a
}

func (s *academyService) sanitizeMaterialProgress(m entity.AcademyMaterial) entity.AcademyMaterial {
	if m.ContentsCount > 0 {
		completed := m.AcademyMaterialProgress.TotalCompletedContents

		realProgress := (float64(completed) / float64(m.ContentsCount)) * 100
		realProgress = math.Round(realProgress*100) / 100

		m.AcademyMaterialProgress.Progress = realProgress

		if completed == 0 {
			m.AcademyMaterialProgress.Status = entity.StatusNotStarted
			m.AcademyMaterialProgress.CompletedAt = nil
		} else if int64(completed) < m.ContentsCount {
			m.AcademyMaterialProgress.Status = entity.StatusInProgress
			m.AcademyMaterialProgress.CompletedAt = nil
		} else if int64(completed) >= m.ContentsCount {
			m.AcademyMaterialProgress.Status = entity.StatusCompleted
			m.AcademyMaterialProgress.Progress = 100
		}
	}
	return m
}