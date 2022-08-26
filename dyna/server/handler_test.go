package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jasonzbao/dyna-take-home/config"
	"github.com/jasonzbao/dyna-take-home/dao/rdb"
	"github.com/jasonzbao/dyna-take-home/dao/redis"
	"github.com/jasonzbao/dyna-take-home/dynaerrors"
)

type mockRedisClient struct {
	redis.Clientiface

	messages []string
}

func (m *mockRedisClient) Publish(ctx context.Context, message, topic string) error {
	m.messages = append(m.messages, message)
	return nil
}

func TestHandleSocketMessage(t *testing.T) {
	cfg := &config.Config{
		DBConnection: "host=localhost user=dyna password=board dbname=dyna port=5431 sslmode=disable TimeZone=America/Los_Angeles",
	}

	dao, err := rdb.NewClient(cfg.DBConnection)
	require.Nil(t, err)

	rClient := &mockRedisClient{messages: []string{}}

	s := NewServer(cfg, dao, rClient)

	v1Name := "jason"
	conn := &V1Connection{
		Name: &v1Name,
	}

	err = s.handleSocketMessage(context.TODO(), &WSMessage{
		Message: "/invalid",
	}, conn)
	require.True(t, errors.Is(err, dynaerrors.ErrorInvalidInstruction))

	err = s.handleSocketMessage(context.TODO(), &WSMessage{
		Message: "/name",
	}, conn)
	require.True(t, errors.Is(err, dynaerrors.ErrorInvalidInstruction))

	err = s.handleSocketMessage(context.TODO(), &WSMessage{
		Message: "/name jose",
	}, conn)
	require.Nil(t, err)
	require.Equal(t, *conn.Name, "jose")

	err = s.handleSocketMessage(context.TODO(), &WSMessage{
		Message: "/exit",
	}, conn)
	require.True(t, errors.Is(err, dynaerrors.ErrorExitChat))

	require.Equal(t, len(rClient.messages), 0)
	err = s.handleSocketMessage(context.TODO(), &WSMessage{
		Message: "Hello World",
	}, conn)
	require.Nil(t, err)
	require.Equal(t, len(rClient.messages), 1)
}
