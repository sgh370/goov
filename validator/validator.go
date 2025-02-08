package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sgh370/goov/validator/rules"
)

type Validator struct {
	rules map[string]rules.Rule
}

func New() *Validator {
	return &Validator{
		rules: make(map[string]rules.Rule),
	}
}

func (v *Validator) AddRule(name string, rule rules.Rule) {
	v.rules[name] = rule
}

func (v *Validator) Validate(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("value is nil")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("value must be a struct or pointer to struct")
	}

	return v.validateStruct(val)
}

func (v *Validator) validateStruct(val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	valType := val.Type()
	var parent interface{}
	if val.CanAddr() {
		parent = val.Addr().Interface()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := valType.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Handle nested struct validation
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				if err := v.validateField(field, tag, parent); err != nil {
					return fmt.Errorf("%s: %v", fieldType.Name, err)
				}
				continue
			}
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := v.validateStruct(field); err != nil {
				return fmt.Errorf("%s: %v", fieldType.Name, err)
			}
		}

		if err := v.validateField(field, tag, parent); err != nil {
			return fmt.Errorf("%s: %v", fieldType.Name, err)
		}
	}

	return nil
}

func (v *Validator) validateField(field reflect.Value, tag string, parent interface{}) error {
	if tag == "" {
		return nil
	}

	for _, rule := range strings.Split(tag, ",") {
		parts := strings.Split(rule, "=")
		ruleName := parts[0]
		var ruleValue string
		if len(parts) > 1 {
			ruleValue = parts[1]
		}

		switch ruleName {
		case "slice":
			if err := v.validateSlice(field, rule); err != nil {
				return err
			}
		case "min":
			val, err := strconv.ParseFloat(ruleValue, 64)
			if err != nil {
				return fmt.Errorf("invalid min value: %s", ruleValue)
			}
			minRule := rules.Min{Value: val}
			if err := minRule.Validate(field.Interface()); err != nil {
				return err
			}
		default:
			rule := v.rules[ruleName]
			if rule == nil {
				return fmt.Errorf("unknown validation rule: %s", ruleName)
			}
			if setter, ok := rule.(interface{ SetParent(interface{}) }); ok {
				setter.SetParent(parent)
			}
			if err := rule.Validate(field.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validator) validateSlice(field reflect.Value, tag string) error {
	if tag == "" {
		return nil
	}

	parts := strings.Split(tag, "=")
	if len(parts) != 2 {
		return fmt.Errorf("invalid slice validation format: %s", tag)
	}

	rule := v.rules[parts[1]]
	if rule == nil {
		return fmt.Errorf("unknown validation rule: %s", parts[1])
	}

	if field.Kind() != reflect.Slice {
		return fmt.Errorf("field is not a slice")
	}

	if field.IsNil() {
		return fmt.Errorf("slice is nil")
	}

	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
		if item.Kind() == reflect.Ptr && !item.IsNil() {
			item = item.Elem()
		}

		if item.Kind() == reflect.Struct {
			if err := v.validateStruct(item); err != nil {
				return fmt.Errorf("item at index %d: %v", i, err)
			}
		} else {
			if err := rule.Validate(item.Interface()); err != nil {
				return fmt.Errorf("item at index %d: %v", i, err)
			}
		}
	}

	return nil
}

func (v *Validator) ValidateAll(value interface{}) []error {
	var errors []error

	if value == nil {
		return append(errors, fmt.Errorf("value is nil"))
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return append(errors, fmt.Errorf("value is nil"))
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return append(errors, fmt.Errorf("value must be a struct or pointer to struct"))
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				for _, ruleName := range strings.Split(tag, ",") {
					rule, ok := v.rules[ruleName]
					if !ok {
						errors = append(errors, fmt.Errorf("%s: unknown validation rule: %s", fieldType.Name, ruleName))
						continue
					}

					if setter, ok := rule.(interface{ SetParent(interface{}) }); ok {
						setter.SetParent(val.Addr().Interface())
					}

					if err := rule.Validate(nil); err != nil {
						errors = append(errors, fmt.Errorf("%s: %v", fieldType.Name, err))
					}
				}
				continue
			}
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := v.validateStruct(field); err != nil {
				errors = append(errors, fmt.Errorf("%s: %v", fieldType.Name, err))
			}
			continue
		}

		if field.Kind() == reflect.Slice {
			if err := v.validateSlice(field, tag); err != nil {
				errors = append(errors, fmt.Errorf("%s: %v", fieldType.Name, err))
			}
			continue
		}

		if err := v.validateField(field, tag, val.Addr().Interface()); err != nil {
			errors = append(errors, fmt.Errorf("%s: %v", fieldType.Name, err))
		}
	}

	return errors
}

type Pattern struct {
	rules []rules.Rule
}

func NewPattern(rules ...rules.Rule) *Pattern {
	return &Pattern{rules: rules}
}

func (p *Pattern) Validate(value interface{}) error {
	for _, rule := range p.rules {
		if err := rule.Validate(value); err != nil {
			return err
		}
	}
	return nil
}
