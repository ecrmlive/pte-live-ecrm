<script setup lang="ts">
import ConfigItemPage from './ConfigItemPage.vue';
</script>

<template>
  <ConfigItemPage
    title="组合数据"
    create-label="新增组合"
    item-type="group_data"
    read-code="maintain.group_data"
    manage-code="maintain.group_data.manage"
    description="组合数据配置落在 qixi_crm_a_config_item，不含密钥。"
  />
</template>
