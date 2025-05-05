package domain

import "github.com/google/uuid"

type TelegramUser interface {
	ID() uuid.UUID
	TelegramID() int64
	Username() string
	Cwd() string
	UpdateCwd(path string) error
	UpdateUsername(username string) error
	HaveReadAccess(path string) (bool, error)
	HaveWriteAccess(path string) (bool, error)
	HaveConfigAccess() (bool, error)
}
