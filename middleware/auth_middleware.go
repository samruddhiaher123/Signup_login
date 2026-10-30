package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func UserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(userIDKey).(uint)
	return id, ok && id > 0
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			parts := strings.SplitN(
				r.Header.Get("Authorization"), " ", 2,
			)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") ||
				strings.TrimSpace(parts[1]) == "" {
				writeError(w, http.StatusUnauthorized, "missing or invalid token")
				return
			}

			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(
				strings.TrimSpace(parts[1]),
				claims,
				func(token *jwt.Token) (any, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, jwt.ErrSignatureInvalid
					}
					return []byte(secret), nil
				},
				jwt.WithValidMethods([]string{"HS256"}),
			)

			if err != nil || token == nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			rawID, ok := claims["user_id"].(float64)
			if !ok || rawID < 1 || rawID != float64(uint(rawID)) {
				writeError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			userID := uint(rawID)
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
