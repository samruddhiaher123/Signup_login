package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"Signup_login/config"
	"Signup_login/handlers"
	"Signup_login/repository"
	"Signup_login/routes"
	"Signup_login/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using environment variables")
	}

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must contain at least 32 characters")
	}

	// Connect database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	// ---------------- USER AUTHENTICATION ----------------

	userRepo := repository.NewUserRepository(db)

	authService := services.NewAuthService(
		userRepo,
		jwtSecret,
	)

	authHandler := handlers.NewAuthHandler(authService)

	// ---------------- TASK API ----------------

	taskRepo := repository.NewTaskRepository(db)

	taskService := services.NewTaskService(taskRepo)

	taskHandler := handlers.NewTaskHandler(taskService)

	// ---------------- ROUTES ----------------

	router := mux.NewRouter()

	// Authentication routes
	routes.RegisterAuthRoutes(router, authHandler, jwtSecret)

	// Task routes
	routes.TaskRoutes(router, taskHandler, jwtSecret)

	// ---------------- SERVER ----------------

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
