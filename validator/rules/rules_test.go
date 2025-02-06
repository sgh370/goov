package rules

import (
	"fmt"
	"testing"
	"time"
)

func TestLength(t *testing.T) {
	tests := []struct {
		name    string
		rule    Length
		value   interface{}
		wantErr bool
	}{
		{"valid slice", Length{Min: 1, Max: 3}, []int{1, 2}, false},
		{"valid string", Length{Min: 1, Max: 3}, "ab", false},
		{"valid map", Length{Min: 1, Max: 3}, map[string]int{"a": 1}, false},
		{"too short", Length{Min: 2, Max: 3}, []int{1}, true},
		{"too long", Length{Min: 1, Max: 2}, []int{1, 2, 3}, true},
		{"invalid type", Length{Min: 1, Max: 3}, 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Length.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEach(t *testing.T) {
	positiveRule := Positive{}
	tests := []struct {
		name    string
		rule    Each
		value   interface{}
		wantErr bool
	}{
		{"valid numbers", Each{Rule: positiveRule}, []int{1, 2, 3}, false},
		{"invalid numbers", Each{Rule: positiveRule}, []int{1, -2, 3}, true},
		{"invalid type", Each{Rule: positiveRule}, 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Each.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeFormat(t *testing.T) {
	tests := []struct {
		name    string
		rule    TimeFormat
		value   interface{}
		wantErr bool
	}{
		{"valid RFC3339", TimeFormat{Layout: time.RFC3339}, "2025-02-06T17:00:00Z", false},
		{"invalid format", TimeFormat{Layout: time.RFC3339}, "2025-02-06", true},
		{"invalid type", TimeFormat{Layout: time.RFC3339}, 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("TimeFormat.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name    string
		rule    URL
		value   interface{}
		wantErr bool
	}{
		{"valid http", URL{AllowedSchemes: []string{"http", "https"}}, "http://example.com", false},
		{"valid https", URL{AllowedSchemes: []string{"http", "https"}}, "https://example.com", false},
		{"invalid scheme", URL{AllowedSchemes: []string{"https"}}, "http://example.com", true},
		{"invalid format", URL{}, "not-a-url", true},
		{"invalid type", URL{}, 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("URL.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid object", `{"key": "value"}`, false},
		{"valid array", `[1, 2, 3]`, false},
		{"invalid json", `{key: value}`, true},
		{"invalid type", 123, true},
	}

	j := JSON{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := j.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("JSON.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	tests := []struct {
		name    string
		rule    OneOf
		value   interface{}
		wantErr bool
	}{
		{"valid string", OneOf{Values: []interface{}{"a", "b", "c"}}, "b", false},
		{"valid int", OneOf{Values: []interface{}{1, 2, 3}}, 2, false},
		{"invalid value", OneOf{Values: []interface{}{"a", "b", "c"}}, "d", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("OneOf.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustom(t *testing.T) {
	evenNumber := Custom{
		Fn: func(value interface{}) error {
			num, ok := value.(int)
			if !ok {
				return fmt.Errorf("value must be an integer")
			}
			if num%2 != 0 {
				return fmt.Errorf("value must be even")
			}
			return nil
		},
	}

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid even", 2, false},
		{"invalid odd", 3, true},
		{"invalid type", "2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := evenNumber.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Custom.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	tests := []struct {
		name      string
		phone     Phone
		value     interface{}
		wantError bool
	}{
		{
			name:      "valid phone",
			phone:     Phone{},
			value:     "+1234567890",
			wantError: false,
		},
		{
			name:      "valid phone without plus",
			phone:     Phone{},
			value:     "1234567890",
			wantError: false,
		},
		{
			name:      "invalid phone - too short",
			phone:     Phone{},
			value:     "123456",
			wantError: true,
		},
		{
			name:      "invalid phone - contains letters",
			phone:     Phone{},
			value:     "+1234abc890",
			wantError: true,
		},
		{
			name:      "empty phone - not allowed",
			phone:     Phone{AllowEmpty: false},
			value:     "",
			wantError: true,
		},
		{
			name:      "empty phone - allowed",
			phone:     Phone{AllowEmpty: true},
			value:     "",
			wantError: false,
		},
		{
			name:      "invalid type",
			phone:     Phone{},
			value:     123,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.phone.Validate(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("Phone.Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		name      string
		uuid      UUID
		value     interface{}
		wantError bool
	}{
		{
			name:      "valid UUID",
			uuid:      UUID{},
			value:     "123e4567-e89b-12d3-a456-426614174000",
			wantError: false,
		},
		{
			name:      "invalid UUID - wrong format",
			uuid:      UUID{},
			value:     "123e4567-e89b-12d3-a456",
			wantError: true,
		},
		{
			name:      "invalid UUID - contains invalid chars",
			uuid:      UUID{},
			value:     "123e4567-e89b-12d3-a456-42661417400g",
			wantError: true,
		},
		{
			name:      "invalid type",
			uuid:      UUID{},
			value:     123,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.uuid.Validate(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("UUID.Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDate(t *testing.T) {
	format := "2006-01-02"
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		date      Date
		value     interface{}
		wantError bool
	}{
		{
			name: "valid date",
			date: Date{Format: format},
			value: "2023-01-01",
			wantError: false,
		},
		{
			name: "valid date with range",
			date: Date{Format: format, Min: min, Max: max},
			value: "2023-01-01",
			wantError: false,
		},
		{
			name: "invalid date - before min",
			date: Date{Format: format, Min: min},
			value: "2019-12-31",
			wantError: true,
		},
		{
			name: "invalid date - after max",
			date: Date{Format: format, Max: max},
			value: "2026-01-01",
			wantError: true,
		},
		{
			name: "invalid date format",
			date: Date{Format: format},
			value: "2023/01/01",
			wantError: true,
		},
		{
			name: "empty date - not allowed",
			date: Date{Format: format, AllowEmpty: false},
			value: "",
			wantError: true,
		},
		{
			name: "empty date - allowed",
			date: Date{Format: format, AllowEmpty: true},
			value: "",
			wantError: false,
		},
		{
			name: "invalid type",
			date: Date{Format: format},
			value: 123,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.date.Validate(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("Date.Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name      string
		required  Required
		value     interface{}
		wantError bool
	}{
		{
			name:      "non-empty string",
			required:  Required{},
			value:     "test",
			wantError: false,
		},
		{
			name:      "empty string",
			required:  Required{},
			value:     "",
			wantError: true,
		},
		{
			name:      "non-empty slice",
			required:  Required{},
			value:     []int{1, 2, 3},
			wantError: false,
		},
		{
			name:      "empty slice",
			required:  Required{},
			value:     []int{},
			wantError: true,
		},
		{
			name:      "non-empty map",
			required:  Required{},
			value:     map[string]int{"a": 1},
			wantError: false,
		},
		{
			name:      "empty map",
			required:  Required{},
			value:     map[string]int{},
			wantError: true,
		},
		{
			name:      "non-nil pointer",
			required:  Required{},
			value:     &struct{}{},
			wantError: false,
		},
		{
			name:      "nil pointer",
			required:  Required{},
			value:     (*struct{})(nil),
			wantError: true,
		},
		{
			name:      "nil value",
			required:  Required{},
			value:     nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.required.Validate(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("Required.Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
