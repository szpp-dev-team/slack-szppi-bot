package handlers

import (
	"encoding/json"
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

func (h *HandlerTimesAll) Handle(w http.ResponseWriter, b []byte) {
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(b), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if eventsAPIEvent.Type != slackevents.CallbackEvent {
		return
	}
	messageEvent := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if messageEvent == nil {
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
		slack.MsgOptionUsername(user.Name),
		slack.MsgOptionIconURL(user.Profile.Image192),
	); err != nil {
		log.Println(err)
		return
	}
}
