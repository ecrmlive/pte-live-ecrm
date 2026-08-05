<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformPresellsApi,
  getPlatformPresellApi,
  updatePlatformPresellApi,
  type PlatformPresell,
} from '#/api/core/platform-presell';

const loading = ref(false);
const rows = ref<PlatformPresell[]>([]);
const total = ref(0);
const canManage = ref(false);
const query = reactive({ limit: 20, mer_id: undefined as number | undefined, page: 1 });
const detailOpen = ref(false);
const detail = ref<PlatformPresell>();

async function load() {
  loading.value = true;
  try {
    const data = await listPlatformPresellsApi(query);
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function showDetail(row: PlatformPresell) {
  detail.value = await getPlatformPresellApi(row.product_presell_id);
  detailOpen.value = true;
}

async function setStatus(row: PlatformPresell, status: number) {
  const action = status === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(`确认${action}预售活动“${row.store_name}”？`, `${action}确认`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await updatePlatformPresellApi(row.product_presell_id, { status });
    ElMessage.success(`活动已${action}`);
    await load();
  } catch {
    // 用户取消或接口错误由统一请求层提示。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi(), load()]);
  canManage.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && permissions.includes('marketing.presell.manage');
});
</script>

<template>
  <Page title="预售监管" description="监管各商户预售活动；可查看完整价格、库存、发货与尾款时间窗。具备运营活动权限的账号可启停，资金字段调整待订单快照影响审计闭环后开放。">
    <el-card shadow="never">
      <el-form inline @submit.prevent="query.page = 1; load()">
        <el-form-item label="商户 ID"><el-input-number v-model="query.mer_id" :min="1" /></el-form-item>
        <el-button type="primary" @click="query.page = 1; load()">查询</el-button>
      </el-form>
      <el-table v-loading="loading" :data="rows">
        <el-table-column label="预售商品" min-width="180" prop="store_name" />
        <el-table-column label="商户" min-width="130"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column>
        <el-table-column label="预售价" width="110"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="时间" min-width="220"><template #default="{ row }">{{ row.start_time }} 至 {{ row.end_time }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column fixed="right" label="操作" :width="canManage ? 146 : 70"><template #default="{ row }"><el-button link type="primary" @click="showDetail(row)">详情</el-button><el-button v-if="canManage" link :type="row.status === 1 ? 'danger' : 'success'" @click="setStatus(row, row.status === 1 ? 0 : 1)">{{ row.status === 1 ? '停用' : '启用' }}</el-button></template></el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
    <el-dialog v-model="detailOpen" title="预售活动详情" width="660px" destroy-on-close>
      <el-descriptions v-if="detail" :column="2" border>
        <el-descriptions-item label="活动名称" :span="2">{{ detail.store_name }}</el-descriptions-item>
        <el-descriptions-item label="商品 / 商户">#{{ detail.product_id }} / {{ detail.mer_name || `商户 #${detail.mer_id}` }}</el-descriptions-item>
        <el-descriptions-item label="活动类型">{{ detail.presell_type === 2 ? '定金预售' : '全款预售' }}</el-descriptions-item>
        <el-descriptions-item label="预售价">¥{{ Number(detail.price).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="定金 / 尾款">{{ detail.presell_type === 2 ? `¥${Number(detail.down_price || 0).toFixed(2)} / ¥${Number(detail.final_price || 0).toFixed(2)}` : '不适用' }}</el-descriptions-item>
        <el-descriptions-item label="活动时间" :span="2">{{ detail.start_time }} 至 {{ detail.end_time }}</el-descriptions-item>
        <el-descriptions-item label="尾款时间" :span="2">{{ detail.presell_type === 2 ? `${detail.final_start_time || '—'} 至 ${detail.final_end_time || '—'}` : '不适用' }}</el-descriptions-item>
        <el-descriptions-item label="库存 / 已售">{{ detail.stock || 0 }} / {{ detail.seles || 0 }}</el-descriptions-item>
        <el-descriptions-item label="发货">{{ detail.delivery_type === 2 ? '发货后' : '付款后' }}{{ detail.delivery_day ? ` ${detail.delivery_day} 天` : '' }}</el-descriptions-item>
        <el-descriptions-item label="前台展示">{{ detail.is_show === 1 ? '上架' : '下架' }}</el-descriptions-item>
        <el-descriptions-item label="活动状态">{{ detail.status === 1 ? '启用' : '停用' }}</el-descriptions-item>
        <el-descriptions-item label="活动说明" :span="2">{{ detail.store_info || '—' }}</el-descriptions-item>
      </el-descriptions>
      <el-alert class="mt-4" type="warning" :closable="false" title="已产生定金或尾款订单的金额、时间窗须保持订单快照。本页不提供资金字段编辑，待订单影响审计闭环后再开放。" />
    </el-dialog>
  </Page>
</template>
