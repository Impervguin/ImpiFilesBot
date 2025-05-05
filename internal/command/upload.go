package command

import (
	"ImpiFilesBot/internal/auth"
	domain "ImpiFilesBot/internal/domain"
	"context"
	"io"
	"path/filepath"
)

type UploadCommandHandler struct {
	fs   domain.FileRepository
	auth auth.AuthService
}

func NewUploadCommandHandler(fs domain.FileRepository, auth auth.AuthService) *UploadCommandHandler {
	return &UploadCommandHandler{fs: fs, auth: auth}
}

type UploadCommand struct {
	TelegramUserID int64
	FileName       string
	FileData       io.Reader
}

func (h *UploadCommandHandler) Handle(ctx context.Context, comm *UploadCommand) error {
	telegramUser, err := h.auth.Authenticate(comm.TelegramUserID)
	if err != nil {
		return err
	}

	file, err := domain.NewFile(comm.FileName, filepath.Join(telegramUser.Cwd(), comm.FileName))
	if err != nil {
		return err
	}
	ok, err := telegramUser.HaveWriteAccess(file.Path)
	if err != nil || !ok {
		return domain.ErrNoAccess
	}

	fileData, err := domain.NewFileData(file, comm.FileData)
	if err != nil {
		return err
	}

	err = h.fs.Save(fileData)
	if err != nil {
		return err
	}

	return nil
}
