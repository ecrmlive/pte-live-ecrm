<script setup lang="ts">
defineProps<{
  modelValue: boolean;
  size?: string;
  title: string;
}>();

defineEmits<{
  'update:modelValue': [value: boolean];
}>();
</script>

<template>
  <el-drawer
    :model-value="modelValue"
    :size="size || '560px'"
    :title="title"
    destroy-on-close
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div class="ecrm-form-drawer flex flex-col gap-4">
      <slot />
    </div>
    <template v-if="$slots.footer" #footer>
      <div class="flex justify-end gap-2">
        <slot name="footer" />
      </div>
    </template>
  </el-drawer>
</template>
