package container

import (
	"github.com/sarulabs/di"
	"gorm.io/gorm"
	"teneta-tg/internal/constants"
	"teneta-tg/internal/repositories/pgsql"
)

func BuildRepositories() []di.Def {
	return []di.Def{
		{
			Name: constants.UserRepository,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnection).(*gorm.DB)

				return pgsql.NewUserRepository(conn), nil
			},
		},
	}
}
