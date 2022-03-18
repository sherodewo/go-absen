package score_models_rules_data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ScoreModelsRulesData struct {
	ID             string `gorm:"type:varchar(50);column:id;primary_key:true" json:"id"`
	Key            string `gorm:"type:varchar(100);column:key" json:"key"`
	Value          string `gorm:"type:varchar(255);column:value" json:"value" `
	Description    string `gorm:"type:varchar(255);column:description" json:"description" `
	CreatedAt      time.Time `gorm:"created_at:c" json:"created_at" `
	ScoreGenerator string `gorm:"type:varchar(50);column:score_generators" json:"score_generators" `
}

type ScoreModelsRulesDataReq struct {
	ID             string `form:"id" json:"id" param:"id"`
	Key            string `form:"key" json:"key" param:"key"`
	Value          string `form:"value" json:"value" json:"value" `
	Description    string `form:"description" json:"description" json:"description" `
	CreatedAt      string `form:"created_at" json:"created_at" json:"created_at" `
	ScoreGenerator string `form:"score_generators" json:"score_generators" json:"score_generators" `
}

func (c ScoreModelsRulesData) TableName() string {
	return "score_models_rules_data"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ScoreModelsRulesData) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
