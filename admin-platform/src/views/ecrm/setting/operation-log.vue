<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { listPlatformOperationLogs, type PlatformOperationLog } from '#/api/core/platform-operation-log';

const rows = ref<PlatformOperationLog[]>([]);
const total = ref(0);
const loading = ref(false);
const canRead = ref(false);
const query = reactive({ page: 1, limit: 20, admin_user_id: undefined as number | undefined, role_code: '', action: '', resource_type: '', dates: [] as string[] });

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const result = await listPlatformOperationLogs({ page: query.page, limit: query.limit, admin_user_id: query.admin_user_id, role_code: query.role_code.trim() || undefined, action: query.action.trim() || undefined, resource_type: query.resource_type.trim() || undefined, start_date: query.dates[0], end_date: query.dates[1] });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally { loading.value = false; }
}
function reset() { Object.assign(query, { page: 1, admin_user_id: undefined, role_code: '', action: '', resource_type: '', dates: [] }); void load(); }

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && codes.includes('setting.operation_log.read');
  await load();
});
</script>

<template>
  <Page title="操作日志" description="只读记录统一后台成功写操作的管理员、角色、路由资源与时间。日志不保存请求体、密码、令牌、手机号或其他敏感数据，业务状态以各领域不可变审计为准。">
    <el-alert v-if="!canRead" title="当前账号没有操作日志查看权限" type="warning" :closable="false" />
    <template v-else><el-card shadow="never"><el-form inline @submit.prevent="query.page = 1; load()"><el-form-item label="管理员 ID"><el-input-number v-model="query.admin_user_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="角色"><el-input v-model="query.role_code" maxlength="32" clearable /></el-form-item><el-form-item label="操作"><el-input v-model="query.action" maxlength="128" clearable placeholder="如 POST /coupons" /></el-form-item><el-form-item label="资源"><el-input v-model="query.resource_type" maxlength="64" clearable /></el-form-item><el-form-item label="日期"><el-date-picker v-model="query.dates" type="daterange" value-format="YYYY-MM-DD" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" /></el-form-item><el-button type="primary" @click="query.page = 1; load()">查询</el-button><el-button @click="reset">重置</el-button></el-form></el-card><el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows"><el-table-column prop="id" label="日志 ID" width="100"/><el-table-column prop="admin_user_id" label="管理员 ID" width="120"/><el-table-column prop="role_code" label="角色" width="130"/><el-table-column prop="action" label="成功操作" min-width="300" show-overflow-tooltip/><el-table-column prop="resource_type" label="资源类型" width="140"/><el-table-column label="资源 ID" width="120"><template #default="{ row }">{{ row.resource_id || '—' }}</template></el-table-column><el-table-column prop="request_id" label="请求号" min-width="220" show-overflow-tooltip/><el-table-column prop="created_at" label="操作时间" min-width="180"/></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10,20,50,100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }"/></div></el-card></template>
  </Page>
</template>
