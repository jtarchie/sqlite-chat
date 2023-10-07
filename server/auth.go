package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dghubble/gologin/v2"
	oauthHandler "github.com/dghubble/gologin/v2/oauth2"
	"github.com/jtarchie/sqlite-chat/server/templates"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func setupAuth(
	clientID string,
	clientSecret string,
	clientEndpoint string,
	server *echo.Echo,
) {
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
