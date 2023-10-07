package services

import "github.com/jmoiron/sqlx"

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
