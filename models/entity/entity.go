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

type EmailVerification struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     uint      `json:"token"`
	AccountId uuid.UUID `json:"account_id"`
	IsExpired bool      `json:"is_expired"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Account   *Account  `gorm:"foreignKey:AccountId"`
}

type ExternalAuth struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OauthID       string    `json:"oauth_id"`
	AccountId     uuid.UUID `json:"account_id"`
	OauthProvider string    `json:"oauth_provider"`
}

type FCM struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountId uuid.UUID `json:"account_id"`
	FCMToken  string    `json:"fcm_token"`
}

type ForgotPassword struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     uint      `json:"token"`
	AccountId uuid.UUID `json:"account_id"`
	IsExpired bool      `json:"is_expired"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Events struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_event"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	StartEvent time.Time `json:"start_event"`
	EndEvent   time.Time `json:"end_event"`
	EventCode  string    `json:"event_code"`
	IsPublic   bool      `json:"is_public"`
}

type Announcement struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_announcement"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
	Publisher string    `json:"publisher"`
	EventId   uuid.UUID `json:"id_event"`
}

type ProblemSet struct {
	Id          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_problem_set"`
	Title       string        `json:"title"`
	Duration    time.Duration `json:"duration"`
	Randomize   uint          `json:"randomize"`
	MC_Count    uint          `json:"mc_count"`
	SA_Count    uint          `json:"sa_count"`
	Essay_Count uint          `json:"essay_count"`
}

type Options struct {
	OptionCategory OptionCategory `json:"option_category"`
	OptionValues   []OptionValues `json:"option_values"`
}

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

type EventAssign struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_assign"`
	AccountId  uuid.UUID `json:"id_account"`
	EventId    uuid.UUID `json:"id_event"`
	AssignedAt time.Time `json:"assigned_at"`

	Account *Account `gorm:"foreignKey:AccountId"`
	Event   *Events  `gorm:"foreignKey:EventId"`
}

type ProblemSetAssign struct {
	Id           uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_problem_set_assign"`
	EventId      uuid.UUID   `json:"id_event"`
	ProblemSetId uuid.UUID   `json:"id_problem_set"`
	Event        *Events     `gorm:"foreignKey:EventId"`
	ProblemSet   *ProblemSet `gorm:"foreignKey:ProblemSetId"`
}

type ExamProgress struct {
	Id             uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_progress"`
	AccountId      uuid.UUID   `json:"id_account"`
	EventId        uuid.UUID   `json:"id_event"`
	ProblemSetId   uuid.UUID   `json:"id_problem_set"`
	CreatedAt      time.Time   `json:"created_at"`
	DueAt          time.Time   `json:"due_at"`
	QuestionsOrder []string    `gorm:"type:text[]" json:"questions_order"`
	Answers        any         `gorm:"type:jsonb" json:"answers"`
	Account        *Account    `gorm:"foreignKey:AccountId"`
	Event          *Events     `gorm:"foreignKey:EventId"`
	ProblemSet     *ProblemSet `gorm:"foreignKey:ProblemSetId"`
}

type Result struct {
	Id            uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id_result"`
	AccountId     uuid.UUID     `json:"id_account"`
	EventId       uuid.UUID     `json:"id_event"`
	ProblemSetId  uuid.UUID     `json:"id_problem_set"`
	ProgressId    uuid.UUID     `json:"id_progress"`
	FinishTime    time.Time     `json:"finish_time"`
	Correct       uint          `json:"correct"`
	Incorrect     uint          `json:"incorrect"`
	Empty         uint          `json:"empty"`
	OnCorrection  uint          `json:"on_correction"`
	ManualScoring float64       `json:"manual_scoring"`
	MCScore       float64       `json:"mc_score"`
	ManualScore   float64       `json:"manual_score"`
	FinalScore    float64       `json:"final_score"`
	Account       *Account      `gorm:"foreignKey:AccountId"`
	Event         *Events       `gorm:"foreignKey:EventId"`
	ProblemSet    *ProblemSet   `gorm:"foreignKey:ProblemSetId"`
	ExamProgress  *ExamProgress `gorm:"foreignKey:ProgressId"`
}

type Academy struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

type AcademyMaterial struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	AcademyId   uint      `json:"academy_id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}

type AcademyContent struct {
	Id                uuid.UUID `gorm:"primaryKey" json:"id"`
	Title             string    `json:"title"`
	Order             uint      `json:"order"`
	AcademyMaterialId uint      `json:"academy_material_id"`
	Contents          string    `json:"contents"`
}

type OptionCategory struct {
	Id         uint   `gorm:"primaryKey" json:"id"`
	OptionName string `json:"option_name"`
	OptionSlug string `json:"option_slug"`
}

type OptionValues struct {
	Id               uint   `gorm:"primaryKey" json:"id"`
	OptionCategoryId uint   `json:"option_category_id"`
	OptionValue      string `json:"option_value"`
}
type AcademyMaterialProgress struct {
	Id                uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountId         uint      `json:"account_id"`
	AcademyMaterialId uint      `json:"academy_material_id"`
	Progress          uint      `json:"progress"`
}

type AcademyContentProgress struct {
	Id        uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountId uuid.UUID `json:"account_id"`
	AcademyId uuid.UUID `json:"academy_id"`
}

type RegionProvince struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type RegionCity struct {
	Id         uint   `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	ProvinceId uint   `json:"province_id"`
}

// Gorm table name settings
func (Account) TableName() string                 { return "account" }
func (AccountDetail) TableName() string           { return "account_details" }
func (EmailVerification) TableName() string       { return "email_verification" }
func (ExternalAuth) TableName() string            { return "external_auth" }
func (FCM) TableName() string                     { return "fcm" }
func (ForgotPassword) TableName() string          { return "forgot_password" }
func (Events) TableName() string                  { return "events" }
func (Announcement) TableName() string            { return "announcement" }
func (ProblemSet) TableName() string              { return "problem_sets" }
func (Questions) TableName() string               { return "questions" }
func (EventAssign) TableName() string             { return "event_assign" }
func (ProblemSetAssign) TableName() string        { return "problem_sets_assign" }
func (Result) TableName() string                  { return "result" }
func (ExamProgress) TableName() string            { return "exam_progress" }
func (Academy) TableName() string                 { return "academy" }
func (AcademyMaterial) TableName() string         { return "academy_materials" }
func (AcademyContent) TableName() string          { return "academy_contents" }
func (AcademyMaterialProgress) TableName() string { return "academy_materials_progress" }
func (AcademyContentProgress) TableName() string  { return "academy_contents_progress" }
func (RegionProvince) TableName() string          { return "region_provinces" }
func (RegionCity) TableName() string              { return "region_cities" }
func (Options) TableName() string                 { return "options" }
func (OptionCategory) TableName() string          { return "option_category" }
func (OptionValues) TableName() string            { return "option_values" }
