import {
    Controller,
    Post,
    Body,
    Get,
    Param,
    BadRequestException
} from '@nestjs/common';
import { ReviewsService } from './reviews.service';
import { CreateReviewDto } from './dto/create-review.dto';

@Controller('clinics/:clinicId/reviews')
export class ReviewsController {
    constructor(private readonly reviewsService: ReviewsService) { }

    @Post()
    async create(
        @Param('clinicId') clinicId: number,
        @Body() createReviewDto: CreateReviewDto
    ) {
        try {
            return await this.reviewsService.createReview(
                createReviewDto.userId,
                {
                    clinicId,
                    rating: createReviewDto.rating,
                    comment: createReviewDto.comment
                }
            );
        } catch (error) {
            throw new BadRequestException(error.message);
        }
    }

    @Get()
    async getReviews(@Param('clinicId') clinicId: number) {
        try {
            return await this.reviewsService.getClinicReviews(clinicId);
        } catch (error) {
            throw new BadRequestException(error.message);
        }
    }
}