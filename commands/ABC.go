package commands

import (
	"github.com/slack-go/slack"
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

	res := MakeUrl(id(), slashCmd.Text)

	_, _, _, err := o.c.SendMessage(slashCmd.ChannelID, slack.MsgOptionText(res, false))

	return err
}
