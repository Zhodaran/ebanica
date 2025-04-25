package facades

import (
	"context"

	"studentgit.kata.academy/Zhodaran/go-kata/internal/infrastructure/postgres"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/usecases/usecasesAuth"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/usecases/usecasesAuthor"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/usecases/usecasesBook"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/usecases/usecasesUser"
)

type LibraryFacade struct {
	AuthService   *usecasesAuth.AuthService
	BookService   *usecasesBook.BookService
	AuthorService *usecasesAuthor.AuthorService
	UserService   *usecasesUser.UserService
	QueryContext  context.Context
}

func NewLibraryFacade(authRepo *postgres.PostgresAuthRepository, bookRepo *postgres.PostgresBookRepository, authorRepo *postgres.PostgresAuthorRepository, userRepo *postgres.PostgresUserRepository) *LibraryFacade {
	return &LibraryFacade{
		AuthService:   usecasesAuth.NewAuthService(authRepo),
		BookService:   usecasesBook.NewBookService(bookRepo),
		AuthorService: usecasesAuthor.NewAuthorService(authorRepo),
		UserService:   usecasesUser.NewUserService(userRepo),
	}
}
