package commands

import "github.com/slack-go/slack"

type SlashCommand interface {
	Handle(slashCmd *slack.SlashCommand) error
	Name() string
}
