export class CreateDoctorDto {
  readonly firstName: string;
  readonly lastName: string;
  readonly specialization: string;
  readonly clinicId: number;
  readonly email: string;
  readonly password: string;
}

export class UpdateDoctorDto {
  readonly firstName?: string;
  readonly lastName?: string;
  readonly specialization?: string;
  readonly clinicId?: number;
  readonly email?: string;
}
