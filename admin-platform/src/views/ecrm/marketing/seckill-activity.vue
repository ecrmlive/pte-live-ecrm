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
  ElSwitch,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformCategoriesApi,
  type PlatformCategory,
} from '#/api/core/platform-catalog';
import {
  clonePlatformSeckillActivityApi,
  deletePlatformSeckillActivityApi,
  listPlatformSeckillActivitiesApi,
  listPlatformSeckillTimeOptionsApi,
  setPlatformSeckillActivityStatusApi,
  type PlatformSeckillActivity,
  type PlatformSeckillTime,
} from '#/api/core/platform-seckill';
import SeckillActivityFormDrawer from '#/components/marketing/seckill-activity-form-drawer.vue';
import SeckillActivityStatsDrawer from '#/components/marketing/seckill-activity-stats-drawer.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const formDrawerRef = ref<InstanceType<typeof SeckillActivityFormDrawer>>();
const statsDrawerRef = ref<InstanceType<typeof SeckillActivityStatsDrawer>>();
const timeOptions = ref<PlatformSeckillTime[]>([]);
const categoryOptions = ref<{ label: string; value: number }[]>([]);

function toDay(value?: string) {
  const s = String(value || '').trim();
  if (s.length >= 10 && s[4] === '-' && s[7] === '-') return s.slice(0, 10);
  return s;
}

function formatDayRange(start?: string, end?: string) {
  const a = toDay(start);
  const b = toDay(end);
  if (!a && !b) return '—';
  return `${a || '—'} - ${b || '—'}`;
}

function statusText(row: PlatformSeckillActivity) {
  if (row.status_text) return row.status_text;
  if (row.active_status === 1) return '进行中';
  if (row.active_status === 0) return '未开始';
  return '已结束';
}

function flattenCategories(
  nodes: PlatformCategory[],
  prefix = '',
): { label: string; value: number }[] {
  const out: { label: string; value: number }[] = [];
  for (const node of nodes) {
    out.push({
      label: `${prefix}${node.cate_name}`,
      value: node.store_category_id,
    });
    if (node.children?.length) {
      out.push(...flattenCategories(node.children, `${prefix}— `));
    }
  }
  return out;
}

/** 筛选：字段与 CRMEB store_seckill/index 一致；按钮顺序按金标准「重置 → 搜索」 */
const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '活动日期',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入活动名称' },
    fieldName: 'name',
    label: '活动名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '正在进行', value: 1 },
        { label: '未开始', value: 0 },
        { label: '已结束', value: -1 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'active_status',
    label: '活动状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启', value: 1 },
        { label: '关闭', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '是否开启',
  },
]);

const gridOptions: VxeGridProps<PlatformSeckillActivity> = {
  columns: [
    { field: 'seckill_activity_id', title: 'ID', width: 60 },
    {
      field: 'name',
      minWidth: 120,
      showOverflow: false,
      title: '活动名称',
    },
    { field: 'product_count', minWidth: 90, title: '商品数量' },
    { field: 'merchant_count', minWidth: 90, title: '店铺数量' },
    {
      field: 'status_text',
      formatter: ({ row }) => statusText(row),
      minWidth: 90,
      title: '活动状态',
    },
    {
      field: 'date_range',
      formatter: ({ row }) => formatDayRange(row.start_day, row.end_day),
      minWidth: 180,
      showOverflow: false,
      title: '活动日期',
    },
    {
      field: 'seckill_time_texts',
      slots: { default: 'time_texts' },
      minWidth: 100,
      showOverflow: false,
      title: '活动时间',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 150,
      title: '创建时间',
    },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: '活动开关',
      width: 90,
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig({ pageSize: 10 }),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const activeRaw = formValues?.active_status;
        const statusRaw = formValues?.status;
        const data = await listPlatformSeckillActivitiesApi({
          page: page.currentPage,
          limit: page.pageSize,
          name: String(formValues?.name ?? '').trim() || undefined,
          date_from: range[0] as string | undefined,
          date_to: range[1] as string | undefined,
          active_status:
            activeRaw === 0 || activeRaw === 1 || activeRaw === -1
              ? Number(activeRaw)
              : undefined,
          status:
            statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'seckill_activity_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function openCreate() {
  formDrawerRef.value?.open('create');
}

function openEdit(row: PlatformSeckillActivity) {
  void formDrawerRef.value?.open('edit', row);
}

function openDetail(row: PlatformSeckillActivity) {
  void formDrawerRef.value?.open('view', row);
}

function openStats(row: PlatformSeckillActivity) {
  statsDrawerRef.value?.open(row.seckill_activity_id, row.name);
}

async function toggleStatus(row: PlatformSeckillActivity) {
  const next = row.status === 1 ? 0 : 1;
  try {
    await setPlatformSeckillActivityStatusApi(row.seckill_activity_id, next);
    row.status = next;
    ElMessage.success(next === 1 ? '已开启' : '已关闭');
  } catch {
    ElMessage.error('切换失败');
    gridApi.reload();
  }
}

async function cloneRow(row: PlatformSeckillActivity) {
  try {
    await ElMessageBox.confirm(
      `复制活动「${row.name}」？`,
      '复制确认',
      { type: 'info' },
    );
    await clonePlatformSeckillActivityApi(row.seckill_activity_id);
    ElMessage.success('已复制');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

async function remove(row: PlatformSeckillActivity) {
  try {
    await ElMessageBox.confirm(
      `删除活动「${row.name}」后不可恢复，是否继续？`,
      '删除确认',
      { type: 'warning' },
    );
    await deletePlatformSeckillActivityApi(row.seckill_activity_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

async function loadOptions() {
  const [times, cats] = await Promise.all([
    listPlatformSeckillTimeOptionsApi(),
    listPlatformCategoriesApi(),
  ]);
  timeOptions.value = times.list || [];
  categoryOptions.value = flattenCategories(cats.list || []);
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const isStaff = profile.roles.some(
    (r) => r === 'platform' || r === 'operations',
  );
  canManage.value =
    isStaff &&
    (codes.includes('marketing.seckill.activity') ||
      codes.includes('marketing.seckill.dir') ||
      codes.includes('marketing.seckill') ||
      codes.includes('marketing.seckill.manage') ||
      codes.includes('marketing.seckill.manage.page'));
  await loadOptions().catch(() => undefined);
  await gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton type="primary" @click="openCreate">
          新增秒杀活动
        </ElButton>
      </template>

      <template #time_texts="{ row }">
        <div v-if="row.seckill_time_texts?.length" class="time-texts">
          <div v-for="(t, i) in row.seckill_time_texts" :key="i">{{ t }}</div>
        </div>
        <span v-else>—</span>
      </template>

      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          @change="() => toggleStatus(row)"
        />
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">查看</ElButton>
        <ElButton link type="primary" @click="openStats(row)">统计</ElButton>
        <ElDropdown v-if="canManage" trigger="click" class="more-dropdown">
          <span class="more-link">
            更多
            <ElIcon class="more-caret"><ArrowDown /></ElIcon>
          </span>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem @click="openEdit(row)">编辑</ElDropdownItem>
              <ElDropdownItem @click="cloneRow(row)">复制</ElDropdownItem>
              <ElDropdownItem @click="remove(row)">删除</ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <SeckillActivityFormDrawer
      ref="formDrawerRef"
      :time-options="timeOptions"
      :category-options="categoryOptions"
      @saved="gridApi.reload()"
    />
    <SeckillActivityStatsDrawer ref="statsDrawerRef" />
  </Page>
</template>

<style scoped>
.time-texts {
  line-height: 1.45;
  font-variant-numeric: tabular-nums;
}

.more-dropdown {
  margin-left: 8px;
  vertical-align: middle;
}

.more-link {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  color: var(--el-color-primary);
  font-size: 14px;
  cursor: pointer;
  line-height: 1;
}

.more-caret {
  font-size: 12px;
}
</style>
