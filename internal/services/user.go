package services

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/repositories"
)

type UserService interface {
	FirstOrCreate(chatID int64, firstName, lastName, language string) (*entities.User, error)
	Save(user *entities.User)
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

		zap.S().Error(err)
	}

	return user, nil
}

func (s *userService) Save(user *entities.User) {
	if err := s.repo.Save(user); err != nil {
		zap.S().Error(err)
	}
}
