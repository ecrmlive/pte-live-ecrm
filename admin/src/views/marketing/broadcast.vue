<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
      <span class="hint">直播监管 · 审核通过后 C 端可见（无微信推流）</span>
    </div>
    <a-table
      row-key="broadcast_room_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'live'">
          {{ record.live_status === 101 ? '直播中' : record.live_status === 103 ? '已结束' : '未开始' }}
        </template>
        <template v-else-if="column.key === 'audit'">
          {{ record.status === 2 ? '已通过' : record.status === -1 ? '已驳回' : '待审核' }}
        </template>
        <template v-else-if="column.key === 'show'">
          {{ record.is_show === 1 ? '显示' : '隐藏' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canAudit" type="link" @click="audit(record, 2)">通过</a-button>
          <a-button v-if="canAudit" type="link" danger @click="audit(record, -1)">驳回</a-button>
          <a-button v-if="canAudit" type="link" @click="toggleShow(record)">
            {{ record.is_show === 1 ? '隐藏' : '显示' }}
          </a-button>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchLiveRooms, updateLiveRoomStatus, type LiveRoom } from '@/api/live';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canAudit = computed(() => auth.hasPerm('broadcast/audit'));

const loading = ref(false);
const list = ref<LiveRoom[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'broadcast_room_id', width: 70 },
  { title: '直播间', dataIndex: 'name' },
  { title: '商户', dataIndex: 'mer_name', width: 120 },
  { title: '主播', dataIndex: 'anchor_name', width: 100 },
  { title: '直播', key: 'live', width: 90 },
  { title: '审核', key: 'audit', width: 90 },
  { title: '展示', key: 'show', width: 80 },
  { title: '操作', key: 'action', width: 220 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchLiveRooms({ page: pagination.current, limit: pagination.pageSize });
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

async function audit(row: LiveRoom, status: number) {
  await updateLiveRoomStatus(row.broadcast_room_id, {
    status,
    refusal: status === -1 ? '不符合规范' : '',
    is_show: status === 2 ? 1 : 0,
  });
  message.success(status === 2 ? '已通过' : '已驳回');
  void load();
}

async function toggleShow(row: LiveRoom) {
  const next = row.is_show === 1 ? 0 : 1;
  await updateLiveRoomStatus(row.broadcast_room_id, { is_show: next });
  message.success(next === 1 ? '已显示' : '已隐藏');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}
.hint {
  color: #888;
  font-size: 13px;
}
</style>
