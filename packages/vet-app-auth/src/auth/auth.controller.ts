// auth.controller.ts
import { Controller, Post, Body, UseGuards, Request } from '@nestjs/common';
import { AuthService } from './auth.service';
import { ClinicService } from 'src/clinics/clinic.service';
import { UserService } from 'src/users/users.service';
import { CreateUserDto } from 'src/users/dto/create-user.dto';
import { CreateClinicDto } from 'src/clinics/dto/create-clinic.dto';
import { LocalAuthGuard } from 'src/guards/local-auth.guard';
import { CreateDoctorDto } from 'src/doctors/dto/create-doctor-dto';
import { DoctorService } from 'src/doctors/doctors.service';

@Controller()
export class AuthController {
  constructor(
    private authService: AuthService,
    private userService: UserService,
    private clinicService: ClinicService,
    private doctorService: DoctorService
  ) {}

  @Post('login')
  @UseGuards(LocalAuthGuard)
  async login(@Request() req) {
    console.log('auth');
    return this.authService.login(req.user);
  }

  @Post('register/user')
  async registerUser(@Body() createUserDto: CreateUserDto) {
    return this.userService.create(createUserDto);
  }

  @Post('register/clinic')
  async registerClinic(@Body() createClinicDto: CreateClinicDto) {
    return this.clinicService.create(createClinicDto);
  }

  @Post('register/doctor')
  async registerDoctor(@Body() createDoctorDto: CreateDoctorDto) {
    return this.doctorService.create(createDoctorDto);
  }
}
