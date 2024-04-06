import { ImageBlockParam, MessageParam, TextBlockParam } from '@anthropic-ai/sdk/resources';
import { Injectable } from '@nestjs/common';

import { CompletionMessage } from './dto/common';
import { MessageRoleEnum } from './dto/enum';

@Injectable()
export class AnthropicUtilsService {
  parseMessage(message: CompletionMessage): MessageParam {
    const content: Array<TextBlockParam | ImageBlockParam> = [];

    if (message.content) {
      content.push({
        type: 'text',
        text: message.content,
      });
    }

    if (message.attachments?.length) {
      for (const attachment of message.attachments) {
        if (attachment.contentType?.split('/').at(0) === 'image') {
          content.push({
            type: 'image',
            source: {
              type: 'base64',
              media_type: attachment.contentType as 'image/jpeg',
              data: Buffer.from(attachment.content).toString('base64'),
            },
          });
        } else {
          content.push({
            type: 'text',
            text: `${attachment.name}${attachment.contentType ? ` ${attachment.contentType}` : ''}:\n\n${attachment.content.toString()}`,
          });
        }
      }
    }

    return {
      role: this.mapRole(message.role),
      content,
    };
  }

  getMessageLength(message: MessageParam) {
    if (typeof message.content === 'string') {
      return message.content.length;
    }

    return message.content.reduce((acc, content) => {
      if (content.type === 'text') {
        return acc + content.text.length;
      }

      if (content.type === 'image') {
        return (
          acc +
          content.source.data.length +
          content.source.media_type.length +
          content.source.type.length
        );
      }
      return acc;
    }, 0);
  }

  mapRole(role: MessageRoleEnum): MessageParam['role'] {
    return (
      {
        [MessageRoleEnum.USER]: 'user' as const,
        [MessageRoleEnum.ASSISTANT]: 'assistant' as const,
      }[role] ?? ('user' as const)
    );
  }
}
