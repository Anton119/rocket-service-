package middleware

import (
	"context"
	"net/http"
	"strings"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/platform/pkg/auth"
)

// SessionValidator проверяет сессию через IAM.
type SessionValidator interface {
	Whoami(ctx context.Context, sessionUUID string) (userUUID string, err error)
}

// AuthMiddleware проверяет Bearer/X-Session-UUID и кладёт user/session UUID в context.
func AuthMiddleware(validator SessionValidator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID, err := sessionUUIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userUUID, err := validator.Whoami(r.Context(), sessionUUID)
		if err != nil {
			http.Error(w, "сессия недействительна или истекла", http.StatusUnauthorized)
			return
		}

		ctx := auth.WithSessionUUID(r.Context(), sessionUUID)
		ctx = auth.WithUserUUID(ctx, userUUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func sessionUUIDFromRequest(r *http.Request) (string, error) {
	if sessionUUID := strings.TrimSpace(r.Header.Get(auth.SessionHeader)); sessionUUID != "" {
		return sessionUUID, nil
	}

	const bearerPrefix = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, bearerPrefix) {
		sessionUUID := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	return "", errs.ErrUnauthorized
}
