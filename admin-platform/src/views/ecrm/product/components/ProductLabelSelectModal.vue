<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElForm,
  ElFormItem,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import {
  fetchProductLabels,
  type ProductLabelRow,
} from '#/api/core/ecrm';

const emit = defineEmits<{
  submit: [string[]];
}>();

const loading = ref(false);
const options = ref<ProductLabelRow[]>([]);
const draftIds = ref<string[]>([]);

const [Modal, modalApi] = useVbenModal({
  title: '选择标签',
  class: 'w-[560px] max-w-[96vw]',
  confirmText: '提交',
  cancelText: '取消',
  destroyOnClose: true,
  onConfirm: async () => {
    emit('submit', [...draftIds.value]);
    modalApi.close();
  },
});

async function open(payload: {
  productId?: number;
  selectedIds?: string[];
  options?: ProductLabelRow[];
}) {
  draftIds.value = [...(payload.selectedIds || [])]
    .map(String)
    .filter(Boolean);
  modalApi.open();
  if (payload.options?.length) {
    options.value = payload.options;
    return;
  }
  loading.value = true;
  modalApi.setState({ loading: true });
  try {
    const data = await fetchProductLabels();
    options.value = (data.list || []).filter((x) => Number(x.status) !== 0);
  } catch {
    options.value = [];
    ElMessage.error('加载商品标签失败');
  } finally {
    loading.value = false;
    modalApi.setState({ loading: false });
  }
}

defineExpose({ open });
</script>

<template>
  <Modal>
    <div v-loading="loading" class="label-select">
      <ElForm label-width="60px">
        <ElFormItem label="标签">
          <ElSelect
            v-model="draftIds"
            multiple
            filterable
            clearable
            collapse-tags
            collapse-tags-tooltip
            placeholder="请选择"
            class="label-select__control"
          >
            <ElOption
              v-for="item in options"
              :key="item.id"
              :label="item.name"
              :value="String(item.id)"
            />
          </ElSelect>
        </ElFormItem>
      </ElForm>
    </div>
  </Modal>
</template>

<style scoped>
.label-select {
  min-height: 80px;
  padding: 8px 0 4px;
}

.label-select__control {
  width: 100%;
}
</style>
