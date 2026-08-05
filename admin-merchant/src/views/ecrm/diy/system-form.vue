<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="系统表单" description="自定义表单收集与数据管理。">
    <el-alert type="info" :closable="false" title="系统表单能力待 diy 域扩展；当前为占位页面。" />
    <el-card class="mt-4" shadow="never"><el-empty description="暂无系统表单" /></el-card>
  </Page>
</template>
