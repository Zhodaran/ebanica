package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
)

func PullSQL() {

	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("Запуск задержки")
	time.Sleep(10 * time.Second)

}

func InitDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return sql.Open("postgres", connStr)
}

func RunMigrations(db *sql.DB) {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(255) NOT NULL,
		deleted_at TIMESTAMP NULL
	);`

	_, err := db.Exec(migrationSQL)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func CreateTableBook(db *sql.DB) []entities.Book {
	table := ` 
	CREATE TABLE IF NOT EXISTS book (
		index SERIAL PRIMARY KEY,
		book VARCHAR(50) NOT NULL,
		author VARCHAR(255) NOT NULL,
		block BOOLEAN,
		take_count INT DEFAULT 0
	);`

	var authors []string
	for i := 0; i < 10; i++ {
		author := gofakeit.Name()
		authors = append(authors, author)
	}
	_, err := db.Exec(table)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	var books []entities.Book
	for i := 1; i < 101; i++ {
		block := false
		book := entities.Book{
			Index:     i,
			Book:      gofakeit.Sentence(1),                        // Генерация названия книги
			Author:    authors[gofakeit.Number(0, len(authors)-1)], // Случайный автор
			Block:     &block,                                      // Устанавливаем значение блокировки
			TakeCount: 0,                                           // Начальное значение take_count
		}
		books = append(books, book) // Добавляем книгу в массив
	}

	// Вставка книг в базу данных
	for _, b := range books {
		_, err := db.Exec("INSERT INTO book (book, author, block) VALUES ($1, $2, $3)", b.Book, b.Author, b.Block)
		if err != nil {
			log.Fatalf("Error inserting book: %v", err)
		}
	}

	return books

}
