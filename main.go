package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	URN            string `help:"urn for the database" default:"file:test.db?cache=shared&mode=memory"`
	Port           int    `help:"port to listen for HTTP connections" default:"8080"`
	ClientID       string `help:"client ID to hello.dev" required:"" env:"OAUTH_CLIENT_ID"`
	ClientSecret   string `help:"client Secret to hello.dev" required:"" env:"OAUTH_CLIENT_SECRET"`
	ClientRedirect string `help:"client endpoint for the redirect" required:"" env:"OAUTH_CLIENT_REDIRECT"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &CLI{}
	ctx := kong.Parse(cli)

	err := ctx.Run()
	if err != nil {
		log.Fatalf("could not execute: %s", err)
	}
}

func (c *CLI) Run() error {
	db, err := NewDB(c.URN)
	if err != nil {
		return fmt.Errorf("could not setup DB: %w", err)
	}

	server, err := NewServer(
		c.Port,
		db,
		c.ClientID,
		c.ClientSecret,
		c.ClientRedirect,
	)
	if err != nil {
		return fmt.Errorf("could not setup server: %w", err)
	}

	return server.Start(fmt.Sprintf("0.0.0.0:%d", c.Port))
}
