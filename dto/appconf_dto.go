package dto

type AppConfDto struct {
	GroupName string `json:"group_name" form:"group_name" validate:"required"`
	Key       string `json:"key" form:"key" validate:"required"`
	IsActive  int    `json:"is_active" form:"is_active"`
	Value     string `json:"value" form:"value"`
}
