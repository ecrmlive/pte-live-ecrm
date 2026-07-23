<script setup lang="ts">
import type { ShopFileItem } from '#/api/core/file';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, ref, watch, nextTick } from 'vue';

import MediaLibraryPanel from '#/components/shop/media-library-panel.vue';

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    limit?: number;
    /** 每次打开弹窗时的默认素材库（resetPickerState） */
    defaultLibrary?: 'merchant' | 'system';
  }>(),
  { defaultLibrary: 'merchant', limit: 1 },
);

const emit = defineEmits<{
  select: [ShopFileItem[]];
}>();

const panelRef = ref<InstanceType<typeof MediaLibraryPanel> | null>(null);

const selectionHint = computed(() => `已选 ${props.limit === 1 ? '1' : `最多 ${props.limit}`} 张`);

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void nextTick(() => {
        panelRef.value?.resetPickerState();
      });
    }
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      modalApi.open();
      return;
    }
    modalApi.close();
  },
  { immediate: true },
);

function confirm() {
  const selected = panelRef.value?.getSelectedFiles() ?? [];
  if (!selected.length) {
    ElMessage.warning('请选择图片');
    return;
  }
  if (selected.length > props.limit) {
    ElMessage.warning(`最多选择 ${props.limit} 张图片`);
    return;
  }
  emit('select', [...selected]);
  open.value = false;
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="image-picker-dialog w-[910px] max-w-[96vw]"
    content-class="image-picker-dialog__content"
    title="选择图片素材"
  >
    <div class="image-picker-dialog__body">
      <MediaLibraryPanel
        ref="panelRef"
        :default-library="defaultLibrary"
        enable-system-library
        file-type="image"
        picker-mode
        :selection-limit="limit"
      />
      <p class="image-picker-dialog__hint">{{ selectionHint }}</p>
    </div>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>

<style lang="scss">
.image-picker-dialog__content {
  display: flex;
  flex-direction: column;
  max-height: min(72vh, 680px);
  overflow: hidden;
  padding: 12px 16px 8px;
}
</style>

<style scoped lang="scss">
.image-picker-dialog__body {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
}

.image-picker-dialog__hint {
  flex-shrink: 0;
  margin: 0;
  font-size: 12px;
  color: hsl(var(--muted-foreground));
  text-align: right;
}
</style>
