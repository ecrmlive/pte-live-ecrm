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
  ElIcon,
  ElImage,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElRate,
  ElSwitch,
} from 'element-plus';
import { ArrowDown } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { fetchPlatformMerchants } from '#/api/core/ecrm';
import {
  forceOffPlatformPresellApi,
  getPlatformPresellTypeFilterApi,
  listPlatformPresellsApi,
  setPlatformPresellLabelsApi,
  setPlatformPresellShowApi,
  setPlatformPresellStarApi,
  type PlatformPresell,
} from '#/api/core/platform-presell';
import PresellAuditDrawer from '#/components/marketing/presell-audit-drawer.vue';
import PresellDetailDrawer from '#/components/marketing/presell-detail-drawer.vue';
import PresellEditDrawer from '#/components/marketing/presell-edit-drawer.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import ProductLabelSelectModal from '#/views/ecrm/product/components/ProductLabelSelectModal.vue';
import ProductPreviewModal from '#/views/ecrm/product/components/ProductPreviewModal.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const TYPE_TABS: Array<{ type: number; name: string }> = [
  { type: 1, name: '全款预售' },
  { type: 2, name: '定金预售' },
];

const tabType = ref(1);
const tabCounts = ref<Record<number, number>>({ 1: 0, 2: 0 });
const canManage = ref(false);
const lastFormValues = ref<Record<string, unknown>>({});
const merchantOptions = ref<{ label: string; value: number }[]>([]);
const detailDrawerRef = ref<InstanceType<typeof PresellDetailDrawer>>();
const editDrawerRef = ref<InstanceType<typeof PresellEditDrawer>>();
const auditDrawerRef = ref<InstanceType<typeof PresellAuditDrawer>>();
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
  const auditRaw = values.product_status;
  const actRaw = values.type;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    presell_type: tabType.value,
    mer_id:
      merRaw === 0 || merRaw === undefined || merRaw === null || merRaw === ''
        ? undefined
        : Number(merRaw),
    keyword: String(values.keyword ?? '').trim() || undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 ? Number(traderRaw) : undefined,
    star: starRaw === 0 || starRaw ? Number(starRaw) : undefined,
    product_status:
      auditRaw === '' || auditRaw === undefined || auditRaw === null
        ? undefined
        : Number(auditRaw),
    type:
      actRaw === '' || actRaw === undefined || actRaw === null
        ? undefined
        : Number(actRaw),
    us_status:
      usRaw === 0 || usRaw === 1 || usRaw === -1 ? Number(usRaw) : undefined,
  };
}

async function loadTabCounts(formValues?: Record<string, unknown>) {
  const values = formValues || lastFormValues.value || {};
  try {
    const data = await getPlatformPresellTypeFilterApi({
      mer_id:
        values.mer_id === 0 ||
        values.mer_id === undefined ||
        values.mer_id === null ||
        values.mer_id === ''
          ? undefined
          : Number(values.mer_id),
      keyword: String(values.keyword ?? '').trim() || undefined,
      is_trader:
        values.is_trader === 0 || values.is_trader === 1
          ? Number(values.is_trader)
          : undefined,
      star:
        values.star === 0 || values.star ? Number(values.star) : undefined,
      product_status:
        values.product_status === '' ||
        values.product_status === undefined ||
        values.product_status === null
          ? undefined
          : Number(values.product_status),
      type:
        values.type === '' ||
        values.type === undefined ||
        values.type === null
          ? undefined
          : Number(values.type),
      us_status:
        values.us_status === 0 ||
        values.us_status === 1 ||
        values.us_status === -1
          ? Number(values.us_status)
          : undefined,
    });
    const next: Record<number, number> = { 1: 0, 2: 0 };
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
        { label: '全部', value: '' },
        { label: '待审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '审核未通过', value: -1 },
        { label: '强制下架', value: -2 },
      ],
      placeholder: '全部',
    },
    defaultValue: '',
    fieldName: 'product_status',
    label: '审核状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '未开始', value: 0 },
        { label: '进行中', value: 1 },
        { label: '已结束', value: 2 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'type',
    label: '活动状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启显示', value: 1 },
        { label: '关闭显示', value: 0 },
        { label: '审核未通过/关闭', value: -1 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'us_status',
    label: '商品状态',
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
      placeholder: '请输入预售商品名称/ID',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
]);

const gridOptions: VxeGridProps<PlatformPresell> = {
  columns: [
    { field: 'product_presell_id', title: 'ID', width: 80 },
    {
      field: 'mer_name',
      minWidth: 120,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) => row.mer_name || `店铺#${row.mer_id}`,
    },
    {
      field: 'store_name',
      minWidth: 160,
      showOverflow: false,
      title: '商品名称',
    },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '商品图',
      width: 90,
    },
    {
      field: 'trader_name',
      title: '店铺类别',
      width: 90,
      formatter: ({ cellValue, row }) =>
        cellValue || (row.is_trader === 1 ? '自营' : '非自营'),
    },
    {
      field: 'price',
      title: '预售价',
      width: 100,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'presell_status_text',
      title: '预售活动状态',
      width: 110,
      formatter: ({ cellValue, row }) => {
        if (cellValue) return cellValue;
        const map: Record<number, string> = {
          0: '未开始',
          1: '进行中',
          2: '已结束',
        };
        return map[Number(row.presell_status)] || '—';
      },
    },
    {
      field: 'start_time',
      minWidth: 200,
      showOverflow: false,
      title: '预售活动日期',
      formatter: ({ row }) =>
        `${formatShanghaiDateTime(row.start_time)} 至 ${formatShanghaiDateTime(row.end_time)}`,
    },
    {
      field: 'success_num',
      title: '成功/参与人数',
      width: 120,
      formatter: ({ row }) =>
        `${row.success_num ?? 0}/${row.attend_num ?? 0}`,
    },
    {
      field: 'seles',
      title: '已售',
      width: 80,
    },
    {
      field: 'stock_count',
      title: '限量',
      width: 80,
      formatter: ({ row }) => {
        const total = Number(row.stock_count || 0);
        if (total > 0) return total;
        return Number(row.stock || 0) + Number(row.seles || 0);
      },
    },
    {
      field: 'stock',
      title: '限量剩余',
      width: 90,
    },
    {
      field: 'is_show',
      slots: { default: 'isShow' },
      title: '显示状态',
      width: 100,
    },
    {
      field: 'star',
      slots: { default: 'star' },
      title: '推荐星',
      width: 140,
    },
    platformListActionColumn({ width: 168 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const data = await listPlatformPresellsApi(
          buildParams(page, formValues),
        );
        await loadTabCounts(formValues);
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_presell_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [ForceOffModal, forceOffModalApi] = useVbenModal({
  title: '强制下架',
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确认下架',
  onConfirm: () => void submitForceOff(),
});

function setTypeTab(type: number) {
  tabType.value = type;
  gridApi.reload();
}

function openDetail(row: PlatformPresell) {
  void detailDrawerRef.value?.open(row.product_presell_id);
}

function openEdit(row: PlatformPresell) {
  void editDrawerRef.value?.open(row.product_presell_id);
}

function openPreview(row: PlatformPresell) {
  previewProductId.value = Number(row.product_id || 0);
  previewProductTitle.value = row.store_name || '';
  previewDisplayPrice.value = Number(row.price || 0);
  previewDisplayOtPrice.value =
    row.ot_price !== undefined && row.ot_price !== null
      ? Number(row.ot_price)
      : undefined;
  previewModalRef.value?.open();
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

function openLabels(row: PlatformPresell) {
  labelEditingID.value = row.product_presell_id;
  labelModalRef.value?.open({
    selectedIds: parseLabelIds(row.sys_labels),
  });
}

async function onLabelSubmit(ids: string[]) {
  if (!labelEditingID.value) return;
  try {
    await setPlatformPresellLabelsApi(
      labelEditingID.value,
      ids.map(String).join(','),
    );
    ElMessage.success('标签已更新');
    void gridApi.reload();
  } catch {
    ElMessage.error('标签更新失败');
  }
}

async function toggleShow(row: PlatformPresell) {
  const next = row.is_show === 1 ? 0 : 1;
  await setPlatformPresellShowApi(row.product_presell_id, next);
  row.is_show = next;
  ElMessage.success(next === 1 ? '已显示' : '已隐藏');
  await loadTabCounts();
}

async function changeStar(row: PlatformPresell, value: number) {
  await setPlatformPresellStarApi(row.product_presell_id, value);
  row.star = value;
}

function canAudit(row: PlatformPresell) {
  return Number(row.product_status) === 0;
}

function canForceOff(row: PlatformPresell) {
  return Number(row.product_status) === 1 && Number(row.is_show) === 1;
}

function openAudit(row: PlatformPresell) {
  void auditDrawerRef.value?.open(row.product_presell_id);
}

function onAudited() {
  void gridApi.reload();
  void loadTabCounts();
}

function openForceOff(row: PlatformPresell) {
  forceOffIds.value = [row.product_presell_id];
  forceOffReason.value = '';
  forceOffModalApi.open();
}

function onMoreCommand(command: string, row: PlatformPresell) {
  switch (command) {
    case 'edit':
      openEdit(row);
      break;
    case 'audit':
      if (canAudit(row)) openAudit(row);
      break;
    case 'labels':
      openLabels(row);
      break;
    case 'forceOff':
      if (canForceOff(row)) openForceOff(row);
      break;
    default:
      break;
  }
}

async function submitForceOff() {
  if (!forceOffReason.value.trim()) {
    ElMessage.warning('请填写下架原因');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `确认强制下架选中的 ${forceOffIds.value.length} 个预售商品？`,
      '强制下架',
      { type: 'warning', confirmButtonText: '确认下架', cancelButtonText: '取消' },
    );
  } catch {
    return;
  }
  forceOffModalApi.lock();
  try {
    await forceOffPlatformPresellApi(
      forceOffIds.value,
      forceOffReason.value.trim(),
    );
    ElMessage.success('已强制下架');
    forceOffModalApi.close();
    gridApi.reload();
  } finally {
    forceOffModalApi.unlock();
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
    codes.includes('marketing.presell.manage');
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
        <div class="presell-toolbar">
          <div class="presell-tabs" role="tablist">
            <button
              v-for="tab in TYPE_TABS"
              :key="tab.type"
              type="button"
              role="tab"
              class="presell-tabs__item"
              :aria-selected="tabType === tab.type"
              :class="{ 'is-active': tabType === tab.type }"
              @click="setTypeTab(tab.type)"
            >
              {{ tab.name }}({{ tabCounts[tab.type] || 0 }})
            </button>
          </div>
        </div>
      </template>

      <template #image="{ row }">
        <ElImage
          v-if="row.image"
          :src="resolveCosMediaUrl(row.image)"
          fit="cover"
          class="presell-thumb"
        >
          <template #error>
            <div class="presell-thumb presell-thumb--empty">无图</div>
          </template>
        </ElImage>
        <div v-else class="presell-thumb presell-thumb--empty">—</div>
      </template>

      <template #star="{ row }">
        <ElRate
          :model-value="Number(row.star || 0)"
          :max="5"
          :disabled="!canManage"
          @change="(v: number) => changeStar(row, v)"
        />
      </template>

      <template #isShow="{ row }">
        <ElSwitch
          :model-value="row.is_show === 1"
          :disabled="!canManage || Number(row.product_status) !== 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="() => toggleShow(row)"
        />
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openPreview(row)">预览</ElButton>
        <ElDropdown
          v-if="canManage"
          trigger="click"
          @command="(cmd: string) => onMoreCommand(cmd, row)"
        >
          <ElButton link type="primary">
            更多
            <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem command="edit">编辑</ElDropdownItem>
              <ElDropdownItem v-if="canAudit(row)" command="audit">
                审核
              </ElDropdownItem>
              <ElDropdownItem command="labels">编辑标签</ElDropdownItem>
              <ElDropdownItem
                v-if="canForceOff(row)"
                command="forceOff"
                divided
              >
                强制下架
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <PresellDetailDrawer ref="detailDrawerRef" />
    <PresellEditDrawer ref="editDrawerRef" @saved="onEditSaved" />
    <PresellAuditDrawer ref="auditDrawerRef" @audited="onAudited" />
    <ProductPreviewModal
      ref="previewModalRef"
      modal-title="预售预览"
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
.presell-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.presell-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 18px;
  border-bottom: 1px solid hsl(var(--border));
  padding-bottom: 8px;
}

.presell-tabs__item {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 6px 0;
  font-size: 13px;
  color: hsl(var(--muted-foreground));
  cursor: pointer;
}

.presell-tabs__item.is-active {
  color: hsl(var(--primary));
  font-weight: 600;
  box-shadow: inset 0 -2px 0 hsl(var(--primary));
}

.presell-thumb {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  overflow: hidden;
}

.presell-thumb--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: hsl(var(--muted) / 0.4);
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}
</style>
