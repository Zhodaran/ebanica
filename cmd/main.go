package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	_ "studentgit.kata.academy/Zhodaran/go-kata/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/controllers"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/facades"
	postgresRepo "studentgit.kata.academy/Zhodaran/go-kata/internal/infrastructure/postgres"
)

type Server struct {
	http.Server
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := postgresRepo.InitDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err) // Обработка ошибки
	}
	defer db.Close()
	postgresRepo.PullSQL()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	postgresRepo.RunMigrations(db)
	books := postgresRepo.CreateTableBook(db)
	librar := controllers.NewLibrary()
	librar.AddBooks(books)

	resp := controllers.NewResponder(logger)

	// Инициализация репозиториев
	authRepo := postgresRepo.NewPostgresAuthRepository(db)
	bookRepo := postgresRepo.NewPostgresBookRepository(db)
	authorRepo := postgresRepo.NewPostgresAuthorRepository(db)
	userRepo := postgresRepo.NewPostgresUserRepository(db)

	// Фасад
	library := facades.NewLibraryFacade(
		authRepo,
		bookRepo,
		authorRepo,
		userRepo,
	)

	// Контроллеры
	authController := controllers.NewAuthController(library)
	userController := controllers.NewUserController(library)
	bookController := controllers.NewBookController(library)
	authorController := controllers.NewAuthorController(library)

	// Роутер
	r := chi.NewRouter()
	controllers.GenerateUsers(50)

	// Middleware

	// Публичные маршруты
	r.Post("/api/register", authController.Register)
	r.Post("/api/login", authController.Login)

	// Приватные маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		// Пользователи
		r.Post("/api/users", userController.CreateUser)
		r.Get("/api/users/{id}", userController.GetUser)
		r.Put("/api/users/{id}", userController.UpdateUser)
		r.Delete("/api/users/{id}", userController.DeleteUser)
		r.Get("/api/users", userController.ListUsers)

		// Книги
		r.Post("/api/book/take/{id}", bookController.TakeBookHandler(resp, db, &books, librar))
		r.Delete("/api/book/return/{id}", bookController.ReturnBook(resp, db, &books, librar))
		r.Post("/api/books", bookController.AddBookHandler(resp, db, librar, &books))
		r.Get("/api/books", bookController.ListBooks)
		r.Put("/api/books/{id}", bookController.UpdateBook(resp, db))

		// Авторы
		r.Post("/api/authors", authorController.AddAuthorHandler(resp, librar))
		r.Get("/api/authors", authorController.ListAuthorsHandler(resp, librar))
	})

	// Запуск сервера
	srv := &Server{
		Server: http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// Создаем Listener
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	// Запускаем сервер в горутине
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	WaitForShutdown(srv)
}

func WaitForShutdown(srv *Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
