<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
      <span class="hint">协议存于 qixi_cache；C 端 GET /agreements/:key</span>
    </div>
    <a-table row-key="key" :loading="loading" :columns="columns" :data-source="list" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'preview'">
          <span class="preview">{{ preview(record.content) }}</span>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canUpdate" type="link" @click="openEdit(record)">编辑</a-button>
          <span v-else class="hint">无保存权限</span>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editing?.label || '编辑协议'"
      :confirm-loading="saving"
      width="720px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="键">
          <a-input :value="editing?.key" disabled />
        </a-form-item>
        <a-form-item label="正文" required>
          <a-textarea v-model:value="form.content" :rows="12" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchAgreements, saveAgreement, type Agreement } from '@/api/content';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canUpdate = computed(() => auth.hasPerm('agreement/update'));

const loading = ref(false);
const saving = ref(false);
const list = ref<Agreement[]>([]);
const modalOpen = ref(false);
const editing = ref<Agreement | null>(null);
const form = reactive({ content: '' });

const columns = [
  { title: '名称', dataIndex: 'label', width: 160 },
  { title: '键', dataIndex: 'key', width: 220 },
  { title: '摘要', key: 'preview' },
  { title: '操作', key: 'action', width: 100 },
];

function preview(text: string) {
  const s = (text || '').replace(/\s+/g, ' ').trim();
  return s.length > 48 ? `${s.slice(0, 48)}…` : s || '（空）';
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAgreements();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

function reload() {
  void load();
}

function openEdit(row: Agreement) {
  editing.value = row;
  form.content = row.content || '';
  modalOpen.value = true;
}

async function submit() {
  if (!editing.value) return;
  if (!canUpdate.value) {
    message.warning('无保存权限');
    return;
  }
  if (!form.content.trim()) {
    message.warning('请填写正文');
    return;
  }
  saving.value = true;
  try {
    await saveAgreement(editing.value.key, form.content.trim());
    message.success('已保存');
    modalOpen.value = false;
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
  gap: 12px;
  margin-bottom: 16px;
}
.hint {
  color: #888;
  font-size: 13px;
}
.preview {
  color: #555;
}
</style>
