package formsgen

import (
	"reflect"
	"testing"
	"time"
)

func TestMapForm(t *testing.T) {
	form := new(TestForm)
	values := map[string][]string{
		"name":     {"username"},
		"age":      {"32"},
		"resident": {"1"},
		"duration": {"1000"},
	}

	err := MapForm(reflect.ValueOf(form), values)
	if err != nil {
		t.Error(err.Error())
	}

	if form.Username != "username" {
		t.Errorf("%s != username", form.Username)
	}

	if form.Age != 32 {
		t.Errorf("%s != 32", form.Age)
	}

	if form.Resident != true {
		t.Errorf("%s != true", form.Resident)
	}

	if form.NoResident != false {
		t.Errorf("%s != false", form.NoResident)
	}

	if form.Token != "" {
		t.Errorf("%s != ''", form.Token)
	}

	if form.Duration != time.Duration(1000) {
		t.Errorf("%s != ''", form.Duration)
	}
}
