package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
	"time"
)

type Dummy struct {
	ID               string    `gorm:"type:varchar(20);column:id;primary_key:true"`
	PhoneNumber      string    `gorm:"type:varchar(200);column:phone_number"`
	Score            float64   `gorm:"column:score"`
	ExperianResponse string    `gorm:"column:experian_response"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (c Dummy) TableName() string {
	return "experian_dummy"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Dummy) BeforeCreate(tx *gorm.DB) (err error) {
	guid := xid.New()
	c.ID = guid.String()

	return
}
