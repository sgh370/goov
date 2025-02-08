package main

import (
	"fmt"
	"github.com/sgh370/goov/validator"
	"github.com/sgh370/goov/validator/rules"
)

// User represents a sample user registration form
type User struct {
	Username     string   `validate:"required,min=3,max=20"`
	Email        string   `validate:"required,email"`
	Password     string   `validate:"required,min=8"`
	Age          int      `validate:"required,min=18"`
	PhoneNumber  string   `validate:"required_if=ContactMethod phone"`
	ContactMethod string  `validate:"required,oneof=email phone"`
	Website      string   `validate:"url,required_if=IsCompany true"`
	IsCompany    bool
	IPAddress    string   `validate:"ip"`
	Interests    []string `validate:"required,min=1,dive,required"`
}

func main() {
	// Create a new validator instance
	v := validator.New()

	// Add validation rules
	v.AddRule("required", rules.Required{})
	v.AddRule("min", &rules.Min{Value: 3}) // For string length and numeric values
	v.AddRule("max", &rules.Range{Max: 20})
	v.AddRule("email", &rules.EmailDNS{})
	v.AddRule("url", &rules.URL{})
	v.AddRule("ip", &rules.IP{})
	v.AddRule("oneof", &rules.OneOf{Values: []interface{}{"email", "phone"}})
	
	// Add conditional rules
	v.AddRule("required_if", &rules.If{
		Field: "ContactMethod",
		Then: rules.Required{},
	})

	// Example 1: Valid user
	validUser := User{
		Username:      "johndoe",
		Email:         "john@example.com",
		Password:      "securepass123",
		Age:          25,
		ContactMethod: "email",
		IsCompany:    false,
		Interests:    []string{"coding", "reading"},
	}

	if err := v.Validate(validUser); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Valid user registration")
	}

	// Example 2: Invalid user (missing required fields and validation failures)
	invalidUser := User{
		Username:     "jo", // too short
		Email:        "invalid-email", // invalid email format
		Password:     "short", // too short
		Age:         16, // under minimum age
		ContactMethod: "phone", // requires phone number
		Website:      "not-a-url", // invalid URL format
		IPAddress:    "256.256.256.256", // invalid IP
		Interests:    []string{}, // empty slice
	}

	if err := v.Validate(invalidUser); err != nil {
		fmt.Println("\nExpected validation errors:")
		fmt.Printf("%v\n", err)
	}

	// Example 3: Company user with conditional validations
	companyUser := User{
		Username:     "techcorp",
		Email:        "contact@techcorp.com",
		Password:     "company123secure",
		Age:         30,
		ContactMethod: "phone",
		PhoneNumber:  "+1234567890",
		IsCompany:    true,
		Website:      "https://techcorp.com",
		IPAddress:    "192.168.1.1",
		Interests:    []string{"technology", "innovation"},
	}

	if err := v.Validate(companyUser); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("\nValid company registration")
	}
}