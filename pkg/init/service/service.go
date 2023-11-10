package service

import (
	"smq/pkg/entity"
	"smq/pkg/services"
	"smq/pkg/smqerrors"
)

func InitServiceDependencies() (*entity.MessageQueue, string) {
	cfg, err := services.NewViperConfig()
	smqerrors.Panic(err)
	port := cfg.GetString(services.AppPort)
	return entity.NewMessageQueue(), port
}
