<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { exportPlatformUsers, fetchPlatformUserDetail, fetchPlatformUsers, type PlatformUserDetail, type PlatformUserRow } from '#/api/core/ecrm';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const detailLoading = ref(false);
const detailOpen = ref(false);
const rows = ref<PlatformUserRow[]>([]);
const total = ref(0);
const detail = ref<PlatformUserDetail>();
const canRead = ref(false);
const canExport = ref(false);
const isPlatform = ref(false);
const query = reactive({ page: 1, limit: 20, id: undefined as number | undefined, keyword: '', status: undefined as 0 | 1 | undefined });

async function load() {
  if (!isPlatform.value) return;
  loading.value = true;
  try {
    const result = await fetchPlatformUsers({ ...query, keyword: query.keyword.trim() || undefined });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function reset() {
  Object.assign(query, { page: 1, id: undefined, keyword: '', status: undefined });
  void load();
}

async function exportRows() {
  try {
    const { value } = await ElMessageBox.prompt(
      '请填写导出原因，导出仅包含脱敏最小字段，最多 5000 行。',
      '导出用户信息',
      { inputPattern: /.{2,}/, inputErrorMessage: '导出原因至少 2 个字符', confirmButtonText: '生成 CSV', cancelButtonText: '取消' },
    );
    const result = await exportPlatformUsers({
      id: query.id,
      keyword: query.keyword.trim() || undefined,
      status: query.status,
      reason: value.trim(),
    });
    const blob = new Blob([result.content], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(`已导出 ${result.row_count} 条脱敏用户记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`);
  } catch {}
}

async function openDetail(row: PlatformUserRow) {
  detailOpen.value = true;
  detail.value = undefined;
  detailLoading.value = true;
  try {
    detail.value = await fetchPlatformUserDetail(row.id);
  } finally {
    detailLoading.value = false;
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  isPlatform.value = profile.roles.includes('platform');
  canRead.value = isPlatform.value && codes.includes('user.list.read');
  canExport.value = isPlatform.value && codes.includes('user.list.export');
  await load();
});

const assetLabels: Record<string, string> = { balance: '余额', points: '积分', commission: '佣金' };
const orderLabels: Record<string, string> = {
  pending_pay: '待支付', paid: '已支付', awaiting_final: '待尾款', final_timeout: '尾款超时', fulfilling: '履约中',
  shipped: '已发货', completed: '已完成', cancelled: '已取消', aftersale: '售后中',
};
const membershipLabels: Record<string, string> = { initial: '初始', upgrade: '升级', downgrade: '降级', manual: '人工调整' };
const couponStatusLabels: Record<string, string> = { unused: '未使用', locked: '锁定', used: '已使用', expired: '已过期' };
</script>

<template>
  <Page title="用户列表" description="仅展示脱敏手机号；详情仅提供订单、资产、会员变更、签到、浏览、优惠券和推荐关系等最小化监管读模型。余额、积分、会员等级、优惠券、推荐关系、分组、标签、推广资格、启停与身份维护均由独立受控页面执行。">
    <el-alert
      v-if="isPlatform && !canRead"
      class="mb-4"
      title="当前账号缺少 user.list.read 权限，列表仍尝试加载；导出与敏感操作已受限。"
      type="warning"
      :closable="false"
    />
    <EcrmListPage title="用户列表" description="平台监管用户只读视图；不含地址、支付凭据或提现账户。">
      <template #filters>
        <el-form inline @submit.prevent="query.page = 1; load()">
          <el-form-item label="用户 ID"><el-input-number v-model="query.id" :min="1" /></el-form-item>
          <el-form-item label="昵称"><el-input v-model="query.keyword" maxlength="64" clearable /></el-form-item>
          <el-form-item label="状态">
            <el-select v-model="query.status" clearable>
              <el-option label="启用" :value="1" />
              <el-option label="停用" :value="0" />
            </el-select>
          </el-form-item>
        </el-form>
      </template>
      <template #actions>
        <el-button type="primary" @click="query.page = 1; load()">查询</el-button>
        <el-button @click="reset">重置</el-button>
        <el-button v-if="canExport" type="success" plain @click="exportRows">导出脱敏信息</el-button>
      </template>
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="90" />
        <el-table-column prop="nickname" label="昵称" min-width="160" />
        <el-table-column prop="mobile" label="手机号（脱敏）" width="150" />
        <el-table-column label="余额" width="110"><template #default="{ row }">¥{{ Number(row.balance).toFixed(2) }}</template></el-table-column>
        <el-table-column prop="points" label="积分" width="100" />
        <el-table-column label="会员等级" min-width="130"><template #default="{ row }">{{ row.level_name || '普通会员' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="90" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="openDetail(row)">监管详情</el-button></template></el-table-column>
      </el-table>
      <template #pager>
        <el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" />
      </template>
    </EcrmListPage>
    <el-drawer v-model="detailOpen" title="用户监管详情" size="780px" destroy-on-close>
      <div v-loading="detailLoading">
        <template v-if="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户 ID">{{ detail.profile.id }}</el-descriptions-item>
            <el-descriptions-item label="昵称">{{ detail.profile.nickname }}</el-descriptions-item>
            <el-descriptions-item label="手机号（脱敏）">{{ detail.profile.mobile || '—' }}</el-descriptions-item>
            <el-descriptions-item label="注册时间">{{ detail.profile.created_at }}</el-descriptions-item>
            <el-descriptions-item label="账户余额">¥{{ Number(detail.profile.balance).toFixed(2) }}</el-descriptions-item>
            <el-descriptions-item label="积分 / 佣金">{{ detail.profile.points }} / ¥{{ Number(detail.profile.commission).toFixed(2) }}</el-descriptions-item>
          </el-descriptions>
          <el-tabs class="mt-4">
            <el-tab-pane label="最近订单">
              <el-table :data="detail.orders" max-height="360">
                <el-table-column prop="order_no" label="订单号" min-width="180" />
                <el-table-column prop="store_name" label="店铺" min-width="130" />
                <el-table-column label="实付" width="100"><template #default="{ row }">¥{{ Number(row.pay_amount).toFixed(2) }}</template></el-table-column>
                <el-table-column prop="total_quantity" label="件数" width="70" />
                <el-table-column label="状态" width="100"><template #default="{ row }">{{ orderLabels[row.status] || row.status }}</template></el-table-column>
                <el-table-column prop="created_at" label="创建时间" width="170" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="资产流水">
              <el-table :data="detail.assets" max-height="360">
                <el-table-column prop="id" label="ID" width="80" />
                <el-table-column label="类型" width="90"><template #default="{ row }">{{ assetLabels[row.asset_type] || row.asset_type }}</template></el-table-column>
                <el-table-column label="变动" width="110"><template #default="{ row }"><span :class="row.amount < 0 ? 'text-red-500' : 'text-green-600'">{{ row.amount > 0 ? '+' : '' }}{{ Number(row.amount).toFixed(2) }}</span></template></el-table-column>
                <el-table-column prop="reference_type" label="业务来源" min-width="120" />
                <el-table-column prop="reference_id" label="业务引用" min-width="120" />
                <el-table-column prop="created_at" label="创建时间" width="170" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="会员变更">
              <el-table :data="detail.membership_logs" max-height="360">
                <el-table-column prop="previous_level_name" label="原等级" min-width="110" />
                <el-table-column prop="level_name" label="当前等级" min-width="110" />
                <el-table-column label="变更类型" width="100"><template #default="{ row }">{{ membershipLabels[row.change_type] || row.change_type }}</template></el-table-column>
                <el-table-column prop="note" label="说明" min-width="220" />
                <el-table-column prop="created_at" label="时间" width="170" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="签到记录">
              <el-table :data="detail.signs" max-height="360">
                <el-table-column prop="sign_date" label="签到日期" min-width="140" />
                <el-table-column prop="points" label="获得积分" width="110" />
                <el-table-column prop="continuous_days" label="连续天数" width="110" />
                <el-table-column prop="created_at" label="记录时间" width="170" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="浏览记录">
              <el-table :data="detail.browse_history" max-height="360">
                <el-table-column prop="product_id" label="商品 ID" width="100" />
                <el-table-column prop="title" label="商品" min-width="180" />
                <el-table-column prop="store_name" label="店铺" min-width="140" />
                <el-table-column prop="viewed_at" label="浏览时间" width="170" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="持有优惠券">
              <el-table :data="detail.coupons" max-height="360">
                <el-table-column prop="coupon_id" label="券 ID" width="90" />
                <el-table-column prop="name" label="优惠券" min-width="160" />
                <el-table-column label="优惠" width="100"><template #default="{ row }">{{ row.discount_type === 'rate' ? `${row.discount_value / 10} 折` : `¥${Number(row.discount_value).toFixed(2)}` }}</template></el-table-column>
                <el-table-column label="门槛" width="100"><template #default="{ row }">¥{{ Number(row.min_amount).toFixed(2) }}</template></el-table-column>
                <el-table-column label="状态" width="100"><template #default="{ row }">{{ couponStatusLabels[row.status] || row.status }}</template></el-table-column>
                <el-table-column prop="obtained_at" label="领取时间" width="170" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="推荐关系">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="上级用户">{{ detail.distribution.parent_user_id ? `${detail.distribution.parent_nickname || '未知用户'}（#${detail.distribution.parent_user_id}）` : '未绑定' }}</el-descriptions-item>
                <el-descriptions-item label="直推用户数">{{ detail.distribution.direct_user_count }}</el-descriptions-item>
                <el-descriptions-item label="推广员资格"><el-tag :type="detail.distribution.promoter_status === 1 ? 'success' : 'info'">{{ detail.distribution.promoter_status === 1 ? '启用' : '未开通或已停用' }}</el-tag></el-descriptions-item>
              </el-descriptions>
            </el-tab-pane>
          </el-tabs>
          <el-alert class="mt-4" type="info" :closable="false" title="为保护个人与交易敏感信息，页面不展示收货地址、发票资料、支付交易号或提现账户快照。" />
        </template>
      </div>
    </el-drawer>
  </Page>
</template>
