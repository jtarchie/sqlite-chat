package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dghubble/gologin/v2"
	oauthHandler "github.com/dghubble/gologin/v2/oauth2"
	"github.com/jmoiron/sqlx"
	"github.com/jtarchie/sqlite-chat/services"
	"github.com/jtarchie/sqlite-chat/templates"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/valyala/quicktemplate/qtc -ext=html -dir=templates/ -skipLineComments

func NewServer(
	db *sqlx.DB,
	clientID string,
	clientSecret string,
	clientEndpoint string,
) (*echo.Echo, error) {
	server := echo.New()
	server.Use(slogecho.New(slog.Default()))

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("%s/callback", clientEndpoint),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://wallet.hello.coop/authorize",
			TokenURL: "https://wallet.hello.coop/oauth/token",
		},
	}

	stateConfig := gologin.DebugOnlyCookieConfig
	server.GET("/login", echo.WrapHandler(
		oauthHandler.StateHandler(
			stateConfig,
			oauthHandler.LoginHandler(config, nil),
		),
	))
	server.GET("/callback", echo.WrapHandler(
		oauthHandler.StateHandler(
			stateConfig,
			oauthHandler.CallbackHandler(config, issueSession(), nil),
		),
	))
	server.GET("/", func(c echo.Context) error {
		return stream(c, func(writer io.Writer) {
			templates.WriteLogin(writer, clientID, clientEndpoint)
		})
	})
	server.GET("/dashboard", func(c echo.Context) error {
		user := services.NewUser(db, "bot@example.com")

		return stream(c, func(writer io.Writer) {
			templates.WriteDashboard(writer, user)
		})
	})

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

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		token, _ := oauthHandler.TokenFromContext(ctx)
		state, _ := oauthHandler.StateFromContext(ctx)
		slog.Info("logged in", slog.String("state", state), slog.String("token", token.Expiry.String()))
	}

	return http.HandlerFunc(fn)
}
