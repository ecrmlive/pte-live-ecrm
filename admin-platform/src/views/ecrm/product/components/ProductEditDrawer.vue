<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import { computed, reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { VbenTiptap } from '@vben/plugins/tiptap';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElButton,
  ElCascader,
  ElCheckbox,
  ElCheckboxGroup,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElPagination,
  ElRadioButton,
  ElRadioGroup,
  ElRate,
  ElSelect,
  ElSkeleton,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';

import { uploadAttachmentApi } from '#/api/core/attachment';
import {
  createPlatformProductAdminApi,
  getPlatformProductEditApi,
  getPlatformProductStoreOptionsApi,
  listPlatformBrandsApi,
  listPlatformCategoriesApi,
  listPlatformProductOperateLogsApi,
  listPlatformProductStoresApi,
  updatePlatformProductAdminApi,
  type PlatformCategory,
  type PlatformProductAdminSaveBody,
  type PlatformProductEditDetail,
  type PlatformProductEditSKU,
  type PlatformProductOperateLog,
  type PlatformProductStoreOption,
} from '#/api/core/platform-catalog';
import { fetchProductLabels, type ProductLabelRow } from '#/api/core/ecrm';
import ImageField from '#/components/shop/image-field.vue';
import ImagesField from '#/components/shop/images-field.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';

import ProductLabelSelectModal from './ProductLabelSelectModal.vue';

type SkuEditRow = {
  key: string;
  sku_id?: number;
  specLabel: string;
  image: string;
  price: number;
  ot_price: number;
  stock: number;
  code: string;
  bar_code: string;
  weight: number;
  volume: number;
  extension_one: number;
};

const emit = defineEmits<{
  saved: [];
}>();

const loading = ref(false);
const saving = ref(false);
const activeTab = ref('basic');
const mode = ref<'create' | 'edit'>('edit');
const detail = ref<PlatformProductEditDetail>();
const labelOptions = ref<ProductLabelRow[]>([]);
const categoryTree = ref<PlatformCategory[]>([]);
const brandOptions = ref<{ label: string; value: string }[]>([]);
const storeOptions = ref<PlatformProductStoreOption[]>([]);
const merCateOptions = ref<Array<{ id: number; name: string }>>([]);
const merLabelOptions = ref<Array<{ id: number; name: string }>>([]);
const skuRows = ref<SkuEditRow[]>([]);
const specType = ref<0 | 1>(0);
const labelModalRef = ref<InstanceType<typeof ProductLabelSelectModal>>();

const form = reactive({
  store_id: undefined as number | undefined,
  title: '',
  store_info: '',
  keyword: '',
  unit_name: '件',
  brand_name: '',
  cate_id: undefined as number | undefined,
  mer_cate_id: undefined as number | undefined,
  mer_label_ids: [] as string[],
  delivery_way: [2] as number[],
  image: '',
  slider_image: [] as string[],
  ot_price: 0,
  star: 0,
  rank: 0,
  is_hot: false,
  is_benefit: false,
  is_best: false,
  is_new: false,
  cate_hot: false,
  content: '',
  refund_switch: true,
  once_min_count: 1,
  sys_labels: [] as string[],
});

const categoryCascaderProps = {
  value: 'store_category_id',
  label: 'cate_name',
  children: 'children',
  emitPath: false,
  checkStrictly: true,
};

const operateLogs = ref<PlatformProductOperateLog[]>([]);
const operateLogTotal = ref(0);
const operateLogLoading = ref(false);
const operateLogPage = ref(1);
const operateLogLimit = ref(10);
const operateLogTerminal = ref('');
const operateLogDates = ref<string[]>([]);

const imageUpload: ImageUploadOptions = {
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const row = await uploadAttachmentApi(file);
    return row.attachment_src;
  },
  onUploadError: () => {
    ElMessage.error('图片上传失败');
  },
};

const headerTitle = computed(
  () => form.title.trim() || detail.value?.title || '未命名商品',
);

const statusLabel = computed(() => {
  const map: Record<number, string> = {
    [-2]: '下架',
    [-1]: '审核未通过',
    0: '待审核',
    1: '上架显示',
  };
  return map[Number(detail.value?.status ?? 0)] || '未知';
});

const productTypeLabel = computed(() => {
  const map: Record<number, string> = {
    0: '普通商品',
    1: '虚拟',
    2: '云盘',
    3: '卡密',
    4: '预约',
    5: '年/次卡',
  };
  return map[Number(detail.value?.product_type ?? 0)] || '普通商品';
});

const formLinkedText = computed(() => {
  const id = Number(detail.value?.mer_form_id || 0);
  return id > 0 ? `#${id}` : '关闭';
});

const recommendValues = computed(() => [
  ...(form.is_hot ? ['hot'] : []),
  ...(form.is_benefit ? ['benefit'] : []),
  ...(form.is_best ? ['best'] : []),
  ...(form.is_new ? ['new'] : []),
]);

const selectedLabelNames = computed(() => {
  const map = new Map(labelOptions.value.map((o) => [String(o.id), o.name]));
  return form.sys_labels
    .map((id) => map.get(String(id)) || String(id))
    .filter(Boolean);
});

const drawerTitle = computed(() =>
  mode.value === 'create' ? '新增商品' : '编辑商品',
);

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  title: '编辑商品',
  onConfirm: async () => {
    const ok = await submit();
    // 校验失败必须阻止 Drawer 关闭（resolve 会被当成成功）
    if (!ok) return false;
  },
});

function newSkuRow(partial?: Partial<SkuEditRow>): SkuEditRow {
  return {
    key: `sku-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    sku_id: partial?.sku_id,
    specLabel: partial?.specLabel ?? '标准',
    image: String(partial?.image ?? ''),
    price: Number(partial?.price ?? 0),
    ot_price: Number(partial?.ot_price ?? 0),
    stock: Number(partial?.stock ?? 0),
    code: String(partial?.code ?? ''),
    bar_code: String(partial?.bar_code ?? ''),
    weight: Number(partial?.weight ?? 0),
    volume: Number(partial?.volume ?? 0),
    extension_one: Number(partial?.extension_one ?? 0),
  };
}

function skusFromDetail(skus: PlatformProductEditSKU[] | undefined) {
  const list = skus?.length ? skus : [];
  if (!list.length) {
    skuRows.value = [newSkuRow({ ot_price: Number(form.ot_price || 0) })];
    specType.value = 0;
    return;
  }
  specType.value = list.length > 1 ? 1 : 0;
  skuRows.value = list.map((s) => {
    const values = Object.values(s.spec || {}).filter(Boolean);
    return newSkuRow({
      sku_id: s.sku_id,
      specLabel: values.join(' / ') || s.spec_text || '标准',
      image: s.image || '',
      price: Number(s.price || 0),
      ot_price: Number(s.ot_price || 0),
      stock: Number(s.stock || 0),
      code: s.code || '',
      bar_code: s.bar_code || '',
      weight: Number(s.weight || 0),
      volume: Number(s.volume || 0),
      extension_one: Number(s.extension_one || 0),
    });
  });
}

function onSpecTypeChange(v: 0 | 1) {
  specType.value = v;
  if (v === 0 && skuRows.value.length > 1) {
    skuRows.value = [skuRows.value[0] || newSkuRow()];
  }
  if (v === 1 && skuRows.value.length < 2) {
    skuRows.value = [...skuRows.value, newSkuRow({ specLabel: '规格2' })];
  }
}

function addSkuRow() {
  const seed = skuRows.value[0];
  skuRows.value.push(
    newSkuRow({
      specLabel: `规格${skuRows.value.length + 1}`,
      ot_price: Number(seed?.ot_price || 0),
      image: String(seed?.image || ''),
    }),
  );
  specType.value = 1;
}

function removeSkuRow(index: number) {
  if (skuRows.value.length <= 1) {
    ElMessage.warning('至少保留一个规格');
    return;
  }
  skuRows.value.splice(index, 1);
  if (skuRows.value.length === 1) specType.value = 0;
}

async function loadStoreOptions(storeId: number) {
  const meta = await getPlatformProductStoreOptionsApi(storeId);
  merCateOptions.value = meta.mer_cate_options || [];
  merLabelOptions.value = meta.mer_label_options || [];
}

async function onStoreChange(storeId: number | undefined) {
  form.mer_cate_id = undefined;
  form.mer_label_ids = [];
  merCateOptions.value = [];
  merLabelOptions.value = [];
  if (!storeId) return;
  try {
    await loadStoreOptions(storeId);
  } catch {
    ElMessage.error('加载店铺分类/标签失败');
  }
}

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '—';
  return String(v);
}

function categoryIdInTree(nodes: PlatformCategory[], id: number): boolean {
  for (const node of nodes) {
    if (Number(node.store_category_id) === id) return true;
    if (node.children?.length && categoryIdInTree(node.children, id)) return true;
  }
  return false;
}

/** 去掉空 children，避免 Cascader 把叶子当成可展开节点 */
function sanitizeCategoryTree(nodes: PlatformCategory[]): PlatformCategory[] {
  return nodes.map((node) => {
    const children = node.children?.length
      ? sanitizeCategoryTree(node.children)
      : undefined;
    return children?.length ? { ...node, children } : { ...node, children: undefined };
  });
}

function fillForm(data: PlatformProductEditDetail) {
  form.title = data.title || '';
  form.store_info = data.store_info || '';
  form.keyword = data.keyword || '';
  form.unit_name = data.unit_name || '';
  form.brand_name = data.brand_name || '';
  const cateId = Number(data.cate_id || 0) || undefined;
  form.cate_id =
    cateId && categoryIdInTree(categoryTree.value, cateId) ? cateId : undefined;
  form.store_id = Number(data.store_id || 0) || undefined;
  form.mer_cate_id = Number(data.mer_cate_id || 0) || undefined;
  form.mer_label_ids = [...(data.mer_label_ids || [])].map(String);
  form.delivery_way = (data.delivery_way?.length ? [...data.delivery_way] : [2]).map(
    Number,
  );
  form.image = data.image || '';
  form.slider_image = [...(data.slider_image || [])];
  form.ot_price = Number(data.ot_price || 0);
  skusFromDetail(data.skus);
  form.star = Number(data.star || 0);
  form.rank = Number(data.rank || 0);
  form.is_hot = Number(data.is_hot) === 1;
  form.is_benefit = Number(data.is_benefit) === 1;
  form.is_best = Number(data.is_best) === 1;
  form.is_new = Number(data.is_new) === 1;
  form.cate_hot = Number(data.cate_hot) === 1;
  form.content = data.content || '';
  form.refund_switch = Number(data.refund_switch ?? 1) === 1;
  form.once_min_count = Number(data.once_min_count || 1);
  form.sys_labels = [...(data.sys_labels || [])].map(String);
  if (!form.unit_name.trim()) form.unit_name = '件';
}

async function loadCommonMeta() {
  const [labels, categories, brands] = await Promise.all([
    fetchProductLabels().catch(() => ({ list: [] as ProductLabelRow[] })),
    listPlatformCategoriesApi().catch(() => ({ list: [] as PlatformCategory[] })),
    listPlatformBrandsApi().catch(() => ({ list: [] })),
  ]);
  labelOptions.value = (labels.list || []).filter((x) => Number(x.status) !== 0);
  categoryTree.value = sanitizeCategoryTree(categories.list || []);
  brandOptions.value = (brands.list || []).map((b) => ({
    label: b.brand_name,
    value: b.brand_name,
  }));
}

async function openCreate() {
  mode.value = 'create';
  activeTab.value = 'basic';
  detail.value = undefined;
  drawerApi.setState({ title: '新增商品', loading: true }).open();
  loading.value = true;
  try {
    const [stores] = await Promise.all([
      listPlatformProductStoresApi().catch(() => ({
        list: [] as PlatformProductStoreOption[],
      })),
      loadCommonMeta(),
    ]);
    storeOptions.value = stores.list || [];
    Object.assign(form, {
      store_id: storeOptions.value[0]?.store_id,
      title: '',
      store_info: '',
      keyword: '',
      unit_name: '件',
      brand_name: '',
      cate_id: undefined,
      mer_cate_id: undefined,
      mer_label_ids: [],
      delivery_way: [2],
      image: '',
      slider_image: [],
      ot_price: 0,
      star: 0,
      rank: 0,
      is_hot: false,
      is_benefit: false,
      is_best: false,
      is_new: false,
      cate_hot: false,
      content: '',
      refund_switch: true,
      once_min_count: 1,
      sys_labels: [],
    });
    skuRows.value = [newSkuRow()];
    specType.value = 0;
    merCateOptions.value = [];
    merLabelOptions.value = [];
    if (form.store_id) await loadStoreOptions(form.store_id);
  } catch {
    ElMessage.error('打开新增商品失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

async function open(productId: number) {
  mode.value = 'edit';
  activeTab.value = 'basic';
  detail.value = undefined;
  drawerApi.setState({ title: '编辑商品', loading: true }).open();
  loading.value = true;
  try {
    const [data] = await Promise.all([
      getPlatformProductEditApi(productId),
      loadCommonMeta(),
    ]);
    detail.value = data;
    merCateOptions.value = data.mer_cate_options || [];
    merLabelOptions.value = data.mer_label_options || [];
    fillForm(data);
    resetOperateLogQuery();
    void loadOperateLogs();
  } catch {
    ElMessage.error('加载商品编辑数据失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

function failBasic(message: string) {
  ElMessage.warning(message);
  activeTab.value = 'basic';
  return false;
}

function validateBasic(): boolean {
  if (mode.value === 'create' && !Number(form.store_id || 0)) {
    return failBasic('请选择所属店铺');
  }
  if (!form.title.trim()) return failBasic('请填写商品名称');
  if (!form.image.trim()) return failBasic('请选择封面图');
  if (!form.slider_image.length) return failBasic('请至少选择一张轮播图');
  if (!form.store_info.trim()) return failBasic('请填写商品简介');
  if (!Number(form.cate_id || 0)) return failBasic('请选择平台分类');
  if (!Number(form.mer_cate_id || 0)) return failBasic('请选择店铺分类');
  if (!form.mer_label_ids.length) return failBasic('请选择商品标签');
  if (!form.sys_labels.length) return failBasic('请选择平台标签');
  if (!form.brand_name.trim()) return failBasic('请选择或填写品牌');
  if (!form.unit_name.trim()) return failBasic('请填写单位');
  if (!form.keyword.trim()) return failBasic('请填写关键字');
  if (!form.delivery_way.length) return failBasic('请至少选择一种配送方式');
  if (!merCateOptions.value.length) {
    return failBasic('该店铺暂无店铺分类，请先在店铺后台维护');
  }
  if (!skuRows.value.length) {
    activeTab.value = 'sku';
    ElMessage.warning('请至少配置一个规格');
    return false;
  }
  for (const row of skuRows.value) {
    if (
      Number(row.price) < 0 ||
      Number(row.ot_price) < 0 ||
      Number(row.stock) < 0 ||
      Number(row.weight) < 0 ||
      Number(row.volume) < 0 ||
      Number(row.extension_one) < 0
    ) {
      activeTab.value = 'sku';
      ElMessage.warning('规格售价/划线价/库存/重量/体积/返佣不能为负数');
      return false;
    }
    if (!String(row.specLabel || '').trim()) {
      activeTab.value = 'sku';
      ElMessage.warning('请填写规格名称');
      return false;
    }
  }
  return true;
}

function buildSaveBody(): PlatformProductAdminSaveBody {
  const firstOt = Number(skuRows.value[0]?.ot_price || 0);
  return {
    title: form.title.trim(),
    store_info: form.store_info.trim(),
    keyword: form.keyword.trim(),
    unit_name: form.unit_name.trim(),
    brand_name: form.brand_name.trim(),
    cate_id: Number(form.cate_id),
    mer_cate_id: Number(form.mer_cate_id),
    mer_label_ids: [...form.mer_label_ids],
    delivery_way: form.delivery_way.map(Number),
    image: form.image.trim(),
    slider_image: [...form.slider_image],
    ot_price: firstOt,
    skus: skuRows.value.map((row) => ({
      sku_id: row.sku_id,
      spec: { 规格: String(row.specLabel || '').trim() || '标准' },
      image: String(row.image || '').trim(),
      price: Number(row.price || 0),
      ot_price: Number(row.ot_price || 0),
      stock: Number(row.stock || 0),
      code: String(row.code || '').trim(),
      bar_code: String(row.bar_code || '').trim(),
      weight: Number(row.weight || 0),
      volume: Number(row.volume || 0),
      extension_one: Number(row.extension_one || 0),
      status: 1,
    })),
    star: form.star,
    rank: form.rank,
    is_hot: form.is_hot ? 1 : 0,
    is_benefit: form.is_benefit ? 1 : 0,
    is_best: form.is_best ? 1 : 0,
    is_new: form.is_new ? 1 : 0,
    cate_hot: form.cate_hot ? 1 : 0,
    content: form.content,
    refund_switch: form.refund_switch ? 1 : 0,
    once_min_count: form.once_min_count,
    sys_labels: form.sys_labels,
  };
}

/** @returns 是否保存成功（false=校验失败或接口失败，Drawer 保持打开） */
async function submit(): Promise<boolean> {
  if (mode.value === 'edit' && !detail.value) return false;
  if (!validateBasic()) return false;

  saving.value = true;
  drawerApi.setState({ confirming: true });
  try {
    const body = buildSaveBody();
    if (mode.value === 'create') {
      await createPlatformProductAdminApi({
        ...body,
        store_id: Number(form.store_id),
      });
      ElMessage.success('新增成功');
    } else {
      await updatePlatformProductAdminApi(detail.value!.product_id, body);
      ElMessage.success('保存成功');
    }
    emit('saved');
    return true;
  } catch {
    return false;
  } finally {
    saving.value = false;
    drawerApi.setState({ confirming: false });
  }
}

function onRecommendChange(vals: Array<string | number | boolean>) {
  const set = new Set(vals.map(String));
  form.is_hot = set.has('hot');
  form.is_benefit = set.has('benefit');
  form.is_best = set.has('best');
  form.is_new = set.has('new');
}

function openLabelSelect() {
  labelModalRef.value?.open({
    productId: detail.value?.product_id,
    selectedIds: form.sys_labels,
    options: labelOptions.value,
  });
}

function onLabelSubmit(ids: string[]) {
  form.sys_labels = [...ids];
}

function resetOperateLogQuery() {
  operateLogPage.value = 1;
  operateLogLimit.value = 10;
  operateLogTerminal.value = '';
  operateLogDates.value = [];
}

async function loadOperateLogs() {
  if (!detail.value) return;
  operateLogLoading.value = true;
  try {
    const data = await listPlatformProductOperateLogsApi(detail.value.product_id, {
      page: operateLogPage.value,
      limit: operateLogLimit.value,
      terminal: operateLogTerminal.value || undefined,
      date_from: operateLogDates.value?.[0],
      date_to: operateLogDates.value?.[1],
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

function onOperateLogFilterChange() {
  operateLogPage.value = 1;
  void loadOperateLogs();
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

defineExpose({ open, openCreate });
</script>

<template>
  <Drawer :title="drawerTitle">
    <ElSkeleton :loading="loading" animated :rows="10">
      <template #default>
        <div class="product-edit">
          <div v-if="mode === 'edit' && detail" class="product-edit__header">
            <div class="product-edit__identity">
              <div class="product-edit__icon">
                <IconifyIcon icon="ant-design:shopping-outlined" />
              </div>
              <div class="product-edit__titles">
                <div class="product-edit__name">{{ headerTitle }}</div>
                <div class="product-edit__sub">
                  商品ID：{{ detail.product_id }}
                </div>
              </div>
            </div>
            <div class="product-edit__meta">
              <div class="product-edit__meta-item">
                <span class="label">类型</span>
                <span class="value">{{ productTypeLabel }}</span>
              </div>
              <div class="product-edit__meta-item">
                <span class="label">状态</span>
                <span class="value">{{ statusLabel }}</span>
              </div>
              <div class="product-edit__meta-item">
                <span class="label">销量</span>
                <span class="value">{{ detail.sales }}</span>
              </div>
              <div class="product-edit__meta-item">
                <span class="label">库存</span>
                <span class="value">{{ detail.stock }}</span>
              </div>
              <div class="product-edit__meta-item">
                <span class="label">创建时间</span>
                <span class="value">{{ formatShanghaiDateTime(detail.create_time) }}</span>
              </div>
            </div>
          </div>

          <ElTabs v-model="activeTab" class="product-edit__tabs">
            <ElTabPane label="基本信息" name="basic">
              <div class="product-edit__section-title">基本信息</div>
              <ElForm
                label-position="left"
                label-width="96px"
                class="product-edit__form product-edit__form--edit"
              >
                <ElFormItem v-if="mode === 'create'" label="所属店铺" required>
                  <ElSelect
                    v-model="form.store_id"
                    filterable
                    placeholder="请选择店铺"
                    class="product-edit__control w-full"
                    @change="onStoreChange"
                  >
                    <ElOption
                      v-for="s in storeOptions"
                      :key="s.store_id"
                      :label="`${s.store_name}（${s.merchant_name}）`"
                      :value="s.store_id"
                    />
                  </ElSelect>
                </ElFormItem>
                <ElFormItem label="商品名称" required>
                  <ElInput
                    v-model="form.title"
                    maxlength="128"
                    show-word-limit
                    placeholder="请输入商品名称"
                  />
                </ElFormItem>
                <ElFormItem label="封面图" required>
                  <ImageField
                    v-model="form.image"
                    :preview-size="64"
                    default-library="system"
                  />
                </ElFormItem>
                <ElFormItem label="轮播图" required>
                  <ImagesField
                    v-model="form.slider_image"
                    :limit="10"
                    :preview-size="64"
                    default-library="system"
                  />
                </ElFormItem>
                <div class="product-edit__field-grid">
                  <ElFormItem label="商品简介" required class="product-edit__span-2">
                    <ElInput
                      v-model="form.store_info"
                      type="textarea"
                      :rows="2"
                      maxlength="500"
                      show-word-limit
                      placeholder="请输入商品简介"
                    />
                  </ElFormItem>
                  <ElFormItem label="平台分类" required>
                    <ElCascader
                      v-model="form.cate_id"
                      :options="categoryTree"
                      :props="categoryCascaderProps"
                      clearable
                      filterable
                      placeholder="请选择平台分类"
                      class="product-edit__control w-full"
                      :show-all-levels="true"
                      separator=" / "
                    />
                    <p
                      v-if="detail?.cate_id && !form.cate_id && detail?.cate_name"
                      class="field-tip"
                    >
                      原分类「{{ detail.cate_path || detail.cate_name }}」已失效，请重新选择
                    </p>
                  </ElFormItem>
                  <ElFormItem label="店铺分类" required>
                    <ElSelect
                      v-model="form.mer_cate_id"
                      clearable
                      filterable
                      placeholder="请选择店铺分类"
                      class="product-edit__control w-full"
                    >
                      <ElOption
                        v-for="item in merCateOptions"
                        :key="item.id"
                        :label="item.name"
                        :value="Number(item.id)"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="商品标签" required>
                    <ElSelect
                      v-model="form.mer_label_ids"
                      multiple
                      clearable
                      filterable
                      collapse-tags
                      collapse-tags-tooltip
                      placeholder="请选择商品标签"
                      class="product-edit__control w-full"
                    >
                      <ElOption
                        v-for="item in merLabelOptions"
                        :key="item.id"
                        :label="item.name"
                        :value="String(item.id)"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="平台标签" required>
                    <div class="product-edit__label-field" @click="openLabelSelect">
                      <div class="product-edit__label-tags">
                        <ElTag
                          v-for="name in selectedLabelNames"
                          :key="name"
                          class="mr-1"
                          size="small"
                          type="warning"
                        >
                          {{ name }}
                        </ElTag>
                        <span
                          v-if="!selectedLabelNames.length"
                          class="product-edit__label-placeholder"
                        >
                          选择标签
                        </span>
                      </div>
                      <ElButton type="primary" link @click.stop="openLabelSelect">
                        编辑标签
                      </ElButton>
                    </div>
                  </ElFormItem>
                  <ElFormItem label="品牌选择" required>
                    <ElSelect
                      v-model="form.brand_name"
                      clearable
                      filterable
                      allow-create
                      default-first-option
                      placeholder="选择或输入品牌"
                      class="product-edit__control w-full"
                    >
                      <ElOption
                        v-for="item in brandOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="单位" required>
                    <ElInput v-model="form.unit_name" maxlength="32" placeholder="如：件" />
                  </ElFormItem>
                  <ElFormItem label="关键字" required>
                    <ElInput
                      v-model="form.keyword"
                      maxlength="255"
                      placeholder="多个关键字可用空格分隔"
                    />
                  </ElFormItem>
                  <ElFormItem label="配送方式" required>
                    <ElCheckboxGroup v-model="form.delivery_way" class="product-edit__delivery">
                      <ElCheckbox :value="2" :label="2">快递配送</ElCheckbox>
                      <ElCheckbox :value="1" :label="1">到店自提</ElCheckbox>
                    </ElCheckboxGroup>
                  </ElFormItem>
                </div>
              </ElForm>
            </ElTabPane>

            <ElTabPane label="规格与价格" name="sku">
              <div class="product-edit__section-title">规格信息</div>
              <div class="product-edit__section-head">
                <span>规格类型：</span>
                <ElRadioGroup
                  :model-value="specType"
                  @update:model-value="(v) => onSpecTypeChange(Number(v) as 0 | 1)"
                >
                  <ElRadioButton :value="0">单规格</ElRadioButton>
                  <ElRadioButton :value="1">多规格</ElRadioButton>
                </ElRadioGroup>
              </div>
              <div class="product-edit__section-title">
                规格列表
                <ElButton
                  v-if="specType === 1"
                  class="ml-2"
                  type="primary"
                  link
                  @click="addSkuRow"
                >
                  新增规格
                </ElButton>
              </div>
              <ElTable :data="skuRows" border class="product-edit__sku-table">
                <ElTableColumn label="规格名称" min-width="130" fixed>
                  <template #default="{ row }">
                    <ElInput v-model="row.specLabel" placeholder="如：红色 / XL" />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="图片" width="88" align="center">
                  <template #default="{ row }">
                    <ImageField v-model="row.image" :preview-size="56" />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="售价" width="130" align="center">
                  <template #default="{ row }">
                    <ElInputNumber
                      v-model="row.price"
                      :min="0"
                      :precision="2"
                      :step="1"
                      controls-position="right"
                      class="w-full"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="划线价" width="130" align="center">
                  <template #default="{ row }">
                    <ElInputNumber
                      v-model="row.ot_price"
                      :min="0"
                      :precision="2"
                      :step="1"
                      controls-position="right"
                      class="w-full"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="库存" width="120" align="center">
                  <template #default="{ row }">
                    <ElInputNumber
                      v-model="row.stock"
                      :min="0"
                      :step="1"
                      controls-position="right"
                      class="w-full"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="规格编码" min-width="110" align="center">
                  <template #default="{ row }">
                    <ElInput v-model="row.code" maxlength="64" placeholder="编码" />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="条形码" min-width="110" align="center">
                  <template #default="{ row }">
                    <ElInput v-model="row.bar_code" maxlength="64" placeholder="条码" />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="重量(KG)" width="120" align="center">
                  <template #default="{ row }">
                    <ElInputNumber
                      v-model="row.weight"
                      :min="0"
                      :precision="2"
                      :step="0.1"
                      controls-position="right"
                      class="w-full"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="体积(m³)" width="120" align="center">
                  <template #default="{ row }">
                    <ElInputNumber
                      v-model="row.volume"
                      :min="0"
                      :precision="2"
                      :step="0.01"
                      controls-position="right"
                      class="w-full"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="一级返佣" width="120" align="center">
                  <template #default="{ row }">
                    <ElInputNumber
                      v-model="row.extension_one"
                      :min="0"
                      :precision="2"
                      :step="1"
                      controls-position="right"
                      class="w-full"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn
                  v-if="specType === 1"
                  label="操作"
                  width="90"
                  align="center"
                  fixed="right"
                >
                  <template #default="{ $index }">
                    <ElButton link type="danger" @click="removeSkuRow($index)">
                      删除
                    </ElButton>
                  </template>
                </ElTableColumn>
              </ElTable>
            </ElTabPane>

            <ElTabPane label="商品详情" name="content">
              <div class="product-edit__section-title">商品详情</div>
              <VbenTiptap
                v-model="form.content"
                :image-upload="imageUpload"
                :max-height="480"
                :min-height="320"
                :previewable="true"
                placeholder="请输入商品详情…"
              />
            </ElTabPane>

            <ElTabPane label="营销信息" name="marketing">
              <div class="product-edit__section-title">营销信息</div>
              <ElForm label-width="108px" class="product-edit__form product-edit__form--edit">
                <ElFormItem label="商品名称" required>
                  <ElInput v-model="form.title" maxlength="128" show-word-limit />
                </ElFormItem>
                <ElFormItem label="星级推荐">
                  <div class="product-edit__control-stack">
                    <ElRate v-model="form.star" :max="5" />
                    <p class="field-tip">
                      最高5星，影响搜索与列表排序展示
                    </p>
                  </div>
                </ElFormItem>
                <ElFormItem label="商品推荐">
                  <div class="product-edit__control-stack">
                    <ElCheckboxGroup
                      :model-value="recommendValues"
                      @change="onRecommendChange"
                    >
                      <ElCheckbox label="hot">热门榜单</ElCheckbox>
                      <ElCheckbox label="benefit">促销单品</ElCheckbox>
                      <ElCheckbox label="best">精品推荐</ElCheckbox>
                      <ElCheckbox label="new">首发新品</ElCheckbox>
                    </ElCheckboxGroup>
                    <p class="field-tip">勾选后商品会出现在对应推荐列表</p>
                  </div>
                </ElFormItem>
                <ElFormItem label="大图推荐">
                  <div class="product-edit__control-stack">
                    <ElSwitch
                      v-model="form.cate_hot"
                      inline-prompt
                      active-text="开启"
                      inactive-text="关闭"
                    />
                    <p class="field-tip">开启后点击商品会进入大图推荐列表页</p>
                  </div>
                </ElFormItem>
                <ElFormItem label="排序">
                  <ElInputNumber v-model="form.rank" :min="0" :max="999999" />
                </ElFormItem>
              </ElForm>
            </ElTabPane>

            <ElTabPane label="其他信息" name="other">
              <div class="product-edit__section-title">其他信息</div>
              <ElForm label-width="120px" class="product-edit__form product-edit__form--edit">
                <ElFormItem label="支持退款">
                  <ElSwitch
                    v-model="form.refund_switch"
                    inline-prompt
                    active-text="开启"
                    inactive-text="关闭"
                  />
                </ElFormItem>
                <ElFormItem label="最少购买件数">
                  <ElInputNumber v-model="form.once_min_count" :min="0" :max="9999" />
                  <span v-if="form.once_min_count <= 0" class="hint ml-2">不限购</span>
                </ElFormItem>
                <ElFormItem v-if="mode === 'edit'" label="店铺商品参数">
                  <ElTable
                    :data="detail?.mer_params || []"
                    border
                    size="small"
                    empty-text="暂无数据"
                    class="product-edit__param-table"
                  >
                    <ElTableColumn prop="name" label="参数名称" min-width="120" />
                    <ElTableColumn prop="value" label="参数值" min-width="160" />
                  </ElTable>
                </ElFormItem>
                <ElFormItem v-if="mode === 'edit'" label="平台商品参数">
                  <ElTable
                    :data="detail?.platform_params || []"
                    border
                    size="small"
                    empty-text="暂无数据"
                    class="product-edit__param-table"
                  >
                    <ElTableColumn prop="name" label="参数名称" min-width="120" />
                    <ElTableColumn prop="value" label="参数值" min-width="160" />
                  </ElTable>
                </ElFormItem>
                <ElFormItem v-if="mode === 'edit'" label="关联系统表单">
                  <span>{{ formLinkedText }}</span>
                </ElFormItem>
              </ElForm>
            </ElTabPane>

            <ElTabPane v-if="mode === 'edit' && detail" label="店铺信息" name="store">
              <div class="product-edit__section-title">店铺信息</div>
              <div class="product-edit__store-grid">
                <div class="product-edit__meta-item">
                  <span class="label">店铺名称</span>
                  <span class="value">{{ dash(detail.store_name || detail.mer_name) }}</span>
                </div>
                <div class="product-edit__meta-item">
                  <span class="label">店铺类别</span>
                  <span class="value">{{ dash(detail.mer_category_name) }}</span>
                </div>
                <div class="product-edit__meta-item">
                  <span class="label">店铺类型</span>
                  <span class="value">{{ dash(detail.mer_type_name) }}</span>
                </div>
              </div>
            </ElTabPane>

            <ElTabPane v-if="mode === 'edit' && detail" label="操作记录" name="logs">
              <div class="product-edit__log-filters">
                <div class="product-edit__log-filter">
                  <span class="filter-label">操作端</span>
                  <ElSelect
                    v-model="operateLogTerminal"
                    clearable
                    placeholder="请选择"
                    style="width: 160px"
                    @change="onOperateLogFilterChange"
                  >
                    <ElOption label="平台" value="platform" />
                    <ElOption label="运营" value="operations" />
                    <ElOption label="店铺" value="merchant" />
                  </ElSelect>
                </div>
                <div class="product-edit__log-filter">
                  <span class="filter-label">操作时间</span>
                  <ElDatePicker
                    v-model="operateLogDates"
                    type="daterange"
                    value-format="YYYY-MM-DD"
                    start-placeholder="开始时间"
                    end-placeholder="结束时间"
                    @change="onOperateLogFilterChange"
                  />
                </div>
              </div>
              <ElTable
                v-loading="operateLogLoading"
                :data="operateLogs"
                border
              >
                <ElTableColumn prop="index" label="序号" width="70" align="center" />
                <ElTableColumn prop="action_label" label="操作记录" min-width="160" />
                <ElTableColumn prop="terminal" label="操作端" width="110" align="center" />
                <ElTableColumn prop="role_name" label="操作角色" width="120" align="center" />
                <ElTableColumn prop="operator_name" label="操作人" min-width="140" />
                <ElTableColumn label="操作时间" width="180" align="center">
                  <template #default="{ row }">
                    {{ formatShanghaiDateTime(row.created_at) }}
                  </template>
                </ElTableColumn>
              </ElTable>
              <div class="product-edit__log-pager">
                <ElPagination
                  background
                  layout="total, prev, pager, next"
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
      </template>
    </ElSkeleton>

    <ProductLabelSelectModal ref="labelModalRef" @submit="onLabelSubmit" />
  </Drawer>
</template>

<style scoped>
.product-edit {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-edit__header {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid hsl(var(--border));
}

.product-edit__identity {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  min-width: 0;
}

.product-edit__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: hsl(var(--primary));
  color: #fff;
  font-size: 22px;
}

.product-edit__name {
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
  word-break: break-word;
}

.product-edit__sub {
  margin-top: 2px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.product-edit__meta {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px 16px;
}

.product-edit__meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.product-edit__meta-item .label {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.product-edit__meta-item .value {
  font-size: 14px;
  line-height: 22px;
}

.product-edit__tabs {
  min-height: 420px;
}

.product-edit__section-title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  color: hsl(var(--foreground));
  font-size: 14px;
  font-weight: 600;
}

.product-edit__section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 2px;
  background: hsl(var(--primary));
}

.product-edit__form {
  max-width: 100%;
}

.product-edit__form--edit {
  max-width: 920px;
}

.product-edit__field-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 12px;
}

.product-edit__field-grid :deep(.el-form-item) {
  margin-bottom: 12px;
  align-items: flex-start;
}

.product-edit__field-grid :deep(.el-form-item__label) {
  height: 32px;
  line-height: 32px !important;
  padding-top: 0 !important;
  align-items: center;
}

.product-edit__field-grid :deep(.el-form-item__content) {
  align-items: center;
  line-height: 32px;
  min-height: 32px;
}

.product-edit__control {
  width: 100%;
}

.product-edit__control :deep(.el-select__wrapper),
.product-edit__control :deep(.el-cascader .el-input__wrapper),
.product-edit__field-grid :deep(.el-select__wrapper),
.product-edit__field-grid :deep(.el-cascader .el-input__wrapper) {
  min-height: 32px;
  align-items: center;
}

.product-edit__control :deep(.el-select__selection),
.product-edit__control :deep(.el-select__selected-item),
.product-edit__control :deep(.el-cascader .el-input__inner),
.product-edit__field-grid :deep(.el-select__selection),
.product-edit__field-grid :deep(.el-select__selected-item),
.product-edit__field-grid :deep(.el-cascader .el-input__inner) {
  display: flex;
  align-items: center;
  line-height: 22px;
  height: auto;
}

.product-edit__delivery {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  align-items: center;
  min-height: 32px;
}

.product-edit__span-2 {
  grid-column: span 2;
}

.product-edit__sku-img {
  width: 64px;
  height: 64px;
  border-radius: 4px;
}

.product-edit__sku-table :deep(.image-field) {
  display: flex;
  justify-content: center;
}

.product-edit__label-field {
  display: flex;
  gap: 12px;
  align-items: center;
  width: 100%;
  min-height: 32px;
  padding: 4px 8px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  cursor: pointer;
}

.product-edit__label-field:hover {
  border-color: hsl(var(--primary));
}

.product-edit__label-tags {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
  min-width: 0;
}

.product-edit__label-placeholder {
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.product-edit__section-head {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 16px;
  font-size: 14px;
}

.product-edit__store-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px 24px;
}

.product-edit__control-stack {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.product-edit__param-table {
  width: 100%;
  max-width: 640px;
}

.product-edit__log-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 16px 24px;
  margin-bottom: 12px;
}

.product-edit__log-filter {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-label {
  color: hsl(var(--muted-foreground));
  font-size: 13px;
  white-space: nowrap;
}

.product-edit__log-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.field-tip {
  margin: 0;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 18px;
}

.hint {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.ml-2 {
  margin-left: 8px;
}

.mr-1 {
  margin-right: 4px;
}

.w-full {
  width: 100%;
}

@media (max-width: 960px) {
  .product-edit__meta,
  .product-edit__store-grid,
  .product-edit__field-grid {
    grid-template-columns: 1fr;
  }

  .product-edit__span-2 {
    grid-column: span 1;
  }
}
</style>
