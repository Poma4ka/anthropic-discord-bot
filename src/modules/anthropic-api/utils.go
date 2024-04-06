package anthropicApi

import (
	"github.com/liushuangls/go-anthropic"
)

func getMessageText(message *anthropic.Message) string {
	firstContent := message.GetFirstContent()

	if firstContent.IsTextContent() {
		if len(firstContent.GetText()) < 100 {
			return firstContent.GetText()
		}
		return firstContent.GetText()[:97] + "..."
	}

	if firstContent.IsImageContent() {
		return "Image message"
	}

	return "Unknown message"
}
