package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type ExperianScoring struct {
	ID             string    `gorm:"type:varchar(20);column:id;primary_key:true" json:"id"`
	IsIndosat      string    `gorm:"type:varchar(100);column:is_indosat" json:"is_indosat"`
	Experian       string    `gorm:"column:experian" json:"experian"`
	Internal       string    `gorm:"column:internal" json:"internal"`
	ScoreLos       string    `gorm:"column:score_los" json:"score_los"`
	FinalScore     string    `gorm:"column:final_score" json:"final_score"`
	Notes          string    `gorm:"column:notes" json:"notes"`
}

func (c ExperianScoring) TableName() string {
	return "metric_combination"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ExperianScoring) BeforeCreate(tx *gorm.DB) (err error) {
	guid := xid.New()
	c.ID = guid.String()

	return
}
