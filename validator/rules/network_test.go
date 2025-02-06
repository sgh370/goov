package rules

import "testing"

func TestCIDR(t *testing.T) {
	tests := []struct {
		name    string
		rule    CIDR
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid CIDR",
			rule:    CIDR{},
			value:   "192.168.1.0/24",
			wantErr: false,
		},
		{
			name:    "valid CIDR with small subnet",
			rule:    CIDR{},
			value:   "10.0.0.0/8",
			wantErr: false,
		},
		{
			name:    "invalid CIDR - wrong format",
			rule:    CIDR{},
			value:   "192.168.1.0",
			wantErr: true,
		},
		{
			name:    "invalid CIDR - invalid IP",
			rule:    CIDR{},
			value:   "256.256.256.0/24",
			wantErr: true,
		},
		{
			name:    "invalid CIDR - invalid subnet",
			rule:    CIDR{},
			value:   "192.168.1.0/33",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    CIDR{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    CIDR{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    CIDR{},
			value:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("CIDR.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMAC(t *testing.T) {
	tests := []struct {
		name    string
		rule    MAC
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid MAC with colons",
			rule:    MAC{},
			value:   "00:11:22:33:44:55",
			wantErr: false,
		},
		{
			name:    "valid MAC with hyphens",
			rule:    MAC{},
			value:   "00-11-22-33-44-55",
			wantErr: false,
		},
		{
			name:    "valid MAC without separators",
			rule:    MAC{},
			value:   "001122334455",
			wantErr: false,
		},
		{
			name:    "invalid MAC - wrong length",
			rule:    MAC{},
			value:   "00:11:22:33:44",
			wantErr: true,
		},
		{
			name:    "invalid MAC - invalid characters",
			rule:    MAC{},
			value:   "00:11:22:33:44:GG",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    MAC{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    MAC{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    MAC{},
			value:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("MAC.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
