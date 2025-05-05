package command

import (
	"ImpiFilesBot/internal/auth"
	domain "ImpiFilesBot/internal/domain"
	"fmt"
)

type RootCommandHandler struct {
	auth auth.AuthService
	fs   domain.FileRepository
}

type RootCommand struct {
	TelegramUserID int64
}

func NewRootCommandHandler(fs domain.FileRepository, auth auth.AuthService) *RootCommandHandler {
	return &RootCommandHandler{fs: fs, auth: auth}
}

func (h *RootCommandHandler) Handle(comm *RootCommand) error {
	telegramUser, err := h.auth.Authenticate(comm.TelegramUserID)
	if err != nil {
		return err
	}

	ok, err := telegramUser.HaveReadAccess(domain.Root)
	if err != nil || !ok {
		return domain.ErrNoAccess
	}

	_, err = h.fs.List(domain.Root)
	if err != nil {
		return fmt.Errorf("can't list root: %w", err)
	}

	_, err = h.auth.UpdateCwd(telegramUser.ID(), domain.Root)
	if err != nil {
		return fmt.Errorf("can't update cwd: %w", err)
	}

	return nil
}
