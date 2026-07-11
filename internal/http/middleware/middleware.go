package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	logger, _ := zap.NewProduction()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		requestID := c.GetString("request_id")

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetString("request_id")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":      "internal server error",
			"request_id": requestID,
		})
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Idempotency-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	// Simple in-memory rate limiter
	// In production, use Redis-based rate limiter
	type clientLimiter struct {
		tokens     float64
		lastRefill time.Time
	}

	limiters := make(map[string]*clientLimiter)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		limiter, exists := limiters[clientIP]
		if !exists {
			limiter = &clientLimiter{
				tokens:     float64(burst),
				lastRefill: now,
			}
			limiters[clientIP] = limiter
		}

		// Refill tokens
		elapsed := now.Sub(limiter.lastRefill).Seconds()
		limiter.tokens += elapsed * requestsPerSecond
		if limiter.tokens > float64(burst) {
			limiter.tokens = float64(burst)
		}
		limiter.lastRefill = now

		if limiter.tokens < 1 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		limiter.tokens--
		c.Next()
	}
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		// Authentication is intentionally minimal for Sprint 0-2; replace with JWT validation in hardening.
		c.Set("user_id", uuid.Nil)
		c.Next()
	}
}
