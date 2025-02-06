package rules

import (
	"fmt"
	"reflect"
)

type ValidationRule interface {
	Validate(value interface{}) error
}

type Length struct {
	Min int
	Max int
}

func (l Length) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		length := v.Len()
		if length < l.Min {
			return fmt.Errorf("length must be at least %d", l.Min)
		}
		if l.Max > 0 && length > l.Max {
			return fmt.Errorf("length must not exceed %d", l.Max)
		}
		return nil
	default:
		return fmt.Errorf("value must be a slice, array, map, or string")
	}
}

type Each struct {
	Rule ValidationRule
}

func (e Each) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("value must be a slice or array")
	}

	for i := 0; i < v.Len(); i++ {
		if err := e.Rule.Validate(v.Index(i).Interface()); err != nil {
			return fmt.Errorf("item at index %d: %v", i, err)
		}
	}
	return nil
}

type Contains struct {
	Value interface{}
}

func (c Contains) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("value must be a slice or array")
	}

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), c.Value) {
			return nil
		}
	}
	return fmt.Errorf("value must contain %v", c.Value)
}

type Unique struct{}

func (u Unique) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("value must be a slice or array")
	}

	seen := make(map[interface{}]bool)
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		if seen[item] {
			return fmt.Errorf("duplicate value found: %v", item)
		}
		seen[item] = true
	}
	return nil
}
