package server

import (
	"context"
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

	chName := string(uuid.New().String()[0:5])
	conn := &V1Connection{
		Name: &chName,
	}

	var msg *rdb.Message
	if msg, err = s.dao.NewMessage("has joined the chat!", *conn.Name); err != nil {
		response.Error = pkgErrors.Wrap(err, "Error sending first message")
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	s.redisClient.Publish(c, msg.FormatMessage())

	defer func() {
		var msg *rdb.Message
		if msg, err = s.dao.NewMessage("has left the chat!", *conn.Name); err != nil {
			fmt.Println("error sending last message")
		}
		s.redisClient.Publish(c, msg.FormatMessage())
	}()

	innerCtx, cancel := context.WithCancel(c)
	defer cancel()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// buffer would be one thing we would need to tune
	// drop messages if rate is too fast
	subCh := make(chan string, 10)
	s.messages[chName] = subCh
	defer func() {
		delete(s.messages, chName)
	}()

	// specific channel for closing messages
	cm := make(chan bool, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-innerCtx.Done():
				return
			case msg := <-subCh:
				err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					fmt.Println(err)
					cancel()
					return
				}
			case <-cm:
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "see ya")
				if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	for {
		//Read Message from client
		select {
		case <-innerCtx.Done():
			return
		default:
			mt, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				cancel()
				break
			}

			var wsMessage WSMessage
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				fmt.Println("Error decoding message")
				cancel()
				break
			}

			if err = s.handleSocketMessage(c, &wsMessage, conn); err != nil {
				if errors.Is(err, dynaerrors.ErrorNameNotSet) {
					ws.WriteMessage(mt, []byte("Name needs to be set before messages can be sent"))
				} else if errors.Is(err, dynaerrors.ErrorInvalidInstruction) {
					ws.WriteMessage(mt, []byte("Unrecognized instruction"))
				} else if errors.Is(err, dynaerrors.ErrorExitChat) {
					cm <- true
					cancel()
					break
				} else {
					fmt.Println("invalid message %v", err)
				}
				continue
			}
		}
	}

}
