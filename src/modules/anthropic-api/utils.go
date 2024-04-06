package anthropicApi

import (
	"unicode/utf8"

	"github.com/liushuangls/go-anthropic"
)

func getMessageText(message *anthropic.Message) string {
	firstContent := message.GetFirstContent()
	text := firstContent.GetText()

	if firstContent.IsTextContent() {
		if utf8.RuneCountInString(text) < 100 {
			return firstContent.GetText()
		}
		return string([]rune(text)[:97]) + "..."
	}

	if firstContent.IsImageContent() {
		return "Image message"
	}

	return "Unknown message"
}
