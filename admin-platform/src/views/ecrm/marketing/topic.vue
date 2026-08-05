<script setup lang="ts">
import MarketingDecorPage from './MarketingDecorPage.vue';
</script>

<template>
  <MarketingDecorPage
    title="活动专题"
    create-label="新增专题"
    decor-type="topic"
    read-code="marketing.topic.read"
    manage-code="marketing.topic.manage"
    description="平台营销专题配置，数据落在 qixi_crm_a_marketing_decor；不含密钥。"
  />
</template>
