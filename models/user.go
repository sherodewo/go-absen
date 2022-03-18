package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// User Struct
type User struct {
	UserID     string    `gorm:"type:varchar(50);column:user_id;primary_key:true"`
	Nik        string    `gorm:"type:varchar(10);column:nik"`
	Name       string    `gorm:"type:varchar(50);column:name"`
	Email      string    `gorm:"type:varchar(100);column:email;unique"`
	Password   string    `gorm:"type:varchar(100);column:password"`
	IsActive   int       `gorm:"column:is_active"`
	TypeUser   int       `gorm:"column:type_user"`
	UserRoleID string    `gorm:"type:varchar(50);column:user_role_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (c User) TableName() string {
	return "user"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *User) BeforeCreate(tx *gorm.DB) (err error) {
	c.UserID = uuid.New().String()

	return
}

type UserViewDetail struct {
	UserID     string `gorm:"type:varchar(20);column:user_id;primary_key:true"`
	Nik        string `gorm:"type:varchar(100);column:nik"`
	Name       string `gorm:"type:varchar(50);column:name"`
	Email      string `gorm:"type:varchar(50);column:email;unique"`
	Password   string `gorm:"type:varchar(200);column:password"`
	IsActive   string `gorm:"column:is_active"`
	TypeUser   int    `gorm:"column:type_user"`
	UserRoleID string `gorm:"column:user_role_id"`
}
