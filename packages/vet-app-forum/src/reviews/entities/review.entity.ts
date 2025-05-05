import { Entity, PrimaryGeneratedColumn, Column } from 'typeorm';

@Entity('reviews')
export class Review {
    @PrimaryGeneratedColumn()
    id!: number;

    @Column({ name: 'clinic_id' })
    clinicId: number;

    @Column('int', { name: 'user_id' })
    authorId!: number;

    @Column('int')
    rating!: number;

    @Column('text', { nullable: true })
    comment?: string;

    @Column('timestamp', {
        name: 'created_at',
        default: () => 'CURRENT_TIMESTAMP'
    })
    createdAt!: Date;

    @Column('timestamp', {
        name: 'updated_at',
        nullable: true,
        onUpdate: 'CURRENT_TIMESTAMP' // Автообновление
    })
    updatedAt?: Date;
}