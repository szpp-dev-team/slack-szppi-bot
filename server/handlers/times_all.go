package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const timesAllChannelID = "C04MX77LHQE"

type HandlerTimesAll struct {
	c *slack.Client
}

func NewHandlerTimesAll(c *slack.Client) *HandlerTimesAll {
	return &HandlerTimesAll{c}
}

func (h *HandlerTimesAll) Handle(c echo.Context, messageEvent *slackevents.MessageEvent) error {
	if isReplyMessage(messageEvent) {
		return nil
	}
	user, err := h.c.GetUserInfo(messageEvent.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if messageEvent.Channel == timesAllChannelID {
		return nil // skip messages in times_all
	}
	msgOptList := []slack.MsgOption{
		slack.MsgOptionUsername(user.Profile.DisplayName),
		slack.MsgOptionIconURL(user.Profile.Image192),
		slack.MsgOptionAttachments(messageEvent.Attachments...),
	}
	// bot による message の場合は text は投稿しない
	// NOTE: 今後「おもしろメッセージ bot」みたいのが出てきた時にどうするか考える必要がある
	if !user.IsBot {
		msgOptList = append(msgOptList, slack.MsgOptionText(messageEvent.Text, false))
	}

	if _, _, err := h.c.PostMessage(timesAllChannelID, msgOptList...); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func isReplyMessage(messageEvent *slackevents.MessageEvent) bool {
	return messageEvent.ThreadTimeStamp != ""
}
