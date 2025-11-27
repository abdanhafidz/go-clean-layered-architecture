package repositories

import (
	"context"
	"errors"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AcademyRepository interface {
	Atomic(ctx context.Context, fn func(r AcademyRepository) error) error

	GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error)
	GetAcademyWithProgress(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error)
	CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error
	ListAcademy(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error)
	GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error)
	CountMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) (int64, error)

	GetMaterialBySlug(ctx context.Context, academy_id uuid.UUID, materialSlug string) (entity.AcademyMaterial, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error)
	GetMaterialWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, slug string) (entity.AcademyMaterial, error)
	CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	DeleteMaterial(ctx context.Context, id uuid.UUID) error
	ListMaterials(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error)
	GetMaterialWithContents(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, []entity.AcademyContent, error)
	GetMaterialsWithContents(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error)

	GetContentBySlug(ctx context.Context, materialId uuid.UUID, order uint) (entity.AcademyContent, error)
	GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error)
	GetContentWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID, order uint) (entity.AcademyContent, error)
	CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error
	CountContentsByMaterialID(ctx context.Context, materialId uuid.UUID) (int64, error)

	GetAcademyProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyProgress, error)
	GetMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error)
	GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID, contentId uuid.UUID) (entity.AcademyContentProgress, error)
	UpsertAcademyProgress(ctx context.Context, p entity.AcademyProgress) (entity.AcademyProgress, error)
	UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error)
	UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error)
	DeleteContentProgressByContentID(ctx context.Context, contentId uuid.UUID) error
	DeleteContentProgressByMaterialID(ctx context.Context, materialId uuid.UUID) error
	DeleteMaterialProgressByMaterialID(ctx context.Context, materialId uuid.UUID) error

	CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error)
	CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error)
	CountStartedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error)
	DecrementMaterialOrdersGreaterThan(ctx context.Context, academyId uuid.UUID, order uint) error
	DecrementContentOrdersGreaterThan(ctx context.Context, materialId uuid.UUID, order uint) error
	GetAccumulatedMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (float64, error)

	BatchRecalculateMaterialProgress(ctx context.Context, materialId uuid.UUID) error
	BatchRecalculateAcademyProgress(ctx context.Context, academyId uuid.UUID) error

	GetMaterialProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialIds []uuid.UUID) (map[uuid.UUID]entity.AcademyMaterialProgress, error)
	GetContentProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, contentIds []uuid.UUID) (map[uuid.UUID]entity.AcademyContentProgress, error)
}

type academyRepository struct{ db *gorm.DB }

func NewAcademyRepository(db *gorm.DB) AcademyRepository {
	return &academyRepository{db: db}
}

func (r *academyRepository) Atomic(ctx context.Context, fn func(r AcademyRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewAcademyRepository(tx)
		return fn(txRepo)
	})
}

func (r *academyRepository) GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error) {
	var a entity.Academy
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return entity.Academy{}, nil, err
	}
	var m []entity.AcademyMaterial
	err := r.db.WithContext(ctx).Where("academy_id = ?", id).Order("\"order\" ASC").Find(&m).Error
	return a, m, err
}

func (r *academyRepository) GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	var a entity.Academy
	return a, r.db.WithContext(ctx).First(&a, "id = ?", id).Error
}

func (r *academyRepository) GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error) {
	var a entity.Academy
	return a, r.db.WithContext(ctx).First(&a, "slug = ?", slug).Error
}

func (r *academyRepository) GetAcademyWithProgress(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error) {
	a, err := r.GetAcademyBySlug(ctx, slug)
	if err != nil {
		return a, err
	}
	ap, err := r.GetAcademyProgress(ctx, accountId, a.Id)
	if err != nil {
		return a, err
	}
	a.AcademyProgresss = ap
	return a, nil
}

func (r *academyRepository) CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error) {
	return a, r.db.WithContext(ctx).Create(&a).Error
}

func (r *academyRepository) UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error) {
	return a, r.db.WithContext(ctx).Save(&a).Error
}

func (r *academyRepository) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Academy{}, "id = ?", id).Error
}

func (r *academyRepository) ListAcademy(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error) {
	var list []entity.Academy
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return list, nil
	}

	academyIDs := make([]uuid.UUID, len(list))
	for i, ac := range list {
		academyIDs[i] = ac.Id
	}

	var progressList []entity.AcademyProgress
	if err := r.db.WithContext(ctx).
		Where("account_id = ?", accountId).
		Where("academy_id IN ?", academyIDs).
		Find(&progressList).Error; err != nil {
		return nil, err
	}

	progressMap := make(map[uuid.UUID]entity.AcademyProgress)
	for _, p := range progressList {
		progressMap[p.AcademyId] = p
	}

	for i := range list {
		if p, exists := progressMap[list[i].Id]; exists {
			list[i].AcademyProgresss = p
		} else {
			list[i].AcademyProgresss = entity.AcademyProgress{
				AccountId: accountId,
				AcademyId: list[i].Id,
				Status:    entity.StatusNotStarted,
			}
		}
	}
	return list, nil
}

func (r *academyRepository) CountMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyMaterial{}).Where("academy_id = ?", academyId).Count(&count).Error
	return count, err
}

func (r *academyRepository) GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	return m, r.db.WithContext(ctx).First(&m, "id = ?", id).Error
}

func (r *academyRepository) GetMaterialBySlug(ctx context.Context, academy_id uuid.UUID, materialSlug string) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	return m, r.db.WithContext(ctx).First(&m, "academy_id = ? AND slug = ?", academy_id, materialSlug).Error
}

func (r *academyRepository) GetMaterialWithContents(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, []entity.AcademyContent, error) {
	var m entity.AcademyMaterial
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return entity.AcademyMaterial{}, nil, err
	}
	var c []entity.AcademyContent
	err := r.db.WithContext(ctx).Where("material_id = ?", id).Order("\"order\" ASC").Find(&c).Error
	return m, c, err
}

func (r *academyRepository) GetMaterialWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, slug string) (entity.AcademyMaterial, error) {
	m, err := r.GetMaterialBySlug(ctx, academyId, slug)
	if err != nil {
		return m, err
	}
	ap, err := r.GetMaterialProgress(ctx, accountId, academyId, m.Id)
	if err != nil {
		return m, err
	}
	m.AcademyMaterialProgress = ap
	return m, nil
}

func (r *academyRepository) CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	return m, r.db.WithContext(ctx).Create(&m).Error
}

func (r *academyRepository) ListMaterials(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error) {
	var list []entity.AcademyMaterial
	return list, r.db.WithContext(ctx).Where("academy_id = ?", academyId).Order("\"order\" ASC").Find(&list).Error
}

func (r *academyRepository) GetMaterialsWithContents(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error) {
	var materials []entity.AcademyMaterial
	return materials, r.db.WithContext(ctx).
		Where("academy_id = ?", academyId).
		Order("\"order\" ASC").
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Find(&materials).Error
}

func (r *academyRepository) UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	return m, r.db.WithContext(ctx).Save(&m).Error
}

func (r *academyRepository) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyMaterial{}, "id = ?", id).Error
}

func (r *academyRepository) GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error) {
	var c entity.AcademyContent
	return c, r.db.WithContext(ctx).First(&c, "id = ?", id).Error
}

func (r *academyRepository) GetContentBySlug(ctx context.Context, materialId uuid.UUID, order uint) (entity.AcademyContent, error) {
	var c entity.AcademyContent
	err := r.db.WithContext(ctx).Where("\"order\" = ? AND material_id = ?", order, materialId).First(&c).Error
	return c, err
}

func (r *academyRepository) GetContentWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID, order uint) (entity.AcademyContent, error) {
	c, err := r.GetContentBySlug(ctx, materialId, order)
	if err != nil {
		return c, err
	}
	ap, err := r.GetContentProgress(ctx, accountId, academyId, materialId, c.Id)
	if err != nil {
		return c, err
	}
	c.AcademyContentProgress = ap
	return c, nil
}

func (r *academyRepository) CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	return c, r.db.WithContext(ctx).Create(&c).Error
}

func (r *academyRepository) UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	return c, r.db.WithContext(ctx).Save(&c).Error
}

func (r *academyRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyContent{}, "id = ?", id).Error
}

func (r *academyRepository) CountContentsByMaterialID(ctx context.Context, materialId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyContent{}).Where("material_id = ?", materialId).Count(&count).Error
	return count, err
}

func (r *academyRepository) GetAcademyProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyProgress, error) {
	var existing entity.AcademyProgress
	err := r.db.WithContext(ctx).Where("account_id = ? AND academy_id = ?", accountId, academyId).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.AcademyProgress{
			AccountId: accountId,
			AcademyId: academyId,
			Status:    entity.StatusNotStarted,
		}, nil
	}
	return existing, err
}

func (r *academyRepository) UpsertAcademyProgress(ctx context.Context, p entity.AcademyProgress) (entity.AcademyProgress, error) {
	return p, r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Save(&p).Error
}

func (r *academyRepository) GetMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error) {
	var existing entity.AcademyMaterialProgress
	err := r.db.WithContext(ctx).Where("account_id = ? AND material_id = ?", accountId, materialId).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.AcademyMaterialProgress{
			AccountId:  accountId,
			AcademyId:  academyId,
			MaterialId: materialId,
			Status:     entity.StatusNotStarted,
		}, nil
	}
	return existing, err
}

func (r *academyRepository) UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error) {
	return p, r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Save(&p).Error
}

func (r *academyRepository) GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID, contentId uuid.UUID) (entity.AcademyContentProgress, error) {
	var existing entity.AcademyContentProgress
	err := r.db.WithContext(ctx).Where("account_id = ? AND content_id = ?", accountId, contentId).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.AcademyContentProgress{
			AccountId:  accountId,
			AcademyId:  academyId,
			MaterialId: materialId,
			ContentId:  contentId,
			Status:     entity.StatusNotStarted,
		}, nil
	}
	return existing, err
}

func (r *academyRepository) UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error) {
	return p, r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Save(&p).Error
}

func (r *academyRepository) DeleteContentProgressByContentID(ctx context.Context, contentId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("content_id = ?", contentId).Delete(&entity.AcademyContentProgress{}).Error
}

func (r *academyRepository) DeleteContentProgressByMaterialID(ctx context.Context, materialId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("material_id = ?", materialId).Delete(&entity.AcademyContentProgress{}).Error
}

func (r *academyRepository) DeleteMaterialProgressByMaterialID(ctx context.Context, materialId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("material_id = ?", materialId).Delete(&entity.AcademyMaterialProgress{}).Error
}

func (r *academyRepository) CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyContentProgress{}).
		Where("account_id = ? AND material_id = ? AND status = ?", accountId, materialId, entity.StatusCompleted).
		Count(&count).Error
	return count, err
}

func (r *academyRepository) CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
		Where("account_id = ? AND academy_id = ? AND status = ?", accountId, academyId, entity.StatusCompleted).
		Count(&count).Error
	return count, err
}

func (r *academyRepository) CountStartedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
		Where("account_id = ? AND academy_id = ? AND status IN ?", accountId, academyId, []string{entity.StatusInProgress, entity.StatusCompleted}).
		Count(&count).Error
	return count, err
}

func (r *academyRepository) DecrementMaterialOrdersGreaterThan(ctx context.Context, academyId uuid.UUID, order uint) error {
	return r.db.WithContext(ctx).Model(&entity.AcademyMaterial{}).
		Where("academy_id = ? AND \"order\" > ?", academyId, order).
		Update("\"order\"", gorm.Expr("\"order\" - 1")).Error
}

func (r *academyRepository) DecrementContentOrdersGreaterThan(ctx context.Context, materialId uuid.UUID, order uint) error {
	return r.db.WithContext(ctx).Model(&entity.AcademyContent{}).
		Where("material_id = ? AND \"order\" > ?", materialId, order).
		Update("\"order\"", gorm.Expr("\"order\" - 1")).Error
}

func (r *academyRepository) GetAccumulatedMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
		Where("account_id = ? AND academy_id = ?", accountId, academyId).
		Select("COALESCE(SUM(progress), 0)").
		Scan(&total).Error
	return total, err
}

func (r *academyRepository) BatchRecalculateMaterialProgress(ctx context.Context, materialId uuid.UUID) error {
	totalContents, err := r.CountContentsByMaterialID(ctx, materialId)
	if err != nil {
		return err
	}

	if totalContents == 0 {
		return r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
			Where("material_id = ?", materialId).
			Updates(map[string]interface{}{
				"progress":                 100,
				"status":                   entity.StatusCompleted,
				"total_completed_contents": 0,
			}).Error
	}

	return r.db.WithContext(ctx).Exec(`
		UPDATE academy_material_progresses amp
		SET 
			total_completed_contents = (
				SELECT COUNT(id) FROM academy_content_progresses acp 
				WHERE acp.material_id = amp.material_id AND acp.account_id = amp.account_id AND acp.status = 'COMPLETED'
			),
			progress = (
				SELECT COUNT(id) FROM academy_content_progresses acp 
				WHERE acp.material_id = amp.material_id AND acp.account_id = amp.account_id AND acp.status = 'COMPLETED'
			)::float / ? * 100,
			status = CASE 
				WHEN (
					SELECT COUNT(id) FROM academy_content_progresses acp 
					WHERE acp.material_id = amp.material_id AND acp.account_id = amp.account_id AND acp.status = 'COMPLETED'
				) >= ? THEN 'COMPLETED' 
				ELSE 'IN_PROGRESS' 
			END,
			completed_at = CASE 
				WHEN (
					SELECT COUNT(id) FROM academy_content_progresses acp 
					WHERE acp.material_id = amp.material_id AND acp.account_id = amp.account_id AND acp.status = 'COMPLETED'
				) >= ? THEN NOW() 
				ELSE NULL 
			END
		WHERE amp.material_id = ?
	`, totalContents, totalContents, totalContents, materialId).Error
}

func (r *academyRepository) BatchRecalculateAcademyProgress(ctx context.Context, academyId uuid.UUID) error {
	totalMaterials, err := r.CountMaterialsByAcademyID(ctx, academyId)
	if err != nil {
		return err
	}

	if totalMaterials == 0 {
		return r.db.WithContext(ctx).Model(&entity.AcademyProgress{}).
			Where("academy_id = ?", academyId).
			Updates(map[string]interface{}{
				"progress":                  100,
				"status":                    entity.StatusCompleted,
				"total_completed_materials": 0,
			}).Error
	}

	return r.db.WithContext(ctx).Exec(`
		UPDATE academy_progress ap
		SET 
			progress = (
				SELECT COALESCE(SUM(amp.progress), 0) FROM academy_material_progresses amp
				WHERE amp.academy_id = ap.academy_id AND amp.account_id = ap.account_id
			)::float / ?,
			total_completed_materials = (
				SELECT COUNT(id) FROM academy_material_progresses amp
				WHERE amp.academy_id = ap.academy_id AND amp.account_id = ap.account_id AND amp.status = 'COMPLETED'
			),
			status = CASE 
				WHEN (
					SELECT COUNT(id) FROM academy_material_progresses amp
					WHERE amp.academy_id = ap.academy_id AND amp.account_id = ap.account_id AND amp.status = 'COMPLETED'
				) >= ? THEN 'COMPLETED' 
				ELSE 'IN_PROGRESS' 
			END,
			completed_at = CASE 
				WHEN (
					SELECT COUNT(id) FROM academy_material_progresses amp
					WHERE amp.academy_id = ap.academy_id AND amp.account_id = ap.account_id AND amp.status = 'COMPLETED'
				) >= ? THEN NOW() 
				ELSE NULL 
			END
		WHERE ap.academy_id = ?
	`, totalMaterials, totalMaterials, totalMaterials, academyId).Error
}

func (r *academyRepository) GetMaterialProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialIds []uuid.UUID) (map[uuid.UUID]entity.AcademyMaterialProgress, error) {
	var progresses []entity.AcademyMaterialProgress
	result := r.db.WithContext(ctx).Where("account_id = ? AND academy_id = ? AND material_id IN ?", accountId, academyId, materialIds).Find(&progresses)

	progressMap := make(map[uuid.UUID]entity.AcademyMaterialProgress)
	for _, p := range progresses {
		progressMap[p.MaterialId] = p
	}

	return progressMap, result.Error
}

func (r *academyRepository) GetContentProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, contentIds []uuid.UUID) (map[uuid.UUID]entity.AcademyContentProgress, error) {
	var progresses []entity.AcademyContentProgress
	result := r.db.WithContext(ctx).Where("account_id = ? AND academy_id = ? AND content_id IN ?", accountId, academyId, contentIds).Find(&progresses)

	progressMap := make(map[uuid.UUID]entity.AcademyContentProgress)
	for _, p := range progresses {
		progressMap[p.ContentId] = p
	}

	return progressMap, result.Error
}
