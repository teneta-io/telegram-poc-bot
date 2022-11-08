package main

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"teneta-tg/internal/bot"
	"teneta-tg/internal/constants"
	"teneta-tg/internal/container"
	"teneta-tg/internal/translator"
	"teneta-tg/utils"
	"time"
)

func main() {
	now := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	app := container.Build(ctx, wg)
	_ = app.Get(constants.LoggerName).(*zap.Logger)
	_ = app.Get(constants.Translator).(*translator.Translator)
	b := app.Get(constants.Bot).(*bot.Bot)

	zap.S().Info("Starting application...")
	zap.S().Infof("Up and running (%s)", time.Since(now))

	go b.Run()

	zap.S().Infof("Up and running (%s)", time.Since(now))
	zap.S().Infof("Got %s signal. Shutting down...", <-utils.WaitTermSignal())

	cancel()
	wg.Wait()

	zap.S().Info("Service stopped.")
}
