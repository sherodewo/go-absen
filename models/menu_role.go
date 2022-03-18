package models

import (
	"time"
)

type MenuRole struct {
	ID         string    `gorm:"column:id;" json:"id"`
	MenuID     string    `gorm:"type:varchar(60);column:menu_id"`
	UserRoleID string    `gorm:"type:varchar(50);column:user_role_id"`
	Create     int       `gorm:"column:create"`
	Read       int       `gorm:"column:read"`
	Update     int       `gorm:"column:update"`
	Delete     int       `gorm:"column:delete"`
	Upload     int       `gorm:"column:upload"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (c *MenuRole) TableName() string {
	return "menu_role"
}
