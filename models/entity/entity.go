package models

import (
	"time"

	"github.com/google/uuid"
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
