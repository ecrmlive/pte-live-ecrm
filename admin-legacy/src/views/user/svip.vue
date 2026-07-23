<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
      <span class="hint">is_svip：0关闭 · 1体验 · 2有效期 · 3永久；用了会员价且商户未开叠加时禁用店铺券</span>
    </div>
    <a-table
      row-key="uid"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'svip'">
          {{ svipLabel(record.is_svip) }}
          <a-tag v-if="record.is_svip_active" color="success">有效</a-tag>
          <a-tag v-else>无效</a-tag>
        </template>
        <template v-else-if="column.key === 'end'">
          {{ record.svip_endtime || '—' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canUpdate" type="link" @click="openEdit(record)">设置</a-button>
          <span v-else class="muted">—</span>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="visible" title="设置付费会员" :confirm-loading="saving" @ok="submit">
      <a-form layout="vertical">
        <a-form-item label="会员类型">
          <a-select v-model:value="form.is_svip" style="width: 100%">
            <a-select-option :value="0">关闭</a-select-option>
            <a-select-option :value="1">体验</a-select-option>
            <a-select-option :value="2">有效期</a-select-option>
            <a-select-option :value="3">永久</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.is_svip === 2" label="到期时间">
          <a-input v-model:value="form.svip_endtime" placeholder="YYYY-MM-DD 或 YYYY-MM-DD HH:mm:ss" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchUsers, setUserSvip, type AppUserRow } from '@/api/svip';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canUpdate = computed(() => auth.hasPerm('svip/update'));

const loading = ref(false);
const saving = ref(false);
const visible = ref(false);
const list = ref<AppUserRow[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const form = reactive({ uid: 0, is_svip: 0, svip_endtime: '' });

const columns = [
  { title: 'UID', dataIndex: 'uid', width: 80 },
  { title: '账号', dataIndex: 'account', width: 120 },
  { title: '昵称', dataIndex: 'nickname' },
  { title: '手机', dataIndex: 'phone', width: 120 },
  { title: '会员', key: 'svip', width: 160 },
  { title: '到期', key: 'end', width: 180 },
  { title: '操作', key: 'action', width: 90 },
];

function svipLabel(v: number) {
  return ({ 0: '关闭', 1: '体验', 2: '有效期', 3: '永久' } as Record<number, string>)[v] || String(v);
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchUsers({ page: pagination.current, limit: pagination.pageSize });
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

function openEdit(row: AppUserRow) {
  form.uid = row.uid;
  form.is_svip = row.is_svip;
  form.svip_endtime = row.svip_endtime ? String(row.svip_endtime).slice(0, 19).replace('T', ' ') : '';
  visible.value = true;
}

async function submit() {
  saving.value = true;
  try {
    await setUserSvip(form.uid, {
      is_svip: form.is_svip,
      svip_endtime: form.is_svip === 2 ? form.svip_endtime : undefined,
    });
    message.success('已更新');
    visible.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}
.hint {
  color: #888;
  font-size: 13px;
}
.muted {
  color: #bbb;
}
</style>
