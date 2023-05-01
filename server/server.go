package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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

	e := echo.New()
	e.Use(middleware.Verify(signingSecret))

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
		messageEvent := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
		if messageEvent == nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return timesAllHandler.Handle(c, messageEvent)
	}, middleware.URLVerification())

	return e
}
