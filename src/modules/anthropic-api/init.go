package anthropicApi

import (
	"github.com/liushuangls/go-anthropic"

	"anthropic-discord-bot/src/config"
	"anthropic-discord-bot/src/logger"
)

func Init() (service *Service, err error) {
	log := logger.New("AnthropicModule")
	log.Info("Initializing module...")

	client := anthropic.NewClient(config.Env.AnthropicApiKey)

	service = &Service{
		logger:        log,
		client:        client,
		model:         config.Env.AnthropicModel,
		maxTokens:     int(config.Env.AnthropicMaxTokens),
		systemMessage: config.Env.SystemMessage,
		temperature:   &config.Env.AnthropicTemperature,
		topP:          &config.Env.AnthropicTopP,
		topK:          &config.Env.AnthropicTopK,
	}

	log.Info("Module initialized")
	return
}
