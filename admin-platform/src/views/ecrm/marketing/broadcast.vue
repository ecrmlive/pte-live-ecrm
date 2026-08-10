<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRate,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { fetchPlatformMerchants } from '#/api/core/ecrm';
import {
  auditPlatformBroadcastApi,
  deletePlatformBroadcastApi,
  getPlatformBroadcastApi,
  listRoomsApi,
  setPlatformBroadcastRecommendApi,
  setPlatformBroadcastShowApi,
  type Room,
} from '#/api/core/platform-broadcast';
import {
  PLATFORM_LIST_GRID_CLASS,
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const current = ref<Room>();
const canManageBroadcast = ref(false);
const merchantOptions = ref<{ label: string; value: number }[]>([]);
const rejectForm = reactive({ refusal: '' });
const recommendForm = reactive({ sort: 0, star: 0 });

function auditInfo(status: number) {
  if (status === 2) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '审核未通过', type: 'danger' as const };
  return { label: '待审核', type: 'warning' as const };
}

function liveStatusText(status: number) {
  if (status === 101) return '直播中';
  if (status === 102) return '未开始';
  if (status === 103) return '已结束';
  return '—';
}

function traderText(row: Room) {
  return row.trader_name || (row.is_trader === 1 ? '自营' : '非自营');
}

function canAudit(row: Room) {
  return canManageBroadcast.value && row.status === 0;
}

function canManageRow() {
  return canManageBroadcast.value;
}

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  const merRaw = values.mer_id;
  const statusRaw = values.status_tag;
  const showRaw = values.show_type;
  const liveRaw = values.live_status;
  const starRaw = values.star;
  const traderRaw = values.is_trader;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(values.keyword ?? '').trim() || undefined,
    mer_id:
      merRaw === 0 || merRaw === undefined || merRaw === null || merRaw === ''
        ? undefined
        : Number(merRaw),
    status_tag:
      statusRaw === 0 || statusRaw === 2 || statusRaw === -1
        ? Number(statusRaw)
        : undefined,
    show_type:
      showRaw === 0 || showRaw === 1 ? Number(showRaw) : undefined,
    live_status:
      liveRaw === 101 || liveRaw === 102 || liveRaw === 103
        ? Number(liveRaw)
        : undefined,
    star:
      starRaw === 0 || starRaw === 1 || starRaw === 2 || starRaw === 3 || starRaw === 4 || starRaw === 5
        ? Number(starRaw)
        : undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 ? Number(traderRaw) : undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '全部', value: '' },
        { label: '待审核', value: 0 },
        { label: '审核通过', value: 2 },
        { label: '审核未通过', value: -1 },
      ],
      placeholder: '全部',
    },
    defaultValue: '',
    fieldName: 'status_tag',
    label: '状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'mer_id',
    label: '店铺名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '自营', value: 1 },
        { label: '非自营', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'is_trader',
    label: '店铺类别',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '显示', value: 1 },
        { label: '隐藏', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'show_type',
    label: '显示状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '直播中', value: 101 },
        { label: '未开始', value: 102 },
        { label: '已结束', value: 103 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'live_status',
    label: '直播状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '全部', value: '' },
        ...[0, 1, 2, 3, 4, 5].map((n) => ({
          label: n === 0 ? '未设置' : `${n} 星`,
          value: n,
        })),
      ],
      placeholder: '全部',
    },
    defaultValue: '',
    fieldName: 'star',
    label: '推荐级别',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入直播间名称/ID/主播昵称/微信号',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
]);

const gridOptions: VxeGridProps<Room> = {
  columns: [
    { type: 'seq', title: '序号', width: 60 },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) => row.mer_name || `店铺#${row.mer_id}`,
    },
    {
      field: 'trader_name',
      title: '店铺类别',
      width: 90,
      formatter: ({ row }) => traderText(row),
    },
    { field: 'broadcast_room_id', title: 'ID', width: 80 },
    {
      field: 'name',
      minWidth: 180,
      showOverflow: false,
      title: '直播间名称',
    },
    { field: 'anchor_name', title: '主播昵称', width: 110 },
    {
      field: 'anchor_wechat',
      title: '主播微信号',
      width: 130,
      formatter: ({ cellValue }) => cellValue || '—',
    },
    {
      field: 'start_time',
      minWidth: 170,
      title: '直播开始时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'end_time',
      minWidth: 170,
      title: '直播计划结束时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'live_status',
      title: '直播状态',
      width: 100,
      formatter: ({ cellValue }) => liveStatusText(Number(cellValue)),
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'star',
      slots: { default: 'star' },
      title: '推荐级别',
      width: 130,
    },
    { field: 'sort', title: '排序', width: 70 },
    platformListActionColumn({ width: 260 }),
  ],
  pagerConfig: platformListPagerConfig(),
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

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridClass: PLATFORM_LIST_GRID_CLASS,
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

const [RejectDrawer, rejectDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => reject(),
});

const [RecommendDrawer, recommendDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => saveRecommend(),
});

async function openDetail(row: Room) {
  current.value = await getPlatformBroadcastApi(row.broadcast_room_id);
  detailDrawerApi.setState({ title: '直播间详情' }).open();
}

async function approve(row: Room) {
  try {
    await confirm({
      title: '提示',
      content: '确认通过该直播间审核？通过后可按显示开关对 C 端生效。',
      icon: 'warning',
    });
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
  rejectDrawerApi.setState({ title: '驳回直播间' }).open();
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  rejectDrawerApi.lock();
  saving.value = true;
  try {
    await auditPlatformBroadcastApi(current.value.broadcast_room_id, {
      is_show: 0,
      refusal,
      status: -1,
    });
    rejectDrawerApi.close();
    ElMessage.success('直播间已驳回');
    gridApi.reload();
  } finally {
    saving.value = false;
    rejectDrawerApi.unlock();
  }
}

async function toggleShow(row: Room, value: boolean | string | number) {
  if (!canManageRow()) return;
  const next = value ? 1 : 0;
  const prev = row.is_show;
  row.is_show = next;
  try {
    await setPlatformBroadcastShowApi(row.broadcast_room_id, next as 0 | 1);
    ElMessage.success(next === 1 ? '已显示' : '已隐藏');
  } catch {
    row.is_show = prev;
  }
}

function openRecommend(row: Room) {
  current.value = row;
  recommendForm.sort = Number(row.sort || 0);
  recommendForm.star = Number(row.star || 0);
  recommendDrawerApi.setState({ title: '推荐设置' }).open();
}

async function saveRecommend() {
  if (!current.value) return;
  if (recommendForm.star < 0 || recommendForm.star > 5) {
    ElMessage.warning('请选择 0-5 星推荐级别');
    return;
  }
  recommendDrawerApi.lock();
  saving.value = true;
  try {
    await setPlatformBroadcastRecommendApi(current.value.broadcast_room_id, {
      sort: Number(recommendForm.sort || 0),
      star: Number(recommendForm.star || 0),
    });
    recommendDrawerApi.close();
    ElMessage.success('推荐设置已保存');
    gridApi.reload();
  } finally {
    saving.value = false;
    recommendDrawerApi.unlock();
  }
}

async function removeRoom(row: Room) {
  try {
    await confirm({
      title: '提示',
      content: `确认删除直播间「${row.name}」？`,
      icon: 'warning',
    });
    await deletePlatformBroadcastApi(row.broadcast_room_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

async function loadMerchants() {
  const data = await fetchPlatformMerchants({ page: 1, limit: 200 });
  merchantOptions.value = (data.list || []).map((row) => ({
    label: row.mer_name || `店铺#${row.mer_id}`,
    value: row.mer_id,
  }));
  await gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'mer_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantOptions.value,
        placeholder: '请选择',
      },
    },
  ]);
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManageBroadcast.value =
    (profile.roles.includes('platform') ||
      profile.roles.includes('operations')) &&
    permissions.includes('marketing.broadcast.audit');
  await loadMerchants().catch(() => undefined);
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #star="{ row }">
        <ElRate :model-value="Number(row.star || 0)" disabled />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">查看</ElButton>
        <template v-if="canAudit(row)">
          <ElButton link type="success" @click="approve(row)">通过</ElButton>
          <ElButton link type="danger" @click="openReject(row)">驳回</ElButton>
        </template>
        <template v-if="canManageRow()">
          <ElButton link type="primary" @click="openRecommend(row)">
            推荐
          </ElButton>
          <ElSwitch
            :model-value="row.is_show === 1"
            class="ml-1"
            inline-prompt
            active-text="显"
            inactive-text="隐"
            @change="(v) => toggleShow(row, v)"
          />
          <ElButton link type="danger" @click="removeRoom(row)">删除</ElButton>
        </template>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="current">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="直播间名称">
            {{ current.name }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="ID">
            {{ current.broadcast_room_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺名称">
            {{ current.mer_name || `店铺#${current.mer_id}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺类别">
            {{ traderText(current) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="主播昵称">
            {{ current.anchor_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="主播微信号">
            {{ current.anchor_wechat || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="审核状态">
            <ElTag :type="auditInfo(current.status).type">
              {{ auditInfo(current.status).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="直播状态">
            {{ liveStatusText(current.live_status) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="直播开始时间">
            {{ formatShanghaiDateTime(current.start_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="直播计划结束时间">
            {{ formatShanghaiDateTime(current.end_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="显示状态">
            {{ current.is_show === 1 ? '显示' : '隐藏' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="推荐级别">
            {{ current.star || 0 }} 星
          </ElDescriptionsItem>
          <ElDescriptionsItem label="排序">
            {{ current.sort || 0 }}
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
          <ElTableColumn
            label="商品名称"
            min-width="180"
            prop="store_name"
            show-overflow-tooltip
          />
          <ElTableColumn label="价格" width="110">
            <template #default="{ row }">
              ¥{{ Number(row.price || 0).toFixed(2) }}
            </template>
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

    <RejectDrawer>
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
    </RejectDrawer>

    <RecommendDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="推荐级别" required>
          <ElRate v-model="recommendForm.star" :max="5" />
        </ElFormItem>
        <ElFormItem label="排序" required>
          <ElInputNumber v-model="recommendForm.sort" :min="0" :max="9999" />
        </ElFormItem>
      </ElForm>
    </RecommendDrawer>
  </Page>
</template>
