export class AnthropicConfig {
  systemMessage?: string;
  maxAttachmentSize?: number;
  maxContextLength: number;

  anthropic: {
    apiKeys: string[];
    model: string;
    maxTokens: number;
    temperature?: number;
    topK?: number;
    topP?: number;
  };
}
