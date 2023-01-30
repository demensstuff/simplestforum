package helpers

import (
	"reflect"
)

// ProcessExportedNonEmptyFields invokes cb on each exported non-empty field of the struct in.
func ProcessExportedNonEmptyFields(in interface{}, cb func(reflect.Value, reflect.StructField)) {
	// If this is a pointer, dereference it
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If this is not a struct, panic
	if val.Kind() != reflect.Struct {
		panic("ProcessExportedNonEmptyFields: argument is not a struct")
	}

	t := val.Type()

	// Iterating over the fields
	for i := 0; i < val.NumField(); i++ {
		structField := t.Field(i)

		// If the fieldVal is not exported, skip
		if !structField.IsExported() {
			continue
		}

		fieldVal := val.Field(i)

		// If the fieldVal is not a pointer, pass to the callback now
		if fieldVal.Kind() != reflect.Ptr {
			cb(fieldVal, structField)

			continue
		}

		// If it is nil, skip
		if fieldVal.IsNil() {
			continue
		}

		cb(fieldVal.Elem(), structField)
	}
}
