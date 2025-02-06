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
		ID:            "not-a-uuid",
		Phone:         "not-a-phone",
		Email:         "not-an-email",
		DateOfBirth:   "2010-01-01",
		Premium:       true,
		PremiumUntil:  "2024-01-01",
		ContactMethod: "invalid",
	}

	if err := v.Validate(invalidUser); err != nil {
		fmt.Println("Invalid user with errors:")
		for _, e := range err.(ValidationErrors) {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	}

	// Output:
	// Validation error: PremiumUntil: date must not be before 2025-02-07
	// Invalid user with errors:
	// - ID: invalid UUID format
	// - Name: value is required
	// - Phone: invalid phone number format
	// - Email: value does not match pattern ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	// - DateOfBirth: date must not be after 2007-02-07
	// - PremiumUntil: date must not be before 2025-02-07
	// - ContactMethod: value must be one of: [email phone both]
}

func Example_advancedValidationRules() {
	// Create a struct to validate
	type Server struct {
		IPAddress  string
		Domain     string
		Password   string
		BackupCard string
	}

	// Create a new validator
	v := New()

	// Add IP validation
	v.AddRule("IPAddress", rules.IP{
		AllowV4: true,
		AllowV6: true,
	})

	// Add domain validation
	v.AddRule("Domain", rules.Domain{
		AllowSubdomains: true,
	})

	// Add password validation
	v.AddRule("Password", rules.Password{
		MinLength:      12,
		MaxLength:      50,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: true,
	})

	// Add credit card validation
	v.AddRule("BackupCard", rules.CreditCard{
		AllowEmpty: true,
	})

	// Validate a valid server
	validServer := &Server{
		IPAddress:  "192.168.1.1",
		Domain:     "example.com",
		Password:   "SecurePass123!@#",
		BackupCard: "4111111111111111",
	}

	if err := v.Validate(validServer); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Valid server")
	}

	// Validate an invalid server
	invalidServer := &Server{
		IPAddress:  "invalid-ip",
		Domain:     "invalid..domain",
		Password:   "weak",
		BackupCard: "1234-5678-9012-3456",
	}

	if err := v.Validate(invalidServer); err != nil {
		fmt.Println("Invalid server with errors:")
		for _, e := range err.(ValidationErrors) {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	}

	// Output:
	// Valid server
	// Invalid server with errors:
	// - IPAddress: invalid IP address format
	// - Domain: invalid domain name format
	// - Password: password must be at least 12 characters
	// - BackupCard: invalid credit card number format
}

func Example_networkAndFormatValidation() {
	type NetworkConfig struct {
		Network     string `validate:"cidr"`
		DeviceMAC   string `validate:"mac"`
		Location    string `validate:"latlong"`
		ThemeColor  string `validate:"color"`
	}

	v := New()
	v.AddRule("cidr", rules.CIDR{})
	v.AddRule("mac", rules.MAC{})
	v.AddRule("latlong", rules.LatLong{})
	v.AddRule("color", rules.Color{AllowHEX: true, AllowRGB: true})

	// Valid configuration
	validConfig := &NetworkConfig{
		Network:     "192.168.1.0/24",
		DeviceMAC:   "00:11:22:33:44:55",
		Location:    "40.7128,-74.0060",
		ThemeColor:  "#ff0000",
	}

	if err := v.Validate(validConfig); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}
	fmt.Println("Valid configuration")

	// Invalid configuration
	invalidConfig := &NetworkConfig{
		Network:     "192.168.1.0", // Missing subnet mask
		DeviceMAC:   "00:11:22:33", // Invalid MAC address
		Location:    "91.0000,0.0000", // Invalid latitude
		ThemeColor:  "red", // Invalid color format
	}

	if err := v.Validate(invalidConfig); err != nil {
		fmt.Printf("Validation errors:\n%v\n", err)
		return
	}

	// Output:
	// Valid configuration
}

func Example_internetValidation() {
	type ServerConfig struct {
		Hostname    string `validate:"hostname{AllowWildcard: true}"`
		Port        int    `validate:"port{AllowPrivileged: true}"`
		AdminEmail  string `validate:"emaildns{CheckDNS: true}"`
		Version     string `validate:"semver{AllowPrefix: true, AllowPrerelease: true}"`
		SubnetCIDR  string `validate:"cidr"`
		MACAddress  string `validate:"mac"`
	}

	// Register validation rules
	v := New()
	v.AddRule("Hostname", rules.Hostname{AllowWildcard: true})
	v.AddRule("Port", rules.Port{AllowPrivileged: true})
	v.AddRule("AdminEmail", rules.EmailDNS{CheckDNS: true})
	v.AddRule("Version", rules.SemVer{AllowPrefix: true, AllowPrerelease: true})
	v.AddRule("SubnetCIDR", rules.CIDR{})
	v.AddRule("MACAddress", rules.MAC{})

	config := ServerConfig{
		Hostname:    "*.example.com",
		Port:        443,
		AdminEmail:  "admin@example.com",
		Version:     "v1.2.3-beta.1",
		SubnetCIDR:  "192.168.1.0/24",
		MACAddress:  "00:1A:2B:3C:4D:5E",
	}

	err := v.Validate(config)
	fmt.Printf("Validation error: %v\n", err)

	// Invalid config
	invalidConfig := ServerConfig{
		Hostname:    "-invalid.com",
		Port:        70000,
		AdminEmail:  "not-an-email",
		Version:     "1.2",
		SubnetCIDR:  "300.168.1.0/24",
		MACAddress:  "invalid-mac",
	}

	err = v.Validate(invalidConfig)
	fmt.Printf("Validation errors:\n%v\n", err)

	// Output:
	// Validation error: AdminEmail: domain does not have valid MX records
	// Validation errors:
	// Hostname: invalid hostname format; Port: port must be between 1 and 65535; AdminEmail: invalid email format; Version: version must be in format X.Y.Z; SubnetCIDR: invalid CIDR format; MACAddress: invalid MAC address format
}
