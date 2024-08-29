package server

import "github.com/gin-gonic/gin"

type Server struct {
	gin    *gin.Engine
	Router *gin.RouterGroup
}

func NewServer() *Server {
	return &Server{
		gin:    gin.Default(),
		Router: gin.New().Group("/api"),
	}
}

func (s *Server) Run() {
	s.gin.SetTrustedProxies([]string{"127.0.0.1"})
	s.gin.Run(":8080")
}
