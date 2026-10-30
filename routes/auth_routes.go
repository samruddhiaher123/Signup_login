package routes

import (
	"Signup_login/handlers"
	"Signup_login/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(
	router *mux.Router,
	handler *handlers.AuthHandler,
	jwtSecret string,
) {
	router.HandleFunc("/api/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/login", handler.Login).Methods(http.MethodPost)
	router.Handle("/api/profile", middleware.Authenticate(jwtSecret)(http.HandlerFunc(handler.Profile))).Methods(http.MethodGet)
}
