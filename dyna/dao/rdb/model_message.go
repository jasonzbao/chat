package rdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

type Message struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	Contents string `json:"content"`
	Name     string `json:"name"`
}

func (m *Message) FormatMessage(name string) string {
	return fmt.Sprintf("%s %s%s%s: %s", m.CreatedAt, Red, name, Reset, m.Contents)
}

func (c *Client) NewMessage(message, name string) (*Message, error) {
	msg := &Message{
		ID:       uuid.New(),
		Name:     name,
		Contents: message,
	}

	if err := c.db.Create(msg).Error; err != nil {
		return nil, errors.Wrap(err, "error talking to db")
	}
	return msg, nil
}

// retrieves the last 10 messages
func (c *Client) RetrieveLastMessages() ([]*Message, error) {
	var messages []*Message
	if err := c.db.Limit(10).Order("created_at desc").Find(&messages).Error; err != nil {
		return nil, errors.Wrap(err, "error talking to db")
	}
	return messages, nil
}
