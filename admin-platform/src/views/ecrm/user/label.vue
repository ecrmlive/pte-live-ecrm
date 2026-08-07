<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElSelect,
  ElOption,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createUserLabel,
  deleteUserLabel,
  fetchUserLabels,
  markUserLabels,
  updateUserLabel,
  type UserLabelRow,
} from '#/api/core/ecrm';

const form = reactive({ label_name: '', sort: 0 });
const editingId = ref(0);
const markForm = reactive<{ uid: number; label_ids: number[] }>({
  uid: 1,
  label_ids: [],
});
const allLabels = ref<UserLabelRow[]>([]);

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (!form.label_name.trim()) {
      ElMessage.warning('请填写标签名');
      return;
    }
    const payload = { label_name: form.label_name.trim(), sort: form.sort };
    if (editingId.value) await updateUserLabel(editingId.value, payload);
    else await createUserLabel(payload);
    ElMessage.success('已保存');
    formDrawerApi.close();
    gridApi.reload();
  },
});

const [MarkDrawer, markDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (!markForm.uid) {
      ElMessage.warning('请填写用户 UID');
      return;
    }
    await markUserLabels(markForm.uid, markForm.label_ids);
    ElMessage.success('已打标');
    markDrawerApi.close();
  },
});

const gridOptions: VxeGridProps<UserLabelRow> = {
  border: true,
  columns: [
    { field: 'label_id', title: 'ID', width: 80 },
    { field: 'label_name', minWidth: 160, title: '标签名' },
    { field: 'sort', title: '排序', width: 80 },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 220 },
  ],
  height: 'auto',
  pagerConfig: { enabled: true, pageSize: 10 },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await fetchUserLabels({
          page: page.currentPage,
          limit: page.pageSize,
        });
        allLabels.value = data.list || [];
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'label_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function openCreate() {
  editingId.value = 0;
  form.label_name = '';
  form.sort = 0;
  formDrawerApi.setState({ title: '新建标签' });
  formDrawerApi.open();
}

function openEdit(row: UserLabelRow) {
  editingId.value = row.label_id;
  form.label_name = row.label_name;
  form.sort = row.sort;
  formDrawerApi.setState({ title: '编辑标签' });
  formDrawerApi.open();
}

function openMark(row: UserLabelRow) {
  markForm.uid = 1;
  markForm.label_ids = [row.label_id];
  markDrawerApi.setState({ title: '给用户打标' });
  markDrawerApi.open();
}

async function onDelete(row: UserLabelRow) {
  try {
    await ElMessageBox.confirm(`删除标签「${row.label_name}」？`, '提示', {
      type: 'warning',
    });
  } catch {
    return;
  }
  await deleteUserLabel(row.label_id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height title="用户标签">
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">新建标签</ElButton>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="openMark(row)">打标</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>
    <FormDrawer>
      <ElForm label-position="top">
        <ElFormItem label="标签名" required>
          <ElInput v-model="form.label_name" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
    <MarkDrawer>
      <ElForm label-position="top">
        <ElFormItem label="用户 UID" required>
          <ElInputNumber v-model="markForm.uid" :min="1" class="w-full" />
        </ElFormItem>
        <ElFormItem label="标签">
          <ElSelect v-model="markForm.label_ids" class="w-full" multiple>
            <ElOption
              v-for="item in allLabels"
              :key="item.label_id"
              :label="item.label_name"
              :value="item.label_id"
            />
          </ElSelect>
        </ElFormItem>
      </ElForm>
    </MarkDrawer>
  </Page>
</template>
