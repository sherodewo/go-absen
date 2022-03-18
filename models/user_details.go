package models

import (
	"time"
)

// UserDetail Struct
type UserDetail struct {
	UserID      string `gorm:"column:user_id"`
	Nik         string `gorm:"column:nik"`
	Name        string `gorm:"column:name"`
	BranchID    string `gorm:"column:BranchID"`
	BranchName  string `gorm:"column:branch_name"`
	BranchIDOri string `gorm:"column:branch_id_ori"`
	Type        string `gorm:"column:type"`
	Role        string `gorm:"column:role"`
	//UserType	    TypeUser  `gorm:"foreignKey:Type"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (c UserDetail) TableName() string {
	return "user_details"
}
