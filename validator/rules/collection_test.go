package rules

import (
	"testing"
)

type TestMapStruct struct {
	Data map[string]int `validate:"map"`
}

func TestMap(t *testing.T) {
	tests := []struct {
		name    string
		rule    Map
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid map",
			rule: Map{
				Key:   Required{},
				Value: Required{},
			},
			value: map[string]int{
				"key1": 1,
				"key2": 2,
			},
			wantErr: false,
		},
		{
			name: "nil map",
			rule: Map{
				Key:   Required{},
				Value: Required{},
			},
			value:   nil,
			wantErr: true,
		},
		{
			name: "empty map",
			rule: Map{
				Key:   Required{},
				Value: Required{},
			},
			value:   map[string]int{},
			wantErr: false,
		},
		{
			name: "invalid key",
			rule: Map{
				Key:   Required{},
				Value: Required{},
			},
			value: map[string]int{
				"": 1,
			},
			wantErr: true,
		},
		{
			name: "invalid value",
			rule: Map{
				Key:   Required{},
				Value: Min{Value: 10},
			},
			value: map[string]int{
				"key1": 5,
			},
			wantErr: true,
		},
		{
			name: "non-map value",
			rule: Map{
				Key:   Required{},
				Value: Required{},
			},
			value:   "not a map",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name    string
		rule    Slice
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid slice",
			rule: Slice{
				Rule: Required{},
			},
			value:   []int{1, 2, 3},
			wantErr: false,
		},
		{
			name: "nil slice",
			rule: Slice{
				Rule: Required{},
			},
			value:   nil,
			wantErr: true,
		},
		{
			name: "empty slice",
			rule: Slice{
				Rule: Required{},
			},
			value:   []int{},
			wantErr: false,
		},
		{
			name: "invalid item",
			rule: Slice{
				Rule: Min{Value: 10},
			},
			value:   []int{5, 15, 8},
			wantErr: true,
		},
		{
			name: "non-slice value",
			rule: Slice{
				Rule: Required{},
			},
			value:   "not a slice",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slice.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEachMulti_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    EachMulti
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid items",
			rule: EachMulti{
				Rules: []Rule{Required{}, Length{Min: 1}},
			},
			value:   []string{"valid", "test"},
			wantErr: false,
		},
		{
			name: "invalid items",
			rule: EachMulti{
				Rules: []Rule{Required{}, Length{Min: 5}},
			},
			value:   []string{"", "test"},
			wantErr: true,
		},
		{
			name: "non-slice value",
			rule: EachMulti{
				Rules: []Rule{Required{}},
			},
			value:   "not a slice",
			wantErr: true,
		},
		{
			name: "nil value",
			rule: EachMulti{
				Rules: []Rule{Required{}},
			},
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("EachMulti.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeys_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    Keys
		value   interface{}
		wantErr bool
	}{
		{
			name: "valid keys",
			rule: Keys{
				Rules: []Rule{Required{}, Length{Min: 1}},
			},
			value: map[string]string{
				"valid": "value",
				"test":  "value",
			},
			wantErr: false,
		},
		{
			name: "invalid keys",
			rule: Keys{
				Rules: []Rule{Required{}, Length{Min: 5}},
			},
			value: map[string]string{
				"":     "value",
				"test": "value",
			},
			wantErr: true,
		},
		{
			name: "non-map value",
			rule: Keys{
				Rules: []Rule{Required{}},
			},
			value:   "not a map",
			wantErr: true,
		},
		{
			name: "nil value",
			rule: Keys{
				Rules: []Rule{Required{}},
			},
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keys.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContains_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    Contains
		value   interface{}
		wantErr bool
	}{
		{
			name:    "slice contains value",
			rule:    Contains{Value: 2},
			value:   []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "slice does not contain value",
			rule:    Contains{Value: 4},
			value:   []int{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "string slice contains value",
			rule:    Contains{Value: "world"},
			value:   []string{"hello", "world", "!"},
			wantErr: false,
		},
		{
			name:    "string slice does not contain value",
			rule:    Contains{Value: "xyz"},
			value:   []string{"hello", "world", "!"},
			wantErr: true,
		},
		{
			name:    "non-slice value",
			rule:    Contains{Value: 1},
			value:   123,
			wantErr: true,
		},
		{
			name:    "nil value",
			rule:    Contains{Value: 1},
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnique_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    Unique
		value   interface{}
		wantErr bool
	}{
		{
			name:    "unique slice",
			rule:    Unique{},
			value:   []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "non-unique slice",
			rule:    Unique{},
			value:   []int{1, 2, 2, 3},
			wantErr: true,
		},
		{
			name:    "unique string slice",
			rule:    Unique{},
			value:   []string{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "non-unique string slice",
			rule:    Unique{},
			value:   []string{"a", "b", "b", "c"},
			wantErr: true,
		},
		{
			name:    "non-slice value",
			rule:    Unique{},
			value:   "not a slice",
			wantErr: true,
		},
		{
			name:    "nil value",
			rule:    Unique{},
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unique.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
