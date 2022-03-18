package dto

type ExperianDto struct {
	PhoneNumber string `json:"phone_number" validate:"required,number"`
}
