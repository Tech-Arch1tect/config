package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type DefaultSetter interface {
	SetDefaults()
}

func Load[T DefaultSetter](cfg T) error {
	cfg.SetDefaults()

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}
		if envValue, ok := os.LookupEnv(envTag); ok {
			f := v.Field(i)
			if !f.CanSet() {
				continue
			}
			switch f.Kind() {
			case reflect.String:
				// skip if empty string
				if envValue == "" {
					continue
				}
				f.SetString(envValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intVal, err := strconv.ParseInt(envValue, 10, 64); err == nil {
					f.SetInt(intVal)
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(envValue); err == nil {
					f.SetBool(boolVal)
				}
			default:
				return fmt.Errorf("unsupported type: %s", f.Kind())
			}
		}
	}

	if err := ValidateStruct(cfg); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
