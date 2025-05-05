package handlers

import (
	"ImpiFilesBot/bot/lib"
	"ImpiFilesBot/bot/lib/filters"
	"ImpiFilesBot/internal/command"
	"ImpiFilesBot/internal/domain"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
)

const (
	fileUrlFormat = "https://api.telegram.org/file/bot%s/%s"
)

func uploadFile(ctx context.Context, token, filePath string) (io.Reader, error) {
	fileBuf := make([]byte, 1024)
	code, body, err := fasthttp.Get(fileBuf, fmt.Sprintf(fileUrlFormat, token, filePath))
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("uploadFile: bad response code: %d", code)
	}
	return bytes.NewReader(body), nil
}

type FileUploadHandler struct {
	fileUploadHandler
	filters.FileFilter
}

func NewFileUploadHandler() *FileUploadHandler {
	return &FileUploadHandler{}
}

type fileUploadHandler struct {
}

func (h *fileUploadHandler) Handle(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) error {
	req := command.UploadCommand{
		TelegramUserID: update.Message.From.ID,
		FileName:       update.Message.Document.FileName,
	}
	file, err := bot.Bot.GetFile(ctx, &telego.GetFileParams{
		FileID: update.Message.Document.FileID,
	})
	if err != nil {
		_, sendErr := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "Не удалось загрузить файл",
		})
		if sendErr != nil {
			return fmt.Errorf("fileUploadHandler: can't send message on error %w: %w", err, sendErr)
		}
		return fmt.Errorf("fileUploadHandler: can't get file: %w", err)
	}

	reader, err := uploadFile(ctx, bot.Bot.Token(), file.FilePath)
	if err != nil {
		_, sendErr := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "Не удалось загрузить файл",
		})
		if sendErr != nil {
			return fmt.Errorf("fileUploadHandler: can't send message on error %w: %w", err, sendErr)
		}
		return fmt.Errorf("fileUploadHandler: can't upload file: %w", err)
	}

	req.FileData = reader
	err = bot.Service.UploadCommand.Handle(ctx, &req)
	if errors.Is(err, domain.ErrNoAccess) {
		_, err := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "У вас нет доступа  на загрузку файлов в эту директорию",
		})
		if err != nil {
			return fmt.Errorf("fileUploadHandler: can't send message: %w", err)
		}
	} else if err != nil {
		_, sendErr := bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   "Не удалось загрузить файл",
		})
		if sendErr != nil {
			return fmt.Errorf("fileUploadHandler: can't send message on error %w: %w", err, sendErr)
		}
		return fmt.Errorf("fileUploadHandler: can't upload file: %w", err)
	}

	_, err = bot.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: update.Message.Chat.ChatID(),
		Text:   fmt.Sprintf("Файл загружен: %s", update.Message.Document.FileName),
	})
	if err != nil {
		return fmt.Errorf("fileUploadHandler: can't send message: %w", err)
	}

	return nil
}
