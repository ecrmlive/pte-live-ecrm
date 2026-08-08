<script setup lang="ts">
import type { AttachmentItem } from '#/api/core/attachment';

/**
 * 商户端「选择视频素材」弹窗（与平台 image-picker kind=video 对齐）。
 */
import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';

type PickerItem = AttachmentItem & { file_id: number; file_path: string };

const open = defineModel<boolean>('open', { default: false });
const props = withDefaults(
  defineProps<{
    limit?: number;
  }>(),
  { limit: 1 },
);
const emit = defineEmits<{ select: [PickerItem[]] }>();

function onSelect(rows: PickerItem[]) {
  emit('select', rows);
}
</script>

<template>
  <ImagePickerDialog
    v-model:open="open"
    kind="video"
    :limit="props.limit"
    @select="onSelect"
  />
</template>
