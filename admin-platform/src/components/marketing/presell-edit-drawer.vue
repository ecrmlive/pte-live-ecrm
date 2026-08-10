<script setup lang="ts">
import { reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElRate,
  ElSelect,
  ElSkeleton,
  ElSwitch,
  ElTag,
} from 'element-plus';

import type { PlatformProductEditDetail } from '#/api/core/platform-catalog';
import {
  updatePlatformPresellApi,
  type PlatformPresell,
} from '#/api/core/platform-presell';
import { formatShanghaiDateTime } from '#/utils/date-time';

import { loadPresellProductBundle } from './load-presell-product';
import PresellProductTabs from './presell-product-tabs.vue';

const emit = defineEmits<{ saved: [] }>();

const loading = ref(false);
const saving = ref(false);
const activeTab = ref('presell');
const productMissing = ref(false);
const presellId = ref(0);
const detail = ref<PlatformPresell>();
const product = ref<PlatformProductEditDetail>();

const form = reactive({
  store_name: '',
  store_info: '',
  price: 0,
  down_price: 0,
  final_price: 0,
  stock: 0,
  stock_count: 0,
  pay_count: 0,
  delivery_type: 1,
  delivery_day: 0,
  dateRange: [] as string[],
  finalRange: [] as string[],
  is_show: 1,
  star: 0,
});

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '编辑预售商品',
});

function toPickerValue(raw?: string) {
  if (!raw) return '';
  const formatted = formatShanghaiDateTime(raw);
  return formatted === '—' ? '' : formatted;
}

function fillForm(row: PlatformPresell) {
  form.store_name = row.store_name || '';
  form.store_info = row.store_info || '';
  form.price = Number(row.price || 0);
  form.down_price = Number(row.down_price || 0);
  form.final_price = Number(row.final_price || 0);
  form.stock = Number(row.stock || 0);
  form.stock_count = Number(row.stock_count || 0);
  form.pay_count = Number(row.pay_count || 0);
  form.delivery_type = Number(row.delivery_type || 1);
  form.delivery_day = Number(row.delivery_day || 0);
  form.dateRange = [
    toPickerValue(row.start_time),
    toPickerValue(row.end_time),
  ].filter(Boolean) as string[];
  form.finalRange =
    row.presell_type === 2
      ? ([
          toPickerValue(row.final_start_time),
          toPickerValue(row.final_end_time),
        ].filter(Boolean) as string[])
      : [];
  form.is_show = Number(row.is_show ?? 1);
  form.star = Number(row.star || 0);
}

async function open(id: number) {
  presellId.value = id;
  activeTab.value = 'presell';
  productMissing.value = false;
  detail.value = undefined;
  product.value = undefined;
  drawerApi.setState({ loading: true, title: '编辑预售商品' }).open();
  loading.value = true;
  try {
    const bundle = await loadPresellProductBundle(id);
    detail.value = bundle.presell;
    product.value = bundle.product;
    productMissing.value = bundle.productMissing;
    fillForm(bundle.presell);
  } catch {
    ElMessage.error('加载预售编辑数据失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

function close() {
  drawerApi.close();
}

async function save() {
  if (!presellId.value) return;
  if (!form.store_name.trim()) {
    ElMessage.warning('请填写预售标题');
    activeTab.value = 'presell';
    return;
  }
  if (!form.dateRange?.[0] || !form.dateRange?.[1]) {
    ElMessage.warning('请选择预售活动日期');
    activeTab.value = 'presell';
    return;
  }
  saving.value = true;
  drawerApi.lock();
  try {
    await updatePlatformPresellApi(presellId.value, {
      store_name: form.store_name.trim(),
      store_info: form.store_info.trim(),
      price: form.price,
      down_price: form.down_price,
      final_price: form.final_price,
      stock: form.stock,
      stock_count: form.stock_count,
      pay_count: form.pay_count,
      delivery_type: form.delivery_type,
      delivery_day: form.delivery_day,
      start_time: form.dateRange[0],
      end_time: form.dateRange[1],
      final_start_time: form.finalRange[0] || '',
      final_end_time: form.finalRange[1] || '',
      is_show: form.is_show,
      star: form.star,
    });
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

defineExpose({ open, close });
</script>

<template>
  <Drawer>
    <div v-loading="loading || saving" class="presell-edit">
      <ElSkeleton :loading="loading && !detail" animated :rows="10">
        <template #default>
          <template v-if="detail">
            <PresellProductTabs
              v-model="activeTab"
              show-presell-tab
              :presell="detail"
              :product="product"
              :product-missing="productMissing"
              product-missing-tip="关联商品编辑信息暂不可用，仍可编辑预售活动信息。"
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
                  <ElFormItem label="预售类型">
                    <span>{{
                      detail.presell_type === 2 ? '定金预售' : '全款预售'
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
                  <ElFormItem label="预售标题" required>
                    <ElInput
                      v-model="form.store_name"
                      maxlength="128"
                      show-word-limit
                    />
                  </ElFormItem>
                  <ElFormItem label="预售价" required>
                    <ElInputNumber
                      v-model="form.price"
                      :min="0"
                      :precision="2"
                      :step="1"
                    />
                  </ElFormItem>
                  <template v-if="detail.presell_type === 2">
                    <ElFormItem label="定金">
                      <ElInputNumber
                        v-model="form.down_price"
                        :min="0"
                        :precision="2"
                        :step="1"
                      />
                    </ElFormItem>
                    <ElFormItem label="尾款">
                      <ElInputNumber
                        v-model="form.final_price"
                        :min="0"
                        :precision="2"
                        :step="1"
                      />
                    </ElFormItem>
                    <ElFormItem label="尾款支付时间">
                      <ElDatePicker
                        v-model="form.finalRange"
                        type="datetimerange"
                        value-format="YYYY-MM-DD HH:mm:ss"
                        start-placeholder="开始"
                        end-placeholder="结束"
                        class="!w-full"
                      />
                    </ElFormItem>
                  </template>
                  <ElFormItem label="活动日期" required>
                    <ElDatePicker
                      v-model="form.dateRange"
                      type="datetimerange"
                      value-format="YYYY-MM-DD HH:mm:ss"
                      start-placeholder="开始"
                      end-placeholder="结束"
                      class="!w-full"
                    />
                  </ElFormItem>
                  <ElFormItem label="限量总数">
                    <ElInputNumber
                      v-model="form.stock_count"
                      :min="0"
                      :step="1"
                    />
                  </ElFormItem>
                  <ElFormItem label="限量剩余">
                    <ElInputNumber v-model="form.stock" :min="0" :step="1" />
                  </ElFormItem>
                  <ElFormItem label="限购数量">
                    <ElInputNumber v-model="form.pay_count" :min="0" :step="1" />
                    <span class="ml-2 text-muted-foreground text-xs"
                      >0 表示不限</span
                    >
                  </ElFormItem>
                  <ElFormItem label="发货类型">
                    <ElSelect v-model="form.delivery_type" class="!w-[220px]">
                      <ElOption :value="1" label="付款后" />
                      <ElOption :value="2" label="预售结束后" />
                    </ElSelect>
                  </ElFormItem>
                  <ElFormItem label="发货天数">
                    <ElInputNumber
                      v-model="form.delivery_day"
                      :min="0"
                      :step="1"
                    />
                  </ElFormItem>
                  <ElFormItem label="显示状态">
                    <ElSwitch
                      v-model="form.is_show"
                      :active-value="1"
                      :inactive-value="0"
                      inline-prompt
                      active-text="显示"
                      inactive-text="隐藏"
                    />
                  </ElFormItem>
                  <ElFormItem label="推荐级别">
                    <ElRate v-model="form.star" :max="5" />
                  </ElFormItem>
                  <ElFormItem label="成功 / 参与">
                    <span
                      >{{ detail.success_num ?? 0 }} /
                      {{ detail.attend_num ?? 0 }}</span
                    >
                  </ElFormItem>
                  <ElFormItem label="已售">
                    <span>{{ detail.seles ?? 0 }}</span>
                  </ElFormItem>
                  <ElFormItem v-if="detail.refusal" label="拒绝/下架原因">
                    <span>{{ detail.refusal }}</span>
                  </ElFormItem>
                  <ElFormItem label="活动说明">
                    <ElInput
                      v-model="form.store_info"
                      type="textarea"
                      :rows="3"
                      maxlength="500"
                      show-word-limit
                    />
                  </ElFormItem>
                </ElForm>
              </template>
            </PresellProductTabs>

            <div class="presell-edit__footer">
              <ElButton type="primary" :loading="saving" @click="save">
                保存
              </ElButton>
              <ElButton :disabled="saving" @click="close">取消</ElButton>
            </div>
          </template>
        </template>
      </ElSkeleton>
    </div>
  </Drawer>
</template>

<style scoped>
.presell-edit {
  display: flex;
  flex-direction: column;
  min-height: 420px;
}

.presell-edit__footer {
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
