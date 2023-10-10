package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"github.com/dghubble/sessions"
	"github.com/jtarchie/sqlite-chat/server/templates"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

//nolint
var sessionStore = sessions.NewCookieStore[string](
	sessions.DebugCookieConfig,
	[]byte("something-secret-this-way-comes"),
	nil,
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
		RedirectURL:  fmt.Sprintf("%s/auth/github/callback", clientEndpoint),
		Endpoint:     githubOAuth2.Endpoint,
	}

	stateConfig := gologin.DebugOnlyCookieConfig
	server.GET("/auth/github/login", echo.WrapHandler(
		github.StateHandler(
			stateConfig,
			github.LoginHandler(config, nil),
		),
	))

	server.GET("/auth/github/callback", echo.WrapHandler(
		github.StateHandler(
			stateConfig,
			github.CallbackHandler(config, issueSession(), nil),
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

		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			slog.Error("could not log in user", slog.String("error", err.Error()))
			http.Error(w, "could not log in user", http.StatusInternalServerError)

			return
		}
		// Register user
		// save session for the user
		session := sessionStore.New("chat-app")
		session.Set("email", *githubUser.Email)

		err = session.Save(w)
		if err != nil {
			slog.Error("could not set session", slog.String("error", err.Error()))
			http.Error(w, "could not log in user", http.StatusInternalServerError)

			return
		}

		slog.Info("logged in", slog.String("login", *githubUser.Login))
		http.Redirect(w, req, "/dashboard", http.StatusFound)
	}

	return http.HandlerFunc(fn)
}
