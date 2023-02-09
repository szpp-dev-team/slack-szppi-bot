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
	if messageEvent.Channel == TimesAllChannelID {
		return
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

	if _, _, err := h.c.PostMessage(
		TimesAllChannelID,
		msgOptList...,
	); err != nil {
		log.Println(err)
		return
	}
}

func isReplyMessage(messageEvent *slackevents.MessageEvent) bool {
	return messageEvent.ThreadTimeStamp != ""
}
