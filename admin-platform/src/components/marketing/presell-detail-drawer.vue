<script setup lang="ts">
import { ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElMessage,
  ElRate,
  ElSkeleton,
  ElTag,
} from 'element-plus';

import type { PlatformProductEditDetail } from '#/api/core/platform-catalog';
import type { PlatformPresell } from '#/api/core/platform-presell';
import { formatShanghaiDateTime } from '#/utils/date-time';

import { loadPresellProductBundle } from './load-presell-product';
import PresellProductTabs from './presell-product-tabs.vue';

const loading = ref(false);
const activeTab = ref('basic');
const productMissing = ref(false);
const detail = ref<PlatformPresell>();
const product = ref<PlatformProductEditDetail>();

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '预售商品详情',
});

async function open(id: number) {
  activeTab.value = 'basic';
  productMissing.value = false;
  detail.value = undefined;
  product.value = undefined;
  drawerApi.setState({ loading: true, title: '预售商品详情' }).open();
  loading.value = true;
  try {
    const bundle = await loadPresellProductBundle(id);
    detail.value = bundle.presell;
    product.value = bundle.product;
    productMissing.value = bundle.productMissing;
  } catch {
    ElMessage.error('加载预售详情失败');
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
    <div v-loading="loading" class="presell-detail">
      <ElSkeleton :loading="loading && !detail" animated :rows="10">
        <template #default>
          <template v-if="detail">
            <PresellProductTabs
              v-model="activeTab"
              show-presell-tab
              :presell="detail"
              :product="product"
              :product-missing="productMissing"
              product-missing-tip="关联商品编辑信息暂不可用，仍可查看预售活动信息。"
            >
              <template #presell>
                <ElForm label-position="left" label-width="120px" class="pr-2">
                  <ElFormItem label="预售 ID">
                    <span>{{ detail.product_presell_id }}</span>
                  </ElFormItem>
                  <ElFormItem label="关联商品 ID">
                    <span>{{ detail.product_id || '—' }}</span>
                  </ElFormItem>
                  <ElFormItem label="店铺">
                    <span>{{
                      detail.mer_name || `店铺#${detail.mer_id}`
                    }}</span>
                  </ElFormItem>
                  <ElFormItem label="店铺类别">
                    <span>{{
                      detail.trader_name ||
                      (detail.is_trader === 1 ? '自营' : '非自营')
                    }}</span>
                  </ElFormItem>
                  <ElFormItem label="预售类型">
                    <span>{{
                      detail.presell_type === 2 ? '定金预售' : '全款预售'
                    }}</span>
                  </ElFormItem>
                  <ElFormItem label="预售标题">
                    <span>{{ detail.store_name || '—' }}</span>
                  </ElFormItem>
                  <ElFormItem label="预售价">
                    <span>¥{{ Number(detail.price || 0).toFixed(2) }}</span>
                  </ElFormItem>
                  <ElFormItem label="定金 / 尾款">
                    <span>{{
                      detail.presell_type === 2
                        ? `¥${Number(detail.down_price || 0).toFixed(2)} / ¥${Number(detail.final_price || 0).toFixed(2)}`
                        : '不适用'
                    }}</span>
                  </ElFormItem>
                  <ElFormItem label="活动状态">
                    <ElTag size="small">{{
                      detail.presell_status_text || '—'
                    }}</ElTag>
                  </ElFormItem>
                  <ElFormItem label="审核状态">
                    <span>{{ detail.product_status_name || '—' }}</span>
                  </ElFormItem>
                  <ElFormItem label="预售活动日期">
                    <span
                      >{{ formatShanghaiDateTime(detail.start_time) }} 至
                      {{ formatShanghaiDateTime(detail.end_time) }}</span
                    >
                  </ElFormItem>
                  <ElFormItem label="尾款支付时间">
                    <span>{{
                      detail.presell_type === 2
                        ? `${formatShanghaiDateTime(detail.final_start_time)} 至 ${formatShanghaiDateTime(detail.final_end_time)}`
                        : '不适用'
                    }}</span>
                  </ElFormItem>
                  <ElFormItem label="限量总数 / 剩余">
                    <span
                      >{{ detail.stock_count ?? 0 }} /
                      {{ detail.stock ?? 0 }}</span
                    >
                  </ElFormItem>
                  <ElFormItem label="限购数量">
                    <span>{{
                      Number(detail.pay_count || 0) > 0
                        ? detail.pay_count
                        : '不限'
                    }}</span>
                  </ElFormItem>
                  <ElFormItem label="已售 / 成功 / 参与">
                    <span
                      >{{ detail.seles ?? 0 }} / {{ detail.success_num ?? 0 }} /
                      {{ detail.attend_num ?? 0 }}</span
                    >
                  </ElFormItem>
                  <ElFormItem label="发货">
                    <span
                      >{{
                        detail.delivery_type === 2 ? '预售结束后' : '付款后'
                      }}{{
                        detail.delivery_day ? ` ${detail.delivery_day} 天` : ''
                      }}</span
                    >
                  </ElFormItem>
                  <ElFormItem label="显示状态">
                    <span>{{ detail.is_show === 1 ? '显示' : '隐藏' }}</span>
                  </ElFormItem>
                  <ElFormItem label="推荐级别">
                    <ElRate
                      :model-value="Number(detail.star || 0)"
                      disabled
                      :max="5"
                    />
                  </ElFormItem>
                  <ElFormItem v-if="detail.refusal" label="拒绝/下架原因">
                    <span>{{ detail.refusal }}</span>
                  </ElFormItem>
                  <ElFormItem label="活动说明">
                    <span>{{ detail.store_info || '—' }}</span>
                  </ElFormItem>
                </ElForm>
              </template>
            </PresellProductTabs>

            <div class="presell-detail__footer">
              <ElButton @click="close">关闭</ElButton>
            </div>
          </template>
        </template>
      </ElSkeleton>
    </div>
  </Drawer>
</template>

<style scoped>
.presell-detail {
  display: flex;
  flex-direction: column;
  min-height: 420px;
}

.presell-detail__footer {
  position: sticky;
  bottom: 0;
  z-index: 2;
  display: flex;
  gap: 8px;
  margin-top: 8px;
  padding: 12px 4px 0;
  border-top: 1px solid hsl(var(--border));
  background: hsl(var(--background));
}
</style>
