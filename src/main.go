package main

import (
	"anthropic-discord-bot/src/config"
	"anthropic-discord-bot/src/logger"
	"anthropic-discord-bot/src/modules/anthropic-api"
	"anthropic-discord-bot/src/modules/cache"
	"anthropic-discord-bot/src/modules/discord"
)

func main() {
	logger.Init(getLogLevel(), "App")
	logger.SetPrefix("AnthropicDiscordBot")
	logger.Info("Starting application...")

	cacheService, err := cache.Init()
	if err != nil {
		logger.Fatal("BootstrapError", err)
	}

	anthropicApiService, err := anthropicApi.Init()
	if err != nil {
		logger.Fatal("BootstrapError", err)
	}

	_, err = discord.Init(anthropicApiService, cacheService)
	if err != nil {
		logger.Fatal("BootstrapError", err)
	}

	logger.Info("Application started!")
	<-make(chan struct{})
}

func getLogLevel() logger.Level {
	switch config.Env.LogLevel {
	case "info", "log":
		return logger.InfoLevel
	case "warn":
		return logger.WarnLevel
	case "error":
		return logger.ErrorLevel
	case "debug", "verbose":
		return logger.DebugLevel
	default:
		return logger.InfoLevel
	}
}
