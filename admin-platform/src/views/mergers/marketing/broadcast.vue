<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  auditPlatformBroadcastApi,
  getPlatformBroadcastApi,
  listRoomsApi,
  type Room,
} from '#/api/core/platform-broadcast';

const loading = ref(false);
const saving = ref(false);
const rows = ref<Room[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const current = ref<Room>();
const isPlatformOperator = ref(false);
const canManageBroadcast = ref(false);
const query = reactive({ limit: 20, page: 1 });
const rejectForm = reactive({ refusal: '' });

function auditInfo(status: number) {
  if (status === 2) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  return { label: '待审核', type: 'warning' as const };
}

function liveStatus(status: number) {
  return status === 101 ? '直播中' : status === 102 ? '未开始' : '已结束';
}

function canAudit(row: Room) {
  return canManageBroadcast.value && row.status === 0;
}

function canChangeVisibility(row: Room) {
  return canManageBroadcast.value && row.status === 2;
}

async function load() {
  loading.value = true;
  try {
    const result = await listRoomsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

async function openDetail(row: Room) {
  current.value = await getPlatformBroadcastApi(row.broadcast_room_id);
  detailOpen.value = true;
}

async function approve(row: Room) {
  try {
    await ElMessageBox.confirm(
      '确认通过该直播间审核？通过后将按商户提交的显示状态对 C 端生效。',
      '审核通过确认',
      { cancelButtonText: '取消', confirmButtonText: '确认通过', type: 'warning' },
    );
    await auditPlatformBroadcastApi(row.broadcast_room_id, { status: 2, is_show: row.is_show === 1 ? 1 : 0 });
    ElMessage.success('直播间已审核通过');
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

function openReject(row: Room) {
  current.value = row;
  rejectForm.refusal = '';
  rejectOpen.value = true;
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  saving.value = true;
  try {
    await auditPlatformBroadcastApi(current.value.broadcast_room_id, { is_show: 0, refusal, status: -1 });
    rejectOpen.value = false;
    ElMessage.success('直播间已驳回');
    await load();
  } finally {
    saving.value = false;
  }
}

async function setVisibility(row: Room, isShow: 0 | 1) {
  const action = isShow === 1 ? '显示' : '隐藏';
  try {
    await ElMessageBox.confirm(`确认${action}该已审核通过的直播间？`, `${action}直播间`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await auditPlatformBroadcastApi(row.broadcast_room_id, { is_show: isShow, status: 0 });
    ElMessage.success(`直播间已${action}`);
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  isPlatformOperator.value = profile.is_agent !== 1;
  canManageBroadcast.value = isPlatformOperator.value && permissions.includes('broadcast/audit');
  await load();
});
</script>

<template>
  <Page title="直播监管" description="有“直播审房”按钮权限的平台账号可审核、驳回及显示/隐藏直播间；无权限账号仅查看。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="broadcast_room_id">
        <el-table-column label="房间" min-width="180" prop="name" show-overflow-tooltip />
        <el-table-column label="商户" min-width="130"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column>
        <el-table-column label="主播" prop="anchor_name" width="120" />
        <el-table-column label="开播时间" min-width="170" prop="start_time" />
        <el-table-column label="审核状态" width="110"><template #default="{ row }"><el-tag :type="auditInfo(row.status).type">{{ auditInfo(row.status).label }}</el-tag></template></el-table-column>
        <el-table-column label="C 端显示" width="100"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag></template></el-table-column>
        <el-table-column label="直播状态" width="100"><template #default="{ row }">{{ liveStatus(row.live_status) }}</template></el-table-column>
        <el-table-column fixed="right" label="操作" width="238">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <template v-if="canAudit(row)">
              <el-button link type="success" @click="approve(row)">通过</el-button>
              <el-button link type="danger" @click="openReject(row)">驳回</el-button>
            </template>
            <template v-else-if="canChangeVisibility(row)">
              <el-button v-if="row.is_show !== 1" link type="success" @click="setVisibility(row, 1)">显示</el-button>
              <el-button v-else link type="warning" @click="setVisibility(row, 0)">隐藏</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination :current-page="query.page" :page-size="query.limit" :total="total" background layout="total, prev, pager, next" @current-change="(page: number) => { query.page = page; load(); }" />
      </div>
    </el-card>

    <el-drawer v-model="detailOpen" :with-header="false" size="640px">
      <template v-if="current">
        <div class="mb-5 text-lg font-medium">直播间详情</div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="房间名称">{{ current.name }}</el-descriptions-item>
          <el-descriptions-item label="商户">{{ current.mer_name || `商户 #${current.mer_id}` }}</el-descriptions-item>
          <el-descriptions-item label="主播">{{ current.anchor_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="审核状态"><el-tag :type="auditInfo(current.status).type">{{ auditInfo(current.status).label }}</el-tag></el-descriptions-item>
          <el-descriptions-item label="开播时间">{{ current.start_time || '—' }}</el-descriptions-item>
          <el-descriptions-item label="直播状态">{{ liveStatus(current.live_status) }}</el-descriptions-item>
          <el-descriptions-item label="C 端显示">{{ current.is_show === 1 ? '显示' : '隐藏' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ current.create_time }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="播放地址">{{ current.play_url || '—' }}</el-descriptions-item>
          <el-descriptions-item v-if="current.refusal" :span="2" label="驳回原因">{{ current.refusal }}</el-descriptions-item>
          <el-descriptions-item v-if="current.mark" :span="2" label="备注">{{ current.mark }}</el-descriptions-item>
        </el-descriptions>
        <div class="mb-3 mt-6 text-base font-medium">直播挂货</div>
        <el-table :data="current.goods || []" border empty-text="暂无挂货商品">
          <el-table-column label="商品 ID" min-width="100" prop="product_id" />
          <el-table-column label="商品名称" min-width="180" prop="store_name" show-overflow-tooltip />
          <el-table-column label="价格" width="110"><template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template></el-table-column>
          <el-table-column label="状态" width="88"><template #default="{ row }"><el-tag :type="row.on_sale === 1 ? 'success' : 'info'">{{ row.on_sale === 1 ? '上架' : '下架' }}</el-tag></template></el-table-column>
        </el-table>
      </template>
    </el-drawer>

    <el-dialog v-model="rejectOpen" destroy-on-close title="驳回直播间" width="480px">
      <el-form label-width="84px"><el-form-item label="驳回原因" required><el-input v-model="rejectForm.refusal" :rows="4" maxlength="200" placeholder="请填写可供商户查看的驳回原因" show-word-limit type="textarea" /></el-form-item></el-form>
      <template #footer><el-button @click="rejectOpen = false">取消</el-button><el-button :loading="saving" type="danger" @click="reject">确认驳回</el-button></template>
    </el-dialog>
  </Page>
</template>
