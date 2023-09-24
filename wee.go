package main

import (
	"log"
	"net/http"
	"text/template"
)

type UserData struct {
	Name  string
	Email string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("webpages/mul.gohtml"))
}

func main() {

	http.HandleFunc("/", foo)
	http.ListenAndServe(":8080", nil)
}

func foo(reswt http.ResponseWriter, req *http.Request) {
	user := UserData{
		Name:  "John Doe",
		Email: "johnddddoe@example.com",
	}
	err := tpl.ExecuteTemplate(reswt, "mul.gohtml", user)
	if err != nil {
		log.Fatalln(err)
	}

}
