package rules

import (
	"testing"
	"strings"
	"fmt"
)

func TestIP(t *testing.T) {
	tests := []struct {
		name    string
		rule    IP
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid IPv4",
			rule: IP{AllowV4: true, AllowV6: true},
			value: "192.168.1.1",
			wantErr: false,
		},
		{
			name: "valid IPv6",
			rule: IP{AllowV4: true, AllowV6: true},
			value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			wantErr: false,
		},
		{
			name: "invalid IP",
			rule: IP{AllowV4: true, AllowV6: true},
			value: "invalid-ip",
			wantErr: true,
		},
		{
			name: "IPv4 not allowed",
			rule: IP{AllowV4: false, AllowV6: true},
			value: "192.168.1.1",
			wantErr: true,
		},
		{
			name: "IPv6 not allowed",
			rule: IP{AllowV4: true, AllowV6: false},
			value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			wantErr: true,
		},
		{
			name: "empty allowed",
			rule: IP{AllowV4: true, AllowV6: true, AllowEmpty: true},
			value: "",
			wantErr: false,
		},
		{
			name: "empty not allowed",
			rule: IP{AllowV4: true, AllowV6: true, AllowEmpty: false},
			value: "",
			wantErr: true,
		},
		{
			name: "invalid type",
			rule: IP{AllowV4: true, AllowV6: true},
			value: 123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("IP.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	tests := []struct {
		name    string
		rule    Domain
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid domain",
			rule: Domain{AllowSubdomains: true},
			value: "example.com",
			wantErr: false,
		},
		{
			name: "valid subdomain",
			rule: Domain{AllowSubdomains: true},
			value: "sub.example.com",
			wantErr: false,
		},
		{
			name: "subdomain not allowed",
			rule: Domain{AllowSubdomains: false},
			value: "sub.example.com",
			wantErr: true,
		},
		{
			name: "invalid domain - no TLD",
			rule: Domain{AllowSubdomains: true},
			value: "example",
			wantErr: true,
		},
		{
			name: "invalid domain - numeric TLD",
			rule: Domain{AllowSubdomains: true},
			value: "example.123",
			wantErr: true,
		},
		{
			name: "invalid domain - too long label",
			rule: Domain{AllowSubdomains: true},
			value: strings.Repeat("a", 64) + ".com",
			wantErr: true,
		},
		{
			name: "empty allowed",
			rule: Domain{AllowSubdomains: true, AllowEmpty: true},
			value: "",
			wantErr: false,
		},
		{
			name: "empty not allowed",
			rule: Domain{AllowSubdomains: true, AllowEmpty: false},
			value: "",
			wantErr: true,
		},
		{
			name: "invalid type",
			rule: Domain{AllowSubdomains: true},
			value: 123,
			wantErr: true,
		},
		{
			name: "invalid length",
			rule: Domain{AllowSubdomains: true},
			value: strings.Repeat("a", 265) + ".com",
			wantErr: true,
		},
		{
			name: "invalid format",
			rule: Domain{AllowSubdomains: true},
			value: ".",
			wantErr: true,
		},
		{
			name: "invalid format2",
			rule: Domain{AllowSubdomains: true},
			value: "#.s",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Domain.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		name    string
		rule    Password
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid password - all requirements",
			rule: Password{
				MinLength: 8,
				MaxLength: 20,
				RequireUpper: true,
				RequireLower: true,
				RequireDigit: true,
				RequireSpecial: true,
			},
			value: "Test123!@",
			wantErr: false,
		},
		{
			name: "too short",
			rule: Password{MinLength: 8},
			value: "Test123",
			wantErr: true,
		},
		{
			name: "too long",
			rule: Password{MinLength: 8, MaxLength: 20},
			value: "Test123!@#$%^&*()+_)(*&^%$#@!",
			wantErr: true,
		},
		{
			name: "missing uppercase",
			rule: Password{MinLength: 8, RequireUpper: true},
			value: "test123!@",
			wantErr: true,
		},
		{
			name: "missing lowercase",
			rule: Password{MinLength: 8, RequireLower: true},
			value: "TEST123!@",
			wantErr: true,
		},
		{
			name: "missing digit",
			rule: Password{MinLength: 8, RequireDigit: true},
			value: "TestTest!@",
			wantErr: true,
		},
		{
			name: "missing special",
			rule: Password{MinLength: 8, RequireSpecial: true},
			value: "Test1234",
			wantErr: true,
		},
		{
			name: "invalid type",
			rule: Password{MinLength: 8},
			value: 123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Password.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreditCard(t *testing.T) {
	tests := []struct {
		name    string
		rule    CreditCard
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid Visa",
			rule: CreditCard{},
			value: "4111111111111111",
			wantErr: false,
		},
		{
			name: "valid MasterCard",
			rule: CreditCard{},
			value: "5555555555554444",
			wantErr: false,
		},
		{
			name: "valid American Express",
			rule: CreditCard{},
			value: "378282246310005",
			wantErr: false,
		},
		{
			name: "invalid - wrong format",
			rule: CreditCard{},
			value: "1234",
			wantErr: true,
		},
		{
			name: "invalid - fails Luhn",
			rule: CreditCard{},
			value: "4532815137901852",
			wantErr: true,
		},
		{
			name: "empty allowed",
			rule: CreditCard{AllowEmpty: true},
			value: "",
			wantErr: false,
		},
		{
			name: "empty not allowed",
			rule: CreditCard{AllowEmpty: false},
			value: "",
			wantErr: true,
		},
		{
			name: "invalid type",
			rule: CreditCard{},
			value: 123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("CreditCard.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCrossField_SetParent(t *testing.T) {
	parent := struct{}{}
	cf := CrossField{ValidateFn: func(_, _ interface{}) error { return nil }}
	cf.SetParent(parent)
	if cf.parent != parent {
		t.Error("Expected parent to be set")
	}
}

func TestCrossField_Validate_NoValidationFn(t *testing.T) {
	cf := CrossField{}
	err := cf.Validate(nil)
	if err == nil || err.Error() != "validation function not provided" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCrossField_Validate_NoParent(t *testing.T) {
	cf := CrossField{ValidateFn: func(_, _ interface{}) error { return nil }}
	err := cf.Validate(nil)
	if err == nil || err.Error() != "parent not set" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCrossField_Validate_Success(t *testing.T) {
	parent := struct{ Value int }{42}
	cf := CrossField{
		ValidateFn: func(p, v interface{}) error {
			if p.(struct{ Value int }).Value != 42 {
				return fmt.Errorf("invalid parent value")
			}
			return nil
		},
	}
	cf.SetParent(parent)
	if err := cf.Validate(nil); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
