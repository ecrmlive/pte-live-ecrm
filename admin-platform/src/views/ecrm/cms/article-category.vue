<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createArticleCategoryApi,
  deleteArticleCategoryApi,
  getArticleCategoryListApi,
  updateArticleCategoryApi,
  type ArticleCategoryOption,
} from '#/api/core/plus-article';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<ArticleCategoryOption[]>([]);
const open = ref(false);
const editing = ref<ArticleCategoryOption>();
const canManage = ref(false);
const form = reactive({ sort: 0, status: 1, title: '' });

async function load() {
  loading.value = true;
  try {
    rows.value = (await getArticleCategoryListApi()).list || [];
  } finally {
    loading.value = false;
  }
}

function add() {
  editing.value = undefined;
  Object.assign(form, { sort: 0, status: 1, title: '' });
  open.value = true;
}

function edit(row: ArticleCategoryOption) {
  editing.value = row;
  Object.assign(form, { sort: 0, status: row.status, title: row.title });
  open.value = true;
}

async function save() {
  if (!form.title.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const payload = { sort: form.sort, status: form.status, title: form.title.trim() };
  if (editing.value) await updateArticleCategoryApi(editing.value.cid, payload);
  else await createArticleCategoryApi(payload);
  open.value = false;
  ElMessage.success('已保存');
  await load();
}

async function remove(row: ArticleCategoryOption) {
  try {
    await ElMessageBox.confirm(`删除分类“${row.title}”后不可恢复，是否继续？`, '删除分类', { type: 'warning' });
    await deleteArticleCategoryApi(row.cid);
    ElMessage.success('已删除');
    await load();
  } catch {
    // 取消
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('content.article_category.manage');
});
</script>

<template>
  <Page title="文章分类" description="维护 CMS 文章分类；写入经 article API，不含密钥。">
    <EcrmListPage title="分类列表" description="分类供文章管理页引用。">
      <template #actions><el-button v-if="canManage" type="primary" @click="add">新增分类</el-button></template>
      <el-table v-loading="loading" :data="rows" row-key="cid">
        <el-table-column label="ID" prop="cid" width="80" />
        <el-table-column label="名称" min-width="180" prop="title" />
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="150"><template #default="{ row }"><el-button v-if="canManage" link type="primary" @click="edit(row)">编辑</el-button><el-button v-if="canManage" link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
    </EcrmListPage>
    <el-dialog v-model="open" :title="editing ? '编辑分类' : '新增分类'" width="520px" destroy-on-close>
      <el-form label-width="84px">
        <el-form-item label="名称" required><el-input v-model="form.title" maxlength="64" show-word-limit /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="open = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
