package services

import (
	"context"
	"errors"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
)

var (
	ErrProblemSetNotFound = errors.New("problem set not found")
	ErrQuestionNotFound   = errors.New("question not found")
)

type ProblemSetService interface {
	CreateProblemSet(ctx context.Context, ps entity.ProblemSet) error
	GetProblemSet(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error)
	ListProblemSets(ctx context.Context) ([]entity.ProblemSet, error)
	UpdateProblemSet(ctx context.Context, ps entity.ProblemSet) error
	DeleteProblemSet(ctx context.Context, id uuid.UUID) error

	AddQuestion(ctx context.Context, q entity.Questions) error
	UpdateQuestion(ctx context.Context, q entity.Questions) error
	DeleteQuestion(ctx context.Context, qID uuid.UUID) error
	ListQuestions(ctx context.Context, psID uuid.UUID) ([]entity.Questions, error)

	AssignProblemSetToExam(ctx context.Context, examId uuid.UUID, problemSetId uuid.UUID) error
	RemoveAssignedProblemSet(ctx context.Context, assignId uuid.UUID) error
	GetQuestionById(ctx context.Context, qID uuid.UUID) (entity.Questions, error)
	ListQuestionsByExam(ctx context.Context, examId uuid.UUID) ([]entity.Questions, error)
}

type problemSetService struct {
	problemSetRepository           repositories.ProblemSetRepository
	questionsRepository            repositories.QuestionsRepository
	problemSetExamAssignRepository repositories.ProblemSetExamAssignRepository
}

func NewProblemSetService(
	problemSetRepository repositories.ProblemSetRepository,
	questionsRepository repositories.QuestionsRepository,
	problemSetExamAssignRepository repositories.ProblemSetExamAssignRepository,
) ProblemSetService {
	return &problemSetService{
		problemSetRepository:           problemSetRepository,
		questionsRepository:            questionsRepository,
		problemSetExamAssignRepository: problemSetExamAssignRepository,
	}
}

// ---------------- Problem Set CRUD ----------------

func (s *problemSetService) CreateProblemSet(ctx context.Context, ps entity.ProblemSet) error {
	return s.problemSetRepository.Create(ctx, ps)
}

func (s *problemSetService) GetProblemSet(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error) {
	ps, err := s.problemSetRepository.Get(ctx, id)
	if err != nil {
		return entity.ProblemSet{}, ErrProblemSetNotFound
	}
	return ps, nil
}

func (s *problemSetService) ListProblemSets(ctx context.Context) ([]entity.ProblemSet, error) {
	return s.problemSetRepository.List(ctx)
}

func (s *problemSetService) UpdateProblemSet(ctx context.Context, ps entity.ProblemSet) error {
	_, err := s.problemSetRepository.Get(ctx, ps.Id)
	if err != nil {
		return ErrProblemSetNotFound
	}
	return s.problemSetRepository.Update(ctx, ps)
}

func (s *problemSetService) DeleteProblemSet(ctx context.Context, id uuid.UUID) error {
	_, err := s.problemSetRepository.Get(ctx, id)
	if err != nil {
		return ErrProblemSetNotFound
	}
	return s.problemSetRepository.Delete(ctx, id)
}

// ---------------- Questions ----------------

func (s *problemSetService) AddQuestion(ctx context.Context, q entity.Questions) error {
	_, err := s.problemSetRepository.Get(ctx, q.ProblemSetId)
	if err != nil {
		return ErrProblemSetNotFound
	}
	return s.questionsRepository.Create(ctx, q)
}

func (s *problemSetService) UpdateQuestion(ctx context.Context, q entity.Questions) error {
	_, err := s.questionsRepository.Get(ctx, q.Id)
	if err != nil {
		return ErrQuestionNotFound
	}
	return s.questionsRepository.Update(ctx, q)
}

func (s *problemSetService) DeleteQuestion(ctx context.Context, qID uuid.UUID) error {
	_, err := s.questionsRepository.Get(ctx, qID)
	if err != nil {
		return ErrQuestionNotFound
	}
	return s.questionsRepository.Delete(ctx, qID)
}

func (s *problemSetService) ListQuestions(ctx context.Context, psID uuid.UUID) ([]entity.Questions, error) {
	_, err := s.problemSetRepository.Get(ctx, psID)
	if err != nil {
		return nil, ErrProblemSetNotFound
	}
	return s.questionsRepository.ListByProblemSet(ctx, psID)
}

// ---------------- Exam ↔ Problem Set (Mapping Table) ----------------

func (s *problemSetService) AssignProblemSetToExam(ctx context.Context, examId uuid.UUID, problemSetId uuid.UUID) error {
	_, err := s.problemSetRepository.Get(ctx, problemSetId)
	if err != nil {
		return ErrProblemSetNotFound
	}

	assign := entity.ProblemSetExamAssign{
		Id:           uuid.New(),
		ExamId:       examId,
		ProblemSetId: problemSetId,
	}

	return s.problemSetExamAssignRepository.Create(ctx, assign)
}

func (s *problemSetService) RemoveAssignedProblemSet(ctx context.Context, assignId uuid.UUID) error {
	return s.problemSetExamAssignRepository.Delete(ctx, assignId)
}

func (s *problemSetService) GetQuestionById(ctx context.Context, qID uuid.UUID) (entity.Questions, error) {
	question, err := s.questionsRepository.Get(ctx, qID)
	if err != nil {
		return entity.Questions{}, err
	}
	return question, err
}

func (s *problemSetService) ListQuestionsByExam(ctx context.Context, examId uuid.UUID) ([]entity.Questions, error) {
	assign, err := s.problemSetExamAssignRepository.GetByExam(ctx, examId)
	if err != nil {
		return []entity.Questions{}, err
	}
	return s.questionsRepository.ListByProblemSet(ctx, assign.ProblemSetId)
}

func (s *problemSetService) ListProblemSets(ctx context.Context) ([]entity.ProblemSet, error) {
	return s.problemSetRepository.List(ctx)
}

func (s *problemSetService) UpdateProblemSet(ctx context.Context, ps entity.ProblemSet) error {
	_, err := s.problemSetRepository.Get(ctx, ps.Id)
	if err != nil {
		return ErrProblemSetNotFound
	}
	return s.problemSetRepository.Update(ctx, ps)
}

func (s *problemSetService) DeleteProblemSet(ctx context.Context, id uuid.UUID) error {
	_, err := s.problemSetRepository.Get(ctx, id)
	if err != nil {
		return ErrProblemSetNotFound
	}
	return s.problemSetRepository.Delete(ctx, id)
}

// ---------------- Questions ----------------

func (s *problemSetService) AddQuestion(ctx context.Context, q entity.Questions) error {
	_, err := s.problemSetRepository.Get(ctx, q.ProblemSetId)
	if err != nil {
		return ErrProblemSetNotFound
	}
	return s.questionsRepository.Create(ctx, q)
}

func (s *problemSetService) UpdateQuestion(ctx context.Context, q entity.Questions) error {
	_, err := s.questionsRepository.Get(ctx, q.Id)
	if err != nil {
		return ErrQuestionNotFound
	}
	return s.questionsRepository.Update(ctx, q)
}

func (s *problemSetService) DeleteQuestion(ctx context.Context, qID uuid.UUID) error {
	_, err := s.questionsRepository.Get(ctx, qID)
	if err != nil {
		return ErrQuestionNotFound
	}
	return s.questionsRepository.Delete(ctx, qID)
}

func (s *problemSetService) ListQuestions(ctx context.Context, psID uuid.UUID) ([]entity.Questions, error) {
	_, err := s.problemSetRepository.Get(ctx, psID)
	if err != nil {
		return nil, ErrProblemSetNotFound
	}
	return s.questionsRepository.ListByProblemSet(ctx, psID)
}

// ---------------- Exam ↔ Problem Set (Mapping Table) ----------------

func (s *problemSetService) AssignProblemSetToExam(ctx context.Context, examId uuid.UUID, problemSetId uuid.UUID) error {
	_, err := s.problemSetRepository.Get(ctx, problemSetId)
	if err != nil {
		return ErrProblemSetNotFound
	}

	assign := entity.ProblemSetExamAssign{
		Id:           uuid.New(),
		ExamId:       examId,
		ProblemSetId: problemSetId,
	}

	return s.problemSetExamAssignRepository.Create(ctx, assign)
}

func (s *problemSetService) RemoveAssignedProblemSet(ctx context.Context, assignId uuid.UUID) error {
	return s.problemSetExamAssignRepository.Delete(ctx, assignId)
}

func (s *problemSetService) GetQuestionById(ctx context.Context, qID uuid.UUID) (entity.Questions, error) {
	question, err := s.questionsRepository.Get(ctx, qID)
	if err != nil {
		return entity.Questions{}, err
	}
	return question, err
}

func (s *problemSetService) ListQuestionsByExam(ctx context.Context, examId uuid.UUID) ([]entity.Questions, error) {
	assign, err := s.problemSetExamAssignRepository.GetByExam(ctx, examId)
	if err != nil {
		return []entity.Questions{}, err
	}
	return s.questionsRepository.ListByProblemSet(ctx, assign.ProblemSetId)
}
