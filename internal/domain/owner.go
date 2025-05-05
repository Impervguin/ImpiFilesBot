package domain

import (
	"fmt"

	"github.com/google/uuid"
)

var _ TelegramUser = &Owner{}

type Owner struct {
	id         uuid.UUID
	username   string
	telegramID int64
	cwd        string
}

func NewOwner(username string, telegramID int64, cwd string) *Owner {
	return &Owner{
		id:         uuid.New(),
		username:   username,
		telegramID: telegramID,
		cwd:        cwd,
	}
}

func (o *Owner) ID() uuid.UUID {
	return o.id
}

func (o *Owner) Username() string {
	return o.username
}

func (o *Owner) TelegramID() int64 {
	return o.telegramID
}

func (o *Owner) UpdateCwd(path string) error {
	o.cwd = path
	return nil
}

func (o *Owner) Cwd() string {
	return o.cwd
}

func (o *Owner) UpdateUsername(TelegramUsername string) error {
	if TelegramUsername == "" {
		return fmt.Errorf("TelegramUsername cannot be empty")
	}
	o.username = TelegramUsername
	return nil
}

func (o *Owner) HaveReadAccess(path string) (bool, error) {
	return true, nil
}

func (o *Owner) HaveWriteAccess(path string) (bool, error) {
	return true, nil
}

func (o *Owner) HaveConfigAccess() (bool, error) {
	return true, nil
}
