// file: auth_controller_test.go
package controllers

import (
	"bytes"
	"encoding/json"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/service"
	"github.com/go-chi/jwtauth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthController_RegisterHandler(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	authService := service.NewAuthService()
	authController := NewAuthController(authService, tokenAuth)

	// Создаем запрос с телом JSON
	user := models.User{Username: "testuser", Password: "testpassword"}
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.RegisterHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthController_LoginHandler(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	authService := service.NewAuthService()
	authController := NewAuthController(authService, tokenAuth)

	// Регистрируем пользователя
	user := models.User{Username: "testuser", Password: "testpassword"}
	authService.RegisterUser(user.Username, user.Password)

	// Создаем запрос с телом JSON для логина
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.LoginHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем содержимое ответа
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if _, exists := response["token"]; !exists {
		t.Error("Expected token in response, got none")
	}
}
func TestAuthController_LoginHandler_InvalidUser(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	authService := service.NewAuthService()
	authController := NewAuthController(authService, tokenAuth)

	// Создаем запрос с телом JSON для логина несуществующего пользователя
	user := models.User{Username: "nonexistentuser", Password: "wrongpassword"}
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.LoginHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Проверяем содержимое ответа
	expectedError := "Error encoding token"
	actualError := strings.TrimSpace(rr.Body.String())
	if actualError != expectedError {
		t.Errorf("handler returned unexpected body: got %v want %v", actualError, expectedError)
	}
}
func TestAuthController_RegisterHandler_DuplicateUser(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	authService := service.NewAuthService()
	authController := NewAuthController(authService, tokenAuth)

	// Регистрируем пользователя
	user := models.User{Username: "testuser", Password: "testpassword"}
	authService.RegisterUser(user.Username, user.Password)

	// Создаем запрос с телом JSON для регистрации того же пользователя
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authController.RegisterHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Проверяем содержимое ответа
	expectedError := "username already exists"
	actualError := strings.TrimSpace(rr.Body.String())
	if actualError != expectedError {
		t.Errorf("handler returned unexpected body: got %v want %v", actualError, expectedError)
	}
}
