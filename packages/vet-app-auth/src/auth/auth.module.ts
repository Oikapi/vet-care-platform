import { ConfigurableModuleBuilder, Module } from '@nestjs/common';
import { AuthService } from './auth.service';
import { AuthController } from './auth.controller';
import { JwtModule } from '@nestjs/jwt';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from 'src/users/user.entity';
import { Clinic } from 'src/clinics/clinic.entity';
import { ClinicService } from 'src/clinics/clinic.service';
import { LocalStrategy } from './local.strategy';
import { UserService } from 'src/users/users.service';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true, // чтобы доступ был везде через process.env
    }),
    TypeOrmModule.forFeature([User, Clinic]),
    JwtModule.register({
      secret: process.env.JWT_SECRET,
      signOptions: { expiresIn: '24h' },
    }),
  ],
  providers: [AuthService, UserService, ClinicService, LocalStrategy],
  controllers: [AuthController],
  exports: [AuthService],
})
export class AuthModule {}
