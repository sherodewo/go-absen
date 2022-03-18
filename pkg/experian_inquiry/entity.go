package experian_inquiry

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Experian struct {
	ID             string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	ProspectID     *string   `gorm:"type:varchar(20);column:ProspectID"`
	Result         *string   `gorm:"type:varchar(10);column:result"`
	Status         *string   `gorm:"type:varchar(20);column:status"`
	PhoneNumber    *string   `gorm:"type:varchar(20);column:phone_number"`
	Score          *float64  `gorm:"column:score"`
	Experian       *string   `gorm:"type:varchar;column:experian"`
	ScorePro       *string   `gorm:"type:varchar;column:score_pro"`
	ExperianScore  *float64  `gorm:"column:experian_score"`
	ExperianResult *string   `gorm:"type:varchar(20);column:experian_result"`
	InternalScore  *float64  `gorm:"column:internal_score"`
	InternalResult *string   `gorm:"type:varchar(20);column:internal_result"`
	TransID        *string   `gorm:"type:varchar(34);column:trans_id"`
	Type           *string   `gorm:"type:varchar(5);column:type"`
	DataInquiry    string    `gorm:"type:varchar(50);data_inquiry" json:"data_inquiry"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (c Experian) TableName() string {
	return "scorepro_experian_inquiry"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Experian) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}
