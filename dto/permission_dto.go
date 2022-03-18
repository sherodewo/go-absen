package dto

type PermissionDto struct {
	Feature     string `json:"feature" form:"feature" validate:"required"`
	Url         string `json:"url" form:"url" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

type PermissionUpdateDto struct {
	Feature     string `json:"feature" form:"feature" validate:"required"`
	Url         string `json:"url" form:"url" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}
