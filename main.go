package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/szpp-dev-team/szpp-slack-bot/server"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

func main() {
	botUserOauthToken := os.Getenv("BOT_USER_OAUTH_TOKEN")
	signingSecret := os.Getenv("SIGNING_SECRET")
	customsearchApiKey := os.Getenv("CUSTOM_SEARCH_API_KEY")

	port := getenvOr("PORT", "8080")

	client := slack.New(botUserOauthToken)
	customsearchService, err := customsearch.NewService(context.Background(), option.WithAPIKey(customsearchApiKey))
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(client, customsearchService, signingSecret)
	if err := srv.Start(fmt.Sprintf(":%v", port)); err != nil {
		log.Fatal(err)
	}
}

func getenvOr(key, altValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = altValue
	}
	return value
}
