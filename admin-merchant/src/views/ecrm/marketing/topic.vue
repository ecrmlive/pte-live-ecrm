<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="专题活动" description="店铺专题页与活动分组。">
    <el-alert type="info" :closable="false" title="专题活动配置待 diy/营销域接口上线；当前为占位页面。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无专题活动" /></el-card>
  </Page>
</template>
