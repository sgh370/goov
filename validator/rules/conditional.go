package rules

import (
	"fmt"
	"reflect"
)

type When struct {
	Condition func(interface{}) bool
	Then      ValidationRule
	Else      ValidationRule
	parent    interface{}
}

func (w *When) SetParent(parent interface{}) {
	w.parent = parent
	// Also set parent for Then and Else rules if they implement SetParent
	if w.Then != nil {
		if setter, ok := w.Then.(interface{ SetParent(interface{}) }); ok {
			setter.SetParent(parent)
		}
	}
	if w.Else != nil {
		if setter, ok := w.Else.(interface{ SetParent(interface{}) }); ok {
			setter.SetParent(parent)
		}
	}
}

func (w When) Validate(value interface{}) error {
	if w.parent == nil {
		return fmt.Errorf("parent not set")
	}
	
	if w.Condition(w.parent) {
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
	ValidateFn func(parent, value interface{}) error
	parent     interface{}
}

func (c *CrossField) SetParent(parent interface{}) {
	c.parent = parent
}

func (c CrossField) Validate(value interface{}) error {
	if c.ValidateFn == nil {
		return fmt.Errorf("validation function not provided")
	}

	if c.parent == nil {
		return fmt.Errorf("parent not set")
	}

	return c.ValidateFn(c.parent, value)
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
