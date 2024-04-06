package cache

import (
	"anthropic-discord-bot/src/config"
	"anthropic-discord-bot/src/logger"
)

func Init() (service *Service, err error) {
	log := logger.New("CacheModule")
	log.Info("Initializing module...")

	service = &Service{
		logger:      log,
		cacheDir:    config.Env.CacheDir,
		cacheMaxAge: parseDuration(config.Env.CacheMaxAge),
	}

	err = mkdir(service.cacheDir)
	if err != nil {
		log.Error("Error create cache directory", err)
		return
	}

	go service.startCacheCleanup()

	log.Info("Module initialized")
	return
}
