import { Module } from '@nestjs/common';

import { AnthropicModule } from './modules/anthropic';
import { DiscordModule } from './modules/discord';

@Module({
  imports: [
    AnthropicModule.forRoot({
      systemMessage: process.env.SYSTEM_MESSAGE,
      maxAttachmentSize: process.env.MAX_ATTACHMENT_SIZE
        ? Number(process.env.MAX_ATTACHMENT_SIZE)
        : undefined,
      maxContextLength: Number(process.env.ANTHROPIC_MAX_CONTEXT_LENGTH),
      anthropic: {
        apiKeys: process.env.ANTHROPIC_API_KEY.split(','),
        model: process.env.ANTHROPIC_MODEL,
        maxTokens: Number(process.env.ANTHROPIC_MAX_TOKENS),
        temperature: process.env.ANTHROPIC_TEMPERATURE
          ? Number(process.env.ANTHROPIC_TEMPERATURE)
          : undefined,
        topK: process.env.ANTHROPIC_TOP_K ? Number(process.env.ANTHROPIC_TOP_K) : undefined,
        topP: process.env.ANTHROPIC_TOP_P ? Number(process.env.ANTHROPIC_TOP_P) : undefined,
      },
    }),
    DiscordModule.register({
      botToken: process.env.DISCORD_BOT_TOKEN,
    }),
  ],
})
export class AppModule {}
