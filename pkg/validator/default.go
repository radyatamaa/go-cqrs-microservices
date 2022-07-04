package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	validatorGo "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslator "github.com/go-playground/validator/v10/translations/id"
	"gorm.io/gorm"
)

type defaultValidator struct {
	once       sync.Once
	db         *gorm.DB
	validate   *validatorGo.Validate
	translator *ut.UniversalTranslator
}

type ValidateDynamicStructError struct {
	Field string                 `json:"field"`
	Msg   validatorGo.FieldError `json:"msg"`
}

func (m *ValidateDynamicStructError) Error() string {
	return m.Msg.Error()
}

func (v *defaultValidator) SetDatabaseConnection(db *gorm.DB) {
	v.db = db
}

func (v *defaultValidator) GetTranslator(locale string) (ut.Translator, bool) {
	return v.translator.GetTranslator(locale)
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

var _ structValidator = &defaultValidator{}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.validateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func (v *defaultValidator) ValidateDynamicStruct(dynamicStruct map[string]interface{}, expectedStruct interface{}) error {
	t := reflect.TypeOf(expectedStruct)
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		if _, ok := dynamicStruct[field.Tag.Get("json")]; ok {
			tag := field.Tag.Get("validate")
			err := v.ValidateVar(dynamicStruct[field.Tag.Get("json")], field.Name, tag)
			if err != nil {
				return &ValidateDynamicStructError{
					Field: field.Tag.Get("json"),
					Msg:   err.(validatorGo.ValidationErrors)[0],
				}
			}
		}

	}

	return nil
}

func (v *defaultValidator) ValidateMatchingDynamicStruct(body io.ReadCloser, expectedStruct interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err := dec.Decode(expectedStruct)
	if err != nil {
		return err
	}
	return nil
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj interface{}) error {
	v.lazyInit()
	return v.validate.Struct(obj)
}

func (v *defaultValidator) ValidateVar(val interface{}, field interface{}, tag string) error {
	v.lazyInit()
	return v.validate.VarWithValue(val, field, tag)
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {

		v.validate = validatorGo.New()
		v.validate.SetTagName("validate")

		// add any custom validations etc. here
		registerCustomValidation(v.db, v.validate)

		englishTranslate := en.New()
		v.translator = ut.New(englishTranslate, englishTranslate, id.New())
		if indonesiaTranslator, found := v.translator.GetTranslator("id"); found {
			if err := idTranslator.RegisterDefaultTranslations(v.validate, indonesiaTranslator); err != nil {
				panic(err)
			}
			registerCustomIndonesianTranslator(v.validate, indonesiaTranslator)
		}
		if englishTranslator, found := v.translator.GetTranslator("en"); found {
			if err := enTranslations.RegisterDefaultTranslations(v.validate, englishTranslator); err != nil {
				panic(err)
			}
			registerCustomEnglishTranslator(v.validate, englishTranslator)
		}
	})
}
