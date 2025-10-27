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

	// Material
	CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error)
	GetMaterialBySlug(ctx context.Context, slug string) (entity.AcademyMaterial, error)
	ListMaterialsByAcademyID(ctx context.Context, academyId uint) ([]entity.AcademyMaterial, error)
	UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	DeleteMaterial(ctx context.Context, id uuid.UUID) error

	// Content
	CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error)
	ListContentsByMaterialID(ctx context.Context, materialId uint) ([]entity.AcademyContent, error)
	UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	// Progress
	UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error)
	GetMaterialProgress(ctx context.Context, accountId uint, materialId uint) (entity.AcademyMaterialProgress, error)

	UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error)
	GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyContentProgress, error)
}

type academyRepository struct {
	db *gorm.DB
}

func NewAcademyRepository(db *gorm.DB) AcademyRepository {
	return &academyRepository{db: db}
}

// Academy
func (r *academyRepository) CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error) {
	if err := r.db.WithContext(ctx).Create(&a).Error; err != nil {
		return entity.Academy{}, err
	}
	return a, nil
}

func (r *academyRepository) GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error) {
	var a entity.Academy
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return entity.Academy{}, err
	}
	return a, nil
}

func (r *academyRepository) GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error) {
	var a entity.Academy
	if err := r.db.WithContext(ctx).First(&a, "slug = ?", slug).Error; err != nil {
		return entity.Academy{}, err
	}
	return a, nil
}

func (r *academyRepository) ListAcademy(ctx context.Context) ([]entity.Academy, error) {
	var list []entity.Academy
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *academyRepository) UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error) {
	if err := r.db.WithContext(ctx).Save(&a).Error; err != nil {
		return entity.Academy{}, err
	}
	return a, nil
}

func (r *academyRepository) DeleteAcademy(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Academy{}, "id = ?", id).Error
}

// Material
func (r *academyRepository) CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.AcademyMaterial{}, err
	}
	return m, nil
}

func (r *academyRepository) GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return entity.AcademyMaterial{}, err
	}
	return m, nil
}

func (r *academyRepository) GetMaterialBySlug(ctx context.Context, slug string) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	if err := r.db.WithContext(ctx).First(&m, "slug = ?", slug).Error; err != nil {
		return entity.AcademyMaterial{}, err
	}
	return m, nil
}

func (r *academyRepository) ListMaterialsByAcademyID(ctx context.Context, academyId uint) ([]entity.AcademyMaterial, error) {
	var list []entity.AcademyMaterial
	if err := r.db.WithContext(ctx).Where("academy_id = ?", academyId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *academyRepository) UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return entity.AcademyMaterial{}, err
	}
	return m, nil
}

func (r *academyRepository) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyMaterial{}, "id = ?", id).Error
}

// Content
func (r *academyRepository) CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return entity.AcademyContent{}, err
	}
	return c, nil
}

func (r *academyRepository) GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error) {
	var c entity.AcademyContent
	if err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		return entity.AcademyContent{}, err
	}
	return c, nil
}

func (r *academyRepository) ListContentsByMaterialID(ctx context.Context, materialId uint) ([]entity.AcademyContent, error) {
	var list []entity.AcademyContent
	if err := r.db.WithContext(ctx).Where("academy_material_id = ?", materialId).Order("order ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *academyRepository) UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	if err := r.db.WithContext(ctx).Save(&c).Error; err != nil {
		return entity.AcademyContent{}, err
	}
	return c, nil
}

func (r *academyRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyContent{}, "id = ?", id).Error
}

// Progress
func (r *academyRepository) UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error) {
	var existing entity.AcademyMaterialProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_material_id = ?", p.AccountId, p.AcademyMaterialId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
				return entity.AcademyMaterialProgress{}, err
			}
			return p, nil
		}
		return entity.AcademyMaterialProgress{}, err
	}
	if err := r.db.WithContext(ctx).Model(&existing).Updates(p).Error; err != nil {
		return entity.AcademyMaterialProgress{}, err
	}
	return existing, nil
}

func (r *academyRepository) GetMaterialProgress(ctx context.Context, accountId uint, materialId uint) (entity.AcademyMaterialProgress, error) {
	var res entity.AcademyMaterialProgress
	if err := r.db.WithContext(ctx).First(&res, "account_id = ? AND academy_material_id = ?", accountId, materialId).Error; err != nil {
		return entity.AcademyMaterialProgress{}, err
	}
	return res, nil
}

func (r *academyRepository) UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error) {
	var existing entity.AcademyContentProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ?", p.AccountId, p.AcademyId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
				return entity.AcademyContentProgress{}, err
			}
			return p, nil
		}
		return entity.AcademyContentProgress{}, err
	}
	if err := r.db.WithContext(ctx).Model(&existing).Updates(p).Error; err != nil {
		return entity.AcademyContentProgress{}, err
	}
	return existing, nil
}

func (r *academyRepository) GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyContentProgress, error) {
	var res entity.AcademyContentProgress
	if err := r.db.WithContext(ctx).First(&res, "account_id = ? AND academy_id = ?", accountId, academyId).Error; err != nil {
		return entity.AcademyContentProgress{}, err
	}
	return res, nil
}
