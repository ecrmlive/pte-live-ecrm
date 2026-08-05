<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformDiscountsApi,
  setPlatformDiscountStatusApi,
  type PlatformDiscount,
} from '#/api/core/platform-discount';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<PlatformDiscount[]>([]);
const total = ref(0);
const canRead = ref(false);
const canManage = ref(false);
const query = reactive({
  keyword: '',
  limit: 20,
  page: 1,
  status: undefined as number | undefined,
  store_id: undefined as number | undefined,
});

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const page = await listPlatformDiscountsApi({
      keyword: query.keyword.trim() || undefined,
      limit: query.limit,
      page: query.page,
      status: query.status,
      store_id: query.store_id,
    });
    rows.value = page.list || [];
    total.value = page.total || 0;
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
  query.store_id = undefined;
  query.page = 1;
  void load();
}

async function setStatus(row: PlatformDiscount, status: 0 | 1) {
  const action = status === 1 ? '上架投影' : '下架投影';
  try {
    await ElMessageBox.confirm(
      `确认对「${row.name}」执行${action}？仅更新业务投影，不直连商户库。`,
      action,
      { type: 'warning' },
    );
    await setPlatformDiscountStatusApi(row.activity_id, status);
    ElMessage.success(`已${action}`);
    await load();
  } catch {
    /* cancel or toast */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const roleOK = profile.roles.some((role) => role === 'platform' || role === 'operations');
  canRead.value = roleOK && permissions.includes('marketing.discounts.read');
  canManage.value = roleOK && permissions.includes('marketing.discounts.manage');
  await load();
});
</script>

<template>
  <Page
    title="优惠套餐"
    description="只读监管店铺优惠套餐投影（qixi_crm_b_marketing_activity_view）；创建与编辑由商户后台完成。"
  >
    <el-alert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有优惠套餐监管权限"
      type="warning"
      :closable="false"
    />
    <EcrmListPage v-else title="优惠套餐列表">
      <template #filters>
        <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
          <el-form-item label="店铺 ID">
            <el-input-number v-model="query.store_id" :min="1" controls-position="right" />
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="query.keyword" clearable maxlength="64" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="query.status" clearable class="w-28" placeholder="全部">
              <el-option label="上架" :value="1" />
              <el-option label="下架" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="search">查询</el-button>
            <el-button @click="reset">重置</el-button>
          </el-form-item>
        </el-form>
      </template>

      <el-table v-loading="loading" :data="rows" row-key="activity_id">
        <el-table-column prop="activity_id" label="活动 ID" width="100" />
        <el-table-column prop="store_id" label="店铺 ID" width="100" />
        <el-table-column prop="name" label="名称" min-width="180" />
        <el-table-column label="套餐价" width="110">
          <template #default="{ row }">¥{{ Number(row.package_price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="商品数" width="90">
          <template #default="{ row }">{{ (row.product_ids || []).length }}</template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="有效期" min-width="220">
          <template #default="{ row }">
            {{ row.starts_at || '—' }} ~ {{ row.ends_at || '—' }}
          </template>
        </el-table-column>
        <el-table-column v-if="canManage" fixed="right" label="操作" width="160">
          <template #default="{ row }">
            <el-button
              v-if="row.status !== 1"
              link
              type="success"
              @click="setStatus(row, 1)"
            >
              上架投影
            </el-button>
            <el-button v-else link type="warning" @click="setStatus(row, 0)">下架投影</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #pager>
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          background
          layout="total, sizes, prev, pager, next"
          @current-change="(page: number) => { query.page = page; load(); }"
          @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }"
        />
      </template>
    </EcrmListPage>
  </Page>
</template>
