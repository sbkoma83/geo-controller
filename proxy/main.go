// main.go
package main

import (
	"database/sql"
	"fmt"
	"geo-controller/proxy/internal/controllers"
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
)

var tokenAuth *jwtauth.JWTAuth
var db *sql.DB

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

// func setupRouter() *chi.Mux {

// 	return r

// }

func main() {
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

	// Формируем строку подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")

	// Создаем таблицу users
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL
        )
    `)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	log.Println("Users table created or already exists")
}
