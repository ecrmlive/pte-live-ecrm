<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="分账申请" description="微信支付分账与特约商户进件。">
    <el-alert type="info" :closable="false" title="分账申请需对接支付渠道进件接口；当前为占位页面，请通过平台侧完成特约商户审核。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无分账申请记录" /></el-card>
  </Page>
</template>
