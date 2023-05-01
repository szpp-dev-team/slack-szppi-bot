package commands

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

type CommandShowTimes struct {
	c *slack.Client
}

func NewCommandShowTimes(c *slack.Client) *CommandShowTimes {
	return &CommandShowTimes{c}
}

func (o *CommandShowTimes) Name() string {
	return "show-times"
}

func (o *CommandShowTimes) Handle(slashCmd *slack.SlashCommand) error {
	params := &slack.GetConversationsParameters{
		Types:           []string{"public_channel"},
		Limit:           1000,
		ExcludeArchived: true,
	}
	var channels []slack.Channel
	for {
		channels2, nextCursor, err := o.c.GetConversations(params)
		if err != nil {
			return err
		}
		if nextCursor != "" {
			params.Cursor = nextCursor
		} else {
			channels = channels2
			break
		}
	}

	textBuilder := &strings.Builder{}
	for _, channel := range channels {
		if _, err := textBuilder.WriteString(fmt.Sprintf("<#%s>\n", channel.ID)); err != nil {
			return err
		}
	}

	_, _, _, err := o.c.SendMessage(slashCmd.ChannelID, slack.MsgOptionText(textBuilder.String(), false))

	return err
}
