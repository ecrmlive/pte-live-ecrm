<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue';

withDefaults(
  defineProps<{
    disabled?: boolean;
    plain?: boolean;
    showIcon?: boolean;
    size?: 'default' | 'large' | 'small';
    text?: string;
  }>(),
  {
    disabled: false,
    plain: false,
    showIcon: true,
    size: 'default',
    text: '选择会员',
  },
);

const emit = defineEmits<{ click: [] }>();
</script>

<template>
  <ElButton
    :disabled="disabled"
    :icon="showIcon ? Plus : undefined"
    :plain="plain"
    :size="size"
    type="primary"
    @click="emit('click')"
  >
    {{ text }}
  </ElButton>
</template>
