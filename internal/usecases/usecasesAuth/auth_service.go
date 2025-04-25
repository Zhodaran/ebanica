package usecasesAuth

import (
	"studentgit.kata.academy/Zhodaran/go-kata/internal/infrastructure/postgres"
)

type AuthService struct {
	UserRepo *postgres.PostgresAuthRepository
}

func NewAuthService(repo *postgres.PostgresAuthRepository) *AuthService {
	return &AuthService{UserRepo: repo}
}
