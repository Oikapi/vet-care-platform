import { Controller, Get, Param } from '@nestjs/common';
import { ReviewsService } from './reviews.service';
import { ReviewResponseDto } from './dto/review-response.dto';

@Controller('clinics/:clinicId/reviews')
export class ReviewsController {
    constructor(private readonly reviewsService: ReviewsService) { }

    @Get()
    async getReviews(
        @Param('clinicId') clinicId: number,
    ): Promise<ReviewResponseDto[]> {
        console.log(123)
        return this.reviewsService.getReviews(clinicId);
    }
}