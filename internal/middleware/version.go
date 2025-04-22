package middleware

import (
	"ChatLogger-API-go/internal/version"
	"github.com/gin-gonic/gin"
)

// VersionHeader adds the API version to response headers.
func VersionHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", version.Version)
		c.Next()
	}
}
