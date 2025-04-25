package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
)

// @Summary Get Geo Coordinates by Address
// @Description This endpoint allows you to get geo coordinates by address.
// @Tags User
// @Accept json
// @Produce json
// @Param index path int true "Book INDEX"
// @Param Authorization header string true "Bearer Token"
// @Param body body TakeBookRequest true "Request body"
// @Success 200 {object} Response "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/take/{index} [post]
func (l *BookController) TakeBookHandler(resp Responder, db *sql.DB, Books *[]entities.Book, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid index"))
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET block = $1, take_count = take_count + 1 WHERE index = $2 AND block = $3", true, index, false)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		var bookFind entities.Book
		found := false

		// Поиск книги по индексу
		for i, book := range *Books {
			if index == book.Index {
				bookFind = book
				// Удаление книги из массива
				*Books = append((*Books)[:i], (*Books)[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, fmt.Sprintf("book with index %d not found", index), http.StatusNotFound)
			return
		}

		// Проверка, была ли книга успешно обновлена
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("book not found or already taken"))
			return
		}

		var requestBody TakeBookRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		// Проверка, был ли передан username
		if requestBody.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		// Добавление книги к пользователю
		library.Books[requestBody.Username] = append(library.Books[requestBody.Username], bookFind)
		resp.OutputJSON(w, map[string]string{"message": "Book taken successfully"})
	}
}

// @Summary Get Geo Coordinates by Address
// @Description This endpoint allows you to get geo coordinates by address.
// @Tags User
// @Accept json
// @Produce json
// @Param index path int true "Book INDEX"
// @Param Authorization header string true "Bearer Token"
// @Param body body TakeBookRequest true "Request body"
// @Success 200 {object} Response "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/return/{index} [delete]
func (l *BookController) ReturnBook(resp Responder, db *sql.DB, Books *[]entities.Book, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid index"))
			return
		}

		var requestBody TakeBookRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		if requestBody.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		userBooks, userExists := library.Books[requestBody.Username]
		if !userExists {
			http.Error(w, "User has no books", http.StatusNotFound)
			return
		}

		found := false
		var bookFind entities.Book

		// Поиск книги у пользователя
		for i, book := range userBooks {
			if book.Index == index {
				bookFind = book
				library.Books[requestBody.Username] = append(userBooks[:i], userBooks[i+1:]...) // Удаляем книгу из списка пользователя
				found = true
				break
			}
		}

		if !found {
			http.Error(w, fmt.Sprintf("book with index %d not found for user", index), http.StatusNotFound)
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET block = $1 WHERE index = $2 AND block = $3", false, index, true)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("book not found or already returned"))
			return
		}

		// Добавление книги обратно в общий список книг
		*Books = append(*Books, bookFind) // Добавляем книгу обратно в общий список
		resp.OutputJSON(w, map[string]string{"message": "Book returned successfully"})
	}
}

// @Summary Обновление информации о книге
// @Description Этот эндпоинт позволяет обновить информацию о книге по индексу.
// @Tags Books
// @Accept json
// @Produce json
// @Param index path int true "Индекс книги"
// @Param Authorization header string true "Bearer Token"
// @Param body body models.Book true "Обновленная информация о книге"
// @Success 200 {object} models.Book "Успешное обновление книги"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 404 {object} mErrorResponse "Книга не найдена"
// @Failure 500 {object} mErrorResponse "Ошибка сервера"
// @Router /api/book/{index} [put]
func (l *BookController) UpdateBook(resp Responder, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			resp.ErrorBadRequest(w, errors.New("недопустимый индекс"))
			return
		}

		var updatedBook entities.Book
		if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
			resp.ErrorBadRequest(w, errors.New("недопустимый формат данных"))
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET book = $1, author = $2, block = $3 WHERE index = $4",
			updatedBook.Book, updatedBook.Author, updatedBook.Block, index)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		// Проверка, была ли книга успешно обновлена
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("книга не найдена или не обновлена"))
			return
		}

		// Возвращаем обновленную книгу
		resp.OutputJSON(w, updatedBook)
	}
}

// @Summary Add a new book to the library
// @Description This endpoint allows you to add a new book to the library.
// @Tags Books
// @Accept json
// @Produce json
// @Param book body repository.AddaderBook false "Book details"
// @Success 201 {object} models.Book "Book added successfully"
// @Failure 400 {object} mErrorResponse "Invalid request"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/book [post]
func (l *BookController) AddBookHandler(resp Responder, db *sql.DB, library *Library, Books *[]entities.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addaderBook AddaderBook
		if err := json.NewDecoder(r.Body).Decode(&addaderBook); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		// Проверка на существование книги
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM book WHERE book = $1 AND author = $2)", addaderBook.Book, addaderBook.Author).Scan(&exists)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}
		if exists {
			resp.ErrorBadRequest(w, errors.New("book already exists"))
			return
		}

		var newBook entities.Book
		newBook.Book = addaderBook.Book
		newBook.Author = addaderBook.Author

		bloc := false
		newBook.Block = &bloc

		// Вставка новой книги в базу данных
		_, err = db.Exec("INSERT INTO book (book, author, block) VALUES ($1, $2, $3)", newBook.Book, newBook.Author, newBook.Block)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}
		bookPtr := &Books

		// Получаем последний элемент
		lastElement := (**bookPtr)[len(**bookPtr)-1].Index
		newBook.Index = lastElement + 1

		library.AddBook(newBook)
		*Books = append(*Books, newBook)
		resp.OutputJSON(w, newBook) // Возвращаем добавленную книгу
	}
}

// @Summary List SQL book
// @Description This description created new SQL user
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {object} CreateResponse "List successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/books [get]
func (uc *BookController) ListBooks(w http.ResponseWriter, r *http.Request) {
	// Получаем список книг из базы данных
	books, err := uc.getBooksFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Устанавливаем статус 200 OK

	// Кодируем и отправляем список книг
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uc *BookController) getBooksFromDB() ([]entities.Book, error) {
	query := "SELECT index, book, author, block, take_count FROM book"
	rows, err := uc.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []entities.Book
	for rows.Next() {
		var book entities.Book
		if err := rows.Scan(&book.Index, &book.Book, &book.Author, &book.Block, &book.TakeCount); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (l *Library) AddBook(book entities.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Получаем список книг автора
	booksByAuthor := l.Books[book.Author]

	// Находим первый свободный индекс
	newIndex := 1
	for {
		found := false
		for _, b := range booksByAuthor {
			if b.Index == newIndex {
				found = true
				break
			}
		}
		if !found {
			break
		}
		newIndex++
	}

	// Присваиваем книге новый индекс
	book.Index = newIndex

	// Добавляем книгу в список
	l.Books[book.Author] = append(booksByAuthor, book)

	// Добавляем автора в список, если его там еще нет
	if !contains(l.Authors, book.Author) {
		l.Authors = append(l.Authors, book.Author)
	}
}

func contains(authors []string, author string) bool {
	for _, a := range authors {
		if a == author {
			return true
		}
	}
	return false
}
