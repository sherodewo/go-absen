package models

import (
	"time"
)

// TypeUser Struct
type TypeUser struct {
	Type        string    `gorm:"column:type;primary_key:true"`
	Description string    `gorm:"type:nvarchar(100);column:description;unique"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (c TypeUser) TableName() string {
	return "user_type"
}
