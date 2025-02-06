package validator

import (
	"testing"

	"goov/validator/rules"
	"goov/validator/testdata"
)

func TestValidator(t *testing.T) {
	v := New()

	v.AddRule("Username", StringLength{Min: 3, Max: 20})
	v.AddRule("Age", rules.Range{Min: 18, Max: 100})
	emailPattern, _ := NewPattern(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	v.AddRule("Email", emailPattern)

	tests := []struct {
		name    string
		user    testdata.TestUser
		wantErr bool
	}{
		{
			name: "valid user",
			user: testdata.TestUser{
				Username: "johndoe",
				Age:      25,
				Email:    "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid username",
			user: testdata.TestUser{
				Username: "jo",
				Age:      25,
				Email:    "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid age",
			user: testdata.TestUser{
				Username: "johndoe",
				Age:      15,
				Email:    "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			user: testdata.TestUser{
				Username: "johndoe",
				Age:      25,
				Email:    "invalid-email",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
