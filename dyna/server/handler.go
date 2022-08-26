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
	InstTypeName    InstType = "message_name"
	InstTypeExit    InstType = "message_exit"
	InstTypeColor   InstType = "message_color"
	InstTypePrivate InstType = "message_private"
)

var validInstructions = map[string]InstType{
	"/name":    InstTypeName,
	"/exit":    InstTypeExit,
	"/color":   InstTypeColor,
	"/private": InstTypePrivate,
}

const privateTxt = "\033[36m PRIVATE MESSAGE FROM %s:\033[0m"

func (s *Server) handleSocketMessage(ctx context.Context, message *WSMessage, conn *V1Connection) (err error) {
	if string(message.Message[0]) == "/" {
		symbols := strings.Split(message.Message, " ")
		inst, ok := validInstructions[symbols[0]]
		if !ok {
			return dynaerrors.ErrorInvalidInstruction
		}
		switch inst {
		case InstTypeName:
			if len(symbols) != 2 {
				return dynaerrors.ErrorInvalidInstruction
			}
			conn.Name = &symbols[1]
			return nil
		case InstTypeExit:
			if len(symbols) != 1 {
				return dynaerrors.ErrorInvalidInstruction
			}
			return dynaerrors.ErrorExitChat
		case InstTypeColor:
			if len(symbols) != 2 {
				return dynaerrors.ErrorInvalidInstruction
			}
			conn.Color = symbols[1]
			return nil
		case InstTypePrivate:
			if len(symbols) <= 2 {
				return dynaerrors.ErrorInvalidInstruction
			}
			chanName := symbols[1]
			txt := fmt.Sprintf(privateTxt, *conn.Name) + strings.Join(symbols[2:], " ")
			fmt.Println("private txt", txt)
			if err := s.redisClient.Publish(ctx, txt, chanName); err != nil {
				fmt.Printf("Had issues publishing to pubsub: %v", err)
			}
			return nil
		}
	}

	if conn.Name == nil {
		return dynaerrors.ErrorNameNotSet
	}

	var msg *rdb.Message
	if msg, err = s.dao.NewMessage(message.Message, *conn.Name, conn.Color, conn.ChName); err != nil {
		return err
	}
	if err := s.redisClient.Publish(ctx, msg.FormatMessage(), "all"); err != nil {
		fmt.Printf("Had issues publishing to pubsub: %v", err)
	}
	return nil
}
