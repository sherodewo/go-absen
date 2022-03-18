package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AppConfig struct {
	ID        string    `gorm:"column:id;primary_key:true"`
	GroupName string    `gorm:"column:group_name"`
	Key       string    `gorm:"column:key"`
	Value     string    `gorm:"column:value"`
	IsActive  int       `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (c *AppConfig) TableName() string {
	return "app_config"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *AppConfig) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
