package main

import (
	"./formsgen"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
)

const (
	MAX_MEMORY = int64(1024 * 1024 * 10)
)

type MyForm struct {
	Username   string `field:"name" name:"Имя пользователя" type:"text" required:"true" `
	Password   string `field:"password" name:"Пароль" type:"password" required:"true" `
	Resident   bool   `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident bool   `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
	Gender     string `field:"gender" name:"Пол" type:"select" select:"Не указан=0;selected,М=1,Ж=2"`
	Age        int64  `field:"age" name:"Возраст" type:"text" default:"true"`
	Agree      bool   `field:"agree" type:"checkbox" name:"Согласен с условиями" default:"true"`
	Token      string `field:"token" type:"hidden" default:"true"`
}

func FormRead(form *MyForm, request *http.Request) error {
	if err := request.ParseMultipartForm(MAX_MEMORY); err != nil {
		return err
	}
	if err := formsgen.MapForm(reflect.ValueOf(form), request.Form); err != nil {
		return err
	}

	fmt.Println("Received data:")
	fmt.Println(form)

	return nil
}

func IndexHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if err := FormRead(&MyForm{}, request); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		} else {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("ok"))
		}
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	form := MyForm{
		Age:   20,
		Token: "a628228b1089458da6ff6e58d979bb65",
		Agree: true,
	}

	formHtml, err := formsgen.GenerateForm(form)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	fmt.Println("Generated form:")
	fmt.Println(formHtml)

	rw.Header().Set("Content-Type", "text/html")
	tmpl.Execute(rw, map[string]interface{}{
		"form": template.HTML(formHtml),
	})

}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}
