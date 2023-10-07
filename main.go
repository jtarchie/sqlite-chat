package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/sqlite-chat/services"
)

type CLI struct {
	URN            string `default:"file:test.db?cache=shared&mode=memory" help:"urn for the database"`
	Port           int    `default:"8080"                                  help:"port to listen for HTTP connections"`
	ClientID       string `env:"OAUTH_CLIENT_ID"                           help:"client ID to hello.dev"              required:""`
	ClientSecret   string `env:"OAUTH_CLIENT_SECRET"                       help:"client Secret to hello.dev"          required:""`
	ClientRedirect string `env:"OAUTH_CLIENT_REDIRECT"                     help:"client endpoint for the redirect"    required:""`
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
	db, err := services.NewDB(c.URN)
	if err != nil {
		return fmt.Errorf("could not setup DB: %w", err)
	}

	server, err := NewServer(
		db,
		c.ClientID,
		c.ClientSecret,
		c.ClientRedirect,
	)
	if err != nil {
		return fmt.Errorf("could not setup server: %w", err)
	}

	err = server.Start(fmt.Sprintf("0.0.0.0:%d", c.Port))
	if err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}
