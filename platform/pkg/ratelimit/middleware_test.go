package ratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Anton119/rocket-service-/platform/pkg/ratelimit"
)

func newTestHandler(limiter *redis_rate.Limiter, limit redis_rate.Limit) http.Handler {
	return ratelimit.HTTPMiddleware(limiter, limit)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func TestHTTPMiddleware_ExceedBurstReturns429(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	limit := redis_rate.Limit{Rate: 1, Burst: 2, Period: time.Second}
	handler := newTestHandler(redis_rate.NewLimiter(rdb), limit)

	for i := range 2 {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "запрос %d из burst должен пройти", i+1)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestHTTPMiddleware_FailOpenWhenRedisDown(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	t.Cleanup(func() { _ = rdb.Close() })

	limit := redis_rate.Limit{Rate: 1, Burst: 1, Period: time.Second}
	handler := newTestHandler(redis_rate.NewLimiter(rdb), limit)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "при недоступном Redis запрос должен пройти")
}
