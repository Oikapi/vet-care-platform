import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Clinic } from 'src/clinics/clinic.entity';
import { Doctor } from 'src/doctors/doctors.entity';
import { User } from 'src/users/user.entity';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'mysql',
      host: process.env.DB_HOST,
      port: 3306,
      username: process.env.DB_USERNAME,
      password: process.env.DB_PASSWORD,
      database: process.env.DB_NAME,
      entities: [User, Clinic, Doctor],
      synchronize: true,
    }),
    TypeOrmModule.forFeature([User, Clinic, Doctor]),
  ],
  exports: [TypeOrmModule],
})
export class DatabaseModule {}
