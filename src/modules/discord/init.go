package discord

import (
	"github.com/bwmarrin/discordgo"

	"anthropic-discord-bot/src/config"
	"anthropic-discord-bot/src/logger"
	"anthropic-discord-bot/src/modules/anthropic-api"
)

func Init(anthropic *anthropicApi.Service) (service *Service, err error) {
	log := logger.New("DiscordModule")
	log.Info("Initializing module...")

	log.Info("Starting discord bot...")

	client, err := discordgo.New("Bot " + config.Env.DiscordBotToken)
	if err != nil {
		log.Error("Error create discord bot", err)
		return
	}

	service = &Service{
		Anthropic:         anthropic,
		logger:            log,
		maxAttachmentSize: config.Env.MaxAttachmentSize,
		maxContextSize:    config.Env.MaxContextSize,
	}

	controller := Controller{
		Service: service,
	}

	client.AddHandler(controller.messageCreate)
	// todo:
	//client.AddHandler(controller.messageUpdate)
	//client.AddHandler(controller.messageDelete)

	err = client.Open()
	if err != nil {
		log.Error("Error start discord bot", err)
		return
	}

	log.Info("Discord bot started")

	log.Info("Module initialized")
	return
}