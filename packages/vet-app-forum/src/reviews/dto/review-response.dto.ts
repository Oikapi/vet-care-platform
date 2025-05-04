export class ReviewResponseDto {
    id: number;          // ID отзыва
    rating: number;      // Оценка 1-5
    comment?: string;    // Опциональный комментарий
    clinicId: number;    // ID клиники
    authorId: number;    // ID автора
    createdAt: Date;     // Дата создания
}