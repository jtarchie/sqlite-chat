package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

//go:generate go run github.com/valyala/quicktemplate/qtc -ext=html -dir=templates/ -skipLineComments

func New(
	db *sqlx.DB,
	clientID string,
	clientSecret string,
	clientEndpoint string,
) (*echo.Echo, error) {
	server := echo.New()
	server.Use(slogecho.New(slog.Default()))

	setupAuth(clientID, clientSecret, clientEndpoint, server)
	setupChannels(db, server)

	return server, nil
}

func stream(c echo.Context, fn func(io.Writer)) error {
	reader, writer := io.Pipe()

	go func() {
		fn(writer)
		_ = writer.Close()
	}()

	err := c.Stream(http.StatusOK, "text/html; charset=utf-8", reader)
	if err != nil {
		return fmt.Errorf("could not stream: %w", err)
	}

	return nil
}
