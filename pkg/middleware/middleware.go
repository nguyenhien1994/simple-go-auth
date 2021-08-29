package middleware

import (
	"github.com/gin-gonic/gin"
	"simple-go-auth/pkg/auth"
	"net/http"
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
		if err := m.token.TokenValid(c.Request); err != nil  {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

