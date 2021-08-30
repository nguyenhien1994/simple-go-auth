package server

import (
	"simple-go-auth/pkg/api"
	"simple-go-auth/pkg/middleware"
)

func (s *Server) InitializeRoutes() {
	s.Router.POST("/login", api.Login)
	authorized := s.Router.Group("/")
	authorized.Use(middleware.TokenAuthMiddleware())
	{
		authorized.POST("/api/todo", middleware.Authorize("resource", "write", s.FileAdapter), api.CreateTodo)
		authorized.GET("/api/todo", middleware.Authorize("resource", "read", s.FileAdapter), api.GetTodo)
		authorized.POST("/logout", api.Logout)
	}
	s.Router.POST("/refresh", api.Refresh)
}

