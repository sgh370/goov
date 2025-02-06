package rules

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"time"
)

type TimeFormat struct {
	Layout string
}

func (t TimeFormat) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	_, err := time.Parse(t.Layout, str)
	if err != nil {
		return fmt.Errorf("invalid time format: must match layout %s", t.Layout)
	}
	return nil
}

type URL struct {
	AllowedSchemes []string
}

func (u URL) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	parsed, err := url.Parse(str)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid URL format")
	}

	if len(u.AllowedSchemes) > 0 {
		valid := false
		for _, scheme := range u.AllowedSchemes {
			if parsed.Scheme == scheme {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("URL scheme must be one of: %v", u.AllowedSchemes)
		}
	}
	return nil
}

type JSON struct{}

func (j JSON) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	var js interface{}
	if err := json.Unmarshal([]byte(str), &js); err != nil {
		return fmt.Errorf("invalid JSON format")
	}
	return nil
}

type OneOf struct {
	Values []interface{}
}

func (o OneOf) Validate(value interface{}) error {
	for _, v := range o.Values {
		if reflect.DeepEqual(value, v) {
			return nil
		}
	}
	return fmt.Errorf("value must be one of: %v", o.Values)
}

type Custom struct {
	Fn func(interface{}) error
}

func (c Custom) Validate(value interface{}) error {
	return c.Fn(value)
}

// Phone validates phone numbers
type Phone struct {
	AllowEmpty bool
}

func (p Phone) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	if str == "" && p.AllowEmpty {
		return nil
	}

	// Basic phone validation: +1234567890 or 1234567890
	matched, _ := regexp.MatchString(`^\+?\d{10,15}$`, str)
	if !matched {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

// UUID validates UUID strings
type UUID struct{}

func (u UUID) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, str)
	if !matched {
		return fmt.Errorf("invalid UUID format")
	}
	return nil
}

// Date validates date strings
type Date struct {
	Format     string
	Min        time.Time
	Max        time.Time
	AllowEmpty bool
}

func (d Date) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	if str == "" && d.AllowEmpty {
		return nil
	}

	t, err := time.Parse(d.Format, str)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	if !d.Min.IsZero() && t.Before(d.Min) {
		return fmt.Errorf("date must not be before %v", d.Min.Format(d.Format))
	}

	if !d.Max.IsZero() && t.After(d.Max) {
		return fmt.Errorf("date must not be after %v", d.Max.Format(d.Format))
	}

	return nil
}

// Required validates that a value is not empty
type Required struct{}

func (r Required) Validate(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value is required")
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if v.String() == "" {
			return fmt.Errorf("value is required")
		}
	case reflect.Slice, reflect.Map:
		if v.Len() == 0 {
			return fmt.Errorf("value is required")
		}
	case reflect.Ptr:
		if v.IsNil() {
			return fmt.Errorf("value is required")
		}
	}
	return nil
}
