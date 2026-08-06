<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElSkeleton,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  auditPlatformProductApi,
  getPlatformProductApi,
  listPlatformProductsApi,
  type PlatformProduct,
} from '#/api/core/platform-catalog';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const current = ref<PlatformProduct>();
const detailLoading = ref(false);
const rejectOpen = ref(false);
const rejecting = ref(false);
const canAudit = ref(false);
const rejectForm = reactive({ refusal: '' });

function statusInfo(status: number) {
  return (
    (
      {
        '-2': { label: '已下架', type: 'info' },
        '-1': { label: '审核驳回', type: 'danger' },
        0: { label: '待审核', type: 'warning' },
        1: { label: '已通过', type: 'success' },
      } as Record<number, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }>
    )[status] || { label: '未知', type: 'info' as const }
  );
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const merIdRaw = String(formValues?.mer_id ?? '').trim();
  const statusRaw = formValues?.status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    mer_id: merIdRaw ? Number(merIdRaw) : undefined,
    status:
      statusRaw === 0 ||
      statusRaw === 1 ||
      statusRaw === -1 ||
      statusRaw === -2
        ? Number(statusRaw)
        : undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '名称 / 关键词' },
    fieldName: 'keyword',
    label: '商品搜索',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商户 ID' },
    fieldName: 'mer_id',
    label: '商户 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '已通过', value: 1 },
        { label: '已驳回', value: -1 },
        { label: '已下架', value: -2 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<PlatformProduct> = {
  columns: [
    { field: 'product_id', title: 'ID', width: 80 },
    {
      field: 'title',
      minWidth: 180,
      showOverflow: false,
      title: '商品名称',
    },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'cate_name',
      minWidth: 110,
      showOverflow: false,
      title: '分类',
      formatter: ({ cellValue }) => cellValue || '—',
    },
    {
      field: 'price',
      title: '售价',
      width: 108,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'stock',
      title: '库存/销量',
      width: 110,
      formatter: ({ row }) => `${row.stock} / ${row.sales}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 98,
    },
    {
      field: 'refusal',
      minWidth: 150,
      showOverflow: false,
      title: '驳回原因',
      formatter: ({ cellValue }) => cellValue || '—',
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 162 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformProductsApi(buildListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function openDetail(row: PlatformProduct) {
  current.value = undefined;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: '商品详情', loading: true }).open();
  try {
    current.value = await getPlatformProductApi(row.product_id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function reloadGrid() {
  await gridApi.reload();
}

async function approve(row: PlatformProduct) {
  try {
    await ElMessageBox.confirm(`确认审核通过商品“${row.title}”？`, '商品审核', {
      type: 'warning',
    });
    await auditPlatformProductApi(row.product_id, { status: 1 });
    ElMessage.success('商品已审核通过');
    await reloadGrid();
  } catch {
    /* 用户取消或统一请求错误提示。 */
  }
}

function openReject(row: PlatformProduct) {
  current.value = row;
  rejectForm.refusal = '';
  rejectOpen.value = true;
}

async function submitReject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  rejecting.value = true;
  try {
    await auditPlatformProductApi(current.value.product_id, {
      status: -1,
      refusal,
    });
    rejectOpen.value = false;
    ElMessage.success('商品已驳回');
    await reloadGrid();
  } finally {
    rejecting.value = false;
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canAudit.value =
    profile.roles.includes('platform') &&
    permissions.includes('product.audit.submit');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <template v-if="canAudit && row.status === 0">
          <ElButton link type="success" @click="approve(row)">通过</ElButton>
          <ElButton link type="danger" @click="openReject(row)">驳回</ElButton>
        </template>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="current">
            <ElDescriptions :column="1" border>
              <ElDescriptionsItem label="商品名称">
                {{ current.title }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="店铺">
                {{ current.store_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="商户">
                {{ current.mer_name || `商户 #${current.mer_id}` }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="分类">
                {{ current.cate_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="售价">
                ¥{{ Number(current.price).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="库存 / 销量">
                {{ current.stock }} / {{ current.sales }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="商品简介">
                {{ current.store_info || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="审核状态">
                <ElTag :type="statusInfo(current.status).type">
                  {{ statusInfo(current.status).label }}
                </ElTag>
              </ElDescriptionsItem>
              <ElDescriptionsItem v-if="current.refusal" label="驳回原因">
                {{ current.refusal }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="创建时间">
                {{ formatShanghaiDateTime(current.create_time) }}
              </ElDescriptionsItem>
            </ElDescriptions>
          </template>
        </template>
      </ElSkeleton>
    </DetailDrawer>

    <ElDialog
      v-model="rejectOpen"
      title="驳回商品"
      width="480px"
      destroy-on-close
    >
      <ElForm label-width="84px">
        <ElFormItem label="驳回原因" required>
          <ElInput
            v-model="rejectForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请向商户说明驳回原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="rejectOpen = false">取消</ElButton>
        <ElButton :loading="rejecting" type="danger" @click="submitReject">
          确认驳回
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>
