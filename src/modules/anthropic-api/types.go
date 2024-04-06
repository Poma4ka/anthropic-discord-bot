package anthropicApi

type MessageRole = string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type ContentType = string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
)

type SourceType = string

const (
	SourceTypeBase64 SourceType = "base64"
)

type CompletionChunk struct {
	Text  *string
	Delta *string
}
