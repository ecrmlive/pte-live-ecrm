<script setup lang="ts">
import { Page } from '@vben/common-ui';
</script>

<template>
  <Page
    title="一号通 / 短信配置"
    description="店铺侧仅展示入口；真实短信与电子面单凭据由平台受控配置，本页不保存密钥。"
  >
    <el-alert
      :closable="false"
      show-icon
      type="info"
      title="当前环境使用平台 stub 短信通道。如需开通一号通能力，请联系平台管理员在统一后台配置，勿在店铺后台粘贴 API Key。"
    />
  </Page>
</template>
