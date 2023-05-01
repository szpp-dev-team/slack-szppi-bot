package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
	"github.com/szpp-dev-team/szpp-slack-bot/commands"
)

type SlashCommandHandler struct {
	cmdByName map[string](func(slashCmd *slack.SlashCommand) error)
}

func NewSlashCommandHandler(cmds ...commands.SlashCommand) *SlashCommandHandler {
	cmdByName := make(map[string](func(slashCmd *slack.SlashCommand) error), len(cmds))
	for _, cmd := range cmds {
		cmdByName[cmd.Name()] = cmd.Handle
	}
	return &SlashCommandHandler{cmdByName}
}

func (s *SlashCommandHandler) Handle(c echo.Context, slashCmd *slack.SlashCommand) error {
	rw := c.Response().Writer

	tokens := strings.Fields(slashCmd.Text)
	if len(tokens) == 0 {
		return nil
	}

	handle, ok := s.cmdByName[tokens[0]]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	rw.Header().Add("Content-Type", "application/json")
	msg := &slack.Msg{ResponseType: slack.ResponseTypeInChannel} // 打ったコマンドを表示させる
	if err := json.NewEncoder(rw).Encode(msg); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	go func() {
		if err := handle(slashCmd); err != nil {
			log.Println(err)
		}
	}()

	return nil
}
