import { Anthropic, APIError } from '@anthropic-ai/sdk';
import { AnthropicError } from '@anthropic-ai/sdk/error';
import { MessageParam } from '@anthropic-ai/sdk/resources';
import { Inject, Injectable, Logger } from '@nestjs/common';
import { Observable, Subject } from 'rxjs';

import { AppError } from '../../common/errors';

import { AnthropicUtilsService } from './anthropic-utils.service';
import { AnthropicConfig } from './anthropic.config';
import { CompletionMessage } from './dto/common';
import { CreateCompletionOptionsDto, CreateCompletionResultDto } from './dto/internal';

@Injectable()
export class AnthropicService {
  private logger = new Logger(this.constructor.name);
  private client: Anthropic;

  constructor(
    @Inject(AnthropicConfig)
    private config: AnthropicConfig,
    @Inject(AnthropicUtilsService)
    private anthropicUtilsService: AnthropicUtilsService,
  ) {
    this.client = this.createClient();
  }

  validateAttachment(
    size: number,
    contentType: string = 'application/octet-stream',
    name: string = '',
  ): boolean {
    if (size > (this.config.maxAttachmentSize ?? 0)) {
      return false;
    }

    if (name.length > 100) {
      return false;
    }

    const [type] = contentType.split('/');

    if (['application', 'text', 'image'].includes(type)) {
      return true;
    }

    return false;
  }

  async createCompletion({
    message,
    signal,
    getPreviousMessage,
  }: CreateCompletionOptionsDto): Promise<Observable<CreateCompletionResultDto>> {
    const messages = await this.prepareMessages(message, getPreviousMessage);

    const stream = this.client.messages.stream(
      {
        model: this.config.anthropic.model,
        max_tokens: this.config.anthropic.maxTokens,
        temperature: this.config.anthropic.temperature,
        top_k: this.config.anthropic.topK,
        top_p: this.config.anthropic.topP,
        system: this.config.systemMessage,
        messages,
      },
      {
        signal,
      },
    );

    const subject = new Subject<CreateCompletionResultDto>();

    stream.on('text', (chunk) =>
      subject.next({
        chunk,
      }),
    );

    stream.on('end', () => subject.complete());

    stream.on('abort', (error) => subject.error(this.handleError(error)));
    stream.on('error', (error) => subject.error(this.handleError(error)));

    return subject.asObservable();
  }

  private handleError(error: AnthropicError): AppError {
    if (error instanceof APIError) {
      switch (error.status) {
        case 429: {
          this.logger.warn('429 error: swapping api key...');
          this.swapKey();
          break;
        }
        case 401: {
          this.logger.warn('401 error: Removing current api key...');

          this.removeCurrentKey();
          break;
        }
        default: {
          this.logger.error(error.message, error.stack ?? error.cause);
        }
      }
    } else {
      this.logger.error(error.message, error.stack ?? error.cause);
    }

    return new AppError('Произошла ошибка при запросе к Anthropic API');
  }

  private async prepareMessages(
    message: CompletionMessage,
    getPreviousMessage?: CreateCompletionOptionsDto['getPreviousMessage'],
  ): Promise<MessageParam[]> {
    const result: MessageParam[] = [];

    const parsedMessage = this.anthropicUtilsService.parseMessage(message);

    const messageLength = this.anthropicUtilsService.getMessageLength(parsedMessage);

    while (true) {
      const previousMessage = await getPreviousMessage?.().catch(() => null);

      if (!previousMessage) {
        break;
      }

      const parsedMessage = this.anthropicUtilsService.parseMessage(previousMessage);

      const previousMessageLength = this.anthropicUtilsService.getMessageLength(parsedMessage);

      if (messageLength + previousMessageLength >= this.config.maxContextLength) {
        break;
      }

      result.unshift(parsedMessage);
    }

    result.push(parsedMessage);

    return result;
  }

  private apiKey?: string;

  private createClient(): Anthropic {
    if (this.config.anthropic.apiKeys.length === 0) {
      this.logger.error('Api keys not found');
    }

    if (!this.apiKey) {
      this.apiKey = this.config.anthropic.apiKeys[0];
    }

    return new Anthropic({
      apiKey: `${this.apiKey}`,
      maxRetries: 10,
      timeout: 30000,
    });
  }

  private removeCurrentKey() {
    this.config.anthropic.apiKeys = this.config.anthropic.apiKeys.filter(
      (key) => key !== this.client.apiKey,
    );
    this.client = this.createClient();
  }

  private swapKey() {
    if (this.config.anthropic.apiKeys.length === 0) {
      return;
    }

    const index = this.config.anthropic.apiKeys.indexOf(this.client.apiKey as string);
    if (index === -1) {
      this.apiKey = this.config.anthropic.apiKeys[0];
    } else {
      this.apiKey = this.config.anthropic.apiKeys[index + 1] ?? this.config.anthropic.apiKeys[0];
    }

    this.client = this.createClient();
  }
}
