<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

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
  ElSelect,
  ElSwitch,
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
import { formatShanghaiDateTime } from '#/utils/date-time';
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
      // 覆盖全局 grid.align:'center'：树按钮绝对定位在左，居中会把名称顶到中间形成大空隙
      align: 'left',
      field: 'cate_name',
      headerAlign: 'left',
      minWidth: 260,
      showOverflow: false,
      slots: { default: 'cate_name' },
      title: '分类名称',
      treeNode: true,
    },
    { align: 'center', field: 'sort', title: '排序', width: 90 },
    {
      align: 'center',
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
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
  treeConfig: {
    childrenField: 'children',
    expandAll: true,
    // 左对齐后 indent 只需区分层级；10px 三级仍可辨
    indent: 10,
  },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
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
  formDrawerApi.setState({ title: '新增品牌分类' }).open();
}

function openEdit(row: PlatformBrandCategory) {
  editing.value = row;
  Object.assign(form, {
    cate_name: row.cate_name,
    is_show: row.is_show,
    pid: row.pid,
    sort: row.sort,
  });
  formDrawerApi.setState({ title: '编辑品牌分类' }).open();
}

async function save() {
  if (!form.cate_name.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  formDrawerApi.lock();
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
    formDrawerApi.close();
    ElMessage.success('品牌分类已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: PlatformBrandCategory, enabled: boolean) {
  const before = row.is_show === 1;
  row.is_show = enabled ? 1 : 0;
  try {
    await updatePlatformBrandCategoryApi(row.brand_category_id, {
      cate_name: row.cate_name,
      is_show: enabled ? 1 : 0,
      pid: row.pid,
      sort: row.sort,
    });
  } catch {
    row.is_show = before ? 1 : 0;
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
          新增品牌分类
        </ElButton>
      </template>

      <template #cate_name="{ row }">
        <span>{{ row.cate_name }} [ {{ row.brand_category_id }} ]</span>
      </template>

      <template #is_show="{ row }">
        <ElSwitch
          :model-value="row.is_show === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeShow(row, Boolean(enabled))
          "
        />
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
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
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.is_show"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="显示"
            inactive-text="隐藏"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
/*
 * 根因：adapter 全局 align:'center' → .col--center 对 .vxe-cell--wrapper
 * justify-content:center，树按钮绝对定位在左、名称被整块居中，中间空白巨大。
 * 列已设 align/headerAlign:'left'；此处只收紧图标与文案间距。
 * （旧版类名 vxe-tree--btn-wrapper 在 v4.19 已变为 vxe-cell--tree-btn）
 */
:deep(.vxe-cell--tree-node) {
  text-align: left;
}

:deep(.vxe-cell--tree-btn) {
  width: 1em;
  height: 1em;
  margin: 0;
}

:deep(.vxe-tree-cell) {
  /* 默认 1.5em，略大于 1em 图标宽；收紧到紧贴箭头 */
  padding-left: 1em;
}
</style>
