package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

// ValidationError wraps the validators FieldError so we do not
// expose this to our code
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

// Validate the item
// for more detail the returned error can be cast into a
// validator.ValidationErrors collection
//
// if ve, ok := err.(validator.ValidationErrors); ok {
//			fmt.Println(ve.Namespace())
//			fmt.Println(ve.Field())
//			fmt.Println(ve.StructNamespace())
//			fmt.Println(ve.StructField())
//			fmt.Println(ve.Tag())
//			fmt.Println(ve.ActualTag())
//			fmt.Println(ve.Kind())
//			fmt.Println(ve.Type())
//			fmt.Println(ve.Value())
//			fmt.Println(ve.Param())
//			fmt.Println()
//	}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	var returnErrs []ValidationError

	if errs, ok := v.validate.Struct(i).(validator.ValidationErrors); ok {
			for _, err := range errs {
				if fieldErr, ok := err.(validator.FieldError); ok {
					valErr := ValidationError{fieldErr}
					returnErrs = append(returnErrs, valErr)
				}
			}
	}

	return returnErrs
}

// validateSKU
func validateSKU(fl validator.FieldLevel) bool {
	// SKU must be in the format abc-abc-abc
	regex := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	sku := regex.FindAllString(fl.Field().String(), -1)

	return len(sku) == 1

}