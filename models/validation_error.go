package models

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

type ValidationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
	Message   string `json:"message"`
}

type ErrorValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func WrapValidationErrors(errs validator.ValidationErrors) []ValidationError {
	validationErrors := make([]ValidationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, ValidationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
			Message:   formatMessage(validationErr),
		})
	}

	return validationErrors
}

func WrapValidationApiErrors(errs validator.ValidationErrors) []ErrorValidation {
	validationErrors := make([]ErrorValidation, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, ErrorValidation{
			Field:   ToSnakeCase(validationErr.Field()),
			Message: formatMessage(validationErr),
		})
	}

	return validationErrors
}

func formatMessage(err validator.FieldError) string {
	field := err.Field()
	param := err.Param()

	field = ToSnakeCase(field)

	message := fmt.Sprintf("Field validation for '%s' failed on the '%s'", err.Field(), err.Tag())

	switch err.Tag() {
	case "required":
		message = fmt.Sprintf("The %s field is required.", field)
	case "numeric":
		message = fmt.Sprintf("The %s must be a number.", field)
	case "email":
		message = fmt.Sprintf("The %s must be a valid email address.", field)
	case "gt":
		message = fmt.Sprintf("The %s must be greater than %s.", field, param)
	case "gte":
		message = fmt.Sprintf("The %s must be greater than or equal %s.", field, param)
	case "lt":
		message = fmt.Sprintf("The %s must be less than %s.", field, param)
	case "lte":
		message = fmt.Sprintf("The %s must be less than or equal %s.", field, param)
	case "phone_number":
		message = fmt.Sprintf("The %s must be a valid phone number.", field)
	case "min":
		message = fmt.Sprintf("The %s must be at least %s", field, param)
	case "max":
		message = fmt.Sprintf("The %s may not be greater than %s.", field, param)
	case "len":
		message = fmt.Sprintf("The %s must be a length %s.", field, param)
	case "eq":
		message = fmt.Sprintf("The %s must be a equals %s.", field, param)
	case "mobile_phone":
		message = fmt.Sprintf("The %s must be valid format %s.", field, "08xxxxxxxxxx")
	case "date_only":
		message = fmt.Sprintf("The %s must be valid format %s.", field, "yyyy-mm-dd")
	case "rfe":
		param := strings.Split(err.Param(), `:`)
		paramField := param[0]
		paramValue := param[1]
		message = fmt.Sprintf("The %s is required if %s = %s.", field, paramField, paramValue)
	case "rfg":
		param := strings.Split(err.Param(), `:`)
		paramField := param[0]
		paramValue := param[1]
		message = fmt.Sprintf("The %s is required if %s > %s.", field, paramField, paramValue)
	case "rfne":
		param := strings.Split(err.Param(), `:`)
		paramField := param[0]
		paramValue := param[1]
		message = fmt.Sprintf("The %s is required if %s != %s.", field, paramField, paramValue)
	}

	return message
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake  = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}