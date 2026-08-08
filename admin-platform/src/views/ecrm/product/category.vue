<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
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
  createPlatformCategoryApi,
  deletePlatformCategoryApi,
  listPlatformCategoriesApi,
  updatePlatformCategoryApi,
  updatePlatformCategoryRecommendApi,
  updatePlatformCategoryStatusApi,
  type PlatformCategory,
} from '#/api/core/platform-catalog';
import ImageField from '#/components/shop/image-field.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const treeRows = ref<PlatformCategory[]>([]);
const editing = ref<PlatformCategory>();
const form = reactive({
  cate_name: '',
  is_hot: 0,
  is_show: 1,
  pic: '',
  pid: 0,
  sort: 0,
});

const options = computed(() => flattenParents(treeRows.value));

function flattenParents(
  items: PlatformCategory[],
  prefix = '',
  depth = 0,
): Array<{ label: string; value: number }> {
  // 上级最多选到二级，保证可落到三级
  if (depth >= 2) return [];
  return items.flatMap((item) => [
    { label: `${prefix}${item.cate_name}`, value: item.store_category_id },
    ...flattenParents(item.children || [], `${prefix}— `, depth + 1),
  ]);
}

function filterTree(
  nodes: PlatformCategory[],
  keyword: string,
  status?: number,
): PlatformCategory[] {
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
    .filter((node): node is PlatformCategory => node !== null);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('分类名称'),
  LIST_ENABLE_STATUS_FIELD('显示状态'),
]);

const gridOptions: VxeGridProps<PlatformCategory> = {
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
    {
      align: 'center',
      field: 'pic',
      slots: { default: 'pic' },
      title: '分类图标',
      width: 100,
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
      align: 'center',
      field: 'is_hot',
      slots: { default: 'is_hot' },
      title: '是否推荐',
      width: 130,
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
        let list = (await listPlatformCategoriesApi()).list || [];
        treeRows.value = list;
        if (keyword || status !== undefined) {
          list = filterTree(list, keyword, status);
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'store_category_id' },
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
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm(parentID = 0) {
  editing.value = undefined;
  Object.assign(form, {
    cate_name: '',
    is_hot: 0,
    is_show: 1,
    pic: '',
    pid: parentID,
    sort: 0,
  });
}

function openCreate(parentID = 0) {
  resetForm(parentID);
  formDrawerApi.setState({ title: '新增商品分类' }).open();
}

function openEdit(row: PlatformCategory) {
  editing.value = row;
  Object.assign(form, {
    cate_name: row.cate_name,
    is_hot: row.is_hot ?? 0,
    is_show: row.is_show,
    pic: row.pic || '',
    pid: row.pid,
    sort: row.sort,
  });
  formDrawerApi.setState({ title: '编辑商品分类' }).open();
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
      is_hot: form.is_hot,
      is_show: form.is_show,
      pic: form.pic.trim(),
      pid: form.pid,
      sort: form.sort,
    };
    if (editing.value) {
      await updatePlatformCategoryApi(editing.value.store_category_id, body);
    } else {
      await createPlatformCategoryApi(body);
    }
    formDrawerApi.close();
    ElMessage.success('分类已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: PlatformCategory, enabled: boolean) {
  const before = row.is_show === 1;
  row.is_show = enabled ? 1 : 0;
  try {
    await updatePlatformCategoryStatusApi(row.store_category_id, enabled);
  } catch {
    row.is_show = before ? 1 : 0;
  }
}

async function changeHot(row: PlatformCategory, enabled: boolean) {
  const before = row.is_hot === 1;
  row.is_hot = enabled ? 1 : 0;
  try {
    await updatePlatformCategoryRecommendApi(row.store_category_id, enabled);
  } catch {
    row.is_hot = before ? 1 : 0;
  }
}

async function remove(row: PlatformCategory) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.cate_name}”后不可恢复，是否继续？`,
      '删除分类',
      { type: 'warning' },
    );
    await deletePlatformCategoryApi(row.store_category_id);
    ElMessage.success('分类已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      class="mb-3"
      type="warning"
      :closable="false"
      title="平台商品的分类应添加至三级，否则店铺添加商品时无分类可选"
    />
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate(0)">
          新增商品分类
        </ElButton>
      </template>

      <template #cate_name="{ row }">
        <span>{{ row.cate_name }} [{{ row.store_category_id }}]</span>
      </template>

      <template #pic="{ row }">
        <ElImage
          v-if="row.pic"
          :src="resolveCosMediaUrl(row.pic)"
          class="cate-icon"
          fit="cover"
          :preview-src-list="[resolveCosMediaUrl(row.pic)]"
          preview-teleported
        />
        <span v-else class="text-muted">—</span>
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

      <template #is_hot="{ row }">
        <ElSwitch
          :model-value="row.is_hot === 1"
          inline-prompt
          active-text="推荐"
          inactive-text="不推荐"
          @change="
            (enabled: string | number | boolean) =>
              changeHot(row, Boolean(enabled))
          "
        />
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="上级分类">
          <ElSelect v-model="form.pid" class="w-full" clearable placeholder="顶级分类">
            <ElOption label="顶级分类" :value="0" />
            <ElOption
              v-for="item in options.filter(
                (x) => x.value !== editing?.store_category_id,
              )"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="分类名称" required>
          <ElInput v-model="form.cate_name" placeholder="请输入分类名称" />
        </ElFormItem>
        <ElFormItem label="分类图标">
          <ImageField
            v-model="form.pic"
            default-library="system"
            hint="建议尺寸：110×110px"
            :preview-size="72"
          />
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
        <ElFormItem label="是否推荐">
          <ElSwitch
            v-model="form.is_hot"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="推荐"
            inactive-text="不推荐"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.cate-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
}

.text-muted {
  color: var(--el-text-color-placeholder);
}

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
