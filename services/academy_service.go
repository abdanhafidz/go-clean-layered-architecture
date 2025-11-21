package services

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils" 
)

type AcademyRepositoryExtensions interface {
	BatchRecalculateAcademyProgress(ctx context.Context, academyId uuid.UUID) error
	BatchRecalculateMaterialProgress(ctx context.Context, materialId uuid.UUID) error
}

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
}

type academyService struct {
	academyRepo repositories.AcademyRepository
}

func NewAcademyService(academyRepo repositories.AcademyRepository) AcademyService {
	return &academyService{academyRepo: academyRepo}
}

// ================= ACADEMY =================

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
		Description:    req.Description,
		ImageUrl:       req.ImageUrl,
		MaterialsCount: 0,
	}
	return s.academyRepo.CreateAcademy(ctx, a)
}

func (s *academyService) UpdateAcademy(ctx context.Context, id uuid.UUID, req dto.UpdateAcademyRequest) (entity.Academy, error) {
	existing, err := s.academyRepo.GetAcademyByID(ctx, id)
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
	return s.academyRepo.UpdateAcademy(ctx, existing)
}

func (s *academyService) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	_, mats, err := s.academyRepo.GetAcademyWithMaterials(ctx, id)
	if err != nil {
		return http_error.ACADEMY_NOT_FOUND
	}
	if len(mats) > 0 {
		return http_error.ACADEMY_HAS_MATERIALS
	}
	return s.academyRepo.DeleteAcademy(ctx, id)
}

func (s *academyService) ListAcademies(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error) {
	// Logic list tetap sama, namun karena progress di DB sudah konsisten (akibat recalculate),
	// kita tidak perlu sanitasi manual yang berat di sini.
	return s.academyRepo.ListAcademy(ctx, accountId)
}

// ================= MATERIAL (CRITICAL LOGIC HERE) =================

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
		return entity.AcademyMaterial{}, http_error.ACADEMY_ID_REQUIRED
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

		// 1. Create Material
		res, err := txRepo.CreateMaterial(ctx, m)
		if err != nil {
			return err
		}
		createdMaterial = res

		// 2. Update Parent Count (Academy)
		realCount, err := txRepo.CountMaterialsByAcademyID(ctx, req.AcademyId)
		if err != nil {
			return err
		}
		academy, _ := txRepo.GetAcademyByID(ctx, req.AcademyId)
		academy.MaterialsCount = int64(realCount)
		if _, err := txRepo.UpdateAcademy(ctx, academy); err != nil {
			return err
		}

		// 3. VALIDASI KRUSIAL: Recalculate Progress User
		// Karena MaterialsCount bertambah (misal 3 -> 4), user yang sebelumnya 100% (3/3)
		// sekarang menjadi 75% (3/4). Status harus berubah jadi 'InProgress'.
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
		// 1. Cleanup Child Progress & Content
		if err := txRepo.DeleteContentProgressByMaterialID(ctx, id); err != nil { return err }
		if err := txRepo.DeleteMaterialProgressByMaterialID(ctx, id); err != nil { return err }
		
		// 2. Delete Material
		if err := txRepo.DeleteMaterial(ctx, id); err != nil { return err }

		// 3. Reorder
		if err := txRepo.DecrementMaterialOrdersGreaterThan(ctx, m.AcademyId, m.Order); err != nil { return err }

		// 4. Update Parent Count
		realCount, err := txRepo.CountMaterialsByAcademyID(ctx, m.AcademyId)
		if err != nil { return err }
		
		academy, _ := txRepo.GetAcademyByID(ctx, m.AcademyId)
		academy.MaterialsCount = int64(realCount)
		if _, err := txRepo.UpdateAcademy(ctx, academy); err != nil { return err }

		// 5. VALIDASI KRUSIAL: Recalculate Progress
		// Karena MaterialsCount berkurang (misal 4 -> 3), user yang sebelumnya 75% (3/4)
		// sekarang menjadi 100% (3/3). Status harus berubah jadi 'Completed'.
		if err := txRepo.BatchRecalculateAcademyProgress(ctx, m.AcademyId); err != nil { return err }

		return nil
	})
}

// ================= CONTENT (CRITICAL LOGIC HERE) =================

func (s *academyService) CreateContent(ctx context.Context, req dto.CreateContentRequest) (entity.AcademyContent, error) {
	if req.MaterialId == uuid.Nil {
		return entity.AcademyContent{}, http_error.MATERIAL_ID_REQUIRED
	}

	var createdContent entity.AcademyContent

	err := s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
		m, err := txRepo.GetMaterialByID(ctx, req.MaterialId)
		if err != nil { return http_error.MATERIAL_NOT_FOUND }

		count, _ := txRepo.CountContentsByMaterialID(ctx, req.MaterialId)
		
		c := entity.AcademyContent{
			Id:         uuid.New(),
			MaterialId: req.MaterialId,
			Title:      req.Title,
			Contents:   req.Contents,
			Order:      uint(count + 1),
			// Progress status untuk content baru default-nya NotStarted
		}

		// 1. Create Content
		res, err := txRepo.CreateContent(ctx, c)
		if err != nil { return err }
		createdContent = res

		// 2. Update Parent Count (Material)
		realCount, err := txRepo.CountContentsByMaterialID(ctx, req.MaterialId)
		if err != nil { return err }
		
		m.ContentsCount = realCount
		if _, err := txRepo.UpdateMaterial(ctx, m); err != nil { return err }

		// 3. VALIDASI KRUSIAL: Recalculate Progress Material
		// Konten bertambah -> Material user yang 'Completed' harus jadi 'InProgress'.
		// Ini juga akan men-trigger efek domino ke Academy Progress (accumulated value berubah).
		if err := txRepo.BatchRecalculateMaterialProgress(ctx, req.MaterialId); err != nil { return err }
		
		// Update Academy juga karena bobot material berubah
		if err := txRepo.BatchRecalculateAcademyProgress(ctx, m.AcademyId); err != nil { return err }

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
		// 1. Delete Progress & Content
		if err := txRepo.DeleteContentProgressByContentID(ctx, id); err != nil { return err }
		if err := txRepo.DeleteContent(ctx, id); err != nil { return err }

		// 2. Reorder
		if err := txRepo.DecrementContentOrdersGreaterThan(ctx, c.MaterialId, c.Order); err != nil { return err }

		// 3. Update Parent Count (Material)
		realCount, err := txRepo.CountContentsByMaterialID(ctx, c.MaterialId)
		if err != nil { return err }

		material, _ := txRepo.GetMaterialByID(ctx, c.MaterialId)
		material.ContentsCount = realCount
		if _, err := txRepo.UpdateMaterial(ctx, material); err != nil { return err }

		// 4. VALIDASI KRUSIAL: Recalculate
		// Konten berkurang -> Material user 'InProgress' bisa jadi 'Completed'.
		if err := txRepo.BatchRecalculateMaterialProgress(ctx, c.MaterialId); err != nil { return err }
		
		// Update Academy juga
		if err := txRepo.BatchRecalculateAcademyProgress(ctx, material.AcademyId); err != nil { return err }

		return nil
	})
}

// ================= UPDATE PROGRESS (USER INTERACTION) =================
// Logic ini dijalankan saat user menyelesaikan satu konten.

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

		// 1. Tandai Konten Selesai
		existingACP, _ := txRepo.GetContentProgress(ctx, accountId, academy.Id, material.Id, content.Id)
		acpID := existingACP.Id
		if acpID == uuid.Nil { acpID = uuid.New() }

		acp = entity.AcademyContentProgress{
			Id:          acpID,
			AccountId:   accountId,
			AcademyId:   academy.Id,
			MaterialId:  material.Id,
			ContentId:   content.Id,
			Status:      entity.StatusCompleted,
			CompletedAt: utils.Ptr(time.Now()),
		}
		if _, err := txRepo.UpsertContentProgress(ctx, acp); err != nil { return err }

		totalContentsCompleted, _ := txRepo.CountCompletedContentsByMaterialAndAccount(ctx, accountId, material.Id)
		m, _ := txRepo.GetMaterialByID(ctx, material.Id) // Ambil count terbaru
		
		matStatus := entity.StatusInProgress
		var matCompletedAt *time.Time
		progressPct := 0.0

		if m.ContentsCount > 0 {
			progressPct = (float64(totalContentsCompleted) / float64(m.ContentsCount)) * 100
			progressPct = math.Round(progressPct*100) / 100
			if totalContentsCompleted >= m.ContentsCount {
				matStatus = entity.StatusCompleted
				matCompletedAt = utils.Ptr(time.Now())
				progressPct = 100 // Force 100
			}
		} else {
			// Edge case: Material tanpa konten dianggap completed
			matStatus = entity.StatusCompleted
			progressPct = 100
		}

		existingAMP, _ := txRepo.GetMaterialProgress(ctx, accountId, academy.Id, material.Id)
		ampID := existingAMP.Id
		if ampID == uuid.Nil { ampID = uuid.New() }

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
		if _, err := txRepo.UpsertMaterialProgress(ctx, amp); err != nil { return err }

		// 3. Hitung Ulang Progress Academy (Aggregation)
		// Logic: Average dari seluruh Material Progress
		// Atau Logic Sederhana: (Completed Material / Total Material) -> Ini kurang akurat jika material bobotnya sama.
		// Logic Lebih Akurat: Sum(MaterialProgress) / TotalMaterial
		
		accumulatedProgress, _ := txRepo.GetAccumulatedMaterialProgress(ctx, accountId, academy.Id)
		a, _ := txRepo.GetAcademyByID(ctx, academy.Id) // Ambil count terbaru

		acadStatus := entity.StatusNotStarted
		var acadCompletedAt *time.Time
		acadProgressPct := 0.0

		if a.MaterialsCount > 0 {
			// Rumus: Total Akumulasi Persen Material / Jumlah Material
			// Contoh: Mat A (100%) + Mat B (50%) = 150%. Total Mat = 2. Academy Progress = 75%.
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

		// Hitung berapa material yang full completed (opsional, untuk display)
		totalMaterialsCompleted, _ := txRepo.CountCompletedMaterialsByAcademyAndAccount(ctx, accountId, academy.Id)

		existingAP, _ := txRepo.GetAcademyProgress(ctx, accountId, academy.Id)
		apID := existingAP.Id
		if apID == uuid.Nil { apID = uuid.New() }

		ap = entity.AcademyProgress{
			Id:                      apID,
			AccountId:               accountId,
			AcademyId:               academy.Id,
			Progress:                acadProgressPct,
			TotalCompletedMaterials: uint(totalMaterialsCompleted),
			Status:                  acadStatus,
			CompletedAt:             acadCompletedAt,
		}
		if _, err := txRepo.UpsertAcademyProgress(ctx, ap); err != nil { return err }

		return nil
	})

	return acp, amp, ap, err
}