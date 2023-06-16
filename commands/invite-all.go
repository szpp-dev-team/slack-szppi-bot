package commands

import (
	"log"

	"github.com/slack-go/slack"
)

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
	userIDset := make(map[string]struct{}, len(users))
	for _, user := range users {
		if user.IsBot {
			continue
		}
		userIDset[user.ID] = struct{}{}
	}
	existUserIDs, _, err := c.slackClient.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: slackCmd.ChannelID,
	})
	if err != nil {
		return err
	}
	log.Println(existUserIDs)
	for _, userID := range existUserIDs {
		delete(userIDset, userID)
	}
	userIDs := make([]string, len(userIDset))
	for userID := range userIDset {
		userIDs = append(userIDs, userID)
	}
	if _, err := c.slackClient.InviteUsersToConversation(slackCmd.ChannelID, userIDs...); err != nil {
		return err
	}
	return nil
}
