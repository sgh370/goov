package validator

import (
	"testing"

	"goov/validator/rules"
)

type TestUser struct {
	Name     string `validate:"required"`
	Email    string `validate:"email"`
	Age      int    `validate:"range"`
	Username string `validate:"length"`
	Premium  bool
	Nested   *NestedStruct `validate:"required"`
}

type NestedStruct struct {
	Field string `validate:"required"`
}

func TestValidator_Validate(t *testing.T) {
	v := New()
	v.AddRule("required", rules.Required{})
	v.AddRule("email", rules.EmailDNS{})
	v.AddRule("range", rules.Range{Min: 18, Max: 100})
	v.AddRule("length", rules.Length{Min: 3, Max: 50})

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid user",
			value: &TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Username: "johndoe",
				Nested:   &NestedStruct{Field: "value"},
			},
			wantErr: false,
		},
		{
			name: "invalid user - missing required field",
			value: &TestUser{
				Email:    "john@example.com",
				Age:      25,
				Username: "johndoe",
				Nested:   &NestedStruct{Field: "value"},
			},
			wantErr: true,
		},
		{
			name: "invalid user - invalid email",
			value: &TestUser{
				Name:     "John Doe",
				Email:    "invalid-email",
				Age:      25,
				Username: "johndoe",
				Nested:   &NestedStruct{Field: "value"},
			},
			wantErr: true,
		},
		{
			name: "invalid user - age out of range",
			value: &TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      15,
				Username: "johndoe",
				Nested:   &NestedStruct{Field: "value"},
			},
			wantErr: true,
		},
		{
			name: "invalid user - username too short",
			value: &TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Username: "jo",
				Nested:   &NestedStruct{Field: "value"},
			},
			wantErr: true,
		},
		{
			name: "invalid user - nil nested struct",
			value: &TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Username: "johndoe",
				Nested:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid user - empty nested field",
			value: &TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Username: "johndoe",
				Nested:   &NestedStruct{Field: ""},
			},
			wantErr: true,
		},
		{
			name:    "non-struct value",
			value:   "not a struct",
			wantErr: true,
		},
		{
			name:    "nil value",
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateAll(t *testing.T) {
	v := New()
	v.AddRule("required", rules.Required{})
	v.AddRule("email", rules.EmailDNS{})
	v.AddRule("range", rules.Range{Min: 18, Max: 100})
	v.AddRule("length", rules.Length{Min: 3, Max: 50})

	invalidUser := &TestUser{
		Name:     "",
		Email:    "invalid-email",
		Age:      15,
		Username: "jo",
		Nested:   &NestedStruct{Field: ""},
	}

	errors := v.ValidateAll(invalidUser)
	if len(errors) != 5 {
		t.Errorf("ValidateAll() got %d errors, want 5", len(errors))
	}

	nonStruct := "not a struct"
	errors = v.ValidateAll(nonStruct)
	if len(errors) != 1 {
		t.Errorf("ValidateAll() got %d errors, want 1", len(errors))
	}
}

func TestPattern(t *testing.T) {
	pattern := NewPattern(
		rules.Required{},
		rules.Length{Min: 3, Max: 50},
		rules.EmailDNS{},
	)

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid email",
			value:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid email",
			value:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "too short value",
			value:   "a@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pattern.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pattern.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
