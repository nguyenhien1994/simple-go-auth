package middleware

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-go-auth/pkg/auth"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.GetTokenService().TokenValid(c.Request); err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := auth.GetTokenService().ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		// casbin enforces policy
		ok := enforcer.Enforce(metadata.Username, obj, act)
		if !ok {
			c.AbortWithStatusJSON(403, "forbidden")
			return
		}
		c.Next()
	}
}
