import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { PrismaService } from './prisma/prisma.service';
import { PrismaModule } from './prisma/prisma.module';
import { ExperimentsService } from './experiments/experiments.service';
import { CollectionsService } from './collections/collections.service';
import { CollectionsModule } from './collections/collections.module';
import { ExperimentsModule } from './experiments/experiments.module';

@Module({
  imports: [
    PrismaModule,
    CollectionsModule,
    ExperimentsModule
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
