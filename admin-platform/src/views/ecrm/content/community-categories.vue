<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createCommunityCategoryApi,
  deleteCommunityCategoryApi,
  listCommunityCategoriesApi,
  updateCommunityCategoryApi,
  updateCommunityCategoryStatusApi,
  type CommunityCategory,
} from '#/api/core/platform-community';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';

const canRead = ref(false);
const canManage = ref(false);
const editing = ref<CommunityCategory>();
const form = reactive({
  cate_name: '',
  is_show: 1,
  sort: 0,
});

const gridOptions: VxeGridProps<CommunityCategory> = {
  columns: [
    {
      field: 'cate_name',
      minWidth: 200,
      showOverflow: false,
      title: '分类名称',
    },
    {
      field: 'sort',
      title: '排序',
      width: 100,
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '是否显示',
      width: 120,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const list = (await listCommunityCategoriesApi()).list || [];
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'category_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { cate_name: '', is_show: 1, sort: 0 });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加社区分类' }).open();
}

function openEdit(row: CommunityCategory) {
  editing.value = row;
  Object.assign(form, {
    cate_name: row.cate_name || '',
    is_show: row.is_show === 1 ? 1 : 0,
    sort: Number(row.sort || 0),
  });
  formDrawerApi.setState({ title: '编辑社区分类' }).open();
}

async function save() {
  const cateName = form.cate_name.trim();
  if (!cateName) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const sort = Math.max(0, Math.min(99999, Number(form.sort) || 0));
  formDrawerApi.lock();
  try {
    const payload = {
      cate_name: cateName,
      is_show: form.is_show === 1 ? 1 : 0,
      sort,
    };
    if (editing.value) {
      await updateCommunityCategoryApi(editing.value.category_id, payload);
    } else {
      await createCommunityCategoryApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: CommunityCategory, enabled: boolean) {
  if (!canManage.value) return;
  const before = row.is_show;
  row.is_show = enabled ? 1 : 0;
  try {
    await updateCommunityCategoryStatusApi(row.category_id, enabled ? 1 : 0);
  } catch {
    row.is_show = before;
  }
}

async function remove(row: CommunityCategory) {
  try {
    await confirm({
      content: `删除分类“${row.cate_name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteCommunityCategoryApi(row.category_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('content.community_category.manage');
  canRead.value =
    canManage.value ||
    permissions.includes('content.community_category.read');
  if (canRead.value) {
    gridApi.reload();
  }
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
          添加社区分类
        </ElButton>
      </template>
      <template #is_show="{ row }">
        <ElSwitch
          :model-value="row.is_show === 1"
          :disabled="!canManage"
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
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="primary" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="分类名称" required>
          <ElInput
            v-model="form.cate_name"
            maxlength="32"
            show-word-limit
            placeholder="请输入分类名称"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="99999"
            :precision="0"
            class="w-full"
          />
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
