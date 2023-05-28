package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DBHost     = "localhost"
	DBPort     = 5432
	DBUser     = "your_username"
	DBPassword = "your_password"
	DBName     = "your_database_name"
)

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName))
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
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author_id INTEGER NOT NULL,
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
		err := rows.Scan(&book.ID, &book.Title, &book.AuthorID)
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
	id := params["id"]

	var book Book
	err := db.QueryRow("SELECT * FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.AuthorID)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("INSERT INTO books (title, author_id) VALUES ($1, $2) RETURNING id", book.Title, book.AuthorID).Scan(&book.ID)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	_, err := db.Exec("UPDATE books SET title = $1, author_id = $2 WHERE id = $3", book.Title, book.AuthorID, id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

