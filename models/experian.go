package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
	"time"
)

type Experian struct {
	ExperianID     string    `gorm:"type:varchar(20);column:experian_id;primary_key:true"`
	ProspectID     string    `gorm:"type:varchar(100);column:ProspectID"`
	Result         *string   `gorm:"type:varchar(50);column:result"`
	Status         string    `gorm:"type:varchar(200);column:status"`
	PhoneNumber    string    `gorm:"type:varchar(200);column:phone_number"`
	Score          *float64  `gorm:"column:score"`
	Experian       string    `gorm:"column:experian"`
	ScorePro       string    `gorm:"column:score_pro"`
	ExperianScore  *float64  `gorm:"column:experian_score"`
	ExperianResult *string   `gorm:"column:experian_result"`
	InternalScore  *float64  `gorm:"column:internal_score"`
	InternalResult string    `gorm:"column:internal_result"`
	TransID        *string   `gorm:"column:trans_id"`
	Type           string    `gorm:"column:type"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (c Experian) TableName() string {
	return "experian"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Experian) BeforeCreate(tx *gorm.DB) (err error) {
	guid := xid.New()
	c.ExperianID = guid.String()

	return
}
