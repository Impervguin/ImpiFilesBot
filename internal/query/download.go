package query

import (
	"ImpiFilesBot/internal/auth"
	domain "ImpiFilesBot/internal/domain"
	"path/filepath"
)

type DownloadQueryHandler struct {
	auth auth.AuthService
	fs   domain.FileRepository
}

func NewDownloadQueryHandler(auth auth.AuthService, fs domain.FileRepository) *DownloadQueryHandler {
	return &DownloadQueryHandler{auth: auth, fs: fs}
}

type DownloadQuery struct {
	TelegramUserID int64
	FileName       string
}

type DownloadResponse struct {
	File *domain.FileData
}

func (h *DownloadQueryHandler) Handle(query *DownloadQuery) (*DownloadResponse, error) {
	telegramUser, err := h.auth.Authenticate(query.TelegramUserID)
	if err != nil {
		return nil, err
	}

	fpath := filepath.Join(telegramUser.Cwd(), query.FileName)

	ok, err := telegramUser.HaveReadAccess(fpath)
	if err != nil || !ok {
		return nil, domain.ErrNoAccess
	}

	file, err := h.fs.Download(fpath)
	if err != nil {
		return nil, err
	}

	return &DownloadResponse{File: file}, nil
}
