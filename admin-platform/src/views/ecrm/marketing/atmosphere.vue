<script setup lang="ts">
import MarketingDecorPage from './MarketingDecorPage.vue';
</script>

<template>
  <MarketingDecorPage
    title="活动氛围"
    create-label="新增氛围"
    decor-type="atmosphere"
    read-code="marketing.atmosphere.read"
    manage-code="marketing.atmosphere.manage"
    description="平台活动氛围图配置，数据落在 qixi_crm_a_marketing_decor；不含密钥。"
  />
</template>
