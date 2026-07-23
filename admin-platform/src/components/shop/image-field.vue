<script setup lang="ts">
import { ElInput } from 'element-plus';
const model = defineModel<string>({ default: '' });
</script>
<template>
  <ElInput v-model="model" placeholder="图片 URL" />
</template>
