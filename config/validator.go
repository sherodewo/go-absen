package config

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

type Validator struct {
	validator *validator.Validate
	sync      sync.Mutex
}

func (v *Validator) Validate(i interface{}) error {
	v.sync.Lock()
	_ = v.validator.RegisterValidation("date_only", validateDateOnly)
	_ = v.validator.RegisterValidation(`rfe`, validateRequireIfAnotherField)
	_ = v.validator.RegisterValidation(`rfne`, validateRequireIfNotEqualAnotherField)
	_ = v.validator.RegisterValidation(`rfg`, validateRequireIfGreaterThanAnotherField)
	v.sync.Unlock()

	return v.validator.Struct(i)
}

func validateDateOnly(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)
	return regex.MatchString(fl.Field().String())
}

func validateRequireIfAnotherField(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}

	if isEq(paramFieldValue, paramValue) == false {
		return true
	}
	return hasValue(fl)
}

func validateRequireIfNotEqualAnotherField(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	// param field reflect.Value.
	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}

	if !isEq(paramFieldValue, paramValue) == false {
		return true
	}
	return hasValue(fl)
}

func validateRequireIfGreaterThanAnotherField(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	// param field reflect.Value.
	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}
	if isGte(paramFieldValue, paramValue) == false {
		return true
	}
	return hasValue(fl)
}

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isEq(field reflect.Value, value string) bool {
	switch field.Kind() {

	case reflect.String:
		return field.String() == value

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func isGte(field reflect.Value, value string) bool {
	timeType := reflect.TypeOf(time.Time{})

	switch field.Kind() {

	case reflect.String:
		p := asInt(value)

		return int64(utf8.RuneCountInString(field.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)

		return int64(field.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)

		return field.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)

		return field.Float() > p

	case reflect.Struct:

		if field.Type() == timeType {

			return field.Interface().(time.Time).After(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
