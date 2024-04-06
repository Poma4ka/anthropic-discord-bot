package anthropicApi

type MessageRole = string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type CompletionChunk struct {
	Text  *string
	Delta *string
}
