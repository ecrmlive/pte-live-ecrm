<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { PlatformCategory, PlatformProduct } from '#/api/core/platform-catalog';
import type {
  MarketingDecor,
  MarketingDecorScopeType,
} from '#/api/core/platform-marketing-decor';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElCascader,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSwitch,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  fetchPlatformMerchants,
  fetchProductLabels,
  type ProductLabelRow,
} from '#/api/core/ecrm';
import { listPlatformCategoriesApi } from '#/api/core/platform-catalog';
import {
  createMarketingDecorApi,
  deleteMarketingDecorApi,
  listMarketingDecorApi,
  setMarketingDecorStatusApi,
  updateMarketingDecorApi,
} from '#/api/core/platform-marketing-decor';
import StorePickerModal, {
  type PickedStore,
} from '#/components/ecrm/store-picker-modal.vue';
import ImageField from '#/components/shop/image-field.vue';
import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';
import ProductLabelSelectModal from '#/views/ecrm/product/components/ProductLabelSelectModal.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import {
  LIST_DATE_RANGE_FIELD,
  buildStandardListParams,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type CascaderOption = {
  children?: CascaderOption[];
  label: string;
  value: number;
};

type PickedProduct = {
  image?: string;
  product_id: number;
  store_name?: string;
};

const EXAMPLE_PIC =
  'https://picsum.photos/seed/qixi-atmosphere-example/750/152';

const canRead = ref(false);
const canManage = ref(false);
const drawerTab = ref('basic');
const editing = ref<MarketingDecor>();
const productPickerOpen = ref(false);
const storePickerOpen = ref(false);
const labelModalRef = ref<{
  open: (payload: {
    options?: ProductLabelRow[];
    selectedIds?: string[];
  }) => void;
}>();

const categoryTree = ref<CascaderOption[]>([]);
const labelOptions = ref<ProductLabelRow[]>([]);
const pickedProducts = ref<PickedProduct[]>([]);
const pickedStores = ref<PickedStore[]>([]);
const pickedLabelIds = ref<number[]>([]);

const form = reactive({
  name: '',
  cover_url: '',
  time_range: [] as string[],
  status: 1,
  scope_type: 0 as MarketingDecorScopeType,
  cate_ids: [] as number[],
});

const cascaderProps = {
  multiple: true,
  checkStrictly: false,
  emitPath: false,
  value: 'value',
  label: 'label',
  children: 'children',
};

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '创建时间',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '未开始', value: 0 },
        { label: '进行中', value: 1 },
        { label: '已结束', value: -1 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'activity_status',
    label: '活动状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入活动名称/关键字',
    },
    fieldName: 'keyword',
    label: '活动搜索',
  },
]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues);
  const actRaw = formValues?.activity_status;
  const activityStatus =
    actRaw === 0 || actRaw === 1 || actRaw === -1 || actRaw === '0' || actRaw === '1' || actRaw === '-1'
      ? Number(actRaw)
      : undefined;
  return {
    page: base.page,
    limit: base.limit,
    keyword: base.keyword,
    date_from: base.date_from,
    date_to: base.date_to,
    activity_status: activityStatus,
  };
}

const gridOptions: VxeGridProps<MarketingDecor> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'name',
      minWidth: 160,
      showOverflow: false,
      title: '活动名称',
    },
    {
      field: 'cover_url',
      slots: { default: 'cover' },
      title: '氛围图',
      width: 100,
    },
    {
      field: 'starts_at',
      formatter: ({ row }) => formatActivityDate(row),
      minWidth: 280,
      showOverflow: 'tooltip',
      title: '活动日期',
    },
    {
      field: 'activity_status_text',
      formatter: ({ row }) =>
        row.activity_status_text ||
        ({ 0: '未开始', 1: '进行中', [-1]: '已结束' } as Record<number, string>)[
          Number(row.activity_status)
        ] ||
        '—',
      title: '活动状态',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否开启',
      width: 110,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue, row }) =>
        formatShanghaiDateTime(cellValue || row.updated_at) || '—',
      minWidth: 170,
      showOverflow: false,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listMarketingDecorApi(
          'atmosphere',
          buildParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const [ExampleModal, exampleModalApi] = useVbenModal({
  title: '氛围图示例',
  class: 'w-[820px] max-w-[96vw]',
  footer: false,
});

const pickedLabelNames = computed(() => {
  const map = new Map(labelOptions.value.map((x) => [x.id, x.name]));
  return pickedLabelIds.value.map((id) => ({
    id,
    name: map.get(id) || `标签#${id}`,
  }));
});

function coverOf(row: MarketingDecor) {
  return resolveCosMediaUrl(String(row.cover_url || '').trim());
}

function formatActivityDate(row: MarketingDecor) {
  const start = row.starts_at ? formatShanghaiDateTime(row.starts_at) : '';
  const end = row.ends_at ? formatShanghaiDateTime(row.ends_at) : '';
  if (start && end) return `${start} - ${end}`;
  return start || end || '—';
}

function toCascaderOptions(rows: PlatformCategory[] = []): CascaderOption[] {
  const out: CascaderOption[] = [];
  for (const row of rows) {
    const children = toCascaderOptions(row.children || []);
    const option: CascaderOption = {
      label: row.cate_name,
      value: Number(row.store_category_id),
    };
    if (children.length) option.children = children;
    out.push(option);
  }
  return out;
}

function resetForm() {
  editing.value = undefined;
  drawerTab.value = 'basic';
  Object.assign(form, {
    name: '',
    cover_url: '',
    time_range: [],
    status: 1,
    scope_type: 0,
    cate_ids: [],
  });
  pickedProducts.value = [];
  pickedStores.value = [];
  pickedLabelIds.value = [];
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增活动氛围图' }).open();
}

function openEdit(row: MarketingDecor) {
  editing.value = row;
  drawerTab.value = 'basic';
  const scopeType = Number(row.scope_type ?? row.payload?.scope_type ?? 0) as MarketingDecorScopeType;
  Object.assign(form, {
    name: row.name || '',
    cover_url: row.cover_url || '',
    time_range:
      row.starts_at && row.ends_at ? [row.starts_at, row.ends_at] : [],
    status: row.status === 1 ? 1 : 0,
    scope_type: ([0, 1, 2, 3, 4].includes(scopeType) ? scopeType : 0) as MarketingDecorScopeType,
    cate_ids: [...(row.cate_ids || asNumberIds(row.payload?.cate_ids))],
  });
  pickedProducts.value = (row.spu_ids || asNumberIds(row.payload?.spu_ids)).map(
    (id) => ({ product_id: id, store_name: `商品#${id}` }),
  );
  pickedStores.value = (row.mer_ids || asNumberIds(row.payload?.mer_ids)).map(
    (id) => ({
      mer_id: id,
      mer_name: `店铺#${id}`,
      real_name: '',
      mer_phone: '',
    }),
  );
  pickedLabelIds.value = [
    ...(row.label_ids || asNumberIds(row.payload?.label_ids)),
  ];
  void hydrateScopeNames();
  formDrawerApi.setState({ title: '编辑活动氛围图' }).open();
}

function asNumberIds(raw: unknown): number[] {
  if (!Array.isArray(raw)) return [];
  return raw
    .map((x) => Number(x))
    .filter((n) => Number.isFinite(n) && n > 0);
}

async function hydrateScopeNames() {
  if (form.scope_type === 3 && pickedStores.value.length) {
    try {
      const res = await fetchPlatformMerchants({
        page: 1,
        limit: 100,
        status: 1,
      });
      const map = new Map(
        (res.list || []).map((m) => [m.mer_id, m.mer_name || `店铺#${m.mer_id}`]),
      );
      pickedStores.value = pickedStores.value.map((s) => ({
        ...s,
        mer_name: map.get(s.mer_id) || s.mer_name,
      }));
    } catch {
      /* keep placeholders */
    }
  }
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请填写活动名称');
    drawerTab.value = 'basic';
    return;
  }
  if (!form.time_range?.[0] || !form.time_range?.[1]) {
    ElMessage.warning('请设置活动时间');
    drawerTab.value = 'basic';
    return;
  }
  if (form.scope_type === 1 && !pickedProducts.value.length) {
    ElMessage.warning('请选择指定商品');
    drawerTab.value = 'scope';
    return;
  }
  if (form.scope_type === 2 && !form.cate_ids.length) {
    ElMessage.warning('请选择指定分类');
    drawerTab.value = 'scope';
    return;
  }
  if (form.scope_type === 3 && !pickedStores.value.length) {
    ElMessage.warning('请选择指定店铺');
    drawerTab.value = 'scope';
    return;
  }
  if (form.scope_type === 4 && !pickedLabelIds.value.length) {
    ElMessage.warning('请选择指定商品标签');
    drawerTab.value = 'scope';
    return;
  }

  const body = {
    name,
    cover_url: form.cover_url.trim(),
    starts_at: form.time_range[0],
    ends_at: form.time_range[1],
    status: form.status === 1 ? 1 : 0,
    scope_type: form.scope_type,
    spu_ids:
      form.scope_type === 1
        ? pickedProducts.value.map((p) => p.product_id)
        : [],
    cate_ids: form.scope_type === 2 ? [...form.cate_ids] : [],
    mer_ids:
      form.scope_type === 3 ? pickedStores.value.map((s) => s.mer_id) : [],
    label_ids: form.scope_type === 4 ? [...pickedLabelIds.value] : [],
    payload: {},
  };

  formDrawerApi.lock();
  try {
    if (editing.value) {
      await updateMarketingDecorApi('atmosphere', editing.value.id, body);
    } else {
      await createMarketingDecorApi('atmosphere', body);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: MarketingDecor, enabled: boolean) {
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await setMarketingDecorStatusApi(
      'atmosphere',
      row.id,
      enabled ? 1 : 0,
    );
  } catch {
    row.status = before;
  }
}

async function removeRow(row: MarketingDecor) {
  try {
    await confirm({
      content: `确定删除活动氛围图「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteMarketingDecorApi('atmosphere', row.id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

function openExample() {
  exampleModalApi.open();
}

function onProductsPicked(list: PlatformProduct[]) {
  const map = new Map(pickedProducts.value.map((p) => [p.product_id, p]));
  for (const p of list) {
    map.set(p.product_id, {
      product_id: p.product_id,
      store_name: p.store_name || p.title || `商品#${p.product_id}`,
      image: p.image,
    });
  }
  pickedProducts.value = [...map.values()];
}

function removeProduct(id: number) {
  pickedProducts.value = pickedProducts.value.filter((p) => p.product_id !== id);
}

function onStoresPicked(stores: PickedStore[]) {
  const map = new Map(pickedStores.value.map((s) => [s.mer_id, s]));
  for (const s of stores) map.set(s.mer_id, s);
  pickedStores.value = [...map.values()];
}

function removeStore(id: number) {
  pickedStores.value = pickedStores.value.filter((s) => s.mer_id !== id);
}

function openLabelPicker() {
  labelModalRef.value?.open({
    selectedIds: pickedLabelIds.value.map(String),
    options: labelOptions.value,
  });
}

function onLabelsPicked(ids: string[]) {
  pickedLabelIds.value = ids.map(Number).filter((n) => n > 0);
}

function removeLabel(id: number) {
  pickedLabelIds.value = pickedLabelIds.value.filter((x) => x !== id);
}

onMounted(async () => {
  const [profile, codes, cateRes, labelRes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    listPlatformCategoriesApi().catch(() => ({ list: [] as PlatformCategory[] })),
    fetchProductLabels().catch(() => ({ list: [] as ProductLabelRow[] })),
  ]);
  const cateList = Array.isArray(cateRes)
    ? cateRes
    : (cateRes as { list?: PlatformCategory[] })?.list || [];
  categoryTree.value = toCascaderOptions(cateList);
  labelOptions.value = (labelRes.list || []).filter((x) => Number(x.status) !== 0);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('marketing.atmosphere.read') ||
      codes.includes('marketing.atmosphere.manage'));
  canManage.value = roleOK && codes.includes('marketing.atmosphere.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增活动氛围图
        </ElButton>
      </template>
      <template #cover="{ row }">
        <ElImage
          v-if="coverOf(row)"
          :src="coverOf(row)"
          fit="cover"
          class="atmosphere-cover"
          :preview-src-list="[coverOf(row)]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="开启"
          inactive-text="关闭"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
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
          type="primary"
          @click="removeRow(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElTabs v-model="drawerTab">
        <ElTabPane label="基础设置" name="basic">
          <ElForm label-width="110px" class="atmosphere-form">
            <ElFormItem label="活动名称" required>
              <ElInput
                v-model="form.name"
                maxlength="128"
                show-word-limit
                placeholder="请输入活动名称"
              />
            </ElFormItem>
            <ElFormItem label="活动时间" required>
              <div class="w-full">
                <ElDatePicker
                  v-model="form.time_range"
                  type="datetimerange"
                  value-format="YYYY-MM-DD HH:mm:ss"
                  start-placeholder="开始时间"
                  end-placeholder="结束时间"
                  class="w-full"
                />
                <div class="field-hint">设置活动氛围图在商城展示时间</div>
              </div>
            </ElFormItem>
            <ElFormItem label="活动氛围图">
              <div class="atmosphere-pic">
                <ImageField v-model="form.cover_url" />
                <div class="atmosphere-pic__side">
                  <ElButton link type="primary" @click="openExample">
                    查看示例
                  </ElButton>
                  <div class="field-hint">宽750px，高152px</div>
                </div>
              </div>
            </ElFormItem>
            <ElFormItem label="是否开启">
              <ElSwitch
                v-model="form.status"
                :active-value="1"
                :inactive-value="0"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
              />
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="使用范围" name="scope">
          <ElForm label-width="110px" class="atmosphere-form">
            <ElFormItem label="使用范围" required>
              <ElRadioGroup v-model="form.scope_type" class="scope-radios">
                <ElRadio :label="0">全部商品参与</ElRadio>
                <ElRadio :label="1">指定商品参与</ElRadio>
                <ElRadio :label="2">指定分类参与</ElRadio>
                <ElRadio :label="3">指定店铺参与</ElRadio>
                <ElRadio :label="4">指定商品标签</ElRadio>
              </ElRadioGroup>
            </ElFormItem>

            <ElFormItem v-if="form.scope_type === 1" label="选择商品" required>
              <div class="scope-picker">
                <ElButton type="primary" plain @click="productPickerOpen = true">
                  选择商品
                </ElButton>
                <div v-if="pickedProducts.length" class="scope-tags">
                  <ElTag
                    v-for="p in pickedProducts"
                    :key="p.product_id"
                    closable
                    class="scope-tag"
                    @close="removeProduct(p.product_id)"
                  >
                    {{ p.store_name || `商品#${p.product_id}` }}
                  </ElTag>
                </div>
              </div>
            </ElFormItem>

            <ElFormItem v-if="form.scope_type === 2" label="选择分类" required>
              <ElCascader
                v-model="form.cate_ids"
                class="w-full"
                :options="categoryTree"
                :props="cascaderProps"
                clearable
                filterable
                placeholder="请选择分类"
              />
            </ElFormItem>

            <ElFormItem v-if="form.scope_type === 3" label="选择店铺" required>
              <div class="scope-picker">
                <ElButton type="primary" plain @click="storePickerOpen = true">
                  选择店铺
                </ElButton>
                <div v-if="pickedStores.length" class="scope-tags">
                  <ElTag
                    v-for="s in pickedStores"
                    :key="s.mer_id"
                    closable
                    class="scope-tag"
                    @close="removeStore(s.mer_id)"
                  >
                    {{ s.mer_name || `店铺#${s.mer_id}` }}
                  </ElTag>
                </div>
              </div>
            </ElFormItem>

            <ElFormItem v-if="form.scope_type === 4" label="选择标签" required>
              <div class="scope-picker">
                <ElButton type="primary" plain @click="openLabelPicker">
                  选择标签
                </ElButton>
                <div v-if="pickedLabelNames.length" class="scope-tags">
                  <ElTag
                    v-for="item in pickedLabelNames"
                    :key="item.id"
                    closable
                    class="scope-tag"
                    @close="removeLabel(item.id)"
                  >
                    {{ item.name }}
                  </ElTag>
                </div>
              </div>
            </ElFormItem>
          </ElForm>
        </ElTabPane>
      </ElTabs>
    </FormDrawer>

    <ExampleModal>
      <div class="example-wrap">
        <ElImage :src="EXAMPLE_PIC" fit="contain" class="example-img" />
        <div class="field-hint">建议尺寸：宽 750px × 高 152px</div>
      </div>
    </ExampleModal>

    <ProductPickerDialog
      v-model:open="productPickerOpen"
      multiple
      @confirm="onProductsPicked"
    />
    <StorePickerModal
      v-model:open="storePickerOpen"
      :selected="pickedStores"
      @confirm="onStoresPicked"
    />
    <ProductLabelSelectModal ref="labelModalRef" @submit="onLabelsPicked" />
  </Page>
</template>

<style scoped>
.atmosphere-cover {
  display: block;
  width: 64px;
  height: 36px;
  border-radius: 4px;
}

.atmosphere-form {
  padding-top: 8px;
}

.field-hint {
  margin-top: 6px;
  color: hsl(var(--foreground) / 55%);
  font-size: 12px;
  line-height: 1.4;
}

.atmosphere-pic {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-start;
}

.atmosphere-pic__side {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 4px;
}

.scope-radios {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.scope-picker {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.scope-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.scope-tag {
  max-width: 100%;
}

.example-wrap {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: center;
  padding: 8px 0 16px;
}

.example-img {
  width: 100%;
  max-width: 750px;
  min-height: 120px;
  border-radius: 4px;
  background: hsl(var(--accent));
}
</style>
