package score_generator_meta_data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ScoreGeneratorsMetadata struct {
	ID                  string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	ParamRequestPayload string    `gorm:"type:varchar(100);column:param_request_payload" json:"param_request_payload"`
	FieldPickle         string    `gorm:"type:varchar(255);column:field_pickle" json:"field_pickle"`
	DataType            string    `gorm:"type:varchar(255);column:data_type" json:"data_type"`
	Length              string    `gorm:"type:varchar(255);column:length" json:"length"`
	IsRequired          string    `gorm:"type:varchar(10);column:is_required" json:"is_required"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
	ScoreGenerators     string    `gorm:"type:varchar(50);column:score_generators" json:"score_generators"`
}

func (c ScoreGeneratorsMetadata) TableName() string {
	return "score_generators_metadata"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ScoreGeneratorsMetadata) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
