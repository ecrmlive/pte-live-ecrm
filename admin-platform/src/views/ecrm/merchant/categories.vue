<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantCategory,
  deleteMerchantCategory,
  fetchMerchantCategories,
  updateMerchantCategory,
  type MerchantCategoryRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';

const rows = ref<MerchantCategoryRow[]>([]);
const loading = ref(false);
const open = ref(false);
const canManage = ref(false);
const editing = ref<MerchantCategoryRow>();
const form = reactive({ category_name: '', commission_rate: 0 });

async function load() {
  loading.value = true;
  try {
    rows.value = (await fetchMerchantCategories()).list || [];
  } finally {
    loading.value = false;
  }
}

function add() {
  editing.value = undefined;
  Object.assign(form, { category_name: '', commission_rate: 0 });
  open.value = true;
}

function edit(row: MerchantCategoryRow) {
  editing.value = row;
  Object.assign(form, { category_name: row.category_name, commission_rate: Number(row.commission_rate) });
  open.value = true;
}

async function save() {
  if (!form.category_name.trim() || form.commission_rate < 0 || form.commission_rate > 100) {
    ElMessage.warning('请填写分类名称，并将佣金比例限制在 0% 至 100%');
    return;
  }
  const input = { category_name: form.category_name.trim(), commission_rate: form.commission_rate };
  if (editing.value) await updateMerchantCategory(editing.value.merchant_category_id, input);
  else await createMerchantCategory(input);
  open.value = false;
  ElMessage.success('商户分类已保存');
  await load();
}

async function remove(row: MerchantCategoryRow) {
  try {
    await ElMessageBox.confirm(`删除商户分类“${row.category_name}”后不可恢复，是否继续？`, '删除商户分类', { type: 'warning' });
    await deleteMerchantCategory(row.merchant_category_id);
    ElMessage.success('商户分类已删除');
    await load();
  } catch {
    // 取消或请求失败由统一请求层反馈。
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('merchant.category.manage');
});
</script>

<template>
  <Page title="商户分类" description="维护商户入驻分类及平台佣金比例，不与商品分类混用。">
    <template #extra><el-button v-if="canManage" type="primary" @click="add">新增分类</el-button></template>
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="merchant_category_id">
        <el-table-column label="ID" prop="merchant_category_id" width="90" />
        <el-table-column label="商户分类" min-width="240" prop="category_name" />
        <el-table-column label="平台佣金比例" width="170"><template #default="{ row }">{{ Number(row.commission_rate).toFixed(2) }}%</template></el-table-column>
        <el-table-column label="操作" width="150"><template #default="{ row }"><el-button v-if="canManage" link type="primary" @click="edit(row)">编辑</el-button><el-button v-if="canManage" link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
    </el-card>
    <el-dialog v-model="open" :title="editing ? '编辑商户分类' : '新增商户分类'" width="480px" destroy-on-close>
      <el-form label-width="108px"><el-form-item label="分类名称" required><el-input v-model="form.category_name" maxlength="128" /></el-form-item><el-form-item label="佣金比例" required><el-input-number v-model="form.commission_rate" :max="100" :min="0" :precision="2" class="w-full" /><span class="ml-2">%</span></el-form-item></el-form>
      <template #footer><el-button @click="open = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
