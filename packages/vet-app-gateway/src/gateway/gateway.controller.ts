import {
  Controller,
  Get,
  Post,
  Req,
  Res,
  Body,
  Param,
  All,
} from '@nestjs/common';
import { Request, Response } from 'express';
import { GatewayService } from './gateway.service';

@Controller()
export class GatewayController {
  constructor(private readonly gatewayService: GatewayService) {}

  @All('auth/*path')
  handleAuthRequests(@Req() req: Request, @Param('path') path: string) {
    const fullPath = path ? `/${path}` : '/';
    return this.gatewayService.forwardToAuth(
      fullPath,
      req.body,
      req.method,
      req.query
    );
  }

  @All('management/*path')
  handleManagementRequests(@Req() req: Request, @Param('path') path: string) {
    const fullPath = path ? `/${path}` : '/';
    return this.gatewayService.forwardToManagement(
      fullPath,
      req.body,
      req.method,
      req.query
    );
  }

  @Post('lab/analyze')
  async sendLabData(@Body() body: any) {
    return this.gatewayService.forwardToLab('/analyze', body);
  }
}
