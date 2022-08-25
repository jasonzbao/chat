package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/jasonzbao/dyna-take-home/dao/rdb"
	"github.com/jasonzbao/dyna-take-home/dynaerrors"
)

type InstType string

const (
	InstTypeName InstType = "message_name"
	InstTypeExit InstType = "message_exit"
)

var validInstructions = map[string]InstType{
	"/name": InstTypeName,
	"/exit": InstTypeExit,
}

func (s *Server) handleSocketMessage(ctx context.Context, message *WSMessage, conn *V1Connection) (err error) {
	if string(message.Message[0]) == "/" {
		symbols := strings.Split(message.Message, " ")
		inst, ok := validInstructions[symbols[0]]
		if !ok {
			return dynaerrors.ErrorInvalidInstruction
		}
		switch inst {
		case InstTypeName:
			conn.Name = &symbols[1]
			return nil
		case InstTypeExit:
			var msg *rdb.Message
			if msg, err = s.dao.NewMessage("has left the chat!", *conn.Name); err != nil {
				fmt.Println("error sending last message")
			}
			s.redisClient.Publish(ctx, msg.FormatMessage())
			return nil
		}
	}

	if conn.Name == nil {
		return dynaerrors.ErrorNameNotSet
	}

	var msg *rdb.Message
	if msg, err = s.dao.NewMessage(message.Message, *conn.Name); err != nil {
		return err
	}
	if err := s.redisClient.Publish(ctx, msg.FormatMessage()); err != nil {
		fmt.Printf("Had issues publishing to pubsub: %v", err)
	}
	return nil
}
