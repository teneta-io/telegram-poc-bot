package container

import (
	"context"
	"fmt"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"sync"
	"teneta-tg/internal/config"
	"teneta-tg/internal/constants"
	"teneta-tg/internal/translator"
	"teneta-tg/pkg/pgsql"
)

var (
	container di.Container
	once      sync.Once
)

func Build(ctx context.Context, wg *sync.WaitGroup) di.Container {
	once.Do(func() {
		builder, _ := di.NewBuilder()
		defs := []di.Def{
			{
				Name: constants.ConfigName,
				Build: func(ctn di.Container) (interface{}, error) {
					return config.New()
				},
			},
			{
				Name: constants.LoggerName,
				Build: func(ctn di.Container) (interface{}, error) {
					logger, err := zap.NewProduction()

					if err != nil {
						panic(fmt.Sprintf("can't initialize zap logger: %v", err))
					}
					zap.ReplaceGlobals(logger)

					return logger, nil
				},
			},
			{
				Name: constants.PgSQLConnection,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return pgsql.NewPgsqlConnection(cfg.PgSQLConfig)
				},
			},
			{
				Name: constants.Translator,
				Build: func(ctn di.Container) (interface{}, error) {
					return translator.NewTranslator(), nil
				},
			},
		}

		defs = append(defs, BuildRepositories()...)
		defs = append(defs, BuildServices()...)
		defs = append(defs, BuildBot(ctx, wg)...)

		if err := builder.Add(defs...); err != nil {
			panic(err)
		}

		container = builder.Build()
	})

	return container
}
