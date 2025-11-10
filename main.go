package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/lichess4/books-api-go/internal/service"
	"github.com/lichess4/books-api-go/internal/transport"
	"github.com/lichess4/books-api-go/store"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL
	);`

	if _, err := db.Exec(q); err != nil {
		log.Fatal(err)
	}

	// Here you would typically initialize your store, service, and transport layers
	// and start your HTTP server. This is omitted for brevity.

	bookStore := store.New(db)
	bookService := service.New(bookStore)
	bookHandler := transport.New(bookService)

	// Set up your HTTP server and routes here (omitted for brevity)
	http.HandleFunc("/books", bookHandler.HandleBooks)
	http.HandleFunc("/books/", bookHandler.HandleBookByID)

	fmt.Println("Server is running on port 8080")
	fmt.Println("API Endpoints:")
	fmt.Println(" GET    /books          - Get all books")
	fmt.Println(" POST   /books          - Create a new book")
	fmt.Println(" GET    /books/{id}     - Get a book by ID")
	fmt.Println(" PUT    /books/{id}     - Update a book by ID")
	fmt.Println(" DELETE /books/{id}     - Delete a book by ID")

	log.Fatal(http.ListenAndServe(":8000", nil))
}
