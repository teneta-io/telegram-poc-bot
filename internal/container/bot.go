package container

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
	"sync"
	"teneta-tg/internal/bot"
	"teneta-tg/internal/config"
	"teneta-tg/internal/constants"
	"teneta-tg/internal/services"
	"teneta-tg/internal/translator"
)

func BuildBot(ctx context.Context, wg *sync.WaitGroup) []di.Def {
	return []di.Def{
		{
			Name: constants.TelegramAPI,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				api, err := tgbotapi.NewBotAPI(cfg.TelegramConfig.Token)
				if err != nil {
					return nil, err
				}

				api.Debug = true

				return api, nil
			},
		},
		{
			Name: constants.Bot,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				api := ctn.Get(constants.TelegramAPI).(*tgbotapi.BotAPI)
				userService := ctn.Get(constants.UserService).(services.UserService)
				t := ctn.Get(constants.Translator).(*translator.Translator)

				return bot.New(ctx, wg, cfg.TelegramConfig, api, userService, t)
			},
		},
	}
}
