<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page
    title="付费会员（店铺）"
    description="店铺付费会员开关与权益说明；套餐主数据由平台统一维护。"
  >
    <el-alert
      :closable="false"
      show-icon
      type="info"
      title="SVIP 套餐、权益与购买记录请在平台后台「用户 / 付费会员」维护。本店仅可查看经营侧说明，不在此保存支付密钥。"
    />
  </Page>
</template>
