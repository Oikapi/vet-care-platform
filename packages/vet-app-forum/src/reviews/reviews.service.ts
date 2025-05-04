import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Review } from './entities/review.entity';
import { ReviewResponseDto } from './dto/review-response.dto';

@Injectable()
export class ReviewsService {
    constructor(
        @InjectRepository(Review)
        private readonly repo: Repository<Review>,
    ) { }

    async getReviews(clinicId: number): Promise<ReviewResponseDto[]> {
        const reviews = await this.repo.find({ where: { clinicId } });
        return reviews.map(review => ({
            id: review.id,
            rating: review.rating,
            comment: review.comment,
            clinicId: review.clinicId,
            authorId: review.authorId,
            createdAt: review.createdAt
        }));
    }
}