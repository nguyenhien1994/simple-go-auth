package controller

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-go-auth/services/auth"
)

func AuthenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessDetails, err := auth.GetTokenService().AccessDetailsFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		c.Set(auth.ContextAccessDetailsKey, accessDetails)
		c.Next()
	}
}

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessDetails := c.MustGet(auth.ContextAccessDetailsKey).(*auth.AccessDetails)
		// casbin enforces policy
		ok := enforcer.Enforce(accessDetails.Username, obj, act)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, "forbidden")
			return
		}
		c.Next()
	}
}
