package service

import (
	"ImpiFilesBot/internal/auth"
	"ImpiFilesBot/internal/command"
	"ImpiFilesBot/internal/domain"
	"ImpiFilesBot/internal/query"
)

type Service struct {
	LsQuery       *query.LsQueryHandler
	DownloadQuery *query.DownloadQueryHandler
	UploadCommand *command.UploadCommandHandler
	ChDirCommand  *command.ChDirCommandHandler
	RootCommand   *command.RootCommandHandler
}

func NewService(auth auth.AuthService, fs domain.FileRepository) *Service {
	return &Service{
		LsQuery:       query.NewLsQueryHandler(auth, fs),
		DownloadQuery: query.NewDownloadQueryHandler(auth, fs),
		UploadCommand: command.NewUploadCommandHandler(fs, auth),
		ChDirCommand:  command.NewChDirCommandHandler(fs, auth),
		RootCommand:   command.NewRootCommandHandler(fs, auth),
	}
}
