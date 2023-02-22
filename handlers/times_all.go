package handlers

import (
	"context"
	"log"
	"net/http"
	"regexp"

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

	innerEvent := eventsAPIEvent.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.ReactionAddedEvent:
		name, _ := h.c.GetUserInfo(ev.User)
		log.Println("reaction add")
		log.Printf("%#v", ev)
		h.reflectedReaction(name.Name, ev.Item.Timestamp)
		return
		// 今後eventを拡張する際には、この下にどんどん書いてく？
		// 今回はreactionEventだけを試しに書いたけど、ここのcaseでmessageEventも書いたほうがいいかな？(うまく動けば)
	}

	messageEvent := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if messageEvent == nil {
		return
	}
	log.Println(eventsAPIEvent.Type)
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

func (h *HandlerTimesAll) reflectedReaction(user string, timeStamp string) {
	channelID := ""
	channels, _, err := h.c.GetConversations(&slack.GetConversationsParameters{})
	if err != nil {
		log.Println(err)
		return
	}

	for _, channel := range channels {
		//log.Println(channel.Name, regexp.MustCompile("times_").MatchString(channel.Name))
		if regexp.MustCompile("times_").MatchString(channel.Name) && channel.Name[7:] == user {
			channelID = channel.ID
			log.Println(channelID)
		}
	}
	history, err := h.c.GetConversationHistoryContext(context.Background(), &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     100,
	})
	for _, message := range history.Messages {
		if message.Timestamp == timeStamp {
			log.Println(message.Text)
		}
	}
}

func isReplyMessage(messageEvent *slackevents.MessageEvent) bool {
	return messageEvent.ThreadTimeStamp != ""
}
