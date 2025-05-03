import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Clinic } from './clinic.entity';
import { ClinicService } from './clinic.service';
import { ClinicController } from './clinic.controller';

@Module({
  providers: [ClinicService],
  controllers: [ClinicController],
  imports: [TypeOrmModule.forFeature([Clinic])],
  exports: [ClinicService],
})
export class ClinicsModule {}
