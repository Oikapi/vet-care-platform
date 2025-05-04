import { Entity, PrimaryGeneratedColumn, Column } from 'typeorm';

@Entity('reviews')
export class Review {
    @PrimaryGeneratedColumn()
    id!: number;

    @Column('int')
    clinicId!: number;

    @Column('int')
    authorId!: number;

    @Column('int')
    rating!: number;

    @Column('text', { nullable: true })
    comment?: string;

    @Column('timestamp', { default: () => 'CURRENT_TIMESTAMP' })
    createdAt!: Date;
}