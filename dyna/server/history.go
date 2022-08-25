package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonzbao/dyna-take-home/dao/rdb"
)

type V1HistoryResponse struct {
	Messages []*rdb.Message `json:"past_messages"`
}

func (s *Server) handleHistory(c *gin.Context) {
	resp := &V1HistoryResponse{}

	messages, err := s.dao.RetrieveLastMessages()
	if err != nil {
		fmt.Println("Error: %v", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	resp.Messages = messages
	c.JSON(http.StatusOK, resp)
}
