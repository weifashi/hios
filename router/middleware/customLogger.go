package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		statusCode := c.Writer.Status()
		logLevel := "INFO"

		switch {
		case statusCode >= 500:
			logLevel = "ERROR"
		case statusCode >= 400:
			logLevel = "WARN"
		case statusCode >= 300:
			logLevel = "INFO"
		}

		if strings.HasPrefix(path, "/api/v1") && !strings.HasPrefix(path, "/src") && !strings.Contains(path, "/node_modules") && !strings.Contains(path, "/@") {
			log.Printf("[GIN] [%s] %s %s %s %d\n", logLevel, c.Request.Method, path, query, statusCode)
		}
	}
}
