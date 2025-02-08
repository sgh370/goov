package rules

import "testing"

func TestEmailDNS(t *testing.T) {
	tests := []struct {
		name    string
		rule    EmailDNS
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid email",
			rule:    EmailDNS{},
			value:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with DNS check disabled",
			rule:    EmailDNS{CheckDNS: false},
			value:   "test@gmail.com",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			rule:    EmailDNS{},
			value:   "test.example.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			rule:    EmailDNS{},
			value:   "test@",
			wantErr: true,
		},
		{
			name:    "invalid email - invalid chars",
			rule:    EmailDNS{},
			value:   "test!@example.com",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    EmailDNS{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    EmailDNS{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    EmailDNS{},
			value:   123,
			wantErr: true,
		},
		{
			name:    "invalid email - invalid lookup",
			rule:    EmailDNS{CheckDNS: true},
			value:   "test@asdasd.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("EmailDNS.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHostname(t *testing.T) {
	tests := []struct {
		name    string
		rule    Hostname
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid hostname",
			rule:    Hostname{},
			value:   "example.com",
			wantErr: false,
		},
		{
			name:    "valid hostname with subdomain",
			rule:    Hostname{},
			value:   "sub.example.com",
			wantErr: false,
		},
		{
			name:    "valid hostname with wildcard",
			rule:    Hostname{AllowWildcard: true},
			value:   "*.example.com",
			wantErr: false,
		},
		{
			name:    "invalid hostname - too long",
			rule:    Hostname{},
			value:   string(make([]byte, 256)),
			wantErr: true,
		},
		{
			name:    "invalid hostname - invalid chars",
			rule:    Hostname{},
			value:   "test!.example.com",
			wantErr: true,
		},
		{
			name:    "invalid hostname - starts with hyphen",
			rule:    Hostname{},
			value:   "-example.com",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    Hostname{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    Hostname{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    Hostname{},
			value:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Hostname.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPort(t *testing.T) {
	tests := []struct {
		name    string
		rule    Port
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid port",
			rule:    Port{},
			value:   8080,
			wantErr: false,
		},
		{
			name:    "valid port string",
			rule:    Port{},
			value:   "8080",
			wantErr: false,
		},
		{
			name:    "valid privileged port",
			rule:    Port{AllowPrivileged: true},
			value:   80,
			wantErr: false,
		},
		{
			name:    "invalid privileged port",
			rule:    Port{},
			value:   80,
			wantErr: true,
		},
		{
			name:    "invalid port - too high",
			rule:    Port{},
			value:   65536,
			wantErr: true,
		},
		{
			name:    "invalid port - too low",
			rule:    Port{},
			value:   0,
			wantErr: true,
		},
		{
			name:    "invalid port string",
			rule:    Port{},
			value:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    Port{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    Port{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    Port{},
			value:   []int{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Port.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSemVer(t *testing.T) {
	tests := []struct {
		name    string
		rule    SemVer
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid version",
			rule:    SemVer{},
			value:   "1.2.3",
			wantErr: false,
		},
		{
			name:    "valid version with v prefix",
			rule:    SemVer{AllowPrefix: true},
			value:   "v1.2.3",
			wantErr: false,
		},
		{
			name:    "valid version with prerelease",
			rule:    SemVer{AllowPrerelease: true},
			value:   "1.2.3-alpha.1",
			wantErr: false,
		},
		{
			name:    "valid version with build",
			rule:    SemVer{AllowBuild: true},
			value:   "1.2.3+001",
			wantErr: false,
		},
		{
			name:    "valid version with prerelease and build",
			rule:    SemVer{AllowPrerelease: true, AllowBuild: true},
			value:   "1.2.3-beta.2+build.123",
			wantErr: false,
		},
		{
			name:    "invalid version - wrong format",
			rule:    SemVer{},
			value:   "1.2",
			wantErr: true,
		},
		{
			name:    "invalid version - non-numeric",
			rule:    SemVer{},
			value:   "1.2.x",
			wantErr: true,
		},
		{
			name:    "invalid version - prerelease not allowed",
			rule:    SemVer{},
			value:   "1.2.3-alpha",
			wantErr: true,
		},
		{
			name:    "invalid version - build not allowed",
			rule:    SemVer{},
			value:   "1.2.3+001",
			wantErr: true,
		},
		{
			name:    "invalid version - v prefix not allowed",
			rule:    SemVer{},
			value:   "v1.2.3",
			wantErr: true,
		},
		{
			name:    "invalid version - v prefix required",
			rule:    SemVer{AllowPrefix: true, RequirePrefix: true},
			value:   "1.2.3",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    SemVer{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    SemVer{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    SemVer{},
			value:   123,
			wantErr: true,
		},
		{
			name: "Invalid prerelease with trailing dot",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: true,
				AllowBuild:      false,
				AllowEmpty:      false,
			},
			value:   "1.0.0-alpha.",
			wantErr: true,
		},
		{
			name: "Valid version with build metadata",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      true, // Build metadata allowed
				AllowEmpty:      false,
			},
			value:   "1.0.0+build.123",
			wantErr: false,
		},
		{
			name: "Invalid build metadata with trailing dot",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      true, // Build metadata allowed
				AllowEmpty:      false,
			},
			value:   "1.0.0+build.",
			wantErr: true,
		},
		{
			name: "Invalid build metadata with leading dot",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      true, // Build metadata allowed
				AllowEmpty:      false,
			},
			value:   "1.0.0+.build",
			wantErr: true,
		},
		{
			name: "Invalid build metadata with special characters",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      true, // Build metadata allowed
				AllowEmpty:      false,
			},
			value:   "1.0.0+build@123",
			wantErr: true,
		},
		{
			name: "Empty build metadata part",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      true, // Build metadata allowed
				AllowEmpty:      false,
			},
			value:   "1.0.0+",
			wantErr: true,
		},
		{
			name: "Build metadata not allowed",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      false, // Build metadata not allowed
				AllowEmpty:      false,
			},
			value:   "1.0.0+build",
			wantErr: true,
		},
		{
			name: "Valid version without build metadata",
			rule: SemVer{
				AllowPrefix:     false,
				RequirePrefix:   false,
				AllowPrerelease: false,
				AllowBuild:      false,
				AllowEmpty:      false,
			},
			value:   "1.0.0",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("SemVer.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
