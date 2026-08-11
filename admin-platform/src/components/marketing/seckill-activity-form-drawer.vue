<script setup lang="ts">
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { computed, reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { ArrowDown } from '@element-plus/icons-vue';
import {
  ElButton,
  ElDatePicker,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElIcon,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElPagination,
  ElPopover,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
} from 'element-plus';

import {
  createPlatformSeckillActivityApi,
  getPlatformSeckillActivityApi,
  listPlatformSeckillActivityGoodsApi,
  savePlatformSeckillActivityProductsApi,
  updatePlatformSeckillActivityApi,
  type PlatformSeckillActivity,
  type PlatformSeckillActivityGoods,
  type PlatformSeckillTime,
} from '#/api/core/platform-seckill';
import { getPlatformProductEditApi } from '#/api/core/platform-catalog';
import ImageField from '#/components/shop/image-field.vue';
import SeckillActivityProductPicker from '#/components/marketing/seckill-activity-product-picker.vue';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

export type ActivityFormMode = 'view' | 'edit' | 'create';

interface DraftSKU {
  sku: string;
  image?: string;
  price: number;
  seckill_price: number;
  stock: number;
  limit_stock: number;
}

interface DraftProduct {
  _key: string;
  product_id: number;
  store_name: string;
  image?: string;
  cate_name?: string;
  mer_id: number;
  mer_name?: string;
  price: number;
  product_stock: number;
  seckill_price: number;
  stock: number;
  status: number;
  sort: number;
  children: DraftSKU[];
}

const props = defineProps<{
  timeOptions: PlatformSeckillTime[];
  categoryOptions: { label: string; value: number }[];
}>();

const emit = defineEmits<{
  saved: [];
}>();

const mode = ref<ActivityFormMode>('create');
const activeTab = ref('basic');
const saving = ref(false);
const editingID = ref<number>();
const goodsLoading = ref(false);
const goodsRows = ref<PlatformSeckillActivityGoods[]>([]);
const goodsTotal = ref(0);
const goodsFilter = reactive({
  product_status: undefined as number | undefined,
  keyword: '',
  page: 1,
  limit: 10,
});

const draftRows = ref<DraftProduct[]>([]);
const draftSelection = ref<DraftProduct[]>([]);
const pickerOpen = ref(false);

const form = reactive({
  name: '',
  time_ids: [] as number[],
  date_range: [] as string[],
  once_pay_count: 1,
  all_pay_count: 0,
  category_ids: [] as number[],
  border_pic: '',
  status: 1,
});

const readonly = computed(() => mode.value === 'view');
const showDraftTab = computed(
  () => mode.value === 'create' || mode.value === 'edit',
);
const showGoodsTab = computed(
  () => mode.value === 'view' || mode.value === 'edit',
);
const excludeProductIds = computed(() =>
  draftRows.value.map((r) => r.product_id),
);
const batchDisabled = computed(() => draftSelection.value.length === 0);

const drawerTitle = computed(() => {
  if (mode.value === 'view') return '查看秒杀活动';
  if (mode.value === 'edit') return '编辑秒杀活动';
  return '新增秒杀活动';
});

function formatHourLabel(start: number, end: number) {
  const startLabel = start === 0 ? '00:00' : `${start}:00`;
  return `${startLabel} - ${end}:00`;
}

const timeSelectOptions = computed(() =>
  props.timeOptions.map((t) => ({
    label: `${t.title} | ${formatHourLabel(t.start_time, t.end_time)}`,
    value: t.seckill_time_id,
  })),
);

function toDay(value?: string) {
  const s = String(value || '').trim();
  if (s.length >= 10 && s[4] === '-' && s[7] === '-') return s.slice(0, 10);
  return s;
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    name: '',
    time_ids: [],
    date_range: [],
    once_pay_count: 1,
    all_pay_count: 0,
    category_ids: [],
    border_pic: '',
    status: 1,
  });
  draftRows.value = [];
  draftSelection.value = [];
}

function fillForm(row: PlatformSeckillActivity) {
  editingID.value = row.seckill_activity_id;
  const times = String(row.seckill_time_ids || '')
    .split(',')
    .map((x) => Number(x.trim()))
    .filter((n) => n > 0);
  const cats = String(row.product_category_ids || '')
    .split(',')
    .map((x) => Number(x.trim()))
    .filter((n) => n > 0);
  Object.assign(form, {
    name: row.name || '',
    time_ids: times,
    date_range: [toDay(row.start_day), toDay(row.end_day)].filter(Boolean),
    once_pay_count: Number(row.once_pay_count || 1),
    all_pay_count: Number(row.all_pay_count || 0),
    category_ids: cats,
    border_pic: row.border_pic || '',
    status: Number(row.status ?? 1),
  });
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  placement: 'right',
  onConfirm: () => void save(),
  onCancel: () => {
    if (
      showDraftTab.value &&
      activeTab.value === 'draft' &&
      !readonly.value
    ) {
      activeTab.value = 'basic';
      return;
    }
    drawerApi.close();
  },
});

function syncFooterState() {
  if (readonly.value) {
    drawerApi.setState({
      showConfirmButton: false,
      cancelText: '关闭',
    });
    return;
  }
  const onDraft = activeTab.value === 'draft';
  drawerApi.setState({
    showConfirmButton: true,
    confirmText: '保存',
    cancelText: onDraft ? '上一步' : '取消',
  });
}

watch(activeTab, () => {
  syncFooterState();
  if (activeTab.value === 'goods' && showGoodsTab.value) {
    void loadGoods();
  }
});

async function loadGoods() {
  if (!editingID.value) {
    goodsRows.value = [];
    goodsTotal.value = 0;
    return;
  }
  goodsLoading.value = true;
  try {
    const data = await listPlatformSeckillActivityGoodsApi(editingID.value, {
      page: goodsFilter.page,
      limit: goodsFilter.limit,
      keyword: goodsFilter.keyword.trim() || undefined,
      product_status:
        goodsFilter.product_status === 0 ||
        goodsFilter.product_status === 1 ||
        goodsFilter.product_status === -1
          ? goodsFilter.product_status
          : undefined,
    });
    goodsRows.value = data.list || [];
    goodsTotal.value = data.total || 0;
  } catch {
    ElMessage.error('加载已加商品失败');
  } finally {
    goodsLoading.value = false;
  }
}

function searchGoods() {
  goodsFilter.page = 1;
  void loadGoods();
}

function resetGoodsFilter() {
  goodsFilter.product_status = undefined;
  goodsFilter.keyword = '';
  goodsFilter.page = 1;
  void loadGoods();
}

async function open(next: ActivityFormMode, row?: PlatformSeckillActivity) {
  mode.value = next;
  activeTab.value = 'basic';
  goodsFilter.product_status = undefined;
  goodsFilter.keyword = '';
  goodsFilter.page = 1;
  goodsRows.value = [];
  goodsTotal.value = 0;
  draftRows.value = [];
  draftSelection.value = [];

  if (next === 'create') {
    resetForm();
    drawerApi
      .setState({
        title: drawerTitle.value,
        showConfirmButton: true,
        confirmText: '保存',
        cancelText: '取消',
      })
      .open();
    return;
  }

  if (!row?.seckill_activity_id) {
    ElMessage.error('活动不存在');
    return;
  }
  try {
    const detail = await getPlatformSeckillActivityApi(row.seckill_activity_id);
    fillForm(detail);
    drawerApi
      .setState({
        title: drawerTitle.value,
        showConfirmButton: next !== 'view',
        confirmText: '保存',
        cancelText: next === 'view' ? '关闭' : '取消',
      })
      .open();
  } catch {
    ElMessage.error('加载活动失败');
  }
}

function validateBasic(): boolean {
  if (!form.name.trim()) {
    ElMessage.warning('请输入活动名称');
    activeTab.value = 'basic';
    return false;
  }
  if (!form.time_ids.length) {
    ElMessage.warning('请选择活动场次');
    activeTab.value = 'basic';
    return false;
  }
  const start = form.date_range?.[0];
  const end = form.date_range?.[1];
  if (!start || !end) {
    ElMessage.warning('请选择活动日期');
    activeTab.value = 'basic';
    return false;
  }
  if (
    form.all_pay_count > 0 &&
    form.once_pay_count > 0 &&
    form.once_pay_count > form.all_pay_count
  ) {
    ElMessage.warning('单次限购不能大于活动限购');
    activeTab.value = 'basic';
    return false;
  }
  for (const row of draftRows.value) {
    if (!(row.seckill_price > 0)) {
      ElMessage.warning(`请填写商品 ${row.product_id} 的秒杀价`);
      activeTab.value = 'draft';
      return false;
    }
  }
  return true;
}

function buildActivityBody() {
  const start = form.date_range?.[0];
  const end = form.date_range?.[1];
  return {
    name: form.name.trim(),
    seckill_time_ids: form.time_ids.join(','),
    start_day: toDay(start),
    end_day: toDay(end),
    once_pay_count: Math.max(0, Number(form.once_pay_count) || 0),
    all_pay_count: Math.max(0, Number(form.all_pay_count) || 0),
    product_category_ids: form.category_ids.join(','),
    border_pic: form.border_pic || '',
    status: form.status,
  };
}

async function save() {
  if (readonly.value) return;
  if (!validateBasic()) return;

  drawerApi.lock();
  saving.value = true;
  try {
    const body = buildActivityBody();
    let activityId = editingID.value;
    if (activityId) {
      await updatePlatformSeckillActivityApi(activityId, body);
    } else {
      const created = await createPlatformSeckillActivityApi(body);
      activityId = created.seckill_activity_id;
      editingID.value = activityId;
    }

    if (draftRows.value.length && activityId) {
      await savePlatformSeckillActivityProductsApi(
        activityId,
        draftRows.value.map((r) => ({
          product_id: r.product_id,
          seckill_price: Number(r.seckill_price),
          stock: Number(r.stock) || 0,
          status: Number(r.status) ? 1 : 0,
          sort: Number(r.sort) || 0,
        })),
      );
    }

    ElMessage.success(mode.value === 'create' ? '已添加秒杀活动' : '已更新秒杀活动');
    drawerApi.close();
    emit('saved');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    saving.value = false;
    drawerApi.unlock();
  }
}

function mediaUrl(path?: string) {
  return resolveCosMediaUrl(String(path || '').trim());
}

function formatMoney(value?: number) {
  const n = Number(value || 0);
  return n.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

function onDraftSelection(list: DraftProduct[]) {
  draftSelection.value = list;
}

function removeDraft(row: DraftProduct) {
  draftRows.value = draftRows.value.filter((r) => r._key !== row._key);
}

function batchDeleteDraft() {
  if (!draftSelection.value.length) return;
  const keys = new Set(draftSelection.value.map((r) => r._key));
  draftRows.value = draftRows.value.filter((r) => !keys.has(r._key));
  draftSelection.value = [];
}

async function applyBatchSet(cmd: string) {
  if (!draftSelection.value.length) {
    ElMessage.warning('请先勾选商品');
    return;
  }
  try {
    if (cmd === 'price') {
      const { value } = await ElMessageBox.prompt('请输入秒杀价', '批量设置', {
        inputPattern: /^\d+(\.\d{1,2})?$/,
        inputErrorMessage: '请输入有效金额',
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      });
      const price = Number(value);
      for (const row of draftSelection.value) {
        row.seckill_price = price;
        for (const sku of row.children) sku.seckill_price = price;
      }
    } else if (cmd === 'stock') {
      const { value } = await ElMessageBox.prompt('请输入限量', '批量设置', {
        inputPattern: /^\d+$/,
        inputErrorMessage: '请输入非负整数',
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      });
      const stock = Number(value);
      for (const row of draftSelection.value) {
        row.stock = stock;
        for (const sku of row.children) sku.limit_stock = stock;
      }
    } else if (cmd === 'sort') {
      const { value } = await ElMessageBox.prompt('请输入排序', '批量设置', {
        inputPattern: /^\d+$/,
        inputErrorMessage: '请输入非负整数',
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      });
      const sort = Number(value);
      for (const row of draftSelection.value) row.sort = sort;
    } else if (cmd === 'on') {
      for (const row of draftSelection.value) row.status = 1;
    } else if (cmd === 'off') {
      for (const row of draftSelection.value) row.status = 0;
    }
  } catch {
    /* cancel */
  }
}

async function onProductsPicked(products: PlatformProduct[]) {
  const existing = new Set(draftRows.value.map((r) => r.product_id));
  for (const p of products) {
    if (existing.has(p.product_id)) continue;
    const draft: DraftProduct = {
      _key: `p-${p.product_id}-${Date.now()}`,
      product_id: p.product_id,
      store_name: p.store_name || p.title || '',
      image: p.image,
      cate_name: p.cate_name,
      mer_id: p.mer_id,
      mer_name: p.mer_name,
      price: Number(p.price || 0),
      product_stock: Number(p.stock || 0),
      seckill_price: Number(p.price || 0),
      stock: Number(p.stock || 0),
      status: 1,
      sort: 0,
      children: [
        {
          sku: '默认',
          image: p.image,
          price: Number(p.price || 0),
          seckill_price: Number(p.price || 0),
          stock: Number(p.stock || 0),
          limit_stock: Number(p.stock || 0),
        },
      ],
    };
    draftRows.value.push(draft);
    existing.add(p.product_id);
    void enrichDraftSKUs(draft);
  }
  if (products.length) {
    activeTab.value = 'draft';
  }
}

async function enrichDraftSKUs(draft: DraftProduct) {
  try {
    const detail = await getPlatformProductEditApi(draft.product_id);
    const skus = detail.skus || [];
    if (!skus.length) return;
    draft.price = Number(detail.price || draft.price);
    draft.product_stock = Number(detail.stock || draft.product_stock);
    draft.children = skus.map((sku) => ({
      sku: sku.spec_text || '默认',
      image: sku.image || draft.image,
      price: Number(sku.price || 0),
      seckill_price: Number(draft.seckill_price),
      stock: Number(sku.stock || 0),
      limit_stock: Number(draft.stock || 0),
    }));
  } catch {
    /* keep default sku */
  }
}

function goNextFromBasic() {
  if (!validateBasic()) return;
  if (showDraftTab.value) {
    activeTab.value = 'draft';
    return;
  }
  void save();
}

defineExpose({ open });
</script>

<template>
  <Drawer>
    <ElTabs v-model="activeTab" class="activity-tabs">
      <ElTabPane label="基础设置" name="basic">
        <ElForm label-width="120px" class="activity-form" v-loading="saving">
          <ElFormItem label="活动名称" required>
            <ElInput
              v-model="form.name"
              maxlength="64"
              show-word-limit
              placeholder="请输入活动名称"
              class="form-control"
              :disabled="readonly"
            />
          </ElFormItem>
          <ElFormItem label="活动日期" required>
            <div>
              <ElDatePicker
                v-model="form.date_range"
                type="daterange"
                value-format="YYYY-MM-DD"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                class="form-control"
                :disabled="readonly"
              />
              <div class="form-tip">
                设置活动开始日期与结束日期，用户可以在有效时间内参与秒杀
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="秒杀场次" required>
            <div class="form-control">
              <ElSelect
                v-model="form.time_ids"
                multiple
                collapse-tags
                collapse-tags-tooltip
                placeholder="请选择秒杀场次"
                class="w-full"
                :disabled="readonly"
              >
                <ElOption
                  v-for="opt in timeSelectOptions"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </ElSelect>
              <div class="form-tip">
                选择商品开始时间段，该时间段内用户可参与购买；其它时间段会显示活动未开始或已结束，可多选
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="活动限购">
            <div>
              <ElInputNumber
                v-model="form.all_pay_count"
                :min="0"
                :max="99999"
                :controls="false"
                class="form-control"
                :disabled="readonly"
              />
              <div class="form-tip">
                活动有效期内每个用户可购买该商品总数限制。例如设置为
                4，表示本次活动有效期内每个用户最多可购买总数
                4 个；0 为不限购
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="单次限购">
            <div>
              <ElInputNumber
                v-model="form.once_pay_count"
                :min="0"
                :max="99999"
                :controls="false"
                class="form-control"
                :disabled="readonly"
              />
              <div class="form-tip">
                用户参与秒杀时一次购买最大数量限制。例如设置为
                2，表示一次最多选 2 个；0 为不限购
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="商品范围">
            <div class="form-control">
              <ElSelect
                v-model="form.category_ids"
                multiple
                filterable
                clearable
                collapse-tags
                collapse-tags-tooltip
                placeholder="请选择商品分类（可选）"
                class="w-full"
                :disabled="readonly"
              >
                <ElOption
                  v-for="opt in categoryOptions"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </ElSelect>
              <div class="form-tip">
                设置秒杀活动可以参与的商品分类，可多选，不选为全品类商品。
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="活动边框图">
            <div>
              <div class="border-pic-row">
                <ImageField
                  v-model="form.border_pic"
                  :disabled="readonly"
                  :preview-size="86"
                  default-library="system"
                />
                <div class="border-pic-meta">
                  <div class="form-tip !mt-0">宽 750px，高 750px</div>
                  <ElPopover placement="bottom-start" :width="220" trigger="hover">
                    <template #reference>
                      <ElButton link type="primary">查看示例</ElButton>
                    </template>
                    <div class="example-tip">
                      商品列表角标/边框示意图（750×750）。上传后展示在商品列表的活动边框图位置。
                    </div>
                  </ElPopover>
                </div>
              </div>
              <div class="form-tip">展示在商品列表的活动边框图</div>
            </div>
          </ElFormItem>
          <ElFormItem label="是否开启">
            <ElSwitch
              :model-value="form.status === 1"
              :disabled="readonly"
              @change="(v: string | number | boolean) => (form.status = v ? 1 : 0)"
            />
          </ElFormItem>
        </ElForm>

        <div v-if="!readonly && showDraftTab" class="tab-actions">
          <ElButton type="primary" @click="goNextFromBasic">下一步</ElButton>
        </div>
      </ElTabPane>

      <ElTabPane v-if="showDraftTab" label="秒杀商品" name="draft">
        <div class="draft-toolbar">
          <ElButton type="primary" @click="pickerOpen = true">新增商品</ElButton>
          <ElDropdown :disabled="batchDisabled" @command="applyBatchSet">
            <ElButton :disabled="batchDisabled">
              批量设置
              <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
            </ElButton>
            <template #dropdown>
              <ElDropdownMenu>
                <ElDropdownItem command="price">批量设置秒杀价</ElDropdownItem>
                <ElDropdownItem command="stock">批量设置限量</ElDropdownItem>
                <ElDropdownItem command="sort">批量设置排序</ElDropdownItem>
                <ElDropdownItem command="on">批量开启</ElDropdownItem>
                <ElDropdownItem command="off">批量关闭</ElDropdownItem>
              </ElDropdownMenu>
            </template>
          </ElDropdown>
          <ElButton :disabled="batchDisabled" @click="batchDeleteDraft">
            批量删除
          </ElButton>
        </div>

        <ElTable
          :data="draftRows"
          row-key="_key"
          class="goods-table"
          @selection-change="onDraftSelection"
        >
          <ElTableColumn type="selection" width="48" />
          <ElTableColumn type="expand">
            <template #default="{ row }">
              <div class="sku-wrap">
                <table class="sku-table">
                  <thead>
                    <tr>
                      <th>规格</th>
                      <th>售价</th>
                      <th>秒杀价</th>
                      <th>库存</th>
                      <th>限量</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="(sku, idx) in row.children"
                      :key="`${row._key}-${idx}`"
                    >
                      <td>
                        <div class="sku-info">
                          <ElImage
                            v-if="mediaUrl(sku.image)"
                            :src="mediaUrl(sku.image)"
                            fit="cover"
                            class="sku-info__img"
                          />
                          <span>{{ sku.sku || '默认' }}</span>
                        </div>
                      </td>
                      <td>{{ formatMoney(sku.price) }}</td>
                      <td>
                        <ElInputNumber
                          v-model="sku.seckill_price"
                          :min="0.01"
                          :precision="2"
                          :controls="false"
                          class="cell-input"
                          @change="() => (row.seckill_price = sku.seckill_price)"
                        />
                      </td>
                      <td>{{ sku.stock }}</td>
                      <td>
                        <ElInputNumber
                          v-model="sku.limit_stock"
                          :min="0"
                          :controls="false"
                          class="cell-input"
                          @change="() => (row.stock = sku.limit_stock)"
                        />
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="product_id" label="商品ID" min-width="90" />
          <ElTableColumn label="商品信息" min-width="200">
            <template #default="{ row }">
              <div class="goods-info">
                <ElImage
                  v-if="mediaUrl(row.image)"
                  :src="mediaUrl(row.image)"
                  fit="cover"
                  class="goods-info__img"
                />
                <div v-else class="goods-info__img is-empty" />
                <div class="goods-info__name">{{ row.store_name || '—' }}</div>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn label="商品分类" min-width="100">
            <template #default="{ row }">{{ row.cate_name || '—' }}</template>
          </ElTableColumn>
          <ElTableColumn label="店铺名称" min-width="100">
            <template #default="{ row }">
              {{ row.mer_name || row.mer_id || '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="售价" min-width="90">
            <template #default="{ row }">{{ formatMoney(row.price) }}</template>
          </ElTableColumn>
          <ElTableColumn label="秒杀价" min-width="110">
            <template #default="{ row }">
              <ElInputNumber
                v-model="row.seckill_price"
                :min="0.01"
                :precision="2"
                :controls="false"
                class="cell-input"
              />
            </template>
          </ElTableColumn>
          <ElTableColumn label="库存" min-width="80">
            <template #default="{ row }">{{ row.product_stock }}</template>
          </ElTableColumn>
          <ElTableColumn label="限量" min-width="100">
            <template #default="{ row }">
              <ElInputNumber
                v-model="row.stock"
                :min="0"
                :controls="false"
                class="cell-input"
              />
            </template>
          </ElTableColumn>
          <ElTableColumn label="是否开启" min-width="90">
            <template #default="{ row }">
              <ElSwitch
                :model-value="row.status === 1"
                @change="(v: string | number | boolean) => (row.status = v ? 1 : 0)"
              />
            </template>
          </ElTableColumn>
          <ElTableColumn label="排序" min-width="90">
            <template #default="{ row }">
              <ElInputNumber
                v-model="row.sort"
                :min="0"
                :controls="false"
                class="cell-input"
              />
            </template>
          </ElTableColumn>
          <ElTableColumn label="操作" width="80" fixed="right">
            <template #default="{ row }">
              <ElButton link type="danger" @click="removeDraft(row)">
                删除
              </ElButton>
            </template>
          </ElTableColumn>
        </ElTable>
      </ElTabPane>

      <ElTabPane v-if="showGoodsTab" label="已加商品" name="goods">
        <div class="goods-filter">
          <ElSelect
            v-model="goodsFilter.product_status"
            clearable
            placeholder="商品审核状态"
            class="goods-filter__status"
          >
            <ElOption label="待审核" :value="0" />
            <ElOption label="审核通过" :value="1" />
            <ElOption label="审核失败" :value="-1" />
          </ElSelect>
          <ElInput
            v-model="goodsFilter.keyword"
            clearable
            placeholder="请输入商品名称或ID"
            class="goods-filter__keyword"
            @keyup.enter="searchGoods"
          />
          <ElButton @click="resetGoodsFilter">重置</ElButton>
          <ElButton type="primary" @click="searchGoods">搜索</ElButton>
        </div>

        <ElTable
          v-loading="goodsLoading"
          :data="goodsRows"
          row-key="seckill_active_id"
          class="goods-table"
        >
          <ElTableColumn type="expand">
            <template #default="{ row }">
              <div class="sku-wrap">
                <table class="sku-table">
                  <thead>
                    <tr>
                      <th>规格</th>
                      <th>售价</th>
                      <th>秒杀价</th>
                      <th>库存</th>
                      <th>限量</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="(sku, idx) in row.children || []"
                      :key="`${row.seckill_active_id}-${idx}`"
                    >
                      <td>
                        <div class="sku-info">
                          <ElImage
                            v-if="mediaUrl(sku.image || row.image)"
                            :src="mediaUrl(sku.image || row.image)"
                            fit="cover"
                            class="sku-info__img"
                          />
                          <span>{{ sku.sku || '默认' }}</span>
                        </div>
                      </td>
                      <td>{{ formatMoney(sku.price ?? row.price) }}</td>
                      <td>{{ formatMoney(sku.seckill_price) }}</td>
                      <td>{{ sku.stock }}</td>
                      <td>{{ sku.limit_stock }}</td>
                    </tr>
                    <tr v-if="!(row.children || []).length">
                      <td>默认</td>
                      <td>{{ formatMoney(row.price) }}</td>
                      <td>{{ formatMoney(row.seckill_price) }}</td>
                      <td>{{ row.product_stock }}</td>
                      <td>{{ row.stock }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="product_id" label="商品ID" min-width="90" />
          <ElTableColumn label="商品信息" min-width="220">
            <template #default="{ row }">
              <div class="goods-info">
                <ElImage
                  v-if="mediaUrl(row.image)"
                  :src="mediaUrl(row.image)"
                  fit="cover"
                  class="goods-info__img"
                />
                <div v-else class="goods-info__img is-empty" />
                <div class="goods-info__name">
                  {{ row.store_name || row.name || '—' }}
                </div>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn label="商品分类" min-width="100">
            <template #default="{ row }">
              {{ row.cate_name || '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="店铺名称" min-width="100">
            <template #default="{ row }">
              {{ row.mer_name || row.mer_id || '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="售价" min-width="90">
            <template #default="{ row }">
              {{ formatMoney(row.price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="秒杀价" min-width="90">
            <template #default="{ row }">
              {{ formatMoney(row.seckill_price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="库存" min-width="80">
            <template #default="{ row }">
              {{ row.product_stock }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="限量" min-width="80">
            <template #default="{ row }">
              {{ row.stock }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="sort" label="排序" min-width="80" />
        </ElTable>

        <div class="goods-pager">
          <ElPagination
            v-model:current-page="goodsFilter.page"
            v-model:page-size="goodsFilter.limit"
            background
            layout="total, prev, pager, next, jumper"
            :total="goodsTotal"
            @current-change="loadGoods"
            @size-change="
              () => {
                goodsFilter.page = 1;
                loadGoods();
              }
            "
          />
        </div>
      </ElTabPane>
    </ElTabs>

    <SeckillActivityProductPicker
      v-model:open="pickerOpen"
      :exclude-ids="excludeProductIds"
      :category-options="categoryOptions"
      @select="onProductsPicked"
    />
  </Drawer>
</template>

<style scoped>
.activity-tabs {
  margin-top: -4px;
}

.form-control {
  width: 100%;
  max-width: 460px;
}

.form-tip {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.5;
  color: hsl(var(--muted-foreground));
}

.activity-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.border-pic-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.border-pic-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 4px;
}

.example-tip {
  font-size: 12px;
  line-height: 1.5;
  color: hsl(var(--muted-foreground));
}

.tab-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.draft-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.goods-filter {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.goods-filter__status {
  width: 160px;
}

.goods-filter__keyword {
  width: 220px;
}

.goods-table {
  width: 100%;
}

.goods-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.goods-info__img {
  width: 44px;
  height: 44px;
  border-radius: 4px;
  flex-shrink: 0;
  background: hsl(var(--muted) / 0.35);
}

.goods-info__img.is-empty {
  border: 1px dashed hsl(var(--border));
}

.goods-info__name {
  font-size: 13px;
  line-height: 1.4;
  word-break: break-word;
}

.sku-wrap {
  padding: 8px 12px 8px 48px;
}

.sku-table {
  width: 100%;
  max-width: 720px;
  border-collapse: collapse;
  font-size: 13px;
}

.sku-table th,
.sku-table td {
  padding: 6px 10px;
  border-bottom: 1px solid hsl(var(--border));
  text-align: left;
}

.sku-table th {
  color: hsl(var(--muted-foreground));
  font-weight: 500;
  background: hsl(var(--muted) / 0.35);
}

.sku-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sku-info__img {
  width: 36px;
  height: 36px;
  border-radius: 4px;
  flex-shrink: 0;
}

.cell-input {
  width: 96px;
}

.goods-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
</style>
