package data_inquiry

import (
	"time"
)

type DataInquiry struct {
	ID                string    `gorm:"type:varchar(50);column:id;primary_key:true"`
	ProspectID        string    `gorm:"type:varchar(50);column:prospect_id"`
	StatusKonsumen    string    `gorm:"type:varchar(5);column:status_konsumen"`
	MobilePhone       string    `gorm:"type:varchar(15);column:mobile_phone"`
	SupplierID        *string   `gorm:"type:varchar(100);column:supplier_id"`
	MetricCombination *string   `gorm:"type:varchar(50);column:metric_combination" json:"metric_combination"`
	Journey           string    `gorm:"type:varchar(5);column:journey"`
	TypeModel         string    `gorm:"type:varchar(50);column:type_model"`
	ScoreproScore     string    `gorm:"type:varchar(50);column:scorepro_score"`
	CreatedAt         time.Time `gorm:"column:created_at;"`
	ScoreproRequestor string    `gorm:"type:varchar(50);column:scorepro_requestor"`
}

func (c DataInquiry) TableName() string {
	return "data_inquiry"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
//func (c *DataInquiry) BeforeCreate(tx *gorm.DB) (err error) {
//	c.ID = uuid.New().String()
//
//	return
//}
