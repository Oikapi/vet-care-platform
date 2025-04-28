import { IsEmail, IsString, MinLength } from 'class-validator';

export class CreateClinicDto {
  @IsString()
  name: string;

  @IsString()
  phone: string;

  @IsEmail()
  email: string;

  @IsString()
  @MinLength(6)
  password: string;

  @IsString()
  photo?: string;
}
