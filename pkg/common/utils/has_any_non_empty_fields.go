package utils

import (
	"errors"
	"reflect"
)

func HasAnyNonEmptyFields(structPointer any, skipFields ...string) (bool, error) {
	v := reflect.ValueOf(structPointer)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false, errors.New("ожидалась структура или указатель на структуру")
	}

	skip := make(map[string]struct{}, len(skipFields))
	for i := 0; i < len(skip); i++ {
		skip[skipFields[i]] = struct{}{}
	}

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := v.Type().Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		if _, skipped := skip[fieldType.Name]; skipped {
			continue
		}

		if !fieldVal.IsZero() {
			return true, nil
		}
	}

	return false, nil
}
