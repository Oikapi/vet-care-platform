import { PassportStrategy } from '@nestjs/passport';
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { AuthService } from './auth.service';
import { Strategy } from 'passport-local';

@Injectable()
export class LocalStrategy extends PassportStrategy(Strategy) {
  constructor(private authService: AuthService) {
    super({ usernameField: 'email' }); // Меняем стандартное поле 'username' на 'email'
  }

  async validate(email: string, password: string) {
    const auth = await this.authService.validateAuth(email, password);
    if (!auth) {
      throw new UnauthorizedException('Неверные учетные данные');
    }
    return auth;
  }
}
