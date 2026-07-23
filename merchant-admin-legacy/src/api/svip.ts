import { http } from '@/api/http';

export interface SvipConfig {
  mer_id: number;
  svip_coupon_merge: number;
}

export function fetchSvipConfig() {
  return http.get<SvipConfig>('/setting/svip');
}

export function updateSvipConfig(svip_coupon_merge: number) {
  return http.put<SvipConfig>('/setting/svip', { svip_coupon_merge });
}
