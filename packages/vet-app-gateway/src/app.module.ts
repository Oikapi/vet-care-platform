import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { HttpModule } from '@nestjs/axios';
import { GatewayController } from './gateway/gateway.controller';
import { GatewayService } from './gateway/gateway.service';

@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true }), HttpModule],
  controllers: [GatewayController],
  providers: [GatewayService],
})
export class AppModule {}
