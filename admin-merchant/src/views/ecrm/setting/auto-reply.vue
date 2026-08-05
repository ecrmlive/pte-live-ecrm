<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="自动回复" description="客服关键词自动回复配置。">
    <el-alert type="info" :closable="false" title="自动回复规则待 IM/客服域接口上线；当前为占位页面。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无自动回复规则" /></el-card>
  </Page>
</template>
