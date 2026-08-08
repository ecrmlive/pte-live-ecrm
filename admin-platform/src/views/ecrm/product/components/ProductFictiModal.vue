<script setup lang="ts">
import { reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
} from 'element-plus';

import { setPlatformProductFictiApi } from '#/api/core/platform-catalog';

const emit = defineEmits<{
  success: [payload: { ficti: number; productId: number }];
}>();

const submitting = ref(false);
const productId = ref(0);
const form = reactive({
  current: 0,
  type: 1 as 1 | 2,
  ficti: 0,
});

const [Modal, modalApi] = useVbenModal({
  title: '修改已售数量',
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确定',
  showCancelButton: false,
  destroyOnClose: true,
  onConfirm: async () => {
    if (submitting.value) return;
    const amount = Number(form.ficti || 0);
    if (!Number.isFinite(amount) || amount <= 0) {
      ElMessage.warning('已售数量必须大于0');
      return;
    }
    if (form.type === 2 && amount > Number(form.current || 0)) {
      ElMessage.warning('已售数量不足');
      return;
    }
    submitting.value = true;
    modalApi.setState({ confirmLoading: true });
    try {
      await setPlatformProductFictiApi(productId.value, {
        type: form.type,
        ficti: amount,
      });
      const next =
        form.type === 1
          ? Number(form.current || 0) + amount
          : Number(form.current || 0) - amount;
      ElMessage.success('已售数量已更新');
      emit('success', { productId: productId.value, ficti: next });
      modalApi.close();
    } catch {
      /* requestClient 已提示 */
    } finally {
      submitting.value = false;
      modalApi.setState({ confirmLoading: false });
    }
  },
});

function open(payload: { productId: number; ficti?: number }) {
  productId.value = Number(payload.productId || 0);
  form.current = Math.max(0, Number(payload.ficti || 0));
  form.type = 1;
  form.ficti = 0;
  modalApi.open();
}

defineExpose({ open });
</script>

<template>
  <Modal>
    <ElForm label-width="120px" class="product-ficti-form" @submit.prevent>
      <ElFormItem label="现有已售数量:">
        <ElInput :model-value="String(form.current)" disabled />
      </ElFormItem>
      <ElFormItem label="修改类型:">
        <ElRadioGroup v-model="form.type">
          <ElRadio :label="1">增加</ElRadio>
          <ElRadio :label="2">减少</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="修改已售数量:">
        <ElInputNumber
          v-model="form.ficti"
          :min="0"
          :precision="0"
          :step="1"
          class="product-ficti-form__stepper"
        />
      </ElFormItem>
    </ElForm>
  </Modal>
</template>

<style scoped>
.product-ficti-form :deep(.el-input.is-disabled .el-input__wrapper) {
  background-color: var(--el-fill-color-light);
}

.product-ficti-form__stepper {
  width: 180px;
}

/* 对齐设计图：左右 -/ + 步进器 */
.product-ficti-form__stepper :deep(.el-input-number) {
  width: 180px;
}

.product-ficti-form :deep(.el-input-number) {
  width: 180px;
}

.product-ficti-form :deep(.el-input-number .el-input-number__decrease),
.product-ficti-form :deep(.el-input-number .el-input-number__increase) {
  background: var(--el-fill-color-blank);
}
</style>
