package usecasesBook

import "studentgit.kata.academy/Zhodaran/go-kata/internal/infrastructure/postgres"

type BookService struct {
	UserRepo *postgres.PostgresBookRepository
}

func NewBookService(repo *postgres.PostgresBookRepository) *BookService {
	return &BookService{UserRepo: repo}
}
