package structs

import "reflect"

// CallbackFunc on the struct field
// example:
// structValue := reflect.ValueOf(s)
// ...
// field := structValue.Field(i)
// fieldType := structValue.Type().Field(i)
type CallbackFunc func(reflect.Value, reflect.StructField)

// Walk traverses a struct and executes a callback function on each field in the struct.
// The interface{} passed to the function should be a pointer to a struct
func Walk(s interface{}, callback CallbackFunc) {
	structValue := reflect.ValueOf(s)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}
	if structValue.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structValue.Type().Field(i)
		if !fieldType.IsExported() {
			continue
		}
		if field.Kind() == reflect.Struct {
			Walk(field.Addr().Interface(), callback)
		} else if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			Walk(field.Interface(), callback)
		} else {
			callback(field, fieldType)
		}
	}
}
