<script setup lang="ts">
import { ElDialog, ElButton } from 'element-plus';
import { ref, watch } from 'vue';
const props = defineProps<{ modelValue?: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [boolean]; confirm: [unknown[]] }>();
const open = ref(false);
watch(() => props.modelValue, (v) => (open.value = !!v));
watch(open, (v) => emit('update:modelValue', v));
</script>
<template>
  <ElDialog v-model="open" title="选择店铺" width="480px">
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="emit('confirm', []); open = false">确定</ElButton>
    </template>
  </ElDialog>
</template>
