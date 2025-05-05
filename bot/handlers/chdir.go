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

var _ lib.Handler = &ChDirHandler{}

type ChDirHandler struct {
	chDirHandler
	filters.DirectoryFilter
}

func NewChDirHandler() *ChDirHandler {
	return &ChDirHandler{}
}

type chDirHandler struct {
}

func (h *chDirHandler) Handle(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) error {
	request := command.ChDirCommand{
		TelegramUserID: update.Message.From.ID,
		DirName:        update.Message.Text,
	}

	err := bot.Service.ChDirCommand.Handle(ctx, &request)
	if errors.Is(err, domain.ErrNoAccess) {
		_, err := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "У вас нет доступа к этой директории",
		})
		if err != nil {
			return fmt.Errorf("chDirHandler: can't send message: %w", err)
		}
	} else if errors.Is(err, domain.ErrDirNotFound) {
		_, err := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "Директория не найдена",
		})
		if err != nil {
			return fmt.Errorf("chDirHandler: can't send message: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("chDirHandler: can't handle chdir command: %w", err)
	}

	lsRequest := query.LsQuery{
		TelegramUserID: update.Message.From.ID,
	}
	resp, err := bot.Service.LsQuery.Handle(&lsRequest)
	if err != nil {
		return fmt.Errorf("chDirHandler: can't handle ls query: %w", err)
	}
	keyboard := keyboards.NewListDirKeyboard(resp.Dir)
	_, err = bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        fmt.Sprintf("Переход к директории %s", resp.Dir.Name),
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("chDirHandler: can't send message: %w", err)
	}

	return nil
}
