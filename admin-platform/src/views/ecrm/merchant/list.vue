<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref, watch } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElCascader,
  ElCheckbox,
  ElCheckboxGroup,
  ElDatePicker,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElPagination,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElImage,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
} from 'element-plus';
import { MoreFilled, Plus, Shop } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import {
  createPlatformMerchant,
  fetchBusinessZoneOptions,
  fetchMerchantCategories,
  fetchMerchantOperateLogs,
  fetchMerchantTypes,
  fetchPlatformMerchant,
  fetchPlatformMerchants,
  fetchStoreGroups,
  updatePlatformMerchant,
  updatePlatformMerchantRecommend,
  updatePlatformMerchantStatus,
  type BusinessZoneOptionNode,
  type MerchantCategoryRow,
  type MerchantOperateLogRow,
  type MerchantTypeRow,
  type PlatformMerchantRow,
  type PlatformMerchantSaveInput,
  type StoreGroupRow,
} from '#/api/core/ecrm';
import {
  listPlatformCategoriesApi,
  type PlatformCategory,
} from '#/api/core/platform-catalog';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import ImageField from '#/components/shop/image-field.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';

type DrawerMode = 'create' | 'edit' | 'view';

const GOODS_TYPE_OPTIONS = [
  { label: '普通商品', value: 0 },
  { label: '虚拟', value: 1 },
  { label: '云盘', value: 2 },
  { label: '卡密', value: 3 },
  { label: '预约', value: 4 },
  { label: '年/次卡', value: 5 },
] as const;

const STAR_OPTIONS = [1, 2, 3, 4, 5].map((n) => ({
  label: `${n}星`,
  value: n,
}));

const openCount = ref(0);
const closedCount = ref(0);
const statusFilter = ref<number>(1);
const canManage = ref(false);
const categories = ref<MerchantCategoryRow[]>([]);
const types = ref<MerchantTypeRow[]>([]);
const regions = ref<{ label: string; value: number }[]>([]);
const businessOptions = ref<BusinessZoneOptionNode[]>([]);
const storeGroupOptions = ref<{ label: string; value: number }[]>([]);
/** Cascader 用标准 label/value/children，避免空 children[] 被当成可展开父节点。 */
type PlatformCategoryCascaderOption = {
  label: string;
  value: number;
  children?: PlatformCategoryCascaderOption[];
};
const platformCategoryTree = ref<PlatformCategoryCascaderOption[]>([]);
const platformCategoryOptions = ref<{ label: string; value: number }[]>([]);
const platformCategoryCascaderProps = {
  multiple: true,
  checkStrictly: false,
  emitPath: false,
  value: 'value',
  label: 'label',
  children: 'children',
};
const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const activeTab = ref('basic');
const headerSnapshot = ref<Partial<PlatformMerchantRow>>({});
const operateLogs = ref<MerchantOperateLogRow[]>([]);
const operateLogTotal = ref(0);
const operateLogLoading = ref(false);
const operateLogPage = ref(1);
const operateLogLimit = ref(10);
const operateLogTerminal = ref<string>('');
const operateLogDates = ref<string[]>([]);

const form = reactive<
  PlatformMerchantSaveInput & { mer_password?: string; store_group_ids: number[] }
>({
  mer_name: '',
  business_id: undefined,
  real_name: '',
  mer_phone: '',
  mer_address: '',
  mer_info: '',
  mer_keyword: '',
  mark: '',
  category_id: undefined,
  type_id: undefined,
  region_id: undefined,
  is_best: false,
  offline_pay: false,
  is_trader: false,
  is_audit: true,
  is_bro_room: false,
  is_bro_goods: false,
  commission_switch: false,
  commission_rate: 0,
  mer_account: '',
  mer_password: '',
  sub_mchid: '',
  applyment_id: '',
  care_count: 0,
  care_ficti: 0,
  sort: 0,
  status: true,
  store_group_ids: [],
  goods_types: [] as number[],
  platform_category_ids: [] as number[],
  mer_star: 5,
  mer_avatar: '',
});

const isReadonly = computed(() => drawerMode.value === 'view');
const selectedType = computed(() =>
  types.value.find((t) => t.id === Number(form.type_id)),
);
const selectedCategory = computed(() =>
  categories.value.find(
    (c) => c.merchant_category_id === Number(form.category_id),
  ),
);
const selectedBusinessLabel = computed(() => {
  const id = Number(form.business_id || 0);
  if (!id) return headerSnapshot.value.owner_name || '—';
  const hit = flattenZoneOptions(businessOptions.value).find((x) => x.value === id);
  return hit?.label || headerSnapshot.value.owner_name || '—';
});
const selectedRegionLabel = computed(() => {
  const id = Number(form.region_id || 0);
  if (!id) return headerSnapshot.value.region_name || '—';
  return (
    regions.value.find((r) => r.value === id)?.label ||
    headerSnapshot.value.region_name ||
    '—'
  );
});
const selectedStoreGroupLabels = computed(() => {
  const ids = new Set((form.store_group_ids || []).map(Number));
  if (!ids.size) return '—';
  const labels = storeGroupOptions.value
    .filter((item) => ids.has(item.value))
    .map((item) => item.label);
  return labels.length ? labels.join('，') : '—';
});
const selectedPlatformCategoryLabels = computed(() => {
  const ids = new Set((form.platform_category_ids || []).map(Number));
  if (!ids.size) return '—';
  const labels = platformCategoryOptions.value
    .filter((item) => ids.has(item.value))
    .map((item) => item.label);
  return labels.length ? labels.join('，') : '—';
});
const selectedGoodsTypeLabels = computed(() => {
  const ids = new Set((form.goods_types || []).map(Number));
  if (!ids.size) return '—';
  const labels = GOODS_TYPE_OPTIONS.filter((item) => ids.has(item.value)).map(
    (item) => item.label,
  );
  return labels.length ? labels.join('，') : '—';
});
const categoryCommissionHint = computed(() => {
  const rate = selectedCategory.value?.commission_rate;
  if (rate === undefined || rate === null) return '';
  return `该分类下的店铺手续费是${Number(rate)}%`;
});
const depositStatusText = computed(() => {
  switch (headerSnapshot.value.deposit_state) {
    case 'funded':
      return '已缴';
    case 'pending':
      return '待缴';
    case 'shortfall':
      return '已缴';
    case 'refund_pending':
      return '退款中';
    case 'refunded':
      return '已退';
    case 'not_required':
    default:
      return depositText(headerSnapshot.value.deposit_state);
  }
});
const depositNeedsTopUp = computed(
  () => headerSnapshot.value.deposit_state === 'shortfall',
);
const marginAmountText = computed(() => {
  if (!(selectedType.value?.is_margin || headerSnapshot.value.type_is_margin)) {
    return '无';
  }
  const amount =
    selectedType.value?.margin ?? headerSnapshot.value.type_margin ?? 0;
  return Number(amount).toFixed(2);
});

const isEnabled = (row: PlatformMerchantRow) =>
  row.status === 1 && row.mer_state === 1;
const statusText = (row: Partial<PlatformMerchantRow>) =>
  row.status === 1 ? '开启' : '关闭';
const formatTime = (value?: string) =>
  value ? formatShanghaiDateTime(value) : '—';
const depositText = (state?: string) => {
  switch (state) {
    case 'funded':
      return '已缴';
    case 'pending':
      return '待缴';
    case 'shortfall':
      return '不足';
    case 'refund_pending':
      return '退款中';
    case 'refunded':
      return '已退';
    case 'not_required':
    default:
      return '无';
  }
};
const yesNo = (value: boolean | number | undefined) =>
  value === true || value === 1 ? '是' : '否';
const auditText = (needAudit: boolean) => (needAudit ? '需审核' : '免审核');
const displayOrDash = (value?: string | number | null) => {
  if (value === undefined || value === null || value === '') return '—';
  return String(value);
};

function flattenStoreGroups(
  rows: StoreGroupRow[],
  prefix = '',
): { label: string; value: number }[] {
  const out: { label: string; value: number }[] = [];
  for (const row of rows || []) {
    const label = prefix ? `${prefix} / ${row.name}` : row.name;
    out.push({ label, value: Number(row.id) });
    if (row.children?.length) {
      out.push(...flattenStoreGroups(row.children, label));
    }
  }
  return out;
}

function toPlatformCategoryCascaderOptions(
  rows: PlatformCategory[] = [],
): PlatformCategoryCascaderOption[] {
  const out: PlatformCategoryCascaderOption[] = [];
  for (const row of rows || []) {
    const value = Number(row.store_category_id);
    const label = String(row.cate_name || '').trim();
    if (!Number.isFinite(value) || value <= 0 || !label) continue;
    // 店铺表单只展示启用分类；隐藏节点仍保留其子树中已启用的后代。
    const children = toPlatformCategoryCascaderOptions(row.children || []);
    if (Number(row.is_show) === 0 && !children.length) continue;
    const option: PlatformCategoryCascaderOption = { label, value };
    if (children.length) option.children = children;
    if (Number(row.is_show) !== 0) {
      out.push(option);
    } else {
      out.push(...children);
    }
  }
  return out;
}

function flattenPlatformCategories(
  rows: PlatformCategoryCascaderOption[],
  prefix = '',
): { label: string; value: number }[] {
  const out: { label: string; value: number }[] = [];
  for (const row of rows || []) {
    const label = prefix ? `${prefix} / ${row.label}` : row.label;
    out.push({ label, value: Number(row.value) });
    if (row.children?.length) {
      out.push(...flattenPlatformCategories(row.children, label));
    }
  }
  return out;
}

async function loadPlatformCategories() {
  try {
    const res = await listPlatformCategoriesApi();
    const list = Array.isArray(res?.list)
      ? res.list
      : Array.isArray(res)
        ? (res as PlatformCategory[])
        : [];
    platformCategoryTree.value = toPlatformCategoryCascaderOptions(list);
    platformCategoryOptions.value = flattenPlatformCategories(
      platformCategoryTree.value,
    );
  } catch {
    // 保留已有选项；打开抽屉时会再试一次
  }
}

function parseGoodsTypes(row?: PlatformMerchantRow): number[] {
  if (row?.goods_types?.length) {
    return [...row.goods_types];
  }
  const raw = row?.goods_type?.trim();
  if (!raw) return [];
  return raw
    .split(',')
    .map((part) => Number(part.trim()))
    .filter((value) => Number.isFinite(value));
}

function parsePlatformCategoryIds(row?: PlatformMerchantRow): number[] {
  if (row?.platform_category_id_list?.length) {
    return row.platform_category_id_list.map((id) => Number(id));
  }
  const raw = row?.platform_category_ids?.trim();
  if (!raw) return [];
  return raw
    .split(',')
    .map((part) => Number(part.trim()))
    .filter((value) => Number.isFinite(value));
}

async function loadCounts() {
  const [opened, closed] = await Promise.all([
    fetchPlatformMerchants({ page: 1, limit: 1, status: 1 }),
    fetchPlatformMerchants({ page: 1, limit: 1, status: 0 }),
  ]);
  openCount.value = opened.total || 0;
  closedCount.value = closed.total || 0;
}

function flattenZoneOptions(
  nodes: BusinessZoneOptionNode[],
): { label: string; value: number }[] {
  const out: { label: string; value: number }[] = [];
  const walk = (list: BusinessZoneOptionNode[]) => {
    for (const node of list || []) {
      out.push({ label: node.label, value: Number(node.value) });
      if (node.children?.length) walk(node.children);
    }
  };
  walk(nodes);
  return out;
}

async function loadFilterOptions() {
  try {
    const [cats, typeList, regionTree, businessTree, groups] = await Promise.all([
      fetchMerchantCategories().catch(() => ({ list: [] as MerchantCategoryRow[] })),
      fetchMerchantTypes().catch(() => ({ list: [] as MerchantTypeRow[] })),
      fetchBusinessZoneOptions(0).catch(() => ({ list: [] as BusinessZoneOptionNode[] })),
      fetchBusinessZoneOptions(1).catch(() => ({ list: [] as BusinessZoneOptionNode[] })),
      fetchStoreGroups().catch(() => ({ list: [] as StoreGroupRow[] })),
      loadPlatformCategories(),
    ]);
    categories.value = cats.list || [];
    types.value = typeList.list || [];
    regions.value = flattenZoneOptions(regionTree.list || []);
    businessOptions.value = businessTree.list || [];
    storeGroupOptions.value = flattenStoreGroups(groups.list || []);
  } catch {
    /* 筛选项失败不阻断列表 */
  }
}

function yesNoOptions() {
  return [
    { label: '是', value: 1 },
    { label: '否', value: 0 },
  ];
}

const formOptions: VbenFormProps = {
  collapsed: false,
  commonConfig: { componentProps: { class: 'w-full' } },
  showCollapseButton: false,
  wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
  schema: [
    {
      component: 'DatePicker',
      componentProps: {
        type: 'daterange',
        valueFormat: 'YYYY-MM-DD',
        startPlaceholder: '开始时间',
        endPlaceholder: '结束时间',
      },
      fieldName: 'date_range',
      label: '选择时间',
    },
    {
      component: 'Select',
      componentProps: { clearable: true, options: [], placeholder: '请选择' },
      fieldName: 'category_id',
      label: '店铺分类',
    },
    {
      component: 'Select',
      componentProps: { clearable: true, options: [], placeholder: '请选择' },
      fieldName: 'type_id',
      label: '店铺类型',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: yesNoOptions(),
        placeholder: '请选择',
      },
      fieldName: 'is_best',
      label: '是否推荐',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: yesNoOptions(),
        placeholder: '请选择',
      },
      fieldName: 'offline_pay',
      label: '线下支付',
    },
    {
      component: 'Select',
      componentProps: { clearable: true, options: [], placeholder: '请选择' },
      fieldName: 'region_id',
      label: '店铺区域',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入店铺名称/联系人/联系电话',
      },
      fieldName: 'keyword',
      label: '店铺搜索',
    },
  ],
};

const gridOptions: VxeGridProps<PlatformMerchantRow> = {
  columns: [
    { field: 'mer_id', title: 'ID', width: 72 },
    {
      field: 'mer_avatar',
      slots: { default: 'cover' },
      title: '封面',
      width: 80,
    },
    {
      field: 'mer_name',
      minWidth: 140,
      showOverflow: false,
      slots: { default: 'mer_name' },
      title: '店铺名称',
    },
    {
      field: 'owner_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      showOverflow: false,
      title: '所属商户',
    },
    {
      field: 'real_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 100,
      showOverflow: false,
      title: '联系人',
    },
    {
      field: 'mer_phone',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      showOverflow: false,
      title: '联系手机号',
    },
    {
      field: 'type_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      showOverflow: false,
      title: '店铺类型',
    },
    {
      field: 'deposit_state',
      formatter: ({ cellValue }) => depositText(cellValue),
      title: '保证金',
      width: 88,
    },
    {
      field: 'is_best',
      slots: { default: 'recommend' },
      title: '推荐',
      width: 90,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '开启/关闭',
      width: 100,
    },
    { field: 'sort', title: '排序', width: 72 },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 160,
      showOverflow: false,
      title: '创建时间',
    },
    platformListActionColumn({ width: 168 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await fetchPlatformMerchants({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status: statusFilter.value,
          category_id: formValues?.category_id
            ? Number(formValues.category_id)
            : undefined,
          type_id: formValues?.type_id ? Number(formValues.type_id) : undefined,
          region_id: formValues?.region_id
            ? Number(formValues.region_id)
            : undefined,
          is_best:
            formValues?.is_best === 0 || formValues?.is_best === 1
              ? Number(formValues.is_best)
              : undefined,
          offline_pay:
            formValues?.offline_pay === 0 || formValues?.offline_pay === 1
              ? Number(formValues.offline_pay)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'mer_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [ShopDrawer, shopDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (drawerMode.value === 'view') {
      switchToEditFromDetail();
      return;
    }
    if (!form.mer_name?.trim()) {
      ElMessage.warning('请填写店铺名称');
      activeTab.value = 'basic';
      return;
    }
    if (!form.mer_avatar?.trim()) {
      ElMessage.warning('请从素材库选择店铺封面');
      activeTab.value = 'basic';
      return;
    }
    if (!form.category_id) {
      ElMessage.warning('请选择店铺分类');
      activeTab.value = 'basic';
      return;
    }
    if (!form.type_id) {
      ElMessage.warning('请选择店铺类型');
      activeTab.value = 'basic';
      return;
    }
    if (!form.goods_types?.length) {
      ElMessage.warning('请选择商品类型');
      activeTab.value = 'basic';
      return;
    }
    if (!form.platform_category_ids?.length) {
      ElMessage.warning('请选择商品分类');
      activeTab.value = 'basic';
      return;
    }
    const platformCategoryIds = Array.from(
      new Set(
        (form.platform_category_ids || [])
          .map((id) => Number(id))
          .filter((id) => Number.isFinite(id) && id > 0),
      ),
    );
    if (!platformCategoryIds.length) {
      ElMessage.warning('请选择商品分类');
      activeTab.value = 'basic';
      return;
    }
    if (!form.mer_account?.trim()) {
      ElMessage.warning('请填写店铺账号');
      activeTab.value = 'account';
      return;
    }
    if (drawerMode.value === 'create' && !form.mer_password?.trim()) {
      ElMessage.warning('请填写登录密码');
      activeTab.value = 'account';
      return;
    }
    if (drawerMode.value === 'create' && !form.mer_phone?.trim()) {
      ElMessage.warning('请填写联系电话');
      activeTab.value = 'account';
      return;
    }
    const payload: PlatformMerchantSaveInput = {
      ...form,
      mer_name: form.mer_name.trim(),
      real_name: form.real_name?.trim(),
      mer_phone: form.mer_phone?.trim(),
      mer_address: form.mer_address?.trim(),
      mer_info: form.mer_info?.trim(),
      mer_keyword: form.mer_keyword?.trim(),
      mark: form.mark?.trim(),
      mer_account: form.mer_account?.trim(),
      mer_password: form.mer_password?.trim() || undefined,
      sub_mchid: form.sub_mchid?.trim(),
      applyment_id: form.applyment_id?.trim(),
      business_id: form.business_id ? Number(form.business_id) : 0,
      store_group_ids: [...(form.store_group_ids || [])],
      goods_types: [...(form.goods_types || [])],
      platform_category_ids: platformCategoryIds,
      mer_star: form.mer_star ?? 5,
      mer_avatar: form.mer_avatar?.trim() || '',
    };
    shopDrawerApi.lock();
    try {
      if (editingId.value) {
        await updatePlatformMerchant(editingId.value, payload);
        ElMessage.success('店铺已更新');
      } else {
        await createPlatformMerchant(payload);
        ElMessage.success('店铺已新增');
      }
      shopDrawerApi.close();
      await loadCounts();
      gridApi.reload();
  } finally {
      shopDrawerApi.unlock();
    }
  },
});

function setStatusTab(status: number) {
  if (statusFilter.value === status) return;
  statusFilter.value = status;
  gridApi.reload();
}

function resetForm(row?: PlatformMerchantRow) {
  editingId.value = row?.mer_id || 0;
  headerSnapshot.value = row ? { ...row } : {};
  const isCreate = !row;
  Object.assign(form, {
    mer_name: row?.mer_name || '',
    business_id: row?.business_id || undefined,
    real_name: row?.real_name || '',
    mer_phone: row?.mer_phone || '',
    mer_address: row?.mer_address || '',
    mer_info: row?.mer_info || '',
    mer_keyword: row?.mer_keyword || '',
    mark: row?.mark || '',
    category_id: row?.category_id || undefined,
    type_id: row?.type_id || undefined,
    region_id: row?.region_id || undefined,
    is_best: row ? row.is_best === 1 : false,
    offline_pay: row ? row.offline_pay === 1 : false,
    is_trader: row ? row.is_trader === 1 : false,
    is_audit: row ? row.is_audit !== 0 : isCreate ? false : true,
    is_bro_room: row ? row.is_bro_room === 1 : false,
    is_bro_goods: row ? row.is_bro_goods === 1 : false,
    commission_switch: row ? row.commission_switch === 1 : false,
    commission_rate: row?.commission_rate ?? 0,
    mer_account: row?.mer_account || '',
    mer_password: '',
    sub_mchid: row?.sub_mchid || '',
    applyment_id: row?.applyment_id || '',
    care_count: row?.care_count ?? 0,
    care_ficti: row?.care_ficti ?? 0,
    sort: row?.sort ?? 0,
    status: row ? isEnabled(row) : true,
    store_group_ids: [...(row?.store_group_ids || [])],
    goods_types: isCreate ? [] : parseGoodsTypes(row),
    platform_category_ids: isCreate ? [] : parsePlatformCategoryIds(row),
    mer_star: row?.mer_star ?? 5,
    mer_avatar: row?.mer_avatar || '',
  });
}

function openCreate() {
  drawerMode.value = 'create';
  activeTab.value = 'basic';
  resetForm();
  void loadPlatformCategories();
  shopDrawerApi
    .setState({
      title: '新增店铺',
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
}

async function openEdit(row: PlatformMerchantRow) {
  drawerMode.value = 'edit';
  activeTab.value = 'basic';
  resetForm(row);
  void loadPlatformCategories();
  shopDrawerApi
    .setState({
      loading: true,
      title: `编辑店铺 · ${row.mer_name}`,
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
  try {
    const detail = await fetchPlatformMerchant(row.mer_id);
    resetForm(detail);
    headerSnapshot.value = detail;
  } finally {
    shopDrawerApi.setState({ loading: false });
  }
}

async function openDetail(row: PlatformMerchantRow) {
  drawerMode.value = 'view';
  activeTab.value = 'basic';
  resetOperateLogQuery();
  resetForm(row);
  shopDrawerApi
    .setState({
      loading: true,
      title: '店铺详情',
      showConfirmButton: false,
      cancelText: '关闭',
    })
    .open();
  try {
    const detail = await fetchPlatformMerchant(row.mer_id);
    resetForm(detail);
    headerSnapshot.value = detail;
  } finally {
    shopDrawerApi.setState({ loading: false });
  }
}

function switchToEditFromDetail() {
  if (!canManage.value || !editingId.value) return;
  drawerMode.value = 'edit';
  if (activeTab.value === 'logs') {
    activeTab.value = 'basic';
  }
  shopDrawerApi.setState({
    title: `编辑店铺 · ${form.mer_name || ''}`,
    showConfirmButton: true,
    confirmText: '保存',
    cancelText: '取消',
  });
}

function resetOperateLogQuery() {
  operateLogs.value = [];
  operateLogTotal.value = 0;
  operateLogPage.value = 1;
  operateLogLimit.value = 10;
  operateLogTerminal.value = '';
  operateLogDates.value = [];
}

async function loadOperateLogs() {
  if (drawerMode.value !== 'view' || !editingId.value) return;
  operateLogLoading.value = true;
  try {
    const range = Array.isArray(operateLogDates.value)
      ? operateLogDates.value
      : [];
    const data = await fetchMerchantOperateLogs(editingId.value, {
      page: operateLogPage.value,
      limit: operateLogLimit.value,
      terminal: operateLogTerminal.value || undefined,
      start_date: range[0],
      end_date: range[1],
    });
    operateLogs.value = data.list || [];
    operateLogTotal.value = data.total || 0;
  } catch {
    operateLogs.value = [];
    operateLogTotal.value = 0;
  } finally {
    operateLogLoading.value = false;
  }
}

function onOperateLogPageChange(page: number) {
  operateLogPage.value = page;
  void loadOperateLogs();
}

function onOperateLogSizeChange(size: number) {
  operateLogLimit.value = size;
  operateLogPage.value = 1;
  void loadOperateLogs();
}

function onOperateLogFilterChange() {
  operateLogPage.value = 1;
  void loadOperateLogs();
}

function openLogin(_row?: PlatformMerchantRow) {
  const origin = String(import.meta.env.VITE_MERCHANT_ADMIN_URL || '')
    .trim()
    .replace(/\/+$/, '');
  if (!origin) {
    ElMessage.error('未配置 VITE_MERCHANT_ADMIN_URL（店铺管理系统域名）');
    return;
  }
  // 新标签直接打开店铺后台登录页（hash），不带 from/token。
  window.open(`${origin}/#/auth/login`, '_blank');
}

function onDetailMoreCommand(command: string) {
  if (!editingId.value) return;
  const row = {
    mer_id: editingId.value,
    mer_name: form.mer_name || headerSnapshot.value.mer_name || '',
    status: form.status ? 1 : 0,
    mer_state: form.status ? 1 : 0,
    is_best: form.is_best ? 1 : 0,
  } as PlatformMerchantRow;
  if (command === 'login') {
    openLogin(row);
    return;
  }
  if (command === 'toggle-status' && canManage.value) {
    void (async () => {
      const enabled = !form.status;
      try {
        await confirm({
          content: `${enabled ? '开启' : '关闭'}后将立即影响该店铺的经营状态，是否继续？`,
          icon: 'warning',
          title: '提示',
        });
        await updatePlatformMerchantStatus(editingId.value, enabled);
        form.status = enabled;
        headerSnapshot.value = {
          ...headerSnapshot.value,
          status: enabled ? 1 : 0,
          mer_state: enabled ? 1 : 0,
        };
        ElMessage.success('店铺状态已更新');
        await loadCounts();
        gridApi.reload();
      } catch {
        // cancelled or failed
      }
    })();
  }
}

watch(activeTab, (tab) => {
  if (tab === 'logs' && drawerMode.value === 'view') {
    void loadOperateLogs();
  }
});

async function changeStatus(row: PlatformMerchantRow, enabled: boolean) {
  const before = isEnabled(row);
  try {
    await confirm({
      content: `${enabled ? '开启' : '关闭'}后将立即影响该店铺的经营状态，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await updatePlatformMerchantStatus(row.mer_id, enabled);
    row.status = enabled ? 1 : 0;
    row.mer_state = enabled ? 1 : 0;
    ElMessage.success('店铺状态已更新');
    await loadCounts();
    if (Number(isEnabled(row)) !== statusFilter.value) {
      gridApi.reload();
    }
  } catch {
    row.status = before ? 1 : 0;
    row.mer_state = before ? 1 : 0;
  }
}

async function changeRecommend(row: PlatformMerchantRow, enabled: boolean) {
  const before = row.is_best === 1;
  try {
    await updatePlatformMerchantRecommend(row.mer_id, enabled);
    row.is_best = enabled ? 1 : 0;
    ElMessage.success(enabled ? '已设为推荐' : '已取消推荐');
  } catch {
    row.is_best = before ? 1 : 0;
  }
}

function syncFilterSelectOptions() {
  gridApi.formApi?.updateSchema([
    {
      fieldName: 'category_id',
      componentProps: {
        clearable: true,
        options: categories.value.map((c) => ({
          label: c.category_name,
          value: c.merchant_category_id,
        })),
        placeholder: '请选择',
      },
    },
    {
      fieldName: 'type_id',
      componentProps: {
        clearable: true,
        options: types.value.map((t) => ({ label: t.name, value: t.id })),
        placeholder: '请选择',
      },
    },
    {
      fieldName: 'region_id',
      componentProps: {
        clearable: true,
        options: regions.value,
        placeholder: '请选择',
      },
    },
  ]);
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('merchant.status.manage');
  await Promise.all([loadCounts(), loadFilterOptions()]);
  syncFilterSelectOptions();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="merchant-toolbar">
          <div class="merchant-status-tabs" role="tablist">
            <button
              type="button"
              role="tab"
              class="merchant-status-tabs__item"
              :aria-selected="statusFilter === 1"
              :class="{ 'is-active': statusFilter === 1 }"
              @click="setStatusTab(1)"
            >
              正常开启的店铺({{ openCount }})
            </button>
            <button
              type="button"
              role="tab"
              class="merchant-status-tabs__item"
              :aria-selected="statusFilter === 0"
              :class="{ 'is-active': statusFilter === 0 }"
              @click="setStatusTab(0)"
            >
              已关闭店铺({{ closedCount }})
            </button>
        </div>
          <div class="merchant-toolbar__actions">
            <ElButton
              v-if="canManage"
              :icon="Plus"
              type="primary"
              @click="openCreate"
            >
              新增店铺
            </ElButton>
      </div>
        </div>
      </template>

      <template #cover="{ row }">
        <ElImage
          v-if="row.mer_avatar"
          class="shop-list-cover"
          :src="resolveCosMediaUrl(row.mer_avatar)"
          fit="cover"
          alt="店铺封面"
        >
          <template #error>
            <span>—</span>
          </template>
        </ElImage>
        <span v-else>—</span>
      </template>

      <template #mer_name="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">
          {{ row.mer_name || '—' }}
        </ElButton>
      </template>

      <template #recommend="{ row }">
        <ElSwitch
          v-if="canManage"
          :model-value="row.is_best === 1"
          @change="(enabled: string | number | boolean) => changeRecommend(row, Boolean(enabled))"
        />
        <ElTag v-else :type="row.is_best === 1 ? 'success' : 'info'">
          {{ row.is_best === 1 ? '是' : '否' }}
        </ElTag>
      </template>

      <template #status="{ row }">
        <ElSwitch
          v-if="canManage"
          :model-value="isEnabled(row)"
          @change="(enabled: string | number | boolean) => changeStatus(row, Boolean(enabled))"
        />
        <ElTag v-else :type="isEnabled(row) ? 'success' : 'info'">
          {{ statusText(row) }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openLogin(row)">登录</ElButton>
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          :type="isEnabled(row) ? 'danger' : 'primary'"
          @click="changeStatus(row, !isEnabled(row))"
        >
          {{ isEnabled(row) ? '关闭' : '开启' }}
        </ElButton>
      </template>
    </Grid>

    <ShopDrawer>
      <div class="shop-drawer">
        <div v-if="drawerMode !== 'create'" class="shop-drawer__header">
          <div class="shop-drawer__brand-row">
            <div class="shop-drawer__brand">
              <div class="shop-drawer__avatar">
                <img
                  v-if="form.mer_avatar"
                  :src="form.mer_avatar"
                  alt="店铺封面"
                />
                <Shop v-else />
              </div>
              <div class="shop-drawer__titles">
                <div class="shop-drawer__name-row">
                  <span class="shop-drawer__name">
                    {{ form.mer_name || '未命名店铺' }}
                  </span>
                  <ElTag
                    v-if="selectedType?.name || headerSnapshot.type_name"
                    size="small"
                    type="danger"
                  >
                    {{ selectedType?.name || headerSnapshot.type_name }}
                  </ElTag>
                </div>
                <div class="shop-drawer__sub">
                  {{ selectedBusinessLabel }}
                </div>
              </div>
            </div>
            <div v-if="drawerMode === 'view'" class="shop-drawer__actions">
              <ElButton
                v-if="canManage"
                type="primary"
                @click="switchToEditFromDetail"
              >
                编辑
              </ElButton>
              <ElDropdown trigger="click" @command="onDetailMoreCommand">
                <ElButton :icon="MoreFilled" />
                <template #dropdown>
                  <ElDropdownMenu>
                    <ElDropdownItem command="login">登录店铺</ElDropdownItem>
                    <ElDropdownItem
                      v-if="canManage"
                      command="toggle-status"
                    >
                      {{ form.status ? '关闭店铺' : '开启店铺' }}
                    </ElDropdownItem>
                  </ElDropdownMenu>
                </template>
              </ElDropdown>
            </div>
          </div>
          <div class="shop-drawer__meta">
            <div class="shop-drawer__meta-item">
              <span class="label">联系人</span>
              <span class="value">{{ form.real_name || '—' }}</span>
            </div>
            <div class="shop-drawer__meta-item">
              <span class="label">联系电话</span>
              <span class="value">{{ form.mer_phone || '—' }}</span>
            </div>
            <div class="shop-drawer__meta-item">
              <span class="label">状态</span>
              <span class="value">{{ form.status ? '开启' : '关闭' }}</span>
            </div>
            <div class="shop-drawer__meta-item">
              <span class="label">入驻时间</span>
              <span class="value">{{ formatTime(headerSnapshot.create_time) }}</span>
            </div>
          </div>
        </div>

        <ElTabs v-model="activeTab" class="shop-drawer__tabs">
          <ElTabPane label="基本信息" name="basic">
            <template v-if="drawerMode === 'view'">
              <div class="shop-section">
                <div class="shop-section__title">基础信息</div>
                <div class="shop-desc-grid">
                  <div class="shop-desc">
                    <span class="label">店铺名称</span>
                    <span class="value">{{ displayOrDash(form.mer_name) }}</span>
                  </div>
                  <div class="shop-desc shop-desc--full">
                    <span class="label">店铺封面</span>
                    <span class="value">
                      <img
                        v-if="form.mer_avatar"
                        class="shop-cover-thumb"
                        :src="form.mer_avatar"
                        alt="店铺封面"
                      />
                      <template v-else>—</template>
                    </span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺类型</span>
                    <span class="value">{{ form.is_trader ? '自营' : '非自营' }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺分类</span>
                    <span class="value">
                      {{
                        selectedCategory?.category_name ||
                        headerSnapshot.category_name ||
                        '—'
                      }}
                      <span v-if="categoryCommissionHint" class="hint-danger">
                        {{ categoryCommissionHint }}
                      </span>
                    </span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺分组</span>
                    <span class="value">{{ selectedStoreGroupLabels }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺区域</span>
                    <span class="value">{{ selectedRegionLabel }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">所属商户</span>
                    <span class="value">{{ selectedBusinessLabel }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">推荐店铺</span>
                    <span class="value">{{ yesNo(form.is_best) }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺类型</span>
                    <span class="value">
                      {{ selectedType?.name || headerSnapshot.type_name || '—' }}
                    </span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">商品分类</span>
                    <span class="value">{{ selectedPlatformCategoryLabels }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">商品类型</span>
                    <span class="value">{{ selectedGoodsTypeLabels }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺状态</span>
                    <span class="value">{{ form.status ? '开启' : '关闭' }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">店铺星级</span>
                    <span class="value">{{ form.mer_star || 5 }}星</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">排序</span>
                    <span class="value">{{ form.sort ?? 0 }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">更新时间</span>
                    <span class="value">{{ formatTime(headerSnapshot.create_time) }}</span>
                  </div>
                  <div class="shop-desc shop-desc--full">
                    <span class="label">备注</span>
                    <span class="value">{{ displayOrDash(form.mark) }}</span>
                  </div>
                </div>
              </div>
            </template>
            <template v-else>
            <div class="shop-section">
              <div class="shop-section__title">基础信息</div>
              <ElForm
                class="shop-form-columns"
                label-position="top"
                :disabled="isReadonly"
              >
                <div class="shop-form-col">
                  <ElFormItem label="店铺名称" required>
                    <ElInput v-model="form.mer_name" maxlength="128" />
                    <div class="field-help">
                      支持中英文、数字与符号，将作为店铺展示名称
                    </div>
                  </ElFormItem>
                  <ElFormItem label="店铺封面" required>
                    <ImageField
                      v-model="form.mer_avatar"
                      :disabled="isReadonly"
                      :preview-size="120"
                    />
                  </ElFormItem>
                  <ElFormItem label="所属商户">
                    <ElCascader
                      v-model="form.business_id"
                      :options="businessOptions"
                      :props="{
                        checkStrictly: true,
                        emitPath: false,
                        value: 'value',
                        label: 'label',
                        children: 'children',
                      }"
            clearable
                      filterable
                      class="w-full"
                      placeholder="请选择所属商户"
                    />
                    <div class="field-help">
                      选择店铺归属的商户主体，决定店铺管理权限
                    </div>
                  </ElFormItem>
                  <ElFormItem label="店铺分类" required>
                    <ElSelect v-model="form.category_id" clearable class="w-full">
                      <ElOption
                        v-for="item in categories"
                        :key="item.merchant_category_id"
                        :label="item.category_name"
                        :value="item.merchant_category_id"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="商品分类" required>
                    <ElCascader
                      v-model="form.platform_category_ids"
                      :options="platformCategoryTree"
                      :props="platformCategoryCascaderProps"
                      clearable
                      filterable
                      collapse-tags
                      collapse-tags-tooltip
                      class="w-full"
                      placeholder="请选择商品分类"
                      separator=" / "
                    />
                  </ElFormItem>
                  <ElFormItem label="商品类型" required>
                    <ElCheckboxGroup v-model="form.goods_types">
                      <ElCheckbox
                        v-for="item in GOODS_TYPE_OPTIONS"
                        :key="item.value"
                        :label="item.value"
                      >
                        {{ item.label }}
                      </ElCheckbox>
                    </ElCheckboxGroup>
                  </ElFormItem>
                  <ElFormItem label="店铺类型" required>
                    <ElSelect v-model="form.type_id" clearable class="w-full">
                      <ElOption
                        v-for="item in types"
                        :key="item.id"
                        :label="item.name"
                        :value="item.id"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="推荐店铺">
                    <ElSwitch v-model="form.is_best" />
                  </ElFormItem>
                  <ElFormItem label="备注">
                    <ElInput
                      v-model="form.mark"
                      :rows="2"
                      maxlength="500"
                      type="textarea"
                    />
                  </ElFormItem>
                </div>
                <div class="shop-form-col">
                  <ElFormItem label="店铺地址">
                    <ElInput v-model="form.mer_address" maxlength="255" />
                  </ElFormItem>
                  <ElFormItem label="店铺分组">
                    <ElSelect
                      v-model="form.store_group_ids"
                      clearable
                      collapse-tags
                      collapse-tags-tooltip
                      multiple
                      class="w-full"
                    >
                      <ElOption
                        v-for="item in storeGroupOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="店铺区域">
                    <ElSelect v-model="form.region_id" clearable class="w-full">
                      <ElOption
                        v-for="item in regions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="店铺星级">
                    <ElSelect v-model="form.mer_star" class="w-full">
                      <ElOption
                        v-for="item in STAR_OPTIONS"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="自营类型">
                    <ElRadioGroup v-model="form.is_trader">
                      <ElRadio :label="true">自营</ElRadio>
                      <ElRadio :label="false">非自营</ElRadio>
                    </ElRadioGroup>
                  </ElFormItem>
                  <ElFormItem label="排序">
                    <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
                    <div class="field-help">数值越小越靠前。</div>
                  </ElFormItem>
                  <ElFormItem label="店铺状态">
                    <ElSwitch v-model="form.status" />
                  </ElFormItem>
                </div>
              </ElForm>
            </div>
            </template>
          </ElTabPane>

          <ElTabPane label="经营信息" name="operate">
            <template v-if="drawerMode === 'view'">
              <div class="shop-section">
                <div class="shop-section__title">费用信息</div>
                <div class="shop-desc-grid">
                  <div class="shop-desc">
                    <span class="label">手续费单独设置</span>
                    <span class="value">{{ form.commission_switch ? '开启' : '关闭' }}</span>
        </div>
                  <div class="shop-desc">
                    <span class="label">手续费</span>
                    <span class="value">
                      {{ Number(form.commission_rate || 0).toFixed(2) }}%
                      <span class="hint-danger">
                        （注：此处如未设置手续费，系统会自动读取店铺分类下对应手续费；此处已设置，则优先以此处设置为准）
                      </span>
                    </span>
      </div>
                  <div class="shop-desc">
                    <span class="label">店铺保证金</span>
                    <span class="value">{{ marginAmountText }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">保证金支付状态</span>
                    <span class="value">
                      {{ depositStatusText }}
                      <span v-if="depositNeedsTopUp" class="hint-danger">
                        （需补缴）
                      </span>
                    </span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">保证金余额</span>
                    <span class="value">
                      {{ Number(headerSnapshot.deposit_available ?? 0).toFixed(2) }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="shop-section">
                <div class="shop-section__title">经营数据</div>
                <div class="shop-desc-grid">
                  <div class="shop-desc">
                    <span class="label">实际关注人数</span>
                    <span class="value">{{ form.care_count ?? 0 }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">已关注人数</span>
                    <span class="value">{{ form.care_ficti ?? 0 }}</span>
                  </div>
                </div>
              </div>
              <div class="shop-section">
                <div class="shop-section__title">审核信息</div>
                <div class="shop-desc-grid">
                  <div class="shop-desc">
                    <span class="label">商品审核</span>
                    <span class="value">{{ auditText(!!form.is_audit) }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">直播间审核</span>
                    <span class="value">{{ auditText(!!form.is_bro_room) }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">直播商品</span>
                    <span class="value">{{ auditText(!!form.is_bro_goods) }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">线下支付</span>
                    <span class="value">{{ form.offline_pay ? '开启' : '关闭' }}</span>
                  </div>
                </div>
              </div>
              <div class="shop-section">
                <div class="shop-section__title">其他信息</div>
                <div class="shop-desc-grid">
                  <div class="shop-desc">
                    <span class="label">搜索店铺关键字</span>
                    <span class="value">{{ displayOrDash(form.mer_keyword) }}</span>
                  </div>
                  <div class="shop-desc shop-desc--full">
                    <span class="label">店铺简介</span>
                    <span class="value">{{ displayOrDash(form.mer_info) }}</span>
                  </div>
                </div>
              </div>
          </template>
            <template v-else>
            <div class="shop-section">
              <div class="shop-section__title">费用信息</div>
              <div
                v-if="drawerMode !== 'create'"
                class="shop-kv-grid"
              >
                <div class="shop-kv">
                  <span class="label">店铺保证金</span>
                  <span class="value">
                    {{
                      selectedType?.is_margin || headerSnapshot.type_is_margin
                        ? `${selectedType?.margin ?? headerSnapshot.type_margin ?? 0} 元`
                        : '无'
                    }}
                  </span>
                </div>
                <div class="shop-kv">
                  <span class="label">支付状态</span>
                  <span class="value">{{ depositText(headerSnapshot.deposit_state) }}</span>
                </div>
                <div class="shop-kv">
                  <span class="label">保证金余额</span>
                  <span class="value">
                    {{ headerSnapshot.deposit_available ?? 0 }} 元
                  </span>
                </div>
              </div>
              <div class="shop-kv shop-kv--form">
                <span class="label">手续费设置</span>
                <div class="inline-fee" :class="{ 'is-disabled': isReadonly }">
                  <ElSwitch
                    v-model="form.commission_switch"
                    :disabled="isReadonly"
                  />
                  <ElInputNumber
                    v-model="form.commission_rate"
                    :disabled="isReadonly || !form.commission_switch"
                    :min="0"
                    :max="100"
                    :precision="2"
                  />
                  <span>%</span>
                </div>
                <div class="field-help">
                  开启后按店铺独立手续费；关闭则沿用分类默认比例。
                </div>
              </div>
            </div>

            <div class="shop-section">
              <div class="shop-section__title">审核信息</div>
              <div class="shop-audit-grid">
                <div class="shop-audit">
                  <div class="shop-audit__row">
                    <span>商品审核</span>
                    <ElSwitch v-model="form.is_audit" :disabled="isReadonly" />
                  </div>
                  <div class="field-help">开启后新增商品需平台审核。</div>
                </div>
                <div class="shop-audit">
                  <div class="shop-audit__row">
                    <span>直播间审核</span>
                    <ElSwitch v-model="form.is_bro_room" :disabled="isReadonly" />
                  </div>
                  <div class="field-help">开启后创建直播间需平台审核。</div>
                </div>
                <div class="shop-audit">
                  <div class="shop-audit__row">
                    <span>直播商品审核</span>
                    <ElSwitch v-model="form.is_bro_goods" :disabled="isReadonly" />
                  </div>
                  <div class="field-help">开启后直播商品需平台审核。</div>
                </div>
                <div class="shop-audit">
                  <div class="shop-audit__row">
                    <span>线下支付</span>
                    <ElSwitch v-model="form.offline_pay" :disabled="isReadonly" />
                  </div>
                  <div class="field-help">开启后支持线下收款核销。</div>
                </div>
              </div>
            </div>

            <div v-if="drawerMode !== 'create'" class="shop-section">
              <div class="shop-section__title">经营数据</div>
              <ElForm class="shop-form-grid" label-position="top" :disabled="isReadonly">
                <ElFormItem label="实际关注人数">
                  <ElInputNumber v-model="form.care_count" :min="0" class="w-full" />
                </ElFormItem>
                <ElFormItem label="已关注人数（含虚拟）">
                  <ElInputNumber v-model="form.care_ficti" :min="0" class="w-full" />
                </ElFormItem>
              </ElForm>
            </div>

            <div class="shop-section">
              <div class="shop-section__title">其他信息</div>
              <ElForm label-position="top" :disabled="isReadonly">
                <ElFormItem label="店铺关键字">
                  <ElInput v-model="form.mer_keyword" maxlength="255" />
                  <div class="field-help">
                    用于店铺搜索与 SEO，多个关键字可用空格分隔。
                  </div>
                </ElFormItem>
                <ElFormItem v-if="drawerMode !== 'create'" label="店铺简介">
                  <ElInput
                    v-model="form.mer_info"
                    :rows="3"
                    maxlength="1000"
                    type="textarea"
                  />
                </ElFormItem>
              </ElForm>
            </div>
          </template>
          </ElTabPane>

          <ElTabPane label="账号信息" name="account">
            <template v-if="drawerMode === 'view'">
              <div class="shop-section">
                <div class="shop-section__title">登录账号</div>
                <div class="shop-desc-grid">
                  <div class="shop-desc">
                    <span class="label">店铺账号</span>
                    <span class="value">{{ displayOrDash(form.mer_account) }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">登录密码</span>
                    <span class="value">**********</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">联系人</span>
                    <span class="value">{{ displayOrDash(form.real_name) }}</span>
                  </div>
                  <div class="shop-desc">
                    <span class="label">联系电话</span>
                    <span class="value">{{ displayOrDash(form.mer_phone) }}</span>
                  </div>
                </div>
              </div>
          </template>
            <template v-else>
            <div class="shop-section">
              <div class="shop-section__title">登录账号</div>
              <ElForm class="shop-form-grid" label-position="top" :disabled="isReadonly">
                <ElFormItem label="店铺账号" required>
                  <ElInput
                    v-model="form.mer_account"
                    maxlength="64"
                    :disabled="isReadonly || drawerMode === 'edit'"
                  />
                  <div class="field-help">用于店铺管理员登录后台，创建后不可随意变更。</div>
                </ElFormItem>
                <ElFormItem
                  v-if="drawerMode === 'create'"
                  label="登录密码"
                  required
                >
                  <ElInput
                    v-model="form.mer_password"
                    maxlength="32"
                    show-password
                    type="password"
                  />
                  <div class="field-help">建议 8-16 位，含字母与数字。</div>
                </ElFormItem>
                <ElFormItem label="联系人">
                  <ElInput v-model="form.real_name" maxlength="64" />
                </ElFormItem>
                <ElFormItem
                  label="联系电话"
                  :required="drawerMode === 'create'"
                >
                  <ElInput v-model="form.mer_phone" maxlength="32" />
                </ElFormItem>
              </ElForm>
      </div>

            <div class="shop-section">
              <div class="shop-section__title">财务账号</div>
              <ElForm class="shop-form-grid" label-position="top" :disabled="isReadonly">
                <ElFormItem label="平台收付通">
                  <ElInput v-model="form.sub_mchid" maxlength="64" />
                  <div class="field-help">微信收付通 / 分账子商户号。</div>
                </ElFormItem>
                <ElFormItem label="服务商特约店铺">
                  <ElInput v-model="form.applyment_id" maxlength="64" />
                  <div class="field-help">服务商特约商户申请单号或标识。</div>
                </ElFormItem>
              </ElForm>
      </div>
            </template>
          </ElTabPane>

          <ElTabPane
            v-if="drawerMode === 'view'"
            label="操作记录"
            name="logs"
          >
            <div class="shop-log-filters">
              <ElSelect
                v-model="operateLogTerminal"
                clearable
                placeholder="操作端"
                style="width: 160px"
                @change="onOperateLogFilterChange"
              >
                <ElOption label="平台操作" value="platform" />
                <ElOption label="商户操作" value="merchant" />
              </ElSelect>
              <ElDatePicker
                v-model="operateLogDates"
                type="daterange"
                value-format="YYYY-MM-DD"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                @change="onOperateLogFilterChange"
              />
            </div>
            <ElTable
              v-loading="operateLogLoading"
              :data="operateLogs"
              border
              class="shop-log-table"
            >
              <ElTableColumn
                type="index"
                label="序号"
                width="70"
                :index="(index: number) => (operateLogPage - 1) * operateLogLimit + index + 1"
              />
              <ElTableColumn prop="action_label" label="操作记录" min-width="160" />
              <ElTableColumn prop="terminal" label="操作端" width="120" />
              <ElTableColumn prop="role_label" label="操作角色" width="140" />
              <ElTableColumn label="操作人" min-width="160">
                <template #default="{ row }">
                  {{ row.operator_name }}/ID:{{ row.operator_id }}
                </template>
              </ElTableColumn>
              <ElTableColumn label="操作时间" width="180">
                <template #default="{ row }">
                  {{ formatTime(row.created_at) }}
                </template>
              </ElTableColumn>
            </ElTable>
            <div class="shop-log-pagination">
              <ElPagination
                background
                layout="total, prev, pager, next, jumper"
                :current-page="operateLogPage"
                :page-size="operateLogLimit"
                :total="operateLogTotal"
                @current-change="onOperateLogPageChange"
                @size-change="onOperateLogSizeChange"
              />
            </div>
          </ElTabPane>
        </ElTabs>
      </div>
    </ShopDrawer>
  </Page>
</template>

<style scoped>
.merchant-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.merchant-status-tabs {
  display: flex;
  gap: 28px;
  border-bottom: 1px solid hsl(var(--border));
}

.merchant-status-tabs__item {
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

.merchant-status-tabs__item:hover {
  color: hsl(var(--primary));
}

.merchant-status-tabs__item.is-active {
  border-bottom-color: hsl(var(--primary));
  color: hsl(var(--primary));
  font-weight: 600;
}

.merchant-toolbar__actions {
  display: flex;
  justify-content: flex-start;
}

.shop-drawer {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 100%;
}

.shop-drawer__header {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid hsl(var(--border));
}

.shop-drawer__brand-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.shop-drawer__brand {
  display: flex;
  gap: 12px;
  align-items: center;
  min-width: 0;
}

.shop-drawer__actions {
  display: flex;
  flex-shrink: 0;
  gap: 8px;
  align-items: center;
}

.shop-drawer__avatar {
  display: flex;
  width: 48px;
  height: 48px;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 10px;
  background: hsl(var(--primary) / 12%);
  color: hsl(var(--primary));
  font-size: 22px;
}

.shop-drawer__avatar img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-cover-thumb {
  display: block;
  width: 96px;
  height: 96px;
  border-radius: 8px;
  object-fit: cover;
  border: 1px solid hsl(var(--border));
}

.shop-list-cover {
  display: block;
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid hsl(var(--border));
}

.shop-drawer__name-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.shop-drawer__name {
  font-size: 18px;
  font-weight: 600;
  line-height: 28px;
}

.shop-drawer__sub {
  margin-top: 2px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.shop-drawer__meta {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.shop-drawer__meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.shop-drawer__meta-item .label,
.shop-kv .label,
.shop-desc .label,
.field-help {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 18px;
}

.shop-drawer__meta-item .value,
.shop-kv .value,
.shop-desc .value {
  font-size: 14px;
  line-height: 22px;
  word-break: break-word;
}

.shop-desc-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px 32px;
}

.shop-desc {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.shop-desc--full {
  grid-column: 1 / -1;
}

.hint-danger {
  margin-left: 6px;
  color: hsl(var(--destructive));
  font-size: 12px;
}

.shop-log-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.shop-log-table {
  width: 100%;
}

.shop-log-table :deep(thead th) {
  background: hsl(var(--primary) / 8%) !important;
}

.shop-log-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.shop-section {
  margin-bottom: 20px;
}

.shop-section__title {
  position: relative;
  margin-bottom: 14px;
  padding-left: 10px;
  font-size: 15px;
  font-weight: 600;
  line-height: 24px;
}

.shop-section__title::before {
  position: absolute;
  top: 4px;
  left: 0;
  width: 3px;
  height: 16px;
  border-radius: 2px;
  background: hsl(var(--primary));
  content: '';
}

.shop-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 20px;
}

.shop-form-columns {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 20px;
}

.shop-form-col {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.shop-form-columns :deep(.el-checkbox-group) {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}

.shop-form-grid :deep(.span-2),
.shop-form-grid .span-2 {
  grid-column: 1 / -1;
}

.shop-kv-grid,
.shop-audit-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px 24px;
}

.shop-kv {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.inline-fee {
  display: flex;
  gap: 10px;
  align-items: center;
}

.shop-audit {
  padding: 12px 14px;
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
}

.shop-audit__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 14px;
}

@media (max-width: 900px) {
  .shop-drawer__meta,
  .shop-form-grid,
  .shop-form-columns,
  .shop-kv-grid,
  .shop-audit-grid,
  .shop-desc-grid {
    grid-template-columns: 1fr;
  }
}
</style>
