<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="卡密库" description="虚拟商品卡密管理。">
    <el-alert type="info" :closable="false" title="卡密库功能待 api-merchant 卡密域接口上线；当前仅保留菜单与页面占位。" />
    <el-card class="mt-4" shadow="never">
      <el-empty description="暂无卡密数据" />
    </el-card>
  </Page>
</template>
