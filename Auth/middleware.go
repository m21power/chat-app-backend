package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)
type JWTConfig struct {
	SECRET_KEY string
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		SECRET_KEY: getEnv("SECRET_KEY",""),
	}
}
type contextKey string

const (
	ContextUserID  contextKey = "userID"
	ContextUserRole contextKey = "userRole"
)

func RoleMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			_, claims, err := VerifyToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			userID := claims.Subject
			roleSlice := claims.Audience
			if len(roleSlice) == 0 {
				http.Error(w, "Missing role in token", http.StatusForbidden)
				return
			}
			userRole := roleSlice[0]

			for _, role := range requiredRoles {
				if userRole == role {
					ctx := context.WithValue(r.Context(), ContextUserID, userID)
					ctx = context.WithValue(ctx, ContextUserRole, userRole)

					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			http.Error(w, "Forbidden - Insufficient role", http.StatusForbidden)
		})
	}
}


func getEnv(key, defaultValue string) string {
	err :=godotenv.Load()
	if err != nil{
		log.Fatal(err)
	}
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
