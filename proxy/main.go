// main.go
package main

import (
	"fmt"
	"geo-controller/proxy/internal/controllers"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/repositories"
	"geo-controller/proxy/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth
var db *gorm.DB

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	initDB()
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
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
	r.Use(middleware.Logger)
	userRepo := repositories.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
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
	r.Get("/api/users/{id}", authController.GetByID)
	r.Get("/api/users/list", authController.ListHandler)
	r.Post("/api/users/update/{id}", authController.UpdateByID)
	r.Delete("/api/users/delete/{id}", authController.DeleteByID)

	return r
}

func main() {
	r := setupRouter()
	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", r)

}

func initDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate models")
	}

}
