<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElMessage,
  ElMessageBox,
  ElSkeleton,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  exportPlatformUsers,
  fetchPlatformUserDetail,
  fetchPlatformUsers,
  type PlatformUserDetail,
  type PlatformUserRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const detail = ref<PlatformUserDetail>();
const detailLoading = ref(false);
const canRead = ref(false);
const canExport = ref(false);
const isPlatform = ref(false);
const lastQueryParams = ref<Record<string, unknown>>({});

const assetLabels: Record<string, string> = {
  balance: '余额',
  points: '积分',
  commission: '佣金',
};
const orderLabels: Record<string, string> = {
  pending_pay: '待支付',
  paid: '已支付',
  awaiting_final: '待尾款',
  final_timeout: '尾款超时',
  fulfilling: '履约中',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消',
  aftersale: '售后中',
};
const membershipLabels: Record<string, string> = {
  initial: '初始',
  upgrade: '升级',
  downgrade: '降级',
  manual: '人工调整',
};
const couponStatusLabels: Record<string, string> = {
  unused: '未使用',
  locked: '锁定',
  used: '已使用',
  expired: '已过期',
};

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'id',
    label: '用户 ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '昵称关键词' },
    fieldName: 'nickname',
    label: '昵称',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '手机号（完整或后四位）' },
    fieldName: 'phone',
    label: '手机号',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '昵称 / 手机号 / 用户 ID 综合搜索',
    },
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

const gridOptions: VxeGridProps<PlatformUserRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    { field: 'nickname', minWidth: 160, showOverflow: false, title: '昵称' },
    { field: 'mobile', title: '手机号（脱敏）', width: 150 },
    {
      field: 'balance',
      title: '余额',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'points', title: '积分', width: 100 },
    {
      field: 'level_name',
      minWidth: 130,
      showOverflow: false,
      title: '会员等级',
      formatter: ({ cellValue }) => cellValue || '普通会员',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    {
      field: 'created_at',
      title: '注册时间',
      width: 180,
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 100 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!isPlatform.value) {
          return { items: [], total: 0 };
        }
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const params = {
          page: page.currentPage,
          limit: page.pageSize,
          id: formValues?.id ? Number(formValues.id) : undefined,
          keyword:
            String(formValues?.keyword ?? '').trim() ||
            String(formValues?.nickname ?? '').trim() ||
            undefined,
          nickname: String(formValues?.nickname ?? '').trim() || undefined,
          phone: String(formValues?.phone ?? '').trim() || undefined,
          status:
            formValues?.status === 0 || formValues?.status === 1
              ? (Number(formValues.status) as 0 | 1)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        };
        lastQueryParams.value = params;
        const result = await fetchPlatformUsers(params);
        return { items: result.list || [], total: result.total || 0 };
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

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[960px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function openDetail(row: PlatformUserRow) {
  detail.value = undefined;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: '用户监管详情', loading: true }).open();
  try {
    detail.value = await fetchPlatformUserDetail(row.id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function exportRows() {
  try {
    const { value } = await ElMessageBox.prompt(
      '请填写导出原因，导出仅包含脱敏最小字段，最多 5000 行。',
      '导出用户信息',
      {
        inputPattern: /.{2,}/,
        inputErrorMessage: '导出原因至少 2 个字符',
        confirmButtonText: '生成 CSV',
        cancelButtonText: '取消',
      },
    );
    const params = lastQueryParams.value;
    const result = await exportPlatformUsers({
      id: params.id as number | undefined,
      keyword: (params.keyword as string | undefined) || undefined,
      nickname: (params.nickname as string | undefined) || undefined,
      phone: (params.phone as string | undefined) || undefined,
      status: params.status as 0 | 1 | undefined,
      date_from: params.date_from as string | undefined,
      date_to: params.date_to as string | undefined,
      reason: value.trim(),
    });
    const blob = new Blob([result.content], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条脱敏用户记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`,
    );
  } catch {
    // 用户取消
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  isPlatform.value = profile.roles.includes('platform');
  canRead.value = isPlatform.value && codes.includes('user.list.read');
  canExport.value = isPlatform.value && codes.includes('user.list.export');
  if (isPlatform.value) {
    gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="isPlatform && !canRead"
      class="mb-4"
      title="当前账号缺少 user.list.read 权限，列表仍尝试加载；导出与敏感操作已受限。"
      type="warning"
      :closable="false"
    />
    <Grid>
      <template #toolbar-actions>
        <ElButton v-if="canExport" type="success" plain @click="exportRows">
          导出脱敏信息
        </ElButton>
      </template>

      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">监管详情</ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="detail">
            <ElDescriptions :column="2" border>
              <ElDescriptionsItem label="用户 ID">
                {{ detail.profile.id }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="昵称">
                {{ detail.profile.nickname }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="手机号（脱敏）">
                {{ detail.profile.mobile || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="注册时间">
                {{ formatShanghaiDateTime(detail.profile.created_at) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="账户余额">
                ¥{{ Number(detail.profile.balance).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="积分 / 佣金">
                {{ detail.profile.points }} / ¥{{
                  Number(detail.profile.commission).toFixed(2)
                }}
              </ElDescriptionsItem>
            </ElDescriptions>
            <ElTabs class="mt-4">
              <ElTabPane label="最近订单">
                <ElTable :data="detail.orders" max-height="360">
                  <ElTableColumn label="订单号" min-width="180" prop="order_no" />
                  <ElTableColumn label="店铺" min-width="130" prop="store_name" />
                  <ElTableColumn label="实付" width="100">
                    <template #default="{ row }">
                      ¥{{ Number(row.pay_amount).toFixed(2) }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="件数" prop="total_quantity" width="70" />
                  <ElTableColumn label="状态" width="100">
                    <template #default="{ row }">
                      {{ orderLabels[row.status] || row.status }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="创建时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="资产流水">
                <ElTable :data="detail.assets" max-height="360">
                  <ElTableColumn prop="id" label="ID" width="80" />
                  <ElTableColumn label="类型" width="90">
                    <template #default="{ row }">
                      {{ assetLabels[row.asset_type] || row.asset_type }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="变动" width="110">
                    <template #default="{ row }">
                      <span
                        :class="
                          row.amount < 0 ? 'text-red-500' : 'text-green-600'
                        "
                      >
                        {{ row.amount > 0 ? '+' : ''
                        }}{{ Number(row.amount).toFixed(2) }}
                      </span>
                    </template>
                  </ElTableColumn>
                  <ElTableColumn
                    prop="reference_type"
                    label="业务来源"
                    min-width="120"
                  />
                  <ElTableColumn
                    prop="reference_id"
                    label="业务引用"
                    min-width="120"
                  />
                  <ElTableColumn label="创建时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="会员变更">
                <ElTable :data="detail.membership_logs" max-height="360">
                  <ElTableColumn
                    prop="previous_level_name"
                    label="原等级"
                    min-width="110"
                  />
                  <ElTableColumn prop="level_name" label="当前等级" min-width="110" />
                  <ElTableColumn label="变更类型" width="100">
                    <template #default="{ row }">
                      {{ membershipLabels[row.change_type] || row.change_type }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn prop="note" label="说明" min-width="220" />
                  <ElTableColumn label="时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="签到记录">
                <ElTable :data="detail.signs" max-height="360">
                  <ElTableColumn prop="sign_date" label="签到日期" min-width="140" />
                  <ElTableColumn prop="points" label="获得积分" width="110" />
                  <ElTableColumn
                    prop="continuous_days"
                    label="连续天数"
                    width="110"
                  />
                  <ElTableColumn label="记录时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.created_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="浏览记录">
                <ElTable :data="detail.browse_history" max-height="360">
                  <ElTableColumn prop="product_id" label="商品 ID" width="100" />
                  <ElTableColumn prop="title" label="商品" min-width="180" />
                  <ElTableColumn prop="store_name" label="店铺" min-width="140" />
                  <ElTableColumn label="浏览时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.viewed_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="持有优惠券">
                <ElTable :data="detail.coupons" max-height="360">
                  <ElTableColumn prop="coupon_id" label="券 ID" width="90" />
                  <ElTableColumn prop="name" label="优惠券" min-width="160" />
                  <ElTableColumn label="优惠" width="100">
                    <template #default="{ row }">
                      {{
                        row.discount_type === 'rate'
                          ? `${row.discount_value / 10} 折`
                          : `¥${Number(row.discount_value).toFixed(2)}`
                      }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="门槛" width="100">
                    <template #default="{ row }">
                      ¥{{ Number(row.min_amount).toFixed(2) }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="状态" width="100">
                    <template #default="{ row }">
                      {{ couponStatusLabels[row.status] || row.status }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="领取时间" width="170">
                    <template #default="{ row }">
                      {{ formatShanghaiDateTime(row.obtained_at) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>
              <ElTabPane label="推荐关系">
                <ElDescriptions :column="1" border>
                  <ElDescriptionsItem label="上级用户">
                    {{
                      detail.distribution.parent_user_id
                        ? `${detail.distribution.parent_nickname || '未知用户'}（#${detail.distribution.parent_user_id}）`
                        : '未绑定'
                    }}
                  </ElDescriptionsItem>
                  <ElDescriptionsItem label="直推用户数">
                    {{ detail.distribution.direct_user_count }}
                  </ElDescriptionsItem>
                  <ElDescriptionsItem label="推广员资格">
                    <ElTag
                      :type="
                        detail.distribution.promoter_status === 1
                          ? 'success'
                          : 'info'
                      "
                    >
                      {{
                        detail.distribution.promoter_status === 1
                          ? '启用'
                          : '未开通或已停用'
                      }}
                    </ElTag>
                  </ElDescriptionsItem>
                </ElDescriptions>
              </ElTabPane>
            </ElTabs>
            <ElAlert
              class="mt-4"
              type="info"
              :closable="false"
              title="为保护个人与交易敏感信息，页面不展示收货地址、发票资料、支付交易号或提现账户快照。"
            />
          </template>
        </template>
      </ElSkeleton>
    </DetailDrawer>
  </Page>
</template>
