package handlers

import (
	"context"
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
	if eventsAPIEvent.Type == string(slackevents.ReactionAdded) {
		reaction := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
		h.reflectedReaction(reaction.ItemUser, reaction.Item.Message.Timestamp)
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
	); err != nil {
		log.Println(err)
		return
	}
}

func (h *HandlerTimesAll) reflectedReaction(user string, timeStamp string) {
	channelID := ""
	channels, _, err := h.c.GetConversations(&slack.GetConversationsParameters{
		Types: []string{"public"},
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(user)

	for _, channel := range channels {
		log.Println(channel.User)
		if channel.User == user {
			channelID = channel.ID
		}
	}
	history, err := h.c.GetConversationHistoryContext(context.Background(), &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
	})
	for _, message := range history.Messages {
		if message.Timestamp == timeStamp {
			log.Println(message.Text)
		}
	}
}
