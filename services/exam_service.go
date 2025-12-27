package services

import (
	"context"

	dto "abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
)

type ExamService interface {
	CreateExam(ctx context.Context, req dto.CreateExamRequest) (entity.Exam, error)
	UpdateExam(ctx context.Context, id uuid.UUID, req dto.CreateExamRequest) (entity.Exam, error)
	DeleteExam(ctx context.Context, id uuid.UUID) error
	GetExamList(ctx context.Context, p entity.Pagination) ([]entity.Exam, int64, error)
	GetExamDetail(ctx context.Context, id uuid.UUID) (entity.Exam, error)
	AssignExamToEvent(ctx context.Context, examId uuid.UUID, eventId uuid.UUID) error
	AssignExamToAcademy(ctx context.Context, examId uuid.UUID, academyId uuid.UUID) error
}

type examService struct {
	examRepo              repositories.ExamRepository
	eventExamAssignRepo   repositories.EventExamAssignRepository
	academyExamAssignRepo repositories.AcademyExamAssignRepository
}

func NewExamService(
	examRepo repositories.ExamRepository,
	eventExamAssignRepo repositories.EventExamAssignRepository,
	academyExamAssignRepo repositories.AcademyExamAssignRepository,
) ExamService {
	return &examService{
		examRepo:              examRepo,
		eventExamAssignRepo:   eventExamAssignRepo,
		academyExamAssignRepo: academyExamAssignRepo,
	}
}

func (s *examService) CreateExam(ctx context.Context, req dto.CreateExamRequest) (entity.Exam, error) {

	exam := entity.Exam{
		Slug:        req.Slug,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Randomize:   req.Randomize,
		Configuration: entity.ExamConfiguration{
			AllowRetake: req.AllowRetake,
			AllowReview: req.AllowReview,
			EnableTimer: req.EnableTimer,
		},
		Proctoring: entity.ExamProctoring{
			EnableWebCam:       req.EnableWebCam,
			EnableVAD:          req.EnableVAD,
			EnableTabBlock:     req.EnableTabBlock,
			RequiredFullScreen: req.RequiredFullScreen,
			EnableEyeTracking:  req.EnableEyeTracking,
			DisableCopyPaste:   req.DisableCopyPaste,
			EnableExamBrowser:  req.EnableExamBrowser,
		},
	}

	if err := s.examRepo.Create(ctx, &exam); err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}

func (s *examService) UpdateExam(ctx context.Context, id uuid.UUID, req dto.CreateExamRequest) (entity.Exam, error) {
	exam, err := s.examRepo.Get(ctx, id)
	if err != nil {
		return entity.Exam{}, err
	}

	exam.Slug = req.Slug
	exam.Title = req.Title
	exam.Description = req.Description
	exam.Duration = req.Duration
	exam.Randomize = req.Randomize

	// Update Configuration
	exam.Configuration.AllowRetake = req.AllowRetake
	exam.Configuration.AllowReview = req.AllowReview
	exam.Configuration.EnableTimer = req.EnableTimer

	// Update Proctoring
	exam.Proctoring.EnableWebCam = req.EnableWebCam
	exam.Proctoring.EnableVAD = req.EnableVAD
	exam.Proctoring.EnableTabBlock = req.EnableTabBlock
	exam.Proctoring.RequiredFullScreen = req.RequiredFullScreen
	exam.Proctoring.EnableEyeTracking = req.EnableEyeTracking
	exam.Proctoring.DisableCopyPaste = req.DisableCopyPaste
	exam.Proctoring.EnableExamBrowser = req.EnableExamBrowser

	if err := s.examRepo.Update(ctx, exam); err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}

func (s *examService) DeleteExam(ctx context.Context, id uuid.UUID) error {
	return s.examRepo.Delete(ctx, id)
}

func (s *examService) GetExamList(ctx context.Context, p entity.Pagination) ([]entity.Exam, int64, error) {
	return s.examRepo.ListWithPagination(ctx, p)
}

func (s *examService) GetExamDetail(ctx context.Context, id uuid.UUID) (entity.Exam, error) {
	return s.examRepo.Get(ctx, id)
}

func (s *examService) AssignExamToEvent(ctx context.Context, examId uuid.UUID, eventId uuid.UUID) error {
	// Check if already assigned
	if err := s.eventExamAssignRepo.Check(ctx, eventId, examId); err == nil {
		return nil // Already assigned
	}

	assign := entity.EventExamAssign{
		ExamId:  examId,
		EventId: eventId,
	}

	return s.eventExamAssignRepo.Create(ctx, assign)
}

func (s *examService) AssignExamToAcademy(ctx context.Context, examId uuid.UUID, academyId uuid.UUID) error {
	// Check if already assigned
	if err := s.academyExamAssignRepo.Check(ctx, academyId, examId); err == nil {
		return nil // Already assigned
	}

	assign := entity.AcademyExamAssign{
		ExamId:    examId,
		AcademyId: academyId,
	}

	return s.academyExamAssignRepo.Create(ctx, assign)
}
