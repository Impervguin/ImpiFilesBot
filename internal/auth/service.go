package auth

import (
	domain "ImpiFilesBot/internal/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	TelegramUserRepository domain.TelegramUserRepository
}

func NewAuthService(TelegramUserRepository domain.TelegramUserRepository) *AuthService {
	return &AuthService{TelegramUserRepository: TelegramUserRepository}
}

func (s *AuthService) Authenticate(telegramID int64) (domain.TelegramUser, error) {
	TelegramUser, err := s.TelegramUserRepository.GetByTelegramID(telegramID)
	if err != nil {
		return nil, err
	}
	return TelegramUser, nil
}

func (s *AuthService) UpdateTelegramUsername(userID uuid.UUID, username string) (domain.TelegramUser, error) {
	return s.TelegramUserRepository.Update(userID, func(user domain.TelegramUser) (domain.TelegramUser, error) {
		err := user.UpdateUsername(username)
		return user, err
	})
}

func (S *AuthService) UpdateCwd(userID uuid.UUID, path string) (domain.TelegramUser, error) {
	return S.TelegramUserRepository.Update(userID, func(user domain.TelegramUser) (domain.TelegramUser, error) {
		err := user.UpdateCwd(path)
		return user, err
	})
}
