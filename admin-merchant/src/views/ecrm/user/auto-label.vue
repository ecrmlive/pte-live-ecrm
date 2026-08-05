<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { listStoreAutoLabelRulesApi, saveStoreAutoLabelRulesApi, type StoreAutoLabelRule } from '#/api/core/merchant-user-label';

const rows = ref<StoreAutoLabelRule[]>([]);
const loading = ref(false);
const saving = ref(false);

async function load() {
  loading.value = true;
  try {
    rows.value = (await listStoreAutoLabelRulesApi()).list ?? [];
  } finally {
    loading.value = false;
  }
}

function addRow() {
  rows.value.push({ rule_id: rows.value.length + 1, name: '新规则', rule_type: 'order_count', status: 0 });
}

async function save() {
  saving.value = true;
  try {
    const result = await saveStoreAutoLabelRulesApi({ list: rows.value });
    rows.value = result.list ?? rows.value;
    ElMessage.success('自动标签规则已保存（规则引擎待后续接入）');
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="自动标签" description="配置自动打标规则占位；执行引擎待后续任务接入。">
    <template #extra><el-button @click="addRow">新增规则</el-button><el-button type="primary" :loading="saving" @click="save">保存</el-button></template>
    <el-alert class="mb-4" type="info" :closable="false" title="当前仅保存规则配置，不会自动执行打标。" />
    <el-card v-loading="loading" shadow="never">
      <el-table :data="rows">
        <el-table-column label="规则名称" min-width="180"><template #default="{ row }"><el-input v-model="row.name" /></template></el-table-column>
        <el-table-column label="规则类型" width="180"><template #default="{ row }"><el-select v-model="row.rule_type"><el-option label="下单次数" value="order_count" /><el-option label="累计消费" value="total_pay" /></el-select></template></el-table-column>
        <el-table-column label="启用" width="100"><template #default="{ row }"><el-switch v-model="row.status" :active-value="1" :inactive-value="0" /></template></el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
