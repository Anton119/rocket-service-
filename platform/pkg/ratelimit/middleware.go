package ratelimit

import (
	"log/slog"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

// HTTPMiddleware создаёт HTTP middleware с распределённым rate limiter (GCRA в Redis).
//
// Ключ лимита — путь запроса (r.URL.Path). Все инстансы сервиса делят один счётчик в Redis.
// Если Redis недоступен, запрос пропускается (fail-open): лучше временно пропустить
// лишний трафик, чем остановить сервис из-за проблем с Redis.
func HTTPMiddleware(limiter *redis_rate.Limiter, limit redis_rate.Limit) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res, err := limiter.Allow(r.Context(), r.URL.Path, limit)
			if err != nil {
				slog.Error("ошибка проверки rate limit", "path", r.URL.Path, "error", err)
				next.ServeHTTP(w, r)
				return
			}

			if res.Allowed == 0 {
				slog.Warn("запрос отклонён rate limiter",
					"path", r.URL.Path,
					"retry_after", res.RetryAfter,
				)
				http.Error(w, "слишком много запросов, попробуйте позже", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
