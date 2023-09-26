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

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("webpages/mul.html"))
}

func main() {

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "nst",
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

	var userInput int

	for i := 1; i < 10; i++ {
		fmt.Print("select from 1-5 ")
		fmt.Print("	\n1--view all ")
		fmt.Print("	\n2--view by artist ")
		fmt.Print("	\n3--add album ")
		fmt.Print("	\n4--delete album ")
		fmt.Print("	\nEnter anything else to quit :=+>")
		fmt.Scanln(&userInput)

		switch userInput {
		case 1:
			{
				http.HandleFunc("/", View_all)
				http.ListenAndServe(":8080", nil)
			}
		case 2:
			{
				http.HandleFunc("/", view_by_artist)
				http.ListenAndServe(":8080", nil)

			}
		case 3:
			{
				albID, err := addAlbum(Album{
					Title:  "The Modern Sound of Betty Carter",
					Artist: "Betty Carter",
					Price:  49.99,
				})
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("ID of added album: %v\n", albID)

			}
		case 4:
			{
				err := deletealbumsByArtist("Betty Carter")
				if err != nil {
					log.Fatal(err)
				}

			}

		default:

			return
		}

	}
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func deletealbumsByArtist(name string) error {

	var enough bool
	if err := db.QueryRow("SELECT * FROM album WHERE artist = ?", name).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			println(" already deleted")
			return nil
		}
	}

	db.Query("DELETE  FROM album WHERE artist = ?", name)

	println("deleted")
	return nil

}

func View_all(w http.ResponseWriter, r *http.Request) {

	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
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
