package controllers


import (
	"encoding/json"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/service"
	"github.com/go-chi/jwtauth"
	"net/http"
)

type AuthController struct {
	authService *service.AuthService
	tokenAuth   *jwtauth.JWTAuth
}

func NewAuthController(authService *service.AuthService, tokenAuth *jwtauth.JWTAuth) *AuthController {
	return &AuthController{
		authService: authService,
		tokenAuth:   tokenAuth,
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

