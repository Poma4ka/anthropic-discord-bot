package anthropicApi

import (
	"context"
	"errors"

	"anthropic-discord-bot/src/logger"

	"github.com/liushuangls/go-anthropic"
)

type Service struct {
	logger *logger.Logger
	client *anthropic.Client

	model         string
	maxTokens     int
	systemMessage string
	temperature   *float32
	topP          *float32
	topK          *int
}

func (s *Service) CreateCompletionStream(ctx context.Context, message anthropic.Message, history []anthropic.Message, send chan<- CompletionChunk) (err error) {
	defer close(send)

	messageCuted := getMessageText(&message)

	s.logger.Debug("Creating completion for \"" + messageCuted + "\"...")

	if len(history) > 1 && history[0].Role == MessageRoleAssistant {
		history = history[1:]
	}

	var text string

	request := anthropic.MessagesStreamRequest{
		MessagesRequest: anthropic.MessagesRequest{
			Model:       s.model,
			Messages:    append(history, message),
			MaxTokens:   s.maxTokens,
			System:      s.systemMessage,
			Stream:      true,
			Temperature: s.temperature,
			TopP:        s.topP,
			TopK:        s.topK,
		},
		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
			if data.Delta.Type == "text_delta" {
				text = text + data.Delta.Text

				send <- CompletionChunk{
					Text:  &text,
					Delta: &data.Delta.Text,
				}
			}
		},
	}

	_, err = s.client.CreateMessagesStream(ctx, request)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			s.logger.Warn("Completion stream canceled")
			return err
		}
		s.logger.Error("Create message stream error", err, message, history)
		return err
	}

	s.logger.Debug("Completion created successfully for \"" + messageCuted + "\"")

	return
}
