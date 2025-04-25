package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
)

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uc.UserRepo.Create(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Create successful"})
	w.WriteHeader(http.StatusCreated)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := uc.UserRepo.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Great successful"})
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uc.UserRepo.Update(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Update successful"})
	w.WriteHeader(http.StatusNoContent)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uc.UserRepo.Delete(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Delete successful"})
	w.WriteHeader(http.StatusNoContent)
}

func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Установите значение по умолчанию
	offset := 0 // Установите значение по умолчанию
	users, err := uc.UserRepo.List(context.Background(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "List successful"})
	json.NewEncoder(w).Encode(users)
}

func (r *UserController) Create(ctx context.Context, user entities.User) error {
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
	_, err := r.Db.ExecContext(ctx, query, user.ID, user.Name, user.Email)
	return err
}

// GetByID получает пользователя по ID
func (r *UserController) GetByID(ctx context.Context, id string) (entities.User, error) {
	var user entities.User
	query := "SELECT id, name, email, deleted_at FROM users WHERE id = $1"
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

// Update обновляет данные пользователя
func (r *UserController) Update(ctx context.Context, user entities.User) error {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
	_, err := r.Db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	return err
}

// Delete помечает пользователя как удаленного
func (r *UserController) Delete(ctx context.Context, id string) error {
	query := "UPDATE users SET deleted_at = NOW() WHERE id = $1"
	_, err := r.Db.ExecContext(ctx, query, id)
	return err
}

// List возвращает список пользователей с пагинацией
func (r *UserController) List(ctx context.Context, limit, offset int) ([]entities.User, error) {
	query := "SELECT id, name, email, deleted_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2"
	rows, err := r.Db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
