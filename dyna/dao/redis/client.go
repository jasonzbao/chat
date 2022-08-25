package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewClient(ctx context.Context, addr string) *Client {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &Client{
		client: r,
		pubsub: r.Subscribe(ctx, "all"), // for now only subscribe to one channel (all)
	}
}

func (c *Client) Publish(ctx context.Context, message string) error {
	return c.client.Publish(ctx, "all", message).Err()
}

func (c *Client) ReceiveMessage() <-chan *redis.Message {
	return c.pubsub.Channel()
}
