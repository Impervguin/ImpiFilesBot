package query

import (
	"ImpiFilesBot/internal/auth"
	domain "ImpiFilesBot/internal/domain"
	"errors"
	"fmt"
)

type LsQueryHandler struct {
	auth auth.AuthService
	fs   domain.FileRepository
}

func NewLsQueryHandler(auth auth.AuthService, fs domain.FileRepository) *LsQueryHandler {
	return &LsQueryHandler{auth: auth, fs: fs}
}

type LsQuery struct {
	TelegramUserID int64
}
type LsResponse struct {
	Dir *domain.Directory
}

func (h *LsQueryHandler) Handle(query *LsQuery) (*LsResponse, error) {
	telegramUser, err := h.auth.Authenticate(query.TelegramUserID)
	if err != nil {
		return nil, err
	}

	ok, err := telegramUser.HaveReadAccess(telegramUser.Cwd())
	if err != nil || !ok {
		return nil, domain.ErrNoAccess
	}

	files, err := h.fs.List(telegramUser.Cwd())
	if errors.Is(err, domain.ErrDirNotFound) {
		return nil, fmt.Errorf("%w:%v", domain.ErrCwdNotFound, telegramUser.Cwd())
	} else if err != nil {
		return nil, err
	}

	return &LsResponse{Dir: files}, nil
}
