<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page title="平台公告" description="查看平台向本店下发的系统公告。">
    <el-alert
      type="info"
      :closable="false"
      title="公告由平台统一下发；本店暂无独立公告 CRUD。"
      description="如需发布面向消费者的商城公告，请使用平台后台「内容管理 → 公告管理」。本页预留菜单入口，后续可接入平台公告已读列表。"
    />
  </Page>
</template>
