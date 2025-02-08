package rules

import (
	"testing"
)

func TestPositive(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "positive int",
			value:   42,
			wantErr: false,
		},
		{
			name:    "zero int",
			value:   0,
			wantErr: true,
		},
		{
			name:    "negative int",
			value:   -42,
			wantErr: true,
		},
		{
			name:    "positive float",
			value:   42.5,
			wantErr: false,
		},
		{
			name:    "zero float",
			value:   0.0,
			wantErr: true,
		},
		{
			name:    "negative float",
			value:   -42.5,
			wantErr: true,
		},
		{
			name:    "non-numeric value",
			value:   "42",
			wantErr: true,
		},
		{
			name:    "uint negative",
			value:   uint(0),
			wantErr: true,
		},
	}

	rule := Positive{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Positive.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name    string
		rule    Range
		value   interface{}
		wantErr bool
	}{
		{
			name:    "int within range",
			rule:    Range{Min: 1, Max: 100},
			value:   42,
			wantErr: false,
		},
		{
			name:    "int below min",
			rule:    Range{Min: 1, Max: 100},
			value:   0,
			wantErr: true,
		},
		{
			name:    "int above max",
			rule:    Range{Min: 1, Max: 100},
			value:   101,
			wantErr: true,
		},
		{
			name:    "float within range",
			rule:    Range{Min: 1, Max: 100},
			value:   42.5,
			wantErr: false,
		},
		{
			name:    "float below min",
			rule:    Range{Min: 1, Max: 100},
			value:   0.5,
			wantErr: true,
		},
		{
			name:    "float above max",
			rule:    Range{Min: 1, Max: 100},
			value:   100.5,
			wantErr: true,
		},
		{
			name:    "non-numeric value",
			rule:    Range{Min: 1, Max: 100},
			value:   "42",
			wantErr: true,
		},
		{
			name:    "zero max value",
			rule:    Range{Min: 1, Max: 0},
			value:   42,
			wantErr: false,
		},
		{
			name:    "uint",
			rule:    Range{Min: 1, Max: 0},
			value:   uint(42),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Range.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMin_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    Min
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid int",
			rule:    Min{Value: 10},
			value:   15,
			wantErr: false,
		},
		{
			name:    "invalid int",
			rule:    Min{Value: 10},
			value:   5,
			wantErr: true,
		},
		{
			name:    "valid float",
			rule:    Min{Value: 10.5},
			value:   15.5,
			wantErr: false,
		},
		{
			name:    "invalid float",
			rule:    Min{Value: 10.5},
			value:   5.5,
			wantErr: true,
		},
		{
			name:    "valid uint",
			rule:    Min{Value: 10},
			value:   uint(15),
			wantErr: false,
		},
		{
			name:    "invalid uint",
			rule:    Min{Value: 10},
			value:   uint(5),
			wantErr: true,
		},
		{
			name:    "non-numeric value",
			rule:    Min{Value: 10},
			value:   "not a number",
			wantErr: true,
		},
		{
			name:    "nil value",
			rule:    Min{Value: 10},
			value:   nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Min.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
