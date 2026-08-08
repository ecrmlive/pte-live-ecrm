<script setup lang="ts">
import {
  ArrowDownBold,
  ArrowUpBold,
  Delete,
  DocumentCopy,
} from '@element-plus/icons-vue';
import { ElIcon } from 'element-plus';

defineProps<{
  disabled?: boolean;
  canMoveUp?: boolean;
  canMoveDown?: boolean;
}>();

const emit = defineEmits<{
  up: [];
  down: [];
  copy: [];
  remove: [];
}>();
</script>

<template>
  <div class="field-tools">
    <button
      type="button"
      class="tool"
      :disabled="disabled || !canMoveUp"
      title="上移"
      @click.stop="emit('up')"
    >
      <ElIcon :size="14"><ArrowUpBold /></ElIcon>
    </button>
    <button
      type="button"
      class="tool"
      :disabled="disabled || !canMoveDown"
      title="下移"
      @click.stop="emit('down')"
    >
      <ElIcon :size="14"><ArrowDownBold /></ElIcon>
    </button>
    <button
      type="button"
      class="tool"
      :disabled="disabled"
      title="复制"
      @click.stop="emit('copy')"
    >
      <ElIcon :size="14"><DocumentCopy /></ElIcon>
    </button>
    <button
      type="button"
      class="tool"
      :disabled="disabled"
      title="删除"
      @click.stop="emit('remove')"
    >
      <ElIcon :size="14"><Delete /></ElIcon>
    </button>
  </div>
</template>

<style scoped lang="scss">
.field-tools {
  position: absolute;
  top: 0;
  right: -43px;
  display: flex;
  flex-direction: column;
  width: 36px;
  padding: 4px 0;
  border-radius: 4px;
  background: hsl(var(--primary));
  color: #fff;
  font-weight: 700;
  text-align: center;
  cursor: pointer;
}

.tool {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: auto;
  padding: 5px 0;
  border: 0;
  background: transparent;
  color: #fff;
  cursor: pointer;
}

.tool:disabled {
  cursor: not-allowed;
  opacity: 0.4;
}
</style>
