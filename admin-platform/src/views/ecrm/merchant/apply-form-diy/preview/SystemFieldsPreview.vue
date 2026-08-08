<script setup lang="ts">
import { computed } from 'vue';

import { SYSTEM_FIELDS } from '../system-fields';
import type { SystemField } from '../types';
import PreviewUploadSlots from './PreviewUploadSlots.vue';

const props = defineProps<{
  fields?: SystemField[];
}>();

const list = computed(() => props.fields ?? SYSTEM_FIELDS);
</script>

<template>
  <div
    v-for="item in list"
    :key="item.title"
    class="sys-field"
    :class="{
      'sys-field--block': item.kind === 'image' || item.kind === 'textarea',
    }"
  >
    <template v-if="item.kind === 'image'">
      <div class="sys-field__image-title">
        <span v-if="item.required" class="req">*</span>
        {{ item.title }}
      </div>
      <PreviewUploadSlots :count="item.imageSlots || 1" />
    </template>
    <template v-else-if="item.kind === 'textarea'">
      <div class="sys-field__textarea-title">
        <span v-if="item.required" class="req">*</span>
        {{ item.title }}
      </div>
      <div class="sys-field__textarea">
        <span class="place">{{ item.placeholder }}</span>
      </div>
    </template>
    <template v-else>
      <span class="sys-field__label">
        <span v-if="item.required" class="req">*</span>
        {{ item.title }}
      </span>
      <div class="sys-field__value">
        <span class="place">{{
          item.kind === 'select' ? '请选择' : item.placeholder
        }}</span>
        <span v-if="item.kind === 'select'" class="arrow">›</span>
      </div>
    </template>
  </div>
</template>

<style scoped lang="scss">
.sys-field {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 11px 10px 11px 12px;
  border-bottom: 1px solid #eee;
  background: #fff;
  font-size: 15px;
  color: #333;
}

.sys-field--block {
  display: block;
  padding: 0 12px 15px;
}

.sys-field__label {
  width: 95px;
  flex-shrink: 0;
  line-height: 1.4;
}

.sys-field__value {
  width: 250px;
  flex-shrink: 0;
  text-align: right;
}

.place {
  color: #ccc;
  font-weight: 400;
}

.arrow {
  margin-left: 10px;
  color: #999;
  font-size: 14px;
  line-height: 1;
}

.sys-field__image-title,
.sys-field__textarea-title {
  width: 100%;
  padding-top: 14px;
  color: #282828;
  font-size: 15px;
}

.sys-field__textarea-title {
  padding-bottom: 8px;
}

.sys-field__textarea {
  min-height: 100px;
  padding: 8px 10px;
  border: 1px solid #eee;
  border-radius: 4px;
  background: #fafafa;
  line-height: 1.5;
  font-size: 14px;
}

.req {
  margin-right: 2px;
  color: #f56c6c;
}
</style>
