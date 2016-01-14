// Package echologrus provides a middleware for echo that logs request details
// via the logrus logging library
package echologrus

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

// New returns a new middleware handler with a default logger
func New() echo.MiddlewareFunc {
	return NewWithLogger(logrus.StandardLogger())
}

// NewWithLogger returns a new middleware handler with the specified logger
func NewWithLogger(l *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			isError := false

			if err := next(c); err != nil {
				c.Error(err)
				isError = true
			}

			latency := time.Since(start)

			entry := l.WithFields(logrus.Fields{
				"request":     c.Request().RequestURI,
				"method":      c.Request().Method,
				"remote":      c.Request().RemoteAddr,
				"status":      c.Response().Status(),
				"text_status": http.StatusText(c.Response().Status()),
				"took":        latency,
			})

			if reqID := c.Request().Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}
			// Check middleware error
			if isError {
				entry.Error("error by handling request")
			} else {
				entry.Info("request has been successfully processed")
			}

			return nil
		}
	}
}
