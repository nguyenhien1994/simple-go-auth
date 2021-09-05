package controller

import (
	"simple-go-auth/controller/api"
)

func (s *Server) InitializeRoutes() {
	s.Router.POST("/login", api.Login)

	// Group of APIs that need to authen
	authGroup := s.Router.Group("/")
	authGroup.Use(AuthenHandler())
	{
		authGroup.PUT("/api/user", Authorize("resource", "write", s.Enforcer), api.UpdateUser)
		authGroup.GET("/api/user", Authorize("resource", "read", s.Enforcer), api.GetUser)
		authGroup.POST("/logout", api.Logout)
	}

	s.Router.POST("/refresh", api.Refresh)
}
