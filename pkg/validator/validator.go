package validator

import (
	ut "github.com/go-playground/universal-translator"
	"gorm.io/gorm"
	"io"
)

// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request.
// https://github.com/go-playground/validator/tree/v10.6.1.
type structValidator interface {
	// ValidateStruct can receive any kind of type, and it should never panic, even if the configuration is not right.
	// If the received type is a slice|array, the validation should be performed travel on every element.
	// If the received type is not a struct or slice|array, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise, nil must be returned.
	ValidateStruct(interface{}) error

	ValidateDynamicStruct(dynamicStruct map[string]interface{}, expectedStruct interface{}) error

	ValidateMatchingDynamicStruct(body io.ReadCloser, expectedStruct interface{}) error

	SetDatabaseConnection(*gorm.DB)

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() interface{}

	GetTranslator(locale string) (ut.Translator, bool)
}

// Validate is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v10.6.1
// under the hood.
var Validate structValidator = &defaultValidator{}
