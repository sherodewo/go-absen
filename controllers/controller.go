package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/kreditplus/scorepro/utils"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type Controller struct {
}

func NewCommonResponse() *Controller {
	return &Controller{}
}

type ErrorValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ApiResponse struct {
	Message    string      `json:"messages"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors"`
	ServerTime string      `json:"server_time"`
}

func (c *Controller) Ok(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, ApiResponse{
		Message:    http.StatusText(http.StatusOK),
		Errors:     nil,
		Data:       data,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) BadGateWay(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusBadGateway, ApiResponse{
		Message:    http.StatusText(http.StatusBadGateway),
		Errors:     nil,
		Data:       data,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) OkNotEncrypted(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, ApiResponse{
		Message:    http.StatusText(http.StatusOK),
		Data:       data,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) AccessForbidden(ctx echo.Context, url string) error {
	return ctx.JSON(http.StatusForbidden, ApiResponse{
		Message:    url,
		Errors:     "access_forbidden",
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) InternalServerError(ctx echo.Context, err error) error {
	if listErr := handleUnmarshalError(err); len(listErr) != 0 {
		return ctx.JSON(http.StatusInternalServerError, ApiResponse{
			Message:    http.StatusText(http.StatusInternalServerError),
			Errors:     listErr[0],
			ServerTime: utils.GenerateTimeNow(),
		})
	}
	return ctx.JSON(http.StatusInternalServerError, ApiResponse{
		Message:    http.StatusText(http.StatusInternalServerError),
		Errors:     err.Error(),
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) NotFound(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusNotFound, ApiResponse{
		Message:    http.StatusText(http.StatusNotFound),
		Errors:     err.Error(),
		ServerTime: utils.GenerateTimeNow(),
	})
}
func (c *Controller) NewNotFound(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusNotFound, ApiResponse{
		Message:    fmt.Sprintf("%s", message),
		Errors:     nil,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) Unauthorized(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, ApiResponse{
		Message:    http.StatusText(http.StatusUnauthorized),
		Errors:     err.Error(),
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) ServiceUnavailable(ctx echo.Context, url string) error {
	return ctx.JSON(http.StatusServiceUnavailable, ApiResponse{
		Message:    url,
		Errors:     "service_unavailable",
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) NewUnauthorized(ctx echo.Context, msg ...string) error {
	mapStr := map[string]string{}
	if len(msg) == 2 {
		mapStr["url"] = msg[0]
		mapStr["msg"] = msg[1]
	}
	return ctx.JSON(http.StatusUnauthorized, ApiResponse{
		Message:    mapStr["url"],
		Errors:     mapStr["msg"],
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) BadRequest(ctx echo.Context, err error) error {
	if listErr := handleUnmarshalError(err); len(listErr) != 0 {
		return ctx.JSON(http.StatusUnprocessableEntity, ApiResponse{
			Message:    "ERROR",
			Errors:     listErr,
			ServerTime: utils.GenerateTimeNow(),
		})
	}

	return ctx.JSON(http.StatusBadRequest, ApiResponse{
		Message:    http.StatusText(http.StatusBadRequest),
		Errors:     err.Error(),
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) BadRequestErrorValidation(ctx echo.Context, err error) error {
	var errors = make([]ErrorValidation, len(err.(validator.ValidationErrors)))
	for k, v := range err.(validator.ValidationErrors) {
		errors[k] = ErrorValidation{
			Field:   v.Field(),
			Message: v.Tag(),
		}
	}
	return ctx.JSON(http.StatusUnprocessableEntity, ApiResponse{
		Message:    "ERROR",
		Errors:     errors,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func handleUnmarshalError(err error) []ErrorValidation {
	var apiErrors []ErrorValidation
	if he, ok := err.(*echo.HTTPError); ok {
		if ute, ok := he.Internal.(*json.UnmarshalTypeError); ok {
			valError := ErrorValidation{
				Field:   ute.Field,
				Message: ute.Error(),
			}
			apiErrors = append(apiErrors, valError)
		}
		if se, ok := he.Internal.(*json.SyntaxError); ok {
			valError := ErrorValidation{
				Field:   "Syntax Error",
				Message: se.Error(),
			}
			apiErrors = append(apiErrors, valError)
		}
		if iue, ok := he.Internal.(*json.InvalidUnmarshalError); ok {
			valError := ErrorValidation{
				Field:   iue.Type.String(),
				Message: iue.Error(),
			}
			apiErrors = append(apiErrors, valError)
		}
	}
	return apiErrors
}

func (c *Controller) BadRequestWithSpecificFields(ctx echo.Context, msg string, fields ...[]string) error {
	errors := make([]ErrorValidation, len(fields))

	for key, val := range fields {
		errors[key] = formatMessage2(val[0], val[1])
	}

	return ctx.JSON(http.StatusBadRequest, ApiResponse{
		Message:    msg,
		Errors:     errors,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *Controller) BadRequestWithSpecificFieldResponses(ctx echo.Context, msg string, fields ...[]string) error {
	errors := make([]ErrorValidation, len(fields))

	for key, val := range fields {
		errors[key] = ErrorValidation{
			Field:   strcase.ToSnake(val[0]),
			Message: "accepted:" + strcase.ToSnake(val[1]),
		}
	}
	return ctx.JSON(http.StatusBadRequest, ApiResponse{
		Message:    msg,
		Errors:     errors,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func formatMessage2(field string, tag string) ErrorValidation {
	errors := ErrorValidation{
		Field:   strcase.ToSnake(field),
		Message: "required",
	}

	switch tag {
	case "required":
		errors.Message = "required"
	case "mobilephone":
		errors.Message = "accepted:start=08XXXXXXXXXX"
	case "minphone":
		errors.Message = "accepted:min=10"
	case "maxphone":
		errors.Message = "accepted:max=13"
	case "number":
		errors.Message = "accepted:number"
	}
	return errors
}
