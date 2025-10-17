package app

import (
	"context"
	"github.com/ZhuchkovAA/loglib/internal/config"
	"github.com/ZhuchkovAA/loglib/internal/service"
)

type senderClient interface {
	WriteLog(context.Context, *service.LogRequest) (*service.LogResponse, error)
}

func Run(cfg config.Config, senderCli senderClient) (*service.Logger, error) {
	err := cfg.MustLoadConfig()
	if err != nil {
		return nil, err
	}

	sender := service.NewSender(cfg.ServiceName, senderCli)
	fallback := service.NewFallback(cfg.FallbackPath)

	logger := service.NewLogger(
		sender,
		fallback,
	)

	go logger.Run()

	reSender := service.NewReSender(cfg.FallbackPath, sender)
	go reSender.Loop()

	return logger, nil
}
