<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  listPlatformAssistApi,
  getPlatformAssistApi,
  updatePlatformAssistApi,
  type PlatformAssistActive,
} from '#/api/core/platform-assist';

const loading = ref(false);
const rows = ref<PlatformAssistActive[]>([]);
const total = ref(0);
const canManage = ref(false);
const query = reactive({ limit: 20, mer_id: undefined as number | undefined, page: 1 });
const detailOpen = ref(false);
const detail = ref<PlatformAssistActive>();

function time(value: string) {
  return formatShanghaiDateTime(value);
}

async function showDetail(row: PlatformAssistActive) {
  detail.value = await getPlatformAssistApi(row.product_assist_id);
  detailOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listPlatformAssistApi(query);
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function setVisible(row: PlatformAssistActive, isShow: number) {
  const action = isShow === 1 ? '上架' : '下架';
  try {
    await ElMessageBox.confirm(
      `确认${action}好友助力活动“${row.store_name || `#${row.product_assist_id}`}”吗？`,
      `${action}确认`,
      { cancelButtonText: '取消', confirmButtonText: `确认${action}`, type: 'warning' },
    );
    await updatePlatformAssistApi(row.product_assist_id, { is_show: isShow });
    ElMessage.success(`活动已${action}`);
    await load();
  } catch {
    // 用户取消或接口错误由统一请求层提示。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    load(),
  ]);
  canManage.value = profile.roles.some((role) => role === 'platform' || role === 'operations')
    && permissions.includes('marketing.assist.manage');
});
</script>

<template>
  <Page title="好友助力监管" description="查看各商户好友助力活动及完整规则；具备运营权限可上架或下架。删除、改价、库存与时间调整须完成订单影响审计后另行开放。">
    <el-card shadow="never">
      <el-form inline @submit.prevent="query.page = 1; load()">
        <el-form-item label="商户 ID"><el-input-number v-model="query.mer_id" :min="1" /></el-form-item>
        <el-button type="primary" @click="query.page = 1; load()">查询</el-button>
        <el-button @click="query.mer_id = undefined; query.page = 1; load()">重置</el-button>
      </el-form>
      <el-table v-loading="loading" :data="rows">
        <el-table-column label="活动 ID" prop="product_assist_id" width="96" />
        <el-table-column label="活动 / 商品" min-width="180"><template #default="{ row }">{{ row.store_name || `商品 #${row.product_id}` }}</template></el-table-column>
        <el-table-column label="商户" min-width="130"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column>
        <el-table-column label="助力价" width="110"><template #default="{ row }">¥{{ Number(row.assist_price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="助力规则" min-width="130"><template #default="{ row }">{{ row.assist_count }} 人 / 每人最多 {{ row.assist_user_count }} 次</template></el-table-column>
        <el-table-column prop="stock" label="活动库存" width="100" />
        <el-table-column label="活动时间" min-width="240"><template #default="{ row }">{{ time(row.start_time) }} 至 {{ time(row.end_time) }}</template></el-table-column>
        <el-table-column label="展示状态" width="100"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '已上架' : '已下架' }}</el-tag></template></el-table-column>
        <el-table-column fixed="right" label="操作" :width="canManage ? 146 : 70"><template #default="{ row }"><el-button link type="primary" @click="showDetail(row)">详情</el-button><el-button v-if="canManage" link :type="row.is_show === 1 ? 'danger' : 'success'" @click="setVisible(row, row.is_show === 1 ? 0 : 1)">{{ row.is_show === 1 ? '下架' : '上架' }}</el-button></template></el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
    <el-dialog v-model="detailOpen" title="好友助力活动详情" width="640px" destroy-on-close>
      <el-descriptions v-if="detail" :column="2" border>
        <el-descriptions-item label="活动名称" :span="2">{{ detail.store_name }}</el-descriptions-item>
        <el-descriptions-item label="商品 / 商户">#{{ detail.product_id }} / {{ detail.mer_name || `商户 #${detail.mer_id}` }}</el-descriptions-item>
        <el-descriptions-item label="助力价">¥{{ Number(detail.assist_price).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="所需助力">{{ detail.assist_count }} 人</el-descriptions-item>
        <el-descriptions-item label="单人助力次数">最多 {{ detail.assist_user_count }} 次</el-descriptions-item>
        <el-descriptions-item label="活动库存">{{ detail.stock }}</el-descriptions-item>
        <el-descriptions-item label="活动时间" :span="2">{{ time(detail.start_time) }} 至 {{ time(detail.end_time) }}</el-descriptions-item>
        <el-descriptions-item label="前台展示">{{ detail.is_show === 1 ? '上架' : '下架' }}</el-descriptions-item>
        <el-descriptions-item label="活动状态">{{ detail.status === 1 ? '启用' : '停用' }}</el-descriptions-item>
      </el-descriptions>
      <el-alert class="mt-4" type="warning" :closable="false" title="详情不展示参与用户资料。已发起助力单的价格、库存与完成条件须保持订单快照，本页不提供这些字段编辑。" />
    </el-dialog>
  </Page>
</template>
