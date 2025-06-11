package utils

import (
	"errors"
	"reflect"
)

func HasAnyNonEmptyFields(structPointer any) (bool, error) {
	v := reflect.ValueOf(structPointer)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false, errors.New("ожидалась структура или указатель на структуру")
	}

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := v.Type().Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		if !fieldVal.IsZero() {
			return true, nil
		}
	}

	return false, nil
}
