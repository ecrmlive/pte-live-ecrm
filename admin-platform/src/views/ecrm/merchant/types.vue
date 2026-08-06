<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

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
  deleteMerchantType,
  fetchMerchantType,
  fetchMerchantTypes,
  saveMerchantType,
  setMerchantTypeRemark,
  setMerchantTypeStatus,
  type MerchantTypeRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const editing = ref<MerchantTypeRow>();
const form = reactive({
  name: '',
  type_info: '',
  is_margin: false,
  margin: 0,
  description: '',
  remark: '',
  status: true,
  menu_codes: [] as string[],
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '类型名称' },
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

const gridOptions: VxeGridProps<MerchantTypeRow> = {
  columns: [
    { field: 'name', minWidth: 150, showOverflow: false, title: '类型名称' },
    {
      field: 'type_info',
      minWidth: 220,
      showOverflow: false,
      title: '类型简介',
    },
    {
      field: 'margin',
      formatter: ({ row }) =>
        row.is_margin ? `¥${Number(row.margin).toFixed(2)}` : '不要求',
      title: '保证金',
      width: 130,
    },
    {
      field: 'menu_codes',
      formatter: ({ cellValue }) =>
        Array.isArray(cellValue) && cellValue.length
          ? cellValue.join('、')
          : '未配置',
      minWidth: 220,
      showOverflow: false,
      title: '授权菜单',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 240 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await fetchMerchantTypes()).list || [];
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

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function resetForm(row?: MerchantTypeRow) {
  editing.value = row;
  Object.assign(form, {
    name: row?.name || '',
    type_info: row?.type_info || '',
    is_margin: row?.is_margin === 1,
    margin: Number(row?.margin || 0),
    description: row?.description || '',
    remark: row?.remark || '',
    status: row ? row.status === 1 : true,
    menu_codes: row?.menu_codes || [],
  });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: '新增店铺类型' }).open();
}

async function openEdit(row: MerchantTypeRow) {
  resetForm(await fetchMerchantType(row.id));
  formModalApi.setState({ title: '编辑店铺类型' }).open();
}

async function save() {
  if (
    !form.name.trim() ||
    !form.description.trim() ||
    (form.is_margin && form.margin <= 0)
  ) {
    ElMessage.warning('请填写类型名称、说明；启用保证金时金额必须大于 0');
    return;
  }
  formModalApi.lock();
  try {
    await saveMerchantType(editing.value?.id, {
      ...form,
      name: form.name.trim(),
      type_info: form.type_info.trim(),
      description: form.description.trim(),
      remark: form.remark.trim(),
    });
    formModalApi.close();
    ElMessage.success('店铺类型已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function toggle(row: MerchantTypeRow) {
  await setMerchantTypeStatus(row.id, row.status !== 1);
  ElMessage.success('店铺类型状态已更新');
  gridApi.reload();
}

async function mark(row: MerchantTypeRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写不超过 500 字的内部备注。',
      `备注：${row.name}`,
      {
        inputValue: row.remark,
        inputPattern: /^[\s\S]{0,500}$/,
        inputErrorMessage: '备注不超过 500 字',
      },
    );
    await setMerchantTypeRemark(row.id, value.trim());
    ElMessage.success('备注已更新');
    gridApi.reload();
  } catch {
    /* 取消或统一请求层提示 */
  }
}

async function remove(row: MerchantTypeRow) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”将移除它的店铺菜单授权，是否继续？`,
      '删除店铺类型',
      { type: 'warning' },
    );
    await deleteMerchantType(row.id);
    ElMessage.success('店铺类型已删除');
    gridApi.reload();
  } catch {
    /* 取消或统一请求层提示 */
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value = codes.includes('merchant.type.manage');
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
          新增类型
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="primary" @click="mark(row)">备注</ElButton>
          <ElButton link type="primary" @click="toggle(row)">
            {{ row.status ? '停用' : '启用' }}
          </ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="106px">
        <ElFormItem label="类型名称" required>
          <ElInput v-model="form.name" maxlength="128" />
        </ElFormItem>
        <ElFormItem label="类型简介">
          <ElInput v-model="form.type_info" maxlength="500" />
        </ElFormItem>
        <ElFormItem label="要求保证金">
          <ElSwitch v-model="form.is_margin" />
        </ElFormItem>
        <ElFormItem v-if="form.is_margin" label="保证金金额" required>
          <ElInputNumber v-model="form.margin" :min="0.01" :precision="2" />
        </ElFormItem>
        <ElFormItem label="类型说明" required>
          <ElInput
            v-model="form.description"
            :rows="4"
            maxlength="65535"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="店铺菜单授权">
          <ElSelect
            v-model="form.menu_codes"
            multiple
            allow-create
            filterable
            default-first-option
            class="w-full"
            placeholder="输入统一菜单代码，如 merchant.catalog"
          >
            <ElOption label="商户控制台" value="merchant.dashboard" />
            <ElOption label="商品管理" value="merchant.catalog" />
            <ElOption label="订单管理" value="merchant.order" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="内部备注">
          <ElInput
            v-model="form.remark"
            :rows="2"
            maxlength="500"
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="启用状态">
          <ElSwitch v-model="form.status" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
