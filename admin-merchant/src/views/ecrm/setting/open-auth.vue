<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="开放授权" description="第三方应用授权与回调配置。">
    <el-alert type="info" :closable="false" title="开放授权（OAuth）能力待平台统一网关接入；当前为占位页面。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无授权应用" /></el-card>
  </Page>
</template>
