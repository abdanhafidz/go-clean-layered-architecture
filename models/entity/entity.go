package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Id            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name          string     `gorm:"uniqueIndex" json:"name"`
	DateOfBirth   time.Time  `gorm:"type:date" json:"date_of_birth"`
	AccountNumber uint       `gorm:"uniqueIndex" json:"account_number"`
	Balance       uint       `gorm:"column:balance" json:"balance"`
	PIN           int        `gorm:"column:pin" json:"-"`
	IsActive      bool       `gorm:"type:boolean; column:is_active" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"default:null"`
}

// Gorm table name settings
func (Account) TableName() string { return "account" }
