<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ArrowDown } from '@element-plus/icons-vue';
import {
  ElButton,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElIcon,
  ElMessage,
  ElMessageBox,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformAssistApi,
  updatePlatformAssistApi,
  type PlatformAssistActive,
} from '#/api/core/platform-assist';
import {
  getPlatformProductEditApi,
  updatePlatformProductAdminApi,
} from '#/api/core/platform-catalog';
import AssistDetailDrawer from '#/components/marketing/assist-detail-drawer.vue';
import AssistEditDrawer from '#/components/marketing/assist-edit-drawer.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import ProductLabelSelectModal from '#/views/ecrm/product/components/ProductLabelSelectModal.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD,
  LIST_MER_ID_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const detailDrawerRef = ref<InstanceType<typeof AssistDetailDrawer>>();
const editDrawerRef = ref<InstanceType<typeof AssistEditDrawer>>();
const labelModalRef = ref<InstanceType<typeof ProductLabelSelectModal>>();
const labelProductId = ref(0);

function time(value: string) {
  return formatShanghaiDateTime(value);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('活动 / 商品名称'),
  LIST_MER_ID_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已上架', value: 1 },
        { label: '已下架', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '展示状态',
  },
]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues);
  const statusRaw = formValues?.status;
  return {
    ...base,
    status: undefined,
    is_show: statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
  };
}

const gridOptions: VxeGridProps<PlatformAssistActive> = {
  columns: [
    { field: 'product_assist_id', title: '活动 ID', width: 96 },
    {
      field: 'store_name',
      minWidth: 180,
      showOverflow: false,
      title: '活动 / 商品',
      formatter: ({ row }) => row.store_name || `商品 #${row.product_id}`,
    },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'assist_price',
      title: '助力价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'assist_count',
      minWidth: 130,
      showOverflow: false,
      title: '助力规则',
      formatter: ({ row }) => `${row.assist_count} 人 / 每人最多 ${row.assist_user_count} 次`,
    },
    { field: 'stock', title: '活动库存', width: 100 },
    {
      field: 'start_time',
      minWidth: 240,
      showOverflow: false,
      title: '活动时间',
      formatter: ({ row }) => `${time(row.start_time)} 至 ${time(row.end_time)}`,
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '展示状态',
      width: 100,
    },
    platformListActionColumn({ width: 128 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformAssistApi(buildParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_assist_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function showDetail(row: PlatformAssistActive) {
  const assistId = Number(row.product_assist_id || 0);
  if (!assistId) {
    ElMessage.warning('缺少活动 ID，无法查看详情');
    return;
  }
  void detailDrawerRef.value?.open(assistId);
}

function openEdit(row: PlatformAssistActive) {
  const assistId = Number(row.product_assist_id || 0);
  if (!assistId) {
    ElMessage.warning('缺少活动 ID，无法编辑');
    return;
  }
  void editDrawerRef.value?.open(assistId);
}

function onEditSaved() {
  void gridApi.reload();
}

async function openLabels(row: PlatformAssistActive) {
  const productId = Number(row.product_id || 0);
  if (!productId) {
    ElMessage.warning('缺少关联商品，无法编辑标签');
    return;
  }
  labelProductId.value = productId;
  try {
    const edit = await getPlatformProductEditApi(productId);
    labelModalRef.value?.open({
      productId,
      selectedIds: [...(edit.sys_labels || [])].map(String),
    });
  } catch {
    labelModalRef.value?.open({ productId, selectedIds: [] });
  }
}

async function onLabelSubmit(ids: string[]) {
  if (!labelProductId.value) return;
  try {
    await updatePlatformProductAdminApi(labelProductId.value, {
      sys_labels: ids,
    } as Parameters<typeof updatePlatformProductAdminApi>[1]);
    ElMessage.success('标签已更新');
    void gridApi.reload();
  } catch {
    ElMessage.error('标签更新失败');
  }
}

function canForceOff(row: PlatformAssistActive) {
  return Number(row.product_status) === 1 && Number(row.is_show) === 1;
}

async function openForceOff(row: PlatformAssistActive) {
  try {
    await ElMessageBox.confirm(
      `确认强制下架好友助力活动“${row.store_name || `#${row.product_assist_id}`}”吗？`,
      '强制下架',
      { type: 'warning', confirmButtonText: '确认下架', cancelButtonText: '取消' },
    );
    // 平台助力监管目前仅开放 is_show；强制下架落为隐藏前台展示。
    await updatePlatformAssistApi(row.product_assist_id, { is_show: 0 });
    ElMessage.success('已强制下架');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

function onMoreCommand(command: string, row: PlatformAssistActive) {
  switch (command) {
    case 'edit':
      openEdit(row);
      break;
    case 'labels':
      void openLabels(row);
      break;
    case 'forceOff':
      if (canForceOff(row)) void openForceOff(row);
      break;
    default:
      break;
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    permissions.includes('marketing.assist.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '已上架' : '已下架' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="showDetail(row)">详情</ElButton>
        <ElDropdown
          v-if="canManage"
          trigger="click"
          @command="(cmd: string) => onMoreCommand(cmd, row)"
        >
          <ElButton link type="primary">
            更多
            <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem command="edit">编辑</ElDropdownItem>
              <ElDropdownItem command="labels">编辑标签</ElDropdownItem>
              <ElDropdownItem
                v-if="canForceOff(row)"
                command="forceOff"
                divided
              >
                强制下架
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <AssistDetailDrawer ref="detailDrawerRef" />
    <AssistEditDrawer ref="editDrawerRef" @saved="onEditSaved" />
    <ProductLabelSelectModal ref="labelModalRef" @submit="onLabelSubmit" />
  </Page>
</template>
