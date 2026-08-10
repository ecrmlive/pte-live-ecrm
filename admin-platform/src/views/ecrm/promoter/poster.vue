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
  createDistributionPosterApi,
  deleteDistributionPosterApi,
  listDistributionPostersApi,
  setDistributionPosterStatusApi,
  updateDistributionPosterApi,
  type DistributionPoster,
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
  name: '',
  pic_url: '',
  sort: 0,
  status: 1 as 0 | 1,
});

const gridOptions: VxeGridProps<DistributionPoster> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'name',
      minWidth: 160,
      showOverflow: false,
      title: '名称',
    },
    {
      align: 'center',
      field: 'pic_url',
      slots: { default: 'image' },
      title: '背景图(600*1000px)',
      width: 170,
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
        const data = await listDistributionPostersApi({
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
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editingId.value = 0;
  form.name = '';
  form.pic_url = '';
  form.sort = 0;
  form.status = 1;
}

function openCreate() {
  drawerMode.value = 'create';
  resetForm();
  formDrawerApi.setState({ title: '添加数据', confirmText: '确定' }).open();
}

function openEdit(row: DistributionPoster) {
  drawerMode.value = 'edit';
  editingId.value = row.id;
  form.name = row.name || '';
  form.pic_url = row.pic_url || '';
  form.sort = row.sort ?? 0;
  form.status = row.status === 1 ? 1 : 0;
  formDrawerApi.setState({ title: '编辑数据', confirmText: '确定' }).open();
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请输入名称');
    return;
  }
  if (!form.pic_url.trim()) {
    ElMessage.warning('请选择背景图');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      name,
      pic_url: form.pic_url.trim(),
      sort: form.sort ?? 0,
      status: form.status,
    };
    if (drawerMode.value === 'edit' && editingId.value) {
      await updateDistributionPosterApi(editingId.value, payload);
    } else {
      await createDistributionPosterApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: DistributionPoster, enabled: boolean) {
  const before = row.status === 1;
  row.status = enabled ? 1 : 0;
  try {
    await setDistributionPosterStatusApi(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before ? 1 : 0;
  }
}

async function onDelete(row: DistributionPoster) {
  try {
    await confirm({
      content: `确定删除海报「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
  } catch {
    return;
  }
  await deleteDistributionPosterApi(row.id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加数据
        </ElButton>
      </template>
      <template #image="{ row }">
        <ElImage
          v-if="row.pic_url"
          :src="resolveCosMediaUrl(row.pic_url)"
          :preview-src-list="[resolveCosMediaUrl(row.pic_url)]"
          fit="cover"
          style="width: 48px; height: 80px"
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
      <ElForm label-width="150px">
        <ElFormItem label="名称" required>
          <ElInput
            v-model="form.name"
            maxlength="20"
            placeholder="请输入名称"
          />
        </ElFormItem>
        <ElFormItem label="背景图(600*1000px)" required>
          <ImageField
            v-model="form.pic_url"
            :preview-size="120"
            default-library="system"
            hint="建议尺寸 600×1000px"
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
