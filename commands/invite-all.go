package commands

import "github.com/slack-go/slack"

type CommandInviteAll struct {
	slackClient *slack.Client
}

func NewCommandInviteAll(slackClient *slack.Client) *CommandInviteAll {
	return &CommandInviteAll{slackClient}
}

func (c *CommandInviteAll) Name() string {
	return "invite-all"
}

func (c *CommandInviteAll) Handle(slackCmd *slack.SlashCommand) error {
	users, err := c.slackClient.GetUsers()
	if err != nil {
		return err
	}
	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}
	if _, err := c.slackClient.InviteUsersToConversation(slackCmd.ChannelID, userIDs...); err != nil {
		return err
	}
	return nil
}
