package handlers

import (
	"context"
	"log"
	"net/http"
	"regexp"

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
