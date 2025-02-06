package rules

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
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
