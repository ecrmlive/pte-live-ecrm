<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="product_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ record.price }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canConfig" type="link" @click="openConfig(record)">时段配置</a-button>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="open" title="预约时段" :confirm-loading="saving" width="640px" @ok="submit">
      <a-form layout="vertical">
        <a-form-item label="可约天数">
          <a-input-number v-model:value="days" :min="1" :max="30" style="width: 160px" />
        </a-form-item>
        <a-form-item label="时段（JSON 数组）">
          <a-textarea v-model:value="slotsText" :rows="8" placeholder='[{"start_time":"09:00","end_time":"11:00","stock":5}]' />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  fetchReservationConfig,
  fetchReservationProducts,
  saveReservationConfig,
  type ReservationProduct,
} from '@/api/reservation';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canConfig = computed(() => auth.hasPerm('reservation/config'));

const loading = ref(false);
const saving = ref(false);
const open = ref(false);
const list = ref<ReservationProduct[]>([]);
const editingId = ref(0);
const days = ref(7);
const slotsText = ref('[]');
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const columns = [
  { title: '商品ID', dataIndex: 'product_id', width: 90 },
  { title: '名称', dataIndex: 'store_name' },
  { title: '价格', key: 'price', width: 100 },
  { title: '库存', dataIndex: 'stock', width: 80 },
  { title: '操作', key: 'action', width: 120 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchReservationProducts({ page: pagination.current, limit: pagination.pageSize });
    list.value = data.list || [];
    pagination.total = data.total || 0;
  } finally {
    loading.value = false;
  }
}

function reload() {
  pagination.current = 1;
  void load();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  void load();
}

async function openConfig(row: ReservationProduct) {
  editingId.value = row.product_id;
  const { data } = await fetchReservationConfig(row.product_id);
  days.value = data.config?.show_reservation_days || 7;
  slotsText.value = JSON.stringify(data.slots || [], null, 2);
  open.value = true;
}

async function submit() {
  let slots: unknown[];
  try {
    slots = JSON.parse(slotsText.value || '[]');
    if (!Array.isArray(slots)) throw new Error('slots must be array');
  } catch {
    message.error('时段 JSON 无效');
    return;
  }
  saving.value = true;
  try {
    await saveReservationConfig(editingId.value, {
      show_reservation_days: days.value,
      slots,
    });
    message.success('已保存');
    open.value = false;
  } catch (e) {
    message.error((e as Error).message || '保存失败');
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void load();
});
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
