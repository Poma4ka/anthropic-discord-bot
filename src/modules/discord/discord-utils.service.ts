import { Injectable, Logger } from '@nestjs/common';
import { AttachmentBuilder, BaseMessageOptions, Message, TextBasedChannel } from 'discord.js';

@Injectable()
export class DiscordUtilsService {
  private readonly logger = new Logger(this.constructor.name);

  async createTextAttachment(content: string, name: string): Promise<AttachmentBuilder> {
    const attachment = new AttachmentBuilder(Buffer.from(content));
    attachment.setName(name);
    return attachment;
  }

  sendTyping(channel: TextBasedChannel): () => void {
    channel.sendTyping().catch(() => null);
    const interval = setInterval(() => {
      channel.sendTyping().catch(() => null);
    }, 10000);

    return () => clearInterval(interval);
  }

  getMessageAuthorName(message: Message): string {
    const defaultNick: string = 'Unknown Alien';

    try {
      const user = message.guild?.members.cache.get(message.author.id);

      return user?.nickname || message.author.displayName || message.author.username;
    } catch (error: unknown) {
      this.logger.error(error);
      return defaultNick;
    }
  }

  async editOrReplyMessage(
    message: Message,
    content: string,
    reply?: Message,
  ): Promise<Message | null> {
    const payload: BaseMessageOptions =
      content.length > 2000
        ? {
            files: [await this.createTextAttachment(content, 'message.txt')],
            content: '',
          }
        : {
            content,
          };

    if (content) {
      if (reply) {
        return await reply.edit(payload);
      } else {
        return await message.reply(payload);
      }
    }

    return null;
  }
}
