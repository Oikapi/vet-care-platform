export class CreateReviewDto {
    clinicId: number;
    rating: number; // 1-5
    comment?: string;
}