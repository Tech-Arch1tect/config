package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidatorFunc func(field reflect.StructField, value reflect.Value, param string) error

var validators = map[string]ValidatorFunc{}

func RegisterValidator(name string, fn ValidatorFunc) {
	validators[name] = fn
}

func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			var ruleName, param string
			if strings.Contains(rule, "=") {
				parts := strings.SplitN(rule, "=", 2)
				ruleName = parts[0]
				param = parts[1]
			} else {
				ruleName = rule
			}
			validator, exists := validators[ruleName]
			if !exists {
				return fmt.Errorf("no validator registered for rule '%s'", ruleName)
			}
			if err := validator(field, fieldValue, param); err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	RegisterValidator("required", func(field reflect.StructField, value reflect.Value, param string) error {
		if reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface()) {
			return fmt.Errorf("field '%s' is required", field.Name)
		}
		return nil
	})

	RegisterValidator("min", func(field reflect.StructField, value reflect.Value, param string) error {
		minVal, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid min value for field '%s'", field.Name)
		}
		switch value.Kind() {
		case reflect.String:
			if len(value.String()) < int(minVal) {
				return fmt.Errorf("field '%s' must be at least %d characters", field.Name, minVal)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < minVal {
				return fmt.Errorf("field '%s' must be at least %d", field.Name, minVal)
			}
		default:
			return fmt.Errorf("unsupported type for min validation: %s", value.Kind())
		}
		return nil
	})

	RegisterValidator("max", func(field reflect.StructField, value reflect.Value, param string) error {
		maxVal, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max value for field '%s'", field.Name)
		}
		switch value.Kind() {
		case reflect.String:
			if len(value.String()) > int(maxVal) {
				return fmt.Errorf("field '%s' must be at most %d characters", field.Name, maxVal)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > maxVal {
				return fmt.Errorf("field '%s' must be at most %d", field.Name, maxVal)
			}
		default:
			return fmt.Errorf("unsupported type for max validation: %s", value.Kind())
		}
		return nil
	})
}
