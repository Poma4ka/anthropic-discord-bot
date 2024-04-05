import { InjectDiscordClient, On } from '@discord-nestjs/core';
import { Inject, Injectable } from '@nestjs/common';
import { Client, ClientUser, Message } from 'discord.js';

import { DiscordService } from './discord.service';

@Injectable()
export class DiscordGateway {
  constructor(
    @InjectDiscordClient()
    private readonly client: Client,
    @Inject(DiscordService)
    private discordBotService: DiscordService,
  ) {}

  @On('messageCreate')
  async onMessageCreate(message: Message) {
    if (message.system) {
      return;
    }

    if (message.guildId === null) {
      // todo: may be add support
      return;
    }

    if (
      !message.mentions.has((this.client.user as ClientUser).id, {
        ignoreEveryone: true,
        ignoreRoles: true,
        ignoreRepliedUser: false,
      })
    ) {
      return;
    }

    await this.discordBotService.createMessage(message);
  }

  @On('messageUpdate')
  async onMessageUpdate(message: Message) {
    await this.discordBotService.updateMessage(message);
  }

  @On('messageDelete')
  async onMessageDelete(message: Message) {
    await this.discordBotService.deleteMessage(message);
  }
}
