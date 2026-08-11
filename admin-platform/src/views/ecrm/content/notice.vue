<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createPlatformNoticeApi,
  deletePlatformNoticeApi,
  listPlatformNoticesApi,
  updatePlatformNoticeApi,
  type PlatformNotice,
} from '#/api/core/platform-content';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const editing = ref<PlatformNotice>();
const form = reactive({ content: '', is_show: 1, sort: 0, title: '' });

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('公告标题'),
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '展示', value: 1 },
        { label: '隐藏', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'is_show',
    label: '展示状态',
  },
]);

const gridOptions: VxeGridProps<PlatformNotice> = {
  columns: [
    { field: 'notice_id', title: 'ID', width: 80 },
    {
      field: 'title',
      minWidth: 220,
      showOverflow: false,
      title: '标题',
    },
    {
      field: 'content',
      minWidth: 260,
      showOverflow: false,
      title: '正文',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '展示',
      width: 90,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '发布时间',
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const showRaw = formValues?.is_show;
        const result = await listPlatformNoticesApi({
          page: page.currentPage,
          limit: page.pageSize,
        });
        let list = result.list || [];
        if (keyword) {
          list = list.filter(
            (row) =>
              row.title.toLowerCase().includes(keyword) ||
              row.content.toLowerCase().includes(keyword),
          );
        }
        if (showRaw === 0 || showRaw === 1) {
          list = list.filter((row) => row.is_show === Number(showRaw));
        }
        return {
          items: list,
          total: keyword || showRaw === 0 || showRaw === 1 ? list.length : result.total || 0,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'notice_id' },
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

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { content: '', is_show: 1, sort: 0, title: '' });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '发布公告' }).open();
}

function openEdit(row: PlatformNotice) {
  editing.value = row;
  Object.assign(form, {
    content: row.content,
    is_show: row.is_show,
    sort: row.sort,
    title: row.title,
  });
  formDrawerApi.setState({ title: '编辑公告' }).open();
}

async function save() {
  if (!form.title.trim() || !form.content.trim()) {
    ElMessage.warning('请填写公告标题和正文');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      ...form,
      content: form.content.trim(),
      title: form.title.trim(),
    };
    if (editing.value) {
      await updatePlatformNoticeApi(editing.value.notice_id, payload);
    } else {
      await createPlatformNoticeApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('公告已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function remove(row: PlatformNotice) {
  try {
    await ElMessageBox.confirm(
      `删除公告“${row.title}”后不可恢复，是否继续？`,
      '删除公告',
      { type: 'warning' },
    );
    await deletePlatformNoticeApi(row.notice_id);
    ElMessage.success('公告已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('content.notice.manage');
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
          发布公告
        </ElButton>
      </template>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '展示' : '隐藏' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="72px">
        <ElFormItem label="标题" required>
          <ElInput v-model="form.title" maxlength="100" show-word-limit />
        </ElFormItem>
        <ElFormItem label="正文" required>
          <ElInput v-model="form.content" :rows="10" type="textarea" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
        <ElFormItem label="展示">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
