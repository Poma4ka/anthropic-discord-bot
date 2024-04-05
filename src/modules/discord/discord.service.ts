import { InjectDiscordClient } from '@discord-nestjs/core';
import { Inject, Injectable, Logger } from '@nestjs/common';
import axios from 'axios';
import { Client, Message } from 'discord.js';

import {
  AnthropicService,
  CompletionAttachment,
  CompletionMessage,
  GetPreviousMessage,
  MessageRoleEnum,
} from '../anthropic';

import { DiscordUtilsService } from './discord-utils.service';

interface ProcessedMessage {
  abortController: AbortController;
  reply: Message | null;
}

@Injectable()
export class DiscordService {
  private readonly logger = new Logger(DiscordService.name);

  private readonly processedMessages: Map<string, ProcessedMessage> = new Map();

  constructor(
    @Inject(DiscordUtilsService)
    private discordUtilsService: DiscordUtilsService,
    @Inject(AnthropicService)
    private anthropicService: AnthropicService,
    @InjectDiscordClient()
    private readonly client: Client,
  ) {}

  async createMessage(message: Message): Promise<void> {
    const abortController = new AbortController();

    const abortTyping = this.discordUtilsService.sendTyping(message.channel);

    let processedMessage: ProcessedMessage = {
      abortController,
      reply: null,
    };

    if (this.processedMessages.has(message.id)) {
      processedMessage = {
        ...(this.processedMessages.get(message.id) as ProcessedMessage),
        abortController,
      };
    }

    this.processedMessages.set(message.id, processedMessage);

    let reply: Message | null = processedMessage.reply;

    try {
      const completionMessage = await this.getCompletionMessage(message);

      const completion = await this.anthropicService.createCompletion({
        signal: abortController.signal,
        message: completionMessage,
        getPreviousMessage: this.getPreviousMessage(message),
      });

      let content = '';
      let isReplying: boolean = false;

      await completion.forEach((value) => {
        if (abortController.signal.aborted) {
          return;
        }

        content = `${content}${value.chunk}`;

        if (content) {
          if (!isReplying) {
            isReplying = true;

            this.discordUtilsService
              .editOrReplyMessage(message, content, reply ?? undefined)
              .then((message) => {
                if (message) {
                  reply = message;
                }

                isReplying = false;
              });
          }
        }
      });

      if (abortController.signal.aborted) {
        return;
      }

      await this.discordUtilsService.editOrReplyMessage(message, content, reply ?? undefined);
    } catch (error) {
      this.logger.error(error);

      await this.discordUtilsService
        .editOrReplyMessage(
          message,
          'Что-то я затупил, может быть пора отдохнуть 😞',
          reply ?? undefined,
        )
        .catch((error) => this.logger.error(error));
    } finally {
      abortTyping();

      if (!abortController.signal.aborted) {
        this.processedMessages.delete(message.id);
      }
    }
  }

  async updateMessage(message: Message): Promise<void> {
    if (this.processedMessages.has(message.id)) {
      try {
        this.processedMessages.get(message.id)?.abortController.abort();
      } catch (e) {}

      await this.createMessage(message);
    }
  }

  async deleteMessage(message: Message): Promise<void> {
    if (this.processedMessages.has(message.id)) {
      try {
        this.processedMessages.get(message.id)?.abortController.abort();
      } catch (e) {}

      this.processedMessages.delete(message.id);
    }
  }

  private getPreviousMessage(message: Message): GetPreviousMessage {
    let currMessage: Message = message;

    return async () => {
      if (!currMessage.reference) {
        return null;
      }

      currMessage = await currMessage.fetchReference();

      return await this.getCompletionMessage(currMessage);
    };
  }

  private async getCompletionMessage(message: Message): Promise<CompletionMessage> {
    const content = `${message.cleanContent}`;

    const attachments: CompletionAttachment[] = [];

    for (const [, attachment] of message.attachments) {
      try {
        if (
          !this.anthropicService.validateAttachment(
            attachment.size,
            attachment.contentType ?? undefined,
            attachment.name,
          )
        ) {
          continue;
        }

        const content = await axios
          .get(attachment.url, {
            responseType: 'text',
          })
          .then((r) => r.data);

        attachments.push({
          content,
          name: attachment.name,
          contentType: attachment.contentType ?? undefined,
        });
      } catch (error) {
        this.logger.error(error);
      }
    }

    return {
      content,
      attachments,
      role:
        message.author.id === this.client.user?.id
          ? MessageRoleEnum.ASSISTANT
          : MessageRoleEnum.USER,
    };
  }
}
