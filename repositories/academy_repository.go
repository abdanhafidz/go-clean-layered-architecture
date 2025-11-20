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
	CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error)
	GetAcademyWithProgress(ctx context.Context, accountId string, slug string) (entity.Academy, error)
	ListAcademy(ctx context.Context, accountId string) ([]entity.Academy, error)
	UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error

	GetAcademyWithMaterials(ctx context.Context, id uuid.UUID) (entity.Academy, []entity.AcademyMaterial, error)
	CountMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) (int64, error)

	// Material
	GetMaterialBySlug(ctx context.Context, academy_id uuid.UUID, materialSlug string) (entity.AcademyMaterial, error)

	CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error)
	ListMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyMaterial, error)
	UpdateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error)
	DeleteMaterial(ctx context.Context, id uuid.UUID) error

	GetMaterialWithContents(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, []entity.AcademyContent, error)
	

	// Content
	GetContentBySlug(ctx context.Context, materialId uuid.UUID, order uint) (entity.AcademyContent, error)

	CreateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	GetContentByID(ctx context.Context, id uuid.UUID) (entity.AcademyContent, error)
	CountContentsByMaterialID(ctx context.Context, materialId uuid.UUID) (int64, error)
	UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	// Progress
	GetAcademyProgress(ctx context.Context, accountId string, academyId string) (entity.AcademyProgress, error)
	UpsertAcademyProgress(ctx context.Context, p entity.AcademyProgress) (entity.AcademyProgress, error)

	GetMaterialProgress(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (entity.AcademyMaterialProgress, error)
	UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error)

	UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error)
	GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyContentProgress, error)

	CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error)
	CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error)
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

func (r *academyRepository) GetAcademyWithProgress(ctx context.Context, accountId string, slug string) (entity.Academy, error) {
    var a entity.Academy
	var err error
    a,err = r.GetAcademyBySlug(ctx, slug)
	if err != nil {
		return a, err
	}
    
    academyId := a.Id.String() 

    var ap entity.AcademyProgress 

	ap,err = r.GetAcademyProgress(ctx, accountId, academyId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return a, err
	}
    a.AcademyProgresss = ap

    return a, nil
}

func (r *academyRepository) ListAcademy(ctx context.Context, accountId string) ([]entity.Academy, error) {
    var list []entity.Academy

    if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
        return nil, err
    }

    for i := range list {
        academyId := list[i].Id.String()
        ap, err := r.GetAcademyProgress(ctx, accountId, academyId)
        if err != nil && err != gorm.ErrRecordNotFound {
            return nil, err
        }
        list[i].AcademyProgresss = ap 
    }

    return list, nil
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
func (r *academyRepository) CreateMaterial(ctx context.Context, m entity.AcademyMaterial) (entity.AcademyMaterial, error) {
    return m, r.db.WithContext(ctx).Create(&m).Error
}

func (r *academyRepository) GetMaterialByID(ctx context.Context, id uuid.UUID) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	return m, r.db.WithContext(ctx).First(&m, "id = ?", id).Error
}

func (r *academyRepository) GetMaterialBySlug(ctx context.Context, academy_id uuid.UUID, materialSlug string) (entity.AcademyMaterial, error) {
	var m entity.AcademyMaterial
	return m, r.db.WithContext(ctx).First(&m, "academy_id = ? AND slug = ?", academy_id, materialSlug).Error
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

func (r *academyRepository) CountContentsByMaterialID(ctx context.Context, materialId uuid.UUID) (int64, error) {
    var count int64

    query := r.db.WithContext(ctx).
        Where("academy_material_id = ?", materialId)
        
    err := query.Model(&entity.AcademyContent{}). 
        Count(&count). 
        Error
    
    return count, err
}

func (r *academyRepository) UpdateContent(ctx context.Context, c entity.AcademyContent) (entity.AcademyContent, error) {
	return c, r.db.WithContext(ctx).Save(&c).Error
}

func (r *academyRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyContent{}, "id = ?", id).Error
}

func (r *academyRepository) GetContentBySlug(ctx context.Context, materialId uuid.UUID, order uint) (entity.AcademyContent, error) {
    var c entity.AcademyContent
    result := r.db.WithContext(ctx).
        // Escape "order" with backslashes and double quotes: \"order\"
        Where("\"order\" = ?", order). 
        Where("academy_material_id = ?", materialId).
        First(&c)
    
    return c, result.Error
}


// ========== PROGRESS ==========

func (r *academyRepository) GetAcademyProgress(ctx context.Context, accountId string, academyId string) (entity.AcademyProgress, error) {
	var existing entity.AcademyProgress
    
    err := r.db.WithContext(ctx).
        Where("account_id = ? AND academy_id = ?", accountId, academyId).
        First(&existing).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        accUUID, _ := uuid.Parse(accountId)
        acaUUID, _ := uuid.Parse(academyId)

        return entity.AcademyProgress{
            AccountId:               accUUID,
            AcademyId:               acaUUID,
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


func(r *academyRepository) UpsertAcademyProgress(ctx context.Context, p entity.AcademyProgress) (entity.AcademyProgress, error) {
	var existing entity.AcademyProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ?", p.AccountId, p.AcademyId).Error
	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}
	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

func (r *academyRepository) UpsertMaterialProgress(ctx context.Context, p entity.AcademyMaterialProgress) (entity.AcademyMaterialProgress, error) {
	var existing entity.AcademyMaterialProgress
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ? AND academy_material_id = ?", p.AccountId, p.AcademyId,p.AcademyMaterialId).Error

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
	err := r.db.WithContext(ctx).First(&existing, "account_id = ? AND academy_id = ? AND academy_material_id = ? AND content_id = ?", p.AccountId, p.AcademyId,p.AcademyMaterialId,p.ContentId).Error
	
	if err == gorm.ErrRecordNotFound {
		return p, r.db.WithContext(ctx).Create(&p).Error
	}

	return p, r.db.WithContext(ctx).Model(&existing).Updates(p).Error
}

func (r *academyRepository) GetContentProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (entity.AcademyContentProgress, error) {
	var p entity.AcademyContentProgress
	return p, r.db.WithContext(ctx).First(&p, "account_id = ? AND academy_id = ?", accountId, academyId).Error
}

func (r *academyRepository) CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Where("account_id = ? AND academy_material_id = ? AND status = ?", accountId, materialId, "COMPLETED")
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