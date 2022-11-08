package repositories

import "teneta-tg/internal/entities"

type UserRepository interface {
	FindBy(params map[string]interface{}) (user *entities.User, err error)
	Create(user *entities.User) (*entities.User, error)
}
