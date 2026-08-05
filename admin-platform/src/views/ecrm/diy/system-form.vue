<script setup lang="ts">
import ConfigItemPage from '../maintain/ConfigItemPage.vue';
</script>

<template>
  <ConfigItemPage
    title="系统表单"
    create-label="新增表单"
    item-type="system_form"
    read-code="operations.system_form"
    manage-code="operations.system_form.manage"
    :roles="['platform', 'operations']"
    description="系统表单登记落在 qixi_crm_a_config_item；完整 form_builder 渲染后续扩展。"
  />
</template>
