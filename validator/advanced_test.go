package validator

import (
	"testing"

	"github.com/sgh370/goov/validator/rules"
)

type Address struct {
	Street  string `validate:"required"`
	City    string `validate:"required"`
	Country string `validate:"required"`
	ZIP     string `validate:"required"`
}

type OrderItem struct {
	ProductID string `validate:"required"`
	Quantity  int    `validate:"min=0"`
}

type Order struct {
	ID    string      `validate:"required"`
	Items []OrderItem `validate:"slice=required"`
}

func TestNestedValidation(t *testing.T) {
	v := New()
	v.AddRule("required", rules.Required{})
	v.AddRule("min", rules.Min{Value: 0})
	v.AddRule("slice", rules.Slice{Rule: rules.Required{}})

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid order",
			value: &Order{
				ID: "123",
				Items: []OrderItem{
					{
						ProductID: "P1",
						Quantity:  1,
					},
					{
						ProductID: "P2",
						Quantity:  2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid order - no items",
			value: &Order{
				ID:    "123",
				Items: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid order item",
			value: &Order{
				ID: "123",
				Items: []OrderItem{
					{
						ProductID: "",  // Invalid: required
						Quantity:  -1,  // Invalid: min=0
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type PremiumUser struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Premium  bool
	Details  *UserDetails `validate:"if"`
}

type UserDetails struct {
	Email     string `validate:"required,email"`
	Phone     string `validate:"required"`
	BillingID string `validate:"required"`
}

func TestConditionalValidation(t *testing.T) {
	v := New()
	v.AddRule("required", rules.Required{})
	v.AddRule("email", rules.EmailDNS{})
	v.AddRule("if", &rules.If{
		Field: "Premium",
		Then:  rules.Required{},
	})

	tests := []struct {
		name    string
		user    *PremiumUser
		wantErr bool
	}{
		{
			name: "valid premium user",
			user: &PremiumUser{
				Username: "john_doe",
				Password: "password123",
				Premium:  true,
				Details: &UserDetails{
					Email:     "john@example.com",
					Phone:     "123-456-7890",
					BillingID: "B123",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid premium user - missing details",
			user: &PremiumUser{
				Username: "john_doe",
				Password: "password123",
				Premium:  true,
				Details:  nil,
			},
			wantErr: true,
		},
		{
			name: "valid non-premium user - no details required",
			user: &PremiumUser{
				Username: "john_doe",
				Password: "password123",
				Premium:  false,
				Details:  nil,
			},
			wantErr: false,
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
