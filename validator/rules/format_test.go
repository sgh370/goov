package rules

import "testing"

func TestLatLong(t *testing.T) {
	tests := []struct {
		name    string
		rule    LatLong
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid coordinates",
			rule:    LatLong{},
			value:   "40.7128,-74.0060",
			wantErr: false,
		},
		{
			name:    "valid coordinates with spaces",
			rule:    LatLong{},
			value:   "40.7128, -74.0060",
			wantErr: false,
		},
		{
			name:    "invalid format",
			rule:    LatLong{},
			value:   "40.7128",
			wantErr: true,
		},
		{
			name:    "invalid latitude - too high",
			rule:    LatLong{},
			value:   "91.0000,0.0000",
			wantErr: true,
		},
		{
			name:    "invalid latitude - too low",
			rule:    LatLong{},
			value:   "-91.0000,0.0000",
			wantErr: true,
		},
		{
			name:    "invalid longitude - too high",
			rule:    LatLong{},
			value:   "0.0000,181.0000",
			wantErr: true,
		},
		{
			name:    "invalid longitude - too low",
			rule:    LatLong{},
			value:   "0.0000,-181.0000",
			wantErr: true,
		},
		{
			name:    "invalid format - non-numeric",
			rule:    LatLong{},
			value:   "abc,def",
			wantErr: true,
		},
		{
			name:    "empty allowed",
			rule:    LatLong{AllowEmpty: true},
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty not allowed",
			rule:    LatLong{AllowEmpty: false},
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid type",
			rule:    LatLong{},
			value:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("LatLong.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColor(t *testing.T) {
	tests := []struct {
		name    string
		rule    Color
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid HEX - 6 digits",
			rule: Color{AllowHEX: true},
			value: "#ff0000",
			wantErr: false,
		},
		{
			name: "valid HEX - 3 digits",
			rule: Color{AllowHEX: true},
			value: "#f00",
			wantErr: false,
		},
		{
			name: "invalid HEX - wrong format",
			rule: Color{AllowHEX: true},
			value: "#ff00",
			wantErr: true,
		},
		{
			name: "invalid HEX - no hash",
			rule: Color{AllowHEX: true},
			value: "ff0000",
			wantErr: true,
		},
		{
			name: "valid RGB",
			rule: Color{AllowRGB: true},
			value: "rgb(255,0,0)",
			wantErr: false,
		},
		{
			name: "valid RGB with spaces",
			rule: Color{AllowRGB: true},
			value: "rgb(255, 0, 0)",
			wantErr: false,
		},
		{
			name: "invalid RGB - values too high",
			rule: Color{AllowRGB: true},
			value: "rgb(256,0,0)",
			wantErr: true,
		},
		{
			name: "invalid RGB - values too low",
			rule: Color{AllowRGB: true},
			value: "rgb(-1,0,0)",
			wantErr: true,
		},
		{
			name: "valid HSL",
			rule: Color{AllowHSL: true},
			value: "hsl(0,100%,50%)",
			wantErr: false,
		},
		{
			name: "valid HSL with spaces",
			rule: Color{AllowHSL: true},
			value: "hsl(0, 100%, 50%)",
			wantErr: false,
		},
		{
			name: "invalid HSL - hue too high",
			rule: Color{AllowHSL: true},
			value: "hsl(361,100%,50%)",
			wantErr: true,
		},
		{
			name: "invalid HSL - saturation too high",
			rule: Color{AllowHSL: true},
			value: "hsl(0,101%,50%)",
			wantErr: true,
		},
		{
			name: "invalid HSL - lightness too high",
			rule: Color{AllowHSL: true},
			value: "hsl(0,100%,101%)",
			wantErr: true,
		},
		{
			name: "empty allowed",
			rule: Color{AllowHEX: true, AllowEmpty: true},
			value: "",
			wantErr: false,
		},
		{
			name: "empty not allowed",
			rule: Color{AllowHEX: true, AllowEmpty: false},
			value: "",
			wantErr: true,
		},
		{
			name: "invalid type",
			rule: Color{AllowHEX: true},
			value: 123,
			wantErr: true,
		},
		{
			name: "no formats allowed",
			rule: Color{},
			value: "#ff0000",
			wantErr: true,
		},
		{
			name: "multiple formats allowed - HEX",
			rule: Color{AllowHEX: true, AllowRGB: true, AllowHSL: true},
			value: "#ff0000",
			wantErr: false,
		},
		{
			name: "multiple formats allowed - RGB",
			rule: Color{AllowHEX: true, AllowRGB: true, AllowHSL: true},
			value: "rgb(255,0,0)",
			wantErr: false,
		},
		{
			name: "multiple formats allowed - HSL",
			rule: Color{AllowHEX: true, AllowRGB: true, AllowHSL: true},
			value: "hsl(0,100%,50%)",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Color.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
