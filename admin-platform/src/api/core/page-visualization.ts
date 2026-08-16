import { requestClient } from '#/api/request';

export interface PageVisualizationImage {
  name: string;
  url: string;
  link: string;
}

export interface PageVisualizationReward {
  day: number;
  points: number;
  sort: number;
  enabled: boolean;
}

export interface PageVisualizationRange {
  name: string;
  min: number;
  max: number;
  sort: number;
  enabled: boolean;
}

export interface PageVisualizationConfig {
  carousels: Record<string, PageVisualizationImage[]>;
  sign_rewards: PageVisualizationReward[];
  point_ranges: PageVisualizationRange[];
  splash: {
    enabled: boolean;
    display_seconds: number;
    interval_hours: number;
    images: PageVisualizationImage[];
  };
}

function parseConfig(raw: string) {
  return JSON.parse(raw) as PageVisualizationConfig;
}

export function getPageVisualizationConfigApi() {
  return requestClient.get<{ config: string; note: string }>('/setting/page-visualization').then((data) => ({
    config: parseConfig(data.config),
    note: data.note,
  }));
}

export function savePageVisualizationConfigApi(config: PageVisualizationConfig) {
  return requestClient
    .put<{ config: string }>('/setting/page-visualization', { config: JSON.stringify(config) })
    .then((data) => parseConfig(data.config));
}
