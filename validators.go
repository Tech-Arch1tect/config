package config

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	RegisterValidator("required", RequiredValidator)
	RegisterValidator("min", MinValidator)
	RegisterValidator("max", MaxValidator)
	RegisterValidator("email", EmailValidator)
	RegisterValidator("url", URLValidator)
	RegisterValidator("regexp", RegexpValidator)
	RegisterValidator("in", InValidator)
	RegisterValidator("not_in", NotInValidator)
	RegisterValidator("eq", EqValidator)
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

func EmailValidator(field reflect.StructField, value reflect.Value, param string) error {
	_, err := mail.ParseAddress(value.String())
	if err != nil {
		return fmt.Errorf("field '%s' must be a valid email address", field.Name)
	}
	return nil
}

func RegexpValidator(field reflect.StructField, value reflect.Value, param string) error {
	regex, err := regexp.Compile(param)
	if err != nil {
		return fmt.Errorf("invalid regex pattern for field '%s'", field.Name)
	}
	if !regex.MatchString(value.String()) {
		return fmt.Errorf("field '%s' must match the pattern '%s'", field.Name, param)
	}
	return nil
}

func URLValidator(field reflect.StructField, value reflect.Value, param string) error {
	parsed, err := url.ParseRequestURI(value.String())
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("field '%s' must be a valid URL", field.Name)
	}
	return nil
}

func InValidator(field reflect.StructField, value reflect.Value, param string) error {
	allowedValues := strings.Split(param, "|")
	for _, allowedValue := range allowedValues {
		if value.String() == allowedValue {
			return nil
		}
	}
	return fmt.Errorf("field '%s' must be one of the following values: %s", field.Name, param)
}

func NotInValidator(field reflect.StructField, value reflect.Value, param string) error {
	allowedValues := strings.Split(param, "|")
	for _, allowedValue := range allowedValues {
		if value.String() == allowedValue {
			return fmt.Errorf("field '%s' must not be one of the following values: %s", field.Name, param)
		}
	}
	return nil
}

func EqValidator(field reflect.StructField, value reflect.Value, param string) error {
	if value.String() != param {
		return fmt.Errorf("field '%s' must be equal to '%s'", field.Name, param)
	}
	return nil
}
