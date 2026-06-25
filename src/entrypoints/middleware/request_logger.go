package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/observability"
)

const RequestIDKey = "request_id"

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		// Propagamos el request_id al context.Context para que los logs de los usecases
		// (vía slog.*Context) queden correlacionados por el ContextHandler.
		c.Request = c.Request.WithContext(observability.WithRequestID(c.Request.Context(), requestID))

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)

		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Int64("duration_ms", duration.Milliseconds()),
		}

		if orgID, exists := c.Get(OrgIDKey); exists {
			if id, ok := orgID.(uuid.UUID); ok {
				attrs = append(attrs, slog.String("org_id", id.String()))
			}
		}
		if userID, exists := c.Get(UserIDKey); exists {
			if id, ok := userID.(int64); ok {
				attrs = append(attrs, slog.Int64("user_id", id))
			}
		}

		args := make([]any, len(attrs))
		for i, a := range attrs {
			args[i] = a
		}

		switch {
		case status >= 500:
			slog.Error("request completed", args...)
		case status >= 400:
			slog.Warn("request completed", args...)
		default:
			slog.Info("request completed", args...)
		}
	}
}
