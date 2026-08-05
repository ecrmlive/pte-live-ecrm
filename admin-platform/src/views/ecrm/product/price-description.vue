<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  fetchPriceDescriptions,
  savePriceDescriptions,
  type ProductCacheListItem,
} from '#/api/core/platform-product-cache';

const loading = ref(false);
const saving = ref(false);
const dialog = ref(false);
const rows = ref<ProductCacheListItem[]>([]);
const editingIndex = ref<number>();
const form = reactive<ProductCacheListItem>({ id: '', name: '', enabled: true, remark: '' });

function resetForm() {
  editingIndex.value = undefined;
  Object.assign(form, { id: '', name: '', enabled: true, remark: '' });
}

function open(row?: ProductCacheListItem, index?: number) {
  resetForm();
  if (row !== undefined && index !== undefined) {
    editingIndex.value = index;
    Object.assign(form, row);
  }
  dialog.value = true;
}

async function load() {
  loading.value = true;
  try {
    const result = await fetchPriceDescriptions();
    rows.value = result.list || [];
  } finally {
    loading.value = false;
  }
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请填写名称');
    return;
  }
  const next = rows.value.map((item) => ({ ...item }));
  const payload: ProductCacheListItem = {
    id: form.id.trim() || name,
    name,
    enabled: form.enabled,
    remark: form.remark.trim(),
  };
  if (editingIndex.value === undefined) next.push(payload);
  else next[editingIndex.value] = payload;

  saving.value = true;
  try {
    const result = await savePriceDescriptions(next);
    rows.value = result.list || [];
    dialog.value = false;
    ElMessage.success('已保存');
  } finally {
    saving.value = false;
  }
}

async function remove(index: number, name: string) {
  try {
    await ElMessageBox.confirm(`确认删除“${name}”？`, '删除确认', { type: 'warning' });
    const next = rows.value.filter((_, i) => i !== index);
    const result = await savePriceDescriptions(next);
    rows.value = result.list || [];
    ElMessage.success('已删除');
  } catch {
    // 用户取消或请求失败由统一层反馈。
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="价格说明" description="基于 setting_cache 的 stub 配置，维护商品详情页可引用的价格说明条目。">
    <template #extra><el-button type="primary" @click="open()">新增说明</el-button></template>
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="id">
        <el-table-column label="标识" prop="id" min-width="120" />
        <el-table-column label="名称" prop="name" min-width="160" />
        <el-table-column label="备注" prop="remark" min-width="220" show-overflow-tooltip />
        <el-table-column label="启用" width="90">
          <template #default="{ row }"><el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row, $index }">
            <el-button link type="primary" @click="open(row, $index)">编辑</el-button>
            <el-button link type="danger" @click="remove($index, row.name)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialog" :title="`${editingIndex === undefined ? '新增' : '编辑'}价格说明`" width="520px" destroy-on-close>
      <el-form label-width="84px">
        <el-form-item label="名称" required><el-input v-model="form.name" maxlength="64" show-word-limit /></el-form-item>
        <el-form-item label="标识"><el-input v-model="form.id" maxlength="64" placeholder="留空则使用名称" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" maxlength="255" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.enabled" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialog = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
