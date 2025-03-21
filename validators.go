package config

import (
	"fmt"
	"reflect"
	"strconv"
)

func init() {
	RegisterValidator("required", RequiredValidator)
	RegisterValidator("min", MinValidator)
	RegisterValidator("max", MaxValidator)

}

func RequiredValidator(field reflect.StructField, value reflect.Value, param string) error {
	if reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface()) {
		return fmt.Errorf("field '%s' is required", field.Name)
	}
	return nil
}

func MinValidator(field reflect.StructField, value reflect.Value, param string) error {
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
}

func MaxValidator(field reflect.StructField, value reflect.Value, param string) error {
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
}
