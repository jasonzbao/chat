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
	}
	return router
}
