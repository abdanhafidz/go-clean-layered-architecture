package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Account struct {
	Id                uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username          string     `gorm:"uniqueIndex" json:"username,omitempty"`
	Email             string     `gorm:"uniqueIndex" json:"email,omitempty"`
	Role              string     `json:"role,omitempty"`
	Password          string     `json:"-"`
	IsEmailVerified   bool       `json:"is_email_verified,omitempty"`
	IsDetailCompleted bool       `json:"is_detail_completed,omitempty"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
}

func (Account) TableName() string { return "account" }

type AccountDetail struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountId   uuid.UUID `json:"account_id,omitempty"`
	FullName    *string   `json:"full_name,omitempty"`
	SchoolName  *string   `json:"school_name,omitempty"`
	Province    *string   `json:"province,omitempty"`
	City        *string   `json:"city,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Account     *Account  `gorm:"foreignKey:AccountId" json:"account,omitempty"`
}

func (AccountDetail) TableName() string { return "account_details" }

type EmailVerification struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     uint      `json:"token,omitempty"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	IsExpired bool      `json:"is_expired,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	Account   *Account  `gorm:"foreignKey:AccountId" json:"account,omitempty"`
}

func (EmailVerification) TableName() string { return "email_verification" }

type ExternalAuth struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OauthID       string    `json:"oauth_id,omitempty"`
	AccountId     uuid.UUID `json:"account_id,omitempty"`
	OauthProvider string    `json:"oauth_provider,omitempty"`
}

func (ExternalAuth) TableName() string { return "external_auth" }

type FCM struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	FCMToken  string    `json:"fcm_token,omitempty"`
}

func (FCM) TableName() string { return "fcm" }

type ForgotPassword struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     uint      `json:"token,omitempty"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	IsExpired bool      `json:"is_expired,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

func (ForgotPassword) TableName() string { return "forgot_password" }

type Events struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_event"`
	Title      string    `json:"title,omitempty"`
	Slug       string    `json:"slug,omitempty"`
	StartEvent time.Time `json:"start_event,omitempty"`
	EndEvent   time.Time `json:"end_event,omitempty"`
	Overview   string    `json:"overview,omitempty"`
	EventCode  string    `json:"event_code,omitempty"`
	IsPublic   bool      `json:"is_public,omitempty"`
}

func (Events) TableName() string { return "events" }

type Announcement struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_announcement"`
	Title     string    `json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Message   string    `json:"message,omitempty"`
	Publisher string    `json:"publisher,omitempty"`
	EventId   uuid.UUID `json:"id_event,omitempty"`
}

func (Announcement) TableName() string { return "announcement" }

type ProblemSet struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_problem_set"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (ProblemSet) TableName() string { return "problem_set" }

type Exam struct {
	Id          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_exam"`
	Slug        string        `json:"slug,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
	Randomize   uint          `json:"randomize,omitempty"`
}

func (Exam) TableName() string { return "exam" }

type OptionCategory struct {
	Id         uint   `gorm:"primaryKey" json:"id"`
	OptionName string `json:"option_name,omitempty"`
	OptionSlug string `json:"option_slug,omitempty"`
}

func (OptionCategory) TableName() string { return "option_category" }

type OptionValues struct {
	Id               uint   `gorm:"primaryKey" json:"id"`
	OptionCategoryId uint   `json:"option_category_id,omitempty"`
	OptionValue      string `json:"option_value,omitempty"`
}

func (OptionValues) TableName() string { return "option_values" }

type RegionProvince struct {
	Id   uint   `json:"id"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

func (RegionProvince) TableName() string { return "region_provinces" }

type RegionCity struct {
	Id         uint   `json:"id"`
	Type       string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
	Code       string `json:"code,omitempty"`
	FullCode   string `json:"full_code,omitempty"`
	ProvinceId uint   `json:"province_id,omitempty"`
}

func (RegionCity) TableName() string { return "region_cities" }

type Options struct {
	OptionCategory OptionCategory `json:"option_category,omitempty"`
	OptionValues   []OptionValues `json:"option_values,omitempty"`
}

func (Options) TableName() string { return "options" }

type EventAssign struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_assign"`
	AccountId  uuid.UUID `json:"id_account,omitempty"`
	EventId    uuid.UUID `json:"id_event,omitempty"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`

	Account *Account `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Event   *Events  `gorm:"foreignKey:EventId" json:"event,omitempty"`
}

func (EventAssign) TableName() string { return "event_assign" }

type Questions struct {
	Id           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_question"`
	Type         string         `json:"type,omitempty"`
	Question     string         `json:"question,omitempty"`
	Options      pq.StringArray `gorm:"type:text[]" json:"options,omitempty"`
	AnsKey       pq.StringArray `gorm:"type:text[]" json:"ans_key,omitempty"`
	CorrMark     float64        `json:"corr_mark,omitempty"`
	IncorrMark   float64        `json:"incorr_mark,omitempty"`
	NullMark     float64        `json:"null_mark,omitempty"`
	ProblemSetId uuid.UUID      `json:"id_problem_set,omitempty"`

	ProblemSet *ProblemSet `gorm:"foreignKey:ProblemSetId" json:"problem_set,omitempty"`
}

func (Questions) TableName() string { return "questions" }

type ProblemSetExamAssign struct {
	Id           uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_problem_set_exam_assign"`
	ExamId       uuid.UUID   `json:"id_exam,omitempty"`
	ProblemSetId uuid.UUID   `json:"id_problem_set,omitempty"`
	Exam         *Exam       `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	ProblemSet   *ProblemSet `gorm:"foreignKey:ProblemSetId" json:"problem_set,omitempty"`
}

func (ProblemSetExamAssign) TableName() string { return "problem_set_exam_assign" }

type ExamEventAssign struct {
	Id      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_exam_event_assign"`
	EventId uuid.UUID `json:"id_event,omitempty"`
	ExamId  uuid.UUID `json:"id_exam,omitempty"`
	Exam    *Exam     `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	Event   *Events   `gorm:"foreignKey:EventId" json:"event,omitempty"`
}

func (ExamEventAssign) TableName() string { return "exam_event_assign" }

type CPQuestionVerdict struct {
	TimeExecution float32 `json:"time_exec"`
	MemoryUsage   float32 `json:"memory"`
	Verdict       string  `json:"verdict"`
	Score         float32 `json:"score"`
}

type ExamEventAnswer struct {
	Id               uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AttemptId        uuid.UUID         `json:"id_attempt,omitempty" gorm:"index"`
	QuestionId       uuid.UUID         `json:"id_question,omitempty" gorm:"index"`
	Answers          pq.StringArray    `gorm:"type:text[]" json:"answer,omitempty"`
	Score            float32           `json:"score"`
	ExamEventAttempt *ExamEventAttempt `gorm:"foreignKey:AttemptId" json:"exam_attempt,omitempty"`
	CreatedAt        time.Time         `json:"created_at,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at,omitempty"`
}

func (ExamEventAnswer) TableName() string { return "exam_event_answer" }

type ExamEventAttempt struct {
	Id        uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_attempt"`
	AccountId uuid.UUID         `json:"id_account,omitempty"`
	EventId   uuid.UUID         `json:"id_event,omitempty"`
	ExamId    uuid.UUID         `json:"id_exam,omitempty"`
	Questions []Questions       `gorm:"-" json:"questions,omitempty"`
	Answers   []ExamEventAnswer `gorm:"foreignKey:AttemptId;references:Id" json:"answers,omitempty"`
	Account   *Account          `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Event     *Events           `gorm:"foreignKey:EventId" json:"event,omitempty"`
	Exam      *Exam             `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	RemTime   int               `json:"remaining_time,omitempty"`
	Mark      float32           `json:"mark,omitempty"`
	CreatedAt time.Time         `json:"created_at,omitempty"`
	DueAt     time.Time         `json:"due_at,omitempty"`
	Submitted bool              `json:"submitted,omitempty"`
}

func (ExamEventAttempt) TableName() string { return "exam_event_attempt" }

type Result struct {
	Id               uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_result"`
	AttemptId        uuid.UUID         `json:"id_attempt,omitempty"`
	FinalScore       float32           `json:"final_score"`
	ExamEventAttempt *ExamEventAttempt `gorm:"foreignKey:AttemptId" json:"exam_attempt,omitempty"`
}

func (Result) TableName() string { return "result" }

type Academy struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (Academy) TableName() string { return "academy" }

type AcademyMaterial struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	AcademyId   uuid.UUID `json:"academy_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (AcademyMaterial) TableName() string { return "academy_materials" }

type AcademyContent struct {
	Id                uuid.UUID `gorm:"primaryKey" json:"id"`
	Title             string    `json:"title,omitempty"`
	Order             uint      `json:"order,omitempty"`
	AcademyMaterialId uuid.UUID `json:"academy_material_id,omitempty"`
	Contents          string    `json:"contents,omitempty"`
}

func (AcademyContent) TableName() string { return "academy_contents" }

type AcademyMaterialProgress struct {
	Id                uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountId         uint      `json:"account_id,omitempty"`
	AcademyMaterialId uuid.UUID `json:"academy_material_id,omitempty"`
	Progress          uint      `json:"progress,omitempty"`
}

func (AcademyMaterialProgress) TableName() string { return "academy_materials_progress" }

type AcademyContentProgress struct {
	Id        uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	AcademyId uuid.UUID `json:"academy_id,omitempty"`
}

func (AcademyContentProgress) TableName() string { return "academy_contents_progress" }
