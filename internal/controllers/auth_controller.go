package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
)

func (s *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var user entities.UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя из мапы Users
	storedUser, exists := Users[user.Username]
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем совпадение пароля
	if storedUser.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Если авторизация успешна, создаем токен
	claims := map[string]interface{}{
		"user_id": user.Username, // Используем username как user_id
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	fmt.Println(tokenString)
}

func (s *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var user entities.UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := Users[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	Users[user.Username] = entities.UserAuth{
		Username: user.Username,
		Password: user.Password,
	}

	// Используем логин пользователя в качестве user_id

}
