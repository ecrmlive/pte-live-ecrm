<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, h, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { ArrowDown } from '@element-plus/icons-vue';
import {
  ElAvatar,
  ElButton,
  ElDatePicker,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElIcon,
  ElInput,
  ElMessage,
  ElOption,
  ElPagination,
  ElSelect,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  clearDistributionParentApi,
  getDistributionSummaryApi,
  listDistributionChildrenApi,
  listDistributionCommissionsApi,
  listDistributionLevelsApi,
  listDistributionPromotersApi,
  listDistributionSpreadOrdersApi,
  updateDistributionLevelApi,
  type CommissionStatus,
  type DistributionChild,
  type DistributionCommission,
  type DistributionLevel,
  type DistributionPromoter,
  type DistributionSpreadOrder,
  type DistributionSummary,
} from '#/api/core/platform-spread';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const emptySummary = (): DistributionSummary => ({
  active_promoter_count: 0,
  available_commission: 0,
  pending_commission: 0,
  promoter_count: 0,
  settled_commission: 0,
  spread_order_amount: 0,
  spread_order_count: 0,
  spread_user_count: 0,
  unwithdrawn_amount: 0,
  withdrawn_amount: 0,
});

const summary = ref<DistributionSummary>(emptySummary());
const lastFormValues = ref<Record<string, unknown>>({});
const levels = ref<DistributionLevel[]>([]);
const levelOptions = computed(() =>
  levels.value.map((item) => ({ label: item.name, value: item.id })),
);

const currentUser = ref<DistributionPromoter | null>(null);
const levelForm = reactive({ level_id: 0 as number, reason: '后台调整分销员等级' });
const levelSaving = ref(false);

const childrenRows = ref<DistributionChild[]>([]);
const childrenTotal = ref(0);
const childrenPage = ref(1);
const childrenPageSize = ref(10);
const childrenLoading = ref(false);
const childrenFilters = reactive({
  date_range: [] as string[],
  level: 0 as 0 | 1 | 2,
  keyword: '',
});
const childrenSort = reactive<{
  field?: 'spread_user_count';
  order?: 'asc' | 'desc';
}>({});

const commissionRows = ref<DistributionCommission[]>([]);
const commissionTotal = ref(0);
const commissionPage = ref(1);
const commissionLoading = ref(false);

const orderRows = ref<DistributionSpreadOrder[]>([]);
const orderTotal = ref(0);
const orderPage = ref(1);
const orderLoading = ref(false);

const commissionLabels: Record<CommissionStatus, string> = {
  available: '可结算',
  pending: '待结算',
  settled: '已结算',
  voided: '已作废',
};

const summaryCards = computed(() => [
  {
    key: 'promoter_count',
    label: '分销员人数(人)',
    value: String(summary.value.promoter_count || 0),
    icon: 'ant-design:team-outlined',
    tone: 'blue',
  },
  {
    key: 'spread_user_count',
    label: '推广人数(人)',
    value: String(summary.value.spread_user_count || 0),
    icon: 'ant-design:user-add-outlined',
    tone: 'orange',
  },
  {
    key: 'spread_order_count',
    label: '推广订单数',
    value: String(summary.value.spread_order_count || 0),
    icon: 'ant-design:shopping-outlined',
    tone: 'green',
  },
  {
    key: 'spread_order_amount',
    label: '推广订单金额',
    value: Number(summary.value.spread_order_amount || 0).toFixed(2),
    icon: 'ant-design:pay-circle-outlined',
    tone: 'pink',
  },
  {
    key: 'withdrawn_amount',
    label: '已提现金额(元)',
    value: Number(summary.value.withdrawn_amount || 0).toFixed(2),
    icon: 'ant-design:wallet-outlined',
    tone: 'purple',
  },
  {
    key: 'unwithdrawn_amount',
    label: '未提现金额(元)',
    value: Number(summary.value.unwithdrawn_amount || 0).toFixed(2),
    icon: 'mdi:cash',
    tone: 'sky',
  },
]);

function money(v?: number) {
  return Number(v || 0).toFixed(2);
}

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '') return '—';
  return String(v);
}

function maskMobile(mobile?: string) {
  const raw = String(mobile || '').trim();
  if (!raw) return '—';
  if (raw.length < 7) return raw;
  return `${raw.slice(0, 3)}****${raw.slice(-4)}`;
}

function avatarSrc(url?: string) {
  const resolved = resolveCosMediaUrl(url || '');
  return resolved || undefined;
}

const KEYWORD_OPTIONS = [
  { label: '昵称', value: 'nickname' },
  { label: '用户ID', value: 'uid' },
  { label: '手机号', value: 'phone' },
];

function renderSearchPrepend(
  field: 'keyword_type',
  options: Array<{ label: string; value: string }>,
  defaultValue: string,
) {
  return (values: Record<string, any>, api: { setFieldValue: Function }) => ({
    prepend: () =>
      h(
        ElSelect,
        {
          modelValue: values[field] || defaultValue,
          style: { width: '110px' },
          'onUpdate:modelValue': (v: string) => api.setFieldValue(field, v),
        },
        () =>
          options.map((opt) =>
            h(ElOption, { label: opt.label, value: opt.value, key: opt.value }),
          ),
      ),
  });
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const keyword = String(formValues?.keyword ?? '').trim();
  const keywordTypeRaw = String(formValues?.keyword_type ?? 'nickname').trim();
  const keywordType =
    keywordTypeRaw === 'uid' || keywordTypeRaw === 'phone'
      ? keywordTypeRaw
      : 'nickname';
  const levelRaw = formValues?.level_id;
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    level_id:
      levelRaw === 0 || levelRaw
        ? Number(levelRaw) || undefined
        : undefined,
    keyword: keyword || undefined,
    keyword_type: keyword ? keywordType : undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    { ...LIST_DATE_RANGE_FIELD, label: '时间选择' },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [] as Array<{ label: string; value: number }>,
        placeholder: '请选择',
      },
      fieldName: 'level_id',
      label: '等级名称',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '请输入内容' },
      defaultValue: '',
      fieldName: 'keyword',
      label: '分销员搜索',
      renderComponentContent: renderSearchPrepend(
        'keyword_type',
        KEYWORD_OPTIONS,
        'nickname',
      ),
    },
    {
      component: 'Input',
      defaultValue: 'nickname',
      dependencies: { show: () => false, triggerFields: [''] },
      fieldName: 'keyword_type',
      label: '搜索类型',
    },
  ],
  {
    commonConfig: { componentProps: { class: 'w-full' } },
    submitButtonOptions: { content: '搜索' },
    handleSubmit: async (values) => {
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
    handleReset: async () => {
      await formApi.resetForm();
      const values = (await formApi.getValues()) ?? {};
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
  },
);

const [Form, formApi] = useVbenForm(formOptions);

const gridOptions: VxeGridProps<DistributionPromoter> = {
  columns: [
    { field: 'user_id', title: 'ID', width: 90 },
    {
      field: 'avatar_url',
      title: '头像',
      width: 72,
      align: 'center',
      slots: { default: 'avatar' },
    },
    {
      field: 'nickname',
      title: '用户信息',
      minWidth: 160,
      showOverflow: false,
      slots: { default: 'userInfo' },
    },
    {
      field: 'spread_user_count',
      title: '推广用户数量',
      width: 120,
      slots: { default: 'spreadCount' },
    },
    {
      field: 'level_name',
      title: '等级名称',
      minWidth: 110,
      formatter: ({ cellValue }) => dash(cellValue),
    },
    {
      field: 'spread_order_count',
      title: '推广订单数量',
      width: 120,
    },
    {
      field: 'spread_order_amount',
      title: '推广订单金额',
      width: 120,
      formatter: ({ cellValue }) => money(Number(cellValue)),
    },
    {
      field: 'commission_amount',
      title: '佣金金额',
      width: 110,
      sortable: true,
      formatter: ({ cellValue }) => money(Number(cellValue)),
    },
    {
      field: 'withdrawn_amount',
      title: '已提现金额',
      width: 120,
      sortable: true,
      formatter: ({ cellValue }) => money(Number(cellValue)),
    },
    {
      field: 'withdraw_count',
      title: '提现次数',
      width: 96,
    },
    {
      field: 'unwithdrawn_amount',
      title: '未提现金额',
      width: 120,
      sortable: true,
      formatter: ({ cellValue }) => money(Number(cellValue)),
    },
    {
      field: 'parent_nickname',
      title: '推广人',
      minWidth: 110,
      formatter: ({ row }) =>
        row.parent_nickname || (row.parent_user_id ? `#${row.parent_user_id}` : '—'),
    },
    platformListActionColumn({ width: 160 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, formValues) => {
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const filters = buildFilterParams(values);
        const sort = Array.isArray(sorts) ? sorts[0] : undefined;
        const sortField =
          sort?.field === 'commission_amount' ||
          sort?.field === 'withdrawn_amount' ||
          sort?.field === 'unwithdrawn_amount'
            ? sort.field
            : undefined;
        const [data, stats] = await Promise.all([
          listDistributionPromotersApi({
            page: page.currentPage,
            limit: page.pageSize,
            ...filters,
            sort_field: sortField,
            sort_order: sort?.order === 'asc' ? 'asc' : sort?.order === 'desc' ? 'desc' : undefined,
          }),
          getDistributionSummaryApi(filters).catch(() => emptySummary()),
        ]);
        summary.value = stats || emptySummary();
        return { items: data.list || [], total: data.total || 0 };
      },
    },
    sort: true,
  },
  rowConfig: { isHover: true, keyField: 'user_id' },
  sortConfig: { remote: true, trigger: 'cell' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
    enabled: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [LevelDrawer, levelDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '编辑分销员等级',
  onConfirm: async () => {
    if (!currentUser.value) return;
    if (!levelForm.level_id) {
      ElMessage.warning('请选择分销员等级');
      return;
    }
    levelSaving.value = true;
    try {
      await updateDistributionLevelApi(currentUser.value.user_id, {
        level_id: Number(levelForm.level_id),
        reason: levelForm.reason.trim() || '后台调整分销员等级',
        idempotency_key: `promoter-level-${currentUser.value.user_id}-${crypto.randomUUID()}`,
      });
      ElMessage.success('分销员等级已更新');
      levelDrawerApi.close();
      await gridApi.reload(lastFormValues.value);
    } finally {
      levelSaving.value = false;
    }
  },
});

const [ChildrenDrawer, childrenDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  title: '推广人',
});

const [CommissionDrawer, commissionDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  title: '佣金记录',
});

const [OrderDrawer, orderDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  title: '推广订单',
});

async function loadChildren(page = childrenPage.value) {
  if (!currentUser.value) return;
  childrenLoading.value = true;
  childrenPage.value = page;
  try {
    const range = Array.isArray(childrenFilters.date_range)
      ? childrenFilters.date_range
      : [];
    const data = await listDistributionChildrenApi(currentUser.value.user_id, {
      page,
      limit: childrenPageSize.value,
      level: childrenFilters.level,
      keyword: childrenFilters.keyword.trim() || undefined,
      date_from: range[0] || undefined,
      date_to: range[1] || undefined,
      sort_field: childrenSort.field,
      sort_order: childrenSort.order,
    });
    childrenRows.value = data.list || [];
    childrenTotal.value = data.total || 0;
  } finally {
    childrenLoading.value = false;
  }
}

function searchChildren() {
  void loadChildren(1);
}

function resetChildrenFilters() {
  childrenFilters.date_range = [];
  childrenFilters.level = 0;
  childrenFilters.keyword = '';
  childrenSort.field = undefined;
  childrenSort.order = undefined;
  void loadChildren(1);
}

function onChildrenSortChange(payload: {
  prop?: string;
  order?: 'ascending' | 'descending' | null;
}) {
  if (payload.prop === 'spread_user_count' && payload.order) {
    childrenSort.field = 'spread_user_count';
    childrenSort.order = payload.order === 'ascending' ? 'asc' : 'desc';
  } else {
    childrenSort.field = undefined;
    childrenSort.order = undefined;
  }
  void loadChildren(1);
}

async function loadCommissions(page = 1) {
  if (!currentUser.value) return;
  commissionLoading.value = true;
  commissionPage.value = page;
  try {
    const data = await listDistributionCommissionsApi({
      page,
      limit: 10,
      user_id: currentUser.value.user_id,
    });
    commissionRows.value = data.list || [];
    commissionTotal.value = data.total || 0;
  } finally {
    commissionLoading.value = false;
  }
}

async function loadOrders(page = 1) {
  if (!currentUser.value) return;
  orderLoading.value = true;
  orderPage.value = page;
  try {
    const data = await listDistributionSpreadOrdersApi(currentUser.value.user_id, {
      page,
      limit: 10,
    });
    orderRows.value = data.list || [];
    orderTotal.value = data.total || 0;
  } finally {
    orderLoading.value = false;
  }
}

function openChildren(row: DistributionPromoter) {
  currentUser.value = row;
  childrenFilters.date_range = [];
  childrenFilters.level = 0;
  childrenFilters.keyword = '';
  childrenSort.field = undefined;
  childrenSort.order = undefined;
  childrenPage.value = 1;
  childrenDrawerApi.setState({ title: '推广人' });
  childrenDrawerApi.open();
  void loadChildren(1);
}

function openCommissions(row: DistributionPromoter) {
  currentUser.value = row;
  commissionDrawerApi.setState({ title: `佣金记录（#${row.user_id}）` });
  commissionDrawerApi.open();
  void loadCommissions(1);
}

function openOrders(row: DistributionPromoter) {
  currentUser.value = row;
  orderDrawerApi.setState({ title: `推广订单（#${row.user_id}）` });
  orderDrawerApi.open();
  void loadOrders(1);
}

function openEditLevel(row: DistributionPromoter) {
  currentUser.value = row;
  levelForm.level_id = Number(row.level_id || 0);
  levelForm.reason = '后台调整分销员等级';
  levelDrawerApi.open();
}

async function clearParent(row: DistributionPromoter) {
  if (!row.parent_user_id) {
    ElMessage.info('该分销员暂无上级推广人');
    return;
  }
  try {
    await confirm({
      title: '提示',
      content: `确认清除用户 #${row.user_id} 的上级推广人？此操作写入审计记录。`,
      icon: 'warning',
    });
  } catch {
    return;
  }
  await clearDistributionParentApi(row.user_id, {
    reason: '后台清除上级推广人',
    idempotency_key: `promoter-clear-parent-${row.user_id}-${crypto.randomUUID()}`,
  });
  ElMessage.success('已清除上级推广人');
  await gridApi.reload(lastFormValues.value);
}

function onMoreCommand(cmd: string, row: DistributionPromoter) {
  switch (cmd) {
    case 'commissions':
      openCommissions(row);
      break;
    case 'orders':
      openOrders(row);
      break;
    case 'clearParent':
      void clearParent(row);
      break;
    case 'editLevel':
      openEditLevel(row);
      break;
  }
}

onMounted(async () => {
  try {
    const data = await listDistributionLevelsApi();
    levels.value = data.list || [];
    formApi.updateSchema([
      {
        fieldName: 'level_id',
        componentProps: {
          clearable: true,
          options: levelOptions.value,
          placeholder: '请选择',
        },
      },
    ]);
  } catch {
    levels.value = [];
  }
});
</script>

<template>
  <Page auto-content-height>
    <div class="promoter-filter">
      <Form />
    </div>

    <div class="promoter-summary">
      <div class="verify-summary">
        <div
          v-for="card in summaryCards"
          :key="card.key"
          class="verify-summary__card"
          :class="`verify-summary__card--${card.tone}`"
        >
          <div class="verify-summary__icon">
            <IconifyIcon :icon="card.icon" />
          </div>
          <div class="verify-summary__body">
            <div class="verify-summary__value">{{ card.value }}</div>
            <div class="verify-summary__label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <Grid>
      <template #avatar="{ row }">
        <ElAvatar :size="36" :src="avatarSrc(row.avatar_url)">
          {{ (row.nickname || '?').slice(0, 1) }}
        </ElAvatar>
      </template>

      <template #userInfo="{ row }">
        <div class="user-info">
          <div>昵称: {{ dash(row.nickname) }}</div>
          <div>电话: {{ maskMobile(row.mobile) }}</div>
        </div>
      </template>

      <template #spreadCount="{ row }">
        <ElButton link type="primary" @click="openChildren(row)">
          {{ row.spread_user_count ?? row.direct_user_count ?? 0 }}
        </ElButton>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openChildren(row)">推广人</ElButton>
        <ElDropdown
          trigger="click"
          @command="(cmd: string) => onMoreCommand(cmd, row)"
        >
          <ElButton link type="primary">
            更多
            <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem command="commissions">佣金记录</ElDropdownItem>
              <ElDropdownItem command="orders">推广订单</ElDropdownItem>
              <ElDropdownItem command="clearParent">清除上级推广人</ElDropdownItem>
              <ElDropdownItem command="editLevel">编辑分销员等级</ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <LevelDrawer>
      <ElForm label-width="120px" class="px-2">
        <ElFormItem label="分销员">
          #{{ currentUser?.user_id }} {{ currentUser?.nickname || '' }}
        </ElFormItem>
        <ElFormItem label="等级名称" required>
          <ElSelect
            v-model="levelForm.level_id"
            class="w-full"
            clearable
            placeholder="请选择"
            :disabled="levelSaving"
          >
            <ElOption
              v-for="item in levels"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </ElSelect>
        </ElFormItem>
      </ElForm>
    </LevelDrawer>

    <ChildrenDrawer>
      <div class="children-panel">
        <ElForm class="children-filters" label-width="80px" @submit.prevent="searchChildren">
          <div class="children-filters__row">
            <ElFormItem label="时间选择">
              <ElDatePicker
                v-model="childrenFilters.date_range"
                type="daterange"
                value-format="YYYY-MM-DD"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                class="!w-[280px]"
                clearable
              />
            </ElFormItem>
            <ElFormItem label="用户类型">
              <ElSelect
                v-model="childrenFilters.level"
                class="!w-[180px]"
                placeholder="全部"
              >
                <ElOption label="全部" :value="0" />
                <ElOption label="一级推广人" :value="1" />
                <ElOption label="二级推广人" :value="2" />
              </ElSelect>
            </ElFormItem>
          </div>
          <div class="children-filters__row">
            <ElFormItem label="关键字">
              <ElInput
                v-model="childrenFilters.keyword"
                class="!w-[320px]"
                clearable
                placeholder="请输入姓名、电话、用户ID"
                @keyup.enter="searchChildren"
              />
            </ElFormItem>
            <ElFormItem label-width="0">
              <ElButton type="primary" @click="searchChildren">查询</ElButton>
              <ElButton @click="resetChildrenFilters">重置</ElButton>
            </ElFormItem>
          </div>
        </ElForm>

        <ElTable
          v-loading="childrenLoading"
          :data="childrenRows"
          row-key="user_id"
          empty-text="暂无数据"
          @sort-change="onChildrenSortChange"
        >
          <ElTableColumn prop="user_id" label="ID" width="90" />
          <ElTableColumn label="头像" width="72" align="center">
            <template #default="{ row }">
              <ElAvatar :size="36" :src="avatarSrc(row.avatar_url)">
                {{ (row.nickname || '?').slice(0, 1) }}
              </ElAvatar>
            </template>
          </ElTableColumn>
          <ElTableColumn label="用户信息" min-width="160">
            <template #default="{ row }">
              <div class="user-info">
                <div>昵称: {{ dash(row.nickname) }}</div>
                <div>电话: {{ maskMobile(row.mobile) }}</div>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn label="是否推广员" width="110" align="center">
            <template #default="{ row }">
              <ElTag :type="row.is_promoter === 1 ? 'success' : 'info'" size="small">
                {{ row.is_promoter === 1 ? '是' : '否' }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="spread_user_count"
            label="推广人数"
            width="110"
            sortable="custom"
          />
          <ElTableColumn prop="pay_count" label="订单数" width="96" />
          <ElTableColumn label="订单金额" width="110">
            <template #default="{ row }">
              {{ money(row.pay_amount) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="关注时间" min-width="170">
            <template #default="{ row }">
              {{ formatShanghaiDateTime(row.bound_at) }}
            </template>
          </ElTableColumn>
          <template #empty>
            <ElEmpty description="暂无数据" :image-size="72" />
          </template>
        </ElTable>

        <div class="children-pager">
          <ElPagination
            v-model:current-page="childrenPage"
            v-model:page-size="childrenPageSize"
            background
            layout="total, prev, pager, next, jumper"
            :total="childrenTotal"
            :page-sizes="[10, 20, 50]"
            @current-change="loadChildren"
            @size-change="() => loadChildren(1)"
          />
        </div>
      </div>
    </ChildrenDrawer>

    <CommissionDrawer>
      <div v-loading="commissionLoading" class="drawer-table">
        <table>
          <thead>
            <tr>
              <th>流水ID</th>
              <th>订单ID</th>
              <th>佣金金额</th>
              <th>状态</th>
              <th>可结算时间</th>
              <th>创建时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in commissionRows" :key="row.commission_id">
              <td>{{ row.commission_id }}</td>
              <td>{{ row.order_id || '—' }}</td>
              <td>{{ money(row.amount) }}</td>
              <td>{{ commissionLabels[row.status] || row.status }}</td>
              <td>{{ formatShanghaiDateTime(row.available_at) }}</td>
              <td>{{ formatShanghaiDateTime(row.created_at) }}</td>
            </tr>
            <tr v-if="!commissionRows.length">
              <td colspan="6" class="empty">暂无佣金记录</td>
            </tr>
          </tbody>
        </table>
        <div v-if="commissionTotal > 10" class="pager">
          <ElButton
            :disabled="commissionPage <= 1"
            @click="loadCommissions(commissionPage - 1)"
          >
            上一页
          </ElButton>
          <span>
            {{ commissionPage }} /
            {{ Math.max(1, Math.ceil(commissionTotal / 10)) }}
          </span>
          <ElButton
            :disabled="commissionPage * 10 >= commissionTotal"
            @click="loadCommissions(commissionPage + 1)"
          >
            下一页
          </ElButton>
        </div>
      </div>
    </CommissionDrawer>

    <OrderDrawer>
      <div v-loading="orderLoading" class="drawer-table">
        <table>
          <thead>
            <tr>
              <th>订单号</th>
              <th>买家</th>
              <th>实付金额</th>
              <th>佣金</th>
              <th>佣金状态</th>
              <th>支付时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in orderRows" :key="`${row.order_id}-${row.created_at}`">
              <td>{{ dash(row.order_no) }}</td>
              <td>{{ dash(row.buyer_name) }} (#{{ row.buyer_id || 0 }})</td>
              <td>{{ money(row.pay_amount) }}</td>
              <td>{{ money(row.commission) }}</td>
              <td>
                {{
                  commissionLabels[row.status as CommissionStatus] || row.status || '—'
                }}
              </td>
              <td>{{ formatShanghaiDateTime(row.paid_at || row.created_at) }}</td>
            </tr>
            <tr v-if="!orderRows.length">
              <td colspan="6" class="empty">暂无推广订单</td>
            </tr>
          </tbody>
        </table>
        <div v-if="orderTotal > 10" class="pager">
          <ElButton :disabled="orderPage <= 1" @click="loadOrders(orderPage - 1)">
            上一页
          </ElButton>
          <span>{{ orderPage }} / {{ Math.max(1, Math.ceil(orderTotal / 10)) }}</span>
          <ElButton
            :disabled="orderPage * 10 >= orderTotal"
            @click="loadOrders(orderPage + 1)"
          >
            下一页
          </ElButton>
        </div>
      </div>
    </OrderDrawer>
  </Page>
</template>

<style scoped>
.promoter-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.promoter-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.verify-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.verify-summary__card {
  display: flex;
  gap: 16px;
  align-items: center;
  min-height: 88px;
  padding: 20px 22px;
  background: hsl(var(--background));
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 3%);
}

.verify-summary__icon {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  color: #fff;
  font-size: 24px;
  border-radius: 50%;
}

.verify-summary__card--blue .verify-summary__icon {
  background: #409eff;
}

.verify-summary__card--orange .verify-summary__icon {
  background: #e6a23c;
}

.verify-summary__card--green .verify-summary__icon {
  background: #67c23a;
}

.verify-summary__card--pink .verify-summary__icon {
  background: #f56c6c;
}

.verify-summary__card--purple .verify-summary__icon {
  background: #9b59b6;
}

.verify-summary__card--sky .verify-summary__icon {
  background: #409eff;
}

.verify-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.verify-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

.user-info {
  line-height: 1.5;
  font-size: 13px;
}

.children-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 8px 16px;
}

.children-filters__row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  align-items: flex-start;
}

.children-pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 4px;
}

.drawer-table {
  padding: 4px 8px 16px;
}

.drawer-table table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.drawer-table th,
.drawer-table td {
  padding: 10px 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  text-align: left;
}

.drawer-table th {
  color: var(--el-text-color-secondary);
  font-weight: 500;
  background: var(--el-fill-color-lighter);
}

.drawer-table .empty {
  padding: 28px 0;
  color: var(--el-text-color-secondary);
  text-align: center;
}

.pager {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: flex-end;
  margin-top: 16px;
}

@media (min-width: 1600px) {
  .verify-summary {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
}

@media (max-width: 1100px) {
  .verify-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .verify-summary {
    grid-template-columns: 1fr;
  }
}
</style>
