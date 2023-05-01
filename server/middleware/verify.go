package middleware

import (
	"bytes"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

func Verify(signingSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				return err
			}
			sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
			if err != nil {
				log.Println(err)
				return err
			}
			if _, err := sv.Write(b); err != nil {
				log.Println(err)
				return err
			}
			if err := sv.Ensure(); err != nil {
				log.Println(err)
				return err
			}
			r.Body = io.NopCloser(bytes.NewReader(b))

			return next(c)
		}
	}
}
