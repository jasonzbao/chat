package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Clientiface interface {
	Publish(ctx context.Context, message, topic string) error
	ReceiveMessage() <-chan *redis.Message
	Subscribe(ctx context.Context, topic string) error
	Unsubscribe(ctx context.Context, topic string) error
}

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

func (c *Client) Publish(ctx context.Context, message, topic string) error {
	return c.client.Publish(ctx, topic, message).Err()
}

func (c *Client) Subscribe(ctx context.Context, topic string) error {
	return c.pubsub.Subscribe(ctx, topic)
}

func (c *Client) Unsubscribe(ctx context.Context, topic string) error {
	return c.pubsub.Unsubscribe(ctx, topic)
}

func (c *Client) ReceiveMessage() <-chan *redis.Message {
	return c.pubsub.Channel()
}
