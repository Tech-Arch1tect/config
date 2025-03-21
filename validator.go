package config

import (
	"fmt"
	"reflect"
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
