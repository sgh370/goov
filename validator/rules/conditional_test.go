package rules

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Age      int
	Premium  bool
	Username string
	Email    string
}

type TestParent struct {
	Condition bool
	Field     string
}

type TestStruct struct {
	Field     bool
	OtherBool bool
	Value     string
}

func TestWhen(t *testing.T) {
	tests := []struct {
		name    string
		rule    When
		value   interface{}
		field   string
		parent  interface{}
		wantErr bool
	}{
		{
			name: "condition true, then rule passes",
			rule: When{
				Condition: func(v interface{}) bool {
					s := v.(*testStruct)
					return s.Age > 18
				},
				Then: Required{},
			},
			value:   &testStruct{Age: 20, Username: "test"},
			field:   "Username",
			parent:  &testStruct{Age: 20, Username: "test"},
			wantErr: false,
		},
		{
			name: "condition true, then rule fails",
			rule: When{
				Condition: func(v interface{}) bool {
					s := v.(*testStruct)
					return s.Age > 18
				},
				Then: Required{},
			},
			value:   &testStruct{Age: 20, Username: ""},
			field:   "Username",
			parent:  &testStruct{Age: 20, Username: ""},
			wantErr: true,
		},
		{
			name: "condition false, else rule passes",
			rule: When{
				Condition: func(v interface{}) bool {
					s := v.(*testStruct)
					return s.Age > 18
				},
				Then: Required{},
				Else: Length{Max: 10},
			},
			value:   &testStruct{Age: 16, Username: "test"},
			field:   "Username",
			parent:  &testStruct{Age: 16, Username: "test"},
			wantErr: false,
		},
		{
			name: "condition false, else rule fails",
			rule: When{
				Condition: func(v interface{}) bool {
					s := v.(*testStruct)
					return s.Age > 18
				},
				Then: Required{},
				Else: Length{Max: 3},
			},
			value:   &testStruct{Age: 16, Username: "test"},
			field:   "Username",
			parent:  &testStruct{Age: 16, Username: "test"},
			wantErr: true,
		},
		{
			name: "nil parent",
			rule: When{
				Condition: func(v interface{}) bool {
					return true
				},
				Then: Required{},
			},
			value:   "",
			field:   "",
			parent:  nil,
			wantErr: true,
		},
		{
			name: "condition true, then rule passes",
			rule: When{
				Condition: func(parent interface{}) bool {
					return parent.(*TestParent).Condition
				},
				Then: Required{},
			},
			value: "value",
			parent: &TestParent{
				Condition: true,
				Field:     "value",
			},
			wantErr: false,
		},
		{
			name: "condition true, then rule fails",
			rule: When{
				Condition: func(parent interface{}) bool {
					return parent.(*TestParent).Condition
				},
				Then: Required{},
			},
			value: "",
			parent: &TestParent{
				Condition: true,
				Field:     "",
			},
			wantErr: true,
		},
		{
			name: "condition false, else rule passes",
			rule: When{
				Condition: func(parent interface{}) bool {
					return parent.(*TestParent).Condition
				},
				Then: Required{},
				Else: Length{Min: 1, Max: 10},
			},
			value: "value",
			parent: &TestParent{
				Condition: false,
				Field:     "value",
			},
			wantErr: false,
		},
		{
			name: "condition false, else rule fails",
			rule: When{
				Condition: func(parent interface{}) bool {
					return parent.(*TestParent).Condition
				},
				Then: Required{},
				Else: Length{Min: 1, Max: 3},
			},
			value: "too long",
			parent: &TestParent{
				Condition: false,
				Field:     "too long",
			},
			wantErr: true,
		},
		{
			name: "nil parent",
			rule: When{
				Condition: func(parent interface{}) bool {
					return parent != nil
				},
				Then: Required{},
			},
			value:   "",
			parent:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field != "" {
				tt.rule.SetParent(tt.parent)
				v := reflect.ValueOf(tt.parent)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				err := tt.rule.Validate(v.FieldByName(tt.field).Interface())
				if (err != nil) != tt.wantErr {
					t.Errorf("When.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				tt.rule.SetParent(tt.parent)
				err := tt.rule.Validate(tt.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("When.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestIf(t *testing.T) {
	tests := []struct {
		name    string
		rule    If
		value   interface{}
		parent  interface{}
		wantErr bool
	}{
		{
			name: "field is true, rule passes",
			rule: If{
				Field: "Field",
				Then:  Required{},
			},
			value:   "value",
			parent:  &TestStruct{Field: true, Value: "value"},
			wantErr: false,
		},
		{
			name: "field is true, rule fails",
			rule: If{
				Field: "Field",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{Field: true},
			wantErr: true,
		},
		{
			name: "field is false, rule skipped",
			rule: If{
				Field: "Field",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{Field: false},
			wantErr: false,
		},
		{
			name: "invalid field name",
			rule: If{
				Field: "NonExistent",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{},
			wantErr: true,
		},
		{
			name: "non-bool field",
			rule: If{
				Field: "Value",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{Value: "not a bool"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rule.SetParent(tt.parent)
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("If.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnless(t *testing.T) {
	tests := []struct {
		name    string
		rule    Unless
		value   interface{}
		parent  interface{}
		wantErr bool
	}{
		{
			name: "field is false, rule passes",
			rule: Unless{
				Field: "Field",
				Then:  Required{},
			},
			value:   "value",
			parent:  &TestStruct{Field: false, Value: "value"},
			wantErr: false,
		},
		{
			name: "field is false, rule fails",
			rule: Unless{
				Field: "Field",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{Field: false},
			wantErr: true,
		},
		{
			name: "field is true, rule skipped",
			rule: Unless{
				Field: "Field",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{Field: true},
			wantErr: false,
		},
		{
			name: "invalid field name",
			rule: Unless{
				Field: "NonExistent",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{},
			wantErr: true,
		},
		{
			name: "non-bool field",
			rule: Unless{
				Field: "Value",
				Then:  Required{},
			},
			value:   "",
			parent:  &TestStruct{Value: "not a bool"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rule.SetParent(tt.parent)
			err := tt.rule.Validate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unless.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWhen_SetParent(t *testing.T) {
	type TestStruct struct {
		Field    bool
		OtherVal string
	}

	parent := &TestStruct{Field: true, OtherVal: "test"}
	rule := When{
		Condition: func(v interface{}) bool {
			return v.(*TestStruct).Field
		},
		Then: Required{},
	}

	rule.SetParent(parent)
	err := rule.Validate("test")
	if err != nil {
		t.Errorf("When.Validate() error = %v, wantErr false", err)
	}

	// Test with false condition
	parent.Field = false
	err = rule.Validate("")
	if err != nil {
		t.Errorf("When.Validate() with false condition error = %v, wantErr false", err)
	}
}

func TestIf_SetParent(t *testing.T) {
	type TestStruct struct {
		Field    bool
		OtherVal string
	}

	parent := &TestStruct{Field: true, OtherVal: "test"}
	rule := If{
		Field: "Field",
		Then:  Required{},
	}

	rule.SetParent(parent)
	err := rule.Validate("test")
	if err != nil {
		t.Errorf("If.Validate() error = %v, wantErr false", err)
	}

	// Test with false condition
	parent.Field = false
	err = rule.Validate("")
	if err != nil {
		t.Errorf("If.Validate() with false condition error = %v, wantErr false", err)
	}
}

func TestUnless_SetParent(t *testing.T) {
	type TestStruct struct {
		Field    bool
		OtherVal string
	}

	parent := &TestStruct{Field: false, OtherVal: "test"}
	rule := Unless{
		Field: "Field",
		Then:  Required{},
	}

	rule.SetParent(parent)
	err := rule.Validate("test")
	if err != nil {
		t.Errorf("Unless.Validate() error = %v, wantErr false", err)
	}

	// Test with true condition
	parent.Field = true
	err = rule.Validate("")
	if err != nil {
		t.Errorf("Unless.Validate() with true condition error = %v, wantErr false", err)
	}
}
