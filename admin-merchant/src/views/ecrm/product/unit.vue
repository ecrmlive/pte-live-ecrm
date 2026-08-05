<script setup lang="ts">
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { onMounted, ref } from 'vue';
import { listProductUnitsApi, saveProductUnitsApi, type ProductUnit } from '#/api/core/merchant-product-unit';

const rows = ref<ProductUnit[]>([]);
const loading = ref(false);
const saving = ref(false);

async function load() {
  loading.value = true;
  try {
    rows.value = (await listProductUnitsApi()).list ?? [];
  } finally {
    loading.value = false;
  }
}

function addRow() {
  rows.value.push({ unit_id: rows.value.length + 1, name: '', sort: rows.value.length });
}

function removeRow(index: number) {
  rows.value.splice(index, 1);
}

async function save() {
  if (rows.value.some((item) => !item.name.trim())) {
    ElMessage.warning('请填写全部单位名称');
    return;
  }
  saving.value = true;
  try {
    const result = await saveProductUnitsApi({ list: rows.value.map((item, index) => ({ ...item, name: item.name.trim(), sort: index })) });
    rows.value = result.list ?? rows.value;
    ElMessage.success('商品单位已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="商品单位" description="维护本店商品计量单位，新建商品时可选择。">
    <template #extra><el-button type="primary" @click="addRow">新增单位</el-button><el-button type="success" :loading="saving" @click="save">保存</el-button></template>
    <el-card v-loading="loading" shadow="never">
      <el-table :data="rows" row-key="unit_id">
        <el-table-column label="单位名称" min-width="220"><template #default="{ row }"><el-input v-model="row.name" maxlength="16" /></template></el-table-column>
        <el-table-column label="排序" width="120"><template #default="{ $index }">{{ $index }}</template></el-table-column>
        <el-table-column label="操作" width="100"><template #default="{ $index }"><el-button link type="danger" @click="removeRow($index)">删除</el-button></template></el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
