package middleware

import (
	"log"
	"fmt"

	"github.com/gin-gonic/gin"
	"net/http"
	"simple-go-auth/pkg/auth"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
)

// handler struct
type middleware struct {
	token auth.TokenInterface
}

func NewMiddleWare(token auth.TokenInterface) *middleware {
	return &middleware{token}
}

func (m *middleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := m.token.TokenValid(c.Request); err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Authorize determines if current subject has been authorized to take an action on an object.
func (m *middleware) Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := m.token.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}
		metadata, err := m.token.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		// casbin enforces policy
		ok, err := enforce(metadata.Username, obj, act, adapter)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, "error occurred when authorizing user")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "forbidden")
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	enforcer := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	ok := enforcer.Enforce(sub, obj, act)
	return ok, nil
}
