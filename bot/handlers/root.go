package handlers

import (
	"ImpiFilesBot/bot/keyboards"
	"ImpiFilesBot/bot/lib"
	"ImpiFilesBot/bot/lib/filters"
	"ImpiFilesBot/internal/command"
	"ImpiFilesBot/internal/domain"
	"ImpiFilesBot/internal/query"
	"context"
	"errors"
	"fmt"

	"github.com/mymmrac/telego"
)

var _ lib.Handler = &RootHandler{}

const (
	rootCommand = "/root"
)

type RootHandler struct {
	rootHandler
	filters.CommandFilter
}

func NewRootHandler() *RootHandler {
	return &RootHandler{
		CommandFilter: *filters.NewCommandFilter(rootCommand),
	}
}

type rootHandler struct {
}

func (h *rootHandler) Handle(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) error {
	request := command.RootCommand{
		TelegramUserID: update.Message.From.ID,
	}

	err := bot.Service.RootCommand.Handle(&request)
	if errors.Is(err, domain.ErrNoAccess) {
		_, err := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "У вас нет доступа к корню файловой системы",
		})
		if err != nil {
			return fmt.Errorf("rootHandler: can't send message: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("rootHandler:can't handle root command: %w", err)
	}

	lsRequest := query.LsQuery{
		TelegramUserID: update.Message.From.ID,
	}
	resp, err := bot.Service.LsQuery.Handle(&lsRequest)
	if err != nil {
		return fmt.Errorf("rootHandler:can't handle ls query: %w", err)
	}

	keyboard := keyboards.NewListDirKeyboard(resp.Dir)
	_, err = bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        "Переход к корню успешен",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("rootHandler: can't send message: %w", err)
	}

	return nil
}
