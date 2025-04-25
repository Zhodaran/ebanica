package usecasesAuthor

import (
	"studentgit.kata.academy/Zhodaran/go-kata/internal/infrastructure/postgres"
)

type AuthorService struct {
	UserRepo *postgres.PostgresAuthorRepository
}

func NewAuthorService(repo *postgres.PostgresAuthorRepository) *AuthorService {
	return &AuthorService{UserRepo: repo}
}
