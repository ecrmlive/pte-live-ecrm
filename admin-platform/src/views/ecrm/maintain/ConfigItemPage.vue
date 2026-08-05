<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createConfigItemApi,
  deleteConfigItemApi,
  listConfigItemsApi,
  setConfigItemStatusApi,
  updateConfigItemApi,
  type ConfigItem,
  type ConfigItemType,
} from '#/api/core/platform-config-item';
import { EcrmFormDialog, EcrmListPage } from '#/components/ecrm';

const props = withDefaults(
  defineProps<{
    createLabel: string;
    description: string;
    itemType: ConfigItemType;
    manageCode: string;
    nameLabel?: string;
    readCode: string;
    roles?: string[];
    title: string;
  }>(),
  {
    nameLabel: '名称',
    roles: () => ['platform'],
  },
);

const loading = ref(false);
const saving = ref(false);
const rows = ref<ConfigItem[]>([]);
const total = ref(0);
const canRead = ref(false);
const canManage = ref(false);
const dialog = ref(false);
const editingID = ref<number>();
const query = reactive({
  keyword: '',
  limit: 20,
  page: 1,
  status: undefined as number | undefined,
});
const form = reactive({
  code: '',
  name: '',
  remark: '',
  sort: 0,
  status: 1,
});

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const page = await listConfigItemsApi(props.itemType, {
      keyword: query.keyword.trim() || undefined,
      limit: query.limit,
      page: query.page,
      status: query.status,
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
  query.page = 1;
  void load();
}

function openCreate() {
  editingID.value = undefined;
  Object.assign(form, { code: '', name: '', remark: '', sort: 0, status: 1 });
  dialog.value = true;
}

function openEdit(row: ConfigItem) {
  editingID.value = row.id;
  Object.assign(form, {
    code: row.code,
    name: row.name,
    remark: row.remark,
    sort: row.sort,
    status: row.status,
  });
  dialog.value = true;
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning(`请填写${props.nameLabel}`);
    return;
  }
  saving.value = true;
  try {
    const body = {
      code: form.code.trim(),
      name: form.name.trim(),
      payload: {},
      remark: form.remark.trim(),
      sort: form.sort,
      status: form.status,
    };
    if (editingID.value) {
      await updateConfigItemApi(props.itemType, editingID.value, body);
      ElMessage.success('已更新');
    } else {
      await createConfigItemApi(props.itemType, body);
      ElMessage.success('已创建');
    }
    dialog.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggleStatus(row: ConfigItem) {
  const next = row.status === 1 ? 0 : 1;
  const action = next === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(`确认${action}「${row.name}」？`, `${action}确认`, { type: 'warning' });
    await setConfigItemStatusApi(props.itemType, row.id, next as 0 | 1);
    ElMessage.success(`已${action}`);
    await load();
  } catch {
    /* cancel */
  }
}

async function remove(row: ConfigItem) {
  try {
    await ElMessageBox.confirm(`删除「${row.name}」后列表不再展示，确认继续？`, '删除确认', {
      confirmButtonText: '删除',
      type: 'warning',
    });
    await deleteConfigItemApi(props.itemType, row.id);
    ElMessage.success('已删除');
    await load();
  } catch {
    /* cancel */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const roleOK = profile.roles.some((role) => props.roles.includes(role));
  canRead.value = roleOK && permissions.includes(props.readCode);
  canManage.value = roleOK && permissions.includes(props.manageCode);
  await load();
});
</script>

<template>
  <Page :title="title" :description="description">
    <el-alert
      v-if="!canRead"
      class="mb-4"
      :title="`当前账号没有${title}权限`"
      type="warning"
      :closable="false"
    />
    <EcrmListPage v-else :title="title">
      <template #filters>
        <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
          <el-form-item label="关键词">
            <el-input v-model="query.keyword" clearable maxlength="64" placeholder="名称/标识" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="query.status" clearable class="w-28" placeholder="全部">
              <el-option label="启用" :value="1" />
              <el-option label="停用" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="search">查询</el-button>
            <el-button @click="reset">重置</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template #actions>
        <el-button v-if="canManage" type="primary" @click="openCreate">{{ createLabel }}</el-button>
      </template>

      <el-table v-loading="loading" :data="rows" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="nameLabel" min-width="160" />
        <el-table-column prop="code" label="标识" min-width="120" />
        <el-table-column prop="remark" label="备注" min-width="180" show-overflow-tooltip />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="updated_at" label="更新时间" min-width="170" />
        <el-table-column v-if="canManage" fixed="right" label="操作" width="180">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="toggleStatus(row)">
              {{ row.status === 1 ? '停用' : '启用' }}
            </el-button>
            <el-button link type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #pager>
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="(page: number) => { query.page = page; load(); }"
        />
      </template>
    </EcrmListPage>
  </Page>

  <EcrmFormDialog v-model="dialog" :title="editingID ? `编辑${title}` : createLabel" width="520px">
    <el-form label-width="84px">
      <el-form-item :label="nameLabel" required>
        <el-input v-model="form.name" maxlength="128" show-word-limit />
      </el-form-item>
      <el-form-item label="标识">
        <el-input v-model="form.code" maxlength="64" placeholder="可选业务标识" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="form.remark" type="textarea" :rows="2" maxlength="500" />
      </el-form-item>
      <el-form-item label="排序">
        <el-input-number v-model="form.sort" :step="1" />
      </el-form-item>
      <el-form-item label="启用">
        <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog = false">取消</el-button>
      <el-button :loading="saving" type="primary" @click="save">保存</el-button>
    </template>
  </EcrmFormDialog>
</template>
