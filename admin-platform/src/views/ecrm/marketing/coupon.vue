<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRadio,
  ElRadioGroup,
  ElSwitch,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';
import { CloseBold, Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  clonePlatformCouponApi,
  createPlatformCouponApi,
  deletePlatformCouponApi,
  getPlatformCouponDetailApi,
  listPlatformCouponsApi,
  setPlatformCouponStatusApi,
  updatePlatformCouponApi,
  type PlatformCoupon,
  type PlatformCouponDetail,
  type PlatformCouponSaveInput,
} from '#/api/core/platform-promotion';
import CouponUserIssueModal, {
  type CouponIssueMode,
} from '#/components/marketing/coupon-user-issue-modal.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const SEND_RECEIVE = 0;
const SEND_GIFT = 3;
const SEND_ADMIN = 6;

const saving = ref(false);
const issueOpen = ref(false);
const issueCouponId = ref(0);
const issueMode = ref<CouponIssueMode>('receive');
const editingID = ref<number>();
const canManage = ref(false);
const detail = ref<PlatformCouponDetail | null>(null);
const detailLoading = ref(false);
const detailTab = ref('basic');
const useRange = ref<[string, string] | null>(null);
const claimRange = ref<[string, string] | null>(null);

const form = reactive<PlatformCouponSaveInput>({
  coupon_price: 0,
  coupon_time: 30,
  coupon_type: 0,
  full_reduction: 0,
  is_limited: 0,
  is_timeout: 0,
  send_type: SEND_RECEIVE,
  sort: 0,
  status: 1,
  title: '',
  total_count: 0,
  use_min_price: 0,
  use_type: 0,
});

const sendTypeLabel = computed(() => {
  const t = Number(detail.value?.send_type ?? 0);
  if (t === SEND_ADMIN) return '后台发放';
  if (t === SEND_GIFT) return '赠送券';
  return '领取';
});

const typeLabel = computed(() =>
  (detail.value?.type ?? 10) === 10 ? '平台通用券' : '店铺券',
);

function formatTime(value?: string | null) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function formatMinPrice(value?: number) {
  const n = Number(value || 0);
  return n > 0 ? `满 ¥${n.toFixed(2)} 可用` : '无门槛';
}

function formatPublishCount(row: Pick<PlatformCoupon, 'is_limited' | 'total_count'>) {
  return row.is_limited === 1 ? String(row.total_count) : '不限量';
}

function formatRemainCount(row: Pick<PlatformCoupon, 'is_limited' | 'remain_count'>) {
  return row.is_limited === 1 ? String(row.remain_count) : '不限量';
}

function sendTypeText(v?: number) {
  if (v === SEND_ADMIN) return '后台发放';
  if (v === SEND_GIFT) return '赠送券';
  return '领取';
}

function emptyForm(): PlatformCouponSaveInput {
  return {
    coupon_price: 0,
    coupon_time: 30,
    coupon_type: 0,
    full_reduction: 0,
    is_limited: 0,
    is_timeout: 0,
    send_type: SEND_RECEIVE,
    sort: 0,
    status: 1,
    title: '',
    total_count: 0,
    use_min_price: 0,
    use_type: 0,
  };
}

function fillFormFromCoupon(row: PlatformCoupon | PlatformCouponDetail) {
  Object.assign(form, emptyForm(), {
    coupon_price: Number(row.coupon_price || 0),
    coupon_time: Number(row.coupon_time || 30),
    coupon_type: Number(row.coupon_type ?? 0),
    full_reduction: Number(row.full_reduction || 0),
    is_limited: Number(row.is_limited || 0),
    is_timeout: Number(row.is_timeout || 0),
    send_type: Number(row.send_type ?? SEND_RECEIVE),
    sort: Number(row.sort || 0),
    status: Number(row.status ?? 1),
    title: row.title || '',
    total_count: Number(row.total_count || 0),
    use_min_price: Number(row.use_min_price || 0),
    use_type:
      'use_type' in row && typeof row.use_type === 'number'
        ? row.use_type
        : Number(row.use_min_price || 0) > 0
          ? 1
          : 0,
  });
  useRange.value =
    row.use_start_time && row.use_end_time
      ? [String(row.use_start_time), String(row.use_end_time)]
      : null;
  claimRange.value =
    row.start_time && row.end_time
      ? [String(row.start_time), String(row.end_time)]
      : null;
}

function buildSaveBody(): PlatformCouponSaveInput {
  const body: PlatformCouponSaveInput = {
    ...form,
    title: form.title.trim(),
    use_min_price: form.use_type === 1 ? Math.round(Number(form.use_min_price) || 0) : 0,
    total_count: form.is_limited === 1 ? form.total_count : 0,
    full_reduction: form.send_type === SEND_GIFT ? Number(form.full_reduction || 0) : 0,
    use_start_time: '',
    use_end_time: '',
    start_time: '',
    end_time: '',
  };
  if (form.coupon_type === 1 && useRange.value?.length === 2) {
    body.use_start_time = useRange.value[0];
    body.use_end_time = useRange.value[1];
    body.coupon_time = 0;
  }
  if (form.is_timeout === 1 && claimRange.value?.length === 2) {
    body.start_time = claimRange.value[0];
    body.end_time = claimRange.value[1];
  }
  return body;
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const statusRaw = formValues?.status;
  const sendRaw = formValues?.send_type;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status:
      statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
    send_type:
      sendRaw === 0 || sendRaw === SEND_ADMIN || sendRaw === SEND_GIFT
        ? Number(sendRaw)
        : undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '优惠券名称' },
    fieldName: 'keyword',
    label: '优惠券名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '领取', value: SEND_RECEIVE },
        { label: '后台发放', value: SEND_ADMIN },
        { label: '赠送券', value: SEND_GIFT },
      ],
      placeholder: '全部',
    },
    fieldName: 'send_type',
    label: '获取方式',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启', value: 1 },
        { label: '关闭', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<PlatformCoupon> = {
  columns: [
    { field: 'coupon_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 160, showOverflow: false, title: '优惠券名称' },
    {
      field: 'type',
      title: '优惠券类型',
      width: 110,
      formatter: () => '通用券',
    },
    {
      field: 'send_type',
      title: '获取方式',
      width: 100,
      formatter: ({ cellValue }) => sendTypeText(Number(cellValue ?? 0)),
    },
    {
      field: 'coupon_price',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      title: '面额',
      width: 90,
    },
    {
      field: 'use_min_price',
      formatter: ({ cellValue }) => formatMinPrice(Number(cellValue)),
      title: '使用门槛',
      width: 130,
    },
    {
      field: 'coupon_time',
      formatter: ({ row }) =>
        Number(row.coupon_type || 0) === 1
          ? '时间段'
          : `${row.coupon_time || 0}天`,
      title: '使用时间',
      width: 100,
    },
    {
      field: 'remain_count',
      formatter: ({ row }) =>
        row.is_limited === 1
          ? `发布: ${row.total_count} 剩余: ${row.remain_count}`
          : '不限量',
      title: '发布数量',
      width: 150,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    platformListActionColumn({ width: 320 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformCouponsApi(
          buildListParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'coupon_id' },
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
  confirmText: '提交',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, emptyForm());
  useRange.value = null;
  claimRange.value = null;
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加优惠券' }).open();
}

function openEdit(row: PlatformCoupon) {
  editingID.value = row.coupon_id;
  fillFormFromCoupon(row);
  formDrawerApi.setState({ title: '编辑优惠券' }).open();
}

async function openDetail(row: PlatformCoupon) {
  detailTab.value = 'basic';
  detail.value = null;
  detailDrawerApi.setState({ title: '优惠券详情', loading: true }).open();
  detailLoading.value = true;
  try {
    detail.value = await getPlatformCouponDetailApi(row.coupon_id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function openCopy(row: PlatformCoupon) {
  try {
    await ElMessageBox.confirm(
      `确认复制优惠券「${row.title}」？将创建一条新券。`,
      '复制确认',
      { type: 'info' },
    );
    await clonePlatformCouponApi(row.coupon_id);
    ElMessage.success('已复制优惠券');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

function openIssueRecords(mode: CouponIssueMode, couponId?: number) {
  const id = couponId ?? detail.value?.coupon_id;
  if (!id) return;
  issueCouponId.value = id;
  issueMode.value = mode;
  issueOpen.value = true;
}

async function save() {
  if (!form.title.trim()) {
    ElMessage.warning('请输入优惠券名称');
    return;
  }
  if (form.coupon_price <= 0) {
    ElMessage.warning('请输入优惠券面值');
    return;
  }
  if (form.use_type === 1 && form.use_min_price <= 0) {
    ElMessage.warning('请输入优惠券最低消费');
    return;
  }
  if (form.coupon_type === 0 && form.coupon_time <= 0) {
    ElMessage.warning('请输入有效天数');
    return;
  }
  if (form.coupon_type === 1 && (!useRange.value || useRange.value.length < 2)) {
    ElMessage.warning('请选择使用时间段');
    return;
  }
  if (
    form.is_timeout === 1 &&
    (!claimRange.value || claimRange.value.length < 2)
  ) {
    ElMessage.warning('请选择领取时间段');
    return;
  }
  if (form.is_limited === 1 && form.total_count <= 0) {
    ElMessage.warning('限量发放时必须填写发布数量');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const body = buildSaveBody();
    if (editingID.value) await updatePlatformCouponApi(editingID.value, body);
    else await createPlatformCouponApi(body);
    formDrawerApi.close();
    ElMessage.success(editingID.value ? '优惠券已更新' : '优惠券已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function toggle(row: PlatformCoupon) {
  const next = row.status === 1 ? 0 : 1;
  await setPlatformCouponStatusApi(row.coupon_id, next);
  row.status = next;
  ElMessage.success(next === 1 ? '已开启' : '已关闭');
}

async function remove(row: PlatformCoupon) {
  try {
    await ElMessageBox.confirm(
      `删除优惠券「${row.title}」后不可恢复，是否继续？`,
      '删除确认',
      { type: 'warning' },
    );
    await deletePlatformCouponApi(row.coupon_id);
    ElMessage.success('优惠券已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('marketing.coupon.manage');
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
          添加优惠券
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          v-if="canManage"
          :model-value="row.status === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="() => toggle(row)"
        />
        <ElTag v-else :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '开启' : '关闭' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          link
          type="primary"
          @click="openIssueRecords('receive', row.coupon_id)"
        >
          领取记录
        </ElButton>
        <ElButton
          link
          type="primary"
          @click="openIssueRecords('used', row.coupon_id)"
        >
          使用记录
        </ElButton>
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="primary" @click="openCopy(row)">复制</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
      </template>
    </Grid>

    <CouponUserIssueModal
      v-model:open="issueOpen"
      :coupon-id="issueCouponId"
      coupon-scope="platform"
      :mode="issueMode"
    />

    <FormDrawer>
      <ElForm label-width="120px" class="coupon-form">
        <ElFormItem label="优惠券名称" required>
          <ElInput
            v-model="form.title"
            maxlength="40"
            show-word-limit
            placeholder="请输入优惠券名称"
            class="max-w-md"
          />
        </ElFormItem>
        <ElFormItem label="优惠券面值" required>
          <ElInputNumber
            v-model="form.coupon_price"
            :min="0.01"
            :precision="2"
            :step="1"
          />
        </ElFormItem>
        <ElFormItem label="使用门槛">
          <ElRadioGroup v-model="form.use_type">
            <ElRadio :label="0">无门槛</ElRadio>
            <ElRadio :label="1">有门槛</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="form.use_type === 1" label="最低消费" required>
          <ElInputNumber
            v-model="form.use_min_price"
            :min="1"
            :precision="0"
            :step="1"
          />
        </ElFormItem>
        <ElFormItem label="使用有效期">
          <ElRadioGroup v-model="form.coupon_type">
            <ElRadio :label="0">天数</ElRadio>
            <ElRadio :label="1">时间段</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="form.coupon_type === 0" label="天数" required>
          <ElInputNumber v-model="form.coupon_time" :min="1" :precision="0" />
        </ElFormItem>
        <ElFormItem v-else label="时间段" required>
          <ElDatePicker
            v-model="useRange"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            class="!w-[420px]"
          />
        </ElFormItem>
        <ElFormItem label="领取时间">
          <ElRadioGroup v-model="form.is_timeout">
            <ElRadio :label="0">不限时</ElRadio>
            <ElRadio :label="1">限时</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="form.is_timeout === 1" label="领取时段" required>
          <ElDatePicker
            v-model="claimRange"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            class="!w-[420px]"
          />
        </ElFormItem>
        <ElFormItem label="获取方式">
          <ElRadioGroup v-model="form.send_type">
            <ElRadio :label="SEND_RECEIVE">领取</ElRadio>
            <ElRadio :label="SEND_ADMIN">后台发放</ElRadio>
            <ElRadio :label="SEND_GIFT">赠送券</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="form.send_type === SEND_GIFT" label="满赠金额">
          <ElInputNumber
            v-model="form.full_reduction"
            :min="0"
            :precision="2"
          />
        </ElFormItem>
        <ElFormItem label="是否限量">
          <ElRadioGroup v-model="form.is_limited">
            <ElRadio :label="0">不限量</ElRadio>
            <ElRadio :label="1">限量</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="form.is_limited === 1" label="发布数量" required>
          <ElInputNumber v-model="form.total_count" :min="1" :precision="0" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :max="99999" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :label="1">开启</ElRadio>
            <ElRadio :label="0">关闭</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <DetailDrawer>
      <div v-loading="detailLoading" class="coupon-detail">
        <ElTabs v-model="detailTab">
          <ElTabPane label="基本信息" name="basic" />
        </ElTabs>

        <template v-if="detail">
          <section class="coupon-detail__section">
            <div class="coupon-detail__section-title">优惠券信息</div>
            <div class="coupon-detail__grid">
              <div class="coupon-detail__item">
                <span class="label">优惠券名称</span>
                <span class="value">{{ detail.title || '—' }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">优惠券类型</span>
                <span class="value">通用券</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">优惠券面值</span>
                <span class="value">{{
                  Number(detail.coupon_price || 0).toFixed(2)
                }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">使用门槛</span>
                <span class="value">{{
                  formatMinPrice(detail.use_min_price)
                }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">使用有效期</span>
                <span class="value">
                  <template v-if="Number(detail.coupon_type || 0) === 1">
                    {{ formatTime(detail.use_start_time) }} ~
                    {{ formatTime(detail.use_end_time) }}
                  </template>
                  <template v-else>{{ detail.coupon_time || 0 }}天</template>
                </span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">领取时间</span>
                <span class="value">
                  <template v-if="Number(detail.is_timeout || 0) === 1">
                    {{ formatTime(detail.start_time) }} ~
                    {{ formatTime(detail.end_time) }}
                  </template>
                  <template v-else>不限时</template>
                </span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">获取方式</span>
                <span class="value">{{ sendTypeLabel }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">类型</span>
                <span class="value">{{ typeLabel }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">是否限量</span>
                <span class="value">
                  <ElTag
                    v-if="detail.is_limited === 1"
                    type="success"
                    size="small"
                  >
                    是
                  </ElTag>
                  <span v-else class="limit-no">
                    <CloseBold class="limit-no__icon" />
                  </span>
                </span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">已发布总数</span>
                <span class="value">{{ formatPublishCount(detail) }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">剩余总数</span>
                <span class="value">{{ formatRemainCount(detail) }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">创建时间</span>
                <span class="value">{{ formatTime(detail.create_time) }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">状态</span>
                <span class="value">{{
                  detail.status === 1 ? '开启' : '关闭'
                }}</span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">排序</span>
                <span class="value">{{ detail.sort ?? 0 }}</span>
              </div>
            </div>
          </section>

          <section class="coupon-detail__section">
            <div class="coupon-detail__section-title">优惠券情况</div>
            <div class="coupon-detail__grid coupon-detail__grid--stats">
              <div class="coupon-detail__item">
                <span class="label">已领取/发放总数</span>
                <span class="value value--inline">
                  {{ detail.received_total ?? 0 }}
                  <ElButton
                    link
                    type="primary"
                    @click="openIssueRecords('receive')"
                  >
                    领取记录
                  </ElButton>
                </span>
              </div>
              <div class="coupon-detail__item">
                <span class="label">已使用总数</span>
                <span class="value value--inline">
                  {{ detail.used_total ?? 0 }}
                  <ElButton
                    link
                    type="primary"
                    @click="openIssueRecords('used')"
                  >
                    使用记录
                  </ElButton>
                </span>
              </div>
            </div>
          </section>
        </template>
      </div>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.coupon-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.coupon-detail {
  min-height: 240px;
  padding: 0 4px 16px;
}

.coupon-detail :deep(.el-tabs__header) {
  margin-bottom: 16px;
}

.coupon-detail__section + .coupon-detail__section {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px dashed hsl(var(--border));
}

.coupon-detail__section-title {
  position: relative;
  margin-bottom: 16px;
  padding-left: 10px;
  color: hsl(var(--foreground));
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
}

.coupon-detail__section-title::before {
  position: absolute;
  top: 2px;
  left: 0;
  width: 3px;
  height: 14px;
  background: hsl(var(--primary));
  border-radius: 2px;
  content: '';
}

.coupon-detail__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px 32px;
}

.coupon-detail__grid--stats {
  grid-template-columns: 1fr;
}

.coupon-detail__item {
  display: grid;
  grid-template-columns: 120px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
  font-size: 13px;
  line-height: 1.6;
}

.coupon-detail__item .label {
  color: hsl(var(--muted-foreground));
}

.coupon-detail__item .value {
  color: hsl(var(--foreground));
  word-break: break-all;
}

.coupon-detail__item .value--inline {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.limit-no {
  display: inline-flex;
  color: hsl(var(--primary));
}

.limit-no__icon {
  width: 14px;
  height: 14px;
}
</style>
