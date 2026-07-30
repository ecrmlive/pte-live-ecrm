<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { listPlatformSpreadLogsApi, type SpreadLog } from '#/api/core/platform-spread';
const loading = ref(false); const rows = ref<SpreadLog[]>([]); const total = ref(0); const query = reactive({ limit: 20, page: 1 });
async function load() { loading.value = true; try { const data = await listPlatformSpreadLogsApi(query); rows.value = data.list || []; total.value = data.total || 0; } finally { loading.value = false; } }
onMounted(() => void load());
</script>
<template><Page title="分销监管" description="当前已接入推广关系变更日志监管；分销员等级、佣金结算和提现属于后续独立模块。"><el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="user_spread_log_id"><el-table-column label="日志 ID" prop="user_spread_log_id" width="100" /><el-table-column label="用户 ID" prop="uid" width="100" /><el-table-column label="原推荐人" prop="old_spread_uid" width="120" /><el-table-column label="新推荐人" prop="spread_uid" width="120" /><el-table-column label="变更时间" min-width="180" prop="create_time" /></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10,20,50,100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card></Page></template>
