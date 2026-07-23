<script lang="ts" setup>
import { ElAlert } from 'element-plus';

defineOptions({ name: 'ServiceError' });
</script>

<template>
  <main class="service-error">
    <ElAlert
      :closable="false"
      show-icon
      title="内部错误，请稍后重试"
      type="error"
    />
  </main>
</template>

<style scoped>
.service-error {
  width: min(100%, 620px);
  margin: 20vh auto 0;
}
</style>
