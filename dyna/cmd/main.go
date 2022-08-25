package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jasonzbao/dyna-take-home/config"
	"github.com/jasonzbao/dyna-take-home/dao/rdb"
	"github.com/jasonzbao/dyna-take-home/dao/redis"
	"github.com/jasonzbao/dyna-take-home/server"
)

var (
	configFile = flag.String(
		"config",
		"./config.json",
		"config file")
)

func main() {
	flag.Parse()

	fmt.Println(*configFile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("=> proc signal:", sig.String())
			os.Exit(0)
		}
	}()

	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		panic(err)
	}

	dao, err := rdb.NewClient(cfg.DBConnection)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewClient(ctx, cfg.RedisAddr)

	server := server.NewServer(cfg, dao, redisClient)
	err = server.Run(ctx, cfg.Port)
	if err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
