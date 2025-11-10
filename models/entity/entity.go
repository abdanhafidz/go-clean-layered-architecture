package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id                uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username          string     `gorm:"uniqueIndex" json:"username"`
	Email             string     `gorm:"uniqueIndex" json:"email"`
	Role              string     `json:"role"`
	Password          string     `json:"-"`
	IsEmailVerified   bool       `json:"is_email_verified"`
	IsDetailCompleted bool       `json:"is_detail_completed"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"default:null"`
}

func (Account) TableName() string { return "account" }

type AccountDetail struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountId   uuid.UUID `json:"account_id"`
	FullName    *string   `json:"full_name"`
	SchoolName  *string   `json:"school_name"`
	Province    *string   `json:"province"`
	City        *string   `json:"city"`
	Avatar      *string   `json:"avatar"`
	PhoneNumber *string   `json:"phone_number"`
	Account     *Account  `gorm:"foreignKey:AccountId"`
}

func (AccountDetail) TableName() string { return "account_details" }

type EmailVerification struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     uint      `json:"token"`
	AccountId uuid.UUID `json:"account_id"`
	IsExpired bool      `json:"is_expired"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Account   *Account  `gorm:"foreignKey:AccountId"`
}

func (EmailVerification) TableName() string { return "email_verification" }

type ExternalAuth struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OauthID       string    `json:"oauth_id"`
	AccountId     uuid.UUID `json:"account_id"`
	OauthProvider string    `json:"oauth_provider"`
}

func (ExternalAuth) TableName() string { return "external_auth" }

type FCM struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountId uuid.UUID `json:"account_id"`
	FCMToken  string    `json:"fcm_token"`
}

func (FCM) TableName() string { return "fcm" }

type ForgotPassword struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     uint      `json:"token"`
	AccountId uuid.UUID `json:"account_id"`
	IsExpired bool      `json:"is_expired"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (ForgotPassword) TableName() string { return "forgot_password" }

type Events struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_event"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	StartEvent time.Time `json:"start_event"`
	EndEvent   time.Time `json:"end_event"`
	Overview   string    `json:"overview"`
	EventCode  string    `json:"event_code"`
	IsPublic   bool      `json:"is_public"`
}

func (Events) TableName() string { return "events" }

type Announcement struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_announcement"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
	Publisher string    `json:"publisher"`
	EventId   uuid.UUID `json:"id_event"`
}

func (Announcement) TableName() string { return "announcement" }

type ProblemSet struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_problem_set"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (ProblemSet) TableName() string { return "problem_set" }

type Exam struct {
	Id          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_exam"`
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Randomize   uint          `json:"randomize"`
}

func (Exam) TableName() string { return "exam" }

type OptionCategory struct {
	Id         uint   `gorm:"primaryKey" json:"id"`
	OptionName string `json:"option_name"`
	OptionSlug string `json:"option_slug"`
}

func (OptionCategory) TableName() string { return "option_category" }

type OptionValues struct {
	Id               uint   `gorm:"primaryKey" json:"id"`
	OptionCategoryId uint   `json:"option_category_id"`
	OptionValue      string `json:"option_value"`
}

func (OptionValues) TableName() string { return "option_values" }

type RegionProvince struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (RegionProvince) TableName() string { return "region_provinces" }

type RegionCity struct {
	Id         uint   `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	ProvinceId uint   `json:"province_id"`
}

func (RegionCity) TableName() string { return "region_cities" }

type Options struct {
	OptionCategory OptionCategory `json:"option_category"`
	OptionValues   []OptionValues `json:"option_values"`
}

func (Options) TableName() string { return "options" }

type EventAssign struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_assign"`
	AccountId  uuid.UUID `json:"id_account"`
	EventId    uuid.UUID `json:"id_event"`
	AssignedAt time.Time `json:"assigned_at"`

	Account *Account `gorm:"foreignKey:AccountId"`
	Event   *Events  `gorm:"foreignKey:EventId"`
}

func (EventAssign) TableName() string { return "event_assign" }

type Questions struct {
	Id           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_question"`
	Type         string    `json:"type"` //MultChoices, ShortAns, Essay, IntPuzzle, IntType
	Question     string    `json:"question"`
	Options      []string  `gorm:"type:text[]" json:"options"`
	AnsKey       []string  `gorm:"type:text[]" json:"ans_key"`
	CorrMark     float64   `json:"corr_mark"`
	IncorrMark   float64   `json:"incorr_mark"`
	NullMark     float64   `json:"null_mark"`
	ProblemSetId uuid.UUID `json:"id_problem_set"`

	ProblemSet *ProblemSet `gorm:"foreignKey:ProblemSetId"`
}

func (Questions) TableName() string { return "questions" }

type ProblemSetExamAssign struct {
	Id           uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_problem_set_exam_assign"`
	ExamId       uuid.UUID   `json:"id_exam"`
	ProblemSetId uuid.UUID   `json:"id_problem_set"`
	Exam         *Exam       `gorm:"foreignKey:ExamId"`
	ProblemSet   *ProblemSet `gorm:"foreignKey:ProblemSetId"`
}

func (ProblemSetExamAssign) TableName() string { return "problem_set_exam_assign" }

type ExamEventAssign struct {
	Id      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_exam_event_assign"`
	EventId uuid.UUID `json:"id_event"`
	ExamId  uuid.UUID `json:"id_exam"`
	Exam    *Exam     `gorm:"foreignKey:ExamId"`
	Event   *Events   `gorm:"foreignKey:EventId"`
}

func (ExamEventAssign) TableName() string { return "exam_event_assign" }

type CPQuestionVerdict struct {
	TimeExecution float32 `json:"time_exec"`
	MemoryUsage   float32 `json:"memory"`
	Verdict       string  `json:"verdict"` // AC, WA, PE (pending), QE (queued), TLE, RTE
	Score         float32 `json:"score"`
}
type ExamEventAnswer struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AttemptId  uuid.UUID `json:"id_attempt" gorm:"index"`  // FK ke ExamEventAttempt
	QuestionId uuid.UUID `json:"id_question" gorm:"index"` // FK ke Questions
	Answers    []string  `gorm:"type:text[]" json:"answers"`
	Score      float32   `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ExamEventAnswer) TableName() string { return "exam_event_answer" }

type ExamEventAttempt struct {
	Id        uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_attempt"`
	AccountId uuid.UUID          `json:"id_account"`
	EventId   uuid.UUID          `json:"id_event"`
	ExamId    uuid.UUID          `json:"id_exam"`
	Questions *[]Questions       `json:"questions"`
	Answers   *[]ExamEventAnswer `json:"answers"`
	Account   *Account           `gorm:"foreignKey:AccountId"`
	Event     *Events            `gorm:"foreignKey:EventId"`
	Exam      *Exam              `gorm:"foreignKey:ExamId"`
	RemTime   int                `json:"remaining_time"`
	Mark      float32            `json:"-"`
	CreatedAt time.Time          `json:"created_at"`
	DueAt     time.Time          `json:"due_at"`
	Submitted bool               `json:"submitted"`
}

func (ExamEventAttempt) TableName() string { return "exam_event_attempt" }

type Result struct {
	Id                 uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_result"`
	ExamEventAttemptId uuid.UUID         `json:"id_attempt"`
	FinalScore         float32           `json:"final_score"`
	ExamEventAttempt   *ExamEventAttempt `gorm:"foreignKey:ExamEventAttemptId"`
}

func (Result) TableName() string { return "result" }

type Academy struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

func (Academy) TableName() string { return "academy" }

type AcademyMaterial struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	AcademyId   uuid.UUID `json:"academy_id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

func (AcademyMaterial) TableName() string { return "academy_materials" }

type AcademyContent struct {
	Id                uuid.UUID `gorm:"primaryKey" json:"id"`
	Title             string    `json:"title"`
	Order             uint      `json:"order"`
	AcademyMaterialId uuid.UUID `json:"academy_material_id"`
	Contents          string    `json:"contents"`
}

func (AcademyContent) TableName() string { return "academy_contents" }

type AcademyMaterialProgress struct {
	Id                uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountId         uint      `json:"account_id"`
	AcademyMaterialId uuid.UUID `json:"academy_material_id"`
	Progress          uint      `json:"progress"`
}

func (AcademyMaterialProgress) TableName() string { return "academy_materials_progress" }

type AcademyContentProgress struct {
	Id        uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountId uuid.UUID `json:"account_id"`
	AcademyId uuid.UUID `json:"academy_id"`
}

func (AcademyContentProgress) TableName() string { return "academy_contents_progress" }

// Gorm table name settings
