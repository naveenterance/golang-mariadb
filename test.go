package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {

	cfg := mysql.Config{
		User:   "root",
		Passwd: "nst",

		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func view_by_artist(w http.ResponseWriter, r *http.Request) {

	var albums []Album

	if r.Method == "GET" {

		http.ServeFile(w, r, "webpages/submit.html")
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Error parsing form data")
			return
		}

		name := r.Form.Get("name")

		rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
		if err != nil {
			fmt.Printf("error")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var alb Album
			if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
				fmt.Printf("error")
			}

			albums = append(albums, alb)

		}

		tmpl := template.Must(template.ParseFiles("webpages/mul.html"))
		err = tmpl.Execute(w, albums)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
