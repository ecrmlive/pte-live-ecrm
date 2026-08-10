<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { VbenTiptapPreview } from '@vben/plugins/tiptap';
import {
  ElEmpty,
  ElForm,
  ElFormItem,
  ElImage,
  ElMessage,
  ElRate,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';

import {
  getPlatformProductEditApi,
  type PlatformProductEditDetail,
  type PlatformProductEditSKU,
} from '#/api/core/platform-catalog';
import {
  getPlatformSeckillApi,
  type PlatformSeckillActive,
} from '#/api/core/platform-seckill';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const loading = ref(false);
const activeTab = ref('basic');
const seckill = ref<PlatformSeckillActive>();
const product = ref<PlatformProductEditDetail>();

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
  title: '秒杀商品详情',
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

const platformRecommendText = computed(() => {
  const d = product.value;
  if (!d) return '无';
  const tags: string[] = [];
  if (d.is_hot) tags.push('热门榜单');
  if (d.is_benefit) tags.push('促销单品');
  if (d.is_best) tags.push('精品推荐');
  if (d.is_new) tags.push('首发新品');
  return tags.length ? tags.join('/') : '无';
});

const limitText = (n?: number | null) => {
  const v = Number(n || 0);
  return v > 0 ? String(v) : '不限购';
};

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

async function open(seckillActiveId: number) {
  activeTab.value = 'basic';
  seckill.value = undefined;
  product.value = undefined;
  drawerApi.setState({ loading: true, title: '秒杀商品详情' }).open();
  loading.value = true;
  try {
    const sk = await getPlatformSeckillApi(seckillActiveId);
    seckill.value = sk;
    if (sk.product_id) {
      try {
        product.value = await getPlatformProductEditApi(sk.product_id);
      } catch {
        product.value = undefined;
      }
    }
  } catch {
    ElMessage.error('加载秒杀详情失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

function close() {
  drawerApi.close();
}

defineExpose({ open, close });
</script>

<template>
  <Drawer>
    <div v-loading="loading" class="sk-detail">
      <template v-if="seckill">
        <div class="sk-detail__header">
          <div class="sk-detail__identity">
            <div class="sk-detail__icon">
              <IconifyIcon icon="ant-design:shopping-outlined" />
            </div>
            <div class="sk-detail__titles">
              <div class="sk-detail__sub">
                商品ID：{{ product?.product_id || seckill.product_id || '—' }}
              </div>
              <div class="sk-detail__meta">
                <div class="sk-detail__meta-item">
                  <span class="label">类型</span>
                  <span class="value">{{ productTypeLabel }}</span>
                </div>
                <div class="sk-detail__meta-item">
                  <span class="label">状态</span>
                  <span class="value">{{ statusLabel }}</span>
                </div>
                <div class="sk-detail__meta-item">
                  <span class="label">销量</span>
                  <span class="value">{{
                    dash(product?.sales ?? seckill.sales)
                  }}</span>
                </div>
                <div class="sk-detail__meta-item">
                  <span class="label">库存</span>
                  <span class="value">{{
                    dash(product?.stock ?? seckill.stock)
                  }}</span>
                </div>
                <div class="sk-detail__meta-item">
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

        <ElTabs v-model="activeTab" class="sk-detail__tabs">
          <ElTabPane label="基本信息" name="basic">
            <ElForm
              label-position="left"
              label-width="96px"
              class="sk-detail__form"
            >
              <ElFormItem label="封面图">
                <ElImage
                  v-if="product?.image || seckill.image"
                  class="sk-detail__cover"
                  :src="
                    resolveCosMediaUrl(product?.image || seckill.image || '')
                  "
                  fit="cover"
                />
                <div v-else class="sk-detail__cover sk-detail__cover--empty" />
              </ElFormItem>
              <ElFormItem label="轮播图">
                <div
                  v-if="product?.slider_image?.length"
                  class="sk-detail__sliders"
                >
                  <ElImage
                    v-for="(img, idx) in product.slider_image"
                    :key="`${img}-${idx}`"
                    class="sk-detail__cover"
                    :src="resolveCosMediaUrl(img)"
                    fit="cover"
                  />
                </div>
                <span v-else>—</span>
              </ElFormItem>
              <div class="sk-detail__field-grid">
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
            <div class="sk-detail__section-title">参与活动信息</div>
            <div class="sk-detail__kv-grid">
              <div class="sk-detail__kv">
                <span class="label">活动名称</span>
                <span class="value">{{ dash(seckill.name) }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">规格</span>
                <span class="value">{{
                  Number(product?.spec_type) === 1 ? '多规格' : '单规格'
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">活动日期</span>
                <span class="value">
                  <template v-if="seckill.start_day && seckill.end_day">
                    {{ seckill.start_day }} - {{ seckill.end_day }}
                  </template>
                  <template v-else>—</template>
                </span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">审核状态</span>
                <span class="value">{{
                  dash(seckill.product_status_name)
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">单次限购</span>
                <span class="value">{{
                  limitText(seckill.once_pay_count)
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">活动场次</span>
                <span class="value">{{ dash(seckill.time_titles) }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">活动限购</span>
                <span class="value">{{
                  limitText(seckill.all_pay_count)
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">活动状态</span>
                <span class="value">{{
                  dash(seckill.activity_status_text)
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">秒杀价</span>
                <span class="value"
                  >¥{{ Number(seckill.seckill_price || 0).toFixed(2) }}</span
                >
              </div>
            </div>

            <div class="sk-detail__section-title mt">规格列表</div>
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
                    class="sk-detail__sku-img"
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
                  {{ Number(seckill.seckill_price || 0).toFixed(2) }}
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
            <div class="sk-detail__content">
              <VbenTiptapPreview
                v-if="product?.content"
                :content="product.content"
                :min-height="240"
              />
              <ElEmpty v-else description="暂无详情" :image-size="72" />
            </div>
          </ElTabPane>

          <ElTabPane label="营销信息" name="marketing">
            <div class="sk-detail__kv-grid">
              <div class="sk-detail__kv">
                <span class="label">店铺推荐</span>
                <span class="value">{{
                  Number(product?.mer_recommend) === 1 ? '是' : '否'
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">平台推荐</span>
                <span class="value">{{ platformRecommendText }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">收藏人数</span>
                <span class="value"
                  >{{ Number(product?.care_count || 0) }} 人</span
                >
              </div>
              <div class="sk-detail__kv">
                <span class="label">已售数量</span>
                <span class="value">
                  {{ Number(product?.ficti || seckill.sales || 0) }}
                  <span class="hint">指手动添加数量</span>
                </span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">实际销量</span>
                <span class="value">
                  {{ Number(product?.sales || 0) }}
                  <span class="hint">指实际售出数量</span>
                </span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">平台排序</span>
                <span class="value">{{
                  dash(seckill.sort ?? product?.rank)
                }}</span>
              </div>
              <div class="sk-detail__kv">
                <span class="label">推荐星级</span>
                <span class="value">
                  <ElRate
                    :model-value="Number(seckill.star || product?.star || 0)"
                    disabled
                    :max="5"
                  />
                </span>
              </div>
            </div>
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
                  class="sk-detail__param-table"
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
                  class="sk-detail__param-table"
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
.sk-detail {
  min-height: 360px;
}

.sk-detail__header {
  padding-bottom: 12px;
  margin-bottom: 4px;
  border-bottom: 1px solid hsl(var(--border));
}

.sk-detail__identity {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.sk-detail__icon {
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

.sk-detail__sub {
  margin-bottom: 10px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.sk-detail__meta {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px 16px;
}

.sk-detail__meta-item,
.sk-detail__kv {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sk-detail__meta-item .label,
.sk-detail__kv .label {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.sk-detail__meta-item .value,
.sk-detail__kv .value {
  font-size: 14px;
  line-height: 22px;
  word-break: break-all;
}

.sk-detail__tabs {
  min-height: 420px;
}

.sk-detail__section-title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  font-size: 14px;
  font-weight: 600;
}

.sk-detail__section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 2px;
  background: hsl(var(--primary));
}

.sk-detail__section-title.mt {
  margin-top: 24px;
}

.sk-detail__kv-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px 24px;
}

.sk-detail__field-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 12px;
}

.sk-detail__cover {
  width: 64px;
  height: 64px;
  border-radius: 4px;
  overflow: hidden;
}

.sk-detail__cover--empty {
  border: 1px dashed hsl(var(--border));
  background: hsl(var(--muted) / 0.35);
}

.sk-detail__sliders {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.sk-detail__sku-img {
  width: 40px;
  height: 40px;
  border-radius: 4px;
}

.sk-detail__content {
  min-height: 240px;
  padding: 8px 0;
}

.sk-detail__param-table {
  width: 100%;
  max-width: 640px;
}

.hint {
  margin-left: 6px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.mr-1 {
  margin-right: 4px;
}

@media (max-width: 960px) {
  .sk-detail__meta,
  .sk-detail__kv-grid,
  .sk-detail__field-grid {
    grid-template-columns: 1fr;
  }
}
</style>
