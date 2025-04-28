import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Clinic } from './clinic.entity';
import * as bcrypt from 'bcrypt';
import { CreateClinicDto } from './dto/create-clinic.dto';
// import { CreateClinicDto } from './dto/create-clinic.dto';

@Injectable()
export class ClinicService {
  constructor(
    @InjectRepository(Clinic)
    private clinicRepository: Repository<Clinic>
  ) {}

  async findByEmail(email: string): Promise<Clinic | null> {
    return this.clinicRepository.findOne({ where: { email } });
  }

  async create(createClinicDto: CreateClinicDto): Promise<Clinic> {
    const hashedPassword = await bcrypt.hash(createClinicDto.password, 10);
    const clinic = this.clinicRepository.create({
      ...createClinicDto,
      password: hashedPassword,
    });
    return this.clinicRepository.save(clinic);
  }
}
