import { HttpException, HttpStatus, Injectable } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { ConfigService } from '@nestjs/config';
import { firstValueFrom } from 'rxjs';
import axios, { AxiosError } from 'axios';

@Injectable()
export class GatewayService {
  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService
  ) {}

  async forwardToAuth(
    path: string,
    body: any,
    method: string,
    query?: Record<string, any>
  ) {
    try {
      const url = new URL(
        path,
        `${this.configService.get('AUTH_SERVICE_URL')}`
      );

      if (query) {
        Object.entries(query).forEach(([key, value]) => {
          url.searchParams.append(key, String(value));
        });
      }

      const response = await axios<{
        message: string;
      }>({
        method,
        url: url.toString(),
        data: body,
        headers: {
          'Content-Type': 'application/json',
        },
        validateStatus: () => true,
      });
      console.log(response);
      if (response.status >= 400) {
        throw new HttpException(
          response.data?.message || 'Auth service error',
          response.status
        );
      }

      return response.data;
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }

      const axiosError = error as AxiosError<{ message: string }>;

      if (axiosError.response) {
        throw new HttpException(
          axiosError.response.data?.message || 'Auth service error',
          axiosError.response.status
        );
      } else if (axiosError.request) {
        throw new HttpException(
          'Auth service unavailable',
          HttpStatus.SERVICE_UNAVAILABLE
        );
      }
    }
  }

  async forwardToManagement(
    path: string,
    body: any,
    method: string,
    query?: Record<string, any>
  ) {
    try {
      const url = new URL(
        path,
        `${this.configService.get('CLINIC_MANAGEMENT_SERVICE_URL')}`
      );

      if (query) {
        Object.entries(query).forEach(([key, value]) => {
          url.searchParams.append(key, String(value));
        });
      }

      const response = await axios<{
        message: string;
      }>({
        method,
        url: url.toString(),
        data: body,
        headers: {
          'Content-Type': 'application/json',
        },
        validateStatus: () => true,
      });

      if (response.status >= 400) {
        throw new HttpException(
          response.data?.message || 'Management service error',
          response.status
        );
      }

      return response.data;
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }

      const axiosError = error as AxiosError<{ message: string }>;

      if (axiosError.response) {
        throw new HttpException(
          axiosError.response.data?.message || 'Management service error',
          axiosError.response.status
        );
      } else if (axiosError.request) {
        throw new HttpException(
          'Management service unavailable',
          HttpStatus.SERVICE_UNAVAILABLE
        );
      }
    }
  }

  async forwardToLab(path: string, body?: any) {
    const url = `${this.configService.get('LAB_SERVICE_URL')}${path}`;
    const response = this.httpService.post(url, body);
    return await firstValueFrom(response).then((res) => res.data);
  }
}
