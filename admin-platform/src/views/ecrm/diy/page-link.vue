<script lang="ts" setup>
import type { DiyLinkScope, DiyPageCategory, DiyPageLink } from '#/api/core/diy';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElPagination,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
} from 'element-plus';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';

import {
  createDiyPageLinkApi,
  deleteDiyPageLinkApi,
  listDiyPageCategoriesApi,
  listDiyPageLinksApi,
  updateDiyPageLinkApi,
} from '#/api/core/diy';

const route = useRoute();
const scope = computed<DiyLinkScope>(() =>
  route.path.includes('/merLink/') ? 'merchant' : 'platform',
);
const rows = ref<DiyPageLink[]>([]);
const categories = ref<DiyPageCategory[]>([]);
const total = ref(0);
const page = ref(1);
const editingID = ref<number>();
const form = reactive({
  cate_id: 0,
  name: '',
  url: '',
  param: '',
  example: '',
  sort: 0,
  status: true,
});

function flatten(list: DiyPageCategory[], out: DiyPageCategory[] = []) {
  for (const item of list) {
    out.push(item);
    flatten(item.children || [], out);
  }
  return out;
}
const categoryOptions = computed(() => flatten(categories.value));

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => submit(),
});

async function load() {
  const [categoryRes, linkRes] = await Promise.all([
    listDiyPageCategoriesApi(scope.value),
    listDiyPageLinksApi(scope.value, { page: page.value, limit: 20 }),
  ]);
  categories.value = categoryRes.list || [];
  rows.value = linkRes.list || [];
  total.value = linkRes.total || 0;
}

function openCreate() {
  editingID.value = undefined;
  Object.assign(form, {
    cate_id: 0,
    name: '',
    url: '',
    param: '',
    example: '',
    sort: 0,
    status: true,
  });
  formDrawerApi.setState({ title: '新增链接' }).open();
}

function openEdit(row: DiyPageLink) {
  editingID.value = row.id;
  Object.assign(form, {
    cate_id: row.cate_id,
    name: row.name,
    url: row.url,
    param: row.param,
    example: row.example,
    sort: row.sort,
    status: row.status === 1,
  });
  formDrawerApi.setState({ title: '编辑链接' }).open();
}

async function submit() {
  if (!form.cate_id || !form.name.trim() || !form.url.trim()) {
    ElMessage.warning('请选择分类并填写页面名称、页面路径');
    return;
  }
  const body = {
    ...form,
    name: form.name.trim(),
    url: form.url.trim(),
    status: form.status ? 1 : 0,
  };
  formDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateDiyPageLinkApi(editingID.value, scope.value, body);
    } else {
      await createDiyPageLinkApi(scope.value, body);
    }
    ElMessage.success('已保存');
    formDrawerApi.close();
    await load();
  } finally {
    formDrawerApi.unlock();
  }
}

async function remove(row: DiyPageLink) {
  await ElMessageBox.confirm(`确认删除「${row.name}」？`, '删除链接', {
    type: 'warning',
  });
  await deleteDiyPageLinkApi(row.id, scope.value);
  ElMessage.success('已删除');
  await load();
}

onMounted(load);
</script>

<template>
  <Page auto-content-height>
    <div class="mb-4">
      <ElButton type="primary" @click="openCreate">新增链接</ElButton>
    </div>
    <ElTable :data="rows">
      <ElTableColumn prop="name" label="页面名称" min-width="150" />
      <ElTableColumn label="所属分类" min-width="140">
        <template #default="{ row }">{{ row.category?.name || '-' }}</template>
      </ElTableColumn>
      <ElTableColumn
        prop="url"
        label="页面路径"
        min-width="280"
        show-overflow-tooltip
      />
      <ElTableColumn prop="sort" label="排序" width="90" />
      <ElTableColumn label="状态" width="90">
        <template #default="{ row }">
          {{ row.status === 1 ? '启用' : '停用' }}
        </template>
      </ElTableColumn>
      <ElTableColumn label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
      </ElTableColumn>
    </ElTable>
    <div class="mt-4 flex justify-end">
      <ElPagination
        v-model:current-page="page"
        :page-size="20"
        layout="total, prev, pager, next"
        :total="total"
        @current-change="load"
      />
    </div>

    <FormDrawer>
      <ElForm label-width="92px">
        <ElFormItem label="所属分类">
          <ElSelect v-model="form.cate_id" class="w-full" filterable>
            <ElOption
              v-for="item in categoryOptions"
              :key="item.id"
              :value="item.id"
              :label="`${'— '.repeat(Math.max(0, item.level - 1))}${item.name}`"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="页面名称">
          <ElInput v-model="form.name" maxlength="50" />
        </ElFormItem>
        <ElFormItem label="页面路径">
          <ElInput v-model="form.url" placeholder="如 /pages/index/index" />
        </ElFormItem>
        <ElFormItem label="附加参数">
          <ElInput v-model="form.param" placeholder="可选" />
        </ElFormItem>
        <ElFormItem label="示例说明">
          <ElInput v-model="form.example" placeholder="可选" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :max="99999" />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch v-model="form.status" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
