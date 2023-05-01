package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack/slackevents"
)

func URLVerification() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rw := c.Response().Writer
			b, err := io.ReadAll(c.Request().Body)
			if err != nil {
				log.Println(err)
				return err
			}
			eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(b), slackevents.OptionNoVerifyToken())
			if err != nil {
				log.Println(err)
				return err
			}
			switch eventsAPIEvent.Type {
			case slackevents.URLVerification:
				var chalResp slackevents.ChallengeResponse
				if err := json.NewDecoder(bytes.NewReader(b)).Decode(&chalResp); err != nil {
					return err
				}
				return json.NewEncoder(rw).Encode(chalResp)
			default:
				c.Request().Body = io.NopCloser(bytes.NewReader(b))
				c.Set("events_api_event", &eventsAPIEvent)
				return next(c)
			}
		}
	}
}
