package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) NewRouter() *gin.Engine {
	router := gin.Default()
	{
		router.Any("/ping", func(c *gin.Context) {
			c.JSON(200, nil)
		})

		router.GET("/socket", func(c *gin.Context) {
			s.handleSocket(c)
		})

		router.GET("/history", func(c *gin.Context) {
			s.handleHistory(c)
		})
	}
	return router
}
