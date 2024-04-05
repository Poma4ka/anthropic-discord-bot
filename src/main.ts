import 'dotenv/config';

import { NestFactory } from '@nestjs/core';

import { AppModule } from './app.module';
import { AppLogger } from './common/loggers';

const bootstrap = async () => {
  const logger = new AppLogger();

  const app = await NestFactory.createApplicationContext(AppModule, {
    logger,
    abortOnError: true,
  });

  return app.init();
};

void bootstrap();
