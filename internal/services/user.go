package services

import (
	"errors"
	"gorm.io/gorm"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/repositories"
)

type UserService interface {
	FirstOrCreate(chatID int64, firstName, lastName, language string) (*entities.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) FirstOrCreate(chatID int64, firstName, lastName, language string) (*entities.User, error) {
	user, err := s.repo.FindBy(map[string]interface{}{"chat_id": chatID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.repo.Create(&entities.User{
				ChatID:    chatID,
				FirstName: firstName,
				LastName:  lastName,
				Language:  language,
			})
		}
	}

	return user, nil
}
