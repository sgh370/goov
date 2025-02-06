package rules

import (
	"fmt"
	"reflect"
)

type Range struct {
	Min float64
	Max float64
}

func (r Range) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	var num float64

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		num = v.Float()
	default:
		return fmt.Errorf("value must be numeric")
	}

	if num < r.Min {
		return fmt.Errorf("value must be greater than or equal to %v", r.Min)
	}
	if r.Max > 0 && num > r.Max {
		return fmt.Errorf("value must be less than or equal to %v", r.Max)
	}
	return nil
}

type Positive struct{}

func (p Positive) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() <= 0 {
			return fmt.Errorf("value must be positive")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 0 {
			return fmt.Errorf("value must be positive")
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() <= 0 {
			return fmt.Errorf("value must be positive")
		}
	default:
		return fmt.Errorf("value must be numeric")
	}
	return nil
}
