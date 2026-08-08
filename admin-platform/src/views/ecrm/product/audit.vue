<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
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
  ElTag,
} from 'element-plus';
import { ArrowDown, Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  auditPlatformProductApi,
  batchForceOffPlatformProductsApi,
  batchCopyPlatformProductsApi,
  batchLabelsPlatformProductsApi,
  batchRecommendPlatformProductsApi,
  batchShowPlatformProductsApi,
  forceOffPlatformProductApi,
  getPlatformProductEditApi,
  getPlatformProductStatusFilterApi,
  listPlatformBrandsApi,
  listPlatformCategoriesApi,
  listPlatformProductsApi,
  setPlatformProductShowApi,
  updatePlatformProductAdminApi,
  type PlatformCategory,
  type PlatformProduct,
  type PlatformProductStatusFilter,
} from '#/api/core/platform-catalog';
import {
  fetchMerchantCategories,
  fetchMerchantTypes,
  fetchProductLabels,
  type MerchantCategoryRow,
  type MerchantTypeRow,
  type ProductLabelRow,
} from '#/api/core/ecrm';
import { fetchActivityLabels } from '#/api/core/platform-product-cache';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import ProductDetailDrawer from './components/ProductDetailDrawer.vue';
import ProductEditDrawer from './components/ProductEditDrawer.vue';
import ProductFictiModal from './components/ProductFictiModal.vue';
import ProductLabelSelectModal from './components/ProductLabelSelectModal.vue';
import ProductPreviewModal from './components/ProductPreviewModal.vue';

const router = useRouter();

const STATUS_TABS: Array<{ type: number; name: string }> = [
  { type: 1, name: '出售中商品' },
  { type: 2, name: '仓库中商品' },
  { type: 6, name: '待审核商品' },
  { type: 7, name: '审核未通过商品' },
  { type: 5, name: '回收站商品' },
];

const tabType = ref(1);
const tabCounts = ref<Record<number, number>>({
  1: 0,
  2: 0,
  5: 0,
  6: 0,
  7: 0,
});
const selectedIds = ref<number[]>([]);
const canAudit = ref(false);
const current = ref<PlatformProduct>();
const rejecting = ref(false);
const rejectForm = reactive({ refusal: '' });
const forceOffReason = ref('');
const forceOffIds = ref<number[]>([]);
const forceOffSubmitting = ref(false);
const lastFormValues = ref<Record<string, unknown>>({});
const previewProductId = ref(0);
const previewProductTitle = ref('');

const categoryOptions = ref<{ label: string; value: number }[]>([]);
const brandOptions = ref<{ label: string; value: string }[]>([]);
const merchantTypeOptions = ref<{ label: string; value: number }[]>([]);
const merchantCategoryOptions = ref<{ label: string; value: number }[]>([]);
const labelOptions = ref<ProductLabelRow[]>([]);
const activityLabelOptions = ref<{ label: string; value: string }[]>([]);

const detailDrawerRef = ref<InstanceType<typeof ProductDetailDrawer>>();
const editDrawerRef = ref<InstanceType<typeof ProductEditDrawer>>();
const previewModalRef = ref<InstanceType<typeof ProductPreviewModal>>();
const labelModalRef = ref<InstanceType<typeof ProductLabelSelectModal>>();
const fictiModalRef = ref<InstanceType<typeof ProductFictiModal>>();
const labelEditingProductId = ref(0);

function flattenCategories(
  nodes: PlatformCategory[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    {
      label: `${prefix}${node.cate_name}`,
      value: node.store_category_id,
    },
    ...flattenCategories(node.children || [], `${prefix}— `),
  ]);
}

function statusInfo(status: number) {
  return (
    (
      {
        '-2': { label: '已下架', type: 'info' },
        '-1': { label: '审核未通过', type: 'danger' },
        0: { label: '待审核', type: 'warning' },
        1: { label: '出售中', type: 'success' },
      } as Record<
        number,
        { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
      >
    )[status] || { label: '未知', type: 'info' as const }
  );
}

function productTypeLabel(type?: number) {
  return (
    (
      {
        0: '普通商品',
        1: '虚拟',
        2: '云盘',
        3: '卡密',
        4: '预约',
        5: '年/次卡',
      } as Record<number, string>
    )[Number(type ?? 0)] || '普通商品'
  );
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  lastFormValues.value = values;
  const cateId = Number(values.cate_id || 0);
  const merTypeId = Number(values.mer_type_id || 0);
  const merCategoryId = Number(values.mer_category_id || 0);
  const usStatus = values.us_status;
  const star = values.star;
  const svip = values.svip_price_type;
  const productType = values.product_type;
  const isHot = values.is_hot;
  const cateHot = values.cate_hot;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    type: tabType.value,
    keyword: String(values.keyword ?? '').trim() || undefined,
    store_name: String(values.store_name ?? '').trim() || undefined,
    brand_name: String(values.brand_name ?? '').trim() || undefined,
    cate_id: cateId > 0 ? cateId : undefined,
    mer_type_id: merTypeId > 0 ? merTypeId : undefined,
    mer_category_id: merCategoryId > 0 ? merCategoryId : undefined,
    us_status:
      usStatus === 0 ||
      usStatus === 1 ||
      usStatus === -1 ||
      usStatus === -2
        ? Number(usStatus)
        : undefined,
    star: star === 0 || star ? Number(star) : undefined,
    svip_price_type: svip === 0 || svip === 1 || svip === 2 ? Number(svip) : undefined,
    product_type:
      productType === 0 || productType
        ? Number(productType)
        : undefined,
    is_hot: isHot === 0 || isHot === 1 ? Number(isHot) : undefined,
    cate_hot: cateHot === 0 || cateHot === 1 ? Number(cateHot) : undefined,
  };
}

async function loadTabCounts(formValues?: Record<string, unknown>) {
  const values = formValues || lastFormValues.value || {};
  try {
    const data = await getPlatformProductStatusFilterApi({
      keyword: String(values.keyword ?? '').trim() || undefined,
      store_name: String(values.store_name ?? '').trim() || undefined,
      brand_name: String(values.brand_name ?? '').trim() || undefined,
      cate_id: Number(values.cate_id || 0) || undefined,
      mer_type_id: Number(values.mer_type_id || 0) || undefined,
      mer_category_id: Number(values.mer_category_id || 0) || undefined,
      svip_price_type:
        values.svip_price_type === 0 ||
        values.svip_price_type === 1 ||
        values.svip_price_type === 2
          ? Number(values.svip_price_type)
          : undefined,
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
        options: categoryOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'cate_id',
      label: '商品分类',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '店铺名称' },
      fieldName: 'store_name',
      label: '店铺名称',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantCategoryOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'mer_category_id',
      label: '店铺类别',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '出售中', value: 1 },
          { label: '仓库中', value: -2 },
          { label: '待审核', value: 0 },
          { label: '审核未通过', value: -1 },
        ],
        placeholder: '全部',
      },
      fieldName: 'us_status',
      label: '商品状态',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: brandOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'brand_name',
      label: '品牌选择',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '全部', value: '' },
          ...[0, 1, 2, 3, 4, 5].map((n) => ({
            label: n === 0 ? '未设置' : `${n}星`,
            value: n,
          })),
        ],
        placeholder: '全部',
      },
      fieldName: 'star',
      label: '推荐级别',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: labelOptions.value.map((x) => ({
          label: x.name,
          value: String(x.id),
        })),
        placeholder: '请选择（筛选项预留）',
      },
      fieldName: 'sys_labels',
      label: '商品标签',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: activityLabelOptions.value,
        placeholder: '请选择（筛选项预留）',
      },
      fieldName: 'activity_label_ids',
      label: '活动标签',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '无会员价', value: 0 },
          { label: '默认会员价', value: 1 },
          { label: '自定义会员价', value: 2 },
        ],
        placeholder: '全部',
      },
      fieldName: 'svip_price_type',
      label: '会员价',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '是', value: 1 },
          { label: '否', value: 0 },
        ],
        placeholder: '全部',
      },
      fieldName: 'is_hot',
      label: '商品推荐',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '是', value: 1 },
          { label: '否', value: 0 },
        ],
        placeholder: '全部',
      },
      fieldName: 'cate_hot',
      label: '大图推荐',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '普通商品', value: 0 },
          { label: '虚拟', value: 1 },
          { label: '云盘', value: 2 },
          { label: '卡密', value: 3 },
          { label: '预约', value: 4 },
          { label: '年/次卡', value: 5 },
        ],
        placeholder: '全部',
      },
      fieldName: 'product_type',
      label: '商品类型',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '名称 / 关键字 / 货号',
      },
      fieldName: 'keyword',
      label: '商品搜索',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantTypeOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'mer_type_id',
      label: '店铺类型',
    },
  ]),
);

const gridOptions: VxeGridProps<PlatformProduct> = {
  checkboxConfig: { highlight: true, range: true },
  expandConfig: {
    padding: true,
  },
  columns: [
    { type: 'checkbox', width: 48 },
    {
      type: 'expand',
      width: 46,
      slots: { content: 'expandContent' },
    },
    { field: 'product_id', title: 'ID', width: 80 },
    {
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
      minWidth: 130,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) =>
        row.store_name || row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'price',
      title: '商品售价',
      width: 108,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'sales', title: '销量', width: 80 },
    { field: 'stock', title: '库存', width: 80 },
    {
      field: 'star',
      slots: { default: 'star' },
      title: '推荐级别',
      width: 130,
    },
    {
      field: 'rank',
      title: '排序',
      width: 72,
      formatter: ({ cellValue }) => String(cellValue ?? 0),
    },
    {
      field: 'is_show',
      slots: { default: 'show' },
      title: '是否显示',
      width: 96,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '商品状态',
      width: 104,
    },
    platformListActionColumn({ width: 168 }),
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

const [RejectDrawer, rejectDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确认驳回',
  cancelText: '取消',
  placement: 'right',
  title: '驳回商品',
  onConfirm: async () => submitReject(),
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

function syncSelection() {
  const rows = (gridApi.grid?.getCheckboxRecords?.() ?? []) as PlatformProduct[];
  selectedIds.value = rows.map((row) => row.product_id);
}

function requireSelection() {
  if (!selectedIds.value.length) {
    ElMessage.warning('请先选择商品');
    return false;
  }
  return true;
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

function openPreview(row: PlatformProduct) {
  previewProductId.value = row.product_id;
  previewProductTitle.value = row.title || '';
  previewModalRef.value?.open();
}

function openEdit(row: PlatformProduct | number) {
  const id = typeof row === 'number' ? row : row.product_id;
  editDrawerRef.value?.open(id);
}

function openCreate() {
  editDrawerRef.value?.openCreate();
}

function onDetailEdit(productId: number) {
  openEdit(productId);
}

async function onEditSaved() {
  await reloadGrid();
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

async function submitForceOff() {
  const reason = forceOffReason.value.trim();
  if (!reason) {
    ElMessage.warning('请输入强制下架原因');
    return;
  }
  const ids = forceOffIds.value.filter((id) => id > 0);
  if (!ids.length) {
    ElMessage.warning('请选择商品');
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
    /* 失败由全局拦截提示 */
  } finally {
    forceOffSubmitting.value = false;
    forceOffModalApi.setState({ confirming: false });
  }
}

function forceOff(row: PlatformProduct) {
  openForceOff([row.product_id]);
}

async function approve(row: PlatformProduct) {
  try {
    await ElMessageBox.confirm(`确认审核通过商品「${row.title}」？`, '商品审核', {
      type: 'warning',
    });
    await auditPlatformProductApi(row.product_id, { status: 1 });
    ElMessage.success('商品已审核通过');
    await reloadGrid();
  } catch {
    /* 取消或失败 */
  }
}

function openReject(row: PlatformProduct) {
  current.value = row;
  rejectForm.refusal = '';
  rejectDrawerApi.open();
}

async function submitReject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  rejecting.value = true;
  try {
    await auditPlatformProductApi(current.value.product_id, {
      status: -1,
      refusal,
    });
    rejectDrawerApi.close();
    ElMessage.success('商品已驳回');
    await reloadGrid();
  } finally {
    rejecting.value = false;
  }
}

function batchForceOff() {
  if (!requireSelection()) return;
  openForceOff([...selectedIds.value]);
}

async function batchShow(status: 0 | 1) {
  if (!requireSelection()) return;
  try {
    await batchShowPlatformProductsApi(selectedIds.value, status);
    ElMessage.success(status === 1 ? '批量展示成功' : '批量不显示成功');
    await reloadGrid();
  } catch {
    /* 失败 */
  }
}

async function batchLabels() {
  if (!requireSelection()) return;
  try {
    await batchLabelsPlatformProductsApi(selectedIds.value, '');
  } catch {
    ElMessage.warning('批量设置标签尚未接入（TODO）');
  }
}

async function batchRecommend() {
  if (!requireSelection()) return;
  try {
    await batchRecommendPlatformProductsApi(selectedIds.value, { is_hot: 1 });
  } catch {
    ElMessage.warning('批量设置推荐尚未接入（TODO），请用行内「编辑商品」设置星级');
  }
}

async function batchCopy() {
  if (!requireSelection()) return;
  try {
    await batchCopyPlatformProductsApi(selectedIds.value, 0);
  } catch {
    ElMessage.warning('批量复制商品尚未接入（TODO）');
  }
}

function openComments(row: PlatformProduct) {
  void router.push({
    path: '/product/comment',
    query: { keyword: String(row.product_id) },
  });
}

function openFicti(row: PlatformProduct) {
  fictiModalRef.value?.open({
    productId: row.product_id,
    ficti: Number(row.ficti || 0),
  });
}

async function onFictiSuccess() {
  await reloadGrid();
  detailDrawerRef.value?.reload?.();
}

async function openEditLabels(row: PlatformProduct) {
  labelEditingProductId.value = row.product_id;
  try {
    const [edit, labels] = await Promise.all([
      getPlatformProductEditApi(row.product_id),
      labelOptions.value.length
        ? Promise.resolve({ list: labelOptions.value })
        : fetchProductLabels().catch(() => ({ list: [] as ProductLabelRow[] })),
    ]);
    if (!labelOptions.value.length) {
      labelOptions.value = labels.list || [];
    }
    labelModalRef.value?.open({
      productId: row.product_id,
      selectedIds: [...(edit.sys_labels || [])].map(String),
      options: labelOptions.value.filter((x) => Number(x.status) !== 0),
    });
  } catch {
    ElMessage.error('加载商品标签失败');
  }
}

async function onLabelSubmit(ids: string[]) {
  const productId = labelEditingProductId.value;
  if (!productId) return;
  try {
    await updatePlatformProductAdminApi(productId, { sys_labels: ids });
    ElMessage.success('标签已更新');
    gridApi.reload();
  } catch {
    /* 统一错误提示 */
  }
}

async function onMoreCommand(command: string, row: PlatformProduct) {
  switch (command) {
    case 'edit':
      openEdit(row);
      break;
    case 'comments':
      openComments(row);
      break;
    case 'labels':
      openEditLabels(row);
      break;
    case 'ficti':
      openFicti(row);
      break;
    case 'forceOff':
      await forceOff(row);
      break;
    case 'approve':
      await approve(row);
      break;
    case 'reject':
      openReject(row);
      break;
    default:
      break;
  }
}

onMounted(async () => {
  const [profile, permissions, categories, brands, types, merCates, labels, actLabels] =
    await Promise.all([
      getUserInfoApi(),
      getAccessCodesApi(),
      listPlatformCategoriesApi().catch(() => ({ list: [] as PlatformCategory[] })),
      listPlatformBrandsApi().catch(() => ({ list: [] })),
      fetchMerchantTypes().catch(() => ({ list: [] as MerchantTypeRow[] })),
      fetchMerchantCategories().catch(() => ({ list: [] as MerchantCategoryRow[] })),
      fetchProductLabels().catch(() => ({ list: [] as ProductLabelRow[] })),
      fetchActivityLabels().catch(() => ({ list: [] })),
    ]);
  canAudit.value =
    profile.roles.includes('platform') &&
    permissions.includes('product.audit.submit');
  categoryOptions.value = flattenCategories(categories.list || []);
  brandOptions.value = (brands.list || []).map((b) => ({
    label: b.brand_name,
    value: b.brand_name,
  }));
  merchantTypeOptions.value = (types.list || []).map((t) => ({
    label: t.name || String(t.id),
    value: t.id,
  }));
  merchantCategoryOptions.value = (merCates.list || []).map((c) => ({
    label: c.category_name,
    value: c.merchant_category_id,
  }));
  labelOptions.value = labels.list || [];
  activityLabelOptions.value = (actLabels.list || []).map((x) => ({
    label: x.name,
    value: x.id,
  }));
  await loadTabCounts();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="product-toolbar">
          <div class="product-status-tabs" role="tablist">
            <button
              v-for="tab in STATUS_TABS"
              :key="tab.type"
              type="button"
              role="tab"
              class="product-status-tabs__item"
              :aria-selected="tabType === tab.type"
              :class="{ 'is-active': tabType === tab.type }"
              @click="setStatusTab(tab.type)"
            >
              {{ tab.name }}({{ tabCounts[tab.type] || 0 }})
            </button>
          </div>
          <div class="product-toolbar__actions">
            <ElButton
              v-if="canAudit"
              :icon="Plus"
              type="primary"
              @click="openCreate"
            >
              新增商品
            </ElButton>
            <ElButton
              v-if="canAudit"
              :disabled="!selectedIds.length"
              @click="batchForceOff"
            >
              批量强制下架
            </ElButton>
            <ElButton
              v-if="canAudit"
              :disabled="!selectedIds.length"
              @click="batchShow(0)"
            >
              批量不显示
            </ElButton>
            <ElButton
              v-if="canAudit"
              :disabled="!selectedIds.length"
              @click="batchShow(1)"
            >
              批量展示
            </ElButton>
            <ElButton :disabled="!selectedIds.length" @click="batchLabels">
              批量设置标签
            </ElButton>
            <ElButton :disabled="!selectedIds.length" @click="batchRecommend">
              批量设置推荐
            </ElButton>
            <ElButton :disabled="!selectedIds.length" @click="batchCopy">
              批量复制商品
            </ElButton>
          </div>
        </div>
      </template>

      <template #expandContent="{ row }">
        <div class="product-expand">
          <div class="product-expand__item">
            <span class="label">平台分类</span>
            <span class="value">{{ row.cate_name || '—' }}</span>
          </div>
          <div class="product-expand__item">
            <span class="label">商品分类</span>
            <span class="value">{{ row.mer_cate_name || '—' }}</span>
          </div>
          <div class="product-expand__item">
            <span class="label">品牌</span>
            <span class="value">{{ row.brand_name || '—' }}</span>
          </div>
          <div class="product-expand__item">
            <span class="label">划线价</span>
            <span class="value">
              ¥{{ Number(row.ot_price || 0).toFixed(2) }}
            </span>
          </div>
          <div class="product-expand__item">
            <span class="label">收藏</span>
            <span class="value">{{ Number(row.care_count || 0) }}</span>
          </div>
          <div class="product-expand__item">
            <span class="label">已售数量</span>
            <span class="value">
              <ElButton link type="primary" @click="openFicti(row)">
                {{ Number(row.ficti || 0) }}
              </ElButton>
            </span>
          </div>
          <div class="product-expand__item product-expand__item--tags">
            <span class="label">活动标签</span>
            <span class="value">
              <template v-if="row.activity_labels?.length">
                <ElTag
                  v-for="tag in row.activity_labels"
                  :key="tag"
                  class="product-expand__tag"
                  size="small"
                  type="danger"
                >
                  {{ tag }}
                </ElTag>
              </template>
              <template v-else>—</template>
            </span>
          </div>
        </div>
      </template>

      <template #image="{ row }">
        <ElImage
          v-if="row.image"
          class="product-list-cover"
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
        <div class="product-title-cell">
          <ElButton link type="primary" @click="openDetail(row)">
            {{ row.title || '—' }}
          </ElButton>
          <div class="product-title-cell__tags">
            <ElTag v-if="row.spec_type === 1" size="small" type="warning">
              多规格
            </ElTag>
            <ElTag v-if="Number(row.svip_price_type) > 0" size="small">
              会员价
            </ElTag>
            <ElTag size="small" type="info">
              {{ productTypeLabel(row.product_type) }}
            </ElTag>
          </div>
        </div>
      </template>

      <template #star="{ row }">
        <ElRate
          :model-value="Number(row.star || 0)"
          disabled
          :max="5"
          class="product-star"
        />
      </template>

      <template #show="{ row }">
        <ElSwitch
          v-if="canAudit"
          :model-value="Number(row.is_show) === 1"
          @change="(enabled: string | number | boolean) => changeShow(row, Boolean(enabled))"
        />
        <ElTag v-else :type="Number(row.is_show) === 1 ? 'success' : 'info'">
          {{ Number(row.is_show) === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>

      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openPreview(row)">预览</ElButton>
        <ElDropdown
          trigger="click"
          @command="(cmd: string) => onMoreCommand(cmd, row)"
        >
          <ElButton link type="primary">
            更多
            <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem v-if="canAudit" command="edit">
                编辑商品
              </ElDropdownItem>
              <ElDropdownItem command="comments">查看评价</ElDropdownItem>
              <ElDropdownItem command="labels">编辑标签</ElDropdownItem>
              <ElDropdownItem command="ficti">已售数量</ElDropdownItem>
              <ElDropdownItem
                v-if="canAudit && row.status === 1"
                command="forceOff"
                divided
              >
                强制下架
              </ElDropdownItem>
              <template v-if="canAudit && row.status === 0">
                <ElDropdownItem command="approve" divided>
                  审核通过
                </ElDropdownItem>
                <ElDropdownItem command="reject">审核驳回</ElDropdownItem>
              </template>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <ProductDetailDrawer
      ref="detailDrawerRef"
      @edit="onDetailEdit"
      @ficti-updated="onFictiSuccess"
    />
    <ProductEditDrawer ref="editDrawerRef" @saved="onEditSaved" />
    <ProductLabelSelectModal ref="labelModalRef" @submit="onLabelSubmit" />
    <ProductFictiModal ref="fictiModalRef" @success="onFictiSuccess" />
    <ProductPreviewModal
      ref="previewModalRef"
      :product-id="previewProductId"
      :product-title="previewProductTitle"
    />

    <RejectDrawer>
      <ElForm label-width="84px">
        <ElFormItem label="驳回原因" required>
          <ElInput
            v-model="rejectForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请向商户说明驳回原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </RejectDrawer>

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
  </Page>
</template>

<style scoped>
.product-expand {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 28px;
  align-items: flex-start;
  width: 100%;
  padding: 12px 16px;
  background: hsl(var(--muted) / 45%);
  border-radius: 4px;
  box-sizing: border-box;
}

.product-expand__item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 110px;
}

.product-expand__item--tags {
  min-width: 160px;
  flex: 1;
}

.product-expand__item .label {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 18px;
}

.product-expand__item .value {
  color: hsl(var(--foreground));
  font-size: 13px;
  line-height: 20px;
  word-break: break-all;
}

.product-expand__tag {
  margin-right: 4px;
  margin-bottom: 2px;
}

.product-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.product-status-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 28px;
  border-bottom: 1px solid hsl(var(--border));
}

.product-status-tabs__item {
  margin-bottom: -1px;
  padding: 10px 2px 12px;
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: hsl(var(--foreground) / 70%);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
}

.product-status-tabs__item:hover {
  color: hsl(var(--primary));
}

.product-status-tabs__item.is-active {
  border-bottom-color: hsl(var(--primary));
  color: hsl(var(--primary));
  font-weight: 600;
}

.product-toolbar__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-start;
}

.product-list-cover,
.product-detail-cover {
  width: 48px;
  height: 48px;
  border-radius: 4px;
}

.product-detail-cover {
  width: 96px;
  height: 96px;
}

.product-title-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-start;
}

.product-title-cell__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.product-star {
  --el-rate-icon-size: 14px;
  height: 20px;
}

.product-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-preview__cover {
  width: 100%;
  max-height: 320px;
  border-radius: 8px;
}

.product-preview__title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.product-preview__price {
  margin: 0;
  color: hsl(var(--destructive));
  font-size: 20px;
  font-weight: 600;
}

.product-preview__meta,
.product-preview__info,
.edit-hint {
  margin: 0;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
  line-height: 1.6;
}
</style>
