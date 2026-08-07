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

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  deletePlatformSeckillApi,
  getPlatformSeckillApi,
  listPlatformSeckillApi,
  updatePlatformSeckillApi,
  type PlatformSeckillActive,
  type PlatformSeckillInput,
} from '#/api/core/platform-seckill';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  LIST_MER_ID_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const canManage = ref(false);
const editingID = ref<number>();
const form = reactive<Required<PlatformSeckillInput>>({
  name: '',
  seckill_time_ids: '',
  start_day: '',
  end_day: '',
  seckill_price: 0,
  once_pay_count: 1,
  status: 1,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('活动名称'),
  LIST_MER_ID_FIELD,
  LIST_ENABLE_STATUS_FIELD(),
]);

const gridOptions: VxeGridProps<PlatformSeckillActive> = {
  columns: [
    { field: 'seckill_active_id', title: 'ID', width: 80 },
    { field: 'name', minWidth: 160, showOverflow: false, title: '活动' },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'store_name',
      minWidth: 150,
      showOverflow: false,
      title: '商品',
      formatter: ({ row }) => row.store_name || `商品 #${row.product_id}`,
    },
    {
      field: 'seckill_price',
      title: '秒杀价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'start_day',
      minWidth: 200,
      showOverflow: false,
      title: '活动日期',
      formatter: ({ row }) => `${row.start_day} 至 ${row.end_day}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    platformListActionColumn({ width: 172 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformSeckillApi(buildStandardListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'seckill_active_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

async function edit(row: PlatformSeckillActive) {
  const detail = await getPlatformSeckillApi(row.seckill_active_id);
  editingID.value = row.seckill_active_id;
  Object.assign(form, {
    name: detail.name,
    seckill_time_ids: detail.seckill_time_ids || '',
    start_day: detail.start_day,
    end_day: detail.end_day,
    seckill_price: Number(detail.seckill_price),
    once_pay_count: detail.once_pay_count || 1,
    status: detail.status,
  });
  editDrawerApi.setState({ title: '编辑秒杀活动' }).open();
}

async function save() {
  if (
    !editingID.value ||
    !form.name.trim() ||
    !form.start_day ||
    !form.end_day ||
    form.end_day < form.start_day ||
    form.seckill_price <= 0 ||
    form.once_pay_count < 1
  ) {
    ElMessage.warning('请填写活动名称、有效日期、正数秒杀价和限购数量');
    return;
  }
  editDrawerApi.lock();
  saving.value = true;
  try {
    await updatePlatformSeckillApi(editingID.value, {
      ...form,
      name: form.name.trim(),
      seckill_time_ids: form.seckill_time_ids.trim(),
    });
    editDrawerApi.close();
    ElMessage.success('秒杀活动已更新');
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function setStatus(row: PlatformSeckillActive, status: number) {
  const action = status === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(`确认${action}秒杀活动“${row.name}”？`, `${action}确认`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await updatePlatformSeckillApi(row.seckill_active_id, { status });
    ElMessage.success(`活动已${action}`);
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function remove(row: PlatformSeckillActive) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”只会软删除活动配置，已产生订单不会被修改。是否继续？`,
      '删除秒杀活动',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' },
    );
    await deletePlatformSeckillApi(row.seckill_active_id);
    ElMessage.success('秒杀活动已软删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    permissions.includes('marketing.seckill.manage');
});
</script>

<template>
  <Page
    auto-content-height
    description="监管各商户秒杀活动；运营可维护活动配置或软删除，已产生订单、商品归属和历史价格快照不会被改写。"
    title="秒杀监管"
  >
    <Grid>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="edit(row)">编辑</ElButton>
          <ElButton
            link
            :type="row.status === 1 ? 'danger' : 'success'"
            @click="setStatus(row, row.status === 1 ? 0 : 1)"
          >
            {{ row.status === 1 ? '停用' : '启用' }}
          </ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <EditDrawer class="w-[620px]">
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="仅修改活动配置；商品、商户和已产生订单快照不可在此变更。"
      />
      <ElForm label-width="118px">
        <ElFormItem label="活动名称" required>
          <ElInput v-model="form.name" maxlength="128" show-word-limit />
        </ElFormItem>
        <ElFormItem label="秒杀场次 ID">
          <ElInput
            v-model="form.seckill_time_ids"
            placeholder="多个场次用英文逗号分隔，例如 1,2"
          />
        </ElFormItem>
        <ElFormItem label="活动日期" required>
          <ElInput v-model="form.start_day" class="!w-40" placeholder="YYYY-MM-DD" />
          <span class="mx-2">至</span>
          <ElInput v-model="form.end_day" class="!w-40" placeholder="YYYY-MM-DD" />
        </ElFormItem>
        <ElFormItem label="秒杀价" required>
          <ElInputNumber v-model="form.seckill_price" :min="0.01" :precision="2" :step="1" />
        </ElFormItem>
        <ElFormItem label="单次限购" required>
          <ElInputNumber v-model="form.once_pay_count" :min="1" :max="9999" />
        </ElFormItem>
        <ElFormItem label="活动状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">启用</ElRadio>
            <ElRadio :value="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
