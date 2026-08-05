<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { clearUserSearchRecords, exportUserSearchRecords, listUserSearchRecords, type UserSearchRecord } from '#/api/core/platform-user-search';
import { formatShanghaiDateTime } from '#/utils/date-time';

const rows = ref<UserSearchRecord[]>([]);
const total = ref(0);
const loading = ref(false);
const canRead = ref(false);
const canClear = ref(false);
const canExport = ref(false);
const query = reactive({ page: 1, limit: 20, user_id: undefined as number | undefined, keyword: '', source: undefined as 'h5' | 'mini' | 'pc' | undefined, dates: [] as string[] });
const sourceLabel: Record<string, string> = { pc: 'PC', h5: 'H5', mini: '小程序' };
function formatDate(value: string) { return formatShanghaiDateTime(value); }

function filters() { return { user_id: query.user_id, keyword: query.keyword.trim() || undefined, source: query.source, start_date: query.dates[0], end_date: query.dates[1] }; }
async function load() { if (!canRead.value) return; loading.value = true; try { const result = await listUserSearchRecords({ page: query.page, limit: query.limit, ...filters() }); rows.value = result.list || []; total.value = result.total || 0; } finally { loading.value = false; } }
function reset() { Object.assign(query, { page: 1, user_id: undefined, keyword: '', source: undefined, dates: [] }); void load(); }
async function clearRecords() { try { const user = await ElMessageBox.prompt('请输入需清理的用户 ID。本操作仅逻辑删除该用户当前可见搜索记录，不影响订单、收藏或浏览记录。', '按用户清理搜索记录', { inputPattern: /^[1-9]\d*$/, inputErrorMessage: '请输入正整数用户 ID' }); const reason = await ElMessageBox.prompt('请填写清理原因（2 至 500 个字符）。', '清理原因', { inputPattern: /.{2,}/, inputErrorMessage: '清理原因至少 2 个字符' }); const result = await clearUserSearchRecords({ user_id: Number(user.value), reason: reason.value.trim(), idempotency_key: `search-clear-${user.value}-${crypto.randomUUID()}` }); ElMessage.success(`已逻辑清理 ${result.cleared_count} 条搜索记录`); await load(); } catch {} }
async function exportRows() { try { const prompt = await ElMessageBox.prompt('请填写导出原因。导出含用户 ID 与搜索词，最多 5000 条，请按最小必要原则使用。', '导出搜索记录', { inputPattern: /.{2,}/, inputErrorMessage: '导出原因至少 2 个字符' }); const result = await exportUserSearchRecords({ ...filters(), reason: prompt.value.trim() }); const blob = new Blob([result.content], { type: 'text/csv;charset=utf-8' }); const url = URL.createObjectURL(blob); const link = document.createElement('a'); link.href = url; link.download = result.file_name; link.click(); URL.revokeObjectURL(url); ElMessage.success(`已导出 ${result.row_count} 条搜索记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`); } catch {} }
onMounted(async () => { const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]); const platform = profile.roles.includes('platform'); canRead.value = platform && codes.includes('user.search_record.read'); canClear.value = platform && codes.includes('user.search_record.clear'); canExport.value = platform && codes.includes('user.search_record.export'); await load(); });
</script>

<template>
  <Page title="用户搜索记录" description="平台监管 C 端搜索历史，仅展示用户 ID、关键词、来源和时间；不返回账号、手机号、地址或设备标识。清理为可审计的逻辑删除。">
    <el-alert v-if="!canRead" title="当前账号没有搜索记录查看权限" type="warning" :closable="false" />
    <template v-else><el-card shadow="never"><el-form inline @submit.prevent="query.page = 1; load()"><el-form-item label="用户 ID"><el-input-number v-model="query.user_id" :min="1" controls-position="right"/></el-form-item><el-form-item label="关键词"><el-input v-model="query.keyword" maxlength="128" clearable/></el-form-item><el-form-item label="来源"><el-select v-model="query.source" clearable class="w-28"><el-option label="PC" value="pc"/><el-option label="H5" value="h5"/><el-option label="小程序" value="mini"/></el-select></el-form-item><el-form-item label="日期"><el-date-picker v-model="query.dates" type="daterange" value-format="YYYY-MM-DD" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期"/></el-form-item><el-button type="primary" @click="query.page = 1; load()">查询</el-button><el-button @click="reset">重置</el-button><el-button v-if="canClear" type="warning" plain @click="clearRecords">按用户清理</el-button><el-button v-if="canExport" type="success" plain @click="exportRows">导出 CSV</el-button></el-form></el-card><el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows"><el-table-column prop="id" label="记录 ID" width="100"/><el-table-column prop="user_id" label="用户 ID" width="100"/><el-table-column prop="keyword" label="搜索关键词" min-width="220" show-overflow-tooltip/><el-table-column label="来源" width="100"><template #default="{ row }">{{ sourceLabel[row.source] || row.source }}</template></el-table-column><el-table-column label="搜索时间" min-width="180"><template #default="{ row }">{{ formatDate(row.created_at) }}</template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10,20,50,100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }"/></div></el-card><el-alert class="mt-4" title="导出会写入审计。CSV 使用 UTF-8 BOM 并对以 =、+、-、@ 开头的关键词转义，避免表格公式注入。" type="info" :closable="false"/></template>
  </Page>
</template>
