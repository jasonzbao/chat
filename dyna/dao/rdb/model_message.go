package rdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var colors = map[string]string{
	"Reset":  "\033[0m",
	"Red":    "\033[31m",
	"Green":  "\033[32m",
	"Yellow": "\033[33m",
	"Blue":   "\033[34m",
	"Purple": "\033[35m",
	"Cyan":   "\033[36m",
	"Gray":   "\033[37m",
	"White":  "\033[97m",
}

type Message struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	Contents string `json:"content"`
	Name     string `json:"name"`
	Color    string `json:"color"`
}

func (m *Message) FormatMessage() string {
	var textColor = colors["Reset"]
	if color, ok := colors[m.Color]; ok {
		textColor = color
	}
	return fmt.Sprintf("%s %s%s: %s%s%s", m.CreatedAt, colors["Red"], m.Name, textColor, m.Contents, colors["Reset"])
}

func (c *Client) NewMessage(message, name, color string) (*Message, error) {
	msg := &Message{
		ID:       uuid.New(),
		Name:     name,
		Contents: message,
		Color:    color,
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
