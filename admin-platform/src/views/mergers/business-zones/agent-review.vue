<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { auditBusinessZoneAgent, fetchBusinessZoneAgents, type BusinessZoneAgentRow } from '#/api/core/mergers';

const rows = ref<BusinessZoneAgentRow[]>([]);
const total = ref(0);
const loading = ref(false);
const query = reactive({ keyword: '', page: 1, limit: 20, status: 0 });
const statusText = (value: number) => ({ '-1': '已驳回', '0': '待审核', '1': '已通过' }[String(value)] || '未知');

async function load() { loading.value = true; try { const result = await fetchBusinessZoneAgents({ ...query, keyword: query.keyword.trim() || undefined }); rows.value = result.list; total.value = result.total; } finally { loading.value = false; } }
async function approve(row: BusinessZoneAgentRow) { try { await ElMessageBox.confirm(`确认通过“${row.name}”的代理申请？`, '代理审核'); await auditBusinessZoneAgent(row.circle_agent_id, 1); ElMessage.success('审核已通过'); await load(); } catch {} }
async function reject(row: BusinessZoneAgentRow) { try { const { value } = await ElMessageBox.prompt('请填写驳回原因', '驳回代理申请', { inputPattern: /.+/, inputErrorMessage: '驳回原因必填' }); await auditBusinessZoneAgent(row.circle_agent_id, -1, value.trim()); ElMessage.success('已驳回'); await load(); } catch {} }
function search() { query.page = 1; void load(); }
function reset() { query.keyword = ''; query.status = 0; query.page = 1; void load(); }
onMounted(load);
</script>

<template>
  <Page title="代理审核" description="待审核的代理申请必须明确通过或驳回；驳回原因会保留在代理档案中。">
    <el-card shadow="never"><el-form class="grid gap-x-4 md:grid-cols-3" label-width="72px" @submit.prevent="search"><el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="姓名/手机号/商户名" @keyup.enter="search" /></el-form-item><el-form-item label="审核状态"><el-select v-model="query.status" class="w-full"><el-option label="待审核" :value="0" /><el-option label="已通过" :value="1" /><el-option label="已驳回" :value="-1" /></el-select></el-form-item><el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item></el-form>
      <el-table v-loading="loading" :data="rows" border><el-table-column prop="circle_agent_id" label="ID" width="72" /><el-table-column prop="name" label="代理人" min-width="110" /><el-table-column prop="phone" label="手机号" width="140" /><el-table-column prop="business_name" label="关联商户" min-width="140" /><el-table-column prop="qualification" label="资质说明" min-width="180" show-overflow-tooltip /><el-table-column prop="audit_reason" label="审核说明" min-width="160" show-overflow-tooltip /><el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : row.status === -1 ? 'danger' : 'warning'">{{ statusText(row.status) }}</el-tag></template></el-table-column><el-table-column label="操作" width="150" fixed="right"><template #default="{ row }"><template v-if="row.status === 0"><el-button link type="success" @click="approve(row)">通过</el-button><el-button link type="danger" @click="reject(row)">驳回</el-button></template><span v-else>—</span></template></el-table-column></el-table>
      <div class="mt-4 flex justify-end"><el-pagination v-model:current-page="query.page" v-model:page-size="query.limit" layout="total, sizes, prev, pager, next" :total="total" @current-change="load" @size-change="() => { query.page = 1; load(); }" /></div>
    </el-card>
  </Page>
</template>
