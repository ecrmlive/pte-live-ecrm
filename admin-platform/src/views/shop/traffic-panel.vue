<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref, watch } from 'vue';

import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import LiveTrafficApi from '#/api/core/live-traffic';
import {
  RECHARGE_TYPE_FORM_OPTIONS,
  rechargeTypeLabel,
} from '#/constants/traffic-recharge-type';

defineOptions({ name: 'ShopTrafficPanel' });

interface TrafficAccount {
  app_id: number;
  app_name?: string;
  lvb_play_used_gb?: number;
  push_flux_gb_total?: number;
  remain_gb?: number;
  total_gb?: number;
  vod_play_used_gb?: number;
}

interface RechargeRow {
  amount_yuan?: number;
  create_time_text?: string;
  delta_gb?: number;
  invoice_no?: string;
  operator_name?: string;
  recharge_type?: string;
  remark?: string;
  total_gb_after?: number;
}

const props = withDefaults(
  defineProps<{
    appId?: number | string;
    /** 弹窗内嵌模式：表格随内容增高，不做视口高度计算 */
    embedded?: boolean;
  }>(),
  { embedded: false },
);

const loading = ref(false);
const recharging = ref(false);
const account = ref<TrafficAccount | null>(null);

const rechargeForm = reactive({
  amount_yuan: 0,
  delta_gb: 0,
  invoice_no: '',
  recharge_type: 'purchase',
  remark: '',
});

const rechargeTypeOptions = computed(() => RECHARGE_TYPE_FORM_OPTIONS);

function fmtGB(v: null | number | string | undefined) {
  const n = Number(v ?? 0);
  return n.toFixed(n >= 100 ? 1 : 2);
}

async function fetchRechargePage(appId: number, pageSize: number, currentPage: number) {
  const listRes = await LiveTrafficApi.rechargeList({
    app_id: appId,
    list_rows: pageSize,
    page: currentPage,
  });
  const list = (listRes.data?.list ?? {}) as {
    data?: RechargeRow[];
    total?: number;
  };
  return {
    items: list.data ?? [],
    total: Number(list.total ?? 0),
  };
}

const gridOptions: VxeGridProps = {
  border: true,
  columns: [
    { field: 'create_time_text', title: '时间', width: 160 },
    {
      field: 'recharge_type',
      slots: { default: 'rechargeType' },
      title: '类型',
      width: 90,
    },
    {
      field: 'delta_gb',
      slots: { default: 'deltaGb' },
      title: '流量(GB)',
      width: 100,
    },
    {
      field: 'amount_yuan',
      slots: { default: 'amountYuan' },
      title: '金额(元)',
      width: 100,
    },
    {
      field: 'total_gb_after',
      slots: { default: 'totalGbAfter' },
      title: '充值后总量',
      width: 110,
    },
    { field: 'operator_name', title: '操作人', width: 100 },
    { field: 'invoice_no', minWidth: 120, title: '发票/合同' },
    {
      field: 'remark',
      minWidth: 140,
      showOverflow: true,
      title: '备注',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 30] },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const id = Number(props.appId);
        if (id <= 0) {
          return { items: [], total: 0 };
        }
        return fetchRechargePage(id, page.pageSize, page.currentPage);
      },
    },
  },
  rowConfig: { isHover: true },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function reloadRechargeList() {
  gridApi.reload();
}

async function loadAccount(appId: number) {
  loading.value = true;
  try {
    const accRes = await LiveTrafficApi.account({ app_id: appId });
    account.value = (accRes.data ?? null) as TrafficAccount | null;
    reloadRechargeList();
  } catch {
    account.value = null;
    reloadRechargeList();
  } finally {
    loading.value = false;
  }
}

async function submitRecharge() {
  if (!account.value) {
    return;
  }
  const deltaGb = Number(rechargeForm.delta_gb);
  if (!deltaGb) {
    ElMessage.warning('请填写流量 GB');
    return;
  }
  recharging.value = true;
  try {
    await LiveTrafficApi.recharge({
      amount_yuan: rechargeForm.amount_yuan,
      app_id: account.value.app_id,
      delta_gb: deltaGb,
      invoice_no: rechargeForm.invoice_no,
      recharge_type: rechargeForm.recharge_type,
      remark: rechargeForm.remark,
    });
    ElMessage.success('提交成功');
    rechargeForm.delta_gb = 0;
    await loadAccount(account.value.app_id);
  } finally {
    recharging.value = false;
  }
}

watch(
  () => props.appId,
  (val) => {
    const id = Number(val);
    if (id > 0) {
      void loadAccount(id);
      return;
    }
    account.value = null;
    reloadRechargeList();
  },
  { immediate: true },
);
</script>

<template>
  <div
    v-loading="loading"
    :class="['traffic-panel', embedded ? 'traffic-panel--embedded' : '']"
  >
    <div v-if="account" class="traffic-panel__stats">
      <div class="traffic-stat-card">
        <p class="traffic-stat-card__label">充值总量</p>
        <p class="traffic-stat-card__value">{{ fmtGB(account.total_gb) }} GB</p>
      </div>
      <div class="traffic-stat-card">
        <p class="traffic-stat-card__label">云直播已用</p>
        <p class="traffic-stat-card__value is-warn">
          {{ fmtGB(account.lvb_play_used_gb) }} GB
        </p>
      </div>
      <div class="traffic-stat-card">
        <p class="traffic-stat-card__label">云点播已用</p>
        <p class="traffic-stat-card__value is-warn">
          {{ fmtGB(account.vod_play_used_gb) }} GB
        </p>
      </div>
      <div class="traffic-stat-card">
        <p class="traffic-stat-card__label">推流已用</p>
        <p class="traffic-stat-card__value is-warn">
          {{ fmtGB(account.push_flux_gb_total) }} GB
        </p>
      </div>
      <div class="traffic-stat-card">
        <p class="traffic-stat-card__label">剩余</p>
        <p
          class="traffic-stat-card__value"
          :class="Number(account.remain_gb) >= 0 ? 'is-ok' : 'is-warn'"
        >
          {{ fmtGB(account.remain_gb) }} GB
        </p>
      </div>
    </div>

    <div v-if="account" class="traffic-panel__recharge">
      <h4 class="traffic-panel__section-title">充值 / 调账</h4>
      <ElForm :inline="true" :model="rechargeForm" class="traffic-recharge-form" size="small">
        <ElFormItem label="流量(GB)">
          <ElInputNumber v-model="rechargeForm.delta_gb" :step="10" />
        </ElFormItem>
        <ElFormItem label="金额(元)">
          <ElInputNumber
            v-model="rechargeForm.amount_yuan"
            :min="0"
            :precision="2"
          />
        </ElFormItem>
        <ElFormItem label="类型">
          <ElSelect v-model="rechargeForm.recharge_type" style="width: 120px">
            <ElOption
              v-for="opt in rechargeTypeOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="发票/合同">
          <ElInput v-model="rechargeForm.invoice_no" style="width: 140px" />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model="rechargeForm.remark" style="width: 160px" />
        </ElFormItem>
        <ElFormItem>
          <ElButton :loading="recharging" type="primary" @click="submitRecharge">
            提交
          </ElButton>
        </ElFormItem>
      </ElForm>
    </div>

    <div v-if="account" class="traffic-panel__log">
      <h4 class="traffic-panel__section-title">充值流水</h4>
      <Grid class="traffic-recharge-grid">
        <template #rechargeType="{ row }">
          {{ rechargeTypeLabel(row.recharge_type) }}
        </template>

        <template #deltaGb="{ row }">
          {{ fmtGB(row.delta_gb) }}
        </template>

        <template #amountYuan="{ row }">
          {{ Number(row.amount_yuan ?? 0).toFixed(2) }}
        </template>

        <template #totalGbAfter="{ row }">
          {{ fmtGB(row.total_gb_after) }}
        </template>
      </Grid>
    </div>
  </div>
</template>

<style scoped>
.traffic-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.traffic-panel__stats {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
}

.traffic-stat-card {
  padding: 10px 12px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: var(--el-fill-color-light);
}

.traffic-stat-card__label {
  margin: 0 0 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.traffic-stat-card__value {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 18px;
  font-weight: 600;
  line-height: 1.2;
}

.traffic-stat-card__value.is-warn {
  color: #e6a23c;
}

.traffic-stat-card__value.is-ok {
  color: #67c23a;
}

.traffic-panel__recharge,
.traffic-panel__log {
  padding: 12px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: var(--el-bg-color);
}

.traffic-panel__section-title {
  margin: 0 0 10px;
  color: var(--el-text-color-primary);
  font-size: 14px;
  font-weight: 600;
}

.traffic-recharge-form {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 0;
}

.traffic-recharge-grid :deep(.h-full.rounded-md.bg-card),
.traffic-recharge-grid :deep(.vxe-grid) {
  height: auto !important;
  min-height: 0 !important;
}

.traffic-recharge-grid :deep(.vxe-table--body-inner-wrapper),
.traffic-recharge-grid :deep(.vxe-table--body-wrapper),
.traffic-recharge-grid :deep(.vxe-grid--layout-body-content-wrapper) {
  height: auto !important;
  min-height: 0 !important;
  max-height: none !important;
  overflow: visible !important;
}

@media (max-width: 960px) {
  .traffic-panel__stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
