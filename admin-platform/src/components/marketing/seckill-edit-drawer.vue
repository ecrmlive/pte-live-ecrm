<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { VbenTiptapPreview } from '@vben/plugins/tiptap';
import {
  ElCheckbox,
  ElCheckboxGroup,
  ElDatePicker,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRate,
  ElSelect,
  ElOption,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';

import {
  getPlatformProductEditApi,
  updatePlatformProductOpsApi,
  type PlatformProductEditDetail,
  type PlatformProductEditSKU,
} from '#/api/core/platform-catalog';
import {
  getPlatformSeckillApi,
  listPlatformSeckillTimesApi,
  updatePlatformSeckillApi,
  type PlatformSeckillActive,
  type PlatformSeckillTime,
} from '#/api/core/platform-seckill';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const emit = defineEmits<{
  saved: [];
}>();

const loading = ref(false);
const saving = ref(false);
const activeTab = ref('basic');
const seckillId = ref(0);
const seckill = ref<PlatformSeckillActive>();
const product = ref<PlatformProductEditDetail>();
const timeOptions = ref<PlatformSeckillTime[]>([]);

const form = reactive({
  name: '',
  timeIds: [] as number[],
  dateRange: [] as string[],
  seckill_price: 0,
  once_pay_count: 0,
  all_pay_count: 0,
  stock: 0,
  is_show: 1,
  star: 0,
  sort: 0,
  recommend: [] as string[],
});

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  placement: 'right',
  title: '编辑秒杀商品',
  cancelText: '取消',
  confirmText: '完成',
  onConfirm: () => void save(),
});

const productTypeLabel = computed(() => {
  const map: Record<number, string> = {
    0: '普通商品',
    1: '虚拟商品',
    2: '云盘商品',
    3: '卡密商品',
    4: '预约商品',
    5: '年/次卡',
  };
  return map[Number(product.value?.product_type ?? 0)] || '普通商品';
});

const statusLabel = computed(() => {
  const map: Record<number, string> = {
    [-2]: '下架',
    [-1]: '审核未通过',
    0: '待审核',
    1: '上架显示',
  };
  return map[Number(product.value?.status ?? 0)] || '未知';
});

const skuSpecKeys = computed(() => {
  const keys = new Set<string>();
  for (const sku of product.value?.skus || []) {
    Object.keys(sku.spec || {}).forEach((k) => keys.add(k));
  }
  return [...keys];
});

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '—';
  return String(v);
}

function yesNo(v?: number | null) {
  return Number(v) === 1 ? '是' : '否';
}

function skuSpecValue(sku: PlatformProductEditSKU, key: string) {
  return sku.spec?.[key] || '—';
}

function parseTimeIds(raw?: string) {
  if (!raw?.trim()) return [] as number[];
  return raw
    .split(/[,，\s]+/)
    .map((s) => Number(s))
    .filter((n) => Number.isFinite(n) && n > 0);
}

/** API 可能返回 RFC3339，统一成 yyyy-MM-dd 供日期控件与保存使用 */
function toDay(raw?: string) {
  const s = String(raw || '').trim();
  if (!s) return '';
  const m = s.match(/^(\d{4}-\d{2}-\d{2})/);
  return m?.[1] || '';
}

async function loadTimeOptions() {
  try {
    const page = await listPlatformSeckillTimesApi({ page: 1, limit: 200 });
    timeOptions.value = page.list || [];
  } catch {
    timeOptions.value = [];
  }
}

function resetForm() {
  Object.assign(form, {
    name: '',
    timeIds: [],
    dateRange: [],
    seckill_price: 0,
    once_pay_count: 0,
    all_pay_count: 0,
    stock: 0,
    is_show: 1,
    star: 0,
    sort: 0,
    recommend: [],
  });
}

function fillForm(sk: PlatformSeckillActive, p?: PlatformProductEditDetail) {
  form.name = sk.name || '';
  form.timeIds = parseTimeIds(sk.seckill_time_ids);
  const start = toDay(sk.start_day);
  const end = toDay(sk.end_day);
  form.dateRange = start && end ? [start, end] : [];
  form.seckill_price = Number(sk.seckill_price || 0);
  form.once_pay_count = Number(sk.once_pay_count || 0);
  form.all_pay_count = Number(sk.all_pay_count || 0);
  form.stock = Number(sk.stock || 0);
  form.is_show = Number(sk.is_show ?? 1);
  form.star = Number(sk.star || 0);
  form.sort = Number(sk.sort || 0);
  const recommend: string[] = [];
  if (p) {
    if (Number(p.is_hot) === 1) recommend.push('hot');
    if (Number(p.is_benefit) === 1) recommend.push('benefit');
    if (Number(p.is_best) === 1) recommend.push('best');
    if (Number(p.is_new) === 1) recommend.push('new');
  }
  form.recommend = recommend;
}

async function open(id: number) {
  seckillId.value = id;
  activeTab.value = 'marketing';
  seckill.value = undefined;
  product.value = undefined;
  resetForm();
  drawerApi.setState({ loading: true, title: '编辑秒杀商品' }).open();
  loading.value = true;
  try {
    await loadTimeOptions();
    const sk = await getPlatformSeckillApi(id);
    seckill.value = sk;
    if (sk.product_id) {
      try {
        product.value = await getPlatformProductEditApi(sk.product_id);
      } catch {
        product.value = undefined;
      }
    }
    fillForm(sk, product.value);
  } catch {
    ElMessage.error('加载秒杀编辑数据失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

async function save() {
  if (!seckillId.value) return;
  if (!form.name.trim()) {
    ElMessage.warning('请填写活动名称');
    activeTab.value = 'activity';
    return;
  }
  if (form.seckill_price <= 0) {
    ElMessage.warning('请填写有效秒杀价');
    activeTab.value = 'activity';
    return;
  }
  const startDay = toDay(form.dateRange?.[0]);
  const endDay = toDay(form.dateRange?.[1]);
  if (!startDay || !endDay) {
    ElMessage.warning('请选择活动日期');
    activeTab.value = 'activity';
    return;
  }

  drawerApi.lock();
  saving.value = true;
  try {
    await updatePlatformSeckillApi(seckillId.value, {
      name: form.name.trim(),
      seckill_time_ids: form.timeIds.join(','),
      start_day: startDay,
      end_day: endDay,
      seckill_price: form.seckill_price,
      once_pay_count: form.once_pay_count > 0 ? form.once_pay_count : 1,
      all_pay_count: Math.max(0, Number(form.all_pay_count || 0)),
      stock: form.stock,
      is_show: form.is_show,
      star: form.star,
      sort: form.sort,
    });

    if (product.value?.product_id) {
      const set = new Set(form.recommend);
      await updatePlatformProductOpsApi(product.value.product_id, {
        is_hot: set.has('hot') ? 1 : 0,
        is_benefit: set.has('benefit') ? 1 : 0,
        is_best: set.has('best') ? 1 : 0,
        is_new: set.has('new') ? 1 : 0,
      });
    }

    ElMessage.success('已保存');
    drawerApi.close();
    emit('saved');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    saving.value = false;
    drawerApi.unlock();
  }
}

function close() {
  drawerApi.close();
}

defineExpose({ open, close });
</script>

<template>
  <Drawer>
    <div v-loading="loading || saving" class="sk-edit">
      <template v-if="seckill">
        <div class="sk-edit__header">
          <div class="sk-edit__identity">
            <div class="sk-edit__icon">
              <IconifyIcon icon="ant-design:shopping-outlined" />
            </div>
            <div class="sk-edit__titles">
              <div class="sk-edit__sub">
                商品ID：{{ product?.product_id || seckill.product_id || '—' }}
              </div>
              <div class="sk-edit__meta">
                <div class="sk-edit__meta-item">
                  <span class="label">类型</span>
                  <span class="value">{{ productTypeLabel }}</span>
                </div>
                <div class="sk-edit__meta-item">
                  <span class="label">状态</span>
                  <span class="value">{{ statusLabel }}</span>
                </div>
                <div class="sk-edit__meta-item">
                  <span class="label">销量</span>
                  <span class="value">{{
                    dash(product?.sales ?? seckill.sales)
                  }}</span>
                </div>
                <div class="sk-edit__meta-item">
                  <span class="label">库存</span>
                  <span class="value">{{
                    dash(product?.stock ?? seckill.stock)
                  }}</span>
                </div>
                <div class="sk-edit__meta-item">
                  <span class="label">创建时间</span>
                  <span class="value">{{
                    product?.create_time
                      ? formatShanghaiDateTime(product.create_time)
                      : '—'
                  }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <ElTabs v-model="activeTab" class="sk-edit__tabs">
          <ElTabPane label="基本信息" name="basic">
            <ElForm
              label-position="left"
              label-width="96px"
              class="sk-edit__form"
            >
              <ElFormItem label="封面图">
                <ElImage
                  v-if="product?.image || seckill.image"
                  class="sk-edit__cover"
                  :src="
                    resolveCosMediaUrl(product?.image || seckill.image || '')
                  "
                  fit="cover"
                />
                <div v-else class="sk-edit__cover sk-edit__cover--empty" />
              </ElFormItem>
              <ElFormItem label="轮播图">
                <div
                  v-if="product?.slider_image?.length"
                  class="sk-edit__sliders"
                >
                  <ElImage
                    v-for="(img, idx) in product.slider_image"
                    :key="`${img}-${idx}`"
                    class="sk-edit__cover"
                    :src="resolveCosMediaUrl(img)"
                    fit="cover"
                  />
                </div>
                <span v-else>—</span>
              </ElFormItem>
              <div class="sk-edit__field-grid">
                <ElFormItem label="商品简介">
                  <span>{{ dash(product?.store_info) }}</span>
                </ElFormItem>
                <ElFormItem label="平台分类">
                  <span>{{ dash(product?.cate_name) }}</span>
                </ElFormItem>
                <ElFormItem label="商品标签">
                  <template
                    v-if="
                      product?.mer_labels?.length ||
                      product?.sys_label_names?.length
                    "
                  >
                    <ElTag
                      v-for="tag in product?.mer_labels || []"
                      :key="`m-${tag}`"
                      class="mr-1"
                      size="small"
                    >
                      {{ tag }}
                    </ElTag>
                    <ElTag
                      v-for="tag in product?.sys_label_names || []"
                      :key="`s-${tag}`"
                      class="mr-1"
                      size="small"
                      type="warning"
                    >
                      {{ tag }}
                    </ElTag>
                  </template>
                  <span v-else>—</span>
                </ElFormItem>
                <ElFormItem label="品牌选择">
                  <span>{{ dash(product?.brand_name) || '其它' }}</span>
                </ElFormItem>
                <ElFormItem label="单位">
                  <span>{{ dash(product?.unit_name) }}</span>
                </ElFormItem>
                <ElFormItem label="关键字">
                  <span>{{ dash(product?.keyword) }}</span>
                </ElFormItem>
                <ElFormItem label="配送方式">
                  <span v-if="product?.delivery_way?.length">
                    {{
                      product.delivery_way
                        .map((w) =>
                          Number(w) === 1
                            ? '到店自提'
                            : Number(w) === 2
                              ? '快递配送'
                              : String(w),
                        )
                        .join('/')
                    }}
                  </span>
                  <span v-else>—</span>
                </ElFormItem>
              </div>
            </ElForm>
          </ElTabPane>

          <ElTabPane label="活动与规格" name="activity">
            <div class="sk-edit__section-title">参与活动信息</div>
            <ElForm label-width="108px" class="sk-edit__form">
              <div class="sk-edit__field-grid">
                <ElFormItem label="活动名称" required>
                  <ElInput
                    v-model="form.name"
                    maxlength="128"
                    show-word-limit
                    placeholder="请输入活动名称"
                  />
                </ElFormItem>
                <ElFormItem label="规格">
                  <span>{{
                    Number(product?.spec_type) === 1 ? '多规格' : '单规格'
                  }}</span>
                </ElFormItem>
                <ElFormItem label="活动日期" required>
                  <ElDatePicker
                    v-model="form.dateRange"
                    type="daterange"
                    value-format="YYYY-MM-DD"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    class="!w-full"
                  />
                </ElFormItem>
                <ElFormItem label="审核状态">
                  <span>{{ dash(seckill.product_status_name) }}</span>
                </ElFormItem>
                <ElFormItem label="单次限购">
                  <ElInputNumber
                    v-model="form.once_pay_count"
                    :min="0"
                    :max="9999"
                  />
                  <span v-if="form.once_pay_count <= 0" class="hint ml-2"
                    >0 表示按默认 1</span
                  >
                </ElFormItem>
                <ElFormItem label="活动场次">
                  <ElSelect
                    v-model="form.timeIds"
                    multiple
                    collapse-tags
                    collapse-tags-tooltip
                    clearable
                    placeholder="请选择秒杀场次"
                    class="!w-full"
                  >
                    <ElOption
                      v-for="t in timeOptions"
                      :key="t.seckill_time_id"
                      :label="`${t.title}（${t.start_time}:00-${t.end_time}:00）`"
                      :value="t.seckill_time_id"
                    />
                  </ElSelect>
                </ElFormItem>
                <ElFormItem label="活动限购">
                  <ElInputNumber
                    v-model="form.all_pay_count"
                    :min="0"
                    :max="99999"
                  />
                  <span v-if="form.all_pay_count <= 0" class="hint ml-2"
                    >0 表示不限购</span
                  >
                </ElFormItem>
                <ElFormItem label="活动状态">
                  <span>{{ dash(seckill.activity_status_text) }}</span>
                </ElFormItem>
                <ElFormItem label="秒杀价" required>
                  <ElInputNumber
                    v-model="form.seckill_price"
                    :min="0.01"
                    :precision="2"
                    :step="0.1"
                  />
                </ElFormItem>
                <ElFormItem label="限量剩余">
                  <ElInputNumber v-model="form.stock" :min="0" />
                </ElFormItem>
                <ElFormItem label="是否显示">
                  <ElSwitch
                    :model-value="form.is_show === 1"
                    inline-prompt
                    active-text="显示"
                    inactive-text="隐藏"
                    @change="
                      (v: string | number | boolean) =>
                        (form.is_show = v ? 1 : 0)
                    "
                  />
                </ElFormItem>
              </div>
            </ElForm>

            <div class="sk-edit__section-title mt">规格列表</div>
            <ElTable
              :data="product?.skus || []"
              border
              empty-text="暂无数据"
            >
              <ElTableColumn
                v-for="key in skuSpecKeys"
                :key="key"
                :label="key"
                min-width="90"
              >
                <template #default="{ row }">
                  {{ skuSpecValue(row, key) }}
                </template>
              </ElTableColumn>
              <ElTableColumn label="图片" width="72" align="center">
                <template #default="{ row }">
                  <ElImage
                    v-if="row.image"
                    class="sk-edit__sku-img"
                    :src="resolveCosMediaUrl(row.image)"
                    fit="cover"
                  />
                  <span v-else>—</span>
                </template>
              </ElTableColumn>
              <ElTableColumn label="售价" width="96" align="center">
                <template #default="{ row }">
                  {{ Number(row.price || 0).toFixed(2) }}
                </template>
              </ElTableColumn>
              <ElTableColumn label="秒杀价" width="96" align="center">
                <template #default>
                  {{ Number(form.seckill_price || 0).toFixed(2) }}
                </template>
              </ElTableColumn>
              <ElTableColumn label="划线价" width="96" align="center">
                <template #default="{ row }">
                  {{ Number(row.ot_price || 0).toFixed(2) }}
                </template>
              </ElTableColumn>
              <ElTableColumn
                prop="stock"
                label="库存"
                width="80"
                align="center"
              />
            </ElTable>
          </ElTabPane>

          <ElTabPane label="商品详情" name="content">
            <div class="sk-edit__content">
              <VbenTiptapPreview
                v-if="product?.content"
                :content="product.content"
                :min-height="240"
              />
              <ElEmpty v-else description="暂无详情" :image-size="72" />
            </div>
          </ElTabPane>

          <ElTabPane label="营销信息" name="marketing">
            <ElForm label-width="120px" class="sk-edit__form">
              <ElFormItem label="商品推荐">
                <ElCheckboxGroup v-model="form.recommend">
                  <ElCheckbox label="hot">是否热卖</ElCheckbox>
                  <ElCheckbox label="benefit">促销单品</ElCheckbox>
                  <ElCheckbox label="best">是否精品</ElCheckbox>
                  <ElCheckbox label="new">是否新品</ElCheckbox>
                </ElCheckboxGroup>
              </ElFormItem>
              <ElFormItem label="平台推荐星级">
                <div class="sk-edit__control-stack">
                  <ElRate v-model="form.star" :max="5" />
                  <p class="field-tip">
                    备注：5星为最高推荐级别，1星为最低推荐级别，设置后会在商城商品列表、搜索商品列表中体现。
                  </p>
                </div>
              </ElFormItem>
              <ElFormItem label="排序">
                <ElInputNumber v-model="form.sort" :min="0" :max="999999" />
              </ElFormItem>
            </ElForm>
          </ElTabPane>

          <ElTabPane label="其它信息" name="other">
            <ElForm label-position="left" label-width="120px">
              <ElFormItem label="支持退款">
                <span>{{ yesNo(product?.refund_switch) }}</span>
              </ElFormItem>
              <ElFormItem label="店铺商品参数">
                <ElTable
                  :data="product?.mer_params || []"
                  border
                  size="small"
                  empty-text="暂无数据"
                  class="sk-edit__param-table"
                >
                  <ElTableColumn
                    prop="name"
                    label="参数名称"
                    min-width="120"
                  />
                  <ElTableColumn
                    prop="value"
                    label="参数值"
                    min-width="160"
                  />
                </ElTable>
              </ElFormItem>
              <ElFormItem label="平台商品参数">
                <ElTable
                  :data="product?.platform_params || []"
                  border
                  size="small"
                  empty-text="暂无数据"
                  class="sk-edit__param-table"
                >
                  <ElTableColumn
                    prop="name"
                    label="参数名称"
                    min-width="120"
                  />
                  <ElTableColumn
                    prop="value"
                    label="参数值"
                    min-width="160"
                  />
                </ElTable>
              </ElFormItem>
            </ElForm>
          </ElTabPane>
        </ElTabs>
      </template>
    </div>
  </Drawer>
</template>

<style scoped>
.sk-edit {
  min-height: 360px;
}

.sk-edit__header {
  padding-bottom: 12px;
  margin-bottom: 4px;
  border-bottom: 1px solid hsl(var(--border));
}

.sk-edit__identity {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.sk-edit__icon {
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

.sk-edit__sub {
  margin-bottom: 10px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.sk-edit__meta {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px 16px;
}

.sk-edit__meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sk-edit__meta-item .label {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.sk-edit__meta-item .value {
  font-size: 14px;
  line-height: 22px;
  word-break: break-all;
}

.sk-edit__tabs {
  min-height: 420px;
}

.sk-edit__section-title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  font-size: 14px;
  font-weight: 600;
}

.sk-edit__section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 2px;
  background: hsl(var(--primary));
}

.sk-edit__section-title.mt {
  margin-top: 24px;
}

.sk-edit__field-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 12px;
}

.sk-edit__cover {
  width: 64px;
  height: 64px;
  border-radius: 4px;
  overflow: hidden;
}

.sk-edit__cover--empty {
  border: 1px dashed hsl(var(--border));
  background: hsl(var(--muted) / 0.35);
}

.sk-edit__sliders {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.sk-edit__sku-img {
  width: 40px;
  height: 40px;
  border-radius: 4px;
}

.sk-edit__content {
  min-height: 240px;
  padding: 8px 0;
}

.sk-edit__param-table {
  width: 100%;
  max-width: 640px;
}

.sk-edit__control-stack {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-tip {
  margin: 0;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
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

@media (max-width: 960px) {
  .sk-edit__meta,
  .sk-edit__field-grid {
    grid-template-columns: 1fr;
  }
}
</style>
