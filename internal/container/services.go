package container

import (
	"github.com/sarulabs/di"
	"teneta-tg/internal/constants"
	"teneta-tg/internal/repositories"
	"teneta-tg/internal/services"
)

func BuildServices() []di.Def {
	return []di.Def{
		{
			Name: constants.UserService,
			Build: func(ctn di.Container) (interface{}, error) {
				repo := ctn.Get(constants.UserRepository).(repositories.UserRepository)

				return services.NewUserService(repo), nil
			},
		},
	}
}
