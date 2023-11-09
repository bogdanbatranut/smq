package service

import (
	"smq/pkg/entity"
	"smq/pkg/init/config"
)

func InitServiceDependencies() (*entity.MessageQueue, string) {
	cfg := config.CreateConfig()
	port := cfg.GetString(config.AppPort)
	return entity.NewMessageQueue(), port
}
