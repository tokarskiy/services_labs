package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	initializeDatabase()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/api/books", getBooks)
	router.POST("/api/books", postBook)
	router.PUT("/api/books", putBook)
	router.DELETE("/api/books", deleteBook)

	router.Run(":80")
}

// Book - структура, представляющая книгу
type Book struct {
	ID     int64  `json:"id"`
	Name   string `json:"bookName"`
	Author string `json:"authorName"`
}

func getDbObject() (*sql.DB, error) {
	return sql.Open("postgres", "postgres://postgres@db/postgres?sslmode=disable")
}

func initializeDatabase() {
	db, err := getDbObject()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	queryString := `
		CREATE TABLE IF NOT EXISTS "Books" (
			"bookId"     SERIAL       PRIMARY KEY NOT NULL,
			"bookName"   VARCHAR(100)             NOT NULL,
			"authorName" VARCHAR(100)             NOT NULL 
		);
                

                DELETE FROM "Books";
                INSERT INTO "Books" ("bookName", "authorName") VALUES
                       ('Solaris', 'Stanislaw Lem'),
                       ('1984', 'George Orwell'),
                       ('Harry Potter and Philoshoper Stone', 'J.K. Rowling');
	`

	if _, err := db.Query(queryString); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getBooks(c *gin.Context) {
	db, err := getDbObject()
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error connecting database: %q", err))

		return
	}
	defer db.Close()

	queryString := `
		SELECT row_number() OVER()     AS "RowNumber",
			   "Books"."bookId"        AS "ID",
			   "Books"."bookName"      AS "BookName", 
			   "Books"."authorName"    AS "AuthorName"
		  FROM "Books"
		 ORDER BY "Books"."bookName"
	`

	var rows *sql.Rows
	rows, err = db.Query(queryString)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading books: %q", err))

		return
	}
	defer rows.Close()

	var books []Book
	var book Book
	books = make([]Book, 0)
	for rows.Next() {
		var lineNum int
		err = rows.Scan(&lineNum, &book.ID, &book.Name, &book.Author)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error data scanning: %q", err))
		}
		books = append(books, book)
	}

	c.JSON(200, books)
}

func postBook(c *gin.Context) {
	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error binding JSON: %q", err))
		return
	}

	db, err := getDbObject()
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error connecting database: %q", err))

		return
	}
	defer db.Close()

	queryString := `
		INSERT INTO "Books" ("bookName", "authorName")
		     VALUES ($1, $2);
	`

	if _, err = db.Query(queryString, book.Name, book.Author); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error inserting book: %q", err))
		return
	}

	c.String(http.StatusOK, "The book is successfully added")
}

func putBook(c *gin.Context) {
	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error binding JSON: %q", err))
		return
	}

	db, err := getDbObject()
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error connecting database: %q", err))

		return
	}
	defer db.Close()

	queryString := `
		UPDATE "Books"
		   SET "bookName" = $2, 
		       "authorName" = $3
		 WHERE "bookId" = $1
	`

	if _, err = db.Query(queryString, book.ID, book.Name, book.Author); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error updating book: %q", err))
		return
	}

	c.String(http.StatusOK, "The book is successfully updated!")
}

func deleteBook(c *gin.Context) {
	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error binding JSON: %q", err))
		return
	}

	db, err := getDbObject()
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error connecting database: %q", err))

		return
	}
	defer db.Close()

	queryString := `
		DELETE FROM "Books"
		      WHERE "Books"."bookId" = $1;
	`

	if _, err = db.Query(queryString, book.ID); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error updating book: %q", err))
		return
	}

	c.String(http.StatusOK, "The book is successfully removed")
}
