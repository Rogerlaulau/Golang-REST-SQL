package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"fmt"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

)

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	AuthorID  int    `json:"author_id"`
	Published int    `json:"published"`
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables()

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func createTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author_id INTEGER NOT NULL,
			published INTEGER NOT NULL,
			FOREIGN KEY (author_id) REFERENCES authors (id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.AuthorID, &book.Published)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var book Book
	err = db.QueryRow("SELECT * FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title, &book.AuthorID, &book.Published)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("INSERT INTO books (title, author_id, published) VALUES (?, ?, ?)", book.Title, book.AuthorID, book.Published)
	if err != nil {
		log.Fatal(err)
	}

	// Get the last inserted row ID
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	book.ID = int(lastID) // NOTE: int() to prevent "cannot use lastID (type int64) as type int in assignment"

	fmt.Printf("Last inserted row ID: %d\n", lastID)
		
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	_, err = db.Exec("UPDATE books SET title = ?, author_id = ?, published = ? WHERE id = ?", book.Title, book.AuthorID, book.Published, id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
