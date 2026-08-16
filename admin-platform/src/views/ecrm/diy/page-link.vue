<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { DiyLinkScope, DiyPageCategory, DiyPageLink } from '#/api/core/diy';

import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElInputNumber, ElMessage, ElOption, ElSelect, ElSwitch } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { createDiyPageLinkApi, deleteDiyPageLinkApi, listDiyPageCategoriesApi, listDiyPageLinksApi, updateDiyPageLinkApi } from '#/api/core/diy';
import { platformListActionColumn, platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { LIST_ENABLE_STATUS_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const route = useRoute();
const scope = computed<DiyLinkScope>(() => route.path.includes('/merLink/') ? 'merchant' : 'platform');
const scopeLabel = computed(() => scope.value === 'platform' ? '平台页面链接' : '商户页面链接');
const categories = ref<DiyPageCategory[]>([]);
const editingID = ref<number>();
const form = reactive({ cate_id: 0, name: '', url: '', param: '', example: '', sort: 0, status: true });

function flatten(list: DiyPageCategory[], out: DiyPageCategory[] = []) {
  for (const item of list) {
    out.push(item);
    flatten(item.children || [], out);
  }
  return out;
}

const categoryOptions = computed(() => flatten(categories.value));

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

async function loadCategories() {
  const result = await listDiyPageCategoriesApi(scope.value);
  categories.value = result.list || [];
}

const formOptions = listFormOptionsDefaults([
  { component: 'Input', componentProps: { clearable: true, placeholder: '请输入页面名称或链接' }, fieldName: 'name', label: '页面搜索' },
  LIST_ENABLE_STATUS_FIELD('显示状态'),
]);

const gridOptions: VxeGridProps<DiyPageLink> = {
  columns: [
    { field: 'id', title: 'ID', width: 76 },
    { field: 'name', minWidth: 150, title: '页面名称' },
    { className: 'col--remark', field: 'url', minWidth: 220, showOverflow: 'tooltip', title: '页面链接' },
    { className: 'col--remark', field: 'param', formatter: ({ cellValue }) => cellValue || '—', minWidth: 160, showOverflow: 'tooltip', title: '参数' },
    { field: 'category', formatter: ({ row }) => row.category?.name || '—', minWidth: 130, title: '分组' },
    { field: 'add_time', formatter: ({ cellValue }) => formatTime(cellValue), minWidth: 175, title: '添加时间' },
    { field: 'status', slots: { default: 'status' }, title: '是否显示', width: 110 },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, values) => {
        const data = await listDiyPageLinksApi(scope.value, {
          page: page.currentPage,
          limit: page.pageSize,
          name: String(values?.name ?? '').trim() || undefined,
          status: values?.status === 0 || values?.status === 1 ? Number(values.status) : undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: { custom: false, export: false, refresh: false, search: false, zoom: false },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });
const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]', confirmText: '保存', cancelText: '取消', placement: 'right', onConfirm: async () => submit(),
});

function resetForm(row?: DiyPageLink) {
  editingID.value = row?.id;
  Object.assign(form, { cate_id: row?.cate_id || 0, name: row?.name || '', url: row?.url || '', param: row?.param || '', example: row?.example || '', sort: row?.sort || 0, status: row ? row.status === 1 : true });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: `添加${scopeLabel.value}` }).open();
}

function openEdit(row: DiyPageLink) {
  resetForm(row);
  formDrawerApi.setState({ title: `编辑${scopeLabel.value}` }).open();
}

async function submit() {
  if (!form.cate_id || !form.name.trim() || !form.url.trim()) {
    ElMessage.warning('请选择分组并填写页面名称、页面链接');
    return;
  }
  const body = { ...form, name: form.name.trim(), url: form.url.trim(), param: form.param.trim(), example: form.example.trim(), status: form.status ? 1 : 0 };
  formDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateDiyPageLinkApi(editingID.value, scope.value, body);
      ElMessage.success('页面链接已更新');
    } else {
      await createDiyPageLinkApi(scope.value, body);
      ElMessage.success('页面链接已添加');
    }
    formDrawerApi.close();
    await gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function updateStatus(row: DiyPageLink, enabled: boolean) {
  try {
    await updateDiyPageLinkApi(row.id, scope.value, {
      cate_id: row.cate_id,
      example: row.example,
      name: row.name,
      param: row.param,
      sort: row.sort,
      status: enabled ? 1 : 0,
      url: row.url,
    });
    row.status = enabled ? 1 : 0;
    ElMessage.success(enabled ? '已显示' : '已隐藏');
  } catch {
    row.status = enabled ? 0 : 1;
  }
}

async function remove(row: DiyPageLink) {
  try {
    await confirm({ cancelText: '取消', confirmText: '删除', content: `确认删除「${row.name}」吗？`, icon: 'warning', title: '删除页面链接' });
  } catch {
    return;
  }
  await deleteDiyPageLinkApi(row.id, scope.value);
  ElMessage.success('已删除');
  await gridApi.reload();
}

onMounted(async () => {
  await loadCategories();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">添加链接</ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch :model-value="row.status === 1" active-text="显示" inactive-text="隐藏" inline-prompt @change="(value) => updateStatus(row, Boolean(value))" />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>
    <FormDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="分组" required>
          <ElSelect v-model="form.cate_id" class="w-full" filterable>
            <ElOption v-for="item in categoryOptions" :key="item.id" :value="item.id" :label="`${'— '.repeat(Math.max(0, item.level - 1))}${item.name}`" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="页面名称" required><ElInput v-model="form.name" maxlength="50" show-word-limit /></ElFormItem>
        <ElFormItem label="页面链接" required><ElInput v-model="form.url" placeholder="如 /pages/index/index" /></ElFormItem>
        <ElFormItem label="参数"><ElInput v-model="form.param" placeholder="可选，例如 id=1" /></ElFormItem>
        <ElFormItem label="示例说明"><ElInput v-model="form.example" placeholder="可选" /></ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="form.sort" :min="0" :max="99999" /></ElFormItem>
        <ElFormItem label="是否显示"><ElSwitch v-model="form.status" active-text="显示" inactive-text="隐藏" inline-prompt /></ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
