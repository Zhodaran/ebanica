package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/facades"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repositories"
)

type BookController struct {
	facade *facades.LibraryFacade
	DB     *sql.DB
}

func NewBookController(facade *facades.LibraryFacade) *BookController {
	return &BookController{facade: facade}
}

type AuthController struct {
	facade *facades.LibraryFacade
}

func NewAuthController(facade *facades.LibraryFacade) *AuthController {
	return &AuthController{facade: facade}
}

type AuthorController struct {
	facade *facades.LibraryFacade
}

func NewAuthorController(facade *facades.LibraryFacade) *AuthorController {
	return &AuthorController{facade: facade}
}

type UserController struct {
	facade   *facades.LibraryFacade
	UserRepo repositories.UserRepository
	Db       *sql.DB
}

func NewUserController(facade *facades.LibraryFacade) *UserController {
	return &UserController{facade: facade}
}

type Respond struct {
	log *zap.Logger
}

func NewResponder(logger *zap.Logger) Responder {
	return &Respond{log: logger}
}

type Library struct {
	Books   map[string][]entities.Book
	Authors []string
	mu      sync.RWMutex
}

func NewLibrary() *Library {
	return &Library{
		Books: make(map[string][]entities.Book),
	}
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Info("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne Unauthorized", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type mErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

var (
	TokenAuth = jwtauth.New("HS256", []byte("your_secret_key"), nil)
	Users     = make(map[string]entities.UserAuth) // Хранение пользователей
	mu        sync.Mutex
)

type TokenResponse struct {
	Token string `json:"token"`
}

type AuthorRequest struct {
	Name string `json:"name"`
}

type TakeBookRequest struct {
	Username string `json:"username"` // Поле для имени пользователя
}

type AddaderBook struct {
	Book   string `json:"book"`
	Author string `json:"author"`
}

type CreateResponse struct {
	Message string          `json:"message"`
	Books   []entities.Book `json:"books"` // Добавляем поле для списка книг
}
