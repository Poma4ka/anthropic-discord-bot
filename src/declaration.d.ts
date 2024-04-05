declare global {
  namespace NodeJS {
    interface ProcessEnv {
      NODE_ENV: 'development' | 'production';
      LOG_LEVEL: string;

      DISCORD_BOT_TOKEN: string;

      SYSTEM_MESSAGE?: string;
      MAX_ATTACHMENT_SIZE?: string;

      ANTHROPIC_API_KEY: string;
      ANTHROPIC_MODEL: string;
      ANTHROPIC_MAX_TOKENS: string;
      ANTHROPIC_MAX_CONTEXT_LENGTH: string;
      ANTHROPIC_TEMPERATURE?: string;
      ANTHROPIC_TOP_K?: string;
      ANTHROPIC_TOP_P?: string;
    }
  }
}

export {};
