<script setup lang="ts">
import MarketingDecorPage from './MarketingDecorPage.vue';
</script>

<template>
  <MarketingDecorPage
    title="活动报名"
    create-label="新增报名"
    decor-type="application"
    read-code="marketing.application.read"
    manage-code="marketing.application.manage"
    description="平台活动报名配置，数据落在 qixi_crm_a_marketing_decor；不含密钥。"
  />
</template>
