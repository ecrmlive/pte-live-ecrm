<script setup lang="ts">
import { reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElMessage,
  ElSkeleton,
  ElSwitch,
  ElTag,
} from 'element-plus';

import type { PlatformProductEditDetail } from '#/api/core/platform-catalog';
import {
  updatePlatformAssistApi,
  type PlatformAssistActive,
} from '#/api/core/platform-assist';
import ImageField from '#/components/shop/image-field.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';

import { loadAssistProductBundle } from './load-assist-product';

const emit = defineEmits<{ saved: [] }>();

const loading = ref(false);
const saving = ref(false);
const productMissing = ref(false);
const assistId = ref(0);
const detail = ref<PlatformAssistActive>();
const product = ref<PlatformProductEditDetail>();
const cover = ref('');

const form = reactive({
  is_show: 1,
});

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '编辑好友助力活动',
});

function time(value?: string) {
  return formatShanghaiDateTime(value || '');
}

function fillForm(row: PlatformAssistActive) {
  form.is_show = Number(row.is_show ?? 1);
}

async function open(id: number) {
  assistId.value = id;
  productMissing.value = false;
  detail.value = undefined;
  product.value = undefined;
  cover.value = '';
  drawerApi.setState({ loading: true, title: '编辑好友助力活动' }).open();
  loading.value = true;
  try {
    const bundle = await loadAssistProductBundle(id);
    detail.value = bundle.assist;
    product.value = bundle.product;
    productMissing.value = bundle.productMissing;
    cover.value =
      bundle.product?.image ||
      bundle.assist.image ||
      '';
    fillForm(bundle.assist);
  } catch {
    ElMessage.error('加载助力编辑数据失败');
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
  if (!assistId.value) return;
  saving.value = true;
  drawerApi.lock();
  try {
    await updatePlatformAssistApi(assistId.value, {
      is_show: form.is_show,
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
    <div v-loading="loading || saving" class="assist-edit">
      <ElSkeleton :loading="loading && !detail" animated :rows="10">
        <template #default>
          <template v-if="detail">
            <ElAlert
              class="mb-4"
              type="info"
              :closable="false"
              title="平台监管编辑仅开放前台展示状态。助力价、库存与完成条件须保持订单快照，本页只读。"
            />

            <ElForm label-position="left" label-width="120px" class="pr-2">
              <ElFormItem label="活动 ID">
                <span>{{ detail.product_assist_id }}</span>
              </ElFormItem>
              <ElFormItem label="活动 / 商品">
                <span>{{ detail.store_name || `商品 #${detail.product_id}` }}</span>
              </ElFormItem>
              <ElFormItem label="关联商品 ID">
                <span>{{ detail.product_id || '—' }}</span>
              </ElFormItem>
              <ElFormItem label="商户">
                <span>{{ detail.mer_name || `商户 #${detail.mer_id}` }}</span>
              </ElFormItem>
              <ElFormItem label="活动封面">
                <ImageField v-model="cover" disabled />
              </ElFormItem>
              <ElFormItem label="助力价">
                <span>¥{{ Number(detail.assist_price || 0).toFixed(2) }}</span>
              </ElFormItem>
              <ElFormItem label="助力规则">
                <span>
                  {{ detail.assist_count }} 人 / 每人最多
                  {{ detail.assist_user_count }} 次
                </span>
              </ElFormItem>
              <ElFormItem label="活动库存">
                <span>{{ detail.stock }}</span>
              </ElFormItem>
              <ElFormItem label="活动时间">
                <span>{{ time(detail.start_time) }} 至 {{ time(detail.end_time) }}</span>
              </ElFormItem>
              <ElFormItem label="活动状态">
                <ElTag size="small" :type="detail.status === 1 ? 'success' : 'info'">
                  {{ detail.status === 1 ? '启用' : '停用' }}
                </ElTag>
              </ElFormItem>
              <ElFormItem label="前台展示" required>
                <ElSwitch
                  v-model="form.is_show"
                  :active-value="1"
                  :inactive-value="0"
                  inline-prompt
                  active-text="上架"
                  inactive-text="下架"
                />
              </ElFormItem>
            </ElForm>

            <ElAlert
              v-if="productMissing"
              class="mb-4"
              type="warning"
              :closable="false"
              title="关联商品编辑信息暂不可用；仍可保存助力展示状态。可用「编辑标签」在商品源表补齐后重试。"
            />

            <ElDescriptions
              v-if="product"
              class="mb-4"
              :column="2"
              border
              title="关联商品（只读）"
            >
              <ElDescriptionsItem label="商品名称" :span="2">
                {{ product.title || product.store_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="品牌">
                {{ product.brand_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="单位">
                {{ product.unit_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="售价">
                ¥{{ Number(product.price || 0).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="原价">
                ¥{{ Number(product.ot_price || 0).toFixed(2) }}
              </ElDescriptionsItem>
            </ElDescriptions>

            <div class="assist-edit__footer">
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
.assist-edit {
  display: flex;
  flex-direction: column;
  min-height: 420px;
}

.assist-edit__footer {
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
