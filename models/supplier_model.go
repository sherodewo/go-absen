package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type SupplierModel struct {
	ID string `gorm:"type:varchar(20);column:id;primary_key:true" json:"id"`
	SupplierID string `gorm:"column:supplier_id" json:"supplier_id"`
	TypeModel string `gorm:"column:type_model" json:"type_model"`
	IsActive string `gorm:"column:is_active" json:"is_active"`
}

func (c SupplierModel) TableName() string {
	return "supplier_model"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SupplierModel) BeforeCreate(tx *gorm.DB) (err error) {
	guid := xid.New()
	c.ID = guid.String()

	return
}
