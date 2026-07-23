<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="createOpen = true">新建店铺券</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="coupon_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ Number(record.coupon_price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'min'">满 {{ record.use_min_price }} 可用</template>
        <template v-else-if="column.key === 'stock'">
          {{ record.is_limited === 1 ? `${record.remain_count}/${record.total_count}` : '不限' }}
        </template>
        <template v-else-if="column.key === 'status'">{{ record.status === 1 ? '开启' : '关闭' }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canToggle" type="link" @click="toggle(record)">
            {{ record.status === 1 ? '关闭' : '开启' }}
          </a-button>
          <a-popconfirm v-if="canDelete" title="确认删除？" @confirm="onDelete(record.coupon_id)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="createOpen" title="新建店铺券" :confirm-loading="saving" @ok="submitCreate">
      <a-form layout="vertical">
        <a-form-item label="标题" required>
          <a-input v-model:value="form.title" />
        </a-form-item>
        <a-form-item label="面额" required>
          <a-input-number v-model:value="form.coupon_price" :min="0.01" :step="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="最低消费（元取整）">
          <a-input-number v-model:value="form.use_min_price" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="领取后有效天数">
          <a-input-number v-model:value="form.coupon_time" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="发放总量（0=不限）">
          <a-input-number v-model:value="form.total_count" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { createCoupon, deleteCoupon, fetchCoupons, setCouponStatus, type Coupon } from '@/api/coupon';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canToggle = computed(() => auth.hasPerm('coupon/toggle'));
const canCreate = computed(() => auth.hasPerm('coupon/create'));
const canDelete = computed(() => auth.hasPerm('coupon/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<Coupon[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const createOpen = ref(false);
const form = reactive({
  title: '店铺满减券',
  coupon_price: 3,
  use_min_price: 20,
  coupon_time: 30,
  total_count: 500,
});

const columns = [
  { title: 'ID', dataIndex: 'coupon_id', width: 80 },
  { title: '标题', dataIndex: 'title' },
  { title: '面额', key: 'price', width: 100 },
  { title: '门槛', key: 'min', width: 120 },
  { title: '库存', key: 'stock', width: 120 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchCoupons({ page: pagination.current, limit: pagination.pageSize });
    list.value = data.list || [];
    pagination.total = data.total;
  } finally {
    loading.value = false;
  }
}

function reload() {
  pagination.current = 1;
  load();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

async function submitCreate() {
  saving.value = true;
  try {
    await createCoupon({
      title: form.title,
      coupon_price: form.coupon_price,
      use_min_price: form.use_min_price,
      coupon_time: form.coupon_time,
      total_count: form.total_count,
      is_limited: form.total_count > 0 ? 1 : 0,
    });
    message.success('已创建');
    createOpen.value = false;
    reload();
  } catch (e) {
    message.error((e as Error).message || '创建失败');
  } finally {
    saving.value = false;
  }
}

async function toggle(row: Coupon) {
  try {
    await setCouponStatus(row.coupon_id, row.status === 1 ? 0 : 1);
    message.success('已更新');
    load();
  } catch (e) {
    message.error((e as Error).message || '操作失败');
  }
}

async function onDelete(id: number) {
  try {
    await deleteCoupon(id);
    message.success('已删除');
    load();
  } catch (e) {
    message.error((e as Error).message || '删除失败');
  }
}

onMounted(load);
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
