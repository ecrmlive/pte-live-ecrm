<script setup lang="ts">
import type { AttachmentItem } from '#/api/core/attachment';

/**
 * 统一「选择视频素材」弹窗。
 * 结构与图片素材选择器一致：侧栏分类 + 网格 + 上传 + 底部分页 + 取消/确定。
 * 默认系统分类：全部素材、店铺视频、商品视频、其他视频。
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
