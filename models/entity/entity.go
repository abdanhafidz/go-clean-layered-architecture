package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Account struct {
	Id                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username          string         `gorm:"uniqueIndex" json:"username,omitempty"`
	Email             string         `gorm:"uniqueIndex" json:"email,omitempty"`
	Role              string         `json:"role,omitempty"`
	Password          string         `json:"-"`
	IsEmailVerified   bool           `json:"is_email_verified,omitempty"`
	IsDetailCompleted bool           `json:"is_detail_completed,omitempty"`
	CreatedAt         time.Time      `json:"created_at,omitempty"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (Account) TableName() string { return "account" }

type AccountDetail struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
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
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Token     uint      `json:"token,omitempty"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	IsExpired bool      `json:"is_expired,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	Account   *Account  `gorm:"foreignKey:AccountId" json:"account,omitempty"`
}

func (EmailVerification) TableName() string { return "email_verification" }

type ExternalAuth struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OauthID       string    `json:"oauth_id,omitempty"`
	AccountId     uuid.UUID `json:"account_id,omitempty"`
	OauthProvider string    `json:"oauth_provider,omitempty"`
}

func (ExternalAuth) TableName() string { return "external_auth" }

type FCM struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	FCMToken  string    `json:"fcm_token,omitempty"`
}

func (FCM) TableName() string { return "fcm" }

type ForgotPassword struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Token     uint      `json:"token,omitempty"`
	AccountId uuid.UUID `json:"account_id,omitempty"`
	IsExpired bool      `json:"is_expired,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

func (ForgotPassword) TableName() string { return "forgot_password" }

type Events struct {
	Id             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_event"`
	Title          string         `json:"title,omitempty"`
	Slug           string         `json:"slug,omitempty"`
	StartEvent     time.Time      `json:"start_event,omitempty"`
	EndEvent       time.Time      `json:"end_event,omitempty"`
	Overview       string         `json:"overview,omitempty"`
	ImgBanner      string         `json:"img_banner,omitempty"`
	EventCode      string         `json:"event_code,omitempty"`
	IsPublic       bool           `json:"is_public,omitempty"`
	CreatedAt      time.Time      `json:"created_at,omitempty"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Price          float64        `json:"price"`
	RegisterStatus int            `gorm:"-" json:"register_status"`
}

func (Events) TableName() string { return "events" }

type Announcement struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_announcement"`
	Title     string    `json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Message   string    `json:"message,omitempty"`
	Publisher string    `json:"publisher,omitempty"`
	EventId   uuid.UUID `json:"id_event,omitempty"`
}

func (Announcement) TableName() string { return "announcement" }

type ProblemSet struct {
	Id          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_problem_set"`
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (ProblemSet) TableName() string { return "problem_set" }

type Exam struct {
	Id            uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_exam"`
	Slug          string            `json:"slug,omitempty"`
	Title         string            `json:"title,omitempty"`
	Description   string            `json:"description,omitempty"`
	Duration      time.Duration     `json:"duration,omitempty"`
	Randomize     uint              `json:"randomize,omitempty"`
	Configuration ExamConfiguration `gorm:"foreignKey:ExamId;references:Id" json:"configuration,omitempty"`
	Proctoring    ExamProctoring    `gorm:"foreignKey:ExamId;references:Id" json:"proctoring,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	DeletedAt     gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"index"`
}

func (Exam) TableName() string { return "exam" }

type ExamConfiguration struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_result"`
	ExamId      uuid.UUID `json:"id_exam,omitempty"`
	AllowRetake bool      `json:"allow_retake,omitempty"`
	AllowReview bool      `json:"allow_review,omitempty"`
	EnableTimer bool      `json:"enable_timer,omitempty"`
}

func (ExamConfiguration) TableName() string { return "exam_configuration" }

type ExamProctoring struct {
	Id                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_result"`
	ExamId             uuid.UUID `json:"id_exam,omitempty"`
	EnableWebCam       bool      `json:"enable_webcam,omitempty"`
	EnableVAD          bool      `json:"enable_vad,omitempty"`
	EnableTabBlock     bool      `json:"enable_tab_block,omitempty"`
	RequiredFullScreen bool      `json:"enable_full_screen,omitempty"`
	EnableEyeTracking  bool      `json:"enable_eye_tracking,omitempty"`
	DisableCopyPaste   bool      `json:"disable_copy_paste,omitempty"`
	EnableExamBrowser  bool      `json:"enable_exam_browser,omitempty"`
}

func (ExamProctoring) TableName() string { return "exam_proctoring" }

type EventExamProctoringLogs struct {
	Id                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_result"`
	EventId           uuid.UUID `json:"id_event,omitempty"`
	ExamId            uuid.UUID `json:"id_exam,omitempty"`
	AccountId         uuid.UUID `json:"id_account,omitempty"`
	ViolationScore    uint      `json:"violation_score,omitempty"`
	ViolationCategory string    `json:"violation_category,omitempty"`
	Attachement       string    `json:"attachement,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}

func (EventExamProctoringLogs) TableName() string { return "exam_event_proctoring_logs" }

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
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_assign"`
	AccountId  uuid.UUID `json:"id_account,omitempty"`
	EventId    uuid.UUID `json:"id_event,omitempty"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`

	Account *Account `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Event   *Events  `gorm:"foreignKey:EventId" json:"event,omitempty"`
}

func (EventAssign) TableName() string { return "event_assign" }

type Questions struct {
	Id           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_question"`
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
	Id           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_problem_set_exam_assign"`
	ExamId       uuid.UUID   `json:"id_exam,omitempty"`
	ProblemSetId uuid.UUID   `json:"id_problem_set,omitempty"`
	Exam         *Exam       `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	ProblemSet   *ProblemSet `gorm:"foreignKey:ProblemSetId" json:"problem_set,omitempty"`
}

func (ProblemSetExamAssign) TableName() string { return "problem_set_exam_assign" }

type EventExamAssign struct {
	Id      uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_exam_event_assign"`
	EventId uuid.UUID `json:"id_event,omitempty"`
	ExamId  uuid.UUID `json:"id_exam,omitempty"`
	Exam    *Exam     `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	Event   *Events   `gorm:"foreignKey:EventId" json:"event,omitempty"`
}

func (EventExamAssign) TableName() string { return "exam_event_assign" }

type CPQuestionVerdict struct {
	TimeExecution float32 `json:"time_exec"`
	MemoryUsage   float32 `json:"memory"`
	Verdict       string  `json:"verdict"`
	Score         float32 `json:"score"`
}

type EventExamAnswer struct {
	Id               uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AttemptId        uuid.UUID         `json:"id_attempt,omitempty" gorm:"index"`
	QuestionId       uuid.UUID         `json:"id_question,omitempty" gorm:"index"`
	Answers          pq.StringArray    `gorm:"type:text[]" json:"answer,omitempty"`
	Score            float32           `json:"score"`
	EventExamAttempt *EventExamAttempt `gorm:"foreignKey:AttemptId" json:"exam_attempt,omitempty"`
	CreatedAt        time.Time         `json:"created_at,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at,omitempty"`
}

func (EventExamAnswer) TableName() string { return "exam_event_answer" }

type EventExamAttempt struct {
	Id        uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_attempt"`
	AccountId uuid.UUID         `json:"id_account,omitempty"`
	EventId   uuid.UUID         `json:"id_event,omitempty"`
	ExamId    uuid.UUID         `json:"id_exam,omitempty"`
	Questions []Questions       `gorm:"-" json:"questions,omitempty"`
	Answers   []EventExamAnswer `gorm:"foreignKey:AttemptId;references:Id" json:"answers,omitempty"`
	Account   *Account          `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Event     *Events           `gorm:"foreignKey:EventId" json:"event,omitempty"`
	Exam      *Exam             `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	RemTime   int               `json:"remaining_time,omitempty"`
	Mark      float32           `json:"mark,omitempty"`
	CreatedAt time.Time         `json:"created_at,omitempty"`
	DueAt     time.Time         `json:"due_at,omitempty"`
	Submitted bool              `json:"submitted,omitempty"`
}

func (EventExamAttempt) TableName() string { return "exam_event_attempt" }

type Result struct {
	Id               uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_result"`
	AttemptId        uuid.UUID         `json:"id_attempt,omitempty"`
	FinalScore       float32           `json:"final_score"`
	EventExamAttempt *EventExamAttempt `gorm:"foreignKey:AttemptId" json:"exam_attempt,omitempty"`
}

func (Result) TableName() string { return "result" }

type Academy struct {
	Id              uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title           string            `json:"title,omitempty"`
	Slug            string            `gorm:"unique" json:"slug,omitempty"`
	Code            string            `gorm:"unique" json:"code,omitempty"`
	IsPublic        bool              `json:"is_public,omitempty"`
	Description     string            `json:"description,omitempty"`
	ImageUrl        string            `json:"image_url,omitempty"`
	MaterialsCount  int64             `json:"materials_count,omitempty"`
	Materials       []AcademyMaterial `gorm:"foreignKey:AcademyId;references:Id" json:"materials,omitempty"`
	AcademyProgress AcademyProgress   `gorm:"foreignKey:AcademyId;references:Id" json:"academy_progress,omitempty"`
	Price           float64           `json:"price,omitempty"`
	RegisterStatus  int               `gorm:"-" json:"register_status"`
	CreatedAt       time.Time         `json:"created_at,omitempty"`
	DeletedAt       gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"index"`
}

func (Academy) TableName() string { return "academy" }

type AcademyMaterial struct {
	Id                      uuid.UUID               `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AcademyId               uuid.UUID               `json:"academy_id,omitempty"`
	Title                   string                  `json:"title,omitempty"`
	Slug                    string                  `gorm:"unique" json:"slug,omitempty"`
	Description             string                  `json:"description,omitempty"`
	Order                   uint                    `json:"order,omitempty"`
	ContentsCount           int64                   `json:"contents_count,omitempty"`
	Contents                []AcademyContent        `gorm:"foreignKey:MaterialId;references:Id" json:"contents,omitempty"`
	AcademyMaterialProgress AcademyMaterialProgress `gorm:"foreignKey:MaterialId;references:Id" json:"academy_material_progress,omitempty"`
	CreatedAt               time.Time               `json:"created_at,omitempty"`
	DeletedAt               gorm.DeletedAt          `json:"deleted_at,omitempty" gorm:"index"`
}

func (AcademyMaterial) TableName() string { return "academy_materials" }

type AcademyContent struct {
	Id                     uuid.UUID              `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MaterialId             uuid.UUID              `json:"material_id,omitempty"`
	Title                  string                 `json:"title,omitempty"`
	Order                  uint                   `json:"order,omitempty"`
	Contents               string                 `json:"contents,omitempty"`
	AcademyContentProgress AcademyContentProgress `gorm:"foreignKey:ContentId;references:Id" json:"academy_content_progress,omitempty"`
	CreatedAt              time.Time              `json:"created_at,omitempty"`
	DeletedAt              gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

func (AcademyContent) TableName() string { return "academy_contents" }

// Progress

type AcademyProgress struct {
	Id                      uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	AccountId               uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_account_academy" json:"account_id,omitempty"`
	AcademyId               uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_account_academy" json:"academy_id,omitempty"`
	Status                  string     `gorm:"type:varchar(50);default:'not attempted'" json:"status,omitempty"`
	Progress                float64    `gorm:"default:0" json:"progress"`
	TotalCompletedMaterials uint       `gorm:"default:0" json:"total_completed_materials"`
	CompletedAt             *time.Time `json:"completed_at"`
}

func (AcademyProgress) TableName() string { return "academy_progress" }

type AcademyMaterialProgress struct {
	Id                     uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	AccountId              uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_account_material" json:"account_id,omitempty"`
	AcademyId              uuid.UUID  `gorm:"type:uuid;index" json:"academy_id,omitempty"`
	MaterialId             uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_account_material" json:"material_id,omitempty"`
	Progress               float64    `gorm:"default:0" json:"progress,omitempty"`
	TotalCompletedContents uint       `gorm:"default:0" json:"total_completed_contents,omitempty"`
	Status                 string     `gorm:"type:varchar(50);default:'not attempted'" json:"status,omitempty"`
	CompletedAt            *time.Time `json:"completed_at"`
}

func (AcademyMaterialProgress) TableName() string { return "academy_material_progress" }

type AcademyContentProgress struct {
	Id          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	AccountId   uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_account_content" json:"account_id,omitempty"`
	AcademyId   uuid.UUID  `gorm:"type:uuid;index" json:"academy_id,omitempty"`
	MaterialId  uuid.UUID  `gorm:"type:uuid;index" json:"material_id,omitempty"`
	ContentId   uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_account_content" json:"content_id,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'not attempted'" json:"status,omitempty"`
	CompletedAt *time.Time `json:"completed_at"`
}

func (AcademyContentProgress) TableName() string { return "academy_content_progress" }

type AcademyAssign struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountId uuid.UUID `gorm:"type:uuid;index" json:"account_id,omitempty"`
	AcademyId uuid.UUID `gorm:"type:uuid;index" json:"academy_id,omitempty"`
	Account   *Account  `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Academy   *Academy  `gorm:"foreignKey:AcademyId" json:"academy,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (AcademyAssign) TableName() string { return "academy_assign" }

type File struct {
	Id           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OriginalName string    `json:"original_name,omitempty"`
	StoredName   string    `json:"stored_name,omitempty"`
	MimeType     string    `json:"mime_type,omitempty"`
	Size         int64     `json:"size,omitempty"`
	Path         string    `json:"path,omitempty"`
	Context      string    `json:"context,omitempty"`
	AccountId    uuid.UUID `json:"account_id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Account      *Account  `gorm:"foreignKey:AccountId" json:"account,omitempty"`
}

func (File) TableName() string { return "files" }

type AcademyExamAssign struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_exam_academy_assign"`
	AcademyId uuid.UUID `json:"id_academy,omitempty"`
	ExamId    uuid.UUID `json:"id_exam,omitempty"`
	Exam      *Exam     `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	Academy   *Academy  `gorm:"foreignKey:AcademyId" json:"academy,omitempty"`
}

func (AcademyExamAssign) TableName() string { return "exam_academy_assign" }

type AcademyExamAnswer struct {
	Id                 uuid.UUID           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AttemptId          uuid.UUID           `json:"id_attempt,omitempty" gorm:"index"`
	QuestionId         uuid.UUID           `json:"id_question,omitempty" gorm:"index"`
	Answers            pq.StringArray      `gorm:"type:text[]" json:"answer,omitempty"`
	Score              float32             `json:"score"`
	AcademyExamAttempt *AcademyExamAttempt `gorm:"foreignKey:AttemptId" json:"exam_attempt,omitempty"`
	CreatedAt          time.Time           `json:"created_at,omitempty"`
	UpdatedAt          time.Time           `json:"updated_at,omitempty"`
}

func (AcademyExamAnswer) TableName() string { return "exam_academy_answer" }

type AcademyExamAttempt struct {
	Id        uuid.UUID           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_attempt"`
	AccountId uuid.UUID           `json:"id_account,omitempty"`
	AcademyId uuid.UUID           `json:"id_academy,omitempty"`
	ExamId    uuid.UUID           `json:"id_exam,omitempty"`
	Questions []Questions         `gorm:"-" json:"questions,omitempty"`
	Answers   []AcademyExamAnswer `gorm:"foreignKey:AttemptId;references:Id" json:"answers,omitempty"`
	Account   *Account            `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Academy   *Academy            `gorm:"foreignKey:AcademyId" json:"academy,omitempty"`
	Exam      *Exam               `gorm:"foreignKey:ExamId" json:"exam,omitempty"`
	RemTime   int                 `json:"remaining_time,omitempty"`
	Mark      float32             `json:"mark,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty"`
	DueAt     time.Time           `json:"due_at,omitempty"`
	Submitted bool                `json:"submitted,omitempty"`
}

func (AcademyExamAttempt) TableName() string { return "exam_academy_attempt" }

type AcademyExamResult struct {
	Id                 uuid.UUID           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id_result"`
	AttemptId          uuid.UUID           `json:"id_attempt,omitempty"`
	FinalScore         float32             `json:"final_score"`
	AcademyExamAttempt *AcademyExamAttempt `gorm:"foreignKey:AttemptId" json:"exam_attempt,omitempty"`
}

func (AcademyExamResult) TableName() string { return "academy_exam_result" }

type AcademyPaymentTransaction struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountId     uuid.UUID `json:"account_id,omitempty"`
	AcademyId     uuid.UUID `json:"academy_id,omitempty"`
	ExternalId    string    `json:"xendit_transaction_id,omitempty"`
	InvoiceId     string    `json:"invoice_id,omitempty"`
	InvoiceUrl    string    `json:"invoice_url,omitempty"`
	Amount        float64   `json:"amount,omitempty"`
	TransactionAt time.Time `json:"transaction_at,omitempty"`
	ExpiredAt     time.Time `json:"expired_at,omitempty"`
	Status        string    `json:"status,omitempty"`
}

func (AcademyPaymentTransaction) TableName() string { return "academy_payment_transaction" }

type EventPaymentTransaction struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountId     uuid.UUID `json:"account_id,omitempty"`
	EventId       uuid.UUID `json:"event_id,omitempty"`
	ExternalId    string    `json:"xendit_transaction_id,omitempty"`
	InvoiceId     string    `json:"invoice_id,omitempty"`
	InvoiceUrl    string    `json:"invoice_url,omitempty"`
	Amount        float64   `json:"amount,omitempty"`
	TransactionAt time.Time `json:"transaction_at,omitempty"`
	ExpiredAt     time.Time `json:"expired_at,omitempty"`
	Status        string    `json:"status,omitempty"`
}

func (EventPaymentTransaction) TableName() string { return "event_payment_transaction" }

type AcademyCoupon struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AcademyId  uuid.UUID `json:"academy_id,omitempty"`
	Code       string    `gorm:"unique" json:"code,omitempty"`
	Discount   float64   `json:"discount,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}

func (AcademyCoupon) TableName() string { return "academy_coupon" }
