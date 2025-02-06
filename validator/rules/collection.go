package rules

import (
	"fmt"
	"reflect"
)

type Rule interface {
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
	Rule Rule
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

type Map struct {
	Key   Rule
	Value Rule
}

func (m Map) Validate(value interface{}) error {
	if value == nil {
		return fmt.Errorf("map is nil")
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map {
		return fmt.Errorf("value is not a map")
	}

	for _, key := range v.MapKeys() {
		if m.Key != nil {
			if err := m.Key.Validate(key.Interface()); err != nil {
				return fmt.Errorf("invalid map key: %v", err)
			}
		}

		if m.Value != nil {
			if err := m.Value.Validate(v.MapIndex(key).Interface()); err != nil {
				return fmt.Errorf("invalid map value for key %v: %v", key.Interface(), err)
			}
		}
	}

	return nil
}

type Slice struct {
	Rule Rule
}

func (s Slice) Validate(value interface{}) error {
	if value == nil {
		return fmt.Errorf("slice is nil")
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("value is not a slice")
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if s.Rule != nil {
			if err := s.Rule.Validate(item.Interface()); err != nil {
				return fmt.Errorf("invalid item at index %d: %v", i, err)
			}
		}
	}

	return nil
}

type EachMulti struct {
	Rules []Rule
}

func (e EachMulti) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("value must be a slice or array")
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		for _, rule := range e.Rules {
			if err := rule.Validate(item); err != nil {
				return fmt.Errorf("item at index %d failed validation: %v", i, err)
			}
		}
	}

	return nil
}

type Keys struct {
	Rules []Rule
}

func (k Keys) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map {
		return fmt.Errorf("value must be a map")
	}

	keys := v.MapKeys()
	for _, rule := range k.Rules {
		for _, key := range keys {
			if err := rule.Validate(key.Interface()); err != nil {
				return fmt.Errorf("map key %v failed validation: %v", key.Interface(), err)
			}
		}
	}

	return nil
}
