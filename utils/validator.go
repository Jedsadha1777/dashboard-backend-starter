package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// InitValidator initializes the validator
func InitValidator() error {
	// Create new validator
	validate = validator.New()

	// Register function to get tag name from json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Create translator
	english := en.New()
	uni := ut.New(english, english)
	var ok bool
	trans, ok = uni.GetTranslator("en")
	if !ok {
		return fmt.Errorf("failed to get english translator")
	}

	// Register default english translations
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return fmt.Errorf("failed to register default translations: %w", err)
	}

	// Register custom translations
	if err := validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("required", fe.Field())
		if err != nil {
			return fe.Field() + " is required" // Fallback
		}
		return t
	}); err != nil {
		return fmt.Errorf("failed to register custom translation: %w", err)
	}

	return nil
}

// ValidateStruct validates a struct and returns error message
func ValidateStruct(s interface{}) error {
	// Initialize validator if not already initialized
	if validate == nil {
		InitValidator()
	}

	// Validate struct
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Translate error messages
	errs := err.(validator.ValidationErrors)
	var errMessages []string

	for _, e := range errs {
		translatedErr := e.Translate(trans)
		errMessages = append(errMessages, translatedErr)
	}

	// Join error messages
	return fmt.Errorf("%s", strings.Join(errMessages, "; "))
}
