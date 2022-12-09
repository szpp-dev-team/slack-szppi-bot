package commands

import (
	"strings"

	"github.com/slack-go/slack"
	"github.com/szpp-dev-team/szpp-slack-bot/ABC"
)

type SubHandlerABC struct {
	c *slack.Client
}

func NewSubHandlerABC(c *slack.Client) *SubHandlerABC {
	return &SubHandlerABC{c}
}

func (o *SubHandlerABC) Name() string {
	return "ABC"
}

func (o *SubHandlerABC) Handle(slashCmd *slack.SlashCommand) error {

	res := ABC.MakeUrl(ABC.Id(), strings.Fields(slashCmd.Text)[1])

	_, _, _, err := o.c.SendMessage(slashCmd.ChannelID, slack.MsgOptionText(res, false))

	return err
}
