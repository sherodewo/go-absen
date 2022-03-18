package score_generator

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ScoreGenerator struct {
	ID                   string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	Name                 string    `gorm:"type:varchar(100);column:name" json:"name"`
	Description          string    `gorm:"type:varchar(255);column:description" json:"description"`
	Endpoint             string    `gorm:"type:varchar(255);column:endpoint" json:"endpoint"`
	FilePickle           string    `gorm:"type:varchar(255);column:file_pickle" json:"file_pickle"`
	SaveResultTo         string    `gorm:"type:varchar(10);column:save_result_to" json:"save_result_to"`
	SaveResultObjectName string    `gorm:"type:varchar(255);column:save_result_object_name" json:"save_result_object_name"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
	ScoreModelsRules     string    `gorm:"type:varchar(50);column:score_models_rules" json:"score_models_rules"`
}

func (c ScoreGenerator) TableName() string {
	return "score_generators"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ScoreGenerator) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}

type WizardDto struct {
	Name                string `json:"name" param:"name" form:"name"`
	Description         string `json:"description" param:"description" form:"description"`
	Endpoint            string `json:"endpoint" param:"endpoint" form:"endpoint"`
	SaveResultTo        string `json:"save_result_to" param:"save_result_to" form:"save_result_to"`
	SaveResultName      string `json:"save_result_object_name" param:"save_result_object_name" form:"save_result_object_name"`
	ScoreModelsRules    string `json:"score_models_rules" param:"score_models_rules" form:"score_models_rules"`
	FilePickle          string `json:"file_pickle" param:"file_pickle" form:"file_pickle"`
	ParamRequestPayload string `json:"param_request_payload" param:"param_request_payload" form:"param_request_payload"`
	FieldPickle         string `json:"field_pickle" param:"field_pickle" form:"field_pickle"`
	DataType            string `json:"data_type" param:"data_type" form:"data_type"`
	Length              string `json:"length" param:"length" form:"length"`
	IsRequired          string `json:"is_required" param:"is_required" form:"is_required"`
}
