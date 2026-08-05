<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { fetchBusinessZoneAgentSettings, type BusinessZoneAgentSettings } from '#/api/core/ecrm';

const loading = ref(false);
const settings = ref<BusinessZoneAgentSettings>();

async function load() {
  loading.value = true;
  try {
    settings.value = await fetchBusinessZoneAgentSettings();
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page title="代理设置" description="CRMEB 对应只读配置入口。规则由服务端强制执行，页面仅用于核验当前代理治理与安全策略。">
    <el-skeleton v-if="loading && !settings" :rows="8" animated />
    <template v-else-if="settings">
      <el-row :gutter="16" class="mb-4">
        <el-col v-for="item in [
          { label: '待审核', value: settings.status_counts.pending, type: 'warning' },
          { label: '已通过', value: settings.status_counts.approved, type: 'success' },
          { label: '已驳回', value: settings.status_counts.rejected, type: 'danger' },
          { label: '已撤销', value: settings.status_counts.revoked, type: 'info' },
        ]" :key="item.label" :xs="12" :sm="6">
          <el-card shadow="never" class="mb-3"><div class="text-sm text-gray-500">{{ item.label }}</div><el-tag class="mt-2" :type="item.type" size="large">{{ item.value }} 人</el-tag></el-card>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12"><el-card shadow="never" header="审核与绑定"><el-descriptions :column="1" border><el-descriptions-item label="审核机制">{{ settings.review.platform_review_required ? '代理申请必须由平台审核' : '未启用' }}</el-descriptions-item><el-descriptions-item label="驳回说明">{{ settings.review.rejection_reason_required ? '驳回时必须填写原因' : '非必填' }}</el-descriptions-item><el-descriptions-item label="区域账号绑定">{{ settings.security.admin_binding_required ? '区域管理员必须一对一关联已通过代理' : '未启用' }}</el-descriptions-item></el-descriptions></el-card></el-col>
        <el-col :xs="24" :md="12"><el-card shadow="never" header="凭据与撤销保护"><el-descriptions :column="1" border><el-descriptions-item label="结算资料">{{ settings.security.payment_credentials_write_only ? '仅写入，不回显' : '—' }}</el-descriptions-item><el-descriptions-item label="后台密码">{{ settings.security.password_min_length }} 至 {{ settings.security.password_max_length }} 位；重置后旧会话立即失效</el-descriptions-item><el-descriptions-item label="撤销方式">{{ settings.revocation.hard_delete ? '硬删除' : '资格撤销，保留事实记录' }}</el-descriptions-item><el-descriptions-item label="撤销阻断">{{ settings.revocation.blocked_when.join('、') }}</el-descriptions-item></el-descriptions></el-card></el-col>
      </el-row>
    </template>
  </Page>
</template>
