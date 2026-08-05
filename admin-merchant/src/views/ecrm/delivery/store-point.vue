<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import { createStorePickupPointApi, deleteStorePickupPointApi, listStorePickupPointsApi, updateStorePickupPointApi, type StorePickupPoint } from '#/api/core/merchant-pickup-point';

const rows = ref<StorePickupPoint[]>([]);
const loading = ref(false);
const open = ref(false);
const saving = ref(false);
const editID = ref<number>();
const form = ref<Omit<StorePickupPoint, 'id'>>({ contact_name: '', mobile: '', region_code: '', detail: '', is_default: 0 });

async function load() {
  loading.value = true;
  try {
    rows.value = (await listStorePickupPointsApi()).list ?? [];
  } finally {
    loading.value = false;
  }
}

function reset() {
  editID.value = undefined;
  form.value = { contact_name: '', mobile: '', region_code: '', detail: '', is_default: 0 };
}

function create() {
  reset();
  open.value = true;
}

function edit(row: StorePickupPoint) {
  editID.value = row.id;
  form.value = { contact_name: row.contact_name, mobile: row.mobile, region_code: row.region_code, detail: row.detail, is_default: row.is_default };
  open.value = true;
}

async function save() {
  saving.value = true;
  try {
    if (editID.value) await updateStorePickupPointApi(editID.value, form.value);
    else await createStorePickupPointApi(form.value);
    open.value = false;
    await load();
    ElMessage.success('自提点已保存');
  } finally {
    saving.value = false;
  }
}

async function remove(row: StorePickupPoint) {
  try {
    await ElMessageBox.confirm(`确定删除自提点“${row.contact_name}”吗？`, '删除确认', { type: 'warning' });
    await deleteStorePickupPointApi(row.id);
    await load();
    ElMessage.success('自提点已删除');
  } catch {}
}

onMounted(() => void load());
</script>

<template>
  <Page title="自提点管理" description="维护本店到店自提地址。">
    <template #extra><el-button type="primary" @click="create">新增自提点</el-button></template>
    <el-card v-loading="loading" shadow="never">
      <el-table :data="rows">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="contact_name" label="联系人" width="120" />
        <el-table-column prop="mobile" label="手机号" width="140" />
        <el-table-column prop="detail" label="地址" min-width="220" />
        <el-table-column label="默认" width="80"><template #default="{ row }">{{ row.is_default ? '是' : '否' }}</template></el-table-column>
        <el-table-column label="操作" width="130" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
    </el-card>
    <el-dialog v-model="open" :title="editID ? '编辑自提点' : '新增自提点'" width="520px">
      <el-form label-width="88px">
        <el-form-item label="联系人" required><el-input v-model="form.contact_name" /></el-form-item>
        <el-form-item label="手机号" required><el-input v-model="form.mobile" /></el-form-item>
        <el-form-item label="区划代码"><el-input v-model="form.region_code" /></el-form-item>
        <el-form-item label="详细地址" required><el-input v-model="form.detail" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="默认"><el-switch v-model="form.is_default" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="open = false">取消</el-button><el-button type="primary" :loading="saving" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
