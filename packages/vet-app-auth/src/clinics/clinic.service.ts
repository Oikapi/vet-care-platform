import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Clinic } from './clinic.entity';
import * as bcrypt from 'bcrypt';
import { CreateClinicDto } from './dto/create-clinic.dto';
import { UpdateClinicDto } from './dto/update-clinic.dto';

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

  async findAll(): Promise<Clinic[]> {
    return this.clinicRepository.find();
  }

  async findOne(id: number): Promise<Clinic> {
    const clinic = await this.clinicRepository.findOne({ where: { id } });
    if (!clinic) {
      throw new NotFoundException(`Clinic with id ${id} not found`);
    }
    return clinic;
  }

  async update(id: number, updateClinicDto: UpdateClinicDto): Promise<Clinic> {
    const clinic = await this.findOne(id);

    if (updateClinicDto.password) {
      updateClinicDto.password = await bcrypt.hash(
        updateClinicDto.password,
        10
      );
    }

    const updated = Object.assign(clinic, updateClinicDto);
    return this.clinicRepository.save(updated);
  }

  async remove(id: number): Promise<void> {
    const result = await this.clinicRepository.delete(id);
    if (result.affected === 0) {
      throw new NotFoundException(`Clinic with id ${id} not found`);
    }
  }
}
