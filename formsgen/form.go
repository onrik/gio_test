package formsgen

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	FIELD_TAG = "field"
)

var availableTypes = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Uint,
	reflect.Int64,
	reflect.Uint64,
	reflect.Float64,
	reflect.String,
}

func isFieldTypeAvailable(kind reflect.Kind) bool {
	for i := range availableTypes {
		if kind == availableTypes[i] {
			return true
		}
	}

	return false
}

func MapForm(formStruct reflect.Value, values map[string][]string) error {
	if formStruct.Kind() == reflect.Ptr {
		formStruct = formStruct.Elem()
	}

	formType := formStruct.Type()
	for i := 0; i < formType.NumField(); i++ {
		typeField := formType.Field(i)
		structField := formStruct.Field(i)

		if !structField.CanSet() {
			continue
		}

		if fieldName := typeField.Tag.Get(FIELD_TAG); strings.Trim(fieldName, "-") != "" {
			if inputValue, exists := values[fieldName]; exists {
				formField := NewFormField(typeField, structField)
				if err := formField.setValue(inputValue[0]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func GenerateForm(form interface{}) (string, error) {
	formType := reflect.TypeOf(form)
	if formType.Kind() != reflect.Struct {
		return "", fmt.Errorf("Form is not struct")
	}

	fields := []string{}
	formValue := reflect.ValueOf(form)
	for i := 0; i < formType.NumField(); i++ {
		field := formType.Field(i)
		if f := field.Tag.Get(FIELD_TAG); strings.Trim(f, "-") == "" {
			continue
		}

		if !isFieldTypeAvailable(field.Type.Kind()) {
			return "", fmt.Errorf("Field type '%s' is not available", field.Type.Kind())
		}

		formField := NewFormField(field, formValue.Field(i))
		html, err := formField.Generate()
		if err != nil {
			return "", err
		}

		fields = append(fields, html)
	}

	return strings.Join(fields, "<br>"), nil
}
