import { Injectable, BadRequestException } from '@nestjs/common'; // Добавили импорт
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Review } from './entities/review.entity';

@Injectable()
export class ReviewsService {
    constructor(
        @InjectRepository(Review)
        private readonly repo: Repository<Review>,
    ) { }

    async createReview(userId: number, data: {
        clinicId: number;
        rating: number;
        comment?: string
    }) {
        // Проверка входных данных
        if (!userId || userId <= 0) {
            throw new BadRequestException('Неверный ID пользователя');
        }

        if (data.rating < 1 || data.rating > 5) {
            throw new BadRequestException('Рейтинг должен быть от 1 до 5');
        }

        const review = this.repo.create({
            authorId: userId,
            clinicId: data.clinicId,
            rating: data.rating,
            comment: data.comment
        });

        try {
            return await this.repo.save(review);
        } catch (error) {
            throw new BadRequestException('Ошибка при сохранении отзыва');
        }
    }

    async getClinicReviews(clinicId: number) {
        if (!clinicId || clinicId <= 0) {
            throw new BadRequestException('Неверный ID клиники');
        }

        return this.repo.find({
            where: { clinicId },
            order: { createdAt: 'DESC' }
        });
    }
}