// package main
//
// import (
//
//	"fmt"
//	"geo-controller/proxy/internal/controllers"
//	"geo-controller/proxy/internal/service"
//	"net/http"
//
//	"github.com/go-chi/chi"
//	"github.com/go-chi/chi/middleware"
//	"github.com/go-chi/jwtauth"
//	"github.com/rs/cors"
//	httpSwagger "github.com/swaggo/http-swagger"
//
// )
//
// var tokenAuth *jwtauth.JWTAuth
//
//	func init() {
//		// Инициализация JWTAuth с секретным ключом
//		tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
//	}
//
// const daDataApiKey = "c4aab5f0a277fbaa6de6613c4c78930552172d28"
// const daDataSecretKey = "e1e61bbed8ab858bc7153ba44fc8344ba7681526"
//
//	func AuthMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			_, err := jwtauth.VerifyRequest(tokenAuth, r, jwtauth.TokenFromHeader)
//			if err != nil {
//				http.Error(w, "Unauthorized: "+err.Error(), http.StatusForbidden)
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
//
//	func main() {
//		r := chi.NewRouter()
//		r.Use(cors.New(cors.Options{
//			AllowedOrigins:   []string{"http://localhost:1313"},
//			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
//			AllowedHeaders:   []string{"Content-Type", "Authorization"},
//			AllowCredentials: true,
//			MaxAge:           300,
//		}).Handler)
//		r.Use(middleware.Logger)
//
//		authService := service.NewAuthService()
//		authController := controllers.NewAuthController(authService, tokenAuth)
//
//		addressService := service.NewAddressService(daDataApiKey, daDataSecretKey)
//		addressController := controllers.NewAddressController(addressService)
//
//		// Настройка маршрутов
//		r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
//			http.ServeFile(w, r, "./swagger.json")
//		})
//
//		r.Get("/swagger/*", httpSwagger.Handler(
//			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
//		))
//
//		r.Post("/api/register", authController.RegisterHandler)
//		r.Post("/api/login", authController.LoginHandler)
//
//		r.Group(func(r chi.Router) {
//			r.Use(AuthMiddleware)
//			r.Post("/api/address/search", addressController.AddressSearchHandler)
//			r.Post("/api/address/geocode", addressController.GeocodeHandler)
//		})
//
//		fmt.Println("Starting server on port 8080...")
//		http.ListenAndServe(":8080", r)
//	}
//
// main.go
package main

import (
	"fmt"
	"geo-controller/proxy/internal/controllers"
	"geo-controller/proxy/internal/service"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

const daDataApiKey = "c4aab5f0a277fbaa6de6613c4c78930552172d28"
const daDataSecretKey = "e1e61bbed8ab858bc7153ba44fc8344ba7681526"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := jwtauth.VerifyRequest(tokenAuth, r, jwtauth.TokenFromHeader)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1313"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
	r.Use(middleware.Logger)

	authService := service.NewAuthService()
	authController := controllers.NewAuthController(authService, tokenAuth)

	addressService := service.NewAddressService(daDataApiKey, daDataSecretKey)
	addressController := controllers.NewAddressController(addressService)

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger.json")
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Post("/api/register", authController.RegisterHandler)
	r.Post("/api/login", authController.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/api/address/search", addressController.AddressSearchHandler)
		r.Post("/api/address/geocode", addressController.GeocodeHandler)
	})

	return r
}

func main() {
	r := setupRouter()
	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", r)
}
