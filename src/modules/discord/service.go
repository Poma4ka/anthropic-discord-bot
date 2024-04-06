package discord

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"

	"anthropic-discord-bot/src/logger"
	"anthropic-discord-bot/src/modules/anthropic-api"
)

type Service struct {
	Anthropic *anthropicApi.Service
	logger    *logger.Logger

	maxAttachmentSize uint32
	maxContextSize    uint32
}

func (s *Service) MessageCreate(client *discordgo.Session, message *discordgo.Message) (reply *discordgo.Message, err error) {
	currMessage := createAnthropicMessage(*s.logger, client, message, s.maxAttachmentSize)
	history, err := getMessagesHistory(*s.logger, client, message, s.maxAttachmentSize, s.maxContextSize)
	if err != nil {
		return
	}

	stopTyping := sendTyping(*s.logger, client, message.ChannelID)
	defer stopTyping()

	recv := make(chan anthropicApi.CompletionChunk, 1)

	go func() {
		err := s.Anthropic.CreateCompletionStream(context.Background(), currMessage, history, recv)
		if err != nil {
			s.logger.Error("Error create completion stream", err)
		}
	}()

	text := ""
	locked := false

	for chunk := range recv {
		text = *chunk.Text

		go func() {
			if !locked {
				locked = true
				reply, err = editReplyOrReply(client, reply, message, text)
				if err != nil {
					s.logger.Error("Error update message", err)
				}
				time.Sleep(100)
				locked = false
			}
		}()
	}

	reply, err = editReplyOrReply(client, reply, message, text)

	return
}
