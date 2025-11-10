package services

import (
	"context"
	"errors"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamService interface {
	ListExamByEvent(ctx context.Context, eventSlug string, accountId uuid.UUID) ([]entity.Exam, error)
	GetEventExamExisting(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (ev dto.EventDetailResponse, exam entity.Exam, err error)
	GetExamEventAttempt(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.ExamEventAttempt, error)
	AttemptExamEvent(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (entity.ExamEventAttempt, error)
	SetupQuestions(ctx context.Context, eventSlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error)
	SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.ExamEventAnswer, error)
	SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time)
	SubmitExamEvent(ctx context.Context, attemptId uuid.UUID) (result entity.Result, err error)
	AnswerExamEvent(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error)
}
type evaluator func(answer []string) (float32, entity.CPQuestionVerdict)

type examService struct {
	eventService             EventService
	problemSetService        ProblemSetService
	problemSetExamAssignRepo repositories.ProblemSetExamAssignRepository
	examRepo                 repositories.ExamRepository
	examEventAttemptRepo     repositories.ExamEventAttemptRepository
	examEventAnswerRepo      repositories.ExamEventAnswerRepository
	resultRepo               repositories.ResultRepository
}

func NewExamService(eventService EventService, problemSetService ProblemSetService, problemSetExamAssignRepo repositories.ProblemSetExamAssignRepository, examRepo repositories.ExamRepository, examEventAttemptRepo repositories.ExamEventAttemptRepository, examEventAnswerRepo repositories.ExamEventAnswerRepository, resultRepo repositories.ResultRepository) ExamService {
	return &examService{
		eventService:             eventService,
		problemSetService:        problemSetService,
		problemSetExamAssignRepo: problemSetExamAssignRepo,
		examRepo:                 examRepo,
		examEventAttemptRepo:     examEventAttemptRepo,
		examEventAnswerRepo:      examEventAnswerRepo,
		resultRepo:               resultRepo,
	}
}

func ProtectExamEventAttempt(attempt entity.ExamEventAttempt) entity.ExamEventAttempt {

	var cleanQuestions []entity.Questions
	for _, q := range *attempt.Questions {
		qCopy := q
		qCopy.AnsKey = nil // hide answer key
		qCopy.CorrMark = 0
		qCopy.IncorrMark = 0
		qCopy.NullMark = 0

		cleanQuestions = append(cleanQuestions, qCopy)
	}
	attempt.Questions = &cleanQuestions

	// protect answers verdict info
	var cleanAnswers []entity.ExamEventAnswer

	for _, a := range *attempt.Answers {
		aCopy := a
		aCopy.Score = 0 // hide score

		cleanAnswers = append(cleanAnswers, aCopy)
	}

	attempt.Answers = &cleanAnswers

	return attempt
}

func (s *examService) ListExamByEvent(ctx context.Context, eventSlug string, accountId uuid.UUID) ([]entity.Exam, error) {
	ev, err := s.eventService.DetailBySlug(ctx, eventSlug, accountId)

	if err != nil {
		return []entity.Exam{}, err
	}

	exams, err := s.examRepo.ListByEvent(ctx, ev.Data.Id)

	if err != nil {
		return []entity.Exam{}, err
	}

	return exams, err

}
func (s *examService) GetEventExamExisting(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (ev dto.EventDetailResponse, exam entity.Exam, err error) {
	if ev, err = s.eventService.DetailBySlug(ctx, eventSlug, accountId); err != nil {
		return ev, exam, err
	}

	if exam, err = s.examRepo.GetBySlug(ctx, eventSlug); err != nil {
		return ev, exam, err
	}

	return ev, exam, err
}

func (s *examService) GetExamEventAttempt(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.ExamEventAttempt, error) {

	ev, exam, err := s.GetEventExamExisting(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return dto.UserExamStatus{}, entity.ExamEventAttempt{}, err
	}

	examEventAttempt, err := s.examEventAttemptRepo.GetByExamEvent(ctx, ev.Data.Id, exam.Id, accountId)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserExamStatus{}, entity.ExamEventAttempt{}, err
	}

	var attemptStatus dto.UserExamStatus
	attemptStatus.IsNotAttempt = errors.Is(err, gorm.ErrRecordNotFound)
	attemptStatus.IsTimeOut = (utils.CalculateRemainingTime(examEventAttempt.CreatedAt, examEventAttempt.DueAt) == 0) || false
	attemptStatus.IsSubmitted = examEventAttempt.Submitted
	attemptStatus.IsOnAttempt = !attemptStatus.IsNotAttempt && !attemptStatus.IsTimeOut && !attemptStatus.IsSubmitted

	return attemptStatus, examEventAttempt, err
}
func (s *examService) SetupQuestions(ctx context.Context, eventSlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error) {
	examAssign, err := s.problemSetExamAssignRepo.GetByExam(ctx, examId)

	if err != nil {
		return []entity.Questions{}, err
	}
	problemSet, err := s.problemSetService.GetProblemSet(ctx, examAssign.ProblemSetId)

	if err != nil {
		return []entity.Questions{}, err
	}

	questions, err := s.problemSetService.ListQuestions(ctx, problemSet.Id)

	if err != nil {
		return []entity.Questions{}, err
	}

	return questions, err
}

func (s *examService) SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.ExamEventAnswer, error) {
	var examEventAnswers []entity.ExamEventAnswer
	for _, q := range questions {

		blank_ans := entity.ExamEventAnswer{
			AttemptId:  attemptId,
			QuestionId: q.Id,
		}

		err := s.examEventAnswerRepo.Create(ctx, blank_ans)
		if err != nil {
			return []entity.ExamEventAnswer{}, err
		}
		examEventAnswers = append(examEventAnswers, blank_ans)
	}

	return examEventAnswers, nil
}

func (s *examService) SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time) {
	startTime := time.Now()
	dueTime := startTime.Add(exam.Duration)
	return startTime, dueTime
}

func (s *examService) SubmitExamEvent(ctx context.Context, attemptId uuid.UUID) (result entity.Result, err error) {
	attempt, err := s.examEventAttemptRepo.GetById(ctx, attemptId)

	if err != nil {
		return entity.Result{}, err
	}

	for _, ans := range *attempt.Answers {
		result.FinalScore += ans.Score
	}

	result.ExamEventAttemptId = attempt.Id
	result.ExamEventAttempt = &attempt

	if !attempt.Submitted {
		err := s.resultRepo.Create(ctx, result)
		if err != nil {
			return entity.Result{}, err
		}
	} else {
		err := s.resultRepo.Update(ctx, result)
		if err != nil {
			return entity.Result{}, err
		}
	}
	return result, err

}
func (s *examService) AttemptExamEvent(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (entity.ExamEventAttempt, error) {

	ev, exam, err := s.GetEventExamExisting(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return entity.ExamEventAttempt{}, err
	}
	attemptStatus, examEventAttempt, err := s.GetExamEventAttempt(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return entity.ExamEventAttempt{}, err
	}

	if attemptStatus.IsNotAttempt {
		questions, err := s.SetupQuestions(ctx, eventSlug, exam.Id, accountId)

		if err != nil {
			return entity.ExamEventAttempt{}, err
		}

		answers, err := s.SetupAnswer(ctx, questions, examEventAttempt.Id)

		if err != nil {
			return entity.ExamEventAttempt{}, err
		}

		startTime, dueTime := s.SetupExamTimer(ctx, exam)
		examEventAttempt = entity.ExamEventAttempt{
			AccountId: accountId,
			EventId:   ev.Data.Id,
			ExamId:    exam.Id,
			Questions: &questions,
			Answers:   &answers,
			CreatedAt: startTime,
			DueAt:     dueTime,
			Submitted: false,
		}

		if err := s.examEventAttemptRepo.Create(ctx, examEventAttempt); err != nil {
			return entity.ExamEventAttempt{}, err
		}

		return ProtectExamEventAttempt(examEventAttempt), err

	} else if attemptStatus.IsOnAttempt {

		remTime := utils.CalculateRemainingTime(examEventAttempt.CreatedAt, examEventAttempt.DueAt)
		examEventAttempt.RemTime = remTime

		if err := s.examEventAttemptRepo.Update(ctx, examEventAttempt); err != nil {
			return entity.ExamEventAttempt{}, err
		}

		return ProtectExamEventAttempt(examEventAttempt), err

	} else if attemptStatus.IsTimeOut {
		if examEventAttempt.RemTime != 0 {
			remTime := 0
			examEventAttempt.RemTime = remTime
		}

		s.SubmitExamEvent(ctx, examEventAttempt.Id)
		examEventAttempt.Submitted = true
		if err := s.examEventAttemptRepo.Update(ctx, examEventAttempt); err != nil {
			return entity.ExamEventAttempt{}, err
		}

	} else if attemptStatus.IsSubmitted {
		return examEventAttempt, nil
	}
	return entity.ExamEventAttempt{}, http_error.INTERNAL_SERVER_ERROR
}

func (s *examService) EvaluateAnswer(ctx context.Context, question entity.Questions) evaluator {

	nonCPEvaluator := func(answer []string) (float32, entity.CPQuestionVerdict) {
		score := float32(0)
		for i, ans := range answer {
			if ans != question.AnsKey[i] && ans != "" {
				score += float32(question.IncorrMark)
			} else if ans == "" {
				score += float32(question.NullMark)
			} else {
				score += float32(question.CorrMark)
			}
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
func (s *examService) AnswerExamEvent(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error) {

	question, err := s.problemSetService.GetQuestionById(ctx, questionId)
	if err != nil {
		return entity.CPQuestionVerdict{}, err
	}

	score, CPQuestionVerdict := s.EvaluateAnswer(ctx, question)(answer)

	err = s.examEventAnswerRepo.Update(ctx, entity.ExamEventAnswer{
		AttemptId:  attemptId,
		QuestionId: questionId,
		Answers:    answer,
		Score:      score,
	})

	return CPQuestionVerdict, err
}
