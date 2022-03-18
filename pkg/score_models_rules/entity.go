package score_models_rules

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ScoreModelsRules struct {
	ID               string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	RulesType        string    `gorm:"type:varchar(100);column:rules_type" json:"rules_type"`
	Key              string    `gorm:"type:varchar(255);column:key" json:"key"`
	Value            string    `gorm:"type:varchar(255);column:value" json:"value"`
	Sequence         string    `gorm:"type:varchar(255);column:sequence" json:"sequence"`
	IsActive         string    `gorm:"type:varchar(10);column:is_active" json:"is_active"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	ScoreModelsFlows string    `gorm:"type:varchar(50);column:score_models_flows" json:"score_models_flows"`
	Layer            string    `gorm:"type:varchar(50);column:layer" json:"layer"`
}

func (c ScoreModelsRules) TableName() string {
	return "score_models_rules"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ScoreModelsRules) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
