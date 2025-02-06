package rules

import (
	"fmt"
	"testing"
	"time"
)

func TestRange(t *testing.T) {
	tests := []struct {
		name    string
		rule    Range
		value   interface{}
		wantErr bool
	}{
		{"valid int", Range{Min: 0, Max: 10}, 5, false},
		{"valid float", Range{Min: 0, Max: 10}, 5.5, false},
		{"below min", Range{Min: 0, Max: 10}, -1, true},
		{"above max", Range{Min: 0, Max: 10}, 11, true},
		{"invalid type", Range{Min: 0, Max: 10}, "string", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Range.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
