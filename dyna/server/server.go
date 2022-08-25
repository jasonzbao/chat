package server

import (
	"github.com/jasonzbao/dyna-take-home/config"
	"github.com/jasonzbao/dyna-take-home/dao/rdb"
)

type Server struct {
	config *config.Config
	dao    *rdb.Client
}

func NewServer(cfg *config.Config, dao *rdb.Client) *Server {
	return &Server{
		config: cfg,
		dao:    dao,
	}
}

// blocking
func (s *Server) Run(port string) error {
	router := s.NewRouter()
	return router.Run(port)
}
