package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, _ := c.Get("userData")
		userClaim, ok := userData.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		userRole, ok := userClaim["role"].(string)

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "User role not found in token"})
			return
		}

		if userRole != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You don't have permission to perform this action"})
			return
		}

		c.Next()
	}
}
