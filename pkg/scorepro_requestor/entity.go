package scorepro_requestor

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ScoreproRequestor struct {
	ID          string    `gorm:"type:varchar(50);column:id;primary_key:true" json:"id"`
	Name        string    `gorm:"type:varchar(100);column:name" json:"name"`
	Description float64   `gorm:"type:varchar(255);column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (c ScoreproRequestor) TableName() string {
	return "scorepro_requestors"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ScoreproRequestor) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
