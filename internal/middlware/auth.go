package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hajdudev/invoice-flow/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey = contextKey("userID")

type Middleware struct {
	authenticator auth.Authenticator
}

func New(authenticator auth.Authenticator) *Middleware {
	return &Middleware{
		authenticator: authenticator,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		token, err := m.authenticator.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Could not parse token claims", http.StatusInternalServerError)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "User ID not found in token", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
