package pgsql

import (
	"gorm.io/gorm"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/repositories"
)

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) repositories.UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) FindBy(params map[string]interface{}) (user *entities.User, err error) {
	if err = r.conn.Where(params).Preload("ProviderConfig").First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(user *entities.User) (*entities.User, error) {
	if err := r.conn.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Save(user *entities.User) error {
	return r.conn.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}
