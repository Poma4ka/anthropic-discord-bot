package discord

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/liushuangls/go-anthropic"
	"github.com/nfnt/resize"

	"anthropic-discord-bot/src/logger"
	"anthropic-discord-bot/src/modules/anthropic-api"
)

func sendTyping(log logger.Logger, client *discordgo.Session, channelID string) func() {
	interval := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-interval.C:
				err := client.ChannelTyping(channelID)
				if err != nil {
					log.Error("Failed send typing to channel "+channelID, err)
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

func editReplyOrReply(
	client *discordgo.Session, originalReply *discordgo.Message,
	message *discordgo.Message,
	content string,
) (reply *discordgo.Message, err error) {
	if originalReply != nil {
		return client.ChannelMessageEdit(originalReply.ChannelID, originalReply.ID, content)
	} else {
		return client.ChannelMessageSendReply(message.ChannelID, content, message.Reference())
	}
}

func getMessagesHistory(
	log logger.Logger,
	client *discordgo.Session,
	message *discordgo.Message,
	maxAttachmentSize uint32,
	maxContextSize uint32,
) (result []anthropic.Message, err error) {
	var messages []*discordgo.Message

	currReference := message.ReferencedMessage

	// todo: validate contextSize
	for currReference != nil {
		message, err = client.ChannelMessage(currReference.ChannelID, currReference.ID)
		if err != nil {
			return
		}

		messages = append([]*discordgo.Message{message}, messages...)
		currReference = message.ReferencedMessage
	}

	result = make([]anthropic.Message, len(messages))

	for i, msg := range messages {
		result[i] = createAnthropicMessage(log, client, msg, maxAttachmentSize)
	}

	return
}

func createAnthropicMessage(
	log logger.Logger,
	client *discordgo.Session,
	message *discordgo.Message,
	maxAttachmentSize uint32,
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
		if uint32(attachment.Size) > maxAttachmentSize {
			continue
		}

		// todo: may be get from cache???
		data, err := downloadAttachment(log, attachment.URL)
		if err != nil {
			log.Error("Failed download attachment "+attachment.ID, err)
			continue
		}

		contentType := strings.Split(attachment.ContentType, "/")

		if len(contentType) == 2 && contentType[0] == "image" {
			resizedImage, err := resizeImage(data, 960)

			if err != nil {
				log.Error("ResizeImageError", err)
				continue
			}

			content = append(content, anthropic.MessageContent{
				Type: "image",
				Source: &anthropic.MessageContentImageSource{
					Type:      "base64",
					MediaType: "image/jpeg",
					Data:      base64.StdEncoding.EncodeToString(resizedImage),
				},
			})
		} else {
			text := attachment.Filename + " (" + attachment.ContentType + ")"

			text = text + ":\n\n" + string(data)

			content = append(content, anthropic.MessageContent{
				Type: "text",
				Text: &text,
			})
		}
	}

	role := anthropicApi.MessageRoleUser

	if message.Author.ID == client.State.User.ID {
		role = anthropicApi.MessageRoleAssistant
	}

	return anthropic.Message{
		Role:    role,
		Content: content,
	}
}

func resizeImage(imgBuffer []byte, maxSize uint) (result []byte, err error) {
	img, _, err := image.Decode(bytes.NewReader(imgBuffer))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := uint(bounds.Max.X), uint(bounds.Max.Y)
	var newWidth, newHeight uint
	if width > height {
		newWidth = maxSize
		newHeight = uint(float64(height) * float64(maxSize) / float64(width))
	} else {
		newHeight = maxSize
		newWidth = uint(float64(width) * float64(maxSize) / float64(height))
	}

	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 960})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func downloadAttachment(
	log logger.Logger,
	url string,
) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Failed close HTTP request", err)
		}
	}(resp.Body)

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}
