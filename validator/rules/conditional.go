package rules

import (
	"fmt"
	"reflect"
)

type When struct {
	Condition func(interface{}) bool
	Then      Rule
	Else      Rule
	parent    interface{}
}

func (w *When) SetParent(parent interface{}) {
	w.parent = parent
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
	if w.Condition(w.parent) {
		if w.Then != nil {
			return w.Then.Validate(value)
		}
	} else if w.Else != nil {
		return w.Else.Validate(value)
	}
	return nil
}

type If struct {
	Field  string
	Then   Rule
	Else   Rule
	parent interface{}
}

func (i *If) SetParent(parent interface{}) {
	i.parent = parent
	if i.Then != nil {
		if setter, ok := i.Then.(interface{ SetParent(interface{}) }); ok {
			setter.SetParent(parent)
		}
	}
	if i.Else != nil {
		if setter, ok := i.Else.(interface{ SetParent(interface{}) }); ok {
			setter.SetParent(parent)
		}
	}
}

func (i If) Validate(value interface{}) error {
	if i.parent == nil {
		return fmt.Errorf("parent not set")
	}

	v := reflect.ValueOf(i.parent)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("parent must be a struct")
	}

	field := v.FieldByName(i.Field)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", i.Field)
	}

	if field.Kind() != reflect.Bool {
		return fmt.Errorf("field %s is not a boolean", i.Field)
	}

	if field.Bool() {
		if i.Then != nil {
			return i.Then.Validate(value)
		}
	} else if i.Else != nil {
		return i.Else.Validate(value)
	}
	return nil
}

type Unless struct {
	Field  string
	Then   Rule
	Else   Rule
	parent interface{}
}

func (u *Unless) SetParent(parent interface{}) {
	u.parent = parent
	if u.Then != nil {
		if setter, ok := u.Then.(interface{ SetParent(interface{}) }); ok {
			setter.SetParent(parent)
		}
	}
	if u.Else != nil {
		if setter, ok := u.Else.(interface{ SetParent(interface{}) }); ok {
			setter.SetParent(parent)
		}
	}
}

func (u Unless) Validate(value interface{}) error {
	if u.parent == nil {
		return fmt.Errorf("parent not set")
	}

	v := reflect.ValueOf(u.parent)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("parent must be a struct")
	}

	field := v.FieldByName(u.Field)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", u.Field)
	}

	if field.Kind() != reflect.Bool {
		return fmt.Errorf("field %s is not a boolean", u.Field)
	}

	if !field.Bool() {
		if u.Then != nil {
			return u.Then.Validate(value)
		}
	} else if u.Else != nil {
		return u.Else.Validate(value)
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
