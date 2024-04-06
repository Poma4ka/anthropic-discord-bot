package discord

import (
	"context"
	"encoding/base64"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/liushuangls/go-anthropic"

	"anthropic-discord-bot/src/logger"
	"anthropic-discord-bot/src/modules/anthropic-api"
	"anthropic-discord-bot/src/modules/cache"
)

type Service struct {
	Anthropic *anthropicApi.Service
	Cache     *cache.Service

	logger *logger.Logger

	maxAttachmentSize uint32
	maxImageSize      uint32
	maxContextSize    uint32
}

func (s *Service) MessageCreate(client *discordgo.Session, message *discordgo.Message) (reply *discordgo.Message, err error) {
	currMessage := s.createAnthropicMessage(client, message)
	history, err := s.getMessagesHistory(client, message)
	if err != nil {
		return
	}

	completionChannel := make(chan anthropicApi.CompletionChunk, 1)

	go func() {
		err := s.Anthropic.CreateCompletionStream(context.Background(), currMessage, history, completionChannel)
		if err != nil {
			s.logger.Error("Error create completion stream", err)
		}
	}()

	var text = ""
	var completed = atomic.Bool{}

	go func() {
		for chunk := range completionChannel {
			text = *chunk.Text

			if _, ok := <-completionChannel; !ok {
				completed.Store(true)
			}
		}
	}()

	for !completed.Load() {
		if text == "" {
			continue
		}

		editedReply, err := editReplyOrReply(client, reply, message, text)
		if err != nil {
			s.logger.Error("Error update message", err)
			continue
		}

		reply = editedReply
	}

	reply, err = editReplyOrReply(client, reply, message, text)

	return
}

func (s *Service) sendTyping(
	client *discordgo.Session,
	channelID string,
) func() {
	interval := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-interval.C:
				err := client.ChannelTyping(channelID)
				if err != nil {
					s.logger.Error("Failed send typing to channel "+channelID, err)
				}
			case <-done:
				interval.Stop()
				return
			}
		}
	}()

	return func() {
		done <- true
	}
}

func (s *Service) getMessagesHistory(
	client *discordgo.Session,
	message *discordgo.Message,
) (result []anthropic.Message, err error) {
	var currReference = message.ReferencedMessage
	var totalLength uint32 = 0

	for currReference != nil {
		message, err = client.ChannelMessage(currReference.ChannelID, currReference.ID)
		if err != nil {
			return
		}

		anthropicMessage := s.createAnthropicMessage(client, message)

		if totalLength += getAnthropicMessageLength(&anthropicMessage); totalLength >= s.maxContextSize {
			break
		}

		result = append([]anthropic.Message{anthropicMessage}, result...)
		currReference = message.ReferencedMessage
	}

	return
}

func (s *Service) createAnthropicMessage(
	client *discordgo.Session,
	message *discordgo.Message,
) (result anthropic.Message) {
	cleanMessage := message.ContentWithMentionsReplaced()

	var content []anthropic.MessageContent

	if cleanMessage != "" {
		content = append(content, anthropic.MessageContent{
			Type: "text",
			Text: &cleanMessage,
		})
	}

	for _, attachment := range message.Attachments {
		var isImage = isAttachmentImage(attachment)

		if isImage {
			if uint32(attachment.Size) > s.maxImageSize {
				continue
			}

			data, fromCache, err := s.getAttachmentData(attachment)
			if err != nil {
				continue
			}

			resizedImage, err := resizeImage(data, 1024)

			if err != nil {
				s.logger.Error("ResizeImageError", err)
				continue
			}

			content = append(content, anthropic.MessageContent{
				Type: anthropicApi.ContentTypeImage,
				Source: &anthropic.MessageContentImageSource{
					Type:      anthropicApi.SourceTypeBase64,
					MediaType: "image/jpeg",
					Data:      base64.StdEncoding.EncodeToString(resizedImage),
				},
			})

			if !fromCache {
				s.Cache.SaveAttachment(attachment.ID, &resizedImage)
			}
		} else {
			if uint32(attachment.Size) > s.maxAttachmentSize {
				continue
			}

			data, fromCache, err := s.getAttachmentData(attachment)
			if err != nil {
				continue
			}

			text := attachment.Filename + " (" + attachment.ContentType + "):\n\n" + string(data)

			content = append(content, anthropic.MessageContent{
				Type: anthropicApi.ContentTypeText,
				Text: &text,
			})

			if !fromCache {
				s.Cache.SaveAttachment(attachment.ID, &data)
			}
		}
	}

	return anthropic.Message{
		Role:    getMessageRole(client, message),
		Content: content,
	}
}

func (s *Service) getAttachmentData(attachment *discordgo.MessageAttachment) (data []byte, fromCache bool, err error) {
	if cached := s.Cache.GetAttachment(attachment.ID); cached != nil {
		fromCache = true
		data = *cached
	} else {
		data, err = downloadAttachment(attachment.URL)
		if err != nil {
			s.logger.Error("Failed download attachment "+attachment.ID, err)
		}
	}
	return
}
