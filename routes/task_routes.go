package routes

import (
	"Signup_login/handlers"
	"Signup_login/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func TaskRoutes(r *mux.Router, handler *handlers.TaskHandler, jwtSecret string) {
	authenticated := middleware.Authenticate(jwtSecret)
	r.Handle("/api/tasks", authenticated(http.HandlerFunc(handler.CreateTask))).
		Methods(http.MethodPost)
	r.Handle("/api/tasks", authenticated(http.HandlerFunc(handler.ListTasks))).
		Methods(http.MethodGet)
	r.Handle("/api/tasks/{id}", authenticated(http.HandlerFunc(handler.GetTask))).
		Methods(http.MethodGet)
	r.Handle("/api/tasks/{id}", authenticated(http.HandlerFunc(handler.UpdateTask))).
		Methods(http.MethodPut)
	r.Handle("/api/tasks/{id}", authenticated(http.HandlerFunc(handler.DeleteTask))).
		Methods(http.MethodDelete)
}
