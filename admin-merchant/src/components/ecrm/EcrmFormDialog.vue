<script setup lang="ts">
defineProps<{
  modelValue: boolean;
  title: string;
  width?: string;
}>();

defineEmits<{
  'update:modelValue': [value: boolean];
}>();
</script>

<template>
  <el-dialog
    :model-value="modelValue"
    :title="title"
    :width="width || '640px'"
    destroy-on-close
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div class="ecrm-form-dialog">
      <slot />
    </div>
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </el-dialog>
</template>
