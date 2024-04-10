package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type environment struct {
	LogLevel string `env-default:"info" env:"LOG_LEVEL"`

	DiscordBotToken   string   `env-required:"true" env:"DISCORD_BOT_TOKEN"`
	MaxAttachmentSize uint32   `env-default:"1048576" env:"MAX_ATTACHMENT_SIZE"`
	MaxImageSize      uint32   `env-default:"10485760" env:"MAX_IMAGE_SIZE"`
	MaxContextSize    uint32   `env-default:"50000" env:"MAX_CONTEXT_SIZE"`
	DmWhitelist       []string `env:"DM_WHITELIST"`

	SystemMessage        string  `env-default:"You are bot" env:"SYSTEM_MESSAGE"`
	AnthropicApiKey      string  `env-required:"true" env:"ANTHROPIC_API_KEY"`
	AnthropicModel       string  `env-default:"claude-3-haiku-20240307" env:"ANTHROPIC_MODEL"`
	AnthropicMaxTokens   int64   `env-default:"4000" env:"ANTHROPIC_MAX_TOKENS"`
	AnthropicTemperature float32 `env-default:"1" env:"ANTHROPIC_TEMPERATURE"`
	AnthropicTopK        int     `env-default:"250" env:"ANTHROPIC_TOP_K"`
	AnthropicTopP        float32 `env-default:"1" env:"ANTHROPIC_TOP_P"`

	CacheDir    string `env-default:"./tmp" env:"CACHE_DIR"`
	CacheMaxAge string `env-default:"1h" env:"CACHE_MAX_AGE"`
}

var Env environment

func init() {
	err := cleanenv.ReadConfig(".env", &Env)

	if err != nil {
		if os.IsNotExist(err) {
			err = cleanenv.ReadEnv(&Env)

			if err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	}
}
