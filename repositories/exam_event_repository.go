package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamRepository interface {
	Create(ctx context.Context, e entity.Exam) error
	Get(ctx context.Context, id uuid.UUID) (entity.Exam, error)
	GetBySlug(ctx context.Context, slug string) (entity.Exam, error)
	Update(ctx context.Context, e entity.Exam) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]entity.Exam, error)

	// Additional business need
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Exam, error)
}

type examRepository struct {
	db *gorm.DB
}

func NewExamRepository(db *gorm.DB) ExamRepository {
	return &examRepository{db}
}

// ================= CRUD =================

func (r *examRepository) Create(ctx context.Context, e entity.Exam) error {
	return r.db.WithContext(ctx).Create(&e).Error
}

func (r *examRepository) Get(ctx context.Context, id uuid.UUID) (entity.Exam, error) {
	var e entity.Exam
	err := r.db.WithContext(ctx).
		First(&e, "id_exam = ?", id).Error
	return e, err
}

func (r *examRepository) GetBySlug(ctx context.Context, slug string) (entity.Exam, error) {
	var e entity.Exam
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&e).Error
	return e, err
}

func (r *examRepository) Update(ctx context.Context, e entity.Exam) error {
	return r.db.WithContext(ctx).
		Model(&entity.Exam{}).
		Where("id_exam = ?", e.Id).
		Updates(e).Error
}

func (r *examRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id_exam = ?", id).
		Delete(&entity.Exam{}).Error
}

func (r *examRepository) List(ctx context.Context) ([]entity.Exam, error) {
	var list []entity.Exam
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// =========== Business Specific ============

func (r *examRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Exam, error) {
	var exams []entity.Exam

	err := r.db.WithContext(ctx).
		Table("exam").
		Joins(`JOIN exam_event_assign ON exam_event_assign.id_exam = exam.id_exam`).
		Where("exam_event_assign.id_event = ?", eventId).
		Find(&exams).Error

	return exams, err
}
