package validator

import (
	"fmt"

	"goov/validator/rules"
)

func Example_basicValidation() {
	v := New()

	// Add validation rules
	v.AddRule("required", rules.Required{})
	v.AddRule("email", rules.EmailDNS{})
	v.AddRule("length", rules.Length{Min: 3, Max: 50})

	// Create a struct with validation tags
	type User struct {
		Name     string `validate:"required"`
		Email    string `validate:"email"`
		Username string `validate:"length"`
	}

	// Create an invalid user
	invalidUser := User{
		Name:     "",
		Email:    "invalid-email",
		Username: "jo",
	}

	// Validate the user
	if err := v.Validate(invalidUser); err != nil {
		fmt.Printf("Validation error: %s\n", err)
	}

	// Output:
	// Validation error: Name: value is required
}

func Example_advancedValidation() {
	v := New()

	// Add validation rules
	v.AddRule("required", rules.Required{})
	v.AddRule("email", rules.EmailDNS{})
	v.AddRule("domain", rules.Domain{})
	v.AddRule("port", rules.Port{Min: 1, Max: 65535})

	// Create a struct with validation tags
	type Server struct {
		Name     string `validate:"required"`
		Email    string `validate:"email"`
		Domain   string `validate:"domain"`
		Port     int    `validate:"port"`
	}

	// Create an invalid server
	invalidServer := Server{
		Name:     "",
		Email:    "invalid-email",
		Domain:   "invalid-domain",
		Port:     0,
	}

	// Validate the server
	if err := v.Validate(invalidServer); err != nil {
		fmt.Printf("Validation error: %s\n", err)
	}

	// Output:
	// Validation error: Name: value is required
}

func Example_advancedValidationRules() {
	v := New()

	// Add validation rules
	v.AddRule("ip", rules.IP{})
	v.AddRule("domain", rules.Domain{})
	v.AddRule("port", rules.Port{Min: 1, Max: 65535})

	// Create a struct with validation tags
	type Server struct {
		IP     string `validate:"ip"`
		Domain string `validate:"domain"`
		Port   int    `validate:"port"`
	}

	// Create an invalid server
	invalidServer := Server{
		IP:     "invalid-ip",
		Domain: "invalid-domain",
		Port:   0,
	}

	// Validate the server
	if err := v.Validate(invalidServer); err != nil {
		fmt.Printf("Validation error: %s\n", err)
	}

	// Output:
	// Validation error: IP: invalid IP address format
}

func Example_networkAndFormatValidation() {
	v := New()

	// Add validation rules
	v.AddRule("ip", rules.IP{})
	v.AddRule("domain", rules.Domain{})
	v.AddRule("port", rules.Port{Min: 1, Max: 65535})

	// Create a struct with validation tags
	type NetworkConfig struct {
		IP     string `validate:"ip"`
		Domain string `validate:"domain"`
		Port   int    `validate:"port"`
	}

	// Create an invalid config
	invalidConfig := NetworkConfig{
		IP:     "invalid-ip",
		Domain: "invalid-domain",
		Port:   0,
	}

	// Validate the config
	if err := v.Validate(invalidConfig); err != nil {
		fmt.Printf("Validation error: %s\n", err)
	}

	// Output:
	// Validation error: IP: invalid IP address format
}
