<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  claimCustomerServiceThread,
  fetchCustomerServiceThreads,
  type CustomerServiceThread,
} from '#/api/core/customer-service';

const loading = ref(false);
const rows = ref<CustomerServiceThread[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20, mine: false, status: undefined as 'closed' | 'open' | undefined });

function time(value?: string | null) {
  if (!value) return '—';
  const parsed = new Date(value);
  return Number.isNaN(parsed.valueOf()) ? value : parsed.toLocaleString('zh-CN', { hour12: false });
}

function assignment(row: CustomerServiceThread) {
  return row.assigned_admin_id ? `已领取（后台用户 #${row.assigned_admin_id}）` : '待领取';
}

async function load() {
  loading.value = true;
  try {
    const result = await fetchCustomerServiceThreads(query);
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
  query.page = 1;
  query.mine = false;
  query.status = undefined;
  void load();
}

async function claim(row: CustomerServiceThread) {
  try {
    await ElMessageBox.confirm('领取后该会话归入“仅看我的”，是否继续？', '领取客服会话', {
      confirmButtonText: '领取',
      cancelButtonText: '取消',
      type: 'warning',
    });
    const next = await claimCustomerServiceThread(row.id);
    Object.assign(row, next);
    ElMessage.success('客服会话已领取');
  } catch {
    // 请求客户端已经负责展示服务端错误；取消操作不提示。
  }
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

onMounted(load);
</script>

<template>
  <Page title="客服工作台" description="统一后台负责队列、领取和数据范围；会话消息、UserSig 与 WebSocket 由 pte-live-im 提供。">
    <el-alert class="mb-4" type="info" :closable="false" title="仅展示已获授权店铺的会话；领取不会改变 IM 消息内容或 IM SDK AppID。" />
    <el-card shadow="never">
      <el-form class="grid gap-x-4 md:grid-cols-3" label-width="72px" @submit.prevent="search">
        <el-form-item label="会话状态">
          <el-select v-model="query.status" clearable placeholder="全部状态">
            <el-option label="待服务 / 服务中" value="open" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="归属">
          <el-switch v-model="query.mine" active-text="仅看我的" inactive-text="全部队列" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="id">
        <el-table-column prop="id" label="会话编号" width="96" />
        <el-table-column prop="store_name" label="店铺" min-width="150">
          <template #default="{ row }">{{ row.store_name || `店铺 #${row.store_id}` }}</template>
        </el-table-column>
        <el-table-column prop="user_id" label="用户 ID" width="100" />
        <el-table-column label="关联订单" width="116">
          <template #default="{ row }">{{ row.order_id || '—' }}</template>
        </el-table-column>
        <el-table-column prop="im_conversation_id" label="IM 会话 ID" min-width="180" show-overflow-tooltip />
        <el-table-column label="当前 IM SDK" min-width="150">
          <template #default="{ row }">{{ row.im_sdk_app_id || '未配置' }}</template>
        </el-table-column>
        <el-table-column label="状态" width="104">
          <template #default="{ row }"><el-tag :type="row.status === 'open' ? 'warning' : 'info'">{{ row.status === 'open' ? '进行中' : '已关闭' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="领取状态" min-width="164">
          <template #default="{ row }">{{ assignment(row) }}</template>
        </el-table-column>
        <el-table-column label="更新时间" min-width="172">
          <template #default="{ row }">{{ time(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="92">
          <template #default="{ row }">
            <el-button v-if="row.status === 'open' && !row.assigned_admin_id" link type="primary" @click="claim(row)">领取</el-button>
            <span v-else>—</span>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="pageChange" @size-change="limitChange" />
      </div>
    </el-card>
  </Page>
</template>
