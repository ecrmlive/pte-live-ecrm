<script setup lang="ts">
import { ElDialog, ElButton, ElMessage } from 'element-plus';
import { ref, watch } from 'vue';
const props = defineProps<{ modelValue?: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [boolean]; confirm: [unknown[]] }>();
const open = ref(false);
watch(() => props.modelValue, (v) => (open.value = !!v));
watch(open, (v) => emit('update:modelValue', v));
</script>
<template>
  <ElDialog v-model="open" title="选择商品" width="520px">
    <p>平台端请先在接口侧配置商品选择；当前可保存空列表。</p>
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="emit('confirm', []); open = false; ElMessage.success('已确认')">确定</ElButton>
    </template>
  </ElDialog>
</template>
