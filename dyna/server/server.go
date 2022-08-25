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
	redisClient *redis.Client

	messages map[string]chan string
}

func NewServer(cfg *config.Config, dao *rdb.Client, redisClient *redis.Client) *Server {
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
				for _, ch := range s.messages {
					ch <- msg.Payload
				}
			}
		}
	}()

	router := s.NewRouter()
	return router.Run(port)
}
