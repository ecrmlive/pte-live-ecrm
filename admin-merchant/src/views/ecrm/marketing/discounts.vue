<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="限时折扣" description="店铺限时折扣活动。">
    <el-alert type="info" :closable="false" title="限时折扣活动接口待接入；当前保留列表占位。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无折扣活动" /></el-card>
  </Page>
</template>
