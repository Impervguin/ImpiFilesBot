package command

import (
	"ImpiFilesBot/internal/auth"
	domain "ImpiFilesBot/internal/domain"
	"context"
	"fmt"
	"path/filepath"
)

type ChDirCommandHandler struct {
	auth auth.AuthService
	fs   domain.FileRepository
}

func NewChDirCommandHandler(fs domain.FileRepository, auth auth.AuthService) *ChDirCommandHandler {
	return &ChDirCommandHandler{fs: fs, auth: auth}
}

type ChDirCommand struct {
	TelegramUserID int64
	DirName        string
}

func (h *ChDirCommandHandler) Handle(ctx context.Context, comm *ChDirCommand) error {
	telegramUser, err := h.auth.Authenticate(comm.TelegramUserID)
	if err != nil {
		return err
	}

	ok, err := telegramUser.HaveReadAccess(telegramUser.Cwd())
	if err != nil || !ok {
		return domain.ErrNoAccess
	}

	var newCwd string
	if comm.DirName == ".." {
		newCwd = filepath.Dir(telegramUser.Cwd())
	} else {
		newCwd = filepath.Join(telegramUser.Cwd(), comm.DirName)
	}

	_, err = h.fs.List(newCwd)
	if err != nil {
		return fmt.Errorf("can't find new cwd %v: %w", newCwd, err)
	}

	_, err = h.auth.UpdateCwd(telegramUser.ID(), newCwd)
	if err != nil {
		return err
	}

	return nil
}
