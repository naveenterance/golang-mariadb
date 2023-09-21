package main

import (
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("webpages/tpl.gohtml"))
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8080", nil)
}

func foo(reswt http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(reswt, "tpl.gohtml", 3442)
	if err != nil {
		log.Fatalln(err)
	}

}
