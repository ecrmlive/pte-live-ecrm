<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElRate,
  ElSwitch,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { fetchPlatformMerchants } from '#/api/core/ecrm';
import {
  deletePlatformSeckillApi,
  forceOffPlatformSeckillApi,
  getPlatformSeckillStatusFilterApi,
  listPlatformSeckillApi,
  setPlatformSeckillLabelsApi,
  setPlatformSeckillShowApi,
  setPlatformSeckillStarApi,
  type PlatformSeckillActive,
} from '#/api/core/platform-seckill';
import SeckillDetailDrawer from '#/components/marketing/seckill-detail-drawer.vue';
import SeckillEditDrawer from '#/components/marketing/seckill-edit-drawer.vue';
import ProductLabelSelectModal from '#/views/ecrm/product/components/ProductLabelSelectModal.vue';
import ProductPreviewModal from '#/views/ecrm/product/components/ProductPreviewModal.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const STATUS_TABS: Array<{ type: number; name: string }> = [
  { type: 1, name: '出售中秒杀商品' },
  { type: 2, name: '仓库中秒杀商品' },
  { type: 6, name: '待审核秒杀商品' },
  { type: 7, name: '审核未通过秒杀商品' },
  { type: 5, name: '回收站秒杀商品' },
];

const tabType = ref(1);
const tabCounts = ref<Record<number, number>>({
  1: 0,
  2: 0,
  5: 0,
  6: 0,
  7: 0,
});
const canManage = ref(false);
const selectedIds = ref<number[]>([]);
const lastFormValues = ref<Record<string, unknown>>({});
const merchantOptions = ref<{ label: string; value: number }[]>([]);
const detailDrawerRef = ref<InstanceType<typeof SeckillDetailDrawer>>();
const editDrawerRef = ref<InstanceType<typeof SeckillEditDrawer>>();
const previewModalRef = ref<InstanceType<typeof ProductPreviewModal>>();
const labelModalRef = ref<InstanceType<typeof ProductLabelSelectModal>>();
const forceOffIds = ref<number[]>([]);
const forceOffReason = ref('');
const labelEditingID = ref(0);
const previewProductId = ref(0);
const previewProductTitle = ref('');
const previewDisplayPrice = ref<number>();
const previewDisplayOtPrice = ref<number>();

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  lastFormValues.value = values;
  const merRaw = values.mer_id;
  const usRaw = values.us_status;
  const starRaw = values.star;
  const traderRaw = values.is_trader;
  const goodsStatus = Number(values.goods_status);
  // 商品状态下拉可快速切到对应 Tab（与状态条一致）
  const type = [1, 2, 6, 7].includes(goodsStatus)
    ? goodsStatus
    : tabType.value;
  if (type !== tabType.value) {
    tabType.value = type;
  }
  return {
    page: page.currentPage,
    limit: page.pageSize,
    type,
    mer_id:
      merRaw === 0 || merRaw === undefined || merRaw === null || merRaw === ''
        ? undefined
        : Number(merRaw),
    active_name: String(values.active_name ?? '').trim() || undefined,
    keyword: String(values.keyword ?? '').trim() || undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 ? Number(traderRaw) : undefined,
    star: starRaw === 0 || starRaw ? Number(starRaw) : undefined,
    us_status:
      usRaw === 0 || usRaw === 1 || usRaw === -1 || usRaw === -2
        ? Number(usRaw)
        : undefined,
    sys_labels: String(values.sys_labels ?? '').trim() || undefined,
  };
}

async function loadTabCounts(formValues?: Record<string, unknown>) {
  const values = formValues || lastFormValues.value || {};
  try {
    const data = await getPlatformSeckillStatusFilterApi({
      mer_id:
        values.mer_id === 0 ||
        values.mer_id === undefined ||
        values.mer_id === null ||
        values.mer_id === ''
          ? undefined
          : Number(values.mer_id),
      active_name: String(values.active_name ?? '').trim() || undefined,
      keyword: String(values.keyword ?? '').trim() || undefined,
      is_trader:
        values.is_trader === 0 || values.is_trader === 1
          ? Number(values.is_trader)
          : undefined,
      star:
        values.star === 0 || values.star ? Number(values.star) : undefined,
      us_status:
        values.us_status === 0 ||
        values.us_status === 1 ||
        values.us_status === -1 ||
        values.us_status === -2
          ? Number(values.us_status)
          : undefined,
    });
    const next: Record<number, number> = { 1: 0, 2: 0, 5: 0, 6: 0, 7: 0 };
    for (const item of data.list || []) {
      next[item.type] = item.count;
    }
    tabCounts.value = next;
  } catch {
    /* ignore */
  }
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启', value: 1 },
        { label: '关闭显示', value: 0 },
        { label: '审核未通过', value: -1 },
        { label: '强制下架', value: -2 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'us_status',
    label: '活动状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '出售中', value: 1 },
        { label: '仓库中', value: 2 },
        { label: '待审核', value: 6 },
        { label: '审核未通过', value: 7 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'goods_status',
    label: '商品状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '标签 ID，多个逗号分隔',
    },
    fieldName: 'sys_labels',
    label: '标签',
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
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入活动名称',
    },
    fieldName: 'active_name',
    label: '活动名称',
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
      options: [0, 1, 2, 3, 4, 5].map((n) => ({
        label: n === 0 ? '未设置' : `${n} 星`,
        value: n,
      })),
      placeholder: '全部',
    },
    fieldName: 'star',
    label: '推荐级别',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入商品名称、关键字、编号',
    },
    fieldName: 'keyword',
    label: '商品搜索',
  },
]);

const gridOptions: VxeGridProps<PlatformSeckillActive> = {
  checkboxConfig: { highlight: true, reserve: true },
  columns: [
    { type: 'checkbox', width: 48 },
    { field: 'seckill_active_id', title: 'ID', width: 80 },
    {
      field: 'mer_name',
      minWidth: 120,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) => row.mer_name || `店铺#${row.mer_id}`,
    },
    {
      field: 'trader_name',
      title: '店铺类别',
      width: 90,
      formatter: ({ cellValue, row }) =>
        cellValue || (row.is_trader === 1 ? '自营' : '非自营'),
    },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '商品图片',
      width: 90,
    },
    {
      field: 'store_name',
      minWidth: 160,
      showOverflow: false,
      title: '商品名称',
      formatter: ({ row }) => row.store_name || `商品#${row.product_id}`,
    },
    {
      field: 'time_titles',
      minWidth: 100,
      showOverflow: false,
      title: '活动名称',
      formatter: ({ row }) => row.time_titles || row.name || '—',
    },
    {
      field: 'sales',
      title: '已售数量',
      width: 90,
    },
    {
      field: 'stock',
      title: '限量剩余',
      width: 90,
    },
    {
      field: 'seckill_price',
      title: '秒杀价',
      width: 100,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'star',
      slots: { default: 'star' },
      title: '推荐星级',
      width: 130,
    },
    { field: 'sort', title: '排序', width: 70 },
    {
      field: 'product_status_name',
      title: '商品状态',
      width: 110,
      formatter: ({ cellValue, row }) =>
        cellValue ||
        ({
          1: row.is_show === 1 ? '出售中' : '仓库中',
          0: '待审核',
          [-1]: '审核未通过',
          [-2]: '平台关闭',
        }[row.product_status as 1 | 0 | -1 | -2] || '—'),
    },
    {
      field: 'activity_status_text',
      title: '活动状态',
      width: 100,
      formatter: ({ cellValue, row }) => {
        if (cellValue) return cellValue;
        const map: Record<number, string> = {
          0: '未开始',
          1: '进行中',
          [-1]: '已结束',
        };
        return map[Number(row.activity_status)] || '—';
      },
    },
    {
      field: 'is_show',
      slots: { default: 'isShow' },
      title: '是否显示',
      width: 100,
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const data = await listPlatformSeckillApi(buildParams(page, formValues));
        await loadTabCounts(formValues);
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'seckill_active_id' },
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
  gridEvents: {
    checkboxAll({ records }: { records: PlatformSeckillActive[] }) {
      selectedIds.value = records.map((r) => r.seckill_active_id);
    },
    checkboxChange({ records }: { records: PlatformSeckillActive[] }) {
      selectedIds.value = records.map((r) => r.seckill_active_id);
    },
  },
});

const [ForceOffModal, forceOffModalApi] = useVbenModal({
  title: '强制下架',
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确认下架',
  onConfirm: () => void submitForceOff(),
});

function setStatusTab(type: number) {
  tabType.value = type;
  selectedIds.value = [];
  gridApi.reload();
}

function openDetail(row: PlatformSeckillActive) {
  void detailDrawerRef.value?.open(row.seckill_active_id);
}

function openPreview(row: PlatformSeckillActive) {
  previewProductId.value = Number(row.product_id || 0);
  previewProductTitle.value = row.store_name || row.name || '';
  previewDisplayPrice.value = Number(row.seckill_price || 0);
  previewDisplayOtPrice.value =
    row.price !== undefined && row.price !== null
      ? Number(row.price)
      : undefined;
  previewModalRef.value?.open();
}

function openEdit(row: PlatformSeckillActive) {
  void editDrawerRef.value?.open(row.seckill_active_id);
}

function onEditSaved() {
  void gridApi.reload();
  void loadTabCounts();
}

function parseLabelIds(raw?: string) {
  if (!raw?.trim()) return [] as string[];
  return raw
    .split(/[,，\s]+/)
    .map((s) => s.trim())
    .filter(Boolean);
}

function openLabels(row: PlatformSeckillActive) {
  labelEditingID.value = row.seckill_active_id;
  labelModalRef.value?.open({
    selectedIds: parseLabelIds(row.sys_labels),
  });
}

async function onLabelSubmit(ids: string[]) {
  if (!labelEditingID.value) return;
  try {
    await setPlatformSeckillLabelsApi(
      labelEditingID.value,
      ids.map(String).join(','),
    );
    ElMessage.success('标签已更新');
    void gridApi.reload();
  } catch {
    ElMessage.error('标签更新失败');
  }
}

async function toggleShow(row: PlatformSeckillActive) {
  const next = row.is_show === 1 ? 0 : 1;
  await setPlatformSeckillShowApi(row.seckill_active_id, next);
  row.is_show = next;
  ElMessage.success(next === 1 ? '已显示' : '已隐藏');
  await loadTabCounts();
}

async function changeStar(row: PlatformSeckillActive, value: number) {
  await setPlatformSeckillStarApi(row.seckill_active_id, value);
  row.star = value;
}

function openForceOff(ids: number[]) {
  if (!ids.length) {
    ElMessage.warning('请先选择秒杀商品');
    return;
  }
  forceOffIds.value = ids;
  forceOffReason.value = '';
  forceOffModalApi.open();
}

async function submitForceOff() {
  if (!forceOffReason.value.trim()) {
    ElMessage.warning('请填写下架原因');
    return;
  }
  forceOffModalApi.lock();
  try {
    await forceOffPlatformSeckillApi(
      forceOffIds.value,
      forceOffReason.value.trim(),
    );
    ElMessage.success('已强制下架');
    forceOffModalApi.close();
    selectedIds.value = [];
    gridApi.reload();
  } finally {
    forceOffModalApi.unlock();
  }
}

async function remove(row: PlatformSeckillActive) {
  try {
    await ElMessageBox.confirm(
      `将「${row.store_name || row.name}」移入回收站？`,
      '加入回收站',
      { type: 'warning' },
    );
    await deletePlatformSeckillApi(row.seckill_active_id);
    ElMessage.success('已移入回收站');
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
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManage.value =
    profile.roles.some((r) => r === 'platform' || r === 'operations') &&
    (codes.includes('marketing.seckill.manage') ||
      codes.includes('marketing.seckill.manage.page'));
  await loadMerchants().catch(() => undefined);
  if (canManage.value) {
    await loadTabCounts();
    gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="seckill-toolbar">
          <div class="seckill-tabs" role="tablist">
            <button
              v-for="tab in STATUS_TABS"
              :key="tab.type"
              type="button"
              role="tab"
              class="seckill-tabs__item"
              :aria-selected="tabType === tab.type"
              :class="{ 'is-active': tabType === tab.type }"
              @click="setStatusTab(tab.type)"
            >
              {{ tab.name }}({{ tabCounts[tab.type] || 0 }})
            </button>
          </div>
          <div class="seckill-toolbar__actions">
            <ElButton
              v-if="canManage && tabType === 1"
              :disabled="!selectedIds.length"
              @click="openForceOff(selectedIds)"
            >
              批量强制下架
            </ElButton>
          </div>
        </div>
      </template>

      <template #image="{ row }">
        <ElImage
          v-if="row.image"
          :src="resolveCosMediaUrl(row.image)"
          fit="cover"
          class="seckill-thumb"
        >
          <template #error>
            <div class="seckill-thumb seckill-thumb--empty">无图</div>
          </template>
        </ElImage>
        <div v-else class="seckill-thumb seckill-thumb--empty">—</div>
      </template>

      <template #star="{ row }">
        <ElRate
          :model-value="Number(row.star || 0)"
          :max="5"
          @change="(v: number) => changeStar(row, v)"
        />
      </template>

      <template #isShow="{ row }">
        <ElSwitch
          :model-value="row.is_show === 1"
          :disabled="!canManage || tabType === 5"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="() => toggleShow(row)"
        />
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openPreview(row)">预览</ElButton>
        <ElDropdown v-if="canManage && tabType !== 5" trigger="click">
          <ElButton link type="primary">更多</ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem @click="openEdit(row)">编辑</ElDropdownItem>
              <ElDropdownItem @click="openLabels(row)">选择标签</ElDropdownItem>
              <ElDropdownItem
                v-if="tabType === 1"
                @click="openForceOff([row.seckill_active_id])"
              >
                强制下架
              </ElDropdownItem>
              <ElDropdownItem @click="remove(row)">加入回收站</ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <SeckillDetailDrawer ref="detailDrawerRef" />
    <SeckillEditDrawer ref="editDrawerRef" @saved="onEditSaved" />
    <ProductPreviewModal
      ref="previewModalRef"
      modal-title="秒杀预览"
      :product-id="previewProductId"
      :product-title="previewProductTitle"
      :display-price="previewDisplayPrice"
      :display-ot-price="previewDisplayOtPrice"
    />
    <ProductLabelSelectModal ref="labelModalRef" @submit="onLabelSubmit" />

    <ForceOffModal>
      <ElForm label-width="90px">
        <ElFormItem label="下架原因" required>
          <ElInput
            v-model="forceOffReason"
            type="textarea"
            :rows="3"
            maxlength="200"
            show-word-limit
            placeholder="请填写强制下架原因"
          />
        </ElFormItem>
      </ElForm>
    </ForceOffModal>
  </Page>
</template>

<style scoped>
.seckill-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.seckill-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 18px;
  border-bottom: 1px solid hsl(var(--border));
  padding-bottom: 8px;
}

.seckill-tabs__item {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 6px 0;
  font-size: 13px;
  color: hsl(var(--muted-foreground));
  cursor: pointer;
}

.seckill-tabs__item.is-active {
  color: hsl(var(--primary));
  font-weight: 600;
  box-shadow: inset 0 -2px 0 hsl(var(--primary));
}

.seckill-toolbar__actions {
  display: flex;
  gap: 8px;
}

.seckill-thumb {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  overflow: hidden;
}

.seckill-thumb--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: hsl(var(--muted) / 0.4);
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

</style>
