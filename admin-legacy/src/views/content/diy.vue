<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">新建装修页</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ record.status === 1 ? '启用中' : '未启用' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canUpdate" type="link" @click="openEdit(record)">编辑</a-button>
          <a-button v-if="canActive && record.status !== 1" type="link" @click="onActive(record)">启用</a-button>
          <a-button v-if="canDelete" type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑装修页' : '新建装修页'"
      :confirm-loading="saving"
      width="720px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="标题">
          <a-input v-model:value="form.title" />
        </a-form-item>
        <a-form-item label="模板名">
          <a-input v-model:value="form.template_name" placeholder="home" />
        </a-form-item>
        <a-form-item label="装修 JSON（banners / menus）" required>
          <div v-if="canPick" class="json-tools">
            <a-button size="small" @click="openPicker">从素材库选图</a-button>
            <span v-if="pickedUrl" class="picked">已选：{{ pickedUrl }}</span>
            <a-button v-if="pickedUrl" size="small" type="link" @click="applyPicked">写入首个空 banner</a-button>
          </div>
          <a-textarea v-model:value="form.valueJson" :rows="12" />
        </a-form-item>
        <a-form-item label="保存后启用">
          <a-switch v-model:checked="form.activate" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="pickerOpen" title="选择素材" width="640px" :footer="null">
      <a-spin :spinning="pickerLoading">
        <div v-if="pickerList.length" class="picker-grid">
          <div
            v-for="item in pickerList"
            :key="item.attachment_id"
            class="picker-item"
            @click="pick(item.attachment_src)"
          >
            <img :src="item.attachment_src" :alt="item.attachment_name" />
          </div>
        </div>
        <a-empty v-else description="暂无素材，请先到「内容 / 素材库」上传" />
      </a-spin>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  activateDiyPage,
  createDiyPage,
  deleteDiyPage,
  fetchDiyPages,
  updateDiyPage,
  type DiyPage,
} from '@/api/diy';
import { fetchAttachments, type Attachment } from '@/api/attachment';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canCreate = computed(() => auth.hasPerm('diy/create'));
const canUpdate = computed(() => auth.hasPerm('diy/update'));
const canDelete = computed(() => auth.hasPerm('diy/delete'));
const canActive = computed(() => auth.hasPerm('diy/active'));
const canPick = computed(() => auth.hasPerm('diy/pick'));

const defaultValue = {
  banners: [
    { id: 1, title: '夏日秒杀', image: '', url: '/pages/seckill/list' },
    { id: 2, title: '积分好物', image: '', url: '/pages/points/list' },
  ],
  menus: [
    { id: 1, name: '秒杀', icon: '', url: '/pages/seckill/list' },
    { id: 2, name: '积分', icon: '', url: '/pages/points/list' },
  ],
};

const loading = ref(false);
const saving = ref(false);
const list = ref<DiyPage[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  name: '',
  title: '',
  template_name: 'home',
  valueJson: JSON.stringify(defaultValue, null, 2),
  activate: true,
});
const pickerOpen = ref(false);
const pickerLoading = ref(false);
const pickerList = ref<Attachment[]>([]);
const pickedUrl = ref('');

async function openPicker() {
  pickerOpen.value = true;
  pickerLoading.value = true;
  try {
    const { data } = await fetchAttachments({ page: 1, limit: 48 });
    pickerList.value = data.list || [];
  } finally {
    pickerLoading.value = false;
  }
}

function pick(src: string) {
  pickedUrl.value = src;
  pickerOpen.value = false;
  message.success('已选择素材地址');
}

function applyPicked() {
  if (!pickedUrl.value) return;
  try {
    const obj = JSON.parse(form.valueJson || '{}') as {
      banners?: Array<{ image?: string }>;
    };
    if (!Array.isArray(obj.banners)) obj.banners = [];
    const target = obj.banners.find((b) => !b.image) || obj.banners[0];
    if (!target) {
      obj.banners.push({ image: pickedUrl.value });
    } else {
      target.image = pickedUrl.value;
    }
    form.valueJson = JSON.stringify(obj, null, 2);
    message.success('已写入 banner.image');
  } catch {
    message.warning('请先修正 JSON 再写入');
  }
}

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '名称', dataIndex: 'name' },
  { title: '标题', dataIndex: 'title' },
  { title: '模板', dataIndex: 'template_name', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'action', width: 220 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchDiyPages({ page: pagination.current, limit: pagination.pageSize });
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
  form.name = '平台首页';
  form.title = '栖息商城';
  form.template_name = 'home';
  form.valueJson = JSON.stringify(defaultValue, null, 2);
  form.activate = true;
  modalOpen.value = true;
}

function openEdit(row: DiyPage) {
  editingId.value = row.id;
  form.name = row.name;
  form.title = row.title;
  form.template_name = row.template_name || 'home';
  form.valueJson = row.value || JSON.stringify(defaultValue, null, 2);
  try {
    form.valueJson = JSON.stringify(JSON.parse(form.valueJson), null, 2);
  } catch {
    /* keep raw */
  }
  form.activate = row.status === 1;
  modalOpen.value = true;
}

async function submit() {
  if (!form.name.trim()) {
    message.warning('请填写名称');
    return;
  }
  let value: unknown;
  try {
    value = JSON.parse(form.valueJson || '{}');
  } catch {
    message.warning('装修 JSON 格式错误');
    return;
  }
  saving.value = true;
  const payload = {
    name: form.name.trim(),
    title: form.title.trim(),
    template_name: form.template_name.trim() || 'home',
    value,
    status: form.activate ? 1 : 0,
  };
  try {
    if (editingId.value) {
      await updateDiyPage(editingId.value, payload);
    } else {
      await createDiyPage(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function onActive(row: DiyPage) {
  await activateDiyPage(row.id);
  message.success('已启用');
  void load();
}

async function onDelete(row: DiyPage) {
  await deleteDiyPage(row.id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
.json-tools {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}
.picked {
  font-size: 12px;
  color: #666;
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.picker-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}
.picker-item {
  cursor: pointer;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  overflow: hidden;
}
.picker-item img {
  width: 100%;
  height: 90px;
  object-fit: cover;
  display: block;
}
.picker-item:hover {
  border-color: #1677ff;
}
</style>
