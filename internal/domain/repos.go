package domain

import "github.com/google/uuid"

type TelegramUserRepository interface {
	GetByTelegramID(telegramID int64) (TelegramUser, error)
	Update(id uuid.UUID, updateFn func(use TelegramUser) (TelegramUser, error)) (TelegramUser, error)
}

type FileRepository interface {
	GetByPath(path string) (*File, error)
	Save(file *FileData) error
	Download(path string) (*FileData, error)
	Delete(path string) error
	List(path string) (*Directory, error)
}
