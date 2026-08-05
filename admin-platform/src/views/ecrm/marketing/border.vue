<script setup lang="ts">
import MarketingDecorPage from './MarketingDecorPage.vue';
</script>

<template>
  <MarketingDecorPage
    title="活动边框"
    create-label="新增边框"
    decor-type="border"
    read-code="marketing.border.read"
    manage-code="marketing.border.manage"
    description="平台活动边框图配置，数据落在 qixi_crm_a_marketing_decor；不含密钥。"
  />
</template>
