import { liveAdminPost } from '#/api/live-request';

const LiveTrafficApi = {
  account(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/account', data);
  },
  recharge(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/recharge', data);
  },
  rechargeList(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/recharge/list', data);
  },
  sessionList(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/session/list', data);
  },
  liveList(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/live/list', data);
  },
  usageRank(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/usage/rank', data);
  },
  exportRecharge(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/export/recharge', data);
  },
  exportSettlement(data: Record<string, unknown>) {
    return liveAdminPost('/api/v1/admin/traffic/export/settlement', data);
  },
};

export default LiveTrafficApi;
