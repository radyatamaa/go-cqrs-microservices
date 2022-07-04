package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/helper"
	"gorm.io/gorm"

	validatorGo "github.com/go-playground/validator/v10"
)

func registerCustomValidation(db *gorm.DB, v *validatorGo.Validate) {
	if err := v.RegisterValidation("date_only", ValidateDateOnly); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("date_range", ValidateDateRange); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation(`rfe`, ValidateRequireIfAnotherField); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("enum", ValidateEnum); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("no_space", ValidateNoSpace); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("check_fk", func(fl validatorGo.FieldLevel) bool {
		param := strings.Split(fl.Param(), `:`)
		paramFieldValue := param[0]
		paramTable := param[1]
		paramField := param[2]

		if paramField == `` {
			return true
		}

		// param field reflect.Value.
		var paramReflectValue reflect.Value

		if fl.Parent().Kind() == reflect.Ptr {
			paramReflectValue = fl.Parent().Elem().FieldByName(paramFieldValue)
		} else {
			paramReflectValue = fl.Parent().FieldByName(paramFieldValue)
		}

		// todo : check fl.Field() and paramReflectValue type data before execute to db
		var count int64
		if err := db.Table(paramTable).Where(paramField+"=?", paramReflectValue.Int()).Count(&count).Error; err != nil {
			panic(err)
		}
		if count <= 0 {
			return false
		}
		return true
	}); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("unique_store", func(fl validatorGo.FieldLevel) bool {
		param := strings.Split(fl.Param(), `:`)
		paramField := param[0]
		paramTable := param[1]

		if paramField == `` {
			return true
		}

		// todo : check fl.Field() type data before execute to db
		var count int64
		if err := db.Table(paramTable).Where(paramField+"=?", fl.Field().Interface()).Count(&count).Error; err != nil {
			panic(err)
		}
		if count > 0 {
			return false
		}
		return true
	}); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("unique_update", func(fl validatorGo.FieldLevel) bool {
		if fl.Field().String() != "" {
			param := strings.Split(fl.Param(), `:`)
			paramFieldValue := param[0]
			paramTable := param[1]
			paramField := param[2]
			paramFieldCond := param[3]

			if paramFieldValue == `` {
				return true
			}

			// param field reflect.Value.
			var paramReflectValue reflect.Value

			if fl.Parent().Kind() == reflect.Ptr {
				paramReflectValue = fl.Parent().Elem().FieldByName(paramFieldValue)
			} else {
				paramReflectValue = fl.Parent().FieldByName(paramFieldValue)
			}

			// todo : check fl.Field() and paramReflectValue type data before execute to db
			count := int64(0)
			if err := db.Table(paramTable).Where(paramField+" =?", fl.Field().Interface()).Where(paramFieldCond+" <> ?", paramReflectValue.Int()).
				Count(&count).Error; err != nil {
				panic(err)
			}
			if count > 0 {
				return false
			}
		}
		return true
	}); err != nil {
		panic(err)
	}
}

func ValidateDateOnly(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateDateRange(fl validatorGo.FieldLevel) bool {
	parse := helper.StringToDate(fl.Field().String())
	if parse.IsZero() {
		return false
	}
	return true
}

func ValidateEnum(field validatorGo.FieldLevel) bool {

	if field.Param() == `` {
		return true
	}

	// first, clean/remove the comma
	cleaned := strings.Replace(field.Param(), "-", " ", -1)

	// convert 'cleaned' comma separated string to slice
	strSlice := strings.Fields(cleaned)

	if !helper.ItemExists(strSlice, field.Field().String()) {
		return false
	}

	return true
}

func ValidateRequireIfAnotherField(fl validatorGo.FieldLevel) bool {
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

	if isEq(paramFieldValue, paramValue) == false {
		return true
	}
	return hasValue(fl)
}

func hasValue(fl validatorGo.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func ValidateNoSpace(field validatorGo.FieldLevel) bool {
	value := field.Field().String()

	res := strings.TrimSpace(value)
	if len(res) == 0 {
		return false
	}

	return true
}

func requireCheckFieldKind(fl validatorGo.FieldLevel, param string) bool {
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

	case reflect.Bool:
		p := asBoolean(value)

		return field.Bool() == p
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

func asBoolean(param string) bool {

	i, err := strconv.ParseBool(param)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
