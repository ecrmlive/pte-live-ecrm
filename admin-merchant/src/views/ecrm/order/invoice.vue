<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  auditInvoice,
  fetchInvoices,
  type InvoiceRow,
} from '#/api/core/ecrm';
import { EcrmListPage } from '#/components/ecrm';

type StatusFilter = '' | 'requested' | 'issued' | 'rejected';

const loading = ref(false);
const rows = ref<InvoiceRow[]>([]);
const total = ref(0);
const canAudit = ref(false);
const query = reactive({
  limit: 20,
  page: 1,
  status: '' as StatusFilter,
});

const statusMap: Record<number, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }> = {
  [-1]: { label: '已驳回', type: 'danger' },
  0: { label: '待审核', type: 'warning' },
  1: { label: '已开票', type: 'success' },
};

function statusInfo(status: number) {
  return statusMap[status] || { label: '未知', type: 'info' as const };
}

function rowCanAudit(row: InvoiceRow) {
  return canAudit.value && row.status === 0;
}

async function load() {
  loading.value = true;
  try {
    const result = await fetchInvoices({
      page: query.page,
      limit: query.limit,
      status: query.status || undefined,
    });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.status = '';
  query.page = 1;
  void load();
}

async function onAudit(row: InvoiceRow, status: 1 | -1) {
  try {
    const { value } = await ElMessageBox.prompt(
      status === 1 ? '请填写发票号码（可留空自动生成）' : '请填写驳回原因',
      status === 1 ? '确认开票' : '驳回发票',
      {
        confirmButtonText: status === 1 ? '开票' : '驳回',
        cancelButtonText: '取消',
        inputPlaceholder: status === 1 ? '发票号码' : '驳回原因',
        inputValue: status === 1 ? '' : '商户驳回',
        type: status === 1 ? 'info' : 'warning',
      },
    );
    await auditInvoice(row.invoice_id, { status, mark: value?.trim() || undefined });
    ElMessage.success(status === 1 ? '已开票' : '已驳回');
    await load();
  } catch {
    // 取消或接口错误由统一请求层处理。
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canAudit.value = codes.includes('invoice.audit');
  await load();
});
</script>

<template>
  <EcrmListPage
    title="发票管理"
    description="读写 qixi_crm_b_order_invoice（按本店订单 store_id 隔离）；开票/驳回需 invoice.audit 权限。"
  >
    <template #filters>
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部" style="width: 140px">
            <el-option label="待审核" value="requested" />
            <el-option label="已开票" value="issued" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </template>

    <el-table v-loading="loading" :data="rows" border stripe>
      <el-table-column prop="invoice_id" label="ID" width="80" />
      <el-table-column prop="order_id" label="订单" width="100" />
      <el-table-column prop="uid" label="用户" width="90" />
      <el-table-column prop="header" label="抬头" min-width="160" show-overflow-tooltip />
      <el-table-column prop="tax_no" label="税号" width="160" show-overflow-tooltip />
      <el-table-column prop="email" label="邮箱" width="160" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusInfo(row.status).type" size="small">
            {{ statusInfo(row.status).label }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="mark" label="备注/票号" min-width="140" show-overflow-tooltip />
      <el-table-column prop="create_time" label="申请时间" width="170" />
      <el-table-column fixed="right" label="操作" width="140">
        <template #default="{ row }">
          <template v-if="rowCanAudit(row)">
            <el-button link type="primary" @click="onAudit(row, 1)">开票</el-button>
            <el-button link type="danger" @click="onAudit(row, -1)">驳回</el-button>
          </template>
          <span v-else class="text-gray-400">—</span>
        </template>
      </el-table-column>
    </el-table>

    <template #pager>
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.limit"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="load"
        @size-change="
          () => {
            query.page = 1;
            void load();
          }
        "
      />
    </template>
  </EcrmListPage>
</template>
