package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"Signup_login/dto"
	"Signup_login/services"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	// Reject multiple JSON values in one request.
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}

	return nil
}

func (h *AuthHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.RegisterRequest

	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid JSON request"})
		return
	}

	result, err := h.auth.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			writeJSON(w, http.StatusBadRequest,
				map[string]string{"error": "invalid registration details"})
		case errors.Is(err, services.ErrEmailAlreadyExists):
			writeJSON(w, http.StatusConflict,
				map[string]string{"error": "email already registered"})
		default:
			log.Printf("register user: %v", err)
			writeJSON(w, http.StatusInternalServerError,
				map[string]string{"error": "internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.LoginRequest

	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid JSON request"})
		return
	}

	result, err := h.auth.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized,
				map[string]string{"error": "invalid email or password"})
			return
		}

		log.Printf("login user: %v", err)
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Profile(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := r.Context().Value(userIDKey{}).(uint)
	if !ok || userID == 0 {
		writeJSON(w, http.StatusUnauthorized,
			map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.auth.GetProfile(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound,
				map[string]string{"error": "user not found"})
			return
		}

		log.Printf("get profile: %v", err)
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// Used by the authentication middleware to share the authenticated ID.
type userIDKey struct{}

func UserIDContextKey() any {
	return userIDKey{}
}

func ParseBearerToken(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
