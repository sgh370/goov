package validator

import (
	"fmt"
	"testing"

	"goov/validator/rules"
	"goov/validator/testdata"
)

func TestNestedValidation(t *testing.T) {
	v := New()

	// Address validation rules
	v.AddRule("BillingAddr.Street", rules.Length{Min: 5, Max: 100})
	v.AddRule("BillingAddr.City", rules.Length{Min: 2, Max: 50})
	v.AddRule("BillingAddr.Country", rules.Length{Min: 2, Max: 50})
	
	// OrderItem validation rules
	v.AddRule("Items.ProductID", rules.Positive{})
	v.AddRule("Items.Quantity", rules.Range{Min: 1, Max: 100})
	v.AddRule("Items.UnitPrice", rules.Range{Min: 0.01, Max: 1000000})
	
	// Map validation rules
	v.AddRule("Contacts.Value", rules.Length{Min: 5, Max: 100})

	tests := []struct {
		name    string
		order   testdata.Order
		wantErr bool
	}{
		{
			name: "valid order",
			order: testdata.Order{
				ID:         1,
				CustomerID: 100,
				Status:     "pending",
				Items: []testdata.OrderItem{
					{ProductID: 1, Quantity: 2, UnitPrice: 10.99},
				},
				BillingAddr: testdata.Address{
					Street:  "123 Main St",
					City:    "New York",
					Country: "USA",
				},
				Contacts: map[string]testdata.Contact{
					"email": {Type: "email", Value: "customer@example.com"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid nested address",
			order: testdata.Order{
				BillingAddr: testdata.Address{
					Street: "123", // too short
					City:   "NY",  // too short
				},
			},
			wantErr: true,
		},
		{
			name: "invalid order item",
			order: testdata.Order{
				Items: []testdata.OrderItem{
					{ProductID: -1, Quantity: 0}, // invalid product ID and quantity
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConditionalValidation(t *testing.T) {
	v := New()

	// Add conditional validation
	v.AddRule("PremiumDetails", rules.When{
		Condition: func(value interface{}) bool {
			user, ok := value.(*testdata.User)
			return ok && user.Premium
		},
		Then: rules.Length{Min: 1, Max: 100},
	})

	// Add cross-field validation
	v.AddRule("ConfirmPassword", rules.CrossField{
		Field:  "Password",
		Parent: &testdata.User{},
		ValidateFn: func(value, crossValue interface{}) error {
			user, ok := value.(*testdata.User)
			if !ok {
				return fmt.Errorf("expected *testdata.User, got %T", value)
			}
			if user.ConfirmPassword != user.Password {
				return fmt.Errorf("passwords do not match")
			}
			return nil
		},
	})

	tests := []struct {
		name    string
		user    *testdata.User
		wantErr bool
	}{
		{
			name: "valid premium user",
			user: &testdata.User{
				Username:        "john",
				Premium:         true,
				PremiumDetails: "VIP Member",
				Password:       "secret",
				ConfirmPassword: "secret",
			},
			wantErr: false,
		},
		{
			name: "invalid premium user - missing details",
			user: &testdata.User{
				Username: "john",
				Premium:  true,
				Password: "secret",
				ConfirmPassword: "secret",
			},
			wantErr: true,
		},
		{
			name: "valid non-premium user - no details required",
			user: &testdata.User{
				Username: "john",
				Premium:  false,
				Password: "secret",
				ConfirmPassword: "secret",
			},
			wantErr: false,
		},
		{
			name: "invalid password confirmation",
			user: &testdata.User{
				Username:        "john",
				Password:       "secret",
				ConfirmPassword: "different",
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
