package handlers

import (
	"Signup_login/dto"
	"Signup_login/middleware"
	"Signup_login/services"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.TaskRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid JSON request"})
		return
	}

	userID := taskUserID(w, r)
	if userID == 0 {
		return
	}

	task, err := h.Service.CreateTask(userID, req)
	if err != nil {
		writeTaskError(w, "create task", err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID := taskUserID(w, r)
	if userID == 0 {
		return
	}

	tasks, err := h.Service.ListTasks(userID)
	if err != nil {
		writeTaskError(w, "list tasks", err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := taskUserID(w, r)
	if userID == 0 {
		return
	}
	taskID, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	task, err := h.Service.GetTask(userID, taskID)
	if err != nil {
		writeTaskError(w, "get task", err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := taskUserID(w, r)
	if userID == 0 {
		return
	}
	taskID, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	var req dto.TaskRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid JSON request"})
		return
	}

	task, err := h.Service.UpdateTask(userID, taskID, req)
	if err != nil {
		writeTaskError(w, "update task", err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := taskUserID(w, r)
	if userID == 0 {
		return
	}
	taskID, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	if err := h.Service.DeleteTask(userID, taskID); err != nil {
		writeTaskError(w, "delete task", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func taskUserID(w http.ResponseWriter, r *http.Request) uint {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized,
			map[string]string{"error": "unauthorized"})
		return 0
	}
	return userID
}

func parseTaskID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	value := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(value, 10, strconv.IntSize)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "task id must be a positive integer"})
		return 0, false
	}
	return uint(id), true
}

func writeTaskError(w http.ResponseWriter, operation string, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidTaskInput):
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
	case errors.Is(err, services.ErrTaskNotFound):
		writeJSON(w, http.StatusNotFound,
			map[string]string{"error": "task not found"})
	default:
		log.Printf("%s: %v", operation, err)
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "could not " + operation})
	}
}
