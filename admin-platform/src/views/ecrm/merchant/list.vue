<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElDescriptions,
  ElDescriptionsItem,
  ElDrawer,
  ElEmpty,
  ElMessage,
  ElMessageBox,
} from 'element-plus';

import {
  fetchPlatformMerchant,
  fetchPlatformMerchants,
  updatePlatformMerchantStatus,
  type PlatformMerchantRow,
} from '#/api/core/ecrm';

const loading = ref(false);
const rows = ref<PlatformMerchantRow[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const detailLoading = ref(false);
const detail = ref<PlatformMerchantRow>();
const query = reactive({ keyword: '', page: 1, limit: 20, status: undefined as number | undefined });

const statusOptions = [
  { label: '全部状态', value: undefined },
  { label: '正常开启', value: 1 },
  { label: '已关闭', value: 0 },
];

const isEnabled = (row: PlatformMerchantRow) => row.status === 1 && row.mer_state === 1;
const statusText = (row: PlatformMerchantRow) => (isEnabled(row) ? '营业中' : '已关闭');
const auditText = (row: PlatformMerchantRow) => (row.is_audit === 1 ? '已审核' : '未审核');
const contact = (row: PlatformMerchantRow) => row.real_name || '—';
const phone = (row: PlatformMerchantRow) => row.mer_phone || '—';
const formatTime = (value?: string) => {
  if (!value) return '—';
  const date = new Date(value);
  return Number.isNaN(date.valueOf()) ? value : date.toLocaleString('zh-CN', { hour12: false });
};

const drawerTitle = computed(() => detail.value ? `店铺详情 · ${detail.value.mer_name}` : '店铺详情');

async function load() {
  loading.value = true;
  try {
    const result = await fetchPlatformMerchants({ ...query, keyword: query.keyword.trim() || undefined });
    rows.value = result.list;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.keyword = '';
  query.status = undefined;
  query.page = 1;
  void load();
}

function pageChange(page: number) {
  query.page = page;
  void load();
}

function limitChange(limit: number) {
  query.limit = limit;
  query.page = 1;
  void load();
}

async function changeStatus(row: PlatformMerchantRow, enabled: boolean) {
  const before = row.status === 1 && row.mer_state === 1;
  try {
    await ElMessageBox.confirm(
      `${enabled ? '开启' : '关闭'}后将立即影响该店铺的经营状态，是否继续？`,
      `${enabled ? '开启' : '关闭'}店铺`,
      { cancelButtonText: '取消', confirmButtonText: '确定', type: 'warning' },
    );
    await updatePlatformMerchantStatus(row.mer_id, enabled);
    row.status = enabled ? 1 : 0;
    row.mer_state = enabled ? 1 : 0;
    ElMessage.success('店铺状态已更新');
  } catch {
    row.status = before ? 1 : 0;
    row.mer_state = before ? 1 : 0;
  }
}

async function openDetail(row: PlatformMerchantRow) {
  detailOpen.value = true;
  detailLoading.value = true;
  detail.value = undefined;
  try {
    detail.value = await fetchPlatformMerchant(row.mer_id);
  } finally {
    detailLoading.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page
    title="店铺管理"
    description="集中查看店铺入驻、联系人和经营状态；启停操作会立即影响店铺经营。"
    content-class="merchant-list-page"
  >
    <el-card class="merchant-filter-card" shadow="never">
      <div class="section-heading">
        <div>
          <h2>筛选条件</h2>
          <p>按店铺信息和经营状态快速定位记录</p>
        </div>
      </div>
      <el-form class="merchant-filter-form" label-position="top" @submit.prevent="search">
        <el-form-item label="店铺搜索">
          <el-input
            v-model="query.keyword"
            clearable
            placeholder="店铺名称、联系人或手机号"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="经营状态">
          <el-select v-model="query.status" placeholder="全部状态">
            <el-option v-for="item in statusOptions" :key="String(item.value)" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item class="filter-actions" label=" ">
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="merchant-table-card" shadow="never">
      <div class="table-heading">
        <div>
          <h2>店铺列表</h2>
          <p>共 {{ total }} 家店铺</p>
        </div>
        <el-button plain @click="load">刷新列表</el-button>
      </div>
      <el-table class="merchant-table" v-loading="loading" :data="rows" row-key="mer_id">
        <el-table-column align="center" label="ID" prop="mer_id" width="80" />
        <el-table-column label="店铺名称" min-width="190" prop="mer_name" show-overflow-tooltip />
        <el-table-column label="联系人" min-width="130">
          <template #default="{ row }">{{ contact(row) }}</template>
        </el-table-column>
        <el-table-column label="联系电话" min-width="150">
          <template #default="{ row }">{{ phone(row) }}</template>
        </el-table-column>
        <el-table-column label="店铺地址" min-width="220" prop="mer_address" show-overflow-tooltip />
        <el-table-column align="center" label="入驻审核" width="108">
          <template #default="{ row }">
            <el-tag :type="row.is_audit === 1 ? 'success' : 'info'" effect="light">
              {{ auditText(row) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="经营状态" width="116">
          <template #default="{ row }">
            <el-switch
              :model-value="isEnabled(row)"
              active-text="启用"
              inactive-text="停用"
              inline-prompt
              width="52"
              @change="(enabled) => changeStatus(row, Boolean(enabled))"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="180">
          <template #default="{ row }">{{ formatTime(row.create_time) }}</template>
        </el-table-column>
        <el-table-column align="center" fixed="right" label="操作" width="94">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <ElEmpty v-if="!loading && rows.length === 0" description="暂无店铺数据" />
      <div class="table-footer">
        <el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="pageChange" @size-change="limitChange" />
      </div>
    </el-card>

    <ElDrawer v-model="detailOpen" :title="drawerTitle" size="520px">
      <div v-loading="detailLoading">
        <ElDescriptions v-if="detail" :column="1" border>
          <ElDescriptionsItem label="店铺 ID">{{ detail.mer_id }}</ElDescriptionsItem>
          <ElDescriptionsItem label="店铺名称">{{ detail.mer_name }}</ElDescriptionsItem>
          <ElDescriptionsItem label="联系人">{{ contact(detail) }}</ElDescriptionsItem>
          <ElDescriptionsItem label="联系电话">{{ phone(detail) }}</ElDescriptionsItem>
          <ElDescriptionsItem label="店铺地址">{{ detail.mer_address || '—' }}</ElDescriptionsItem>
          <ElDescriptionsItem label="经营状态">{{ statusText(detail) }}</ElDescriptionsItem>
          <ElDescriptionsItem label="审核状态">{{ auditText(detail) }}</ElDescriptionsItem>
          <ElDescriptionsItem label="备注">{{ detail.mark || '—' }}</ElDescriptionsItem>
          <ElDescriptionsItem label="店铺简介">{{ detail.mer_info || '—' }}</ElDescriptionsItem>
          <ElDescriptionsItem label="创建时间">{{ formatTime(detail.create_time) }}</ElDescriptionsItem>
        </ElDescriptions>
      </div>
    </ElDrawer>
  </Page>
</template>

<style scoped>
.merchant-list-page {
  min-width: 0;
}

.merchant-filter-card,
.merchant-table-card {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
}

.merchant-filter-card {
  margin-bottom: 16px;
}

.section-heading,
.table-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.section-heading {
  margin-bottom: 16px;
}

.section-heading h2,
.table-heading h2 {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

.section-heading p,
.table-heading p {
  margin: 4px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 20px;
}

.merchant-filter-form {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 12px 16px;
}

.merchant-filter-form :deep(.el-form-item) {
  margin: 0;
}

.merchant-filter-form :deep(.el-form-item__label) {
  height: auto;
  margin-bottom: 6px;
  color: var(--el-text-color-regular);
  font-size: 13px;
  line-height: 20px;
}

.merchant-filter-form :deep(.el-input),
.merchant-filter-form :deep(.el-select) {
  width: 260px;
}

.merchant-filter-form .filter-actions :deep(.el-form-item__label) {
  color: transparent;
}

.merchant-table-card :deep(.el-card__body) {
  padding: 16px 18px 14px;
}

.table-heading {
  align-items: center;
  margin-bottom: 14px;
}

.merchant-table {
  width: 100%;
}

.merchant-table :deep(.el-table__header th.el-table__cell) {
  height: 48px;
  color: var(--el-text-color-regular);
  font-size: 13px;
  font-weight: 600;
}

.merchant-table :deep(.el-table__cell) {
  height: 58px;
  padding: 8px 0;
}

.table-footer {
  display: flex;
  justify-content: flex-end;
  min-height: 40px;
  margin-top: 16px;
  overflow-x: auto;
}

.table-footer :deep(.el-pagination) {
  flex: 0 0 auto;
  flex-wrap: nowrap;
}

@media (max-width: 960px) {
  .merchant-filter-form :deep(.el-input),
  .merchant-filter-form :deep(.el-select) {
    width: min(100%, 320px);
  }

  .merchant-filter-form :deep(.el-form-item) {
    flex: 1 1 240px;
  }

  .merchant-filter-form .filter-actions {
    flex: 0 0 auto;
  }
}

@media (max-width: 640px) {
  .merchant-filter-card :deep(.el-card__body),
  .merchant-table-card :deep(.el-card__body) {
    padding: 14px;
  }

  .merchant-filter-form :deep(.el-form-item),
  .merchant-filter-form :deep(.el-input),
  .merchant-filter-form :deep(.el-select) {
    width: 100%;
  }

  .table-heading {
    align-items: flex-start;
  }
}
</style>
