<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="小票打印" description="云打印机与小票模板配置。">
    <el-alert type="info" :closable="false" title="打印机对接接口待接入；当前为占位页面。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无打印机配置" /></el-card>
  </Page>
</template>
