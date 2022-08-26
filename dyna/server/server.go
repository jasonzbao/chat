package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/jasonzbao/dyna-take-home/config"
	"github.com/jasonzbao/dyna-take-home/dao/rdb"
	"github.com/jasonzbao/dyna-take-home/dao/redis"
)

type Server struct {
	config      *config.Config
	dao         *rdb.Client
	redisClient redis.Clientiface

	messages map[string]chan string
}

func NewServer(cfg *config.Config, dao *rdb.Client, redisClient redis.Clientiface) *Server {
	return &Server{
		config:      cfg,
		dao:         dao,
		redisClient: redisClient,

		messages: map[string]chan string{},
	}
}

// blocking
func (s *Server) Run(ctx context.Context, port string) error {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-s.redisClient.ReceiveMessage():
				fmt.Println("Got message from redis", msg.Payload)
				if msg.Channel == "all" {
					for _, ch := range s.messages {
						ch <- msg.Payload
					}
				} else {
					// private channel case. For now just support DMs
					if ch, ok := s.messages[msg.Channel]; ok {
						ch <- msg.Payload
					} else {
						fmt.Printf("Got a message to a channel that isn't on this server %s", msg.Channel)
					}
				}
			}
		}
	}()

	router := s.NewRouter()
	return router.Run(port)
}
