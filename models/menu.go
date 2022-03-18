package models

import (
	"time"
)

type Menu struct {
	ID        string    `gorm:"type:varchar(50);column:id;" json:"id"`
	Name      string    `gorm:"type:varchar(50);column:name"`
	ParentID  string    `gorm:"type:varchar(50);column:parent_id"`
	Route     string    `gorm:"type:varchar(50);column:route"`
	IconClass string    `gorm:"type:varchar(50);column:icon_class"`
	IsActive  string    `gorm:"type:varchar(2);column:is_active"`
	Role      string    `gorm:"type:varchar(255);role" json:"role"`
	Prefix    string    `gorm:"type:varchar(50);prefix" json:"prefix"`
	MenuType  string    `gorm:"type:varchar(20);menu_type" json:"menu_type"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (c *Menu) TableName() string {
	return "web_menu"
}
