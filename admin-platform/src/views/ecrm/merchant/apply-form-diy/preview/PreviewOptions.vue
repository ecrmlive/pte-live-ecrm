<script setup lang="ts">
import type { ApplyFormField } from '../types';

defineProps<{ field: ApplyFormField }>();
</script>

<template>
  <div class="preview-options">
    <div class="preview-options__title">
      <span v-if="field.required" class="req">*</span>
      {{ field.title || '未命名' }}
    </div>
    <div class="preview-options__list">
      <div
        v-for="(opt, index) in field.options || []"
        :key="`${opt}-${index}`"
        class="preview-options__item"
      >
        <span
          class="preview-options__ctrl"
          :class="{
            'is-radio': field.type === 'radio',
            on: field.type === 'radio' && index === 0,
          }"
        >
          <span
            v-if="field.type === 'radio' && index === 0"
            class="preview-options__check"
          >✓</span>
        </span>
        <span class="preview-options__label">{{ opt }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.preview-options {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  width: 100%;
  padding: 8px 5px 8px 12px;
  border-bottom: 1px solid #eee;
  background: #fff;
  font-size: 15px;
  color: #333;
}

.preview-options__title {
  width: 95px;
  flex-shrink: 0;
  line-height: 1.4;
}

.preview-options__list {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  align-items: center;
  max-width: 250px;
}

.preview-options__item {
  display: flex;
  align-items: center;
  margin: 0 10px;
  padding: 5px 0;
  color: #282828;
  font-size: 13px;
}

.preview-options__ctrl {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 17px;
  height: 17px;
  margin-right: 7px;
  border: 1px solid #ccc;
  background: #fff;
  flex-shrink: 0;
}

.preview-options__ctrl.is-radio {
  width: 19px;
  height: 19px;
  border-radius: 50%;
}

.preview-options__ctrl.on {
  border-color: #e93323;
  background: #e93323;
  color: #fff;
}

.preview-options__check {
  font-size: 11px;
  line-height: 1;
  transform: scale(0.9);
}

.preview-options__label {
  line-height: 1.3;
}

.req {
  margin-right: 2px;
  color: #f56c6c;
}
</style>
