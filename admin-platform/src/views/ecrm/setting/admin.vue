<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createPlatformAdmin,
  deletePlatformAdmin,
  fetchBusinessZoneAgentOptions,
  fetchPlatformAdmins,
  updatePlatformAdmin,
  type PlatformAdminRow,
} from '#/api/core/ecrm';

const rows = ref<PlatformAdminRow[]>([]);
const total = ref(0);
const loading = ref(false);
const dialogOpen = ref(false);
const editing = ref<PlatformAdminRow>();
const query = reactive({ page: 1, limit: 20 });
const agentOptions = ref<Array<{ circle_agent_id: number; name: string; type: 0 | 1 }>>([]);
const form = reactive({ account: '', password: '', real_name: '', phone: '', roles: '', status: 1, merchant_ids: '', region_ids: '', service_store_ids: '', circle_agent_id: 0 });
const title = computed(() => editing.value ? '编辑管理员' : '新增管理员');
const isRegionAccount = computed(() => { const roles = form.roles.split(',').map((value) => value.trim()); return roles.includes('region') && !roles.includes('platform'); });

async function load() { loading.value = true; try { const result = await fetchPlatformAdmins(query); rows.value = result.list; total.value = result.total; } finally { loading.value = false; } }
async function loadAgentOptions() { agentOptions.value = (await fetchBusinessZoneAgentOptions()).list || []; }
function reset() { Object.assign(form, { account: '', password: '', real_name: '', phone: '', roles: '', status: 1, merchant_ids: '', region_ids: '', service_store_ids: '', circle_agent_id: 0 }); }
function add() { editing.value = undefined; reset(); dialogOpen.value = true; }
function edit(row: PlatformAdminRow) { editing.value = row; Object.assign(form, { account: row.account, password: '', real_name: row.real_name, phone: row.phone, roles: row.roles, status: row.status, merchant_ids: row.merchant_ids, region_ids: row.region_ids, service_store_ids: row.service_store_ids, circle_agent_id: row.circle_agent_id }); dialogOpen.value = true; }
async function save() {
  if (!editing.value && (!form.account.trim() || form.password.length < 8)) { ElMessage.warning('账号与至少 8 位的初始密码必填'); return; }
	if (form.roles.split(',').map((value) => value.trim()).includes('merchant') && !form.merchant_ids.trim()) { ElMessage.warning('商户角色必须填写授权商户 ID，多个以逗号分隔'); return; }
	if (isRegionAccount.value && !form.region_ids.trim()) { ElMessage.warning('区域管理员必须填写可管理区域 ID，多个以逗号分隔'); return; }
	if (isRegionAccount.value && !form.circle_agent_id) { ElMessage.warning('区域管理员必须关联一名已审核通过的区域代理'); return; }
	if (form.roles.split(',').map((value) => value.trim()).includes('customer_service') && !form.service_store_ids.trim()) { ElMessage.warning('客服账号必须填写授权店铺 ID，多个以逗号分隔'); return; }
	const payload = { ...form, account: form.account.trim(), real_name: form.real_name.trim(), phone: form.phone.trim(), roles: form.roles.trim(), merchant_ids: form.merchant_ids.trim(), region_ids: form.region_ids.trim(), service_store_ids: form.service_store_ids.trim() };
  if (editing.value) await updatePlatformAdmin(editing.value.admin_id, payload); else await createPlatformAdmin(payload);
  ElMessage.success('管理员已保存'); dialogOpen.value = false; await load();
}
async function remove(row: PlatformAdminRow) {
  try {
    await ElMessageBox.confirm(`删除“${row.real_name || row.account}”后，该账号将无法登录，历史操作与审计记录会保留。是否继续？`, '逻辑删除管理员', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' });
    await deletePlatformAdmin(row.admin_id);
    ElMessage.success('管理员已逻辑删除并强制失效');
    if (rows.value.length === 1 && query.page > 1) query.page -= 1;
    await load();
  } catch {
    // 取消不提示；请求客户端负责展示服务端保护规则。
  }
}
onMounted(() => { void load(); void loadAgentOptions(); });
</script>

<template>
  <Page title="管理员管理" description="平台管理员拥有全量数据；区域管理员与客服账号必须配置服务端数据范围，前端菜单隐藏不构成授权。">
    <el-card shadow="never"><div class="mb-4"><el-button type="primary" @click="add">新增管理员</el-button></div>
      <el-table v-loading="loading" :data="rows" border><el-table-column prop="admin_id" label="ID" width="72" /><el-table-column prop="account" label="账号" min-width="130" /><el-table-column prop="real_name" label="姓名" min-width="110" /><el-table-column prop="phone" label="手机号" width="140" /><el-table-column prop="roles" label="角色" min-width="144" /><el-table-column prop="merchant_ids" label="授权商户" min-width="120"><template #default="{ row }">{{ row.merchant_ids || '—' }}</template></el-table-column><el-table-column prop="region_ids" label="可管理区域" min-width="132"><template #default="{ row }">{{ row.region_ids || '—' }}</template></el-table-column><el-table-column prop="service_store_ids" label="客服授权店铺" min-width="148"><template #default="{ row }">{{ row.service_store_ids || '—' }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="130" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table>
      <div class="mt-4 flex justify-end"><el-pagination v-model:current-page="query.page" v-model:page-size="query.limit" layout="total, sizes, prev, pager, next" :total="total" @current-change="load" @size-change="() => { query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="title" width="620px"><el-form label-width="120px"><el-form-item label="登录账号" required><el-input v-model="form.account" :disabled="!!editing" /></el-form-item><el-form-item :label="editing ? '重置密码' : '初始密码'" :required="!editing"><el-input v-model="form.password" show-password type="password" :placeholder="editing ? '留空则不修改' : '至少 8 位'" /></el-form-item><el-form-item label="姓名"><el-input v-model="form.real_name" /></el-form-item><el-form-item label="手机号"><el-input v-model="form.phone" /></el-form-item><el-form-item label="角色代码" required><el-input v-model="form.roles" placeholder="例如 merchant；多个以逗号分隔" /></el-form-item><el-form-item v-if="form.roles.split(',').map((value) => value.trim()).includes('merchant')" label="授权商户" required><el-input v-model="form.merchant_ids" placeholder="例如 2001,2002（merchant_id）" /></el-form-item><el-form-item v-if="isRegionAccount" label="可管理区域" required><el-input v-model="form.region_ids" placeholder="例如 1,2,3（商户 region_id）" /></el-form-item><el-form-item v-if="form.roles.split(',').map((value) => value.trim()).includes('customer_service')" label="客服授权店铺" required><el-input v-model="form.service_store_ids" placeholder="例如 1001,1002（店铺 store_id）" /></el-form-item><el-form-item v-if="isRegionAccount" label="关联代理" required><el-select v-model="form.circle_agent_id" class="w-full" filterable placeholder="请选择已审核通过的代理"><el-option v-for="agent in agentOptions" :key="agent.circle_agent_id" :label="`${agent.name}（ID ${agent.circle_agent_id}）`" :value="agent.circle_agent_id" /></el-select></el-form-item><el-form-item label="启用状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
