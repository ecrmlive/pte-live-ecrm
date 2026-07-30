<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  auditPlatformProductApi,
  getPlatformProductApi,
  listPlatformProductsApi,
  type PlatformProduct,
} from '#/api/core/platform-catalog';

const loading = ref(false);
const rows = ref<PlatformProduct[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const current = ref<PlatformProduct>();
const rejecting = ref(false);
const rejectForm = reactive({ refusal: '' });
const query = reactive({ keyword: '', limit: 20, mer_id: undefined as number | undefined, page: 1, status: undefined as number | undefined });

const statusInfo = (status: number) => ({ '-2': { label: '已下架', type: 'info' }, '-1': { label: '审核驳回', type: 'danger' }, 0: { label: '待审核', type: 'warning' }, 1: { label: '已通过', type: 'success' } } as Record<number, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }>)[status] || { label: '未知', type: 'info' as const };

async function load() {
  loading.value = true;
  try {
    const data = await listPlatformProductsApi({ ...query, keyword: query.keyword.trim() || undefined });
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally { loading.value = false; }
}

function search() { query.page = 1; void load(); }
function reset() { Object.assign(query, { keyword: '', mer_id: undefined, page: 1, status: undefined }); void load(); }
async function detail(row: PlatformProduct) { current.value = await getPlatformProductApi(row.product_id); detailOpen.value = true; }

async function approve(row: PlatformProduct) {
  try {
    await ElMessageBox.confirm(`确认审核通过商品“${row.store_name}”？`, '商品审核', { type: 'warning' });
    await auditPlatformProductApi(row.product_id, { status: 1 });
    ElMessage.success('商品已审核通过');
    await load();
  } catch { /* 用户取消或统一请求错误提示。 */ }
}

function reject(row: PlatformProduct) { current.value = row; rejectForm.refusal = ''; rejectOpen.value = true; }
async function submitReject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) { ElMessage.warning('请填写驳回原因'); return; }
  if (!current.value) return;
  rejecting.value = true;
  try {
    await auditPlatformProductApi(current.value.product_id, { status: -1, refusal });
    rejectOpen.value = false;
    ElMessage.success('商品已驳回');
    await load();
  } finally { rejecting.value = false; }
}

onMounted(() => void load());
</script>

<template>
  <Page title="商品审核" description="审核全平台商户商品；通过、驳回均受独立 product/audit 按钮权限控制，区域账号不可访问。">
    <el-card shadow="never"><el-form class="grid gap-x-4 md:grid-cols-4" label-width="72px" @submit.prevent="search"><el-form-item label="商品搜索"><el-input v-model="query.keyword" clearable placeholder="名称 / 关键词" /></el-form-item><el-form-item label="商户 ID"><el-input-number v-model="query.mer_id" :min="1" class="w-full" controls-position="right" /></el-form-item><el-form-item label="审核状态"><el-select v-model="query.status" clearable class="w-full" placeholder="全部"><el-option label="待审核" :value="0" /><el-option label="已通过" :value="1" /><el-option label="已驳回" :value="-1" /><el-option label="已下架" :value="-2" /></el-select></el-form-item><el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item></el-form></el-card>
    <el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows" row-key="product_id"><el-table-column label="ID" prop="product_id" width="80" /><el-table-column label="商品名称" min-width="180" prop="store_name" show-overflow-tooltip /><el-table-column label="商户" min-width="130"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column><el-table-column label="分类" min-width="110" prop="cate_name" /><el-table-column label="售价" width="108"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column><el-table-column label="库存/销量" width="110"><template #default="{ row }">{{ row.stock }} / {{ row.sales }}</template></el-table-column><el-table-column label="状态" width="98"><template #default="{ row }"><el-tag :type="statusInfo(row.status).type">{{ statusInfo(row.status).label }}</el-tag></template></el-table-column><el-table-column label="驳回原因" min-width="150" prop="refusal" show-overflow-tooltip /><el-table-column label="创建时间" min-width="170" prop="create_time" /><el-table-column fixed="right" label="操作" width="162"><template #default="{ row }"><el-button link type="primary" @click="detail(row)">详情</el-button><template v-if="row.status === 0"><el-button link type="success" @click="approve(row)">通过</el-button><el-button link type="danger" @click="reject(row)">驳回</el-button></template></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-drawer v-model="detailOpen" :with-header="false" size="560px"><template v-if="current"><div class="mb-5 text-lg font-medium">商品详情</div><el-descriptions :column="1" border><el-descriptions-item label="商品名称">{{ current.store_name }}</el-descriptions-item><el-descriptions-item label="商户">{{ current.mer_name || `商户 #${current.mer_id}` }}</el-descriptions-item><el-descriptions-item label="分类">{{ current.cate_name || '—' }}</el-descriptions-item><el-descriptions-item label="售价">¥{{ Number(current.price).toFixed(2) }}</el-descriptions-item><el-descriptions-item label="库存 / 销量">{{ current.stock }} / {{ current.sales }}</el-descriptions-item><el-descriptions-item label="商品简介">{{ current.store_info || '—' }}</el-descriptions-item><el-descriptions-item label="审核状态"><el-tag :type="statusInfo(current.status).type">{{ statusInfo(current.status).label }}</el-tag></el-descriptions-item><el-descriptions-item v-if="current.refusal" label="驳回原因">{{ current.refusal }}</el-descriptions-item></el-descriptions></template></el-drawer>
    <el-dialog v-model="rejectOpen" title="驳回商品" width="480px" destroy-on-close><el-form label-width="84px"><el-form-item label="驳回原因" required><el-input v-model="rejectForm.refusal" :rows="4" maxlength="200" placeholder="请向商户说明驳回原因" show-word-limit type="textarea" /></el-form-item></el-form><template #footer><el-button @click="rejectOpen = false">取消</el-button><el-button :loading="rejecting" type="danger" @click="submitReject">确认驳回</el-button></template></el-dialog>
  </Page>
</template>
