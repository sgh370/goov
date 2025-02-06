package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"goov/validator/rules"
)

type Validator struct {
	rules map[string][]rules.ValidationRule
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func New() *Validator {
	return &Validator{
		rules: make(map[string][]rules.ValidationRule),
	}
}

func (v *Validator) AddRule(field string, rule rules.ValidationRule) {
	if v.rules[field] == nil {
		v.rules[field] = make([]rules.ValidationRule, 0)
	}
	v.rules[field] = append(v.rules[field], rule)
}

func (v *Validator) validateValue(field string, value reflect.Value, parent interface{}) []ValidationError {
	var errors []ValidationError

	// Get rules for this field
	rules, exists := v.rules[field]
	if !exists {
		return errors
	}

	// Apply all rules
	for _, rule := range rules {
		// For all validations, pass the field value and update parent if needed
		if cf, ok := rule.(interface{ Validate(interface{}) error }); ok {
			if crossField, ok := cf.(interface{ SetParent(interface{}) }); ok {
				crossField.SetParent(parent)
			}
			if err := cf.Validate(value.Interface()); err != nil {
				errors = append(errors, ValidationError{
					Field:   field,
					Message: err.Error(),
				})
			}
			continue
		}

		// For other validations, pass the field value
		if err := rule.Validate(value.Interface()); err != nil {
			errors = append(errors, ValidationError{
				Field:   field,
				Message: err.Error(),
			})
		}
	}

	return errors
}

func (v *Validator) validateStruct(current reflect.Value, parentField string, parent interface{}) []ValidationError {
	var errors []ValidationError

	t := current.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := current.Field(i)

		// Build the full field path
		fieldName := field.Name
		if parentField != "" {
			fieldName = parentField + "." + fieldName
		}

		// Handle nested structs
		if value.Kind() == reflect.Struct {
			errors = append(errors, v.validateStruct(value, fieldName, parent)...)
		}

		// Handle slice/array of structs
		if (value.Kind() == reflect.Slice || value.Kind() == reflect.Array) && value.Type().Elem().Kind() == reflect.Struct {
			for j := 0; j < value.Len(); j++ {
				nestedErrors := v.validateStruct(value.Index(j), fmt.Sprintf("%s[%d]", fieldName, j), parent)
				errors = append(errors, nestedErrors...)
			}
		}

		// Handle map of structs
		if value.Kind() == reflect.Map && value.Type().Elem().Kind() == reflect.Struct {
			for _, key := range value.MapKeys() {
				nestedErrors := v.validateStruct(value.MapIndex(key), fmt.Sprintf("%s[%v]", fieldName, key.Interface()), parent)
				errors = append(errors, nestedErrors...)
			}
		}

		// Validate the field itself
		errors = append(errors, v.validateValue(fieldName, value, parent)...)
	}

	return errors
}

func (v *Validator) Validate(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return &ValidationError{Field: "input", Message: "input must be a struct"}
	}

	errors := v.validateStruct(val, "", data)
	if len(errors) > 0 {
		return ValidationErrors(errors)
	}

	return nil
}

type Required struct{}

func (r Required) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("value is required")
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		if v.String() == "" {
			return fmt.Errorf("value is required")
		}
	case reflect.Slice, reflect.Map:
		if v.Len() == 0 {
			return fmt.Errorf("value is required")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return fmt.Errorf("value is required")
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return fmt.Errorf("value is required")
		}
	}

	return nil
}

type StringLength struct {
	Min int
	Max int
}

func (s StringLength) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	length := len(str)
	if length < s.Min {
		return fmt.Errorf("length must be at least %d", s.Min)
	}
	if length > s.Max {
		return fmt.Errorf("length must be at most %d", s.Max)
	}

	return nil
}

type Pattern struct {
	Regex *regexp.Regexp
}

func NewPattern(pattern string) (*Pattern, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Pattern{Regex: regex}, nil
}

func (p Pattern) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if !p.Regex.MatchString(str) {
		return fmt.Errorf("value does not match pattern %s", p.Regex.String())
	}

	return nil
}
