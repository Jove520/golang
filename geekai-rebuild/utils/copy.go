package utils

import (
	"errors"
	"reflect"
)

func CopyObject(dst, src any) error {
	if dst == nil || src == nil {
		return errors.New("dst and src cannot be nil")
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Pointer || dstValue.IsNil() {
		return errors.New("dst must be a non-nil pointer")
	}

	dstValue = dstValue.Elem()
	srcValue := reflect.ValueOf(src)

	if srcValue.Kind() == reflect.Pointer {
		if srcValue.IsNil() {
			return errors.New("src must not be nil")
		}

		srcValue = srcValue.Elem()
	}

	if dstValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		return errors.New("dst and src must be structs")
	}

	srcType := srcValue.Type()
	for i := 0; i < srcValue.NumField(); i++ {
		srcFieldType := srcType.Field(i)
		if srcFieldType.PkgPath != "" {
			continue
		}

		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcFieldType.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
			continue
		}

		if srcField.Type().ConvertibleTo(dstField.Type()) {
			dstField.Set(srcField.Convert(dstField.Type()))
		}
	}

	return nil
}
