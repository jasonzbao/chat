package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	pkgErrors "github.com/pkg/errors"

	"github.com/jasonzbao/dyna-take-home/dao/rdb"
	"github.com/jasonzbao/dyna-take-home/dynaerrors"
)

var upgrader = websocket.Upgrader{}

type V1SocketResponse struct {
	Error error `json:"error,omitempty"`
}

type WSMessage struct {
	Message string `json:"message" binding:"required"`
}

type V1Connection struct {
	Name *string `json:"name"`
}

func (s *Server) handleSocket(c *gin.Context) {
	response := &V1SocketResponse{}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.Error = pkgErrors.Wrap(err, "Error upgrading socket request")
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer ws.Close()

	newName := string(uuid.New().String()[0:5])
	conn := &V1Connection{
		Name: &newName,
	}

	var msg *rdb.Message
	if msg, err = s.dao.NewMessage("has joined the chat!", *conn.Name); err != nil {
		response.Error = pkgErrors.Wrap(err, "Error sending first message")
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	s.redisClient.Publish(c, msg.FormatMessage(*conn.Name))

	defer func() {
		var msg *rdb.Message
		if msg, err = s.dao.NewMessage("has left the chat!", *conn.Name); err != nil {
			fmt.Println("error sending last message")
		}
		s.redisClient.Publish(c, msg.FormatMessage(*conn.Name))
	}()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// buffer would be one thing we would need to tune
	// drop messages if rate is too fast
	subCh := make(chan string, 10)
	s.messages = append(s.messages, subCh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-c.Done():
				return
			case msg := <-subCh:
				err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	}()

	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		var wsMessage WSMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			fmt.Println("Error decoding message")
			break
		}

		if err = s.handleSocketMessage(c, &wsMessage, conn); err != nil {
			if errors.Is(err, dynaerrors.ErrorNameNotSet) {
				ws.WriteMessage(mt, []byte("Name needs to be set before messages can be sent"))
			} else if errors.Is(err, dynaerrors.ErrorInvalidInstruction) {
				ws.WriteMessage(mt, []byte("Unrecognized instruction"))
			} else {
				fmt.Println("invalid message %v", err)
			}
			continue
		}
	}

}
