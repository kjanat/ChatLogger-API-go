package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kjanat/ChatLogger-API-go/internal/version"
)

// VersionHeader adds the API version to response headers.
func VersionHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", version.Version)
		c.Next()
	}
}
