package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type ChannelService struct {
	db *sqlx.DB
	id int
}

func NewChannel(
	db *sqlx.DB,
	id int,
) *ChannelService {
	return &ChannelService{
		db: db,
		id: id,
	}
}

type Message struct {
	ID        int        `db:"message_id"`
	Copy      string     `db:"message_copy"`
	Username  string     `db:"user_name"`
	CreatedAt *time.Time `db:"message_time"`
}
type Messages []Message

func (c *ChannelService) Messages() (Messages, error) {
	messages := Messages{}

	err := c.db.Select(&messages, "SELECT message_id, message_copy, user_name, message_time FROM channel_messages WHERE channel_id = ?", c.id)
	if err != nil {
		slog.Error("could not load channels", slog.String("error", err.Error()))

		return nil, fmt.Errorf("could not find channels for user: %w", err)
	}

	return messages, nil
}

func (c *ChannelService) ID() int {
	return c.id
}
