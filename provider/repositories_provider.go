package provider

import "abdanhafidz.com/go-boilerplate/repositories"

type RepositoriesProvider interface {
	ProvideAcademyPaymentRepository() repositories.AcademyPaymentRepository
	ProvideAcademyRepository() repositories.AcademyRepository
	ProvideAcademyResultRepository() repositories.AcademyResultRepository
	ProvideAccountDetailRepository() repositories.AccountDetailRepository
	ProvideAccountRepository() repositories.AccountRepository
	ProvideEmailVerificationRepository() repositories.EmailVerificationRepository
	ProvideEventAssignRepository() repositories.EventAssignRepository
	ProvideEventPaymentRepository() repositories.EventPaymentRepository
	ProvideEventsRepository() repositories.EventsRepository
	ProvideAcademyExamAnswerRepository() repositories.AcademyExamAnswerRepository
	ProvideAcademyExamAssignRepository() repositories.AcademyExamAssignRepository
	ProvideAcademyExamAttemptRepository() repositories.AcademyExamAttemptRepository
	ProvideEventExamAnswerRepository() repositories.EventExamAnswerRepository
	ProvideEventExamAssignRepository() repositories.EventExamAssignRepository
	ProvideEventExamAttemptRepository() repositories.EventExamAttemptRepository
	ProvideEventExamProctoringRepository() repositories.EventExamProctoringRepository
	ProvideExamRepository() repositories.ExamRepository
	ProvideExternalAuthRepository() repositories.ExternalAuthRepository
	ProvideFCMRepository() repositories.FCMRepository
	ProvideFileRepository() repositories.FileRepository
	ProvideForgotPasswordRepository() repositories.ForgotPasswordRepository
	ProvideOptionRepository() repositories.OptionRepository
	ProvideProblemSetExamAssignRepository() repositories.ProblemSetExamAssignRepository
	ProvideProblemSetRepository() repositories.ProblemSetRepository
	ProvideQuestionsRepository() repositories.QuestionsRepository
	ProvideRegionRepository() repositories.RegionRepository
	ProvideResultRepository() repositories.ResultRepository
}

type repositoriesProvider struct {
	academyPaymentRepository       repositories.AcademyPaymentRepository
	academyRepository              repositories.AcademyRepository
	academyResultRepository        repositories.AcademyResultRepository
	accountDetailRepository        repositories.AccountDetailRepository
	accountRepository              repositories.AccountRepository
	emailVerificationRepository    repositories.EmailVerificationRepository
	eventAssignRepository          repositories.EventAssignRepository
	eventPaymentRepository         repositories.EventPaymentRepository
	eventsRepository               repositories.EventsRepository
	academyExamAnswerRepository    repositories.AcademyExamAnswerRepository
	academyExamAssignRepository    repositories.AcademyExamAssignRepository
	academyExamAttemptRepository   repositories.AcademyExamAttemptRepository
	eventExamAnswerRepository      repositories.EventExamAnswerRepository
	eventExamAssignRepository      repositories.EventExamAssignRepository
	eventExamAttemptRepository     repositories.EventExamAttemptRepository
	eventExamProctoringRepository  repositories.EventExamProctoringRepository
	examRepository                 repositories.ExamRepository
	externalAuthRepository         repositories.ExternalAuthRepository
	fCMRepository                  repositories.FCMRepository
	fileRepository                 repositories.FileRepository
	forgotPasswordRepository       repositories.ForgotPasswordRepository
	optionRepository               repositories.OptionRepository
	problemSetExamAssignRepository repositories.ProblemSetExamAssignRepository
	problemSetRepository           repositories.ProblemSetRepository
	questionsRepository            repositories.QuestionsRepository
	regionRepository               repositories.RegionRepository
	resultRepository               repositories.ResultRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	dbConfig := cfg.ProvideDatabaseConfig()
	db := dbConfig.GetInstance()

	academyPaymentRepository := repositories.NewAcaddemyPaymentRepository(db)
	academyRepository := repositories.NewAcademyRepository(db)
	academyResultRepository := repositories.NewAcademyResultRepository(db)
	accountDetailRepository := repositories.NewAccountDetailRepository(db)
	accountRepository := repositories.NewAccountRepository(db)
	emailVerificationRepository := repositories.NewEmailVerificationRepository(db)
	eventAssignRepository := repositories.NewEventAssignRepository(db)
	eventPaymentRepository := repositories.NewEventPaymentRepository(db)
	eventsRepository := repositories.NewEventsRepository(db)
	academyExamAnswerRepository := repositories.NewAcademyExamAnswerRepository(db)
	academyExamAssignRepository := repositories.NewAcademyExamAssignRepository(db)
	academyExamAttemptRepository := repositories.NewAcademyExamAttemptRepository(db)
	eventExamAnswerRepository := repositories.NewEventExamAnswerRepository(db)
	eventExamAssignRepository := repositories.NewEventExamAssignRepository(db)
	eventExamAttemptRepository := repositories.NewEventExamAttemptRepository(db)
	eventExamProctoringRepository := repositories.NewEventExamProctoringRepository(db)
	examRepository := repositories.NewExamRepository(db)
	externalAuthRepository := repositories.NewExternalAuthRepository(db)
	fCMRepository := repositories.NewFCMRepository(db)
	fileRepository := repositories.NewFileRepository(db)
	forgotPasswordRepository := repositories.NewForgotPasswordRepository(db)
	optionRepository := repositories.NewOptionRepository(db)
	problemSetExamAssignRepository := repositories.NewProblemSetExamAssignRepository(db)
	problemSetRepository := repositories.NewProblemSetRepository(db)
	questionsRepository := repositories.NewQuestionsRepository(db)
	regionRepository := repositories.NewRegionRepository(db)
	resultRepository := repositories.NewResultRepository(db)

	return &repositoriesProvider{
		academyPaymentRepository:       academyPaymentRepository,
		academyRepository:              academyRepository,
		academyResultRepository:        academyResultRepository,
		accountDetailRepository:        accountDetailRepository,
		accountRepository:              accountRepository,
		emailVerificationRepository:    emailVerificationRepository,
		eventAssignRepository:          eventAssignRepository,
		eventPaymentRepository:         eventPaymentRepository,
		eventsRepository:               eventsRepository,
		academyExamAnswerRepository:    academyExamAnswerRepository,
		academyExamAssignRepository:    academyExamAssignRepository,
		academyExamAttemptRepository:   academyExamAttemptRepository,
		eventExamAnswerRepository:      eventExamAnswerRepository,
		eventExamAssignRepository:      eventExamAssignRepository,
		eventExamAttemptRepository:     eventExamAttemptRepository,
		eventExamProctoringRepository:  eventExamProctoringRepository,
		examRepository:                 examRepository,
		externalAuthRepository:         externalAuthRepository,
		fCMRepository:                  fCMRepository,
		fileRepository:                 fileRepository,
		forgotPasswordRepository:       forgotPasswordRepository,
		optionRepository:               optionRepository,
		problemSetExamAssignRepository: problemSetExamAssignRepository,
		problemSetRepository:           problemSetRepository,
		questionsRepository:            questionsRepository,
		regionRepository:               regionRepository,
		resultRepository:               resultRepository,
	}
}

func (r *repositoriesProvider) ProvideAcademyPaymentRepository() repositories.AcademyPaymentRepository {
	return r.academyPaymentRepository
}

func (r *repositoriesProvider) ProvideAcademyRepository() repositories.AcademyRepository {
	return r.academyRepository
}

func (r *repositoriesProvider) ProvideAcademyResultRepository() repositories.AcademyResultRepository {
	return r.academyResultRepository
}

func (r *repositoriesProvider) ProvideAccountDetailRepository() repositories.AccountDetailRepository {
	return r.accountDetailRepository
}

func (r *repositoriesProvider) ProvideAccountRepository() repositories.AccountRepository {
	return r.accountRepository
}

func (r *repositoriesProvider) ProvideEmailVerificationRepository() repositories.EmailVerificationRepository {
	return r.emailVerificationRepository
}

func (r *repositoriesProvider) ProvideEventAssignRepository() repositories.EventAssignRepository {
	return r.eventAssignRepository
}

func (r *repositoriesProvider) ProvideEventPaymentRepository() repositories.EventPaymentRepository {
	return r.eventPaymentRepository
}

func (r *repositoriesProvider) ProvideEventsRepository() repositories.EventsRepository {
	return r.eventsRepository
}

func (r *repositoriesProvider) ProvideAcademyExamAnswerRepository() repositories.AcademyExamAnswerRepository {
	return r.academyExamAnswerRepository
}

func (r *repositoriesProvider) ProvideAcademyExamAssignRepository() repositories.AcademyExamAssignRepository {
	return r.academyExamAssignRepository
}

func (r *repositoriesProvider) ProvideAcademyExamAttemptRepository() repositories.AcademyExamAttemptRepository {
	return r.academyExamAttemptRepository
}

func (r *repositoriesProvider) ProvideEventExamAnswerRepository() repositories.EventExamAnswerRepository {
	return r.eventExamAnswerRepository
}

func (r *repositoriesProvider) ProvideEventExamAssignRepository() repositories.EventExamAssignRepository {
	return r.eventExamAssignRepository
}

func (r *repositoriesProvider) ProvideEventExamAttemptRepository() repositories.EventExamAttemptRepository {
	return r.eventExamAttemptRepository
}

func (r *repositoriesProvider) ProvideEventExamProctoringRepository() repositories.EventExamProctoringRepository {
	return r.eventExamProctoringRepository
}

func (r *repositoriesProvider) ProvideExamRepository() repositories.ExamRepository {
	return r.examRepository
}

func (r *repositoriesProvider) ProvideExternalAuthRepository() repositories.ExternalAuthRepository {
	return r.externalAuthRepository
}

func (r *repositoriesProvider) ProvideFCMRepository() repositories.FCMRepository {
	return r.fCMRepository
}

func (r *repositoriesProvider) ProvideFileRepository() repositories.FileRepository {
	return r.fileRepository
}

func (r *repositoriesProvider) ProvideForgotPasswordRepository() repositories.ForgotPasswordRepository {
	return r.forgotPasswordRepository
}

func (r *repositoriesProvider) ProvideOptionRepository() repositories.OptionRepository {
	return r.optionRepository
}

func (r *repositoriesProvider) ProvideProblemSetExamAssignRepository() repositories.ProblemSetExamAssignRepository {
	return r.problemSetExamAssignRepository
}

func (r *repositoriesProvider) ProvideProblemSetRepository() repositories.ProblemSetRepository {
	return r.problemSetRepository
}

func (r *repositoriesProvider) ProvideQuestionsRepository() repositories.QuestionsRepository {
	return r.questionsRepository
}

func (r *repositoriesProvider) ProvideRegionRepository() repositories.RegionRepository {
	return r.regionRepository
}

func (r *repositoriesProvider) ProvideResultRepository() repositories.ResultRepository {
	return r.resultRepository
}
