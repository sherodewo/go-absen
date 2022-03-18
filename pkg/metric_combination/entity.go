package metric_combination

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MetricCombination struct {
	ID               string    `gorm:"type:varchar(50);column:id;primary_key:true" json:"id"`
	IsIndosat        string    `gorm:"type:varchar(3);column:is_indosat;" json:"is_indosat"`
	Experian         string    `gorm:"type:varchar(6);column:experian;" json:"experian"`
	Internal         string    `gorm:"type:varchar(10);column:internal;" json:"internal"`
	ScoreLos         string    `gorm:"type:varchar(20);column:score_los;" json:"score_los"`
	FinalScore       string    `gorm:"type:varchar(20);column:final_score;" json:"final_score"`
	Notes            string    `gorm:"type:varchar(255);column:notes;" json:"notes"`
	IsActive         int       `gorm:"column:is_active;" json:"is_active"`
	CreatedAt        time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;" json:"updated_at"`
	Lob              string    `gorm:"type:varchar(10);column:lob;" json:"lob"`
	ExpectedScoreLos string    `gorm:"type:varchar(15);column:expected_score_los;" json:"expected_score_los"`
}

func (c MetricCombination) TableName() string {
	return "metric_combination"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *MetricCombination) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
