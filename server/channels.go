package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/jtarchie/sqlite-chat/server/templates"
	"github.com/jtarchie/sqlite-chat/services"
	"github.com/labstack/echo/v4"
)

func setupChannels(db *sqlx.DB, server *echo.Echo) {
	server.GET("/dashboard", func(c echo.Context) error {
		session, err := sessionStore.Get(c.Request(), "chat-app")
		if err != nil {
			slog.Error("could not get session", slog.String("error", err.Error()))
		} else {
			email := session.Get("email")
			slog.Info("session email", slog.String("email", email))
		}

		user := services.NewUser(db, "bot@example.com")

		channels, err := user.OccupiedChannels()
		if err != nil {
			return fmt.Errorf("could not load user channels: %w", err)
		}

		return c.Redirect(
			http.StatusTemporaryRedirect,
			fmt.Sprintf("/dashboard/channels/%d", channels[0].ID),
		)
	})
	server.GET("/dashboard/channels/:id", func(c echo.Context) error {
		channelID, _ := strconv.Atoi(c.Param("id"))

		user := services.NewUser(db, "bot@example.com")
		channel := services.NewChannel(db, channelID)

		return stream(c, func(writer io.Writer) {
			templates.WriteDashboard(writer, user, channel)
		})
	})
}
