package controllers

import (
	"encoding/json"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/responder"
	"geo-controller/proxy/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"net/http"
	"strconv"
)

type AuthController struct {
	authService *service.AuthService
	tokenAuth   *jwtauth.JWTAuth
	responder   *responder.Responder
}

func NewAuthController(authService *service.AuthService, tokenAuth *jwtauth.JWTAuth) *AuthController {
	return &AuthController{
		authService: authService,
		tokenAuth:   tokenAuth,
		responder:   responder.NewResponder(),
	}
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.authService.RegisterUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !c.authService.AuthenticateUser(user.Username, user.Password) {
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}

	//Создание JWT токена при успешной аутентификации
	_, tokenString, err := c.tokenAuth.Encode(map[string]interface{}{"username": user.Username})
	if err != nil {
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
func (c *AuthController) ListHandler(w http.ResponseWriter, r *http.Request) {
	users, err := c.authService.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	c.responder.OutputJSON(w, users)
}
func (c *AuthController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
	user, err := c.authService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	c.responder.OutputJSON(w, user)
}
func (c *AuthController) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = id
	err = c.authService.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated"))
}
func (c *AuthController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = c.authService.DeleteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted"))
}
