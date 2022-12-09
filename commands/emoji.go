package commands

import (
	"net/http"
	"strings"

	"github.com/slack-go/slack"
)

type SubHandlerEmoji struct {
	c *slack.Client
}

func NewSubHandlerEmoji(c *slack.Client) *SubHandlerEmoji {
	return &SubHandlerEmoji{c}
}

func (o *SubHandlerEmoji) Name() string {
	return "emoji"
}

func (o *SubHandlerEmoji) Handle(slashCmd *slack.SlashCommand) error {
	tokens := strings.Fields(slashCmd.Text)
	text := strings.Join(tokens[1:], " ")

	req, _ := http.NewRequest(http.MethodGet, "https://emoji-gen.ninja/emoji", nil)
	q := req.URL.Query()
	q.Set("align", "center")
	q.Set("back_color", "00000000")
	q.Set("color", "EC71A1FF")
	q.Set("font", "notosans-mono-bold")
	q.Set("locale", "ja")
	q.Set("public_fg", "true")
	q.Set("size_fixed", "false")
	q.Set("stretch", "true")
	q.Set("text", text)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileParam := slack.FileUploadParameters{
		Reader:   resp.Body,
		Filename: strings.Join(strings.Fields(text), ""),
		Channels: []string{slashCmd.ChannelID},
	}
	if _, err := o.c.UploadFile(fileParam); err != nil {
		return err
	}

	return err
}
