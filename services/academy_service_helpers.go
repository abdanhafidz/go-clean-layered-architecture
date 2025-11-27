package services

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils"
)

func (s *academyService) getOrCreateID(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.New()
	}
	return id
}

func (s *academyService) calculateProgress(completed, total int64) float64 {
	if total <= 0 {
		return 0
	}
	progress := (float64(completed) / float64(total)) * 100
	return math.Round(progress*100) / 100
}

func (s *academyService) getProgressStatus(progress float64, completed, total int64) string {
	if progress >= 100 || (total > 0 && completed >= total) {
		return entity.StatusCompleted
	}
	if progress > 0 {
		return entity.StatusInProgress
	}
	return entity.StatusNotStarted
}

func (s *academyService) upsertContentProgressSimplified(
	ctx context.Context,
	txRepo repositories.AcademyRepository,
	accountId, academyId, materialId, contentId uuid.UUID,
) entity.AcademyContentProgress {
	existing, _ := txRepo.GetContentProgress(ctx, accountId, academyId, materialId, contentId)

	acp := entity.AcademyContentProgress{
		Id:          s.getOrCreateID(existing.Id),
		AccountId:   accountId,
		AcademyId:   academyId,
		MaterialId:  materialId,
		ContentId:   contentId,
		Status:      entity.StatusCompleted,
		CompletedAt: utils.Ptr(time.Now()),
	}

	txRepo.UpsertContentProgress(ctx, acp)
	return acp
}

func (s *academyService) calculateMaterialProgress(
	ctx context.Context,
	txRepo repositories.AcademyRepository,
	accountId, academyId, materialId uuid.UUID,
	totalCompleted, totalContents int64,
) entity.AcademyMaterialProgress {
	existing, _ := txRepo.GetMaterialProgress(ctx, accountId, academyId, materialId)

	progress := s.calculateProgress(totalCompleted, totalContents)
	status := s.getProgressStatus(progress, totalCompleted, totalContents)

	var completedAt *time.Time
	if status == entity.StatusCompleted {
		completedAt = utils.Ptr(time.Now())
		progress = 100
	}

	return entity.AcademyMaterialProgress{
		Id:                     s.getOrCreateID(existing.Id),
		AccountId:              accountId,
		AcademyId:              academyId,
		MaterialId:             materialId,
		Progress:               progress,
		TotalCompletedContents: uint(totalCompleted),
		Status:                 status,
		CompletedAt:            completedAt,
	}
}


func (s *academyService) calculateAcademyProgress(
	ctx context.Context,
	txRepo repositories.AcademyRepository,
	accountId, academyId uuid.UUID,
	accumulatedProgress float64,
	totalMaterials int64,
) entity.AcademyProgress {
	existing, _ := txRepo.GetAcademyProgress(ctx, accountId, academyId)

	var progress float64
	if totalMaterials > 0 {
		progress = math.Round((accumulatedProgress/float64(totalMaterials))*100) / 100
	}

	status := s.getProgressStatus(progress, 0, totalMaterials)

	var completedAt *time.Time
	if status == entity.StatusCompleted {
		completedAt = utils.Ptr(time.Now())
		progress = 100
	}

	totalCompleted, _ := txRepo.CountCompletedMaterialsByAcademyAndAccount(ctx, accountId, academyId)

	return entity.AcademyProgress{
		Id:                      s.getOrCreateID(existing.Id),
		AccountId:               accountId,
		AcademyId:               academyId,
		Progress:                progress,
		TotalCompletedMaterials: uint(totalCompleted),
		Status:                  status,
		CompletedAt:             completedAt,
	}
}

func (s *academyService) updateAcademyMaterialCount(ctx context.Context, txRepo repositories.AcademyRepository, academyId uuid.UUID) error {
	count, err := txRepo.CountMaterialsByAcademyID(ctx, academyId)
	if err != nil {
		return err
	}
	academy, _ := txRepo.GetAcademyByID(ctx, academyId)
	academy.MaterialsCount = count
	_, err = txRepo.UpdateAcademy(ctx, academy)
	return err
}


func (s *academyService) updateMaterialContentCount(ctx context.Context, txRepo repositories.AcademyRepository, materialId uuid.UUID) error {
	count, err := txRepo.CountContentsByMaterialID(ctx, materialId)
	if err != nil {
		return err
	}
	material, _ := txRepo.GetMaterialByID(ctx, materialId)
	material.ContentsCount = count
	_, err = txRepo.UpdateMaterial(ctx, material)
	return err
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format("2006-01-02T15:04:05Z07:00")
	return &formatted
}


func buildAcademyProgressResponse(ap entity.AcademyProgress) *dto.AcademyProgressResponse {
	if ap.Id == uuid.Nil {
		return nil
	}
	return &dto.AcademyProgressResponse{
		Id:                      ap.Id,
		AccountId:               ap.AccountId,
		AcademyId:               ap.AcademyId,
		Status:                  ap.Status,
		Progress:                ap.Progress,
		TotalCompletedMaterials: ap.TotalCompletedMaterials,
		CompletedAt:             formatTime(ap.CompletedAt),
	}
}

func buildAcademyContentResponse(content entity.AcademyContent, progress entity.AcademyContentProgress) dto.AcademyContentResponse {
	status := entity.StatusNotStarted
	if progress.Id != uuid.Nil {
		status = progress.Status
	}
	return dto.AcademyContentResponse{
		Id:     content.Id,
		Order:  content.Order,
		Title:  content.Title,
		Status: status,
	}
}


func buildAcademyMaterialResponse(
	material entity.AcademyMaterial,
	contents []entity.AcademyContent,
	progressMap map[uuid.UUID]entity.AcademyContentProgress,
	progress entity.AcademyMaterialProgress,
) dto.AcademyMaterialResponse {
	status := entity.StatusNotStarted
	if progress.Id != uuid.Nil {
		status = progress.Status
	}

	contentDTOs := make([]dto.AcademyContentResponse, 0)
	for _, content := range contents {
		contentProgress := progressMap[content.Id]
		contentDTOs = append(contentDTOs, buildAcademyContentResponse(content, contentProgress))
	}

	return dto.AcademyMaterialResponse{
		Id:                     material.Id,
		Order:                  material.Order,
		Title:                  material.Title,
		Slug:                   material.Slug,
		Status:                 status,
		Progress:               progress.Progress,
		TotalCompletedContents: progress.TotalCompletedContents,
		ContentsCount:          material.ContentsCount,
		Contents:               contentDTOs,
	}
}

func buildAcademyDetailResponse(
	academy entity.Academy,
	materials []entity.AcademyMaterial,
	academyProgress entity.AcademyProgress,
	ctx context.Context,
	repo repositories.AcademyRepository,
) *dto.AcademyDetailResponse {
	materialDTOs := make([]dto.AcademyMaterialResponse, 0)

	materialIds := make([]uuid.UUID, 0)
	contentIds := make([]uuid.UUID, 0)

	for _, material := range materials {
		materialIds = append(materialIds, material.Id)
		for _, content := range material.Contents {
			contentIds = append(contentIds, content.Id)
		}
	}

	materialProgressMap, _ := repo.GetMaterialProgressBatch(ctx, academyProgress.AccountId, academyProgress.AcademyId, materialIds)
	contentProgressMap, _ := repo.GetContentProgressBatch(ctx, academyProgress.AccountId, academyProgress.AcademyId, contentIds)

	for _, material := range materials {
		materialProgress := materialProgressMap[material.Id]
		contentDTOs := make([]dto.AcademyContentResponse, 0)

		for _, content := range material.Contents {
			contentProgress := contentProgressMap[content.Id]
			contentDTOs = append(contentDTOs, buildAcademyContentResponse(content, contentProgress))
		}

		materialDTOs = append(materialDTOs, dto.AcademyMaterialResponse{
			Id:                     material.Id,
			Order:                  material.Order,
			Title:                  material.Title,
			Slug:                   material.Slug,
			Status:                 getContentStatus(materialProgress),
			Progress:               materialProgress.Progress,
			TotalCompletedContents: materialProgress.TotalCompletedContents,
			ContentsCount:          material.ContentsCount,
			Contents:               contentDTOs,
		})
	}

	return &dto.AcademyDetailResponse{
		Id:             academy.Id,
		Title:          academy.Title,
		Slug:           academy.Slug,
		Description:    academy.Description,
		ImageUrl:       academy.ImageUrl,
		MaterialsCount: academy.MaterialsCount,
		UserProgress:   buildAcademyProgressResponse(academyProgress),
		Materials:      materialDTOs,
	}
}

func getContentStatus(progress entity.AcademyMaterialProgress) string {
	if progress.Id != uuid.Nil {
		return progress.Status
	}
	return entity.StatusNotStarted
}

func buildMaterialProgressResponse(mp entity.AcademyMaterialProgress) *dto.MaterialProgressResponse {
	if mp.Id == uuid.Nil {
		return nil
	}
	return &dto.MaterialProgressResponse{
		Id:                     mp.Id,
		AccountId:              mp.AccountId,
		AcademyId:              mp.AcademyId,
		MaterialId:             mp.MaterialId,
		Progress:               mp.Progress,
		TotalCompletedContents: mp.TotalCompletedContents,
		Status:                 mp.Status,
		CompletedAt:            formatTime(mp.CompletedAt),
	}
}

func buildContentDetailResponse(content entity.AcademyContent, progress entity.AcademyContentProgress) dto.ContentDetailResponse {
	status := entity.StatusNotStarted
	if progress.Id != uuid.Nil {
		status = progress.Status
	}

	return dto.ContentDetailResponse{
		Id:     content.Id,
		Order:  content.Order,
		Title:  content.Title,
		Status: status,
	}
}


func buildMaterialDetailResponse(
	material entity.AcademyMaterial,
	contents []entity.AcademyContent,
	materialProgress entity.AcademyMaterialProgress,
	accountId, academyId uuid.UUID,
	ctx context.Context,
	repo repositories.AcademyRepository,
	academySlug, materialSlug string,
) *dto.MaterialDetailResponse {
	contentDTOs := make([]dto.ContentDetailResponse, 0)
	contentIds := make([]uuid.UUID, len(contents))
	for i, c := range contents {
		contentIds[i] = c.Id
	}

	contentProgressMap, _ := repo.GetContentProgressBatch(ctx, accountId, academyId, contentIds)

	for _, content := range contents {
		contentProgress := contentProgressMap[content.Id]
		contentDTOs = append(contentDTOs, buildContentDetailResponse(content, contentProgress))
	}

	return &dto.MaterialDetailResponse{
		Id:            material.Id,
		AcademyId:     material.AcademyId,
		Title:         material.Title,
		Slug:          material.Slug,
		Description:   material.Description,
		Order:         material.Order,
		ContentsCount: material.ContentsCount,
		Progress:      buildMaterialProgressResponse(materialProgress),
		Contents:      contentDTOs,
		Meta: map[string]string{
			"academy_slug":  academySlug,
			"material_slug": materialSlug,
		},
	}
}
