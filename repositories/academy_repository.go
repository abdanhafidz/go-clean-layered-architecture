package repositories

import (
	"context"
	"errors"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyRepository interface {
	// Academy
	GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error)
	GetAcademyWithProgress(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error)

	CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error

	ListAcademy(ctx context.Context, accountId uuid.UUID) ([]entity.Academy, error)
	GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error)
	CountMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) (int64, error)

	// Material
	GetMaterialBySlug(ctx context.Context, academy_id uuid.UUID, materialSlug string) (entity.AcademyMaterial, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error)
	GetMaterialWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, slug string) (entity.AcademyMaterial, error)

	CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	DeleteMaterial(ctx context.Context, id uuid.UUID) error

	ListMaterials(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error)
	GetMaterialWithContents(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, []entity.AcademyContent, error)

	// Content
	GetContentBySlug(ctx context.Context, materialId uuid.UUID, order uint) (entity.AcademyContent, error)
	GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error)
	GetContentWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID,materialId uuid.UUID, order uint) (entity.AcademyContent, error)

	CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	CountContentsByMaterialID(ctx context.Context, materialId uuid.UUID) (int64, error)

	// Progress
	GetAcademyProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyProgress, error)
	GetMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error)
	GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID, contentId uuid.UUID) (entity.AcademyContentProgress, error)

	UpsertAcademyProgress(ctx context.Context, p entity.AcademyProgress) (entity.AcademyProgress, error)
	UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error)
	UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error)

	CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error)
	CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error)
}

type academyRepository struct{ db *gorm.DB }

func NewAcademyRepository(db *gorm.DB) AcademyRepository {
	return &academyRepository{db: db}
}

// ========== ACADEMY ==========
func (r *academyRepository) GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error) {
	var a entity.Academy
	err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error
	if err != nil {
		return entity.Academy{}, nil, err
	}

	var m []entity.AcademyMaterial
	return a, m, r.db.WithContext(ctx).Where("academy_id = ?", id).Find(&m).Error
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
	var a entity.Academy
	var err error
	a, err = r.GetAcademyBySlug(ctx, slug)
	if err != nil {
		return a, err
	}

	academyId := a.Id

	var ap entity.AcademyProgress

	ap, err = r.GetAcademyProgress(ctx, accountId, academyId)
	if err != nil && err != gorm.ErrRecordNotFound {
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

	for i := range list {
		academyId := list[i].Id
		ap, err := r.GetAcademyProgress(ctx, accountId, academyId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		list[i].AcademyProgresss = ap
	}

	return list, nil
}

func (r *academyRepository) CountMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Where("academy_id = ?", academyId)

	err := query.Model(&entity.AcademyMaterial{}).
		Count(&count).
		Error

	return count, err
}

// ========== MATERIAL ==========
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
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		return entity.AcademyMaterial{}, nil, err
	}

	var c []entity.AcademyContent
	return m, c, r.db.WithContext(ctx).Where("material_id = ?", id).Order("order asc").Find(&c).Error
}

func (r *academyRepository) GetMaterialWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, slug string) (entity.AcademyMaterial, error){
	var m entity.AcademyMaterial
	var err error
	m, err = r.GetMaterialBySlug(ctx, academyId, slug)
	if err != nil {
		return m, err
	}

	MaterialId := m.Id

	var ap entity.AcademyMaterialProgress

	ap, err = r.GetMaterialProgress(ctx, accountId,academyId, MaterialId)
	if err != nil && err != gorm.ErrRecordNotFound {
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
	return list, r.db.WithContext(ctx).Where("academy_id = ?", academyId).Find(&list).Error
}

func (r *academyRepository) UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	return m, r.db.WithContext(ctx).Save(&m).Error
}

func (r *academyRepository) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyMaterial{}, "id = ?", id).Error
}

// ========== CONTENT ==========

func (r *academyRepository) GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error) {
	var c entity.AcademyContent
	return c, r.db.WithContext(ctx).First(&c, "id = ?", id).Error
}

func (r *academyRepository) GetContentBySlug(ctx context.Context, materialId uuid.UUID, order uint) (entity.AcademyContent, error) {
	var c entity.AcademyContent
	result := r.db.WithContext(ctx).
		// Escape "order" with backslashes and double quotes: \"order\"
		Where("\"order\" = ?", order).
		Where("material_id = ?", materialId).
		First(&c)

	return c, result.Error
}

func (r *academyRepository) GetContentWithProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID,materialId uuid.UUID, order uint) (entity.AcademyContent, error){
	var c entity.AcademyContent
	var err error
	c, err = r.GetContentBySlug(ctx,materialId,order)
	if err != nil {
		return c, err
	}

	var ap entity.AcademyContentProgress

	ap, err = r.GetContentProgress(ctx, accountId, academyId, materialId,c.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
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

	query := r.db.WithContext(ctx).
		Where("material_id = ?", materialId)

	err := query.Model(&entity.AcademyContent{}).
		Count(&count).
		Error

	return count, err
}

// ========== PROGRESS ==========

func (r *academyRepository) GetAcademyProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyProgress, error) {
	var existing entity.AcademyProgress

	err := r.db.WithContext(ctx).
		Where("account_id = ? AND academy_id = ?", accountId, academyId).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.AcademyProgress{
			AccountId:               accountId,
			AcademyId:               academyId,
			Status:                  "NOT_STARTED",
			Progress:                0,
			TotalCompletedMaterials: 0,
		}, nil
	}

	if err != nil {
		return existing, err
	}
	return existing, nil
}

func (r *academyRepository) UpsertAcademyProgress(ctx context.Context, p entity.AcademyProgress) (entity.AcademyProgress, error) {
	var existing entity.AcademyProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ?", p.AccountId, p.AcademyId).Error
	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}
	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

func (r *academyRepository) GetMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error) {
	var existing entity.AcademyMaterialProgress

	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ?  AND material_id = ?", accountId, academyId, materialId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		return entity.AcademyMaterialProgress{
			AccountId:              accountId,
			AcademyId:              academyId,
			MaterialId:             materialId,
			Progress:               0,
			TotalCompletedContents: 0,
			Status:                 "NOT_STARTED",
		}, nil
	}

	if err != nil {
		return existing, err
	}
	return existing, nil
}

func (r *academyRepository) UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error) {
	var existing entity.AcademyMaterialProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ? AND material_id = ?", p.AccountId, p.AcademyId, p.MaterialId).Error

	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}

	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

func (r *academyRepository) GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialId uuid.UUID, contentId uuid.UUID) (entity.AcademyContentProgress, error) {
	var existing entity.AcademyContentProgress

	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ?  AND material_id = ? AND content_id = ?", accountId, academyId, materialId,contentId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.AcademyContentProgress{
			AccountId:  accountId,
			AcademyId:  academyId,
			MaterialId: materialId,
			ContentId:  contentId,
			Status:     "NOT_STARTED",
		}, nil
	}

	if err != nil {
		return existing, err
	}
	return existing, nil
}

func (r *academyRepository) UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error) {
	var existing entity.AcademyContentProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ? AND material_id = ? AND content_id = ?", p.AccountId, p.AcademyId, p.MaterialId, p.ContentId).Error

	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}

	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

// UTILS

func (r *academyRepository) CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Where("account_id = ? AND material_id = ? AND status = ?", accountId, materialId, "COMPLETED")
	err := query.Model(&entity.AcademyContentProgress{}).
		Count(&count).
		Error
	return count, err
}

func (r *academyRepository) CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Where("account_id = ? AND academy_id = ? AND status = ?", accountId, academyId, "COMPLETED")
	err := query.Model(&entity.AcademyMaterialProgress{}).
		Count(&count).
		Error
	return count, err
}
