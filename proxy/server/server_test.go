package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Тест для маршрута /test
func TestHandleTest(t *testing.T) {
	// Создаем новый тестовый HTTP-запрос
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый HTTP-ответ
	rr := httptest.NewRecorder()

	// Создаем обработчик и вызываем его
	handler := http.HandlerFunc(HandleTest)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код ответа: ожидается %v, получено %v", http.StatusOK, status)
	}
}

// Тест для маршрута /swagger.yaml
func TestSwaggerFile(t *testing.T) {
	req, err := http.NewRequest("GET", "/swagger.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	server := NewServer()
	server.handler.ServeHTTP(rr, req)

	// Проверяем код ответа
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Неверный код ответа для отсутствующего swagger.yaml: ожидается %v, получено %v", http.StatusNotFound, status)
	}
}

// Тест остановки сервера
func TestServerStop(t *testing.T) {
	server := NewServer()

	errChan := make(chan error, 1)
	// Запускаем сервер в отдельной горутине
	go func() {
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()
	select {
	case err := <-errChan:
		t.Fatalf("Ошибка запуска сервера: %v", err)
	case <-time.After(2 * time.Second):
		// Останавливаем сервер
		err := server.Stop()
		if err != nil {
			t.Fatalf("Ошибка остановки сервера: %v", err)
		}
	}
}

func TestHandleSearch(t *testing.T) {
	// Создание тестового запроса
	requestBody, _ := json.Marshal(SearchRequest{Query: "test query"})
	req, err := http.NewRequest(http.MethodPost, "/api/address/search", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	// Запуск тестового сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleSearch)

	// Вызов обработчика
	handler.ServeHTTP(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// var response []Address
	// err = json.Unmarshal(rr.Body.Bytes(), &response)
	// assert.NoError(t, err)

}

func TestHandleGeocode(t *testing.T) {
	// Создание тестового запроса
	requestBody, _ := json.Marshal(GeocodeRequest{Lat: "55.878", Lng: "37.653"})
	req, err := http.NewRequest(http.MethodPost, "/api/address/geocode", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	// Запуск тестового сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGeocode)

	// Вызов обработчика
	handler.ServeHTTP(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Логирование тела ответа для диагностики
	body := rr.Body.String()
	t.Logf("Response body: %s", body)

	// // Проверка содержимого ответа
	// var response []ResponseAddress
	// err = json.Unmarshal(rr.Body.Bytes(), &response)
	// assert.NoError(t, err)

}

func TestRegisterHandler(t *testing.T) {
	// Очищаем базу данных перед тестом
	usersDB = map[string]string{}

	// Создаем данные для запроса
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	// Создаем HTTP-запрос
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Вызываем обработчик
	registerHandler(w, req)

	// Проверяем код статуса
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем тело ответа как JSON
	var resp map[string]string
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	expected := map[string]string{"message": "User registered successfully"}
	if resp["message"] != expected["message"] {
		t.Errorf("handler returned unexpected body: got %v want %v", resp["message"], expected["message"])
	}

	// Проверяем, был ли добавлен пользователь в базу данных
	if _, exists := usersDB["test@example.com"]; !exists {
		t.Errorf("user was not added to the database")
	}
}

func TestLoginHandler(t *testing.T) {
	// Очищаем базу данных и добавляем тестового пользователя
	usersDB = map[string]string{}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	usersDB["test@example.com"] = string(hashedPassword)

	// Создаем данные для запроса
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	// Создаем HTTP-запрос
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Вызываем обработчик
	loginHandler(w, req)

	// Проверяем код статуса
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что токен был сгенерирован
	var resp map[string]string
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil || resp["token"] == "" {
		t.Errorf("expected a token but got: %v", w.Body.String())
	}
}
