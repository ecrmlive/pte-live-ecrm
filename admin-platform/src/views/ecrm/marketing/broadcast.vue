<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  auditPlatformBroadcastApi,
  getPlatformBroadcastApi,
  listRoomsApi,
  type Room,
} from '#/api/core/platform-broadcast';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD,
  LIST_MER_ID_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const current = ref<Room>();
const canManageBroadcast = ref(false);
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

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues, {
    statusValues: [0, 2, -1],
  });
  return base;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('房间名称'),
  LIST_MER_ID_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '审核通过', value: 2 },
        { label: '已驳回', value: -1 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<Room> = {
  columns: [
    { field: 'name', minWidth: 180, showOverflow: false, title: '房间' },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    { field: 'anchor_name', title: '主播', width: 120 },
    {
      field: 'start_time',
      minWidth: 170,
      title: '开播时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'status',
      slots: { default: 'audit_status' },
      title: '审核状态',
      width: 110,
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: 'C 端显示',
      width: 100,
    },
    {
      field: 'live_status',
      title: '直播状态',
      width: 100,
      formatter: ({ cellValue }) => liveStatus(Number(cellValue)),
    },
    platformListActionColumn({ width: 238 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listRoomsApi(buildParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'broadcast_room_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[640px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

const [RejectModal, rejectModalApi] = useVbenModal({
  onConfirm: async () => reject(),
});

async function openDetail(row: Room) {
  current.value = await getPlatformBroadcastApi(row.broadcast_room_id);
  detailDrawerApi.setState({ title: '直播间详情' }).open();
}

async function approve(row: Room) {
  try {
    await ElMessageBox.confirm(
      '确认通过该直播间审核？通过后将按商户提交的显示状态对 C 端生效。',
      '审核通过确认',
      { cancelButtonText: '取消', confirmButtonText: '确认通过', type: 'warning' },
    );
    await auditPlatformBroadcastApi(row.broadcast_room_id, {
      status: 2,
      is_show: row.is_show === 1 ? 1 : 0,
    });
    ElMessage.success('直播间已审核通过');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

function openReject(row: Room) {
  current.value = row;
  rejectForm.refusal = '';
  rejectModalApi.setState({ title: '驳回直播间' }).open();
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  rejectModalApi.lock();
  saving.value = true;
  try {
    await auditPlatformBroadcastApi(current.value.broadcast_room_id, {
      is_show: 0,
      refusal,
      status: -1,
    });
    rejectModalApi.close();
    ElMessage.success('直播间已驳回');
    gridApi.reload();
  } finally {
    saving.value = false;
    rejectModalApi.unlock();
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
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManageBroadcast.value =
    (profile.roles.includes('platform') || profile.roles.includes('operations')) &&
    permissions.includes('marketing.broadcast.audit');
});
</script>

<template>
  <Page
    auto-content-height
    description="拥有统一后台直播审核权限的账号可审核、驳回及显示/隐藏直播间；无权限账号仅查看。"
    title="直播监管"
  >
    <Grid>
      <template #audit_status="{ row }">
        <ElTag :type="auditInfo(row.status).type">
          {{ auditInfo(row.status).label }}
        </ElTag>
      </template>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <template v-if="canAudit(row)">
          <ElButton link type="success" @click="approve(row)">通过</ElButton>
          <ElButton link type="danger" @click="openReject(row)">驳回</ElButton>
        </template>
        <template v-else-if="canChangeVisibility(row)">
          <ElButton
            v-if="row.is_show !== 1"
            link
            type="success"
            @click="setVisibility(row, 1)"
          >
            显示
          </ElButton>
          <ElButton v-else link type="warning" @click="setVisibility(row, 0)">
            隐藏
          </ElButton>
        </template>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="current">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="房间名称">{{ current.name }}</ElDescriptionsItem>
          <ElDescriptionsItem label="商户">
            {{ current.mer_name || `商户 #${current.mer_id}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="主播">{{ current.anchor_name || '—' }}</ElDescriptionsItem>
          <ElDescriptionsItem label="审核状态">
            <ElTag :type="auditInfo(current.status).type">
              {{ auditInfo(current.status).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="开播时间">
            {{ formatShanghaiDateTime(current.start_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="直播状态">
            {{ liveStatus(current.live_status) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="C 端显示">
            {{ current.is_show === 1 ? '显示' : '隐藏' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="创建时间">
            {{ formatShanghaiDateTime(current.create_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="播放地址">
            {{ current.play_url || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.refusal" :span="2" label="驳回原因">
            {{ current.refusal }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.mark" :span="2" label="备注">
            {{ current.mark }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div class="mb-3 mt-6 text-base font-medium">直播挂货</div>
        <ElTable :data="current.goods || []" border empty-text="暂无挂货商品">
          <ElTableColumn label="商品 ID" min-width="100" prop="product_id" />
          <ElTableColumn label="商品名称" min-width="180" prop="store_name" show-overflow-tooltip />
          <ElTableColumn label="价格" width="110">
            <template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template>
          </ElTableColumn>
          <ElTableColumn label="状态" width="88">
            <template #default="{ row }">
              <ElTag :type="row.on_sale === 1 ? 'success' : 'info'">
                {{ row.on_sale === 1 ? '上架' : '下架' }}
              </ElTag>
            </template>
          </ElTableColumn>
        </ElTable>
      </template>
    </DetailDrawer>

    <RejectModal class="w-[480px]">
      <ElForm label-width="84px">
        <ElFormItem label="驳回原因" required>
          <ElInput
            v-model="rejectForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请填写可供商户查看的驳回原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </RejectModal>
  </Page>
</template>
