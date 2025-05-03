import { Entity, PrimaryGeneratedColumn, Column, ManyToOne } from 'typeorm';
import { Clinic } from '../clinics/clinic.entity';

@Entity()
export class Doctor {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  firstName: string;

  @Column()
  lastName: string;

  @Column()
  specialization: string;

  @ManyToOne(() => Clinic, clinic => clinic.doctors)
  clinic: Clinic;

  // другие поля
}