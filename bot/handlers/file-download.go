package handlers

import (
	"ImpiFilesBot/bot/lib"
	"ImpiFilesBot/bot/lib/filters"
	"ImpiFilesBot/internal/domain"
	"ImpiFilesBot/internal/query"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mymmrac/telego"
)

type downloadFile struct {
	content io.Reader
	name    string
}

func (d *downloadFile) Name() string {
	return d.name
}

func (d *downloadFile) Read(p []byte) (n int, err error) {
	return d.content.Read(p)
}

type FileDownloadHandler struct {
	fileDownloadHandler
	filters.TextMessageFilter
}

func NewFileDownloadHandler() *FileDownloadHandler {
	return &FileDownloadHandler{}
}

type fileDownloadHandler struct {
}

func (h *fileDownloadHandler) Handle(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) error {
	req := query.DownloadQuery{
		TelegramUserID: update.Message.From.ID,
		FileName:       update.Message.Text,
	}
	resp, err := bot.Service.DownloadQuery.Handle(&req)
	if errors.Is(err, domain.ErrNoAccess) {
		_, err := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "У вас нет доступа к этому файлу",
		})
		if err != nil {
			return fmt.Errorf("fileDownloadHandler: can't send message: %w", err)
		}
	} else if errors.Is(err, domain.ErrFileSizeLimit) {
		_, err := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   fmt.Sprintf("Файл слишком большой: %v (максимальный размер: %v)", resp.File.Name, domain.FileSizeLimit),
		})
		if err != nil {
			return fmt.Errorf("fileDownloadHandler: can't send message: %w", err)
		}
	}
	if err != nil {
		return err
	}

	_, err = bot.Bot.SendDocument(ctx, &telego.SendDocumentParams{
		ChatID: update.Message.Chat.ChatID(),
		Document: telego.InputFile{
			File: &downloadFile{name: resp.File.Name, content: resp.File.Content},
		},
	})
	if err != nil {
		return fmt.Errorf("fileDownloadHandler: can't send document: %w", err)
	}
	return nil
}
