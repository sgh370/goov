package validator

import (
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
	v.AddRule("Name", StringLength{Min: 3, Max: 50})

	// Price validation
	v.AddRule("Price", rules.Range{Min: 0.01, Max: 1000000})

	// Tags validation
	v.AddRule("Tags", rules.Length{Min: 1, Max: 10})
	v.AddRule("Tags", rules.Each{Rule: StringLength{Min: 2, Max: 20}})
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
	v.AddRule("Categories", rules.Contains{Value: "default"})

	// Attributes validation
	v.AddRule("Attributes", rules.Length{Min: 0, Max: 20})

	// Status validation
	v.AddRule("Status", rules.OneOf{Values: []interface{}{"draft", "published", "archived"}})

	tests := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: Product{
				ID:          1,
				Name:        "Test Product",
				Price:       99.99,
				Tags:        []string{"electronics", "gadget"},
				CreatedAt:   "2025-02-06T17:00:00Z",
				UpdatedAt:   "2025-02-06T17:00:00Z",
				URL:         "https://example.com/product",
				Metadata:    `{"color": "blue", "size": "medium"}`,
				Categories:  []string{"default", "electronics"},
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
			product: Product{
				ID:          0, // invalid: must be positive
				Name:        "T", // invalid: too short
				Price:       0, // invalid: too low
				Tags:        []string{"a"}, // invalid: tag too short
				CreatedAt:   "invalid-date", // invalid: wrong format
				UpdatedAt:   "invalid-date", // invalid: wrong format
				URL:         "ftp://example.com", // invalid: wrong scheme
				Metadata:    "{invalid-json}", // invalid: wrong JSON
				Categories:  []string{"electronics"}, // invalid: missing default category
				Attributes: map[string]string{}, // valid: empty map allowed
				Status:     "pending", // invalid: not in allowed values
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
