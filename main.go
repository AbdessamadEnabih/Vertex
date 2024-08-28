package main

import (
	"net/http"

	"github.com/AbdessamadEnabih/Vertex/pkg/server"
	"github.com/gin-gonic/gin"
)

func main() {
	server := server.NewServer() // Declare the variable here

	server.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Base route")
	})

	server.Router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, "hello")
		c.String(http.StatusOK, "Route accessed successfully")
	})

	server.Run() // listen and serve on 0.0.0.0:8080
}
