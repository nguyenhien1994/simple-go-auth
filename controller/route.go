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
		authGroup.POST("/api/todo", Authorize("resource", "write", s.Enforcer), api.CreateTodo)
		authGroup.GET("/api/todo", Authorize("resource", "read", s.Enforcer), api.GetTodo)
		authGroup.POST("/logout", api.Logout)
	}

	s.Router.POST("/refresh", api.Refresh)
}
