<script setup lang="ts">
defineProps<{
  title: string;
  required?: boolean;
  value: string;
  showArrow?: boolean;
}>();
</script>

<template>
  <div class="preview-row">
    <span class="preview-row__label">
      <span v-if="required" class="req">*</span>
      {{ title || '未命名' }}
    </span>
    <div class="preview-row__value">
      <span class="place">{{ value }}</span>
      <span v-if="showArrow" class="arrow">›</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
.preview-row {
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

.preview-row__label {
  width: 95px;
  flex-shrink: 0;
  line-height: 1.4;
}

.preview-row__value {
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

.req {
  margin-right: 2px;
  color: #f56c6c;
}
</style>
