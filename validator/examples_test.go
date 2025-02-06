package validator

import (
	"fmt"
	"testing"
	"time"

	"goov/validator/rules"
)

type Product struct {
	ID          int
	Name        string
	Price       float64
	Tags        []string
	CreatedAt   string
	UpdatedAt   string
	URL         string
	Metadata    string
	Categories  []string
	Attributes  map[string]string
	Status      string
}

func TestComplexValidation(t *testing.T) {
	v := New()

	// ID validation
	v.AddRule("ID", rules.Positive{})

	// Name validation
	v.AddRule("Name", rules.Length{Min: 3, Max: 50})

	// Price validation
	v.AddRule("Price", rules.Range{Min: 0.01, Max: 1000000})

	// Tags validation
	v.AddRule("Tags", rules.Length{Min: 1, Max: 10})
	v.AddRule("Tags", rules.Each{Rule: rules.Length{Min: 2, Max: 20}})
	v.AddRule("Tags", rules.Unique{})

	// Time validation
	v.AddRule("CreatedAt", rules.TimeFormat{Layout: time.RFC3339})
	v.AddRule("UpdatedAt", rules.TimeFormat{Layout: time.RFC3339})

	// URL validation
	v.AddRule("URL", rules.URL{AllowedSchemes: []string{"http", "https"}})

	// JSON validation
	v.AddRule("Metadata", rules.JSON{})

	// Categories validation
	v.AddRule("Categories", rules.Length{Min: 1, Max: 5})
	v.AddRule("Categories", rules.Each{Rule: rules.Length{Min: 2, Max: 20}})
	v.AddRule("Categories", rules.Unique{})

	// Attributes validation
	v.AddRule("Attributes", rules.Length{Min: 1, Max: 10})

	// Status validation
	v.AddRule("Status", rules.OneOf{Values: []interface{}{"draft", "published", "archived"}})

	tests := []struct {
		name    string
		product *Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: &Product{
				ID:          1,
				Name:        "Test Product",
				Price:       99.99,
				Tags:        []string{"electronics", "gadget"},
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
				URL:         "https://example.com/product",
				Metadata:    `{"color": "blue", "size": "large"}`,
				Categories:  []string{"electronics", "accessories"},
				Attributes: map[string]string{
					"brand": "TestBrand",
					"model": "X100",
				},
				Status: "published",
			},
			wantErr: false,
		},
		{
			name: "invalid product",
			product: &Product{
				ID:          -1,
				Name:        "A",
				Price:       0,
				Tags:        []string{"a"},
				CreatedAt:   "invalid-time",
				UpdatedAt:   "invalid-time",
				URL:         "invalid-url",
				Metadata:    "invalid-json",
				Categories:  []string{"a"},
				Attributes:  map[string]string{},
				Status:      "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.product)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Example_basicValidation() {
	// Create a struct to validate
	type Product struct {
		ID    int
		Name  string
		Price float64
		Tags  []string
		Time  string
	}

	// Create a new validator
	v := New()

	// ID validation
	v.AddRule("ID", rules.Positive{})

	// Name validation
	v.AddRule("Name", rules.Length{Min: 3, Max: 50})

	// Price validation
	v.AddRule("Price", rules.Range{Min: 0.01, Max: 1000000})

	// Tags validation
	v.AddRule("Tags", rules.Length{Min: 1, Max: 10})
	v.AddRule("Tags", rules.Each{Rule: rules.Length{Min: 2, Max: 20}})

	// Time validation
	v.AddRule("Time", rules.TimeFormat{Layout: time.RFC3339})

	// Validate a valid product
	validProduct := &Product{
		ID:    1,
		Name:  "Test Product",
		Price: 99.99,
		Tags:  []string{"electronics", "gadget"},
		Time:  time.Now().Format(time.RFC3339),
	}

	if err := v.Validate(validProduct); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Valid product")
	}

	// Validate an invalid product
	invalidProduct := &Product{
		ID:    -1,
		Name:  "A",
		Price: 0,
		Tags:  []string{"a"},
		Time:  "invalid-time",
	}

	if err := v.Validate(invalidProduct); err != nil {
		fmt.Println("Invalid product with errors:")
		for _, e := range err.(ValidationErrors) {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	}

	// Output:
	// Valid product
	// Invalid product with errors:
	// - ID: value must be positive
	// - Name: length must be at least 3
	// - Price: value must be greater than or equal to 0.01
	// - Tags: item at index 0: length must be at least 2
	// - Time: invalid time format: must match layout 2006-01-02T15:04:05Z07:00
}

func Example_advancedValidation() {
	// Create a struct to validate
	type User struct {
		ID            string
		Name          string
		Phone         string
		Email         string
		DateOfBirth   string
		Premium       bool
		PremiumUntil  string
		ContactMethod string
	}

	// Create a new validator
	v := New()

	// Add validation rules
	v.AddRule("ID", rules.UUID{})
	v.AddRule("Name", rules.Required{})
	v.AddRule("Phone", rules.Phone{AllowEmpty: true})
	v.AddRule("Email", rules.Required{})
	emailPattern, _ := NewPattern(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	v.AddRule("Email", emailPattern)

	// Add date validation with range
	minAge := time.Now().AddDate(-100, 0, 0)
	maxAge := time.Now().AddDate(-18, 0, 0)
	v.AddRule("DateOfBirth", rules.Date{
		Format: "2006-01-02",
		Min:    minAge,
		Max:    maxAge,
	})

	// Add conditional validation for premium users
	v.AddRule("PremiumUntil", &rules.When{
		Condition: func(value interface{}) bool {
			if user, ok := value.(*User); ok {
				return user.Premium
			}
			return false
		},
		Then: rules.Date{
			Format: "2006-01-02",
			Min:    time.Now(),
		},
	})

	// Add OneOf validation for contact method
	v.AddRule("ContactMethod", rules.OneOf{
		Values: []interface{}{"email", "phone", "both"},
	})

	// Validate a valid user
	validUser := &User{
		ID:            "123e4567-e89b-12d3-a456-426614174000",
		Name:          "John Doe",
		Phone:         "+1234567890",
		Email:         "john@example.com",
		DateOfBirth:   "1990-01-01",
		Premium:       true,
		PremiumUntil:  "2024-12-31",
		ContactMethod: "email",
	}

	if err := v.Validate(validUser); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Valid user")
	}

	// Validate an invalid user
	invalidUser := &User{
		ID:            "invalid-uuid",
		Name:          "",
		Phone:         "invalid-phone",
		Email:         "invalid-email",
		DateOfBirth:   "2010-01-01", // Too young
		Premium:       true,
		PremiumUntil:  "2020-12-31", // Past date
		ContactMethod: "invalid",
	}

	if err := v.Validate(invalidUser); err != nil {
		fmt.Println("Invalid user with errors:")
		for _, e := range err.(ValidationErrors) {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	}

	// Output:
	// Validation error: PremiumUntil: date must not be before 2025-02-06
	// Invalid user with errors:
	// - ID: invalid UUID format
	// - Name: value is required
	// - Phone: invalid phone number format
	// - Email: value does not match pattern ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	// - DateOfBirth: date must not be after 2007-02-06
	// - PremiumUntil: date must not be before 2025-02-06
	// - ContactMethod: value must be one of: [email phone both]
}
