<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElOption,
  ElPagination,
  ElSelect,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
} from 'element-plus';
import {
  CircleCheckFilled,
  ShoppingCart,
  UserFilled,
  WalletFilled,
} from '@element-plus/icons-vue';

import { fetchPlatformMerchants } from '#/api/core/ecrm';
import {
  getPlatformSeckillActivityStatsApi,
  listPlatformSeckillActivityStatOrdersApi,
  listPlatformSeckillActivityStatPeopleApi,
  listPlatformSeckillActivityStatProductsApi,
  type PlatformSeckillActivityProduct,
  type PlatformSeckillActivityStatOrder,
  type PlatformSeckillActivityStatPeople,
  type PlatformSeckillActivityStats,
} from '#/api/core/platform-seckill';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

type StatTab = 'participants' | 'order' | 'product';

const activityID = ref(0);
const activityName = ref('');
const loading = ref(false);
const tab = ref<StatTab>('participants');
const panel = ref<PlatformSeckillActivityStats | null>(null);
const peopleRows = ref<PlatformSeckillActivityStatPeople[]>([]);
const orderRows = ref<PlatformSeckillActivityStatOrder[]>([]);
const productRows = ref<PlatformSeckillActivityProduct[]>([]);
const total = ref(0);
const merchantOptions = ref<{ label: string; value: number }[]>([]);

const filters = reactive({
  keyword: '',
  date_range: [] as string[],
  mer_id: undefined as number | undefined,
  status: undefined as number | undefined,
  page: 1,
  limit: 10,
});

const orderStatusOptions = [
  { label: '待付款', value: 1 },
  { label: '待发货', value: 2 },
  { label: '待收货', value: 3 },
  { label: '待评价', value: 4 },
  { label: '交易完成', value: 5 },
  { label: '已退款', value: 6 },
  { label: '待核销', value: 8 },
];

const kpiCards = computed(() => [
  {
    name: '下单人数',
    count: panel.value?.orders_people_count ?? 0,
    unit: '人',
    icon: UserFilled,
    tone: 'blue',
  },
  {
    name: '支付订单额',
    count: formatMoney(panel.value?.pay_order_money),
    unit: '元',
    icon: WalletFilled,
    tone: 'orange',
  },
  {
    name: '支付人数',
    count: panel.value?.pay_order_people_count ?? 0,
    unit: '人',
    icon: CircleCheckFilled,
    tone: 'green',
  },
  {
    name: '支付订单数',
    count: panel.value?.pay_order_count ?? 0,
    unit: '笔',
    icon: ShoppingCart,
    tone: 'purple',
  },
]);

const keywordPlaceholder = computed(() =>
  tab.value === 'product'
    ? '请输入商品名称/ID'
    : '请输入用户姓名/手机/用户ID',
);

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  cancelText: '关闭',
});

function formatMoney(value?: number) {
  return Number(value || 0).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

function buildParams() {
  const range = Array.isArray(filters.date_range) ? filters.date_range : [];
  return {
    page: filters.page,
    limit: filters.limit,
    keyword: filters.keyword.trim() || undefined,
    date_from:
      tab.value !== 'product' ? (range[0] as string | undefined) : undefined,
    date_to:
      tab.value !== 'product' ? (range[1] as string | undefined) : undefined,
    mer_id: filters.mer_id || undefined,
    status:
      tab.value === 'order' && filters.status
        ? Number(filters.status)
        : undefined,
  };
}

async function loadMerchants() {
  if (merchantOptions.value.length) return;
  try {
    const data = await fetchPlatformMerchants({ page: 1, limit: 200 });
    merchantOptions.value = (data.list || [])
      .map((m) => ({
        label: String(m.mer_name || m.mer_id || ''),
        value: Number(m.mer_id || 0),
      }))
      .filter((m) => m.value > 0 && m.label);
  } catch {
    merchantOptions.value = [
      { label: 'CRM Live服饰旗舰店', value: 1 },
      { label: 'CRM Live居家优选店', value: 2 },
    ];
  }
}

async function loadPanel() {
  panel.value = await getPlatformSeckillActivityStatsApi(activityID.value, {
    mer_id: filters.mer_id || undefined,
  });
}

async function loadList() {
  if (!activityID.value) return;
  loading.value = true;
  try {
    const params = buildParams();
    if (tab.value === 'participants') {
      const data = await listPlatformSeckillActivityStatPeopleApi(
        activityID.value,
        params,
      );
      peopleRows.value = data.list || [];
      total.value = data.total || 0;
    } else if (tab.value === 'order') {
      const data = await listPlatformSeckillActivityStatOrdersApi(
        activityID.value,
        params,
      );
      orderRows.value = data.list || [];
      total.value = data.total || 0;
    } else {
      const data = await listPlatformSeckillActivityStatProductsApi(
        activityID.value,
        params,
      );
      productRows.value = data.list || [];
      total.value = data.total || 0;
    }
  } finally {
    loading.value = false;
  }
}

async function reloadAll() {
  await Promise.all([loadPanel(), loadList()]);
}

function resetFilters() {
  filters.keyword = '';
  filters.date_range = [];
  filters.mer_id = undefined;
  filters.status = undefined;
  filters.page = 1;
  void reloadAll();
}

function search() {
  filters.page = 1;
  void reloadAll();
}

function onTabChange(name: string | number) {
  tab.value = String(name) as StatTab;
  filters.page = 1;
  void loadList();
}

async function open(id: number, name: string) {
  activityID.value = id;
  activityName.value = name;
  tab.value = 'participants';
  filters.keyword = '';
  filters.date_range = [];
  filters.mer_id = undefined;
  filters.status = undefined;
  filters.page = 1;
  drawerApi.setState({ title: `秒杀统计 · ${name}` }).open();
  await loadMerchants();
  await reloadAll();
}

defineExpose({ open });
</script>

<template>
  <Drawer>
    <div class="stats-page" v-loading="loading">
      <ElForm inline class="stats-filter" label-width="85px" @submit.prevent>
        <ElFormItem v-if="tab === 'order'" label="订单状态">
          <ElSelect
            v-model="filters.status"
            clearable
            placeholder="请选择订单状态"
            class="filter-ctl"
          >
            <ElOption
              v-for="opt in orderStatusOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="搜索">
          <ElInput
            v-model="filters.keyword"
            clearable
            :placeholder="keywordPlaceholder"
            class="filter-ctl"
            @keyup.enter="search"
          />
        </ElFormItem>
        <ElFormItem v-if="tab !== 'product'" label="活动日期">
          <ElDatePicker
            v-model="filters.date_range"
            type="daterange"
            value-format="YYYY-MM-DD"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            class="filter-date"
          />
        </ElFormItem>
        <ElFormItem label="店铺名称">
          <ElSelect
            v-model="filters.mer_id"
            clearable
            filterable
            placeholder="请选择"
            class="filter-ctl"
          >
            <ElOption
              v-for="opt in merchantOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton @click="resetFilters">重置</ElButton>
          <ElButton type="primary" @click="search">搜索</ElButton>
        </ElFormItem>
      </ElForm>

      <div class="kpi-row">
        <div
          v-for="card in kpiCards"
          :key="card.name"
          class="kpi-card"
          :class="`tone-${card.tone}`"
        >
          <div class="kpi-icon">
            <component :is="card.icon" />
          </div>
          <div class="kpi-body">
            <div class="kpi-count">
              {{ card.count }}
              <span class="kpi-unit">{{ card.unit }}</span>
            </div>
            <div class="kpi-name">{{ card.name }}</div>
          </div>
        </div>
      </div>

      <div class="stats-table-card">
        <ElTabs v-model="tab" @tab-change="onTabChange">
          <ElTabPane label="活动参与人" name="participants" />
          <ElTabPane label="活动订单" name="order" />
          <ElTabPane label="活动商品" name="product" />
        </ElTabs>

        <ElTable
          v-if="tab === 'participants'"
          :data="peopleRows"
          size="small"
          empty-text="暂无数据"
        >
          <ElTableColumn label="用户姓名" min-width="140">
            <template #default="{ row }">
              <span v-if="row.uid">[{{ row.uid }}]</span>{{ row.nickname }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="sum_total_num" label="购买件数" min-width="90" />
          <ElTableColumn prop="order_count" label="支付订单数" min-width="100" />
          <ElTableColumn label="支付金额" min-width="100">
            <template #default="{ row }">
              {{ formatMoney(row.sum_pay_price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="最近参与时间" min-width="170">
            <template #default="{ row }">
              {{
                row.last_join_time
                  ? formatShanghaiDateTime(row.last_join_time)
                  : '—'
              }}
            </template>
          </ElTableColumn>
        </ElTable>

        <ElTable
          v-else-if="tab === 'order'"
          :data="orderRows"
          size="small"
          empty-text="暂无数据"
        >
          <ElTableColumn prop="order_sn" label="订单号" min-width="150" />
          <ElTableColumn label="用户" min-width="120">
            <template #default="{ row }">
              <span v-if="row.uid">[{{ row.uid }}]</span>{{ row.nickname }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="订单状态" min-width="100">
            <template #default="{ row }">
              {{ row.status_text || '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="订单支付金额" min-width="110">
            <template #default="{ row }">
              {{ formatMoney(row.pay_price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="total_num" label="订单商品数" min-width="100" />
          <ElTableColumn label="下单时间" min-width="170">
            <template #default="{ row }">
              {{
                row.create_time
                  ? formatShanghaiDateTime(row.create_time)
                  : '—'
              }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="支付时间" min-width="170">
            <template #default="{ row }">
              {{
                row.pay_time ? formatShanghaiDateTime(row.pay_time) : '—'
              }}
            </template>
          </ElTableColumn>
        </ElTable>

        <ElTable
          v-else
          :data="productRows"
          size="small"
          empty-text="暂无数据"
        >
          <ElTableColumn prop="product_id" label="ID" width="80" />
          <ElTableColumn label="商品图片" width="90">
            <template #default="{ row }">
              <ElImage
                v-if="row.image"
                :src="resolveCosMediaUrl(row.image)"
                fit="cover"
                class="prod-thumb"
              />
              <div v-else class="prod-thumb prod-thumb--empty">—</div>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="name" label="商品名称" min-width="180" />
          <ElTableColumn label="分类" min-width="100">
            <template #default="{ row }">
              {{ row.category_name || '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="店铺" min-width="120">
            <template #default="{ row }">
              {{ row.mer_name || row.mer_id || '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="售价" min-width="90">
            <template #default="{ row }">
              {{ formatMoney(row.price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="stock" label="限量" min-width="80" />
          <ElTableColumn label="秒杀价" min-width="90">
            <template #default="{ row }">
              {{ formatMoney(row.seckill_price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="sales" label="秒杀销量" min-width="90" />
          <ElTableColumn label="活动场次" min-width="110">
            <template #default="{ row }">
              <div v-if="row.seckill_time_texts?.length" class="time-texts">
                <div
                  v-for="(t, i) in row.seckill_time_texts"
                  :key="i"
                >
                  {{ t }}
                </div>
              </div>
              <span v-else>—</span>
            </template>
          </ElTableColumn>
        </ElTable>

        <div class="pager">
          <ElPagination
            background
            layout="total, prev, pager, next, jumper"
            :total="total"
            :page-size="filters.limit"
            :current-page="filters.page"
            @current-change="
              (p: number) => {
                filters.page = p;
                void loadList();
              }
            "
          />
        </div>
      </div>
    </div>
  </Drawer>
</template>

<style scoped>
.stats-page {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.stats-filter {
  padding: 12px 12px 0;
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
  background: hsl(var(--card));
}

.filter-ctl {
  width: 200px;
}

.filter-date {
  width: 260px;
}

.kpi-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

@media (min-width: 900px) {
  .kpi-row {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

.kpi-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 8px;
  border: 1px solid hsl(var(--border));
  background: hsl(var(--card));
}

.kpi-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 10px;
  color: #fff;
  font-size: 22px;
}

.tone-blue .kpi-icon {
  background: #409eff;
}
.tone-orange .kpi-icon {
  background: #e6a23c;
}
.tone-green .kpi-icon {
  background: #67c23a;
}
.tone-purple .kpi-icon {
  background: #9b59b6;
}

.kpi-count {
  font-size: 22px;
  font-weight: 600;
  line-height: 1.2;
}

.kpi-unit {
  margin-left: 4px;
  font-size: 12px;
  font-weight: 400;
  color: hsl(var(--muted-foreground));
}

.kpi-name {
  margin-top: 4px;
  font-size: 13px;
  color: hsl(var(--muted-foreground));
}

.stats-table-card {
  padding: 8px 12px 12px;
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
  background: hsl(var(--card));
}

.pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.prod-thumb {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  overflow: hidden;
}

.prod-thumb--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: hsl(var(--muted) / 0.4);
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.time-texts {
  line-height: 1.45;
  font-variant-numeric: tabular-nums;
}
</style>
