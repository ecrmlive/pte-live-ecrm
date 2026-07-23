<script setup lang="ts">
import { Picture, Plus } from '@element-plus/icons-vue';
import { computed } from 'vue';

import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const props = withDefaults(
  defineProps<{
    buttonText?: string;
    hint?: string;
    previewUrl?: string;
    size?: number;
  }>(),
  {
    buttonText: '选择图片',
    hint: '',
    previewUrl: '',
    size: 100,
  },
);

const emit = defineEmits<{
  pick: [];
}>();

const resolvedPreviewUrl = computed(() =>
  resolveCosMediaUrl(props.previewUrl ?? ''),
);
</script>

<template>
  <div class="image-picker-trigger">
    <el-button
      class="image-picker-trigger__btn"
      type="primary"
      @click="emit('pick')"
    >
      {{ buttonText }}
    </el-button>
    <button
      class="image-picker-trigger__box"
      :style="{ height: `${size}px`, width: `${size}px` }"
      type="button"
      @click="emit('pick')"
    >
      <img
        v-if="resolvedPreviewUrl"
        :alt="'preview'"
        class="image-picker-trigger__preview"
        :src="resolvedPreviewUrl"
      />
      <div v-else class="image-picker-trigger__placeholder">
        <el-icon :size="28"><Picture /></el-icon>
        <span>暂无图片</span>
      </div>
      <div v-if="resolvedPreviewUrl" class="image-picker-trigger__mask">
        <el-icon :size="18"><Plus /></el-icon>
        <span>更换</span>
      </div>
    </button>
    <p v-if="hint" class="image-picker-trigger__hint">{{ hint }}</p>
  </div>
</template>

<style scoped lang="scss">
.image-picker-trigger {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.image-picker-trigger__btn {
  min-width: 96px;
  border-radius: var(--el-border-radius-base);
}

.image-picker-trigger__box {
  position: relative;
  padding: 0;
  border: 1px dashed hsl(var(--border));
  border-radius: 8px;
  background: hsl(var(--muted) / 18%);
  cursor: pointer;
  overflow: hidden;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.image-picker-trigger__box:hover {
  border-color: hsl(var(--primary) / 55%);
  box-shadow: 0 0 0 2px hsl(var(--primary) / 10%);
}

.image-picker-trigger__preview {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-picker-trigger__placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  height: 100%;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.image-picker-trigger__mask {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  background: hsl(var(--foreground) / 45%);
  color: #fff;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.image-picker-trigger__box:hover .image-picker-trigger__mask {
  opacity: 1;
}

.image-picker-trigger__hint {
  margin: 0;
  font-size: 12px;
  line-height: 1.5;
  color: hsl(var(--muted-foreground));
}
</style>
