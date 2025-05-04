import { Module } from '@nestjs/common';
import { ReviewsService } from './reviews.service';
import { ReviewsController } from './reviews.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Review } from './entities/review.entity';

@Module({
  imports: [
    TypeOrmModule.forFeature([Review]) // Регистрируем сущность для TypeORM
  ],
  controllers: [ReviewsController],    // Подключаем контроллер
  providers: [ReviewsService],        // Подключаем сервис
  exports: [ReviewsService]          // Если сервис нужен в других модулях
})
export class ReviewsModule { }
