package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenAuth *jwtauth.JWTAuth
	usersDB   = map[string]string{} // Простая "база данных" пользователей (email -> хэш пароля)
)

// registerHandler — регистрация пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Декодируем тело запроса
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверка, существует ли пользователь
	if _, exists := usersDB[req.Email]; exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Хэшируем пароль с помощью bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя и хэш пароля в "базу данных"
	usersDB[req.Email] = string(hashedPassword)
	render.JSON(w, r, map[string]string{"message": "User registered successfully"})
}

// loginHandler — аутентификация пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Декодируем тело запроса
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	storedHash, exists := usersDB[req.Email]
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Если аутентификация успешна, создаём JWT-токен
	claims := map[string]interface{}{
		"user_id": req.Email,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Токен истекает через 1 час
	}
	_, tokenString, _ := tokenAuth.Encode(claims)

	// Возвращаем токен клиенту
	render.JSON(w, r, map[string]string{"token": tokenString})
}

// protectedHandler — защищенный маршрут
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	// Проверяем, что claims содержит user_id
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Защищённый контент для пользователя: %v", userID)
}
