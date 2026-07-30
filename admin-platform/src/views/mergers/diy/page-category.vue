<script lang="ts" setup>
import type { DiyLinkScope, DiyPageCategory } from '#/api/core/diy';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
} from 'element-plus';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';

import {
  createDiyPageCategoryApi,
  deleteDiyPageCategoryApi,
  listDiyPageCategoriesApi,
  updateDiyPageCategoryApi,
} from '#/api/core/diy';

const route = useRoute();
const scope = computed<DiyLinkScope>(() =>
  route.path.includes('/merchant/') ? 'merchant' : 'platform',
);
const title = computed(() => (scope.value === 'merchant' ? '商户页面分类' : '平台页面分类'));
const rows = ref<DiyPageCategory[]>([]);
const dialogOpen = ref(false);
const editingID = ref<number>();
const form = reactive({ pid: 0, name: '', sort: 0, status: true });

function flatten(list: DiyPageCategory[], out: DiyPageCategory[] = []) {
  for (const item of list) {
    out.push(item);
    flatten(item.children || [], out);
  }
  return out;
}

const parentOptions = computed(() => flatten(rows.value).filter((item) => item.level < 3));

async function load() {
  const res = await listDiyPageCategoriesApi(scope.value);
  rows.value = res.list || [];
}

function openCreate() {
  editingID.value = undefined;
  Object.assign(form, { pid: 0, name: '', sort: 0, status: true });
  dialogOpen.value = true;
}

function openEdit(row: DiyPageCategory) {
  editingID.value = row.id;
  Object.assign(form, { pid: row.pid, name: row.name, sort: row.sort, status: row.status === 1 });
  dialogOpen.value = true;
}

async function submit() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const body = { ...form, name: form.name.trim(), status: form.status ? 1 : 0 };
  if (editingID.value) await updateDiyPageCategoryApi(editingID.value, scope.value, body);
  else await createDiyPageCategoryApi(scope.value, body);
  ElMessage.success('已保存');
  dialogOpen.value = false;
  await load();
}

async function remove(row: DiyPageCategory) {
  await ElMessageBox.confirm('仅允许删除没有子分类和链接的分类，确认继续？', '删除分类', { type: 'warning' });
  await deleteDiyPageCategoryApi(row.id, scope.value);
  ElMessage.success('已删除');
  await load();
}

onMounted(load);
</script>

<template>
  <Page :title="title" description="配置 DIY 装修链接的三级分类，供页面链接和装修组件选择。">
    <template #extra><ElButton type="primary" @click="openCreate">新增分类</ElButton></template>
    <ElTable :data="rows" row-key="id" default-expand-all :tree-props="{ children: 'children' }">
      <ElTableColumn prop="name" label="分类名称" min-width="240" />
      <ElTableColumn prop="level" label="层级" width="100" />
      <ElTableColumn prop="sort" label="排序" width="100" />
      <ElTableColumn label="状态" width="120"><template #default="{ row }">{{ row.status === 1 ? '启用' : '停用' }}</template></ElTableColumn>
      <ElTableColumn label="操作" width="180" fixed="right"><template #default="{ row }"><ElButton link type="primary" @click="openEdit(row)">编辑</ElButton><ElButton link type="danger" @click="remove(row)">删除</ElButton></template></ElTableColumn>
    </ElTable>
    <ElDialog v-model="dialogOpen" :title="editingID ? '编辑分类' : '新增分类'" width="480px" destroy-on-close>
      <ElForm label-width="92px">
        <ElFormItem label="上级分类"><ElSelect v-model="form.pid" class="w-full"><ElOption :value="0" label="顶级分类" /><ElOption v-for="item in parentOptions.filter((item) => item.id !== editingID)" :key="item.id" :value="item.id" :label="`${'— '.repeat(Math.max(0, item.level - 1))}${item.name}`" /></ElSelect></ElFormItem>
        <ElFormItem label="分类名称"><ElInput v-model="form.name" maxlength="50" show-word-limit /></ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="form.sort" :min="0" :max="99999" /></ElFormItem>
        <ElFormItem label="启用"><ElSwitch v-model="form.status" /></ElFormItem>
      </ElForm>
      <template #footer><ElButton @click="dialogOpen = false">取消</ElButton><ElButton type="primary" @click="submit">保存</ElButton></template>
    </ElDialog>
  </Page>
</template>
