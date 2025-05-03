import { Doctor } from 'src/doctors/doctors.entity';
import { Entity, PrimaryGeneratedColumn, Column, OneToMany } from 'typeorm';

@Entity()
export class Clinic {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name: string;

  @Column()
  phone: string;

  @Column()
  email: string;

  @Column()
  password: string;

  @Column({ nullable: true })
  photo?: string;

  @OneToMany(() => Doctor, (doctor) => doctor.clinic)
  doctors: Doctor[];
}
