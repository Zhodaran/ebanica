package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// @Summary Add a new author to the library
// @Description This endpoint allows you to add a new author to the library.
// @Tags Authors
// @Accept json
// @Produce json
// @Param author body AuthorRequest true "Author name"
// @Success 201 {object} string "Author added successfully"
// @Failure 400 {object} mErrorResponse "Invalid request"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/authors [post]
func (a *AuthorController) AddAuthorHandler(resp Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorRequest AuthorRequest
		if err := json.NewDecoder(r.Body).Decode(&authorRequest); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		if authorRequest.Name == "" {
			resp.ErrorBadRequest(w, errors.New("author name is required"))
			return
		}

		library.mu.Lock()         // Блокируем запись
		defer library.mu.Unlock() // Разблокируем запись после завершения

		// Добавление автора в библиотеку (можно добавить логику для проверки уникальности)
		// Здесь предполагается, что у вас есть структура для хранения авторов
		library.Authors = append(library.Authors, authorRequest.Name)

		resp.OutputJSON(w, map[string]string{"message": "Author added successfully"})
	}
}

// getAuthorsHandler godoc
// @Summary Get all authors
// @Description Get a list of all authors in the library
// @Tags Authors
// @Produce json
// @Success 200 {array} string "List of authors"
// @Failure 404 {object} mErrorResponse "No authors found"
// @Router /api/get-authors [get]
func (a *AuthorController) GetAuthorsHandler(resp Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		library.mu.RLock()         // Блокируем чтение
		defer library.mu.RUnlock() // Разблокируем чтение после завершения

		// Проверяем, есть ли авторы
		if len(library.Authors) == 0 {
			http.Error(w, "No authors found", http.StatusNotFound)
			return
		}

		// Возвращаем список авторов в формате JSON
		resp.OutputJSON(w, library.Authors)
	}
}

func (a *AuthorController) ListAuthorsHandler(resp Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		library.mu.RLock()         // Блокируем чтение
		defer library.mu.RUnlock() // Разблокируем чтение после завершения

		authorsSet := make(map[string]struct{}) // Используем множество для уникальных авторов

		// Проходим по всем книгам в библиотеке и собираем авторов
		for _, books := range library.Books {
			for _, book := range books {
				authorsSet[book.Author] = struct{}{} // Добавляем автора в множество
			}
		}

		// Преобразуем множество в срез
		var authors []string
		for author := range authorsSet {
			authors = append(authors, author)
		}

		resp.OutputJSON(w, authors) // Возвращаем список авторов в формате JSON
	}
}
