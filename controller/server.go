package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"simple-go-auth/services/auth"
)

// server struct
type Server struct {
	Router   *gin.Engine
	Enforcer *casbin.Enforcer
	Auth     auth.AuthInterface
	Token    auth.TokenInterface
}

func (s *Server) Initialize() {
	s.Auth = auth.GetAuthService()
	s.Token = auth.GetTokenService()
	s.Enforcer = casbin.NewEnforcer("config/rbac_model.conf", "config/policy.csv")
	s.Router = gin.Default()
	s.InitializeRoutes()
}

func (s *Server) Run(addr string) {
	fmt.Printf("Listen on port %s \n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Router))

	srv := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	log.Println("Shuting down the server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Failed to shutdown the server:", err)
	}
	log.Println("Finished shutdown")
}

