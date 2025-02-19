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

// mockRule implements Rule interface and SetParent interface for testing
type mockRule struct {
	parent interface{}
}

func (m *mockRule) SetParent(parent interface{}) {
	m.parent = parent
}

func (m *mockRule) Validate(value interface{}) error {
	return nil
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
			name: "nil parent",
			rule: If{
				Field: "Value",
				Then:  Required{},
			},
			value:   "",
			parent:  nil,
			wantErr: true,
		},
		{
			name: "non-struct parent",
			rule: If{
				Field: "NonExistent",
				Then:  Required{},
			},
			value:   "",
			parent:  map[string]string{},
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

	tests := []struct {
		name          string
		rule          When
		parent        interface{}
		thenSetsParent bool
		elseSetsParent bool
	}{
		{
			name: "Then implements SetParent",
			rule: When{
				Condition: func(v interface{}) bool { return true },
				Then:      &mockRule{},
				Else:      Required{},
			},
			parent:        &TestStruct{Field: true},
			thenSetsParent: true,
			elseSetsParent: false,
		},
		{
			name: "Else implements SetParent",
			rule: When{
				Condition: func(v interface{}) bool { return true },
				Then:      Required{},
				Else:      &mockRule{},
			},
			parent:        &TestStruct{Field: true},
			thenSetsParent: false,
			elseSetsParent: true,
		},
		{
			name: "Both implement SetParent",
			rule: When{
				Condition: func(v interface{}) bool { return true },
				Then:      &mockRule{},
				Else:      &mockRule{},
			},
			parent:        &TestStruct{Field: true},
			thenSetsParent: true,
			elseSetsParent: true,
		},
		{
			name: "Neither implements SetParent",
			rule: When{
				Condition: func(v interface{}) bool { return true },
				Then:      Required{},
				Else:      Required{},
			},
			parent:        &TestStruct{Field: true},
			thenSetsParent: false,
			elseSetsParent: false,
		},
		{
			name: "Then is nil",
			rule: When{
				Condition: func(v interface{}) bool { return true },
				Then:      nil,
				Else:      &mockRule{},
			},
			parent:        &TestStruct{Field: true},
			thenSetsParent: false,
			elseSetsParent: true,
		},
		{
			name: "Else is nil",
			rule: When{
				Condition: func(v interface{}) bool { return true },
				Then:      &mockRule{},
				Else:      nil,
			},
			parent:        &TestStruct{Field: true},
			thenSetsParent: true,
			elseSetsParent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rule.SetParent(tt.parent)

			// Check if Then rule's parent was set correctly
			if tt.thenSetsParent {
				thenRule, ok := tt.rule.Then.(*mockRule)
				if !ok {
					t.Error("Then rule should be *mockRule")
				} else if thenRule.parent != tt.parent {
					t.Error("Then rule's parent not set correctly")
				}
			}

			// Check if Else rule's parent was set correctly
			if tt.elseSetsParent {
				elseRule, ok := tt.rule.Else.(*mockRule)
				if !ok {
					t.Error("Else rule should be *mockRule")
				} else if elseRule.parent != tt.parent {
					t.Error("Else rule's parent not set correctly")
				}
			}

			// Verify the When rule's own parent is set
			if tt.rule.parent != tt.parent {
				t.Error("When rule's parent not set correctly")
			}
		})
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

func TestDependentRequired(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Email   string
		Address string
		Active  bool
		Count   int
	}

	tests := []struct {
		name    string
		rule    DependentRequired
		parent  interface{}
		field   string
		wantErr bool
	}{
		{
			name: "valid - field has value",
			rule: DependentRequired{
				Field: "Name",
			},
			parent: &Person{
				Name: "John",
			},
			wantErr: false,
		},
		{
			name: "invalid - field is empty string",
			rule: DependentRequired{
				Field: "Name",
			},
			parent: &Person{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "invalid - field is zero int",
			rule: DependentRequired{
				Field: "Age",
			},
			parent: &Person{
				Age: 0,
			},
			wantErr: true,
		},
		{
			name: "valid - field is non-zero int",
			rule: DependentRequired{
				Field: "Age",
			},
			parent: &Person{
				Age: 25,
			},
			wantErr: false,
		},
		{
			name: "invalid - field is false bool",
			rule: DependentRequired{
				Field: "Active",
			},
			parent: &Person{
				Active: false,
			},
			wantErr: true,
		},
		{
			name: "valid - field is true bool",
			rule: DependentRequired{
				Field: "Active",
			},
			parent: &Person{
				Active: true,
			},
			wantErr: false,
		},
		{
			name: "error - parent not provided",
			rule: DependentRequired{
				Field: "Name",
			},
			parent:  nil,
			wantErr: true,
		},
		{
			name: "error - parent is not a struct",
			rule: DependentRequired{
				Field: "Name",
			},
			parent:  "not a struct",
			wantErr: true,
		},
		{
			name: "error - field doesn't exist",
			rule: DependentRequired{
				Field: "NonExistentField",
			},
			parent: &Person{
				Name: "John",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rule.Parent = tt.parent
			err := tt.rule.Validate(nil) // value is not used in DependentRequired
			if (err != nil) != tt.wantErr {
				t.Errorf("DependentRequired.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
