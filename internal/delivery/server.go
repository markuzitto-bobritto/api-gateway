package delivery

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Serverer interface {
	Run()
	Stop()
	MustImplementServer()
}

type Server struct {
	srv *http.Server
}

func InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/", GetUser)
	return router
}

func NewServer(addr string) *Server {
	routes := InitRoutes()
	s := &http.Server{
		Addr:    addr,
		Handler: routes,
	}
	server := &Server{
		srv: s,
	}
	return server
}

func (s *Server) Run() {
	s.srv.ListenAndServe()
}

func (s *Server) Stop() {
	slog.Info("Server shutted down")
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func (s *Server) MustImplementServer() {}
