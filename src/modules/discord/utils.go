package discord

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/liushuangls/go-anthropic"
	"github.com/nfnt/resize"

	anthropicApi "anthropic-discord-bot/src/modules/anthropic-api"
)

func editReplyOrReply(
	client *discordgo.Session,
	originalReply *discordgo.Message,
	message *discordgo.Message,
	content string,
) (reply *discordgo.Message, err error) {
	var isFile = utf8.RuneCountInString(content) > 2000

	var newReply *discordgo.Message
	var files []*discordgo.File

	if isFile {
		file := bytes.NewBufferString(content)

		files = []*discordgo.File{
			{
				Name:        "message.md",
				ContentType: "text/markdown",
				Reader:      file,
			},
		}

		content = ""
	}

	if originalReply != nil {
		newReply, err = client.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content:     &content,
			Files:       files,
			ID:          originalReply.ID,
			Channel:     message.ChannelID,
			Attachments: &[]*discordgo.MessageAttachment{},
		})
	} else {
		newReply, err = client.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:   content,
			Files:     files,
			Reference: message.Reference(),
		})
	}

	if err != nil {
		return originalReply, err
	}

	return newReply, err
}

func resizeImage(imgBuffer []byte, maxSize uint) (result []byte, err error) {
	img, _, err := image.Decode(bytes.NewReader(imgBuffer))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := uint(bounds.Max.X), uint(bounds.Max.Y)
	if width > maxSize || height > maxSize {
		var newWidth, newHeight uint
		if width > height {
			newWidth = maxSize
			newHeight = uint(float64(height) * float64(maxSize) / float64(width))
		} else {
			newHeight = maxSize
			newWidth = uint(float64(width) * float64(maxSize) / float64(height))
		}

		img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func downloadAttachment(
	url string,
) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(resp.Body)

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}

func getMessageRole(client *discordgo.Session, message *discordgo.Message) anthropicApi.MessageRole {
	if message.Author.ID == client.State.User.ID {
		return anthropicApi.MessageRoleAssistant
	}

	return anthropicApi.MessageRoleUser
}

func isAttachmentImage(attachment *discordgo.MessageAttachment) bool {
	var contentType = strings.Split(attachment.ContentType, "/")
	return len(contentType) == 2 && contentType[0] == "image"
}

func getAnthropicMessageLength(message *anthropic.Message) uint32 {
	length := 0
	for _, content := range message.Content {
		if content.IsTextContent() && content.Text != nil {
			length += utf8.RuneCountInString(*content.Text)
		}
	}
	return uint32(length)
}
