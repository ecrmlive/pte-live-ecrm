<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import { ArrowDown } from '@element-plus/icons-vue';
import {
  ElAlert,
  ElButton,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElIcon,
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
  auditPlatformCombinationApi,
  forceOffPlatformCombinationApi,
  getPlatformCombinationApi,
  listPlatformCombinationsApi,
  updatePlatformCombinationApi,
  type PlatformCombination,
  type PlatformCombinationInput,
} from '#/api/core/platform-combination';
import {
  getPlatformProductEditApi,
  updatePlatformProductAdminApi,
} from '#/api/core/platform-catalog';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import ProductLabelSelectModal from '#/views/ecrm/product/components/ProductLabelSelectModal.vue';
import ProductPreviewModal from '#/views/ecrm/product/components/ProductPreviewModal.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';
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
const detailRow = ref<(PlatformCombination & { time?: number }) | null>(null);
const labelProductId = ref(0);
const forceOffIds = ref<number[]>([]);
const forceOffReason = ref('');
const previewProductId = ref(0);
const previewProductTitle = ref('');
const previewDisplayPrice = ref<number>();
const previewDisplayOtPrice = ref<number>();
const auditID = ref(0);
const auditForm = reactive({ status: 1 as 1 | -1, refusal: '' });
const labelModalRef = ref<InstanceType<typeof ProductLabelSelectModal>>();
const previewModalRef = ref<InstanceType<typeof ProductPreviewModal>>();

const form = reactive<Required<PlatformCombinationInput>>({
  price: 0,
  buying_count_num: 2,
  time: 24,
  start_time: '',
  end_time: '',
  is_show: 1,
  status: 1,
  product_status: 1,
  refusal: '',
});

function dateTime(value: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('商品名称'),
  LIST_MER_ID_FIELD,
  LIST_ENABLE_STATUS_FIELD(),
]);

const gridOptions: VxeGridProps<PlatformCombination> = {
  columns: [
    { field: 'product_group_id', title: 'ID', width: 80 },
    {
      field: 'store_name',
      minWidth: 180,
      showOverflow: false,
      title: '商品',
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
      field: 'price',
      title: '拼团价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'buying_count_num', title: '成团人数', width: 100 },
    {
      field: 'start_time',
      minWidth: 220,
      showOverflow: false,
      title: '活动时间',
      formatter: ({ row }) =>
        `${dateTime(row.start_time)} 至 ${dateTime(row.end_time)}`,
    },
    {
      field: 'product_status',
      slots: { default: 'auditStatus' },
      title: '审核',
      width: 88,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    platformListActionColumn({ width: 168 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformCombinationsApi(
          buildStandardListParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_group_id' },
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

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '拼团商品详情',
});

const [AuditDrawer, auditDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '提交审核',
  cancelText: '取消',
  placement: 'right',
  title: '拼团审核',
  onConfirm: async () => submitAudit(),
});

const [ForceOffModal, forceOffModalApi] = useVbenModal({
  title: '强制下架',
  class: 'w-[520px] max-w-[96vw]',
  confirmText: '确认下架',
  onConfirm: () => void submitForceOff(),
});

function toFormDateTime(value: string) {
  return value ? value.replace('T', ' ').slice(0, 19) : '';
}

function auditLabel(status?: number) {
  switch (Number(status)) {
    case 0:
      return '待审';
    case 1:
      return '通过';
    case -1:
      return '未通过';
    case -2:
      return '强制下架';
    default:
      return '—';
  }
}

function auditTagType(status?: number) {
  switch (Number(status)) {
    case 0:
      return 'warning';
    case 1:
      return 'success';
    case -1:
      return 'danger';
    case -2:
      return 'info';
    default:
      return 'info';
  }
}

function canAudit(row: PlatformCombination) {
  return Number(row.product_status) === 0;
}

function canForceOff(row: PlatformCombination) {
  return Number(row.product_status) === 1 && Number(row.is_show ?? 1) === 1;
}

async function openDetail(row: PlatformCombination) {
  try {
    detailRow.value = await getPlatformCombinationApi(row.product_group_id);
    detailDrawerApi.open();
  } catch {
    ElMessage.error('加载详情失败');
  }
}

async function openEdit(row: PlatformCombination) {
  const detail = await getPlatformCombinationApi(row.product_group_id);
  editingID.value = row.product_group_id;
  Object.assign(form, {
    price: Number(detail.price),
    buying_count_num: detail.buying_count_num,
    time: detail.time || 24,
    start_time: toFormDateTime(detail.start_time),
    end_time: toFormDateTime(detail.end_time),
    is_show: detail.is_show,
    status: detail.status,
    product_status: detail.product_status ?? 1,
    refusal: detail.refusal || '',
  });
  editDrawerApi.setState({ title: '编辑拼团活动' }).open();
}

function openPreview(row: PlatformCombination) {
  const productId = Number(row.product_id || 0);
  if (!productId) {
    ElMessage.warning('缺少关联商品，无法预览');
    return;
  }
  previewProductId.value = productId;
  previewProductTitle.value = row.store_name || '';
  previewDisplayPrice.value = Number(row.price || 0);
  previewDisplayOtPrice.value =
    row.ot_price !== undefined && row.ot_price !== null
      ? Number(row.ot_price)
      : undefined;
  previewModalRef.value?.open();
}

function openAudit(row: PlatformCombination) {
  auditID.value = row.product_group_id;
  auditForm.status = 1;
  auditForm.refusal = '';
  auditDrawerApi.open();
}

async function openLabels(row: PlatformCombination) {
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

function openForceOff(row: PlatformCombination) {
  forceOffIds.value = [row.product_group_id];
  forceOffReason.value = '';
  forceOffModalApi.open();
}

function onMoreCommand(command: string, row: PlatformCombination) {
  switch (command) {
    case 'edit':
      void openEdit(row);
      break;
    case 'audit':
      if (canAudit(row)) openAudit(row);
      break;
    case 'labels':
      void openLabels(row);
      break;
    case 'forceOff':
      if (canForceOff(row)) openForceOff(row);
      break;
    default:
      break;
  }
}

async function save() {
  if (
    !editingID.value ||
    form.price <= 0 ||
    form.buying_count_num < 2 ||
    form.time < 1 ||
    !form.start_time ||
    !form.end_time ||
    new Date(form.end_time).valueOf() < new Date(form.start_time).valueOf()
  ) {
    ElMessage.warning(
      '请填写正数拼团价、至少 2 人、有效时长和正确的活动时间',
    );
    return;
  }
  editDrawerApi.lock();
  saving.value = true;
  try {
    await updatePlatformCombinationApi(editingID.value, {
      price: form.price,
      buying_count_num: form.buying_count_num,
      time: form.time,
      start_time: form.start_time,
      end_time: form.end_time,
      is_show: form.is_show,
      status: form.status,
    });
    editDrawerApi.close();
    ElMessage.success('拼团活动已更新');
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function submitAudit() {
  if (!auditID.value) return;
  if (auditForm.status === -1 && !auditForm.refusal.trim()) {
    ElMessage.warning('请填写拒绝原因');
    return;
  }
  auditDrawerApi.lock();
  try {
    await auditPlatformCombinationApi(
      auditID.value,
      auditForm.status,
      auditForm.status === -1 ? auditForm.refusal.trim() : '',
    );
    ElMessage.success(auditForm.status === 1 ? '已通过审核' : '已拒绝');
    auditDrawerApi.close();
    gridApi.reload();
  } catch {
    ElMessage.error('提交审核失败');
  } finally {
    auditDrawerApi.unlock();
  }
}

async function submitForceOff() {
  if (!forceOffReason.value.trim()) {
    ElMessage.warning('请填写下架原因');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `确认强制下架选中的 ${forceOffIds.value.length} 个拼团商品？`,
      '强制下架',
      { type: 'warning', confirmButtonText: '确认下架', cancelButtonText: '取消' },
    );
    forceOffModalApi.lock();
    await forceOffPlatformCombinationApi(
      forceOffIds.value,
      forceOffReason.value.trim(),
    );
    ElMessage.success('已强制下架');
    forceOffModalApi.close();
    gridApi.reload();
  } catch {
    /* 用户取消或失败 */
  } finally {
    forceOffModalApi.unlock();
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    permissions.includes('marketing.combination.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #auditStatus="{ row }">
        <ElTag :type="auditTagType(row.product_status)" size="small">
          {{ auditLabel(row.product_status) }}
        </ElTag>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openPreview(row)">预览</ElButton>
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
              <ElDropdownItem v-if="canAudit(row)" command="audit">
                审核
              </ElDropdownItem>
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

    <DetailDrawer>
      <template v-if="detailRow">
        <ElForm label-width="118px">
          <ElFormItem label="活动 ID">
            {{ detailRow.product_group_id }}
          </ElFormItem>
          <ElFormItem label="商品">
            {{ detailRow.store_name || `商品 #${detailRow.product_id}` }}
          </ElFormItem>
          <ElFormItem label="商户">
            {{ detailRow.mer_name || `商户 #${detailRow.mer_id}` }}
          </ElFormItem>
          <ElFormItem label="拼团价">
            ¥{{ Number(detailRow.price || 0).toFixed(2) }}
          </ElFormItem>
          <ElFormItem label="成团人数">
            {{ detailRow.buying_count_num }}
          </ElFormItem>
          <ElFormItem label="成团时限">
            {{ detailRow.time || '—' }} 小时
          </ElFormItem>
          <ElFormItem label="活动时间">
            {{ dateTime(detailRow.start_time) }} 至
            {{ dateTime(detailRow.end_time) }}
          </ElFormItem>
          <ElFormItem label="审核状态">
            <ElTag :type="auditTagType(detailRow.product_status)" size="small">
              {{ auditLabel(detailRow.product_status) }}
            </ElTag>
          </ElFormItem>
          <ElFormItem v-if="detailRow.refusal" label="原因">
            {{ detailRow.refusal }}
          </ElFormItem>
          <ElFormItem label="前台展示">
            {{ detailRow.is_show === 1 ? '上架' : '下架' }}
          </ElFormItem>
          <ElFormItem label="活动状态">
            {{ detailRow.status === 1 ? '启用' : '停用' }}
          </ElFormItem>
        </ElForm>
      </template>
    </DetailDrawer>

    <EditDrawer>
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="仅维护拼团配置；商品、商户与已产生团单的成员、价格快照不可在此修改。"
      />
      <ElForm label-width="118px">
        <ElFormItem label="拼团价" required>
          <ElInputNumber
            v-model="form.price"
            :min="0.01"
            :precision="2"
            :step="1"
          />
        </ElFormItem>
        <ElFormItem label="成团人数" required>
          <ElInputNumber v-model="form.buying_count_num" :min="2" :max="9999" />
        </ElFormItem>
        <ElFormItem label="成团时限（小时）" required>
          <ElInputNumber v-model="form.time" :min="1" :max="720" />
        </ElFormItem>
        <ElFormItem label="活动时间" required>
          <ElInput
            v-model="form.start_time"
            class="!w-48"
            placeholder="YYYY-MM-DD HH:mm:ss"
          />
          <span class="mx-2">至</span>
          <ElInput
            v-model="form.end_time"
            class="!w-48"
            placeholder="YYYY-MM-DD HH:mm:ss"
          />
        </ElFormItem>
        <ElFormItem label="前台展示">
          <ElRadioGroup v-model="form.is_show">
            <ElRadio :label="1">上架</ElRadio>
            <ElRadio :label="0">下架</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="活动状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :label="1">启用</ElRadio>
            <ElRadio :label="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </EditDrawer>

    <AuditDrawer>
      <ElForm label-width="90px">
        <ElFormItem label="审核结果" required>
          <ElRadioGroup v-model="auditForm.status">
            <ElRadio :label="1">通过</ElRadio>
            <ElRadio :label="-1">拒绝</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem
          v-if="auditForm.status === -1"
          label="拒绝原因"
          required
        >
          <ElInput
            v-model="auditForm.refusal"
            type="textarea"
            :rows="3"
            maxlength="200"
            show-word-limit
            placeholder="请填写拒绝原因"
          />
        </ElFormItem>
      </ElForm>
    </AuditDrawer>

    <ProductPreviewModal
      ref="previewModalRef"
      modal-title="拼团预览"
      :product-id="previewProductId"
      :product-title="previewProductTitle"
      :display-price="previewDisplayPrice"
      :display-ot-price="previewDisplayOtPrice"
    />
    <ProductLabelSelectModal ref="labelModalRef" @submit="onLabelSubmit" />

    <ForceOffModal>
      <ElForm label-width="90px">
        <ElFormItem label="下架原因" required>
          <ElInput
            v-model="forceOffReason"
            type="textarea"
            :rows="3"
            maxlength="200"
            show-word-limit
            placeholder="请填写强制下架原因"
          />
        </ElFormItem>
      </ElForm>
    </ForceOffModal>
  </Page>
</template>
