<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { listStoreUserLabelsApi, saveStoreUserLabelsApi, type StoreUserLabel } from '#/api/core/merchant-user-label';

const rows = ref<StoreUserLabel[]>([]);
const loading = ref(false);
const saving = ref(false);

async function load() {
  loading.value = true;
  try {
    rows.value = (await listStoreUserLabelsApi()).list ?? [];
  } finally {
    loading.value = false;
  }
}

function addRow() {
  rows.value.push({ label_id: rows.value.length + 1, name: '', sort: rows.value.length, status: 1 });
}

async function save() {
  saving.value = true;
  try {
    const result = await saveStoreUserLabelsApi({ list: rows.value.map((item, index) => ({ ...item, name: item.name.trim(), sort: index })) });
    rows.value = result.list ?? rows.value;
    ElMessage.success('用户标签已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="用户标签" description="本店用户运营标签，保存在店铺配置中。">
    <template #extra><el-button @click="addRow">新增标签</el-button><el-button type="primary" :loading="saving" @click="save">保存</el-button></template>
    <el-card v-loading="loading" shadow="never">
      <el-table :data="rows">
        <el-table-column label="标签名称" min-width="200"><template #default="{ row }"><el-input v-model="row.name" maxlength="32" /></template></el-table-column>
        <el-table-column label="排序" width="100"><template #default="{ $index }">{{ $index }}</template></el-table-column>
        <el-table-column label="启用" width="100"><template #default="{ row }"><el-switch v-model="row.status" :active-value="1" :inactive-value="0" /></template></el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
