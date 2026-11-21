package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware allows the frontend (running on a different port/domain) to call the API.
// Browsers block cross-origin requests unless the server sends these headers.
func CORSMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Max-Age", "600")
		}

		// Preflight request: browser asks permission before the real request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
