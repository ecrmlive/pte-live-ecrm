<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { ElEmpty, ElMessage, ElTag } from 'element-plus';

import {
  fetchProductLabels,
  type ProductLabelRow,
} from '#/api/core/ecrm';

const emit = defineEmits<{
  submit: [string[]];
}>();

const loading = ref(false);
const dropdownOpen = ref(false);
const options = ref<ProductLabelRow[]>([]);
const draftIds = ref<string[]>([]);
const productId = ref(0);

const selectedOptions = computed(() => {
  const map = new Map(options.value.map((o) => [String(o.id), o]));
  return draftIds.value
    .map((id) => map.get(id))
    .filter((o): o is ProductLabelRow => Boolean(o));
});

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
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      dropdownOpen.value = false;
    }
  },
});

async function open(payload: {
  productId?: number;
  selectedIds?: string[];
  options?: ProductLabelRow[];
}) {
  productId.value = Number(payload.productId || 0);
  draftIds.value = [...(payload.selectedIds || [])].map(String);
  dropdownOpen.value = false;
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

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value;
}

function isSelected(id: number | string) {
  return draftIds.value.includes(String(id));
}

function toggleOption(id: number | string) {
  const key = String(id);
  if (draftIds.value.includes(key)) {
    draftIds.value = draftIds.value.filter((x) => x !== key);
  } else {
    draftIds.value = [...draftIds.value, key];
  }
}

function removeSelected(id: number | string) {
  draftIds.value = draftIds.value.filter((x) => x !== String(id));
}

defineExpose({ open });
</script>

<template>
  <Modal>
    <div class="label-select" v-loading="loading">
      <div
        class="label-select__control"
        :class="{ 'is-open': dropdownOpen }"
        @click="toggleDropdown"
      >
        <div class="label-select__tags">
          <ElTag
            v-for="item in selectedOptions"
            :key="item.id"
            class="label-select__tag"
            closable
            type="info"
            effect="plain"
            @click.stop
            @close.stop="removeSelected(item.id)"
          >
            {{ item.name }}
          </ElTag>
          <span v-if="!selectedOptions.length" class="label-select__placeholder">
            请选择标签
          </span>
        </div>
        <IconifyIcon
          class="label-select__arrow"
          :icon="
            dropdownOpen
              ? 'ant-design:up-outlined'
              : 'ant-design:down-outlined'
          "
        />
      </div>

      <div v-show="dropdownOpen" class="label-select__dropdown">
        <button
          v-for="item in options"
          :key="item.id"
          type="button"
          class="label-select__option"
          :class="{ 'is-selected': isSelected(item.id) }"
          @click.stop="toggleOption(item.id)"
        >
          <span>{{ item.name }}</span>
          <IconifyIcon
            v-if="isSelected(item.id)"
            class="label-select__check"
            icon="ant-design:check-outlined"
          />
        </button>
        <ElEmpty
          v-if="!options.length"
          :image-size="64"
          description="暂无标签，请先在商品标签中维护"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.label-select {
  display: flex;
  flex-direction: column;
  gap: 0;
  min-height: 120px;
  padding: 8px 0 4px;
}

.label-select__control {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  min-height: 40px;
  padding: 6px 12px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  transition: border-color 0.15s ease;
}

.label-select__control.is-open,
.label-select__control:hover {
  border-color: hsl(var(--primary));
}

.label-select__tags {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
  min-width: 0;
}

.label-select__tag {
  margin: 0;
  background: #f0f2f5;
  border-color: transparent;
  color: #303133;
}

.label-select__placeholder {
  color: hsl(var(--muted-foreground));
  font-size: 14px;
  line-height: 28px;
}

.label-select__arrow {
  flex-shrink: 0;
  margin-top: 6px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.label-select__dropdown {
  z-index: 2;
  max-height: 280px;
  margin-top: 4px;
  overflow: auto;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  background: #fff;
  box-shadow: 0 6px 16px rgb(0 0 0 / 8%);
}

.label-select__option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 10px 14px;
  border: 0;
  background: transparent;
  color: #303133;
  font-size: 14px;
  line-height: 22px;
  text-align: left;
  cursor: pointer;
}

.label-select__option:hover {
  background: #ecf5ff;
}

.label-select__option.is-selected {
  color: hsl(var(--primary));
}

.label-select__check {
  flex-shrink: 0;
  color: hsl(var(--primary));
  font-size: 16px;
}
</style>
