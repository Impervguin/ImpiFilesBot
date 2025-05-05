package usermap

import (
	domain "ImpiFilesBot/internal/domain"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type TelegramUserMap struct {
	users     map[int64]domain.TelegramUser
	uuidUsers map[uuid.UUID]domain.TelegramUser
	mut       sync.RWMutex
}

func NewTelegramUserMap() *TelegramUserMap {
	return &TelegramUserMap{
		users:     make(map[int64]domain.TelegramUser),
		uuidUsers: make(map[uuid.UUID]domain.TelegramUser),
		mut:       sync.RWMutex{},
	}
}

func (m *TelegramUserMap) GetByTelegramID(telegramID int64) (domain.TelegramUser, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	user, ok := m.users[telegramID]
	if !ok {
		return nil, fmt.Errorf("%w: %v", domain.ErrUserNotFound, telegramID)
	}
	return user, nil
}

func (m *TelegramUserMap) Update(ID uuid.UUID, updateFn func(use domain.TelegramUser) (domain.TelegramUser, error)) (domain.TelegramUser, error) {
	m.mut.Lock()
	defer m.mut.Unlock()
	user, ok := m.uuidUsers[ID]
	if !ok {
		return nil, fmt.Errorf("%w: %v", domain.ErrUserNotFound, ID)
	}
	updatedUser, err := updateFn(user)
	if err != nil {
		return nil, err
	}
	m.users[updatedUser.TelegramID()] = updatedUser
	m.uuidUsers[ID] = updatedUser
	return updatedUser, nil
}

func (m *TelegramUserMap) Set(telegramID int64, user domain.TelegramUser) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.users[telegramID] = user
	m.uuidUsers[user.ID()] = user
}
