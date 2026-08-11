<script setup lang="ts">
import { ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElDatePicker,
  ElImage,
  ElMessage,
  ElOption,
  ElPagination,
  ElSelect,
  ElSkeleton,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
} from 'element-plus';

import {
  getPlatformOrderApi,
  listPlatformOrderLogsApi,
  type PlatformOrder,
  type PlatformOrderLogRow,
} from '#/api/core/platform-trade';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';

const detail = ref<PlatformOrder>();
const detailLoading = ref(false);
const detailTab = ref('info');

const logLoading = ref(false);
const logs = ref<PlatformOrderLogRow[]>([]);
const logTotal = ref(0);
const logPage = ref(1);
const logLimit = ref(10);
const logTerminal = ref('');
const logDates = ref<string[]>([]);

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
  title: '订单详情',
});

function money(v?: number) {
  return `¥${Number(v || 0).toFixed(2)}`;
}

function moneyPlain(v?: number) {
  return Number(v || 0).toFixed(2);
}

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '') return '-';
  return String(v);
}

function noneText(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '') return '无';
  return String(v);
}

function statusWarn(row?: PlatformOrder) {
  if (!row) return false;
  return (
    !!row.user_deleted ||
    row.status_label === '待付款' ||
    row.paid === 0
  );
}

function resetLogs() {
  logs.value = [];
  logTotal.value = 0;
  logPage.value = 1;
  logLimit.value = 10;
  logTerminal.value = '';
  logDates.value = [];
}

async function loadLogs() {
  if (!detail.value?.order_id) return;
  logLoading.value = true;
  try {
    const range = Array.isArray(logDates.value) ? logDates.value : [];
    const data = await listPlatformOrderLogsApi(detail.value.order_id, {
      page: logPage.value,
      limit: logLimit.value,
      terminal: logTerminal.value || undefined,
      date_from: range[0],
      date_to: range[1],
    });
    logs.value = data.list || [];
    logTotal.value = data.total || 0;
  } catch {
    logs.value = [];
    logTotal.value = 0;
  } finally {
    logLoading.value = false;
  }
}

async function open(orderId: number): Promise<boolean> {
  const id = Number(orderId || 0);
  if (!id) {
    ElMessage.warning('缺少订单 ID');
    return false;
  }
  detail.value = undefined;
  detailTab.value = 'info';
  detailLoading.value = true;
  resetLogs();
  drawerApi.setState({ title: '订单详情', loading: true }).open();
  try {
    detail.value = await getPlatformOrderApi(id);
    await loadLogs();
    return true;
  } catch {
    ElMessage.error('加载订单详情失败');
    drawerApi.close();
    return false;
  } finally {
    detailLoading.value = false;
    drawerApi.setState({ loading: false });
  }
}

function close() {
  drawerApi.close();
}

function onDetailTabChange(name: string | number) {
  if (String(name) === 'logs') {
    void loadLogs();
  }
}

defineExpose({ open, close });
</script>

<template>
  <Drawer>
    <ElSkeleton :loading="detailLoading" animated :rows="10">
      <template #default>
        <div v-if="detail" class="order-detail">
          <div class="order-detail__header">
            <div class="order-detail__identity">
              <div class="order-detail__icon">
                <IconifyIcon icon="ant-design:file-text-outlined" />
              </div>
              <div class="order-detail__titles">
                <div class="order-detail__type">
                  {{ detail.order_type_label || '普通订单' }}
                </div>
                <div class="order-detail__sn">
                  订单编号：{{ detail.order_sn }}
                </div>
              </div>
            </div>
            <div class="order-detail__status">
              <div class="order-detail__status-item">
                <span class="label">订单状态</span>
                <span
                  class="value"
                  :class="{ 'is-warn': statusWarn(detail) }"
                >
                  {{ detail.status_label || '-' }}
                  <template v-if="detail.user_deleted">（用户已删除）</template>
                </span>
              </div>
              <div class="order-detail__status-item">
                <span class="label">支付方式</span>
                <span class="value">{{ dash(detail.pay_type_label) }}</span>
              </div>
              <div class="order-detail__status-item">
                <span class="label">支付时间</span>
                <span class="value">{{
                  detail.pay_time
                    ? formatShanghaiDateTime(detail.pay_time)
                    : formatShanghaiDateTime(detail.create_time)
                }}</span>
              </div>
            </div>
          </div>

          <ElTabs
            v-model="detailTab"
            class="order-detail__tabs"
            @tab-change="onDetailTabChange"
          >
            <ElTabPane label="订单信息" name="info">
              <div class="order-section order-section--user">
                <div class="order-section__title">用户信息</div>
                <div class="order-kv-grid">
                  <div class="order-kv">
                    <span class="order-kv__label">用户昵称:</span>
                    <span class="order-kv__value">{{
                      dash(detail.nickname)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">用户ID:</span>
                    <span class="order-kv__value">{{
                      detail.uid || '-'
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">绑定电话:</span>
                    <span class="order-kv__value">{{
                      dash(detail.user_phone_mask)
                    }}</span>
                  </div>
                </div>
              </div>

              <div class="order-section order-section--ship">
                <div class="order-section__title">收货信息</div>
                <div class="order-kv-stack">
                  <div class="order-kv">
                    <span class="order-kv__label">收货人:</span>
                    <span class="order-kv__value">{{
                      dash(detail.real_name)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">收货电话:</span>
                    <span class="order-kv__value">{{
                      dash(detail.user_phone)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">收货地址:</span>
                    <span class="order-kv__value">{{
                      dash(detail.user_address)
                    }}</span>
                  </div>
                </div>
              </div>

              <div class="order-section order-section--order">
                <div class="order-section__title">订单信息</div>
                <div class="order-price-box">
                  <div class="order-price-box__formula">
                    商品总价: <em>{{ money(detail.total_price) }}</em>
                    - 平台优惠券:
                    <em>{{ money(detail.platform_coupon) }}</em>
                    - 商家优惠券:
                    <em>{{ money(detail.merchant_coupon) }}</em>
                    - 积分抵扣:
                    <em>{{ money(detail.points_deduction) }}</em>
                    - 会员优惠:
                    <em>{{ money(detail.member_discount) }}</em>
                    + 支付运费: <em>{{ money(detail.freight_price) }}</em>
                  </div>
                  <div class="order-price-box__total">
                    = 订单总金额: <em>{{ money(detail.pay_price) }}</em>
                  </div>
                </div>
                <div class="order-kv-grid order-kv-grid--meta">
                  <div class="order-kv">
                    <span class="order-kv__label">创建时间:</span>
                    <span class="order-kv__value">{{
                      formatShanghaiDateTime(detail.create_time)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">商品总数:</span>
                    <span class="order-kv__value">{{
                      detail.total_num ?? '-'
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">发货方式:</span>
                    <span class="order-kv__value">{{
                      dash(detail.delivery_type_label)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">一级佣金:</span>
                    <span class="order-kv__value">{{
                      Number(detail.first_brokerage || 0)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">二级佣金:</span>
                    <span class="order-kv__value">{{
                      Number(detail.second_brokerage || 0)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">推广人:</span>
                    <span class="order-kv__value">{{
                      noneText(detail.spread_name)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">上级推广人:</span>
                    <span class="order-kv__value">{{
                      noneText(detail.top_spread_name)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">商品类型:</span>
                    <span class="order-kv__value">{{
                      dash(detail.product_type_label)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">活动类型:</span>
                    <span class="order-kv__value">{{
                      dash(detail.activity_type_label)
                    }}</span>
                  </div>
                </div>
              </div>
            </ElTabPane>

            <ElTabPane label="商品信息" name="products">
              <ElTable
                :data="detail.products || []"
                class="order-detail-table"
                border
              >
                <ElTableColumn label="商品ID" prop="product_id" width="90" />
                <ElTableColumn label="商品信息" min-width="300">
                  <template #default="{ row }">
                    <div class="order-product order-product--detail">
                      <ElImage
                        class="order-product__thumb"
                        :src="resolveCosMediaUrl(row.product_image || '')"
                        fit="cover"
                      />
                      <div class="order-product__text">
                        <div class="order-product__name">
                          {{ row.product_info || '-' }}
                        </div>
                        <div class="order-product__sku">
                          规格：{{ row.product_sku || '-' }}
                        </div>
                      </div>
                    </div>
                  </template>
                </ElTableColumn>
                <ElTableColumn label="售价" width="110" align="center">
                  <template #default="{ row }">
                    {{ moneyPlain(row.product_price) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn label="成本价" width="110" align="center">
                  <template #default="{ row }">
                    {{ moneyPlain(row.cost_price) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn label="实付金额" width="110" align="center">
                  <template #default="{ row }">
                    {{ moneyPlain(row.pay_price ?? row.total_price) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn
                  label="购买数量"
                  prop="product_num"
                  width="100"
                  align="center"
                />
              </ElTable>
            </ElTabPane>

            <ElTabPane label="订单记录" name="logs">
              <div class="order-log-filters">
                <div class="order-log-filters__item">
                  <span class="filter-label">操作端：</span>
                  <ElSelect
                    v-model="logTerminal"
                    clearable
                    placeholder="请选择"
                    style="width: 160px"
                    @change="
                      () => {
                        logPage = 1;
                        loadLogs();
                      }
                    "
                  >
                    <ElOption label="用户" value="user" />
                    <ElOption label="店铺" value="merchant" />
                    <ElOption label="系统" value="system" />
                  </ElSelect>
                </div>
                <div class="order-log-filters__item">
                  <span class="filter-label">操作时间：</span>
                  <ElDatePicker
                    v-model="logDates"
                    type="daterange"
                    value-format="YYYY-MM-DD"
                    start-placeholder="开始时间"
                    end-placeholder="结束时间"
                    style="width: 280px"
                    @change="
                      () => {
                        logPage = 1;
                        loadLogs();
                      }
                    "
                  />
                </div>
              </div>
              <ElTable
                v-loading="logLoading"
                :data="logs"
                class="order-detail-table"
                border
              >
                <ElTableColumn
                  label="订单编号"
                  prop="order_sn"
                  min-width="190"
                />
                <ElTableColumn
                  label="操作记录"
                  prop="content"
                  min-width="140"
                />
                <ElTableColumn label="操作角色" prop="role" width="100" />
                <ElTableColumn
                  label="操作人"
                  prop="operator"
                  min-width="120"
                />
                <ElTableColumn label="操作时间" min-width="170">
                  <template #default="{ row }">
                    {{ formatShanghaiDateTime(row.operate_time) }}
                  </template>
                </ElTableColumn>
              </ElTable>
              <div class="order-log-pager">
                <ElPagination
                  v-model:current-page="logPage"
                  v-model:page-size="logLimit"
                  background
                  layout="prev, pager, next, jumper"
                  :total="logTotal"
                  @current-change="loadLogs"
                  @size-change="
                    () => {
                      logPage = 1;
                      loadLogs();
                    }
                  "
                />
              </div>
            </ElTabPane>

            <ElTabPane label="店铺信息" name="store">
              <div class="order-section order-section--last">
                <div class="order-section__title">店铺信息</div>
                <div class="order-kv-grid">
                  <div class="order-kv">
                    <span class="order-kv__label">店铺名称：</span>
                    <span class="order-kv__value">{{
                      dash(detail.store_name || detail.mer_name)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">店铺类型：</span>
                    <span class="order-kv__value">{{
                      dash(detail.store_type_name)
                    }}</span>
                  </div>
                  <div class="order-kv">
                    <span class="order-kv__label">店铺类别：</span>
                    <span class="order-kv__value">{{
                      dash(detail.store_category_name)
                    }}</span>
                  </div>
                </div>
              </div>
            </ElTabPane>
          </ElTabs>
        </div>
      </template>
    </ElSkeleton>
  </Drawer>
</template>

<style scoped>
.order-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-detail__header {
  padding-bottom: 4px;
}

.order-detail__identity {
  display: flex;
  gap: 12px;
  align-items: center;
}

.order-detail__icon {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  color: #fff;
  font-size: 22px;
  background: var(--el-color-primary);
  border-radius: 8px;
}

.order-detail__type {
  color: var(--el-text-color-primary);
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
}

.order-detail__sn {
  margin-top: 2px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 20px;
}

.order-detail__status {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px 24px;
  margin-top: 16px;
  margin-bottom: 4px;
}

.order-detail__status-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-detail__status-item .label {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 18px;
}

.order-detail__status-item .value {
  color: var(--el-text-color-primary);
  font-size: 14px;
  line-height: 22px;
}

.order-detail__status-item .value.is-warn {
  color: #ed6a0c;
  font-weight: 500;
}

.order-detail__tabs {
  min-height: 420px;
}

.order-detail__tabs :deep(.el-tabs__header) {
  margin-bottom: 18px;
}

.order-detail__tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: var(--el-border-color-lighter);
}

.order-detail__tabs :deep(.el-tabs__item) {
  height: 40px;
  padding: 0 18px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  font-weight: 400;
}

.order-detail__tabs :deep(.el-tabs__item.is-active) {
  color: var(--el-color-primary);
  font-weight: 500;
}

.order-detail__tabs :deep(.el-tabs__active-bar) {
  height: 2px;
  background-color: var(--el-color-primary);
}

.order-section {
  padding: 2px 0 16px;
}

.order-section--user {
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.order-section--ship,
.order-section--order {
  margin-top: 16px;
  border-bottom: 1px dashed var(--el-border-color);
}

.order-section__title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  color: var(--el-text-color-primary);
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.order-section__title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 1px;
  background: var(--el-color-primary);
}

.order-kv-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px 20px;
}

.order-kv-grid--meta {
  margin-top: 16px;
}

.order-kv-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.order-kv {
  display: flex;
  gap: 4px;
  min-width: 0;
  font-size: 13px;
  line-height: 22px;
}

.order-kv__label {
  flex-shrink: 0;
  color: var(--el-text-color-regular);
}

.order-kv__value {
  min-width: 0;
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.order-price-box {
  padding: 14px 16px;
  color: var(--el-text-color-regular);
  font-size: 13px;
  line-height: 22px;
  text-align: right;
  background: #f5f7fa;
  border-radius: 4px;
}

.order-price-box__formula {
  text-align: right;
  word-break: break-word;
}

.order-price-box__formula em,
.order-price-box__total em {
  font-style: normal;
  font-weight: 400;
  color: #ed4014;
}

.order-price-box__total {
  margin-top: 8px;
  text-align: right;
}

.order-product {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.order-product__thumb {
  flex: 0 0 48px;
  width: 48px;
  height: 48px;
  border-radius: 4px;
  background: var(--el-fill-color-light);
}

.order-product__text {
  min-width: 0;
  line-height: 1.45;
  word-break: break-all;
}

.order-product__name {
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.order-product__sku {
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.order-product--detail .order-product__thumb {
  flex: 0 0 56px;
  width: 56px;
  height: 56px;
}

.order-detail-table :deep(.el-table__header th) {
  background: #f5f7fa;
  color: var(--el-text-color-regular);
  font-weight: 500;
}

.order-detail-table :deep(.el-table__cell) {
  padding: 12px 0;
}

.order-log-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 16px 24px;
  margin-bottom: 14px;
}

.order-log-filters__item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-label {
  color: var(--el-text-color-regular);
  font-size: 13px;
  white-space: nowrap;
}

.order-log-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 14px;
}

@media (max-width: 960px) {
  .order-detail__status,
  .order-kv-grid {
    grid-template-columns: 1fr;
  }
}
</style>
