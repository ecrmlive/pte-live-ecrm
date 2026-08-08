<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { VbenTiptapPreview } from '@vben/plugins/tiptap';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElButton,
  ElDatePicker,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElImage,
  ElMessage,
  ElOption,
  ElPagination,
  ElRate,
  ElSelect,
  ElSkeleton,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';

import {
  getPlatformProductEditApi,
  listPlatformProductOperateLogsApi,
  type PlatformProductEditDetail,
  type PlatformProductEditSKU,
  type PlatformProductOperateLog,
} from '#/api/core/platform-catalog';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';

import ProductFictiModal from './ProductFictiModal.vue';

const emit = defineEmits<{
  edit: [productId: number];
  fictiUpdated: [];
}>();

const fictiModalRef = ref<InstanceType<typeof ProductFictiModal>>();

const loading = ref(false);
const activeTab = ref('basic');
const detail = ref<PlatformProductEditDetail>();

const operateLogs = ref<PlatformProductOperateLog[]>([]);
const operateLogTotal = ref(0);
const operateLogLoading = ref(false);
const operateLogPage = ref(1);
const operateLogLimit = ref(10);
const operateLogTerminal = ref('');
const operateLogDates = ref<string[]>([]);

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

const skuSpecKeys = computed(() => {
  const keys = new Set<string>();
  for (const sku of detail.value?.skus || []) {
    Object.keys(sku.spec || {}).forEach((k) => keys.add(k));
  }
  return [...keys];
});

const platformRecommendTags = computed(() => {
  const d = detail.value;
  if (!d) return [] as string[];
  const tags: string[] = [];
  if (d.is_hot) tags.push('热门榜单');
  if (d.is_benefit) tags.push('促销单品');
  if (d.is_best) tags.push('精品推荐');
  if (d.is_new) tags.push('首发新品');
  return tags;
});

const svipPriceText = computed(() => {
  const d = detail.value;
  if (!d) return '不设置';
  if (Number(d.svip_price_type) === 1) return '默认会员价';
  if (Number(d.svip_price_type) === 2) {
    return `¥${Number(d.svip_price || 0).toFixed(2)}`;
  }
  return '不设置';
});

const onceMinText = computed(() => {
  const n = Number(detail.value?.once_min_count || 0);
  return n > 0 ? String(n) : '不限购';
});

const formLinkedText = computed(() => {
  const id = Number(detail.value?.mer_form_id || 0);
  return id > 0 ? `#${id}` : '关闭';
});

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
  title: '商品详情',
});

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '—';
  return String(v);
}

function onOff(v?: number | null) {
  return Number(v) === 1 ? '开启' : '关闭';
}

function yesNo(v?: number | null) {
  return Number(v) === 1 ? '是' : '否';
}

async function open(productId: number) {
  activeTab.value = 'basic';
  detail.value = undefined;
  drawerApi.setState({ loading: true }).open();
  loading.value = true;
  try {
    detail.value = await getPlatformProductEditApi(productId);
    resetOperateLogQuery();
    void loadOperateLogs();
  } catch {
    ElMessage.error('加载商品详情失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

async function reload() {
  const id = Number(detail.value?.product_id || 0);
  if (!id) return;
  loading.value = true;
  try {
    detail.value = await getPlatformProductEditApi(id);
  } catch {
    ElMessage.error('刷新商品详情失败');
  } finally {
    loading.value = false;
  }
}

function close() {
  drawerApi.close();
}

function openFicti() {
  if (!detail.value) return;
  fictiModalRef.value?.open({
    productId: detail.value.product_id,
    ficti: Number(detail.value.ficti || 0),
  });
}

async function onFictiSuccess(payload: { ficti: number; productId: number }) {
  if (detail.value && detail.value.product_id === payload.productId) {
    detail.value = { ...detail.value, ficti: payload.ficti };
  }
  await reload();
  emit('fictiUpdated');
}

function onEdit() {
  if (!detail.value) return;
  const id = detail.value.product_id;
  close();
  emit('edit', id);
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

function skuSpecValue(sku: PlatformProductEditSKU, key: string) {
  return sku.spec?.[key] || '—';
}

defineExpose({ open, close, reload });
</script>

<template>
  <Drawer>
    <ElSkeleton :loading="loading" animated :rows="10">
      <template #default>
        <div v-if="detail" class="product-detail">
          <div class="product-detail__header">
            <div class="product-detail__brand-row">
              <div class="product-detail__identity">
                <div class="product-detail__icon">
                  <IconifyIcon icon="ant-design:shopping-outlined" />
                </div>
                <div class="product-detail__titles">
                  <div class="product-detail__name">
                    {{ detail.title || '未命名商品' }}
                  </div>
                  <div class="product-detail__sub">
                    商品ID：{{ detail.product_id }}
                  </div>
                </div>
              </div>
              <ElButton type="primary" @click="onEdit">编辑</ElButton>
            </div>
            <div class="product-detail__meta">
              <div class="product-detail__meta-item">
                <span class="label">类型</span>
                <span class="value">{{ productTypeLabel }}</span>
              </div>
              <div class="product-detail__meta-item">
                <span class="label">状态</span>
                <span class="value">{{ statusLabel }}</span>
              </div>
              <div class="product-detail__meta-item">
                <span class="label">销量</span>
                <span class="value">{{ detail.sales }}</span>
              </div>
              <div class="product-detail__meta-item">
                <span class="label">库存</span>
                <span class="value">{{ detail.stock }}</span>
              </div>
              <div class="product-detail__meta-item">
                <span class="label">创建时间</span>
                <span class="value">{{ formatShanghaiDateTime(detail.create_time) }}</span>
              </div>
            </div>
          </div>

          <ElTabs v-model="activeTab" class="product-detail__tabs">
            <ElTabPane label="基本信息" name="basic">
              <div class="product-detail__section-title">基本信息</div>
              <ElForm label-position="left" label-width="96px" class="product-detail__form">
                <ElFormItem label="封面图">
                  <ElImage
                    v-if="detail.image"
                    class="product-detail__cover"
                    :src="resolveCosMediaUrl(detail.image)"
                    fit="cover"
                    :preview-src-list="[resolveCosMediaUrl(detail.image)]"
                  />
                  <span v-else>—</span>
                </ElFormItem>
                <ElFormItem label="轮播图">
                  <div v-if="detail.slider_image?.length" class="product-detail__sliders">
                    <ElImage
                      v-for="(img, idx) in detail.slider_image"
                      :key="`${img}-${idx}`"
                      class="product-detail__cover"
                      :src="resolveCosMediaUrl(img)"
                      fit="cover"
                      :preview-src-list="detail.slider_image.map((i) => resolveCosMediaUrl(i))"
                      :initial-index="idx"
                    />
                  </div>
                  <span v-else>—</span>
                </ElFormItem>
                <div class="product-detail__field-grid">
                  <ElFormItem label="商品简介">
                    <span>{{ dash(detail.store_info) }}</span>
                  </ElFormItem>
                  <ElFormItem label="平台分类">
                    <span>{{ dash(detail.cate_name) }}</span>
                  </ElFormItem>
                  <ElFormItem label="店铺分类">
                    <span>{{ dash(detail.mer_cate_name) }}</span>
                  </ElFormItem>
                  <ElFormItem label="商品标签">
                    <template v-if="detail.mer_labels?.length || detail.sys_label_names?.length">
                      <ElTag
                        v-for="tag in detail.mer_labels || []"
                        :key="`m-${tag}`"
                        class="mr-1"
                        size="small"
                      >
                        {{ tag }}
                      </ElTag>
                      <ElTag
                        v-for="tag in detail.sys_label_names || []"
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
                    <span>{{ dash(detail.brand_name) }}</span>
                  </ElFormItem>
                  <ElFormItem label="单位">
                    <span>{{ dash(detail.unit_name) }}</span>
                  </ElFormItem>
                  <ElFormItem label="关键字">
                    <span>{{ dash(detail.keyword) }}</span>
                  </ElFormItem>
                  <ElFormItem label="配送方式">
                    <span v-if="detail.delivery_way?.length">
                      {{ detail.delivery_way.join('/') }}
                    </span>
                    <span v-else>—</span>
                  </ElFormItem>
                </div>
              </ElForm>
            </ElTabPane>

            <ElTabPane label="规格与价格" name="sku">
              <div class="product-detail__section-title">规格信息</div>
              <div class="product-detail__section-head">
                <span>规格：</span>
                <span>{{ detail.spec_type === 1 ? '多规格' : '单规格' }}</span>
              </div>
              <div class="product-detail__section-title">规格列表</div>
              <ElTable :data="detail.skus || []" border>
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
                      class="product-detail__sku-img"
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
                <ElTableColumn label="划线价" width="96" align="center">
                  <template #default="{ row }">
                    {{ Number(row.ot_price || 0).toFixed(2) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn prop="stock" label="库存" width="80" align="center" />
                <ElTableColumn label="规格编码" min-width="100" align="center">
                  <template #default="{ row }">{{ dash(row.code) }}</template>
                </ElTableColumn>
                <ElTableColumn label="条形码" min-width="100" align="center">
                  <template #default="{ row }">{{ dash(row.bar_code) }}</template>
                </ElTableColumn>
                <ElTableColumn label="重量(KG)" width="96" align="center">
                  <template #default="{ row }">
                    {{ Number(row.weight || 0).toFixed(2) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn label="体积(m³)" width="96" align="center">
                  <template #default="{ row }">
                    {{ Number(row.volume || 0).toFixed(2) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn label="一级返佣" width="96" align="center">
                  <template #default="{ row }">
                    {{ Number(row.extension_one || 0).toFixed(2) }}
                  </template>
                </ElTableColumn>
              </ElTable>
            </ElTabPane>

            <ElTabPane label="商品详情" name="content">
              <div class="product-detail__section-title">商品详情</div>
              <div class="product-detail__content">
                <VbenTiptapPreview
                  v-if="detail.content"
                  :content="detail.content"
                  :min-height="240"
                />
                <ElEmpty v-else description="暂无详情" :image-size="72" />
              </div>
            </ElTabPane>

            <ElTabPane label="营销信息" name="marketing">
              <div class="product-detail__section-title">营销信息</div>
              <div class="product-detail__kv-grid">
                <div class="product-detail__kv">
                  <span class="label">店铺推荐</span>
                  <span class="value">{{ onOff(detail.mer_recommend) }}</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">平台推荐</span>
                  <span class="value">
                    <template v-if="platformRecommendTags.length">
                      {{ platformRecommendTags.join('/ ') }}
                    </template>
                    <template v-else>—</template>
                  </span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">分销礼包</span>
                  <span class="value">{{ yesNo(detail.is_gift_bag) }}</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">平台推荐星级</span>
                  <span class="value">
                    <ElRate :model-value="Number(detail.star || 0)" disabled :max="5" />
                  </span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">收藏人数</span>
                  <span class="value">{{ Number(detail.care_count || 0) }}人</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">已售数量</span>
                  <span class="value">
                    <ElButton link type="primary" @click="openFicti">
                      {{ Number(detail.ficti || 0) }}
                    </ElButton>
                    <span class="hint">指手动添加数量</span>
                  </span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">实际销量</span>
                  <span class="value">
                    {{ detail.sales }}
                    <span class="hint">指实际售出数量</span>
                  </span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">佣金设置</span>
                  <span class="value">
                    {{ dash(detail.commission_text) || '默认' }}
                  </span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">付费会员价</span>
                  <span class="value">{{ svipPriceText }}</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">大图推荐</span>
                  <span class="value">{{ onOff(detail.cate_hot) }}</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">活动标签</span>
                  <span class="value">
                    <template v-if="detail.activity_labels?.length">
                      <ElTag
                        v-for="tag in detail.activity_labels"
                        :key="tag"
                        class="mr-1 activity-tag"
                        size="small"
                        effect="plain"
                        type="danger"
                      >
                        {{ tag }}
                      </ElTag>
                    </template>
                    <template v-else>—</template>
                  </span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">排序</span>
                  <span class="value">{{ detail.rank || 0 }}</span>
                </div>
              </div>
            </ElTabPane>

            <ElTabPane label="其他信息" name="other">
              <div class="product-detail__section-title">其他信息</div>
              <ElForm label-position="left" label-width="120px" class="product-detail__form">
                <ElFormItem label="支持退款">
                  <span>{{ yesNo(detail.refund_switch) }}</span>
                </ElFormItem>
                <ElFormItem label="最少购买件数">
                  <span>{{ onceMinText }}</span>
                </ElFormItem>
                <ElFormItem label="店铺商品参数">
                  <ElTable
                    :data="detail.mer_params || []"
                    border
                    size="small"
                    empty-text="暂无数据"
                    class="product-detail__param-table"
                  >
                    <ElTableColumn prop="name" label="参数名称" min-width="120" />
                    <ElTableColumn prop="value" label="参数值" min-width="160" />
                  </ElTable>
                </ElFormItem>
                <ElFormItem label="平台商品参数">
                  <ElTable
                    :data="detail.platform_params || []"
                    border
                    size="small"
                    empty-text="暂无数据"
                    class="product-detail__param-table"
                  >
                    <ElTableColumn prop="name" label="参数名称" min-width="120" />
                    <ElTableColumn prop="value" label="参数值" min-width="160" />
                  </ElTable>
                </ElFormItem>
                <ElFormItem label="关联系统表单">
                  <span>{{ formLinkedText }}</span>
                </ElFormItem>
              </ElForm>
            </ElTabPane>

            <ElTabPane label="店铺信息" name="store">
              <div class="product-detail__section-title">店铺信息</div>
              <div class="product-detail__store-grid">
                <div class="product-detail__kv">
                  <span class="label">店铺名称</span>
                  <span class="value">{{ dash(detail.store_name || detail.mer_name) }}</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">店铺类别</span>
                  <span class="value">{{ dash(detail.mer_category_name) }}</span>
                </div>
                <div class="product-detail__kv">
                  <span class="label">店铺类型</span>
                  <span class="value">{{ dash(detail.mer_type_name) }}</span>
                </div>
              </div>
            </ElTabPane>

            <ElTabPane label="操作记录" name="logs">
              <div class="product-detail__log-filters">
                <div class="product-detail__log-filter">
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
                <div class="product-detail__log-filter">
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
              <ElTable v-loading="operateLogLoading" :data="operateLogs" border>
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
              <div class="product-detail__log-pager">
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
  </Drawer>
  <ProductFictiModal ref="fictiModalRef" @success="onFictiSuccess" />
</template>

<style scoped>
.product-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-detail__header {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid hsl(var(--border));
}

.product-detail__brand-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  justify-content: space-between;
}

.product-detail__identity {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  min-width: 0;
}

.product-detail__icon {
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

.product-detail__name {
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
}

.product-detail__sub {
  margin-top: 2px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.product-detail__meta {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px 16px;
}

.product-detail__meta-item,
.product-detail__kv {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.product-detail__meta-item .label,
.product-detail__kv .label {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.product-detail__meta-item .value,
.product-detail__kv .value {
  font-size: 14px;
  line-height: 22px;
}

.product-detail__kv-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px 24px;
}

.product-detail__store-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px 24px;
}

.product-detail__tabs {
  min-height: 420px;
}

.product-detail__section-title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  color: hsl(var(--foreground));
  font-size: 14px;
  font-weight: 600;
}

.product-detail__section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 2px;
  background: hsl(var(--primary));
}

.product-detail__form {
  max-width: 100%;
}

.product-detail__field-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 12px;
}

.product-detail__field-grid :deep(.el-form-item) {
  margin-bottom: 12px;
}

.product-detail__cover,
.product-detail__sku-img {
  width: 64px;
  height: 64px;
  border-radius: 4px;
}

.product-detail__sliders {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.product-detail__section-head {
  display: flex;
  gap: 4px;
  align-items: center;
  margin-bottom: 16px;
  font-size: 14px;
}

.product-detail__content {
  min-height: 240px;
  padding: 8px 0;
}

.product-detail__param-table {
  width: 100%;
  max-width: 640px;
}

.product-detail__log-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 16px 24px;
  margin-bottom: 12px;
}

.product-detail__log-filter {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-label {
  color: hsl(var(--muted-foreground));
  font-size: 13px;
  white-space: nowrap;
}

.product-detail__log-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.hint {
  margin-left: 6px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.activity-tag {
  margin-bottom: 2px;
}

.mr-1 {
  margin-right: 4px;
}

@media (max-width: 960px) {
  .product-detail__meta,
  .product-detail__kv-grid,
  .product-detail__store-grid,
  .product-detail__field-grid {
    grid-template-columns: 1fr;
  }
}
</style>
