<script setup lang="ts">
import { reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSkeleton,
} from 'element-plus';

import {
  auditPlatformPresellApi,
  type PlatformPresell,
} from '#/api/core/platform-presell';
import type { PlatformProductEditDetail } from '#/api/core/platform-catalog';

import { loadPresellProductBundle } from './load-presell-product';
import PresellProductTabs from './presell-product-tabs.vue';

const emit = defineEmits<{
  audited: [];
}>();

const loading = ref(false);
const submitting = ref(false);
const activeTab = ref('basic');
const productMissing = ref(false);
const presell = ref<PlatformPresell>();
const product = ref<PlatformProductEditDetail>();

const auditForm = reactive({
  status: 1 as 1 | -1,
  refusal: '',
});

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '商品审核',
});

function resetAudit() {
  auditForm.status = 1;
  auditForm.refusal = '';
}

async function open(id: number) {
  activeTab.value = 'basic';
  productMissing.value = false;
  presell.value = undefined;
  product.value = undefined;
  resetAudit();
  drawerApi.setState({ loading: true, title: '商品审核' }).open();
  loading.value = true;
  try {
    const bundle = await loadPresellProductBundle(id);
    presell.value = bundle.presell;
    product.value = bundle.product;
    productMissing.value = bundle.productMissing;
  } catch {
    ElMessage.error('加载预售审核数据失败');
    drawerApi.close();
  } finally {
    loading.value = false;
    drawerApi.setState({ loading: false });
  }
}

function close() {
  drawerApi.close();
}

async function submit() {
  const id = Number(presell.value?.product_presell_id || 0);
  if (!id) return;
  if (auditForm.status === -1 && !auditForm.refusal.trim()) {
    ElMessage.warning('请填写拒绝原因');
    return;
  }
  submitting.value = true;
  drawerApi.lock();
  try {
    await auditPlatformPresellApi(
      id,
      auditForm.status,
      auditForm.status === -1 ? auditForm.refusal.trim() : '',
    );
    ElMessage.success(auditForm.status === 1 ? '已通过审核' : '已拒绝');
    drawerApi.close();
    emit('audited');
  } catch {
    ElMessage.error('提交审核失败');
  } finally {
    submitting.value = false;
    drawerApi.unlock();
  }
}

defineExpose({ open, close });
</script>

<template>
  <Drawer>
    <div v-loading="loading || submitting" class="presell-audit">
      <ElSkeleton :loading="loading && !presell" animated :rows="10">
        <template #default>
          <template v-if="presell">
            <PresellProductTabs
              v-model="activeTab"
              :presell="presell"
              :product="product"
              :product-missing="productMissing"
              product-missing-tip="关联商品编辑信息暂不可用，仍可审核预售活动状态。"
            />

            <div class="presell-audit__footer">
              <ElForm label-width="96px" class="presell-audit__audit-form">
                <ElFormItem label="审核状态" required>
                  <ElRadioGroup v-model="auditForm.status">
                    <ElRadio :value="1">通过</ElRadio>
                    <ElRadio :value="-1">拒绝</ElRadio>
                  </ElRadioGroup>
                </ElFormItem>
                <ElFormItem
                  v-if="auditForm.status === -1"
                  label="拒绝原因"
                  required
                >
                  <ElInput
                    v-model="auditForm.refusal"
                    type="textarea"
                    :rows="3"
                    maxlength="200"
                    show-word-limit
                    placeholder="请填写拒绝原因"
                  />
                </ElFormItem>
                <ElFormItem>
                  <ElButton type="primary" :loading="submitting" @click="submit">
                    提交
                  </ElButton>
                  <ElButton :disabled="submitting" @click="close">关闭</ElButton>
                </ElFormItem>
              </ElForm>
            </div>
          </template>
        </template>
      </ElSkeleton>
    </div>
  </Drawer>
</template>

<style scoped>
.presell-audit {
  display: flex;
  flex-direction: column;
  min-height: 420px;
}

.presell-audit__footer {
  position: sticky;
  bottom: 0;
  z-index: 2;
  margin-top: 8px;
  padding: 12px 4px 0;
  border-top: 1px solid hsl(var(--border));
  background: hsl(var(--background));
}

.presell-audit__audit-form {
  max-width: 720px;
}
</style>
