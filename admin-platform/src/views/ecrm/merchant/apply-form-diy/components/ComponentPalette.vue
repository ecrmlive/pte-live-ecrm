<script setup lang="ts">
import { ElIcon } from 'element-plus';

import { COMPONENT_LIBRARY } from '../registry';
import type { ApplyFieldType } from '../types';

defineProps<{ disabled?: boolean }>();
const emit = defineEmits<{
  add: [type: ApplyFieldType];
}>();
</script>

<template>
  <aside class="diy-left">
    <div class="diy-left__group">
      <div class="diy-left__tips">组件</div>
      <div class="diy-left__palette">
        <button
          v-for="item in COMPONENT_LIBRARY"
          :key="item.type"
          type="button"
          class="palette-item"
          :disabled="disabled"
          @click="emit('add', item.type)"
        >
          <ElIcon class="palette-item__icon" :size="22">
            <component :is="item.icon" />
          </ElIcon>
          <span>{{ item.label }}</span>
        </button>
      </div>
    </div>
  </aside>
</template>

<style scoped lang="scss">
.diy-left {
  position: relative;
  z-index: 2;
  flex-shrink: 0;
  width: 300px;
  background: #fff;
  border-right: 1px solid #eee;
  overflow: auto;
  pointer-events: auto;
}

.diy-left__group {
  padding: 15px;
}

.diy-left__tips {
  margin-bottom: 15px;
  color: #000;
  font-size: 13px;
}

.diy-left__palette {
  display: flex;
  flex-wrap: wrap;
}

.palette-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 74px;
  height: 66px;
  margin-right: 17px;
  margin-bottom: 10px;
  padding: 0;
  border: 0;
  border-radius: 5px;
  background: transparent;
  color: #666;
  font-size: 12px;
  cursor: pointer;
}

.palette-item:nth-child(3n) {
  margin-right: 0;
}

.palette-item:hover:not(:disabled) {
  box-shadow: 0 0 5px 0 rgb(24 144 255 / 30%);
}

.palette-item:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.palette-item__icon {
  margin-bottom: 4px;
  color: hsl(var(--primary));
}
</style>
