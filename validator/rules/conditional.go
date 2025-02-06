package rules

import (
	"fmt"
	"reflect"
)

type When struct {
	Condition func(interface{}) bool
	Then      ValidationRule
	Else      ValidationRule
}

func (w When) Validate(value interface{}) error {
	if w.Condition(value) {
		if w.Then != nil {
			return w.Then.Validate(value)
		}
	} else if w.Else != nil {
		return w.Else.Validate(value)
	}
	return nil
}

type CrossField struct {
	Field      string
	Parent     interface{}
	ValidateFn func(value, crossValue interface{}) error
}

func (c *CrossField) SetParent(parent interface{}) {
	c.Parent = parent
}

func (c CrossField) Validate(fieldValue interface{}) error {
	if c.ValidateFn == nil {
		return fmt.Errorf("validation function not provided")
	}

	// Get the parent value
	parentVal := reflect.ValueOf(c.Parent)
	if parentVal.Kind() == reflect.Ptr {
		parentVal = parentVal.Elem()
	}

	if parentVal.Kind() != reflect.Struct {
		return fmt.Errorf("parent must be a struct")
	}

	// Get the cross-field value
	crossField := parentVal.FieldByName(c.Field)
	if !crossField.IsValid() {
		return fmt.Errorf("field %s not found", c.Field)
	}

	return c.ValidateFn(fieldValue, crossField.Interface())
}

type DependentRequired struct {
	Field  string
	Parent interface{}
}

func (d DependentRequired) Validate(value interface{}) error {
	if d.Parent == nil {
		return fmt.Errorf("parent struct not provided")
	}

	parentVal := reflect.ValueOf(d.Parent)
	if parentVal.Kind() == reflect.Ptr {
		parentVal = parentVal.Elem()
	}

	if parentVal.Kind() != reflect.Struct {
		return fmt.Errorf("parent must be a struct")
	}

	field := parentVal.FieldByName(d.Field)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", d.Field)
	}

	// Check if the field is zero value
	if field.IsZero() {
		return fmt.Errorf("field %s is required", d.Field)
	}

	return nil
}
