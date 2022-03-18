package scorepro_score

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ScoreproScore struct {
	ID              string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	ScoreStart      float64   `gorm:"column:score_start" json:"score_start"`
	ScoreEnd        float64   `gorm:"column:score_end" json:"score_end"`
	ScoreBand		string    `gorm:"type:varchar(100);column:score_band" json:"score_band"`
	Result          string    `gorm:"type:varchar(10);column:result" json:"result"`
	ThresholdOwners string    `gorm:"type:varchar(50);column:threshold_owners" json:"threshold_owners"`
	MaxDSR          int       `gorm:"column:max_dsr" json:"max_dsr"`
	TransactionType string    `gorm:"type:varchar(30);column:transaction_type" json:"transaction_type"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

func (c ScoreproScore) TableName() string {
	return "scorepro_score"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ScoreproScore) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
