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

type ExamService interface {
	ListExamByEvent(ctx context.Context, eventSlug string, accountId uuid.UUID) ([]entity.Exam, error)
	GetEventExamExisting(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (ev dto.EventDetailResponse, exam entity.Exam, err error)
	GetExamEventAttempt(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.ExamEventAttempt, error)
	AttemptExamEvent(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (entity.ExamEventAttempt, error)
	SetupQuestions(ctx context.Context, eventSlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error)
	SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.ExamEventAnswer, error)
	SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time)
	SubmitExamEvent(ctx context.Context, attemptId uuid.UUID) (result entity.Result, err error)
	AnswerExamEvent(ctx context.Context, eventSlug string, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error)
}
type evaluator func(answer []string) (float32, entity.CPQuestionVerdict)

type examService struct {
	eventService             EventService
	problemSetService        ProblemSetService
	problemSetExamAssignRepo repositories.ProblemSetExamAssignRepository
	examRepo                 repositories.ExamRepository
	examEventAttemptRepo     repositories.ExamEventAttemptRepository
	examEventAnswerRepo      repositories.ExamEventAnswerRepository
	examEventAssignRepo      repositories.ExamEventAssignRepository
	resultRepo               repositories.ResultRepository
}

func NewExamService(eventService EventService, problemSetService ProblemSetService, problemSetExamAssignRepo repositories.ProblemSetExamAssignRepository, examRepo repositories.ExamRepository, examEventAttemptRepo repositories.ExamEventAttemptRepository, examEventAssignRepo repositories.ExamEventAssignRepository, examEventAnswerRepo repositories.ExamEventAnswerRepository, resultRepo repositories.ResultRepository) ExamService {
	return &examService{
		eventService:             eventService,
		problemSetService:        problemSetService,
		problemSetExamAssignRepo: problemSetExamAssignRepo,
		examRepo:                 examRepo,
		examEventAttemptRepo:     examEventAttemptRepo,
		examEventAssignRepo:      examEventAssignRepo,
		examEventAnswerRepo:      examEventAnswerRepo,
		resultRepo:               resultRepo,
	}
}

func ProtectExamEventAttempt(attempt entity.ExamEventAttempt) entity.ExamEventAttempt {

	var cleanQuestions []entity.Questions
	for _, q := range attempt.Questions {
		qCopy := q
		qCopy.AnsKey = nil // hide answer key
		qCopy.CorrMark = 0
		qCopy.IncorrMark = 0
		qCopy.NullMark = 0

		cleanQuestions = append(cleanQuestions, qCopy)
	}
	attempt.Questions = cleanQuestions

	// protect answers verdict info
	var cleanAnswers []entity.ExamEventAnswer

	for _, a := range attempt.Answers {
		aCopy := a
		aCopy.Score = 0 // hide score

		cleanAnswers = append(cleanAnswers, aCopy)
	}

	attempt.Answers = cleanAnswers

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

	if exam, err = s.examRepo.GetBySlug(ctx, examSlug); err != nil {
		return ev, exam, err
	}

	if err := s.examEventAssignRepo.Check(ctx, ev.Data.Id, exam.Id); err != nil {
		return dto.EventDetailResponse{}, entity.Exam{}, err
	}

	return ev, exam, err
}

func (s *examService) GetExamEventAttempt(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.ExamEventAttempt, error) {

	ev, exam, err := s.GetEventExamExisting(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return dto.UserExamStatus{}, entity.ExamEventAttempt{}, err
	}

	examEventAttempt, err := s.examEventAttemptRepo.GetByExamEvent(ctx, ev.Data.Id, exam.Id, accountId)
	fmt.Println("Error Exam Event Attempt", errors.Is(err, gorm.ErrRecordNotFound))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserExamStatus{}, entity.ExamEventAttempt{}, err
	}

	var attemptStatus dto.UserExamStatus
	attemptStatus.IsNotAttempt = errors.Is(err, gorm.ErrRecordNotFound)
	attemptStatus.IsTimeOut = (utils.CalculateRemainingTime(examEventAttempt.CreatedAt, examEventAttempt.DueAt) == 0) || false
	attemptStatus.IsSubmitted = examEventAttempt.Submitted
	attemptStatus.IsOnAttempt = !attemptStatus.IsNotAttempt && !attemptStatus.IsTimeOut && !attemptStatus.IsSubmitted
	return attemptStatus, examEventAttempt, nil

}
func (s *examService) SetupQuestions(ctx context.Context, eventSlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error) {
	examAssign, err := s.problemSetExamAssignRepo.GetByExam(ctx, examId)

	if err != nil {
		return []entity.Questions{}, err
	}

	questions, err := s.problemSetService.ListQuestions(ctx, examAssign.ProblemSetId)

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

		err := s.examEventAnswerRepo.Create(ctx, &blank_ans)
		if err != nil {
			return []entity.ExamEventAnswer{}, err
		}
		examEventAnswers = append(examEventAnswers, blank_ans)
	}

	return examEventAnswers, nil
}

func (s *examService) SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time) {
	startTime := time.Now()
	dueTime := startTime.Add(exam.Duration * time.Minute)
	return startTime, dueTime
}

func (s *examService) SubmitExamEvent(ctx context.Context, attemptId uuid.UUID) (result entity.Result, err error) {
	attempt, err := s.examEventAttemptRepo.GetById(ctx, attemptId)
	finalScore := float32(0)
	if err != nil {
		return entity.Result{}, err
	}

	for _, ans := range attempt.Answers {
		finalScore += ans.Score
	}

	if !attempt.Submitted {

		attempt.Submitted = true
		result.AttemptId = attempt.Id
		result.ExamEventAttempt = &attempt
		result.FinalScore = float32(finalScore)

		s.examEventAttemptRepo.Update(ctx, &attempt)
		err := s.resultRepo.Create(ctx, &result)

		if err != nil {
			return entity.Result{}, err
		}
	} else {
		result, err = s.resultRepo.GetByAttemptId(ctx, attempt.Id)

		if err != nil {
			return entity.Result{}, err
		}

		result.FinalScore = float32(finalScore)
		err := s.resultRepo.Update(ctx, &result)
		if err != nil {
			return entity.Result{}, err
		}
	}
	return result, err

}
func (s *examService) AttemptExamEvent(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (entity.ExamEventAttempt, error) {
	eventStatus, err := s.eventService.GetStatus(ctx, eventSlug, accountId)

	if err != nil {
		return entity.ExamEventAttempt{}, err
	}

	if eventStatus.IsHasNotStarted {
		return entity.ExamEventAttempt{}, http_error.EVENT_NOT_STARTED
	}

	ev, exam, err := s.GetEventExamExisting(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return entity.ExamEventAttempt{}, err
	}
	attemptStatus, examEventAttempt, err := s.GetExamEventAttempt(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return entity.ExamEventAttempt{}, err
	}

	questions, err := s.SetupQuestions(ctx, eventSlug, exam.Id, accountId)
	examEventAttempt.Questions = questions

	if err != nil {
		return entity.ExamEventAttempt{}, err
	}
	if attemptStatus.IsNotAttempt {

		if eventStatus.IsFinished {
			return entity.ExamEventAttempt{}, err.EVENT_FINISHED
		}

		startTime, dueTime := s.SetupExamTimer(ctx, exam)
		remTime := utils.CalculateRemainingTime(startTime, dueTime)

		fmt.Println("Rem Time = ", remTime)
		examEventAttempt = entity.ExamEventAttempt{
			AccountId: accountId,
			EventId:   ev.Data.Id,
			ExamId:    exam.Id,
			CreatedAt: startTime,
			DueAt:     dueTime,
			Submitted: false,
			RemTime:   remTime,
			Questions: questions,
		}

		if err := s.examEventAttemptRepo.Create(ctx, &examEventAttempt); err != nil {
			return entity.ExamEventAttempt{}, err
		}

		answers, err := s.SetupAnswer(ctx, questions, examEventAttempt.Id)
		fmt.Println("Answer = ", answers)
		if err != nil {
			return entity.ExamEventAttempt{}, err
		}

		examEventAttempt.Answers = answers
		return ProtectExamEventAttempt(examEventAttempt), err

	} else if attemptStatus.IsOnAttempt {

		if eventStatus.IsFinished {
			s.SubmitExamEvent(ctx, examEventAttempt.Id)
			examEventAttempt.Submitted = true
			if err := s.examEventAttemptRepo.Update(ctx, &examEventAttempt); err != nil {
				return entity.ExamEventAttempt{}, err
			}
			return examEventAttempt, err
		}
		remTime := utils.CalculateRemainingTime(examEventAttempt.CreatedAt, examEventAttempt.DueAt)
		examEventAttempt.RemTime = remTime

		if err := s.examEventAttemptRepo.Update(ctx, &examEventAttempt); err != nil {
			return entity.ExamEventAttempt{}, err
		}

		examEventAttempt.Questions = questions

		return ProtectExamEventAttempt(examEventAttempt), err

	} else if attemptStatus.IsTimeOut {
		if examEventAttempt.RemTime != 0 {
			remTime := 0
			examEventAttempt.RemTime = remTime
		}

		s.SubmitExamEvent(ctx, examEventAttempt.Id)
		examEventAttempt.Submitted = true
		if err := s.examEventAttemptRepo.Update(ctx, &examEventAttempt); err != nil {
			return entity.ExamEventAttempt{}, err
		}
		return examEventAttempt, err

	} else if attemptStatus.IsSubmitted {
		return examEventAttempt, nil
	}
	return entity.ExamEventAttempt{}, http_error.INTERNAL_SERVER_ERROR
}

func (s *examService) EvaluateAnswer(ctx context.Context, question entity.Questions) evaluator {

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
func (s *examService) AnswerExamEvent(ctx context.Context, eventSlug string, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error) {

	attempt, err := s.examEventAttemptRepo.GetById(ctx, attemptId)

	if err != nil {
		return entity.CPQuestionVerdict{}, err
	}

	eventStatus, err := s.eventService.GetStatus(ctx, eventSlug, attempt.AccountId)

	if err != nil {
		return entity.CPQuestionVerdict{}, err
	}

	if eventStatus.IsFinished {
		return entity.CPQuestionVerdict{}, http_error.EVENT_FINISHED
	}

	if attempt.Submitted {
		return entity.CPQuestionVerdict{}, http_error.EXAMS_SUBMITTED
	}

	question, err := s.problemSetService.GetQuestionById(ctx, questionId)
	if err != nil {
		return entity.CPQuestionVerdict{}, err
	}

	score, CPQuestionVerdict := s.EvaluateAnswer(ctx, question)(answer)

	err = s.examEventAnswerRepo.Update(ctx, &entity.ExamEventAnswer{
		AttemptId:  attemptId,
		QuestionId: questionId,
		Answers:    answer,
		Score:      score,
	})

	return CPQuestionVerdict, err
}

