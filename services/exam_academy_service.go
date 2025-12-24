package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyExamService interface {
	ListExamByAcademy(ctx context.Context, academySlug string, accountId uuid.UUID) ([]entity.Exam, error)
	GetAcademyExamExisting(ctx context.Context, academySlug string, examSlug string, accountId uuid.UUID) (entity.Academy, entity.Exam, error)
	GetExamAcademyAttempt(ctx context.Context, academySlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.ExamAcademyAttempt, error)
	AttemptExamAcademy(ctx context.Context, academySlug string, examSlug string, accountId uuid.UUID) (entity.ExamAcademyAttempt, error)
	SetupQuestions(ctx context.Context, academySlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error)
	SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.ExamAcademyAnswer, error)
	SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time)
	SubmitExamAcademy(ctx context.Context, attemptId uuid.UUID) (entity.ExamAcademyResult, error)
	AnswerExamAcademy(ctx context.Context, academySlug string, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error)
}

type academyExamService struct {
	academyService         AcademyService
	problemSetService      ProblemSetService
	examRepo               repositories.ExamRepository
	examAcademyAttemptRepo repositories.ExamAcademyAttemptRepository
	examAcademyAnswerRepo  repositories.ExamAcademyAnswerRepository
	examAcademyAssignRepo  repositories.ExamAcademyAssignRepository
	academyResultRepo      repositories.AcademyResultRepository
}

func NewAcademyExamService(academyService AcademyService, problemSetService ProblemSetService, examRepo repositories.ExamRepository, examAcademyAttemptRepo repositories.ExamAcademyAttemptRepository, examAcademyAssignRepo repositories.ExamAcademyAssignRepository, examAcademyAnswerRepo repositories.ExamAcademyAnswerRepository, academyResultRepo repositories.AcademyResultRepository) AcademyExamService {
	return &academyExamService{
		academyService:         academyService,
		problemSetService:      problemSetService,
		examRepo:               examRepo,
		examAcademyAttemptRepo: examAcademyAttemptRepo,
		examAcademyAssignRepo:  examAcademyAssignRepo,
		examAcademyAnswerRepo:  examAcademyAnswerRepo,
		academyResultRepo:      academyResultRepo,
	}
}

func (s *academyExamService) ListExamByAcademy(ctx context.Context, academySlug string, accountId uuid.UUID) ([]entity.Exam, error) {
	academy, err := s.academyService.GetAcademy(ctx, accountId, academySlug)
	if err != nil {
		return []entity.Exam{}, err
	}
	assigns, err := s.examAcademyAssignRepo.ListByAcademy(ctx, academy.Id)
	if err != nil {
		return []entity.Exam{}, err
	}
	var exams []entity.Exam
	for _, a := range assigns {
		exams = append(exams, *a.Exam)
	}
	return exams, nil
}

func (s *academyExamService) GetAcademyExamExisting(ctx context.Context, academySlug string, examSlug string, accountId uuid.UUID) (entity.Academy, entity.Exam, error) {
	academy, err := s.academyService.GetAcademy(ctx, accountId, academySlug)
	if err != nil {
		return entity.Academy{}, entity.Exam{}, err
	}
	exam, err := s.examRepo.GetBySlug(ctx, examSlug)
	if err != nil {
		return entity.Academy{}, entity.Exam{}, err
	}
	if err := s.examAcademyAssignRepo.Check(ctx, academy.Id, exam.Id); err != nil {
		return entity.Academy{}, entity.Exam{}, err
	}
	return academy, exam, nil
}

func (s *academyExamService) GetExamAcademyAttempt(ctx context.Context, academySlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.ExamAcademyAttempt, error) {
	academy, exam, err := s.GetAcademyExamExisting(ctx, academySlug, examSlug, accountId)
	if err != nil {
		return dto.UserExamStatus{}, entity.ExamAcademyAttempt{}, err
	}
	attempt, err := s.examAcademyAttemptRepo.GetByExamAcademy(ctx, academy.Id, exam.Id, accountId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserExamStatus{}, entity.ExamAcademyAttempt{}, err
	}
	var status dto.UserExamStatus
	status.IsNotAttempt = errors.Is(err, gorm.ErrRecordNotFound)
	status.IsTimeOut = (utils.CalculateRemainingTime(attempt.CreatedAt, attempt.DueAt) == 0) || false
	status.IsSubmitted = attempt.Submitted
	status.IsOnAttempt = !status.IsNotAttempt && !status.IsTimeOut && !status.IsSubmitted
	return status, attempt, nil
}

func (s *academyExamService) AttemptExamAcademy(ctx context.Context, academySlug string, examSlug string, accountId uuid.UUID) (entity.ExamAcademyAttempt, error) {
	academy, exam, err := s.GetAcademyExamExisting(ctx, academySlug, examSlug, accountId)
	if err != nil {
		return entity.ExamAcademyAttempt{}, err
	}
	status, attempt, err := s.GetExamAcademyAttempt(ctx, academySlug, examSlug, accountId)
	if err != nil {
		return entity.ExamAcademyAttempt{}, err
	}
	questions, err := s.SetupQuestions(ctx, academySlug, exam.Id, accountId)
	attempt.Questions = questions
	if err != nil {
		return entity.ExamAcademyAttempt{}, err
	}
	if status.IsNotAttempt {
		startTime, dueTime := s.SetupExamTimer(ctx, exam)
		remTime := utils.CalculateRemainingTime(startTime, dueTime)
		attempt = entity.ExamAcademyAttempt{
			AccountId: accountId,
			AcademyId: academy.Id,
			ExamId:    exam.Id,
			CreatedAt: startTime,
			DueAt:     dueTime,
			Submitted: false,
			RemTime:   remTime,
			Questions: questions,
		}
		if err := s.examAcademyAttemptRepo.Create(ctx, &attempt); err != nil {
			return entity.ExamAcademyAttempt{}, err
		}
		answers, err := s.SetupAnswer(ctx, questions, attempt.Id)
		if err != nil {
			return entity.ExamAcademyAttempt{}, err
		}
		attempt.Answers = answers
		return ProtectExamAcademyAttempt(attempt), nil
	}
	return ProtectExamAcademyAttempt(attempt), nil
}

func (s *academyExamService) SetupQuestions(ctx context.Context, academySlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error) {
	qs, err := s.problemSetService.ListQuestionsByExam(ctx, examId)
	if err != nil {
		return []entity.Questions{}, err
	}
	return qs, nil
}

func (s *academyExamService) SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.ExamAcademyAnswer, error) {
	var answers []entity.ExamAcademyAnswer
	for _, q := range questions {
		ans := entity.ExamAcademyAnswer{Id: uuid.New(), AttemptId: attemptId, QuestionId: q.Id, Score: 0}
		if err := s.examAcademyAnswerRepo.Create(ctx, &ans); err != nil {
			return []entity.ExamAcademyAnswer{}, err
		}
		answers = append(answers, ans)
	}
	return answers, nil
}

func (s *academyExamService) SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time) {
	start := time.Now()
	due := start.Add(exam.Duration * time.Minute)
	return start, due
}

func (s *academyExamService) SubmitExamAcademy(ctx context.Context, attemptId uuid.UUID) (entity.ExamAcademyResult, error) {
	attempt, err := s.examAcademyAttemptRepo.GetById(ctx, attemptId)
	if err != nil {
		return entity.ExamAcademyResult{}, err
	}
	if attempt.Submitted {
		return entity.ExamAcademyResult{}, http_error.EXAMS_SUBMITTED
	}
	answers, err := s.examAcademyAnswerRepo.ListByAttempt(ctx, attemptId)
	if err != nil {
		return entity.ExamAcademyResult{}, err
	}
	var sum float32
	for _, a := range answers {
		sum += a.Score
	}
	rec := entity.ExamAcademyResult{Id: uuid.New(), AttemptId: attemptId, FinalScore: sum}
	if err := s.academyResultRepo.Create(ctx, &rec); err != nil {
		return entity.ExamAcademyResult{}, err
	}
	attempt.Submitted = true
	if err := s.examAcademyAttemptRepo.Update(ctx, &attempt); err != nil {
		return entity.ExamAcademyResult{}, err
	}
	return rec, nil
}

func (s *academyExamService) AnswerExamAcademy(ctx context.Context, academySlug string, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error) {
	attempt, err := s.examAcademyAttemptRepo.GetById(ctx, attemptId)
	if err != nil {
		return entity.CPQuestionVerdict{}, err
	}
	if attempt.Submitted {
		return entity.CPQuestionVerdict{}, http_error.EXAMS_SUBMITTED
	}
	if utils.CalculateRemainingTime(attempt.CreatedAt, attempt.DueAt) == 0 || time.Now().After(attempt.DueAt) {
		return entity.CPQuestionVerdict{}, http_error.EXAMS_TIME_EXCEEDED
	}
	question, err := s.problemSetService.GetQuestionById(ctx, questionId)
	if err != nil {
		return entity.CPQuestionVerdict{}, err
	}
	score, verdict := s.EvaluateAnswer(ctx, question)(answer)
	err = s.examAcademyAnswerRepo.Update(ctx, &entity.ExamAcademyAnswer{AttemptId: attemptId, QuestionId: questionId, Answers: answer, Score: score})
	return verdict, err
}

func (s *academyExamService) EvaluateAnswer(ctx context.Context, question entity.Questions) evaluator {

	nonCPEvaluator := func(answer []string) (float32, entity.CPQuestionVerdict) {
		score := float32(0)
		isCorrect := true
		for i, ans := range answer {
			fmt.Println("User Answer :", ans)
			fmt.Println("Answer Key :", question.AnsKey[i])
			if ans != question.AnsKey[i] && ans != "" {
				score += float32(question.IncorrMark)
				isCorrect = false
				break
			} else if ans == "" {
				score += float32(question.NullMark)
				isCorrect = false
				break
			}
		}

		if isCorrect {
			score += float32(question.CorrMark)
		}

		return score, entity.CPQuestionVerdict{}
	}

	CPEvaluator := func(answer []string) (float32, entity.CPQuestionVerdict) {
		return 0, entity.CPQuestionVerdict{
			TimeExecution: 0.01,
			MemoryUsage:   256.0,
			Verdict:       "AC",
			Score:         100.0,
		}
	}

	var examEvaluator = map[string]evaluator{
		"multiple_choices":         nonCPEvaluator,
		"multiple_choices_complex": nonCPEvaluator,
		"short_answer":             nonCPEvaluator,
		"true_false":               nonCPEvaluator,
		"code_puzzle":              nonCPEvaluator,
		"code_type":                nonCPEvaluator,
		"competitive_programming":  CPEvaluator,
	}

	return examEvaluator[question.Type]
}
func ProtectExamAcademyAttempt(attempt entity.ExamAcademyAttempt) entity.ExamAcademyAttempt {
	var cleanQuestions []entity.Questions
	for _, q := range attempt.Questions {
		qc := q
		qc.AnsKey = nil
		qc.CorrMark = 0
		qc.IncorrMark = 0
		qc.NullMark = 0
		cleanQuestions = append(cleanQuestions, qc)
	}
	attempt.Questions = cleanQuestions
	var cleanAnswers []entity.ExamAcademyAnswer
	for _, a := range attempt.Answers {
		ac := a
		ac.Score = 0
		cleanAnswers = append(cleanAnswers, ac)
	}
	attempt.Answers = cleanAnswers
	return attempt
}
