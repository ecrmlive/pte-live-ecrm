<script setup lang="ts">
import { ElDialog, ElInput, ElButton, ElMessage } from 'element-plus';
import { ref, watch } from 'vue';
const props = defineProps<{ modelValue?: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [boolean]; confirm: [string] }>();
const url = ref('');
const open = ref(false);
watch(() => props.modelValue, (v) => (open.value = !!v));
watch(open, (v) => emit('update:modelValue', v));
function ok() {
  if (!url.value) { ElMessage.warning('请输入图片 URL'); return; }
  emit('confirm', url.value);
  open.value = false;
}
</script>
<template>
  <ElDialog v-model="open" title="选择图片" width="480px">
    <ElInput v-model="url" placeholder="图片 URL" />
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="ok">确定</ElButton>
    </template>
  </ElDialog>
</template>
