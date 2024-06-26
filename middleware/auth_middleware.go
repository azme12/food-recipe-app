// auth_middleware.go

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware is a middleware function to authenticate requests
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("azme07"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID := token.Claims.(jwt.MapClaims)["userID"].(string)
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		// Authentication successful, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
