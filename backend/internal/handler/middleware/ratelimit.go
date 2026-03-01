package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"golang.org/x/time/rate"
)

type RateLimiterConfig struct {
	Rate  rate.Limit
	Burst int
}

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	limiters map[string]*ipLimiter
	config   RateLimiterConfig
}

func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*ipLimiter),
		config:   cfg,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	if l, exists := rl.limiters[ip]; exists {
		l.lastSeen = time.Now()
		return l.limiter
	}
	l := rate.NewLimiter(rl.config.Rate, rl.config.Burst)
	rl.limiters[ip] = &ipLimiter{limiter: l, lastSeen: time.Now()}
	return l
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		for ip, l := range rl.limiters {
			if time.Since(l.lastSeen) > 10*time.Minute {
				delete(rl.limiters, ip)
			}
		}
	}
}

func RateLimit(requestsPerMinute int) echo.MiddlewareFunc {
	rl := NewRateLimiter(RateLimiterConfig{
		Rate:  rate.Limit(float64(requestsPerMinute) / 60.0),
		Burst: requestsPerMinute,
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := rl.getLimiter(ip)

			if !limiter.Allow() {
				resetTime := time.Now().Add(time.Minute).Unix()
				c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")
				c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

				return c.JSON(http.StatusTooManyRequests, model.NewValidationErrorResponse(
					"RATE_LIMITED",
					"リクエスト数の制限を超えました。しばらく待ってから再試行してください",
					[]model.FieldError{
						{Field: "retry_after", Message: fmt.Sprintf("%d秒後に再試行してください", 30)},
					},
				))
			}

			remaining := int(limiter.Tokens())
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

			return next(c)
		}
	}
}
