import { Module } from '@nestjs/common';
import { DoctorService } from './doctors.service';
import { DoctorController } from './doctors.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Doctor } from './doctors.entity';

@Module({
  providers: [DoctorService],
  controllers: [DoctorController],
  imports: [TypeOrmModule.forFeature([Doctor])],
  exports: [DoctorService],
})
export class DoctorsModule {}
