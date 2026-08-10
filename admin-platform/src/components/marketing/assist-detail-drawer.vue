<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElMessage,
  ElSkeleton,
  ElTag,
} from 'element-plus';

import type { PlatformAssistActive } from '#/api/core/platform-assist';
import ImageField from '#/components/shop/image-field.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';

import { loadAssistProductBundle } from './load-assist-product';

const loading = ref(false);
const productMissing = ref(false);
const detail = ref<PlatformAssistActive>();
const cover = ref('');

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '助力商品详情',
});

const activityDateText = computed(() => {
  const row = detail.value;
  if (!row) return '—';
  return `${formatShanghaiDateTime(row.start_time)} – ${formatShanghaiDateTime(row.end_time)}`;
});

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '—';
  return String(v);
}

function payCountText(n?: number | null) {
  const v = Number(n || 0);
  return v > 0 ? String(v) : '不限';
}

function showStatusText(isShow?: number | null) {
  return Number(isShow) === 1 ? '显示' : '隐藏';
}

async function open(id: number) {
  productMissing.value = false;
  detail.value = undefined;
  cover.value = '';
  drawerApi.setState({ loading: true, title: '助力商品详情' }).open();
  loading.value = true;
  try {
    const bundle = await loadAssistProductBundle(id);
    detail.value = bundle.assist;
    productMissing.value = bundle.productMissing;
    cover.value =
      bundle.product?.image ||
      bundle.assist.image ||
      '';
  } catch {
    ElMessage.error('加载助力商品详情失败');
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
    <div v-loading="loading" class="assist-detail">
      <ElSkeleton :loading="loading && !detail" animated :rows="12">
        <template #default>
          <template v-if="detail">
            <section class="assist-detail__section">
              <div class="assist-detail__section-title">基本信息</div>
              <div class="assist-detail__grid">
                <div class="assist-detail__item">
                  <span class="label">店铺名称</span>
                  <span class="value">{{
                    detail.mer_name || `店铺#${detail.mer_id}`
                  }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">店铺类别</span>
                  <span class="value">{{
                    detail.trader_name ||
                    (detail.is_trader === 1 ? '自营' : '非自营')
                  }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">商品ID</span>
                  <span class="value">{{ dash(detail.product_id) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">商品名称</span>
                  <span class="value">{{
                    detail.store_name || `商品 #${detail.product_id}`
                  }}</span>
                </div>
                <div class="assist-detail__item assist-detail__item--image">
                  <span class="label">商品图</span>
                  <span class="value">
                    <ImageField v-model="cover" disabled />
                  </span>
                </div>
              </div>
            </section>

            <section class="assist-detail__section">
              <div class="assist-detail__section-title">助力商品活动信息</div>
              <div class="assist-detail__grid">
                <div class="assist-detail__item assist-detail__item--span2">
                  <span class="label">助力活动简介</span>
                  <span class="value">{{ dash(detail.store_info) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">助力价</span>
                  <span class="value"
                    >¥{{ Number(detail.assist_price || 0).toFixed(2) }}</span
                  >
                </div>
                <div class="assist-detail__item">
                  <span class="label">限量</span>
                  <span class="value">{{
                    dash(detail.stock_count ?? detail.stock)
                  }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">助力人数</span>
                  <span class="value">{{ dash(detail.assist_count) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">限量剩余</span>
                  <span class="value">{{ dash(detail.stock) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">限购件数</span>
                  <span class="value">{{ payCountText(detail.pay_count) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">助力次数</span>
                  <span class="value">{{
                    dash(detail.assist_user_count)
                  }}</span>
                </div>
                <div class="assist-detail__item assist-detail__item--span2">
                  <span class="label">助力活动日期</span>
                  <span class="value">{{ activityDateText }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">已售商品数</span>
                  <span class="value">{{ dash(detail.pay ?? 0) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">助力成功/参与人次</span>
                  <span class="value"
                    >{{ detail.success ?? 0 }} / {{ detail.all ?? 0 }}</span
                  >
                </div>
                <div class="assist-detail__item">
                  <span class="label">审核状态</span>
                  <span class="value">
                    <ElTag size="small">{{
                      detail.product_status_name || '—'
                    }}</ElTag>
                  </span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">助力活动状态</span>
                  <span class="value">
                    <ElTag
                      size="small"
                      :type="
                        detail.assist_status === 1
                          ? 'success'
                          : detail.assist_status === 0
                            ? 'info'
                            : 'warning'
                      "
                    >
                      {{ detail.assist_status_text || '—' }}
                    </ElTag>
                  </span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">显示状态</span>
                  <span class="value">{{ showStatusText(detail.is_show) }}</span>
                </div>
                <div class="assist-detail__item">
                  <span class="label">创建时间</span>
                  <span class="value">{{
                    detail.create_time
                      ? formatShanghaiDateTime(detail.create_time)
                      : '—'
                  }}</span>
                </div>
              </div>
            </section>

            <div
              v-if="productMissing"
              class="assist-detail__hint"
            >
              关联商品编辑信息暂不可用，以上活动字段仍来自助力活动接口。
            </div>

            <div class="assist-detail__footer">
              <ElButton @click="close">关闭</ElButton>
            </div>
          </template>
        </template>
      </ElSkeleton>
    </div>
  </Drawer>
</template>

<style scoped>
.assist-detail {
  display: flex;
  flex-direction: column;
  min-height: 420px;
  padding: 0 4px 8px;
}

.assist-detail__section + .assist-detail__section {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px dashed hsl(var(--border));
}

.assist-detail__section-title {
  position: relative;
  margin-bottom: 16px;
  padding-left: 10px;
  color: hsl(var(--foreground));
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
}

.assist-detail__section-title::before {
  position: absolute;
  top: 2px;
  left: 0;
  width: 3px;
  height: 14px;
  background: hsl(var(--primary));
  border-radius: 2px;
  content: '';
}

.assist-detail__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px 32px;
}

.assist-detail__item {
  display: grid;
  grid-template-columns: 120px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
  font-size: 13px;
  line-height: 1.6;
}

.assist-detail__item--span2 {
  grid-column: 1 / -1;
}

.assist-detail__item--image {
  grid-column: 1 / -1;
  align-items: center;
}

.assist-detail__item .label {
  color: hsl(var(--muted-foreground));
}

.assist-detail__item .value {
  color: hsl(var(--foreground));
  word-break: break-all;
}

.assist-detail__hint {
  margin-top: 12px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
}

.assist-detail__footer {
  position: sticky;
  bottom: 0;
  z-index: 2;
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding: 12px 4px 0;
  border-top: 1px solid hsl(var(--border));
  background: hsl(var(--background));
}
</style>
