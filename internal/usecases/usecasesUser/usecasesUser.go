package usecasesUser

import "studentgit.kata.academy/Zhodaran/go-kata/internal/infrastructure/postgres"

type UserService struct {
	UserRepo *postgres.PostgresUserRepository
}

func NewUserService(repo *postgres.PostgresUserRepository) *UserService {
	return &UserService{UserRepo: repo}
}
