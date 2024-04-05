import { DiscordModule as NestjsDiscordModule } from '@discord-nestjs/core';
import { DynamicModule, Module } from '@nestjs/common';
import { GatewayIntentBits, Partials } from 'discord.js';

import { AnthropicModule } from '../anthropic';

import { DiscordUtilsService } from './discord-utils.service';
import { DiscordConfig } from './discord.config';
import { DiscordGateway } from './discord.gateway';
import { DiscordService } from './discord.service';

@Module({})
export class DiscordModule {
  static register({ botToken }: DiscordConfig): DynamicModule {
    return {
      module: DiscordModule,
      imports: [
        AnthropicModule.forFeature(),
        NestjsDiscordModule.forRootAsync({
          useFactory: () => ({
            token: botToken,
            discordClientOptions: {
              intents: [
                GatewayIntentBits.Guilds,
                GatewayIntentBits.GuildMessages,
                GatewayIntentBits.GuildIntegrations,
                GatewayIntentBits.DirectMessages,
                GatewayIntentBits.DirectMessageTyping,
                GatewayIntentBits.MessageContent,
              ],
              partials: [Partials.Channel, Partials.Channel, Partials.Reaction],
            },
            failOnLogin: true,
            autoLogin: true,
            shutdownOnAppDestroy: true,
          }),
        }),
      ],
      providers: [DiscordUtilsService, DiscordService, DiscordGateway],
    };
  }
}
