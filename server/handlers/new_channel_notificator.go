package handlers

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type HandlerNewChannelNotificator struct {
	c *slack.Client
}

func NewHandlerNewChannelNotificator(c *slack.Client) *HandlerNewChannelNotificator {
	return &HandlerNewChannelNotificator{c}
}

func (h *HandlerNewChannelNotificator) Handle(ev *slackevents.ChannelCreatedEvent) error {
	text := fmt.Sprintf("<#%s> が新たに作られたっぴ！", ev.Channel.ID)
	_, _, err := h.c.PostMessage(timesAllChannelID, slack.MsgOptionText(text, false))
	return err
}
