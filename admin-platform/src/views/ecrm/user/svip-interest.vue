<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRadio,
  ElRadioGroup,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createSvipInterest,
  deleteSvipInterest,
  listSvipInterests,
  updateSvipInterest,
  type SvipInterest,
  type SvipInterestInput,
} from '#/api/core/platform-svip-interest';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const saving = ref(false);
const editing = ref<SvipInterest>();
const form = reactive<SvipInterestInput>({
  name: '',
  description: '',
  icon_url: '',
  status: 1,
  sort: 0,
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    description: '',
    icon_url: '',
    status: 1,
    sort: 0,
    version: undefined,
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '权益名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<SvipInterest> = {
  columns: [
    { field: 'name', minWidth: 140, showOverflow: false, title: '权益名称' },
    {
      field: 'description',
      minWidth: 260,
      showOverflow: false,
      title: '权益说明',
    },
    {
      field: 'icon_url',
      formatter: ({ cellValue }) => cellValue || '—',
      title: '图标',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    { field: 'sort', title: '排序', width: 80 },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await listSvipInterests()).list || [];
        if (keyword) {
          list = list.filter((row) => row.name.toLowerCase().includes(keyword));
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return { items: list, total: list.length };
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增会员权益' }).open();
}

function openEdit(row: SvipInterest) {
  editing.value = row;
  Object.assign(form, {
    name: row.name,
    description: row.description,
    icon_url: row.icon_url,
    status: row.status,
    sort: row.sort,
    version: row.version,
  });
  formDrawerApi.setState({ title: '编辑会员权益' }).open();
}

async function save() {
  const icon = form.icon_url.trim();
  if (
    !form.name.trim() ||
    (icon && !icon.startsWith('/demo/') && !icon.startsWith('https://'))
  ) {
    ElMessage.warning('请填写权益名称；图标仅允许 /demo/ 或 https:// 地址');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const data = {
      ...form,
      name: form.name.trim(),
      description: form.description.trim(),
      icon_url: icon,
    };
    if (editing.value) await updateSvipInterest(editing.value.id, data);
    else await createSvipInterest(data);
    formDrawerApi.close();
    ElMessage.success('会员权益已保存；套餐只能选择启用权益');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function remove(row: SvipInterest) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”前会检查所有启用会员类型；如仍被使用，服务端将拒绝删除。`,
      '删除会员权益',
      { type: 'warning' },
    );
    await deleteSvipInterest(row.id);
    ElMessage.success('会员权益已逻辑删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    codes.includes('user.svip.interest.manage');
  if (canManage.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canManage"
      class="mb-4"
      title="当前账号没有会员权益维护权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增权益
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="88px">
        <ElFormItem label="权益名称" required>
          <ElInput v-model="form.name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="权益说明">
          <ElInput
            v-model="form.description"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="图标地址">
          <ElInput
            v-model="form.icon_url"
            maxlength="1024"
            placeholder="/demo/ 或 https://；可留空"
          />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">启用</ElRadio>
            <ElRadio :value="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
