<script setup lang="ts">
import ConfigItemPage from './ConfigItemPage.vue';
</script>

<template>
  <ConfigItemPage
    title="数据备份"
    create-label="登记备份"
    item-type="backup"
    read-code="maintain.backup"
    manage-code="maintain.backup.manage"
    description="备份登记落在 qixi_crm_a_config_item；真实备份任务仍由运维侧执行，此处仅登记元数据。"
  />
</template>
