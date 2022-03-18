package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
	"time"
)

// Dashboard Struct
type Dashboard struct {
	UserID        string     `gorm:"type:varchar(20);column:user_id;primary_key:true"`
	Nik           string     `gorm:"type:varchar(100);column:nik"`
	Email         string     `gorm:"type:varchar(50);column:email;unique"`
	Password      string     `gorm:"type:varchar(200);column:password"`
	IsActive      int        `gorm:"column:is_active"`
	IsMultiBranch int        `gorm:"column:is_multi_branch;"`
	TypeUser      string     `gorm:"column:type_user"`
	UserDetail    UserDetail `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (c Dashboard) TableName() string {
	return "user"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Dashboard) BeforeCreate(tx *gorm.DB) (err error) {
	guid := xid.New()
	c.UserID = guid.String()

	return
}
