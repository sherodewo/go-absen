package models

import (
	"time"
)

type MenuParent struct {
	ID          string    `gorm:"column:id;" json:"id"`
	Name        string    `gorm:"type:varchar(60);column:name"`
	Description string    `gorm:"type:varchar(50);column:parent_id"`
	IsActive    string    `gorm:"column:is_active"`
	UserRoleID  string    `gorm:"user_role_id" json:"user_role_id"`
	IconParent  string    `gorm:"column:icon_parent"`
	Order       int       `gorm:"order" json:"order"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (c *MenuParent) TableName() string {
	return "web_menu_parent"
}
