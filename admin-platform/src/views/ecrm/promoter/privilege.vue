<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createDistributionPrivilegeApi,
  deleteDistributionPrivilegeApi,
  listDistributionPrivilegesApi,
  setDistributionPrivilegeStatusApi,
  updateDistributionPrivilegeApi,
  type DistributionPrivilege,
} from '#/api/core/platform-spread';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

type DrawerMode = 'create' | 'edit';

const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const form = reactive({
  title: '',
  img_url: '',
  sort: 0,
  status: 1 as 0 | 1,
});

const gridOptions: VxeGridProps<DistributionPrivilege> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'title',
      minWidth: 160,
      showOverflow: false,
      title: '标题',
    },
    {
      align: 'center',
      field: 'img_url',
      slots: { default: 'image' },
      title: '图片(90*90px)',
      width: 140,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '添加时间',
    },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await listDistributionPrivilegesApi({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
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
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editingId.value = 0;
  form.title = '';
  form.img_url = '';
  form.sort = 0;
  form.status = 1;
}

function openCreate() {
  drawerMode.value = 'create';
  resetForm();
  formDrawerApi.setState({ title: '新增分销特权', confirmText: '保存' }).open();
}

function openEdit(row: DistributionPrivilege) {
  drawerMode.value = 'edit';
  editingId.value = row.id;
  form.title = row.title || '';
  form.img_url = row.img_url || '';
  form.sort = row.sort ?? 0;
  form.status = row.status === 1 ? 1 : 0;
  formDrawerApi.setState({ title: '编辑分销特权', confirmText: '保存' }).open();
}

async function save() {
  const title = form.title.trim();
  if (!title) {
    ElMessage.warning('请输入标题');
    return;
  }
  if (!form.img_url.trim()) {
    ElMessage.warning('请选择图片');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      title,
      img_url: form.img_url.trim(),
      sort: form.sort ?? 0,
      status: form.status,
    };
    if (drawerMode.value === 'edit' && editingId.value) {
      await updateDistributionPrivilegeApi(editingId.value, payload);
    } else {
      await createDistributionPrivilegeApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: DistributionPrivilege, enabled: boolean) {
  const before = row.status === 1;
  row.status = enabled ? 1 : 0;
  try {
    await setDistributionPrivilegeStatusApi(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before ? 1 : 0;
  }
}

async function onDelete(row: DistributionPrivilege) {
  try {
    await confirm({
      content: `确定删除特权「${row.title}」吗？`,
      icon: 'warning',
      title: '提示',
    });
  } catch {
    return;
  }
  await deleteDistributionPrivilegeApi(row.id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增分销特权
        </ElButton>
      </template>
      <template #image="{ row }">
        <ElImage
          v-if="row.img_url"
          :src="resolveCosMediaUrl(row.img_url)"
          :preview-src-list="[resolveCosMediaUrl(row.img_url)]"
          fit="contain"
          style="width: 90px; height: 90px"
        />
        <span v-else>—</span>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="120px">
        <ElFormItem label="标题" required>
          <ElInput
            v-model="form.title"
            maxlength="20"
            placeholder="请输入标题"
          />
        </ElFormItem>
        <ElFormItem label="图片(90*90px)" required>
          <ImageField
            v-model="form.img_url"
            :preview-size="90"
            default-library="system"
            hint="建议尺寸 90×90px"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :precision="0" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="开启"
            inactive-text="关闭"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
