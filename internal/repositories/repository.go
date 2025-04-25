package repositories

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
)

type Gnida interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type BookRepository interface {
	AddBooks(books []entities.Book)
	TakeBookHandler(resp Responder, db *sql.DB, Books *[]entities.Book, library *Library) http.HandlerFunc
	ReturnBook(resp Responder, db *sql.DB, Books *[]entities.Book, library *Library) http.HandlerFunc
	UpdateBook(resp Responder, db *sql.DB) http.HandlerFunc
	AddBookHandler(resp Responder, db *sql.DB, library *Library, Books *[]entities.Book) http.HandlerFunc
}

type UserRepository interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	Create(ctx context.Context, user entities.User) error
	GetByID(ctx context.Context, id string) (entities.User, error)
	Update(ctx context.Context, user entities.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]entities.User, error)
}

type AuthorRepository interface {
	GetAuthorsHandler(resp Responder, library *Library) http.HandlerFunc
	ListAuthorsHandler(resp Responder, library *Library) http.HandlerFunc
	AddAuthorHandler(resp Responder, library *Library) http.HandlerFunc
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})
	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type Library struct {
	Books   map[string][]entities.Book
	Authors []string
	mu      sync.RWMutex
}
