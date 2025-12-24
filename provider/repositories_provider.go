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
	ProvideExamAcademyAnswerRepository() repositories.ExamAcademyAnswerRepository
	ProvideExamAcademyAssignRepository() repositories.ExamAcademyAssignRepository
	ProvideExamAcademyAttemptRepository() repositories.ExamAcademyAttemptRepository
	ProvideExamEventAnswerRepository() repositories.ExamEventAnswerRepository
	ProvideExamEventAssignRepository() repositories.ExamEventAssignRepository
	ProvideExamEventAttemptRepository() repositories.ExamEventAttemptRepository
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
	academyPaymentRepository repositories.AcademyPaymentRepository
	academyRepository repositories.AcademyRepository
	academyResultRepository repositories.AcademyResultRepository
	accountDetailRepository repositories.AccountDetailRepository
	accountRepository repositories.AccountRepository
	emailVerificationRepository repositories.EmailVerificationRepository
	eventAssignRepository repositories.EventAssignRepository
	eventPaymentRepository repositories.EventPaymentRepository
	eventsRepository repositories.EventsRepository
	examAcademyAnswerRepository repositories.ExamAcademyAnswerRepository
	examAcademyAssignRepository repositories.ExamAcademyAssignRepository
	examAcademyAttemptRepository repositories.ExamAcademyAttemptRepository
	examEventAnswerRepository repositories.ExamEventAnswerRepository
	examEventAssignRepository repositories.ExamEventAssignRepository
	examEventAttemptRepository repositories.ExamEventAttemptRepository
	examRepository repositories.ExamRepository
	externalAuthRepository repositories.ExternalAuthRepository
	fCMRepository repositories.FCMRepository
	fileRepository repositories.FileRepository
	forgotPasswordRepository repositories.ForgotPasswordRepository
	optionRepository repositories.OptionRepository
	problemSetExamAssignRepository repositories.ProblemSetExamAssignRepository
	problemSetRepository repositories.ProblemSetRepository
	questionsRepository repositories.QuestionsRepository
	regionRepository repositories.RegionRepository
	resultRepository repositories.ResultRepository
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
	examAcademyAnswerRepository := repositories.NewExamAcademyAnswerRepository(db)
	examAcademyAssignRepository := repositories.NewExamAcademyAssignRepository(db)
	examAcademyAttemptRepository := repositories.NewExamAcademyAttemptRepository(db)
	examEventAnswerRepository := repositories.NewExamEventAnswerRepository(db)
	examEventAssignRepository := repositories.NewExamEventAssignRepository(db)
	examEventAttemptRepository := repositories.NewExamEventAttemptRepository(db)
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
		academyPaymentRepository: academyPaymentRepository,
		academyRepository: academyRepository,
		academyResultRepository: academyResultRepository,
		accountDetailRepository: accountDetailRepository,
		accountRepository: accountRepository,
		emailVerificationRepository: emailVerificationRepository,
		eventAssignRepository: eventAssignRepository,
		eventPaymentRepository: eventPaymentRepository,
		eventsRepository: eventsRepository,
		examAcademyAnswerRepository: examAcademyAnswerRepository,
		examAcademyAssignRepository: examAcademyAssignRepository,
		examAcademyAttemptRepository: examAcademyAttemptRepository,
		examEventAnswerRepository: examEventAnswerRepository,
		examEventAssignRepository: examEventAssignRepository,
		examEventAttemptRepository: examEventAttemptRepository,
		examRepository: examRepository,
		externalAuthRepository: externalAuthRepository,
		fCMRepository: fCMRepository,
		fileRepository: fileRepository,
		forgotPasswordRepository: forgotPasswordRepository,
		optionRepository: optionRepository,
		problemSetExamAssignRepository: problemSetExamAssignRepository,
		problemSetRepository: problemSetRepository,
		questionsRepository: questionsRepository,
		regionRepository: regionRepository,
		resultRepository: resultRepository,
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

func (r *repositoriesProvider) ProvideExamAcademyAnswerRepository() repositories.ExamAcademyAnswerRepository {
	return r.examAcademyAnswerRepository
}

func (r *repositoriesProvider) ProvideExamAcademyAssignRepository() repositories.ExamAcademyAssignRepository {
	return r.examAcademyAssignRepository
}

func (r *repositoriesProvider) ProvideExamAcademyAttemptRepository() repositories.ExamAcademyAttemptRepository {
	return r.examAcademyAttemptRepository
}

func (r *repositoriesProvider) ProvideExamEventAnswerRepository() repositories.ExamEventAnswerRepository {
	return r.examEventAnswerRepository
}

func (r *repositoriesProvider) ProvideExamEventAssignRepository() repositories.ExamEventAssignRepository {
	return r.examEventAssignRepository
}

func (r *repositoriesProvider) ProvideExamEventAttemptRepository() repositories.ExamEventAttemptRepository {
	return r.examEventAttemptRepository
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

