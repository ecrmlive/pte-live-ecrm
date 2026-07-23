import { applyAdminAuthorization } from '#/utils/live-api-auth';
import { liveAdminClient } from '#/api/live-request';

const LiveConfigApi = {
  reloadTencent() {
    return liveAdminClient.post('/api/v1/admin/config/reload-tencent');
  },
};

export default LiveConfigApi;
