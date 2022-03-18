package global

import (
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	DBLOG *gorm.DB
)

func NewGlobalService(db *gorm.DB) {
	DB = db
}

func NewGlobalLogService(db *gorm.DB) {
	DBLOG = db
}
