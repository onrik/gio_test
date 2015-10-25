package formsgen

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

type TestForm struct {
	Username   string        `field:"name" name:"Имя пользователя" type:"text" required:"true" `
	Password   string        `field:"password" name:"Пароль пользователя" type:"password" required:"true" `
	Resident   bool          `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident bool          `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
	Gender     string        `field:"gender" name:"Пол" type:"select" select:"Не указан=0;selected,М=1,Ж=2"`
	Agree      bool          `field:"agree" type:"checkbox" name:"Согласен с условиями" default:"true"`
	Age        uint          `field:"age" name:"Возраст" type:"text" default:"true"`
	Token      string        `field:"token" type:"hidden" default:"true"`
	Duration   time.Duration `field:"duration" type:"hidden"`
}

func isFormFieldsEqual(field1, field2 FormField) bool {
	if field1.Name != field2.Name {
		return false
	}

	if field1.Type != field2.Type {
		return false
	}

	if field1.Field != field2.Field {
		return false
	}

	if field1.Required != field2.Required {
		return false
	}

	if field1.Default != field2.Default {
		return false
	}

	if field1.RadioValue != field2.RadioValue {
		return false
	}

	if !reflect.DeepEqual(field1.Value.Interface(), field2.Value.Interface()) {
		return false
	}

	return reflect.DeepEqual(field1.SelectValues, field2.SelectValues)
}

func TestNewFormField(t *testing.T) {
	testForm := TestForm{
		Username: "testuser",
		Age:      20,
		Token:    "a628228b1089458da6ff6e58d979bb65",
		Agree:    true,
	}

	testFields := map[string]FormField{
		"Username": FormField{
			Type:         "text",
			Field:        "name",
			Name:         "Имя пользователя",
			Required:     true,
			SelectValues: []SelectValue{},
			Value:        reflect.ValueOf(testForm.Username),
		},
		"Password": FormField{
			Type:         "password",
			Field:        "password",
			Name:         "Пароль пользователя",
			Required:     true,
			SelectValues: []SelectValue{},
			Value:        reflect.ValueOf(""),
		},
		"Resident": FormField{
			Type:         "radio",
			Field:        "resident",
			Name:         "Резидент РФ",
			Required:     false,
			SelectValues: []SelectValue{},
			RadioValue:   RadioValue{"1", true},
			Value:        reflect.ValueOf(false),
		},
		"NoResident": FormField{
			Type:         "radio",
			Field:        "resident",
			Name:         "Не резидент РФ",
			Required:     false,
			SelectValues: []SelectValue{},
			RadioValue:   RadioValue{"2", false},
			Value:        reflect.ValueOf(false),
		},
		"Gender": FormField{
			Type:  "select",
			Field: "gender",
			Name:  "Пол",
			SelectValues: []SelectValue{
				SelectValue{"0", "Не указан", true},
				SelectValue{"1", "М", false},
				SelectValue{"2", "Ж", false},
			},
			Value: reflect.ValueOf(""),
		},
		"Agree": FormField{
			Type:         "checkbox",
			Field:        "agree",
			Name:         "Согласен с условиями",
			Default:      true,
			SelectValues: []SelectValue{},
			Value:        reflect.ValueOf(testForm.Agree),
		},
		"Age": FormField{
			Type:         "text",
			Field:        "age",
			Name:         "Возраст",
			Default:      true,
			SelectValues: []SelectValue{},
			Value:        reflect.ValueOf(uint(20)),
		},
		"Token": FormField{
			Type:         "hidden",
			Field:        "token",
			Default:      true,
			SelectValues: []SelectValue{},
			Value:        reflect.ValueOf(testForm.Token),
		},
	}

	formType := reflect.TypeOf(testForm)
	formValue := reflect.ValueOf(testForm)

	for fieldName, formField := range testFields {
		field, _ := formType.FieldByName(fieldName)
		value := formValue.FieldByName(fieldName)

		if v := NewFormField(field, value); !isFormFieldsEqual(v, formField) {
			t.Errorf("Invalid FormField for %s", fieldName)
		}
	}
}

func TestRenderFormField(t *testing.T) {
	testFields := map[*FormField]string{
		&FormField{
			Type:     "text",
			Field:    "name",
			Name:     "Имя пользователя",
			Required: true,
		}: `<label for="name">Имя пользователя</label>
    <input type="text" name="name" required>`,
		&FormField{
			Type:  "password",
			Field: "password",
			Name:  "Пароль",
		}: `<label for="password">Пароль</label>
    <input type="password" name="password">`,
		&FormField{
			Type:       "radio",
			Field:      "resident",
			Name:       "Резидент РФ",
			RadioValue: RadioValue{"1", true},
		}: `<label>
        <input type="radio" name="resident" value="1" checked> Резидент РФ
    </label>`,
		&FormField{
			Type:       "radio",
			Field:      "resident",
			Name:       "Не резидент РФ",
			RadioValue: RadioValue{"2", false},
		}: `<label>
        <input type="radio" name="resident" value="2"> Не резидент РФ
    </label>`,
		&FormField{
			Type:    "hidden",
			Field:   "token",
			Default: true,
			Value:   reflect.ValueOf("a628228b1089458da6ff6e58d979bb65"),
		}: `<input type="hidden" name="token" value="a628228b1089458da6ff6e58d979bb65">`,
		&FormField{
			Type:    "checkbox",
			Field:   "agree",
			Name:    "Согласен с условиями",
			Default: true,
			Value:   reflect.ValueOf(true),
		}: `<label>
        <input type="checkbox" name="agree" checked> Согласен с условиями
    </label>`,
		&FormField{
			Type:    "checkbox",
			Field:   "agree",
			Name:    "Согласен с условиями",
			Default: false,
			Value:   reflect.ValueOf(true),
		}: `<label>
        <input type="checkbox" name="agree"> Согласен с условиями
    </label>`,
		&FormField{
			Type:    "checkbox",
			Field:   "agree",
			Name:    "Согласен с условиями",
			Default: true,
			Value:   reflect.ValueOf(false),
		}: `<label>
        <input type="checkbox" name="agree"> Согласен с условиями
    </label>`,
		&FormField{
			Type:  "select",
			Field: "gender",
			Name:  "Пол",
			SelectValues: []SelectValue{
				SelectValue{"0", "Не указан", true},
				SelectValue{"1", "М", false},
				SelectValue{"2", "Ж", false},
			},
		}: `<label for="gender">Пол</label>
    <select name="gender">
        <option value="0" selected>Не указан</option>
        <option value="1" >М</option>
        <option value="2" >Ж</option>
    </select>`,
	}

	for field, html := range testFields {
		generated, err := field.Generate()
		if err != nil {
			t.Error(err.Error())
		} else if strings.TrimSpace(generated) != html {
			t.Errorf("Invalid html for %+v", field)
		}
	}
}

func TestNewRadioValue(t *testing.T) {
	testValues := map[string]RadioValue{
		"":          RadioValue{},
		"1":         RadioValue{"1", false},
		"2;checked": RadioValue{"2", true},
	}

	for tag, value := range testValues {
		if NewRadioValue(tag) != value {
			t.Errorf("Invalid RadioValue for tag '%s'", tag)
		}
	}
}

func TestNewSelectValue(t *testing.T) {
	testValues := map[string]SelectValue{
		"":                     SelectValue{},
		"value":                SelectValue{"value", "value", false},
		"label=value;selected": SelectValue{"value", "label", true},
	}

	for tag, value := range testValues {
		if NewSelectValue(tag) != value {
			t.Errorf("Invalid SelectValue for tag '%s'", tag)
		}
	}
}

func TestNewSelectValues(t *testing.T) {
	selectValues := NewSelectValues("Не указан=0;selected,М=1,Ж=2")
	if len(selectValues) != 3 {
		t.Error("Invalid select values count")
		return
	}

	referenseValues := []SelectValue{
		SelectValue{"0", "Не указан", true},
		SelectValue{"1", "М", false},
		SelectValue{"2", "Ж", false},
	}

	for i := range selectValues {
		if selectValues[i] != referenseValues[i] {
			t.Errorf("%+v != %+v", selectValues[i], referenseValues[i])
			return
		}
	}
}
