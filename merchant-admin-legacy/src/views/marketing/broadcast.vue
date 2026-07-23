<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">新建直播间</a-button>
      <a-button @click="reload">刷新</a-button>
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
          {{ liveText(record.live_status) }}
        </template>
        <template v-else-if="column.key === 'audit'">
          {{ auditText(record.status) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canLive" type="link" @click="setStatus(record, 101)">开播</a-button>
          <a-button v-if="canLive" type="link" @click="setStatus(record, 103)">结束</a-button>
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button v-if="canDelete" type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑直播间' : '新建直播间'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="直播间名称" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="主播昵称">
          <a-input v-model:value="form.anchor_name" />
        </a-form-item>
        <a-form-item v-if="canGoods || !editingId" label="挂货商品 ID（逗号分隔）">
          <a-input v-model:value="form.product_ids_text" placeholder="如 1,2" :disabled="!!editingId && !canGoods" />
        </a-form-item>
        <a-form-item v-if="canLive || !editingId" label="直播状态">
          <a-select v-model:value="form.live_status" style="width: 100%" :disabled="!!editingId && !canLive">
            <a-select-option :value="101">直播中</a-select-option>
            <a-select-option :value="102">未开始</a-select-option>
            <a-select-option :value="103">已结束</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createLiveRoom,
  deleteLiveRoom,
  fetchLiveRoom,
  fetchLiveRooms,
  setLiveStatus,
  updateLiveRoom,
  type LiveRoom,
} from '@/api/live';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canCreate = computed(() => auth.hasPerm('broadcast/create'));
const canDelete = computed(() => auth.hasPerm('broadcast/delete'));
const canLive = computed(() => auth.hasPerm('broadcast/live'));
const canGoods = computed(() => auth.hasPerm('broadcast/goods'));

const loading = ref(false);
const saving = ref(false);
const list = ref<LiveRoom[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  name: '',
  anchor_name: '店主',
  product_ids_text: '1',
  live_status: 102 as number,
});

const columns = [
  { title: 'ID', dataIndex: 'broadcast_room_id', width: 70 },
  { title: '直播间', dataIndex: 'name' },
  { title: '主播', dataIndex: 'anchor_name', width: 100 },
  { title: '直播', key: 'live', width: 90 },
  { title: '审核', key: 'audit', width: 90 },
  { title: '操作', key: 'action', width: 280 },
];

function liveText(s: number) {
  if (s === 101) return '直播中';
  if (s === 103) return '已结束';
  return '未开始';
}

function auditText(s: number) {
  if (s === 2) return '已通过';
  if (s === -1) return '已驳回';
  return '待审核';
}

function parseProductIDs(text: string) {
  return text
    .split(/[,，\s]+/)
    .map((x) => Number(x))
    .filter((n) => n > 0);
}

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

function openCreate() {
  editingId.value = 0;
  form.name = '';
  form.anchor_name = '店主';
  form.product_ids_text = '1';
  form.live_status = 102;
  modalOpen.value = true;
}

async function openEdit(row: LiveRoom) {
  editingId.value = row.broadcast_room_id;
  form.name = row.name;
  form.anchor_name = row.anchor_name || '';
  form.live_status = row.live_status;
  form.product_ids_text = (row.goods || []).map((g) => g.product_id).join(',');
  try {
    const { data } = await fetchLiveRoom(row.broadcast_room_id);
    form.product_ids_text = (data.goods || []).map((g) => g.product_id).join(',');
  } catch {
    /* keep list snapshot */
  }
  modalOpen.value = true;
}

async function submit() {
  saving.value = true;
  try {
    const body: Record<string, unknown> = {
      name: form.name,
      anchor_name: form.anchor_name,
    };
    if (!editingId.value || canLive.value) {
      body.live_status = form.live_status;
    }
    if (!editingId.value || canGoods.value) {
      body.product_ids = parseProductIDs(form.product_ids_text);
    }
    if (editingId.value) {
      await updateLiveRoom(editingId.value, body);
    } else {
      await createLiveRoom(body);
    }
    message.success('已保存（需平台审核后才对 C 端可见）');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function setStatus(row: LiveRoom, live_status: number) {
  await setLiveStatus(row.broadcast_room_id, live_status);
  message.success('已更新直播状态');
  void load();
}

async function onDelete(row: LiveRoom) {
  await deleteLiveRoom(row.broadcast_room_id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
