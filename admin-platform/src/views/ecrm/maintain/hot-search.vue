<script setup lang="ts">
import ConfigItemPage from './ConfigItemPage.vue';
</script>

<template>
  <ConfigItemPage
    title="热门搜索"
    create-label="新增标签"
    name-label="关键词"
    item-type="hot_search"
    read-code="maintain.hot_search"
    manage-code="maintain.hot_search.manage"
    description="热门搜索标签落在 qixi_crm_a_config_item，不含密钥。"
  />
</template>
