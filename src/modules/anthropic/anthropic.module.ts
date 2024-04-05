import { DynamicModule, Module, Provider } from '@nestjs/common';

import { AnthropicUtilsService } from './anthropic-utils.service';
import { AnthropicConfig } from './anthropic.config';
import { AnthropicService } from './anthropic.service';

@Module({})
export class AnthropicModule {
  private static configProvider: Provider;

  static forRoot(config: AnthropicConfig): DynamicModule {
    this.configProvider = {
      provide: AnthropicConfig,
      useValue: config,
    };

    return {
      module: AnthropicModule,
      imports: [],
      providers: [],
      exports: [],
    };
  }

  static forFeature(): DynamicModule {
    return {
      module: AnthropicModule,
      imports: [],
      providers: [AnthropicService, AnthropicUtilsService, this.configProvider],
      exports: [AnthropicService],
    };
  }
}
