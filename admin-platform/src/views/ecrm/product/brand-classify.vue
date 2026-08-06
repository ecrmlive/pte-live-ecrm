<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformBrandCategoryApi,
  deletePlatformBrandCategoryApi,
  listPlatformBrandCategoriesApi,
  updatePlatformBrandCategoryApi,
  type PlatformBrandCategory,
} from '#/api/core/platform-catalog';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const treeRows = ref<PlatformBrandCategory[]>([]);
const editing = ref<PlatformBrandCategory>();
const form = reactive({ cate_name: '', is_show: 1, pid: 0, sort: 0 });

const flatParents = computed(() => flatten(treeRows.value));

function flatten(
  nodes: PlatformBrandCategory[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.brand_category_id },
    ...flatten(node.children || [], `${prefix}— `),
  ]);
}

function filterTree(
  nodes: PlatformBrandCategory[],
  keyword: string,
  status?: number,
): PlatformBrandCategory[] {
  return nodes
    .map((node) => {
      const children = node.children
        ? filterTree(node.children, keyword, status)
        : undefined;
      const nameMatch =
        !keyword || node.cate_name.toLowerCase().includes(keyword);
      const statusMatch =
        status !== 0 && status !== 1 ? true : node.is_show === status;
      if ((nameMatch && statusMatch) || (children && children.length)) {
        return { ...node, children };
      }
      return null;
    })
    .filter((node): node is PlatformBrandCategory => node !== null);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('分类名称'),
  LIST_ENABLE_STATUS_FIELD('显示状态'),
]);

const gridOptions: VxeGridProps<PlatformBrandCategory> = {
  columns: [
    {
      field: 'cate_name',
      minWidth: 220,
      showOverflow: false,
      title: '分类名称',
      treeNode: true,
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '状态',
      width: 100,
    },
    platformListActionColumn({ width: 220 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        const status =
          statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined;
        let list = (await listPlatformBrandCategoriesApi()).list || [];
        treeRows.value = list;
        if (keyword || status !== undefined) {
          list = filterTree(list, keyword, status);
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'brand_category_id' },
  treeConfig: { childrenField: 'children', expandAll: true },
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

function resetForm(parentID = 0) {
  editing.value = undefined;
  Object.assign(form, {
    cate_name: '',
    is_show: 1,
    pid: parentID,
    sort: 0,
  });
}

function openCreate(parentID = 0) {
  resetForm(parentID);
  formModalApi.setState({ title: '新增品牌分类' }).open();
}

function openEdit(row: PlatformBrandCategory) {
  editing.value = row;
  Object.assign(form, {
    cate_name: row.cate_name,
    is_show: row.is_show,
    pid: row.pid,
    sort: row.sort,
  });
  formModalApi.setState({ title: '编辑品牌分类' }).open();
}

async function save() {
  if (!form.cate_name.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  formModalApi.lock();
  try {
    const body = {
      cate_name: form.cate_name.trim(),
      is_show: form.is_show,
      pid: form.pid,
      sort: form.sort,
    };
    if (editing.value) {
      await updatePlatformBrandCategoryApi(
        editing.value.brand_category_id,
        body,
      );
    } else {
      await createPlatformBrandCategoryApi(body);
    }
    formModalApi.close();
    ElMessage.success('品牌分类已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function remove(row: PlatformBrandCategory) {
  try {
    await ElMessageBox.confirm(
      `删除分类“${row.cate_name}”？子分类与品牌须先清空。`,
      '删除品牌分类',
      { type: 'warning' },
    );
    await deletePlatformBrandCategoryApi(row.brand_category_id);
    ElMessage.success('品牌分类已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate(0)">
          新增分类
        </ElButton>
      </template>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openCreate(row.brand_category_id)">
          加子类
        </ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="88px">
        <ElFormItem label="上级分类">
          <ElSelect v-model="form.pid" clearable class="w-full" placeholder="顶级分类">
            <ElOption label="顶级分类" :value="0" />
            <ElOption
              v-for="item in flatParents.filter(
                (x) => x.value !== editing?.brand_category_id,
              )"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="分类名称" required>
          <ElInput v-model="form.cate_name" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="显示">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
