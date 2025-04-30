import { Controller, Get } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import axios from 'axios';

@Controller('docs')
export class DocsController {
  private services: Record<string, string>;

  constructor(private configService: ConfigService) {
    this.services = {
      auth: `${this.configService.get<string>('AUTH_SERVICE_URL')}/docs-json`,
      clinic: `${this.configService.get<string>('CLINIC_MANAGEMENT_SERVICE_URL')}/docs-json`,
      lab: `${this.configService.get<string>('LAB_SERVICE_URL')}/docs-json`,
    };
  }

  @Get('aggregated')
  async getAggregatedDocs() {
    const docs = {};

    for (const [name, url] of Object.entries(this.services)) {
      try {
        const { data } = await axios.get(url);
        docs[name] = data;
      } catch (err) {
        docs[name] = { error: `Failed to fetch ${name} docs` };
      }
    }

    return docs;
  }
}
