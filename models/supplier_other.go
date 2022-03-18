package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type SupplierOther struct {
	ID        string `gorm:"type:varchar(20);column:id;primary_key:true" json:"id"`
	ZipCode   string `gorm:"column:zip_code" json:"zip_code"`
	TypeModel string `gorm:"column:type_model" json:"type_model"`
	IsActive  string `gorm:"column:is_active" json:"is_active"`
}

func (c SupplierOther) TableName() string {
	return "supplier_other"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SupplierOther) BeforeCreate(tx *gorm.DB) (err error) {
	guid := xid.New()
	c.ID = guid.String()

	return
}
