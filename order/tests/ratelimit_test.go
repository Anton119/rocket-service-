package tests

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

func TestRateLimit_BurstThenTooManyRequests(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	limit := redis_rate.Limit{Rate: 1, Burst: 3, Period: time.Second}
	handler := ratelimit.HTTPMiddleware(redis_rate.NewLimiter(rdb), limit)(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	for i := range 3 {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code, "запрос %d из burst должен пройти", i+1)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestRateLimit_FailOpenWhenRedisUnavailable(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	t.Cleanup(func() { _ = rdb.Close() })

	limit := redis_rate.Limit{Rate: 1, Burst: 1, Period: time.Second}
	handler := ratelimit.HTTPMiddleware(redis_rate.NewLimiter(rdb), limit)(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
