package formsgen

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

const (
	TYPE_RADIO    = "radio"
	TYPE_CHECKBOX = "checkbox"
)

type RadioValue struct {
	Value   string
	Checked bool
}

func NewRadioValue(tag string) RadioValue {
	if tag == "" {
		return RadioValue{}
	}

	if splits := strings.Split(tag, ";"); len(splits) > 1 {
		return RadioValue{
			Value:   splits[0],
			Checked: splits[1] == "checked",
		}
	} else {
		return RadioValue{
			Value: tag,
		}
	}
}

type SelectValue struct {
	Value    string
	Label    string
	Selected bool
}

func NewSelectValue(value string) SelectValue {
	selectValue := SelectValue{}
	if splits := strings.Split(value, "="); len(splits) > 1 {
		selectValue.Label = splits[0]
		selectValue.Value = splits[1]
	} else {
		selectValue.Label = value
		selectValue.Value = value
	}

	if strings.Contains(selectValue.Value, ";") {
		splits := strings.Split(selectValue.Value, ";")
		selectValue.Value = splits[0]
		selectValue.Selected = splits[1] == "selected"
	}

	return selectValue
}

func NewSelectValues(tag string) []SelectValue {
	if tag == "" {
		return []SelectValue{}
	}

	values := strings.Split(tag, ",")
	selectValues := make([]SelectValue, len(values))
	for i, value := range values {
		selectValues[i] = NewSelectValue(value)
	}

	return selectValues
}

type FormField struct {
	Type         string
	Field        string
	Name         string
	Required     bool
	Default      bool
	RadioValue   RadioValue
	SelectValues []SelectValue
	Value        reflect.Value
}

func (formField *FormField) Generate() (string, error) {
	buffer := new(bytes.Buffer)
	if err := fieldsTemplate.ExecuteTemplate(buffer, formField.Type, formField); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (formField *FormField) setValue(value string) error {
	switch formField.Value.Type().Kind() {
	case reflect.Int, reflect.Int64:
		if value == "" {
			value = "0"
		}

		if intValue, err := strconv.ParseInt(value, 10, 64); err != nil {
			return err
		} else {
			formField.Value.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint64:
		if value == "" {
			value = "0"
		}

		if uintValue, err := strconv.ParseUint(value, 10, 64); err != nil {
			return err
		} else {
			formField.Value.SetUint(uintValue)
		}
	case reflect.Float64:
		if value == "" {
			value = "0.0"
		}

		if floatValue, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			formField.Value.SetFloat(floatValue)
		}
	case reflect.Bool:
		if formField.Type == TYPE_RADIO {
			formField.Value.SetBool(formField.RadioValue.Value == value)
			return nil
		} else if formField.Type == TYPE_CHECKBOX {
			formField.Value.SetBool(value == "on")
			return nil
		}

		if value == "" {
			value = "false"
		}

		if boolValue, err := strconv.ParseBool(value); err != nil {
			formField.Value.SetBool(false)
		} else {
			formField.Value.SetBool(boolValue)
		}

		return nil
	case reflect.String:
		formField.Value.SetString(value)
	}

	return nil
}

func NewFormField(structField reflect.StructField, value reflect.Value) FormField {
	formField := FormField{
		Type:         structField.Tag.Get("type"),
		Field:        structField.Tag.Get("field"),
		Name:         structField.Tag.Get("name"),
		Required:     structField.Tag.Get("required") == "true",
		Default:      structField.Tag.Get("default") == "true",
		RadioValue:   NewRadioValue(structField.Tag.Get("radio")),
		SelectValues: NewSelectValues(structField.Tag.Get("select")),
		Value:        value,
	}

	return formField
}
