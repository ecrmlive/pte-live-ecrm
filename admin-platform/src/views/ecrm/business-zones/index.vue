<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createBusinessZone,
  deleteBusinessZone,
  fetchBusinessZones,
  updateBusinessZone,
  type BusinessZoneRow,
} from '#/api/core/ecrm';

const rows = ref<BusinessZoneRow[]>([]);
const total = ref(0);
const loading = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const query = reactive({ keyword: '', status: undefined as number | undefined, page: 1, limit: 20 });
const form = reactive({ pid: 0, name: '', circle_agent_id: 0, commission_type: 0, commission_rate: 0, remark: '', sort: 0, status: 1, type: 0, role_id: 0 });

const statusText = (value: number) => value === 1 ? '启用' : '禁用';
const typeText = (value: number) => value === 1 ? '商户型商圈' : '区域商圈';

async function load() {
  loading.value = true;
  try { const result = await fetchBusinessZones({ ...query, keyword: query.keyword.trim() || undefined }); rows.value = result.list; total.value = result.total; }
  finally { loading.value = false; }
}
function resetForm() { Object.assign(form, { pid: 0, name: '', circle_agent_id: 0, commission_type: 0, commission_rate: 0, remark: '', sort: 0, status: 1, type: 0, role_id: 0 }); }
function add() { editingID.value = undefined; resetForm(); dialogOpen.value = true; }
function edit(row: BusinessZoneRow) { editingID.value = row.circle_id; Object.assign(form, { pid: row.pid, name: row.name, circle_agent_id: row.circle_agent_id, commission_type: row.commission_type, commission_rate: row.commission_rate, remark: row.remark, sort: row.sort, status: row.status, type: row.type, role_id: row.role_id }); dialogOpen.value = true; }
async function save() { if (!form.name.trim()) { ElMessage.warning('请填写商圈名称'); return; } if (editingID.value) await updateBusinessZone(editingID.value, form); else await createBusinessZone(form); ElMessage.success('保存成功'); dialogOpen.value = false; await load(); }
async function remove(row: BusinessZoneRow) { try { await ElMessageBox.confirm(`删除“${row.name}”后不可恢复，是否继续？`, '删除区域'); await deleteBusinessZone(row.circle_id); ElMessage.success('已删除'); await load(); } catch {} }
function search() { query.page = 1; void load(); }
function reset() { query.keyword = ''; query.status = undefined; query.page = 1; void load(); }
onMounted(load);
</script>

<template>
  <Page title="区域列表" description="维护区域商圈层级、启停状态和代理提成规则。区域一旦建立父级不可直接变更，避免影响已绑定商户的数据范围。">
    <el-card shadow="never">
      <el-form class="grid gap-x-4 md:grid-cols-3" label-width="72px" @submit.prevent="search">
        <el-form-item label="区域名称"><el-input v-model="query.keyword" clearable placeholder="商圈名称" @keyup.enter="search" /></el-form-item>
        <el-form-item label="状态"><el-select v-model="query.status" clearable class="w-full" placeholder="全部"><el-option label="启用" :value="1" /><el-option label="禁用" :value="0" /></el-select></el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button><el-button type="success" @click="add">新建区域</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" class="mt-2" border>
        <el-table-column prop="circle_id" label="ID" width="72" /><el-table-column prop="name" label="区域/商圈" min-width="150" />
        <el-table-column prop="path" label="层级路径" min-width="150" /><el-table-column label="类型" width="120"><template #default="{ row }">{{ typeText(row.type) }}</template></el-table-column>
        <el-table-column label="提成" width="150"><template #default="{ row }">{{ row.commission_type === 1 ? `${row.commission_rate}%（独立）` : '平台默认' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="150" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination v-model:current-page="query.page" v-model:page-size="query.limit" :page-sizes="[10,20,50,100]" layout="total, sizes, prev, pager, next" :total="total" @current-change="load" @size-change="() => { query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑区域' : '新建区域'" width="620px">
      <el-form label-width="110px"><el-form-item label="上级区域ID"><el-input-number v-model="form.pid" :min="0" /></el-form-item><el-form-item label="区域名称" required><el-input v-model="form.name" maxlength="64" /></el-form-item><el-form-item label="类型"><el-radio-group v-model="form.type"><el-radio :value="0">区域商圈</el-radio><el-radio :value="1">商户型商圈</el-radio></el-radio-group></el-form-item><el-form-item label="提成规则"><el-radio-group v-model="form.commission_type"><el-radio :value="0">平台默认</el-radio><el-radio :value="1">独立比例</el-radio></el-radio-group><el-input-number v-if="form.commission_type === 1" v-model="form.commission_rate" class="ml-2" :max="100" :min="0" :precision="2" /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" /></el-form-item><el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item><el-form-item label="说明"><el-input v-model="form.remark" type="textarea" /></el-form-item></el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
