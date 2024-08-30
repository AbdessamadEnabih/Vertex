package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	gin    *gin.Engine
	Router *gin.RouterGroup
}

func newServer() *Server {
	engine := gin.Default()

	return &Server{
		gin:    engine,
		Router: engine.Group("/api"),
	}
}

func (s *Server) run(PORT string) {
	s.gin.SetTrustedProxies([]string{"127.0.0.1"})
	s.gin.Run(":" + PORT)
}

func initRoutes(server *Server) {
	server.Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Base route"})
	})

	server.Router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello"})
	})
}

func StartServer(PORT string) {
	server := newServer()
	initRoutes(server)
	server.run(PORT)
}
