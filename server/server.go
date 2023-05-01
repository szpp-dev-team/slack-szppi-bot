package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/szpp-dev-team/szpp-slack-bot/commands"
	"github.com/szpp-dev-team/szpp-slack-bot/server/handlers"
	"github.com/szpp-dev-team/szpp-slack-bot/server/middleware"
	"google.golang.org/api/customsearch/v1"
)

func New(slackClient *slack.Client, customsearchService *customsearch.Service, signingSecret string) *echo.Echo {
	slashHandler := handlers.NewSlashCommandHandler(
		commands.NewSubHandlerOmikuji(slackClient),
		commands.NewSubHandlerImage(slackClient, customsearchService),
	)
	timesAllHandler := handlers.NewHandlerTimesAll(slackClient)
	channelNotificatorHandler := handlers.NewHandlerNewChannelNotificator(slackClient)

	e := echo.New()
	e.Use(middleware.Verify(signingSecret), echomiddleware.Logger())

	e.POST("/slack/slash_commands", func(c echo.Context) error {
		slashCmd, err := slack.SlashCommandParse(c.Request())
		if err != nil {
			log.Println("failed to parse slash command:", err)
			return err
		}

		return slashHandler.Handle(c, &slashCmd)
	})
	e.POST("/slack/events", func(c echo.Context) error {
		eventsAPIEvent, _ := c.Get("events_api_event").(*slackevents.EventsAPIEvent)
		if eventsAPIEvent == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "events_api_event not found")
		}
		if eventsAPIEvent.Type != slackevents.CallbackEvent {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		switch ty := eventsAPIEvent.InnerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			return timesAllHandler.Handle(c, ty)
		case *slackevents.ChannelCreatedEvent:
			return channelNotificatorHandler.Handle(ty)
		}

		return echo.NewHTTPError(http.StatusBadRequest)
	}, middleware.URLVerification())

	return e
}
