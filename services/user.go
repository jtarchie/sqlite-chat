package services

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	db    *sqlx.DB
	email string
}

func NewUser(
	db *sqlx.DB,
	email string,
) *UserService {
	return &UserService{
		db:    db,
		email: email,
	}
}

type Channel struct {
	Description string `db:"description"`
	Name        string `db:"channel_name"`
	ID          int    `db:"channel_id"`
	Private     bool   `db:"private"`
}

type Channels []Channel

func (u *UserService) OccupiedChannels() (Channels, error) {
	channels := Channels{}

	err := u.db.Select(&channels, "SELECT channel_name, channel_id, private, description FROM user_channels WHERE email_address = ?", u.email)
	if err != nil {
		slog.Error("could not load channels", slog.String("error", err.Error()))

		return nil, fmt.Errorf("could not find channels for user: %w", err)
	}

	return channels, nil
}
