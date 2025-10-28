package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyRepository interface {
	// Academy
	CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error)
	ListAcademy(ctx context.Context) ([]entity.Academy, error)
	UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error

	GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error)

	// Material
	CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error)
	GetMaterialBySlug(ctx context.Context, slug string) (entity.AcademyMaterial, error)
	ListMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error)
	UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	DeleteMaterial(ctx context.Context, id uuid.UUID) error

	GetMaterialWithContents(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, []entity.AcademyContent, error)

	// Content
	CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error)
	ListContentsByMaterialID(ctx context.Context, materialId uuid.UUID) ([]entity.AcademyContent, error)
	UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	// Progress
	UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error)
	GetMaterialProgress(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error)

	UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error)
	GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyContentProgress, error)
}

type academyRepository struct{ db *gorm.DB }

func NewAcademyRepository(db *gorm.DB) AcademyRepository {
	return &academyRepository{db: db}
}

// ========== ACADEMY ==========
func (r *academyRepository) CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error) {
	return a, r.db.WithContext(ctx).Create(&a).Error
}

func (r *academyRepository) GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	var a entity.Academy
	return a, r.db.WithContext(ctx).First(&a, "id = ?", id).Error
}

func (r *academyRepository) GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error) {
	var a entity.Academy
	return a, r.db.WithContext(ctx).First(&a, "slug = ?", slug).Error
}

func (r *academyRepository) ListAcademy(ctx context.Context) ([]entity.Academy, error) {
	var list []entity.Academy
	return list, r.db.WithContext(ctx).Find(&list).Error
}

func (r *academyRepository) UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error) {
	return a, r.db.WithContext(ctx).Save(&a).Error
}

func (r *academyRepository) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Academy{}, "id = ?", id).Error
}

func (r *academyRepository) GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error) {
	var a entity.Academy
	err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error
	if err != nil {
		return entity.Academy{}, nil, err
	}

	var m []entity.AcademyMaterial
	return a, m, r.db.WithContext(ctx).Where("academy_id = ?", id).Find(&m).Error
}

// ========== MATERIAL ==========
func (r *academyRepository) CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	return m, r.db.WithContext(ctx).Create(&m).Error
}

func (r *academyRepository) GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	return m, r.db.WithContext(ctx).First(&m, "id = ?", id).Error
}

func (r *academyRepository) GetMaterialBySlug(ctx context.Context, slug string) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	return m, r.db.WithContext(ctx).First(&m, "slug = ?", slug).Error
}

func (r *academyRepository) ListMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error) {
	var list []entity.AcademyMaterial
	return list, r.db.WithContext(ctx).Where("academy_id = ?", academyId).Find(&list).Error
}

func (r *academyRepository) UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	return m, r.db.WithContext(ctx).Save(&m).Error
}

func (r *academyRepository) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyMaterial{}, "id = ?", id).Error
}

func (r *academyRepository) GetMaterialWithContents(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, []entity.AcademyContent, error) {
	var m entity.AcademyMaterial
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		return entity.AcademyMaterial{}, nil, err
	}

	var c []entity.AcademyContent
	return m, c, r.db.WithContext(ctx).Where("academy_material_id = ?", id).Order("order asc").Find(&c).Error
}

// ========== CONTENT ==========
func (r *academyRepository) CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	return c, r.db.WithContext(ctx).Create(&c).Error
}

func (r *academyRepository) GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error) {
	var c entity.AcademyContent
	return c, r.db.WithContext(ctx).First(&c, "id = ?", id).Error
}

func (r *academyRepository) ListContentsByMaterialID(ctx context.Context, materialId uuid.UUID) ([]entity.AcademyContent, error) {
	var list []entity.AcademyContent
	return list, r.db.WithContext(ctx).Where("academy_material_id = ?", materialId).Order("order asc").Find(&list).Error
}

func (r *academyRepository) UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	return c, r.db.WithContext(ctx).Save(&c).Error
}

func (r *academyRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyContent{}, "id = ?", id).Error
}

// ========== PROGRESS ==========
func (r *academyRepository) UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error) {
	var existing entity.AcademyMaterialProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_material_id = ?", p.AccountId, p.AcademyMaterialId).Error

	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}

	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

func (r *academyRepository) GetMaterialProgress(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error) {
	var p entity.AcademyMaterialProgress
	return p, r.db.WithContext(ctx).First(&p, "account_id = ? AND academy_material_id = ?", accountId, materialId).Error
}

func (r *academyRepository) UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error) {
	var existing entity.AcademyContentProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ?", p.AccountId, p.AcademyId).Error

	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}

	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

func (r *academyRepository) GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyContentProgress, error) {
	var p entity.AcademyContentProgress
	return p, r.db.WithContext(ctx).First(&p, "account_id = ? AND academy_id = ?", accountId, academyId).Error
}
