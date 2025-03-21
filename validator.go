package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
			switch {
			case rule == "required":
				if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
					return fmt.Errorf("field '%s' is required", field.Name)
				}
			case strings.HasPrefix(rule, "min="):
				minValStr := strings.TrimPrefix(rule, "min=")
				minVal, err := strconv.ParseInt(minValStr, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid min value for field '%s'", field.Name)
				}
				switch fieldValue.Kind() {
				case reflect.String:
					if len(fieldValue.String()) < int(minVal) {
						return fmt.Errorf("field '%s' must be at least %d characters", field.Name, minVal)
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if fieldValue.Int() < minVal {
						return fmt.Errorf("field '%s' must be at least %d", field.Name, minVal)
					}
				default:
					return fmt.Errorf("unsupported type for min validation: %s", fieldValue.Kind())
				}
			case strings.HasPrefix(rule, "max="):
				maxValStr := strings.TrimPrefix(rule, "max=")
				maxVal, err := strconv.ParseInt(maxValStr, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid max value for field '%s'", field.Name)
				}
				switch fieldValue.Kind() {
				case reflect.String:
					if len(fieldValue.String()) > int(maxVal) {
						return fmt.Errorf("field '%s' must be at most %d characters", field.Name, maxVal)
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if fieldValue.Int() > maxVal {
						return fmt.Errorf("field '%s' must be at most %d", field.Name, maxVal)
					}
				default:
					return fmt.Errorf("unsupported type for max validation: %s", fieldValue.Kind())
				}
			}
		}
	}

	return nil
}
