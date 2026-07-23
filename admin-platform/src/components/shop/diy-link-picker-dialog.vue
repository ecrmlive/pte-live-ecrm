<script setup lang="ts">
import { ElDialog, ElInput, ElButton } from 'element-plus';
import { ref, watch } from 'vue';
const props = defineProps<{ modelValue?: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [boolean]; confirm: [Record<string, string>] }>();
const link = ref('');
const open = ref(false);
watch(() => props.modelValue, (v) => (open.value = !!v));
watch(open, (v) => emit('update:modelValue', v));
</script>
<template>
  <ElDialog v-model="open" title="选择链接" width="480px">
    <ElInput v-model="link" placeholder="页面路径，如 /pages/index/index" />
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="emit('confirm', { path: link }); open = false">确定</ElButton>
    </template>
  </ElDialog>
</template>
