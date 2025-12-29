package services

import (
	"context"
	"mime/multipart"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
)

type EventExamProctoringService interface {
	CreateLog(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID, req dto.EventExamProctoringLogsRequest, file *multipart.FileHeader) error
	ListLogs(ctx context.Context, accountId uuid.UUID, examId uuid.UUID, eventId uuid.UUID) ([]entity.EventExamProctoringLogs, error)
	GetLogById(ctx context.Context, id uuid.UUID) (entity.EventExamProctoringLogs, error)
	UpdateLog(ctx context.Context, id uuid.UUID, req dto.EventExamProctoringLogsRequest, file *multipart.FileHeader) error
	DeleteLog(ctx context.Context, id uuid.UUID) error
}

type eventExamProctoringService struct {
	eventExamService EventExamService
	uploadService    UploadService
	repo             repositories.EventExamProctoringRepository
}

func NewEventExamProctoringService(eventExamService EventExamService, uploadService UploadService, repo repositories.EventExamProctoringRepository) EventExamProctoringService {
	return &eventExamProctoringService{
		eventExamService: eventExamService,
		uploadService:    uploadService,
		repo:             repo,
	}
}

func (s *eventExamProctoringService) CreateLog(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID, req dto.EventExamProctoringLogsRequest, file *multipart.FileHeader) error {
	_, attempt, err := s.eventExamService.GetEventExamAttempt(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return err
	}

	var attachmentUrl string
	if file != nil {
		files, err := s.uploadService.UploadFiles(ctx, []*multipart.FileHeader{file}, "submission", accountId)
		if err != nil {
			return err
		}
		if len(files) > 0 {
			attachmentUrl = files[0].Path
		}
	}

	log := entity.EventExamProctoringLogs{
		Id:                uuid.New(),
		EventId:           attempt.EventId,
		ExamId:            attempt.ExamId,
		AccountId:         accountId,
		ViolationScore:    req.ViolationScore,
		ViolationCategory: req.ViolationCategory,
		Attachement:       attachmentUrl,
		CreatedAt:         time.Now(),
	}

	return s.repo.Create(ctx, &log)
}

func (s *eventExamProctoringService) ListLogs(ctx context.Context, accountId uuid.UUID, examId uuid.UUID, eventId uuid.UUID) ([]entity.EventExamProctoringLogs, error) {
	return s.repo.List(ctx, accountId, examId, eventId)
}

func (s *eventExamProctoringService) GetLogById(ctx context.Context, id uuid.UUID) (entity.EventExamProctoringLogs, error) {
	return s.repo.GetById(ctx, id)
}

func (s *eventExamProctoringService) UpdateLog(ctx context.Context, id uuid.UUID, req dto.EventExamProctoringLogsRequest, file *multipart.FileHeader) error {
	log, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	var attachmentUrl = log.Attachement
	if file != nil {
		files, err := s.uploadService.UploadFiles(ctx, []*multipart.FileHeader{file}, "submission", log.AccountId)
		if err != nil {
			return err
		}
		if len(files) > 0 {
			attachmentUrl = files[0].Path
		}
	}

	// Update fields if they are provided (for non-zero values) or logic requires
	// Here I assume we update what's in request.
	if req.ViolationScore != 0 {
		log.ViolationScore = req.ViolationScore
	}
	if req.ViolationCategory != "" {
		log.ViolationCategory = req.ViolationCategory
	}
	log.Attachement = attachmentUrl

	return s.repo.Update(ctx, &log)
}

func (s *eventExamProctoringService) DeleteLog(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
