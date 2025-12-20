package repositories

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AcademyRepository interface {
	Atomic(ctx context.Context, fn func(r AcademyRepository) error) error

	GetAcademyByID(ctx context.Context, id uuid.UUID) (entity.Academy, error)
	GetAcademyBySlug(ctx context.Context, slug string) (entity.Academy, error)
	GetAcademyByCode(ctx context.Context, code string) (entity.Academy, error)
	GetAcademyWithProgress(ctx context.Context, accountId uuid.UUID, slug string) (entity.Academy, error)
	CreateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	UpdateAcademy(ctx context.Context, a entity.Academy) (entity.Academy, error)
	DeleteAcademy(ctx context.Context, id uuid.UUID) error
	ListVisibleAcademy(ctx context.Context, accountId uuid.UUID, p *entity.Pagination) ([]entity.Academy, int64, error)
	ListAssignmentsByAccount(ctx context.Context, accountId uuid.UUID) ([]entity.AcademyAssign, error)
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
	DeleteAcademyProgressByAcademyID(ctx context.Context, academyId uuid.UUID) error
	DeleteAcademyAssignByAcademyID(ctx context.Context, academyId uuid.UUID) error
	DeleteMaterialProgressByAcademyID(ctx context.Context, academyId uuid.UUID) error
	DeleteContentProgressByAcademyID(ctx context.Context, academyId uuid.UUID) error
	DeleteContentsByMaterialID(ctx context.Context, materialId uuid.UUID) error
	DeleteContentsByAcademyID(ctx context.Context, academyId uuid.UUID) error
	DeleteMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) error

	CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error)
	CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error)
	CountStartedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error)
	DecrementMaterialOrdersGreaterThan(ctx context.Context, academyId uuid.UUID, order uint) error
	DecrementContentOrdersGreaterThan(ctx context.Context, materialId uuid.UUID, order uint) error
	GetAccumulatedMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (float64, error)

	GetMaterialProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialIds []uuid.UUID) (map[uuid.UUID]entity.AcademyMaterialProgress, error)
	GetBatchMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (map[uuid.UUID]entity.AcademyMaterialProgress, error)
	GetContentProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, contentIds []uuid.UUID) (map[uuid.UUID]entity.AcademyContentProgress, error)

	ListAccountsByMaterialId(ctx context.Context, materialId uuid.UUID) ([]uuid.UUID, error)
	ListAccountsByContentMaterialId(ctx context.Context, materialId uuid.UUID) ([]uuid.UUID, error)

	BatchRecalculateAcademyProgress(ctx context.Context, academyId uuid.UUID) error
	BatchRecalculateMaterialProgress(ctx context.Context, materialId uuid.UUID) error

	IsAccountAssignedToAcademy(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (bool, error)

	AssignAccountToAcademy(ctx context.Context, assign entity.AcademyAssign) (entity.AcademyAssign, error)
	GetAssignmentByAcademyAndAccount(ctx context.Context, academyId uuid.UUID, accountId uuid.UUID) (entity.AcademyAssign, error)
	ListAssignmentsByAcademy(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyAssign, error)
	RemoveAssignment(ctx context.Context, id uuid.UUID) error
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

func (r *academyRepository) GetAcademyByCode(ctx context.Context, code string) (entity.Academy, error) {
	var a entity.Academy
	return a, r.db.WithContext(ctx).First(&a, "code = ?", code).Error
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
	a.AcademyProgress = ap
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

func (r *academyRepository) ListVisibleAcademy(ctx context.Context, accountId uuid.UUID, p *entity.Pagination) ([]entity.Academy, int64, error) {
	var list []entity.Academy
	var total int64

	sub := r.db.WithContext(ctx).Model(&entity.AcademyAssign{}).Select("academy_id").Where("account_id = ?", accountId)

	q := r.db.WithContext(ctx).Model(&entity.Academy{}).
		Where(r.db.Where("is_public = ?", true).Or("id IN (?)", sub))

	if p != nil {
		if p.RegisterStatus != nil {
			if *p.RegisterStatus == 1 {
				q = q.Where("id IN (?)", sub)
			} else if *p.RegisterStatus == 0 {
				q = q.Where("id NOT IN (?)", sub)
			}
		}

		sortColumn := "title"
		switch p.SortBy {
		case "title":
			sortColumn = "title"
		case "createdAt", "created_at":
			sortColumn = "created_at"
		}

		if s := strings.TrimSpace(p.Search); s != "" {
			s = strings.Trim(s, "\"'")
			s = strings.ToLower(s)
			like := "%" + s + "%"
			q = q.Where("LOWER(title) LIKE ? OR LOWER(slug) LIKE ?", like, like)
		}

		sortDirection := "DESC"
		if strings.ToLower(p.Order) == "asc" {
			sortDirection = "ASC"
		}

		q = q.Order(fmt.Sprintf("%s %s", sortColumn, sortDirection))

		if p.Limit > 0 {
			q = q.Limit(p.Limit)
		}
		if p.Offset > 0 {
			q = q.Offset(p.Offset)
		}
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	if len(list) == 0 {
		return list, total, nil
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
		return nil, 0, err
	}

	progressMap := make(map[uuid.UUID]entity.AcademyProgress)
	for _, p := range progressList {
		progressMap[p.AcademyId] = p
	}

	var assigns []entity.AcademyAssign
	if err := r.db.WithContext(ctx).
		Where("account_id = ?", accountId).
		Where("academy_id IN ?", academyIDs).
		Find(&assigns).Error; err != nil {
		return nil, 0, err
	}
	assignedMap := make(map[uuid.UUID]bool)
	for _, a := range assigns {
		assignedMap[a.AcademyId] = true
	}

	for i := range list {
		if pr, exists := progressMap[list[i].Id]; exists {
			list[i].AcademyProgress = pr
		} else {
			list[i].AcademyProgress = entity.AcademyProgress{
				AccountId: accountId,
				AcademyId: list[i].Id,
				Status:    entity.StatusNotStarted,
			}
		}
		if assignedMap[list[i].Id] {
			list[i].RegisterStatus = 1
		} else {
			list[i].RegisterStatus = 0
		}
	}
	return list, total, nil
}

func (r *academyRepository) ListAssignmentsByAccount(ctx context.Context, accountId uuid.UUID) ([]entity.AcademyAssign, error) {
	var assigns []entity.AcademyAssign
	err := r.db.WithContext(ctx).Where("account_id = ?", accountId).Find(&assigns).Error
	return assigns, err
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
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "academy_id"}},
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
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "material_id"}},
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
			Status:     entity.StatusInProgress,
		}, nil
	}
	return existing, err
}

func (r *academyRepository) UpsertContentProgress(ctx context.Context, p entity.AcademyContentProgress) (entity.AcademyContentProgress, error) {
	return p, r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "content_id"}},
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

func (r *academyRepository) DeleteAcademyProgressByAcademyID(ctx context.Context, academyId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("academy_id = ?", academyId).Delete(&entity.AcademyProgress{}).Error
}

func (r *academyRepository) DeleteAcademyAssignByAcademyID(ctx context.Context, academyId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("academy_id = ?", academyId).Delete(&entity.AcademyAssign{}).Error
}

func (r *academyRepository) DeleteMaterialProgressByAcademyID(ctx context.Context, academyId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("academy_id = ?", academyId).Delete(&entity.AcademyMaterialProgress{}).Error
}

func (r *academyRepository) DeleteContentProgressByAcademyID(ctx context.Context, academyId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("academy_id = ?", academyId).Delete(&entity.AcademyContentProgress{}).Error
}

func (r *academyRepository) DeleteContentsByMaterialID(ctx context.Context, materialId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("material_id = ?", materialId).Delete(&entity.AcademyContent{}).Error
}

func (r *academyRepository) DeleteContentsByAcademyID(ctx context.Context, academyId uuid.UUID) error {
	// Delete contents where material_id matches any material in the academy
	return r.db.WithContext(ctx).Where("material_id IN (?)",
		r.db.Model(&entity.AcademyMaterial{}).Select("id").Where("academy_id = ?", academyId),
	).Delete(&entity.AcademyContent{}).Error
}

func (r *academyRepository) DeleteMaterialsByAcademyID(ctx context.Context, academyId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("academy_id = ?", academyId).Delete(&entity.AcademyMaterial{}).Error
}

func (r *academyRepository) CountCompletedContentsByMaterialAndAccount(ctx context.Context, accountId uuid.UUID, materialId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyContentProgress{}).
		Where("account_id = ? AND material_id = ? AND status = ?", accountId, materialId, entity.StatusFinished).
		Count(&count).Error
	return count, err
}

func (r *academyRepository) CountCompletedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
		Where("account_id = ? AND academy_id = ? AND status = ?", accountId, academyId, entity.StatusFinished).
		Count(&count).Error
	return count, err
}

func (r *academyRepository) CountStartedMaterialsByAcademyAndAccount(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
		Where("account_id = ? AND academy_id = ? AND status IN ?", accountId, academyId, []string{entity.StatusInProgress, entity.StatusFinished}).
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

func (r *academyRepository) GetMaterialProgressBatch(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, materialIds []uuid.UUID) (map[uuid.UUID]entity.AcademyMaterialProgress, error) {
	var progresses []entity.AcademyMaterialProgress
	result := r.db.WithContext(ctx).Where("account_id = ? AND academy_id = ? AND material_id IN ?", accountId, academyId, materialIds).Find(&progresses)

	progressMap := make(map[uuid.UUID]entity.AcademyMaterialProgress)
	for _, p := range progresses {
		progressMap[p.MaterialId] = p
	}

	return progressMap, result.Error
}

func (r *academyRepository) GetBatchMaterialProgress(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (map[uuid.UUID]entity.AcademyMaterialProgress, error) {
	var progresses []entity.AcademyMaterialProgress
	result := r.db.WithContext(ctx).Where("account_id = ? AND academy_id = ?", accountId, academyId).Find(&progresses)

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

func (r *academyRepository) ListAccountsByMaterialId(ctx context.Context, materialId uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entity.AcademyMaterialProgress{}).
		Where("material_id = ?", materialId).
		Distinct().
		Pluck("account_id", &ids).Error
	return ids, err
}

func (r *academyRepository) ListAccountsByContentMaterialId(ctx context.Context, materialId uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entity.AcademyContentProgress{}).
		Where("material_id = ?", materialId).
		Distinct().
		Pluck("account_id", &ids).Error
	return ids, err
}

func (r *academyRepository) BatchRecalculateMaterialProgress(ctx context.Context, materialId uuid.UUID) error {
	var m entity.AcademyMaterial
	if err := r.db.WithContext(ctx).Select("id, contents_count").First(&m, "id = ?", materialId).Error; err != nil {
		return err
	}

	type aggResult struct {
		AccountId uuid.UUID
		Count     int64
	}
	var aggResults []aggResult

	// FIX: Gunakan Model() bukan Table() string hardcode
	if err := r.db.WithContext(ctx).Model(&entity.AcademyContentProgress{}).
		Select("account_id, count(*) as count").
		Where("material_id = ? AND status = ?", materialId, entity.StatusFinished).
		Group("account_id").
		Scan(&aggResults).Error; err != nil {
		return err
	}

	completedMap := make(map[uuid.UUID]int64)
	for _, res := range aggResults {
		completedMap[res.AccountId] = res.Count
	}

	var progresses []entity.AcademyMaterialProgress
	if err := r.db.WithContext(ctx).Where("material_id = ?", materialId).Find(&progresses).Error; err != nil {
		return err
	}

	for i, p := range progresses {
		completed := completedMap[p.AccountId]
		pct := 0.0
		status := entity.StatusInProgress
		var completedAt *time.Time

		if m.ContentsCount > 0 {
			pct = (float64(completed) / float64(m.ContentsCount)) * 100
			pct = math.Round(pct*100) / 100
			if pct >= 100 {
				pct = 100
				status = entity.StatusFinished
				completedAt = utils.Ptr(time.Now())
			} else if pct <= 0 {
				status = entity.StatusNotStarted
			}
		} else {
			pct = 100
			status = entity.StatusFinished
			completedAt = utils.Ptr(time.Now())
		}

		progresses[i].TotalCompletedContents = uint(completed)
		progresses[i].Progress = pct
		progresses[i].Status = status
		progresses[i].CompletedAt = completedAt
	}

	if len(progresses) > 0 {
		return r.db.WithContext(ctx).Save(&progresses).Error
	}
	return nil
}

func (r *academyRepository) BatchRecalculateAcademyProgress(ctx context.Context, academyId uuid.UUID) error {
	var a entity.Academy
	if err := r.db.WithContext(ctx).Select("id, materials_count").First(&a, "id = ?", academyId).Error; err != nil {
		return err
	}

	type aggResult struct {
		AccountId     uuid.UUID
		TotalProgress float64
		CompletedMats int64
	}
	var aggResults []aggResult

	// FIX: Gunakan Model() bukan Table() string hardcode
	if err := r.db.WithContext(ctx).Model(&entity.AcademyMaterialProgress{}).
		Select("account_id, COALESCE(SUM(progress), 0) as total_progress, COUNT(CASE WHEN status = ? THEN 1 END) as completed_mats", entity.StatusFinished).
		Where("academy_id = ?", academyId).
		Group("account_id").
		Scan(&aggResults).Error; err != nil {
		return err
	}

	dataMap := make(map[uuid.UUID]aggResult)
	for _, res := range aggResults {
		dataMap[res.AccountId] = res
	}

	var progresses []entity.AcademyProgress
	if err := r.db.WithContext(ctx).Where("academy_id = ?", academyId).Find(&progresses).Error; err != nil {
		return err
	}

	for i, p := range progresses {
		data := dataMap[p.AccountId]

		pct := 0.0
		status := entity.StatusInProgress
		var completedAt *time.Time

		if a.MaterialsCount > 0 {
			pct = data.TotalProgress / float64(a.MaterialsCount)
			pct = math.Round(pct*100) / 100

			if pct >= 100 {
				pct = 100
				status = entity.StatusFinished
				completedAt = utils.Ptr(time.Now())
			} else if pct <= 0 {
				status = entity.StatusNotStarted
			}
		} else {
			pct = 100
			status = entity.StatusFinished
			completedAt = utils.Ptr(time.Now())
		}

		progresses[i].TotalCompletedMaterials = uint(data.CompletedMats)
		progresses[i].Progress = pct
		progresses[i].Status = status
		progresses[i].CompletedAt = completedAt
	}

	if len(progresses) > 0 {
		return r.db.WithContext(ctx).Save(&progresses).Error
	}
	return nil
}

func (r *academyRepository) IsAccountAssignedToAcademy(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.AcademyAssign{}).
		Where("account_id = ? AND academy_id = ?", accountId, academyId).
		Count(&count).Error
	return count > 0, err
}

func (r *academyRepository) AssignAccountToAcademy(ctx context.Context, assign entity.AcademyAssign) (entity.AcademyAssign, error) {
	if err := r.db.WithContext(ctx).Create(&assign).Error; err != nil {
		return entity.AcademyAssign{}, err
	}
	return assign, nil
}

func (r *academyRepository) GetAssignmentByAcademyAndAccount(ctx context.Context, academyId uuid.UUID, accountId uuid.UUID) (entity.AcademyAssign, error) {
	var rec entity.AcademyAssign
	if err := r.db.WithContext(ctx).
		Where("academy_id = ? AND account_id = ?", academyId, accountId).
		First(&rec).Error; err != nil {
		return entity.AcademyAssign{}, err
	}
	return rec, nil
}

func (r *academyRepository) ListAssignmentsByAcademy(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyAssign, error) {
	var list []entity.AcademyAssign
	if err := r.db.WithContext(ctx).Where("academy_id = ?", academyId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *academyRepository) RemoveAssignment(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AcademyAssign{}, "id = ?", id).Error
}
