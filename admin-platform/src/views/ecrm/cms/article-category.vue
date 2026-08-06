<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createArticleCategoryApi,
  deleteArticleCategoryApi,
  getArticleCategoryListApi,
  updateArticleCategoryApi,
  type ArticleCategoryOption,
} from '#/api/core/plus-article';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const editing = ref<ArticleCategoryOption>();
const form = reactive({ sort: 0, status: 1, title: '' });

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('分类名称'),
  LIST_ENABLE_STATUS_FIELD('状态'),
]);

const gridOptions: VxeGridProps<ArticleCategoryOption> = {
  columns: [
    { field: 'cid', title: 'ID', width: 80 },
    { field: 'title', minWidth: 180, title: '名称' },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await getArticleCategoryListApi()).list || [];
        if (keyword) {
          list = list.filter((row) =>
            row.title.toLowerCase().includes(keyword),
          );
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'cid' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { sort: 0, status: 1, title: '' });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: '新增分类' }).open();
}

function openEdit(row: ArticleCategoryOption) {
  editing.value = row;
  Object.assign(form, { sort: 0, status: row.status, title: row.title });
  formModalApi.setState({ title: '编辑分类' }).open();
}

async function save() {
  if (!form.title.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  formModalApi.lock();
  try {
    const payload = {
      sort: form.sort,
      status: form.status,
      title: form.title.trim(),
    };
    if (editing.value) {
      await updateArticleCategoryApi(editing.value.cid, payload);
    } else {
      await createArticleCategoryApi(payload);
    }
    formModalApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function remove(row: ArticleCategoryOption) {
  try {
    await ElMessageBox.confirm(
      `删除分类“${row.title}”后不可恢复，是否继续？`,
      '删除分类',
      { type: 'warning' },
    );
    await deleteArticleCategoryApi(row.cid);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('content.article_category.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增分类
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>
  </Page>
</template>
