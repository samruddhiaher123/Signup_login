package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"Signup_login/config"
	"Signup_login/handlers"
	"Signup_login/model"
	"Signup_login/repository"
	"Signup_login/routes"
	"Signup_login/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must contain at least 32 characters")
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	// For a new local demo database only. For your existing
	// Study Planner, prefer your versioned migrations.
	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := db.AutoMigrate(&model.User{}); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router := mux.NewRouter()
	routes.RegisterAuthRoutes(router, authHandler, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}
