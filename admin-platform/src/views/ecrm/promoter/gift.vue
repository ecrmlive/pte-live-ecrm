<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  auditPlatformProductApi,
  batchForceOffPlatformProductsApi,
  forceOffPlatformProductApi,
  getPlatformProductStatusFilterApi,
  listPlatformProductsApi,
  setPlatformProductShowApi,
  type PlatformProduct,
  type PlatformProductStatusFilter,
} from '#/api/core/platform-catalog';
import {
  fetchMerchantCategories,
  type MerchantCategoryRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import ProductDetailDrawer from '../product/components/ProductDetailDrawer.vue';

const STATUS_TABS: Array<{ type: number; name: string }> = [
  { type: 1, name: '出售中礼包' },
  { type: 2, name: '仓库中礼包' },
  { type: 6, name: '待审核礼包' },
  { type: 7, name: '审核未通过礼包' },
  { type: 5, name: '回收站礼包' },
];

const tabType = ref(6);
const tabCounts = ref<Record<number, number>>({
  1: 0,
  2: 0,
  5: 0,
  6: 0,
  7: 0,
});
const selectedIds = ref<number[]>([]);
const canAudit = ref(false);
const lastFormValues = ref<Record<string, unknown>>({});
const merchantCategoryOptions = ref<{ label: string; value: number }[]>([]);
const detailDrawerRef = ref<InstanceType<typeof ProductDetailDrawer>>();
const forceOffReason = ref('');
const forceOffIds = ref<number[]>([]);
const forceOffSubmitting = ref(false);
const rejecting = ref(false);
const rejectForm = ref({ refusal: '' });
const current = ref<PlatformProduct>();

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  lastFormValues.value = values;
  const merCategoryId = Number(values.mer_category_id || 0);
  return {
    page: page.currentPage,
    limit: page.pageSize,
    type: tabType.value,
    is_gift_bag: 1 as const,
    keyword: String(values.keyword ?? '').trim() || undefined,
    store_name: String(values.store_name ?? '').trim() || undefined,
    mer_category_id: merCategoryId > 0 ? merCategoryId : undefined,
  };
}

async function loadTabCounts(formValues?: Record<string, unknown>) {
  const values = formValues || lastFormValues.value || {};
  try {
    const data = await getPlatformProductStatusFilterApi({
      is_gift_bag: 1,
      keyword: String(values.keyword ?? '').trim() || undefined,
      store_name: String(values.store_name ?? '').trim() || undefined,
      mer_category_id: Number(values.mer_category_id || 0) || undefined,
    });
    const next: Record<number, number> = { 1: 0, 2: 0, 5: 0, 6: 0, 7: 0 };
    for (const item of (data.list || []) as PlatformProductStatusFilter[]) {
      next[item.type] = item.count;
    }
    tabCounts.value = next;
  } catch {
    /* 统计失败不阻断列表 */
  }
}

const formOptions = computed((): VbenFormProps =>
  listFormOptionsDefaults([
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantCategoryOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'mer_category_id',
      label: '店铺分类',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '请输入店铺名称' },
      fieldName: 'store_name',
      label: '店铺名称',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入商品名称, 关键字, 产品编号',
      },
      fieldName: 'keyword',
      label: '商品搜索',
    },
  ]),
);

const gridOptions: VxeGridProps<PlatformProduct> = {
  checkboxConfig: { highlight: true, range: true },
  columns: [
    { type: 'checkbox', width: 48 },
    { field: 'product_id', title: 'ID', width: 80 },
    {
      align: 'center',
      field: 'image',
      slots: { default: 'image' },
      title: '商品图',
      width: 76,
    },
    {
      field: 'title',
      minWidth: 200,
      showOverflow: false,
      slots: { default: 'title' },
      title: '商品名称',
    },
    {
      field: 'store_name',
      formatter: ({ row }) =>
        row.store_name || row.mer_name || `商户 #${row.mer_id}`,
      minWidth: 140,
      showOverflow: false,
      title: '店铺名称',
    },
    {
      field: 'price',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '商品售价',
      width: 108,
    },
    { field: 'sales', title: '销量', width: 80 },
    { field: 'stock', title: '库存', width: 80 },
    {
      field: 'rank',
      formatter: ({ cellValue }) => String(cellValue ?? 0),
      title: '排序',
      width: 72,
    },
    {
      align: 'center',
      field: 'is_show',
      slots: { default: 'show' },
      title: '是否显示',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 200 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const params = buildListParams(page, formValues);
        const [data] = await Promise.all([
          listPlatformProductsApi(params),
          loadTabCounts(formValues),
        ]);
        selectedIds.value = [];
        gridApi.grid?.clearCheckboxRow?.();
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_id' },
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
  gridEvents: {
    checkboxAll: syncSelection,
    checkboxChange: syncSelection,
  },
  gridOptions,
});

const [ForceOffModal, forceOffModalApi] = useVbenModal({
  title: '强制下架',
  class: 'w-[480px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  destroyOnClose: true,
  onConfirm: async () => submitForceOff(),
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      forceOffReason.value = '';
      forceOffIds.value = [];
      forceOffSubmitting.value = false;
    }
  },
});

const [RejectDrawer, rejectDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确认驳回',
  cancelText: '取消',
  placement: 'right',
  title: '驳回礼包',
  onConfirm: async () => submitReject(),
});

function syncSelection() {
  const rows = (gridApi.grid?.getCheckboxRecords?.() ?? []) as PlatformProduct[];
  selectedIds.value = rows.map((row) => row.product_id);
}

async function reloadGrid() {
  await gridApi.reload();
}

function setStatusTab(type: number) {
  if (tabType.value === type) return;
  tabType.value = type;
  void reloadGrid();
}

function openDetail(row: PlatformProduct) {
  detailDrawerRef.value?.open(row.product_id);
}

async function changeShow(row: PlatformProduct, enabled: boolean) {
  try {
    await setPlatformProductShowApi(row.product_id, enabled ? 1 : 0);
    row.is_show = enabled ? 1 : 0;
    ElMessage.success(enabled ? '已显示' : '已不显示');
  } catch {
    await reloadGrid();
  }
}

function openForceOff(ids: number[]) {
  forceOffIds.value = [...ids];
  forceOffReason.value = '';
  forceOffModalApi.open();
}

function batchForceOff() {
  if (!selectedIds.value.length) {
    ElMessage.warning('请先选择礼包');
    return;
  }
  openForceOff([...selectedIds.value]);
}

async function submitForceOff() {
  const reason = forceOffReason.value.trim();
  if (!reason) {
    ElMessage.warning('请输入强制下架原因');
    return;
  }
  const ids = forceOffIds.value.filter((id) => id > 0);
  if (!ids.length) {
    ElMessage.warning('请选择礼包');
    return;
  }
  forceOffSubmitting.value = true;
  forceOffModalApi.setState({ confirming: true });
  try {
    if (ids.length === 1) {
      await forceOffPlatformProductApi(ids[0]!, reason);
    } else {
      await batchForceOffPlatformProductsApi(ids, reason);
    }
    forceOffModalApi.close();
    ElMessage.success(ids.length === 1 ? '已强制下架' : '批量强制下架成功');
    await reloadGrid();
  } catch {
    /* 全局提示 */
  } finally {
    forceOffSubmitting.value = false;
    forceOffModalApi.setState({ confirming: false });
  }
}

async function approve(row: PlatformProduct) {
  await auditPlatformProductApi(row.product_id, { status: 1 });
  ElMessage.success('审核通过');
  await reloadGrid();
}

function openReject(row: PlatformProduct) {
  current.value = row;
  rejectForm.value.refusal = '';
  rejecting.value = false;
  rejectDrawerApi.open();
}

async function submitReject() {
  const refusal = rejectForm.value.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  rejecting.value = true;
  rejectDrawerApi.lock();
  try {
    await auditPlatformProductApi(current.value.product_id, {
      status: -1,
      refusal,
    });
    rejectDrawerApi.close();
    ElMessage.success('已驳回');
    await reloadGrid();
  } finally {
    rejecting.value = false;
    rejectDrawerApi.unlock();
  }
}

onMounted(async () => {
  const [profile, permissions, merCates] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    fetchMerchantCategories().catch(() => ({ list: [] as MerchantCategoryRow[] })),
  ]);
  canAudit.value =
    profile.roles.includes('platform') &&
    permissions.includes('product.audit.submit');
  merchantCategoryOptions.value = (merCates.list || []).map((c) => ({
    label: c.category_name,
    value: c.merchant_category_id,
  }));
  await loadTabCounts();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="gift-toolbar">
          <div class="gift-status-tabs" role="tablist">
            <button
              v-for="tab in STATUS_TABS"
              :key="tab.type"
              type="button"
              role="tab"
              class="gift-status-tabs__item"
              :aria-selected="tabType === tab.type"
              :class="{ 'is-active': tabType === tab.type }"
              @click="setStatusTab(tab.type)"
            >
              {{ tab.name }}({{ tabCounts[tab.type] || 0 }})
            </button>
          </div>
          <div class="gift-toolbar__actions">
            <ElButton
              v-if="canAudit"
              type="primary"
              :disabled="!selectedIds.length"
              @click="batchForceOff"
            >
              批量强制下架
            </ElButton>
          </div>
        </div>
      </template>

      <template #image="{ row }">
        <ElImage
          v-if="row.image"
          class="gift-list-cover"
          :src="resolveCosMediaUrl(row.image)"
          fit="cover"
          alt="商品图"
        >
          <template #error>
            <span>—</span>
          </template>
        </ElImage>
        <span v-else>—</span>
      </template>

      <template #title="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">
          {{ row.title || '—' }}
        </ElButton>
      </template>

      <template #show="{ row }">
        <ElSwitch
          v-if="canAudit"
          :model-value="Number(row.is_show) === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeShow(row, Boolean(enabled))
          "
        />
        <ElTag v-else :type="Number(row.is_show) === 1 ? 'success' : 'info'">
          {{ Number(row.is_show) === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="canAudit && row.status === 1"
          link
          type="primary"
          @click="openForceOff([row.product_id])"
        >
          强制下架
        </ElButton>
        <template v-if="canAudit && row.status === 0">
          <ElButton link type="primary" @click="approve(row)">通过</ElButton>
          <ElButton link type="primary" @click="openReject(row)">驳回</ElButton>
        </template>
      </template>
    </Grid>

    <ProductDetailDrawer ref="detailDrawerRef" />

    <ForceOffModal>
      <ElInput
        v-model="forceOffReason"
        :autosize="{ minRows: 4, maxRows: 8 }"
        :disabled="forceOffSubmitting"
        maxlength="200"
        placeholder="请输入强制下架原因"
        type="textarea"
      />
    </ForceOffModal>

    <RejectDrawer>
      <ElForm label-width="84px">
        <ElFormItem label="驳回原因" required>
          <ElInput
            v-model="rejectForm.refusal"
            :disabled="rejecting"
            :rows="4"
            maxlength="200"
            placeholder="请向商户说明驳回原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </RejectDrawer>
  </Page>
</template>

<style scoped>
.gift-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.gift-status-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 28px;
  border-bottom: 1px solid hsl(var(--border));
}

.gift-status-tabs__item {
  margin-bottom: -1px;
  padding: 10px 2px 12px;
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: hsl(var(--muted-foreground));
  cursor: pointer;
  font-size: 14px;
  line-height: 20px;
}

.gift-status-tabs__item.is-active {
  border-bottom-color: hsl(var(--primary));
  color: hsl(var(--primary));
  font-weight: 500;
}

.gift-toolbar__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.gift-list-cover {
  width: 48px;
  height: 48px;
  border-radius: 4px;
}
</style>
