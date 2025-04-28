import { Injectable, UnauthorizedException } from '@nestjs/common';
import { UserService } from '../users/users.service';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';
import { ClinicService } from 'src/clinics/clinic.service';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UserService,
    private readonly clinicService: ClinicService,
    private jwtService: JwtService
  ) {}

  async validateAuth(email: string, password: string): Promise<any> {
    // Сначала проверяем в пользователях
    const user = await this.usersService.findByEmail(email);
    if (user && (await bcrypt.compare(password, user.password))) {
      const { password, ...result } = user;
      return { ...result, type: 'user' };
    }

    // Если не найден пользователь, проверяем клиники
    const clinic = await this.clinicService.findByEmail(email);
    if (clinic && (await bcrypt.compare(password, clinic.password))) {
      const { password, ...result } = clinic;
      return { ...result, type: 'clinic' };
    }

    return null;
  }

  async login(user: any) {
    const payload = { email: user.email, sub: user.id };
    return {
      access_token: this.jwtService.sign(payload),
    };
  }
}
