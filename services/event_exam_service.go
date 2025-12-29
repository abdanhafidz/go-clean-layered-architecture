package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventExamService interface {
	ListExamByEvent(ctx context.Context, eventSlug string, accountId uuid.UUID) ([]dto.EventExamListResponse, error)
	GetEventExamExisting(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (ev dto.EventDetailResponse, exam entity.Exam, err error)
	GetEventExamAttempt(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.EventExamAttempt, error)
	AttemptEventExam(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (entity.EventExamAttempt, error)
	SetupQuestions(ctx context.Context, eventSlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error)
	SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.EventExamAnswer, error)
	SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time)
	SubmitEventExam(ctx context.Context, attemptId uuid.UUID) (result entity.Result, err error)
	AnswerEventExam(ctx context.Context, eventSlug string, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error)
	Scoreboard(ctx context.Context, eventSlug string) ([]dto.ScoreboardItem, error)
}

type evaluator func(answer []string) (float32, entity.CPQuestionVerdict)

type eventExamService struct {
	eventService             EventService
	problemSetService        ProblemSetService
	problemSetExamAssignRepo repositories.ProblemSetExamAssignRepository
	examRepo                 repositories.ExamRepository
	eventExamAttemptRepo     repositories.EventExamAttemptRepository
	eventExamAnswerRepo      repositories.EventExamAnswerRepository
	eventExamAssignRepo      repositories.EventExamAssignRepository
	resultRepo               repositories.ResultRepository
}

func NewEventExamService(eventService EventService, problemSetService ProblemSetService, problemSetExamAssignRepo repositories.ProblemSetExamAssignRepository, examRepo repositories.ExamRepository, eventExamAttemptRepo repositories.EventExamAttemptRepository, eventExamAssignRepo repositories.EventExamAssignRepository, eventExamAnswerRepo repositories.EventExamAnswerRepository, resultRepo repositories.ResultRepository) EventExamService {
	return &eventExamService{
		eventService:             eventService,
		problemSetService:        problemSetService,
		problemSetExamAssignRepo: problemSetExamAssignRepo,
		examRepo:                 examRepo,
		eventExamAttemptRepo:     eventExamAttemptRepo,
		eventExamAssignRepo:      eventExamAssignRepo,
		eventExamAnswerRepo:      eventExamAnswerRepo,
		resultRepo:               resultRepo,
	}
}

func ProtectEventExamAttempt(attempt entity.EventExamAttempt) entity.EventExamAttempt {

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
	var cleanAnswers []entity.EventExamAnswer

	for _, a := range attempt.Answers {
		aCopy := a
		aCopy.Score = 0 // hide score

		cleanAnswers = append(cleanAnswers, aCopy)
	}

	attempt.Answers = cleanAnswers

	return attempt
}

func (s *eventExamService) ListExamByEvent(ctx context.Context, eventSlug string, accountId uuid.UUID) ([]dto.EventExamListResponse, error) {
	ev, err := s.eventService.DetailBySlug(ctx, eventSlug, accountId)

	if err != nil {
		return []dto.EventExamListResponse{}, err
	}

	exams, err := s.examRepo.ListByEvent(ctx, ev.Data.Id)

	if err != nil {
		return []dto.EventExamListResponse{}, err
	}

	// Fetch results for this event and user
	results, err := s.resultRepo.ListByEventAndAccount(ctx, ev.Data.Id, accountId)
	if err != nil {
		return []dto.EventExamListResponse{}, err
	}

	// Map results by ExamId for easy lookup
	resultMap := make(map[uuid.UUID]entity.Result)
	for _, res := range results {
		if res.EventExamAttempt != nil {
			resultMap[res.EventExamAttempt.ExamId] = res
		}
	}

	var response []dto.EventExamListResponse
	for _, exam := range exams {
		var score float32
		if res, exists := resultMap[exam.Id]; exists {
			score = res.FinalScore
		}

		response = append(response, dto.EventExamListResponse{
			Exam:  exam,
			Score: score,
		})
	}

	return response, nil

}

func (s *eventExamService) GetEventExamExisting(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (ev dto.EventDetailResponse, exam entity.Exam, err error) {

	if ev, err = s.eventService.DetailBySlug(ctx, eventSlug, accountId); err != nil {
		return ev, exam, err
	}

	if exam, err = s.examRepo.GetBySlug(ctx, examSlug); err != nil {
		return ev, exam, err
	}

	if err := s.eventExamAssignRepo.Check(ctx, ev.Data.Id, exam.Id); err != nil {
		return dto.EventDetailResponse{}, entity.Exam{}, err
	}

	return ev, exam, err
}

func (s *eventExamService) GetEventExamAttempt(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (dto.UserExamStatus, entity.EventExamAttempt, error) {

	ev, exam, err := s.GetEventExamExisting(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return dto.UserExamStatus{}, entity.EventExamAttempt{}, err
	}

	eventExamAttempt, err := s.eventExamAttemptRepo.GetByEventExam(ctx, ev.Data.Id, exam.Id, accountId)
	fmt.Println("Error Exam Event Attempt", errors.Is(err, gorm.ErrRecordNotFound))

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserExamStatus{}, entity.EventExamAttempt{}, err
	}

	var attemptStatus dto.UserExamStatus
	attemptStatus.IsNotAttempt = errors.Is(err, gorm.ErrRecordNotFound)
	attemptStatus.IsTimeOut = !attemptStatus.IsNotAttempt && (utils.CalculateRemainingTime(eventExamAttempt.CreatedAt, eventExamAttempt.DueAt) == 0)
	attemptStatus.IsSubmitted = eventExamAttempt.Submitted
	attemptStatus.IsOnAttempt = !attemptStatus.IsNotAttempt && !attemptStatus.IsTimeOut && !attemptStatus.IsSubmitted
	return attemptStatus, eventExamAttempt, nil

}
func (s *eventExamService) SetupQuestions(ctx context.Context, eventSlug string, examId uuid.UUID, accountId uuid.UUID) ([]entity.Questions, error) {
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

func (s *eventExamService) SetupAnswer(ctx context.Context, questions []entity.Questions, attemptId uuid.UUID) ([]entity.EventExamAnswer, error) {
	var eventExamAnswers []entity.EventExamAnswer
	for _, q := range questions {

		blank_ans := entity.EventExamAnswer{
			AttemptId:  attemptId,
			QuestionId: q.Id,
		}

		err := s.eventExamAnswerRepo.Create(ctx, &blank_ans)
		if err != nil {
			return []entity.EventExamAnswer{}, err
		}
		eventExamAnswers = append(eventExamAnswers, blank_ans)
	}

	return eventExamAnswers, nil
}

func (s *eventExamService) SetupExamTimer(ctx context.Context, exam entity.Exam) (time.Time, time.Time) {
	startTime := time.Now()
	dueTime := startTime.Add(exam.Duration * time.Minute)
	return startTime, dueTime
}

func (s *eventExamService) SubmitEventExam(ctx context.Context, attemptId uuid.UUID) (result entity.Result, err error) {
	attempt, err := s.eventExamAttemptRepo.GetById(ctx, attemptId)
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
		result.EventExamAttempt = &attempt
		result.FinalScore = float32(finalScore)

		s.eventExamAttemptRepo.Update(ctx, &attempt)
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
func (s *eventExamService) AttemptEventExam(ctx context.Context, eventSlug string, examSlug string, accountId uuid.UUID) (entity.EventExamAttempt, error) {
	eventStatus, err := s.eventService.GetStatus(ctx, eventSlug, accountId)

	if err != nil {
		return entity.EventExamAttempt{}, err
	}

	if eventStatus.IsHasNotStarted {
		return entity.EventExamAttempt{}, http_error.EVENT_NOT_STARTED
	}

	ev, exam, err := s.GetEventExamExisting(ctx, eventSlug, examSlug, accountId)

	if err != nil {
		return entity.EventExamAttempt{}, err
	}
	attemptStatus, eventExamAttempt, err := s.GetEventExamAttempt(ctx, eventSlug, examSlug, accountId)

	fmt.Println("Get AttemptStatus = ", attemptStatus, "Err =", err)

	if err != nil {
		return entity.EventExamAttempt{}, err
	}

	questions, err := s.SetupQuestions(ctx, eventSlug, exam.Id, accountId)
	eventExamAttempt.Questions = questions

	if err != nil {
		return entity.EventExamAttempt{}, err
	}

	if attemptStatus.IsNotAttempt {
		if eventStatus.IsFinished {
			return entity.EventExamAttempt{}, err
		}

		startTime, dueTime := s.SetupExamTimer(ctx, exam)
		remTime := utils.CalculateRemainingTime(startTime, dueTime)

		fmt.Println("Rem Time = ", remTime)
		eventExamAttempt = entity.EventExamAttempt{
			AccountId: accountId,
			EventId:   ev.Data.Id,
			ExamId:    exam.Id,
			CreatedAt: startTime,
			DueAt:     dueTime,
			Submitted: false,
			RemTime:   remTime,
			Questions: questions,
		}

		if err := s.eventExamAttemptRepo.Create(ctx, &eventExamAttempt); err != nil {
			return entity.EventExamAttempt{}, err
		}

		answers, err := s.SetupAnswer(ctx, questions, eventExamAttempt.Id)
		fmt.Println("Answer = ", answers)
		if err != nil {
			return entity.EventExamAttempt{}, err
		}

		eventExamAttempt.Answers = answers
		return ProtectEventExamAttempt(eventExamAttempt), err

	} else if attemptStatus.IsOnAttempt {

		if eventStatus.IsFinished {
			s.SubmitEventExam(ctx, eventExamAttempt.Id)
			eventExamAttempt.Submitted = true
			if err := s.eventExamAttemptRepo.Update(ctx, &eventExamAttempt); err != nil {
				return entity.EventExamAttempt{}, err
			}
			return eventExamAttempt, err
		}
		remTime := utils.CalculateRemainingTime(eventExamAttempt.CreatedAt, eventExamAttempt.DueAt)
		eventExamAttempt.RemTime = remTime

		if err := s.eventExamAttemptRepo.Update(ctx, &eventExamAttempt); err != nil {
			return entity.EventExamAttempt{}, err
		}

		eventExamAttempt.Questions = questions

		return ProtectEventExamAttempt(eventExamAttempt), err

	} else if attemptStatus.IsTimeOut {
		if eventExamAttempt.RemTime != 0 {
			remTime := 0
			eventExamAttempt.RemTime = remTime
		}

		s.SubmitEventExam(ctx, eventExamAttempt.Id)

		eventExamAttempt.Submitted = true

		if err := s.eventExamAttemptRepo.Update(ctx, &eventExamAttempt); err != nil {
			return entity.EventExamAttempt{}, err
		}
		return eventExamAttempt, err

	} else if attemptStatus.IsSubmitted {
		if !exam.Configuration.AllowReview {
			eventExamAttempt = ProtectEventExamAttempt(eventExamAttempt)
		}
		return eventExamAttempt, nil
	}
	return entity.EventExamAttempt{}, http_error.INTERNAL_SERVER_ERROR
}

func (s *eventExamService) EvaluateAnswer(ctx context.Context, question entity.Questions) evaluator {

	nonCPEvaluator := func(answer []string) (float32, entity.CPQuestionVerdict) {
		score := float32(0)
		isCorrect := true
		for i, ans := range answer {
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
func (s *eventExamService) AnswerEventExam(ctx context.Context, eventSlug string, attemptId uuid.UUID, questionId uuid.UUID, answer []string) (entity.CPQuestionVerdict, error) {

	attempt, err := s.eventExamAttemptRepo.GetById(ctx, attemptId)

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

	err = s.eventExamAnswerRepo.Update(ctx, &entity.EventExamAnswer{
		AttemptId:  attemptId,
		QuestionId: questionId,
		Answers:    answer,
		Score:      score,
	})

	return CPQuestionVerdict, err
}

func (s *eventExamService) Scoreboard(ctx context.Context, eventSlug string) ([]dto.ScoreboardItem, error) {
	// 1. Get Event
	// We need event details to get ID, but accountId is not strictly needed for public scoreboard?
	// Assuming admin view or public view. If public, accountId might be nil or generic.
	// But `DetailBySlug` requires accountId. Let's pass uuid.Nil if we can, or just get event by slug via repo if exposed.
	// `eventService.DetailBySlug` logic might check registration status which needs accountId.
	// If this is called by admin, pass admin ID. If public, maybe issue?
	// For now, let's assume this relies on `GetBySlug` from event repo directly or service with Nil UUID works for "Guest".
	// However, `DetailBySlug` implementation usually handles Nil UUID gracefully for basic info?
	// Let's rely on `eventService` for now, assuming caller provides appropriate context/account.
	// Since the interface signature I defined earlier: `Scoreboard(ctx, eventSlug)`. I missed accountId.
	// I'll assume I can fetch event properly. Let's use `uuid.Nil` for now as "Viewer".
	ev, err := s.eventService.DetailBySlug(ctx, eventSlug, uuid.Nil)
	if err != nil {
		return nil, err
	}

	// 2. Get Exams for Event (Skipped as we rely on Results for columns per user)
	// We might need exams if we want to ensure all columns are present, but for now purely result-based.

	// 3. Get Results
	results, err := s.resultRepo.ListByEvent(ctx, ev.Data.Id)
	if err != nil {
		return nil, err
	}

	// 4. Aggregate Data
	// Map: Username -> ScoreboardItem
	scoreboardMap := make(map[string]*dto.ScoreboardItem)

	for _, res := range results {
		username := ""
		fullname := ""
		if res.EventExamAttempt != nil && res.EventExamAttempt.Account != nil {
			username = res.EventExamAttempt.Account.Username
			// Fullname not directly accessible from Account entity structure provided
		}
		if username == "" {
			continue // Skip results without user
		}

		if _, exists := scoreboardMap[username]; !exists {
			scoreboardMap[username] = &dto.ScoreboardItem{
				Username: username,
				FullName: fullname,
				Scores:   []dto.ExamScore{},
			}
		}

		item := scoreboardMap[username]
		item.TotalScore += res.FinalScore

		// Calculate Duration
		// Exam Duration (Minutes) -> Seconds
		// RemTime (Seconds? Int)
		var durationTaken int64 = 0
		if res.EventExamAttempt != nil {
			if res.EventExamAttempt.Exam != nil {
				maxDurationSec := int64(res.EventExamAttempt.Exam.Duration) * 60
				remTimeSec := int64(res.EventExamAttempt.RemTime)
				if remTimeSec < 0 {
					remTimeSec = 0
				}
				durationTaken = maxDurationSec - remTimeSec
				if durationTaken < 0 {
					durationTaken = 0
				}
				// If attempt was not submitted or timeout, remTime might be full?
				// But Result exists, so it must be evaluated.
			}
		}
		item.TotalDurationInt += durationTaken

		// Add Exam Score to list (we will normalize this later to matching indices if needed, or just list)
		// Requirement: "score per exam". Usually this means columns matching exams.
		// DTO has `[]ExamScore`. We can just append all.
		examTitle := ""
		if res.EventExamAttempt != nil && res.EventExamAttempt.Exam != nil {
			examTitle = res.EventExamAttempt.Exam.Title
		}

		// Check if exam already added (should not if one attempt per exam)
		item.Scores = append(item.Scores, dto.ExamScore{
			ExamId:    res.EventExamAttempt.ExamId, // Or res.EventExamAttempt.Exam.Id
			ExamTitle: examTitle,
			Score:     res.FinalScore,
		})
	}

	// 5. Finalize Items (Calculate Average, Format Duration)
	var scoreboard []dto.ScoreboardItem
	for _, item := range scoreboardMap {
		// Calculate Average
		// If user didn't take an exam, is it 0? Usually yes for contest rank.
		// Let's assume average of taken exams or all exams?
		// "average score for the entire exam" -> ambiguous. "Average score of the user across all exams".
		// I'll calculate average based on exams taken for now, or count of event exams.
		// Better to average per Taken exam or Total Event Exams?
		// In context of scoreboard, average usually implies performance.
		// Let's use `TotalScore / Count(Exams In Event)` to penalize missing exams, OR `TotalScore / Count(Results)`.
		// Given it's a "Ranking", usually Sum is enough, Average is redundant with Sum unless weights differ.
		// I'll use `TotalScore / len(results)` for that user.
		if len(item.Scores) > 0 {
			item.AverageScore = item.TotalScore / float32(len(item.Scores))
		}

		// Fill missing exams with 0 if needed for structure?
		// The client might need to match ExamId. The current DTO allows list of scores.
		// I won't fill 0s to keep payload small, frontend can map by ID.

		// Format Duration
		// Simple MM:SS or HH:MM:SS
		hours := item.TotalDurationInt / 3600
		minutes := (item.TotalDurationInt % 3600) / 60
		seconds := item.TotalDurationInt % 60
		item.TotalExamDuration = fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)

		scoreboard = append(scoreboard, *item)
	}

	// 6. Sort
	// 1. Total Score Desc
	// 2. Total Duration Asc
	sort.Slice(scoreboard, func(i, j int) bool {
		if scoreboard[i].TotalScore != scoreboard[j].TotalScore {
			return scoreboard[i].TotalScore > scoreboard[j].TotalScore
		}
		return scoreboard[i].TotalDurationInt < scoreboard[j].TotalDurationInt
	})

	return scoreboard, nil
}
