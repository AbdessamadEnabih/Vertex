package server

import (
	"github.com/AbdessamadEnabih/Vertex/internal/state"
	"github.com/AbdessamadEnabih/Vertex/pkg/server/router"
	"github.com/gin-gonic/gin"
)

type Server struct {
	gin    *gin.Engine
	Router *gin.RouterGroup
	state  *state.State
}

func newServer() *Server {
	engine := gin.Default()

	return &Server{
		gin:    engine,
		Router: engine.Group("/api"),
		state:  state.NewState(),
	}
}

func (s *Server) run(PORT string) {
	s.gin.SetTrustedProxies([]string{"127.0.0.1"})
	s.gin.Run(":" + PORT)
}

func StartServer(PORT string) {
	server := newServer()
	router.InitRoutes(*server.Router)
	server.run(PORT)
}
