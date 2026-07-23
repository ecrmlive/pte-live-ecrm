<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

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
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createArticle,
  deleteArticle,
  fetchArticles,
  updateArticle,
  type ArticleRow,
} from '#/api/core/mergers';

const form = reactive({
  title: '',
  author: '',
  synopsis: '',
  content: '',
  cid: 1,
  sort: 0,
  show: true,
});
const editingId = ref(0);

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => {
    if (!form.title.trim() || !form.content.trim()) {
      ElMessage.warning('请填写标题与内容');
      return;
    }
    const payload = {
      title: form.title.trim(),
      author: form.author.trim(),
      synopsis: form.synopsis.trim(),
      content: form.content,
      cid: form.cid,
      sort: form.sort,
      status: form.show ? 1 : 0,
    };
    if (editingId.value) await updateArticle(editingId.value, payload);
    else await createArticle(payload);
    ElMessage.success('已保存');
    formModalApi.close();
    gridApi.reload();
  },
});

const gridOptions: VxeGridProps<ArticleRow> = {
  border: true,
  columns: [
    { field: 'article_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 200, title: '标题' },
    { field: 'author', title: '作者', width: 120 },
    { field: 'visit', title: '浏览', width: 80 },
    {
      field: 'status',
      title: '状态',
      width: 90,
      formatter: ({ cellValue }) => (cellValue === 1 ? '展示' : '隐藏'),
    },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 160 },
  ],
  height: 'auto',
  pagerConfig: { enabled: true, pageSize: 20 },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await fetchArticles({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'article_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function openCreate() {
  editingId.value = 0;
  Object.assign(form, {
    title: '',
    author: '',
    synopsis: '',
    content: '',
    cid: 1,
    sort: 0,
    show: true,
  });
  formModalApi.setState({ title: '新建文章' });
  formModalApi.open();
}

function openEdit(row: ArticleRow) {
  editingId.value = row.article_id;
  Object.assign(form, {
    title: row.title,
    author: row.author,
    synopsis: row.synopsis,
    content: row.content,
    cid: row.cid,
    sort: row.sort,
    show: row.status === 1,
  });
  formModalApi.setState({ title: '编辑文章' });
  formModalApi.open();
}

async function onDelete(row: ArticleRow) {
  try {
    await ElMessageBox.confirm(`删除文章「${row.title}」？`, '提示', {
      type: 'warning',
    });
  } catch {
    return;
  }
  await deleteArticle(row.article_id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height title="文章管理">
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">新建文章</ElButton>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>
    <FormModal class="w-[720px]">
      <ElForm label-position="top">
        <ElFormItem label="标题" required>
          <ElInput v-model="form.title" />
        </ElFormItem>
        <ElFormItem label="作者">
          <ElInput v-model="form.author" />
        </ElFormItem>
        <ElFormItem label="简介">
          <ElInput v-model="form.synopsis" />
        </ElFormItem>
        <ElFormItem label="内容" required>
          <ElInput v-model="form.content" :rows="8" type="textarea" />
        </ElFormItem>
        <ElFormItem label="分类ID">
          <ElInputNumber v-model="form.cid" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="展示">
          <ElSwitch v-model="form.show" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
