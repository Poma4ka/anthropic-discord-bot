import { MessageParam } from '@anthropic-ai/sdk/resources';
import { Injectable } from '@nestjs/common';

import { CompletionMessage } from './dto/common';
import { MessageRoleEnum } from './dto/enum';

@Injectable()
export class AnthropicUtilsService {
  parseMessage(message: CompletionMessage): MessageParam {
    let attachments: string = '';

    if (message.attachments?.length) {
      for (const attachment of message.attachments) {
        attachments = `${attachments}\n\n\n==========File ${attachment.name} start==========\n${attachment.content.toString()}\n==========File ${attachment.name} end==========`;
      }
    }

    return {
      role: message.role,
      content: `${message.content}${attachments}`,
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
