<script setup lang="ts">
import type { ShopFileItem } from '#/api/core/file';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, ref, watch, nextTick } from 'vue';

import MediaLibraryPanel from '#/components/shop/media-library-panel.vue';

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    fileType?: 'image' | 'video';
    limit?: number;
  }>(),
  {
    fileType: 'image',
    limit: 1,
  },
);

const emit = defineEmits<{
  select: [ShopFileItem[]];
}>();

const panelRef = ref<InstanceType<typeof MediaLibraryPanel> | null>(null);

const dialogTitle = computed(() =>
  props.fileType === 'video' ? '选择视频素材' : '选择图片素材',
);

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
    ElMessage.warning(props.fileType === 'video' ? '请选择视频' : '请选择图片');
    return;
  }
  if (selected.length > props.limit) {
    ElMessage.warning(`最多选择 ${props.limit} 个文件`);
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
    class="media-picker-dialog w-[910px] max-w-[96vw]"
    content-class="media-picker-dialog__content"
    :title="dialogTitle"
  >
    <div class="media-picker-dialog__body">
      <MediaLibraryPanel
        ref="panelRef"
        :file-type="fileType"
        picker-mode
        :selection-limit="limit"
      />
    </div>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>

<style lang="scss">
.media-picker-dialog__content {
  display: flex;
  flex-direction: column;
  max-height: min(72vh, 680px);
  overflow: hidden;
  padding: 12px 16px 8px;
}
</style>

<style scoped lang="scss">
.media-picker-dialog__body {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
}
</style>
