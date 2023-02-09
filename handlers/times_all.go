package handlers

import (
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const TimesAllChannelID = "C04MX77LHQE"

type HandlerTimesAll struct {
	c *slack.Client
}

func NewHandlerTimesAll(c *slack.Client) *HandlerTimesAll {
	return &HandlerTimesAll{c}
}

func (h *HandlerTimesAll) Handle(w http.ResponseWriter, eventsAPIEvent *slackevents.EventsAPIEvent) {
	if eventsAPIEvent.Type != slackevents.CallbackEvent {
		return
	}
	messageEvent := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if messageEvent == nil {
		return
	}
	if isReplyMessage(messageEvent) {
		return
	}
	log.Println(messageEvent)
	user, err := h.c.GetUserInfo(messageEvent.User)
	if err != nil {
		log.Println(err)
		return
	}
	// message を送ったのが bot もしくは message の場所が #times_all なら無視
	if user.IsBot || messageEvent.Channel == TimesAllChannelID {
		return
	}
	log.Println(user.Name, user.Profile.Image192)
	if _, _, err := h.c.PostMessage(
		TimesAllChannelID,
		slack.MsgOptionText(messageEvent.Text, false),
		slack.MsgOptionUsername(user.Profile.DisplayName),
		slack.MsgOptionIconURL(user.Profile.Image192),
		slack.MsgOptionAttachments(messageEvent.Attachments...),
	); err != nil {
		log.Println(err)
		return
	}
}

func isReplyMessage(messageEvent *slackevents.MessageEvent) bool {
	return messageEvent.ThreadTimeStamp != ""
}