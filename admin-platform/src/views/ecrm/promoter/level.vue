<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElTable,
  ElTableColumn,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createDistributionLevelApi,
  deleteDistributionLevelApi,
  getDistributionLevelApi,
  listDistributionLevelsApi,
  updateDistributionLevelConfigApi,
  type DistributionLevel,
  type DistributionLevelSaveInput,
  type DistributionLevelTaskItem,
  type DistributionLevelTaskRule,
} from '#/api/core/platform-spread';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit' | 'view';

type TaskKey = keyof DistributionLevelTaskRule;

type TaskRow = {
  info: string;
  key: TaskKey;
  label: string;
  name: string;
  num: number;
  unit: string;
};

const TASK_META: Array<{
  key: TaskKey;
  label: string;
  listLabel: string;
  unit: string;
}> = [
  {
    key: 'spread_user',
    label: '邀请好友成为下线',
    listLabel: '推广人数',
    unit: '人',
  },
  {
    key: 'pay_money',
    label: '自身消费金额',
    listLabel: '自身消费金额',
    unit: '元',
  },
  {
    key: 'pay_num',
    label: '自身消费订单数',
    listLabel: '自身下单',
    unit: '个',
  },
  {
    key: 'spread_money',
    label: '下级消费金额',
    listLabel: '推广订单金额',
    unit: '元',
  },
  {
    key: 'spread_pay_num',
    label: '下级消费订单数',
    listLabel: '推广订单',
    unit: '个',
  },
];

/** 预设徽章（SVG data URL）；选中后写入 icon_url，也可改用素材库。 */
const PRESET_ICONS = [
  svgLineIcon('crown'),
  svgLineIcon('grid'),
  svgLineIcon('user'),
  svgBadge('#F59E0B', '#FB923C', 'L1', 'l1'),
  svgBadge('#EF4444', '#F97316', '双', 'shuang'),
  svgBadge('#F59E0B', '#FBBF24', '庆', 'qing'),
];

function svgData(inner: string) {
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64">${inner}</svg>`,
  )}`;
}

function svgLineIcon(kind: 'crown' | 'grid' | 'user') {
  const stroke = '#F59E0B';
  if (kind === 'crown') {
    return svgData(
      `<path d="M12 44h40v6H12zM14 44l4-22 10 12 4-18 4 18 10-12 4 22"
        fill="none" stroke="${stroke}" stroke-width="3" stroke-linejoin="round" stroke-linecap="round"/>
       <circle cx="18" cy="18" r="3" fill="${stroke}"/>
       <circle cx="32" cy="12" r="3" fill="${stroke}"/>
       <circle cx="46" cy="18" r="3" fill="${stroke}"/>`,
    );
  }
  if (kind === 'grid') {
    return svgData(
      `<rect x="12" y="12" width="16" height="16" rx="3" fill="none" stroke="${stroke}" stroke-width="3"/>
       <rect x="36" y="12" width="16" height="16" rx="3" fill="none" stroke="${stroke}" stroke-width="3"/>
       <rect x="12" y="36" width="16" height="16" rx="3" fill="none" stroke="${stroke}" stroke-width="3"/>
       <rect x="36" y="36" width="16" height="16" rx="3" fill="none" stroke="${stroke}" stroke-width="3"/>`,
    );
  }
  return svgData(
    `<circle cx="32" cy="24" r="10" fill="none" stroke="${stroke}" stroke-width="3"/>
     <path d="M14 52c2-10 10-16 18-16s16 6 18 16" fill="none" stroke="${stroke}" stroke-width="3" stroke-linecap="round"/>`,
  );
}

function svgBadge(c1: string, c2: string, text: string, id: string) {
  // ASCII-only gradient id：中文 encode 会产生 %，在 SVG url(#…) 里失效，徽章整块空白。
  const gid = `lv_badge_${id}`;
  return svgData(
    `<defs><linearGradient id="${gid}" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0%" stop-color="${c1}"/>
      <stop offset="100%" stop-color="${c2}"/>
    </linearGradient></defs>
    <circle cx="32" cy="32" r="28" fill="url(#${gid})"/>
    <circle cx="32" cy="32" r="22" fill="none" stroke="#fff" stroke-width="2" opacity=".55"/>
    <text x="32" y="39" text-anchor="middle" font-size="20" font-family="PingFang SC,Microsoft YaHei,sans-serif" fill="#fff" font-weight="700">${text}</text>`,
  );
}

function isPresetIcon(url: string) {
  return PRESET_ICONS.includes(url);
}

function emptyTaskItem(): DistributionLevelTaskItem {
  return { info: '', name: '', num: 0 };
}

function emptyTaskRule(): DistributionLevelTaskRule {
  return {
    pay_money: emptyTaskItem(),
    pay_num: emptyTaskItem(),
    spread_money: emptyTaskItem(),
    spread_pay_num: emptyTaskItem(),
    spread_user: emptyTaskItem(),
  };
}

function emptyForm(): DistributionLevelSaveInput {
  return {
    extension_one: 0,
    extension_two: 0,
    icon_url: PRESET_ICONS[0] || '',
    name: '',
    rank: 1,
    status: 1,
    task_rule: emptyTaskRule(),
  };
}

const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const form = reactive<DistributionLevelSaveInput>(emptyForm());
const taskRows = ref<TaskRow[]>([]);
const isReadonly = computed(() => drawerMode.value === 'view');

/** 素材库槽：仅展示非预设 URL；选预设时保持为空「+」，避免与第一枚预设重复。 */
const materialIconUrl = computed({
  get() {
    const url = String(form.icon_url || '').trim();
    if (!url || isPresetIcon(url)) return '';
    return url;
  },
  set(value: string) {
    const next = String(value || '').trim();
    form.icon_url = next || PRESET_ICONS[0] || '';
  },
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入等级名称' },
    fieldName: 'name',
    label: '等级名称',
  },
]);

const gridOptions: VxeGridProps<DistributionLevel> = {
  // 列宽策略对齐金标准店铺列表：固定窄列用 width，弹性列用 minWidth；
  // 操作列 fixed:right；scrollX 走 PLATFORM_LIST_GRID_LAYOUT 默认（勿本页改 gt）。
  columns: [
    { field: 'id', title: 'ID', width: 72 },
    {
      field: 'icon_url',
      slots: { default: 'icon' },
      title: '图标',
      width: 80,
    },
    {
      field: 'name',
      minWidth: 120,
      showOverflow: false,
      title: '名称',
    },
    {
      field: 'rank',
      formatter: ({ cellValue }) => `Lv ${Number(cellValue || 0)}`,
      title: '等级',
      width: 80,
    },
    {
      // 弹性主列：吸收容器剩余宽度（依赖表容器 width:100%，见 platform-list-page.scss）
      field: 'task_rule',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'tasks' },
      title: '分销任务',
    },
    {
      field: 'promoter_count',
      formatter: ({ cellValue }) => String(Number(cellValue || 0)),
      minWidth: 108,
      title: '分销员人数',
    },
    {
      field: 'extension_one',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      minWidth: 168,
      title: '一级返佣上浮比例(%)',
    },
    {
      field: 'extension_two',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      minWidth: 168,
      title: '二级返佣上浮比例(%)',
    },
    platformListActionColumn({ width: 168 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const name = String(formValues?.name ?? '').trim() || undefined;
        const data = await listDistributionLevelsApi({
          limit: page.pageSize,
          name,
          page: page.currentPage,
        });
        return {
          items: data.list || [],
          total: Number(data.total || 0),
        };
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

const [LevelDrawer, levelDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (drawerMode.value === 'view') return;
    await save();
  },
});

function syncTaskRowsFromRule(rule: DistributionLevelTaskRule) {
  taskRows.value = TASK_META.map((meta) => {
    const item = rule[meta.key] || emptyTaskItem();
    return {
      info: item.info || '',
      key: meta.key,
      label: meta.label,
      name: item.name || '',
      num: Number(item.num || 0),
      unit: meta.unit,
    };
  });
}

function syncRuleFromTaskRows() {
  const rule = emptyTaskRule();
  for (const row of taskRows.value) {
    rule[row.key] = {
      info: String(row.info || '').trim(),
      name: String(row.name || '').trim(),
      num: Number(row.num || 0),
    };
  }
  form.task_rule = rule;
}

function applyForm(row?: DistributionLevel) {
  const next = emptyForm();
  if (row) {
    next.name = row.name || '';
    next.rank = Number(row.rank || 1);
    next.icon_url = row.icon_url || '';
    next.extension_one = Number(row.extension_one || 0);
    next.extension_two = Number(row.extension_two || 0);
    next.status = row.status ?? 1;
    next.task_rule = {
      ...emptyTaskRule(),
      ...(row.task_rule || {}),
    };
  }
  Object.assign(form, next);
  syncTaskRowsFromRule(form.task_rule);
}

function formatTaskNum(num: number): string {
  if (!Number.isFinite(num)) return '0';
  return Number.isInteger(num) ? String(num) : String(num);
}

function taskLines(row: DistributionLevel): string[] {
  const rule = row.task_rule || emptyTaskRule();
  const lines: string[] = [];
  for (const meta of TASK_META) {
    const item = rule[meta.key];
    const num = Number(item?.num || 0);
    if (num > 0) {
      lines.push(`${meta.listLabel}${formatTaskNum(num)}${meta.unit}`);
    }
  }
  return lines;
}

function iconSrc(url?: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function selectPresetIcon(url: string) {
  if (isReadonly.value) return;
  form.icon_url = url;
}

function openCreate() {
  drawerMode.value = 'create';
  editingId.value = 0;
  applyForm();
  levelDrawerApi
    .setState({
      title: '新增分销员等级',
      showConfirmButton: true,
      confirmText: '保存',
      cancelText: '取消',
    })
    .open();
}

async function openEdit(row: DistributionLevel) {
  drawerMode.value = 'edit';
  editingId.value = row.id;
  levelDrawerApi
    .setState({
      loading: true,
      title: '编辑分销员等级',
      showConfirmButton: true,
      confirmText: '保存',
      cancelText: '取消',
    })
    .open();
  try {
    const detail = await getDistributionLevelApi(row.id);
    applyForm(detail);
  } finally {
    levelDrawerApi.setState({ loading: false });
  }
}

async function openDetail(row: DistributionLevel) {
  drawerMode.value = 'view';
  editingId.value = row.id;
  levelDrawerApi
    .setState({
      loading: true,
      title: '分销员等级详情',
      showConfirmButton: false,
      cancelText: '关闭',
    })
    .open();
  try {
    const detail = await getDistributionLevelApi(row.id);
    applyForm(detail);
  } finally {
    levelDrawerApi.setState({ loading: false });
  }
}

async function save() {
  syncRuleFromTaskRows();
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请输入等级名称');
    return;
  }
  if (!form.rank || form.rank < 1) {
    ElMessage.warning('请输入等级');
    return;
  }
  const rule = form.task_rule;
  const nums = [
    rule.spread_user.num,
    rule.pay_money.num,
    rule.pay_num.num,
    rule.spread_money.num,
    rule.spread_pay_num.num,
  ];
  if (!nums.some((n) => Number(n) > 0)) {
    ElMessage.warning('请至少输入一个等级任务');
    return;
  }
  const pairs: Array<[number, string]> = [
    [rule.spread_user.num, rule.spread_user.name],
    [rule.pay_money.num, rule.pay_money.name],
    [rule.pay_num.num, rule.pay_num.name],
    [rule.spread_money.num, rule.spread_money.name],
    [rule.spread_pay_num.num, rule.spread_pay_num.name],
  ];
  for (const [num, taskName] of pairs) {
    const hasNum = Number(num) > 0;
    const hasName = Boolean(String(taskName || '').trim());
    if (hasNum !== hasName) {
      ElMessage.warning('请输入相对应的任务或数量');
      return;
    }
  }
  if (
    form.extension_one < 0 ||
    form.extension_one > 1000 ||
    form.extension_two < 0 ||
    form.extension_two > 1000
  ) {
    ElMessage.warning('返佣上浮比例须在 0-1000 之间');
    return;
  }
  const payload: DistributionLevelSaveInput = {
    extension_one: Number(form.extension_one || 0),
    extension_two: Number(form.extension_two || 0),
    icon_url: String(form.icon_url || '').trim(),
    name,
    rank: Number(form.rank),
    status: 1,
    task_rule: rule,
  };
  levelDrawerApi.lock();
  try {
    if (drawerMode.value === 'edit' && editingId.value) {
      await updateDistributionLevelConfigApi(editingId.value, payload);
    } else {
      await createDistributionLevelApi(payload);
    }
    levelDrawerApi.close();
    ElMessage.success('分销员等级已保存');
    gridApi.reload();
  } finally {
    levelDrawerApi.unlock();
  }
}

async function remove(row: DistributionLevel) {
  if (Number(row.promoter_count || 0) > 0) {
    ElMessage.warning('该等级下有数据，不能进行删除操作！');
    return;
  }
  try {
    await confirm({
      content: `确定删除等级「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteDistributionLevelApi(row.id);
    ElMessage.success('删除成功');
    gridApi.reload();
  } catch {
    /* cancelled */
  }
}

const oneRateHint =
  '在分销一级佣金基础上上浮 (0-1000之间整数) 百分比，目前一级返佣比率：3%，上浮10%，则返佣比率：一级返佣比率 * (1 + 一级上浮比率) = 3.30%';
const twoRateHint =
  '在分销二级佣金基础上上浮 (0-1000之间整数) 百分比，目前二级返佣比率：2%，上浮1%，则返佣比率：二级返佣比率 * (1 + 二级上浮比率) = 2.02%';
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增分销员等级
        </ElButton>
      </template>

      <template #icon="{ row }">
        <ElImage
          v-if="row.icon_url"
          class="level-list-icon"
          :src="iconSrc(row.icon_url)"
          fit="contain"
          alt="等级图标"
        >
          <template #error>
            <span>—</span>
          </template>
        </ElImage>
        <span v-else>—</span>
      </template>

      <template #tasks="{ row }">
        <template v-if="taskLines(row).length">
          <div v-for="line in taskLines(row)" :key="line">{{ line }}</div>
        </template>
        <span v-else>—</span>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <LevelDrawer>
      <ElForm
        label-width="156px"
        class="level-form"
        require-asterisk-position="left"
      >
        <ElFormItem label="等级名称：" required>
          <ElInput
            v-model="form.name"
            :disabled="isReadonly"
            maxlength="64"
            show-word-limit
            placeholder="请输入等级名称"
          />
        </ElFormItem>
        <ElFormItem label="等级：" required>
          <ElInputNumber
            v-model="form.rank"
            :disabled="isReadonly"
            :min="1"
            :max="999"
            :controls="false"
            class="level-form__rank"
          />
        </ElFormItem>
        <ElFormItem label="图标：">
          <div class="icon-picker">
            <button
              v-for="(url, idx) in PRESET_ICONS"
              :key="`preset-${idx}`"
              type="button"
              class="icon-picker__item"
              :class="{ 'is-active': form.icon_url === url }"
              :disabled="isReadonly"
              :title="`预设图标 ${idx + 1}`"
              @click="selectPresetIcon(url)"
            >
              <img :src="url" alt="" />
            </button>
            <div
              class="icon-picker__custom"
              :class="{ 'is-active': Boolean(materialIconUrl) }"
            >
              <ImageField
                v-model="materialIconUrl"
                :disabled="isReadonly"
                default-library="system"
                :preview-size="48"
              />
              <span class="icon-picker__custom-label">素材库</span>
            </div>
          </div>
        </ElFormItem>
        <ElFormItem label="升级任务：" required>
          <ElTable
            :data="taskRows"
            border
            size="small"
            class="task-table"
          >
            <ElTableColumn label="升级任务" width="148" align="left">
              <template #default="{ row }">
                <span class="task-table__label">{{ row.label }}</span>
              </template>
            </ElTableColumn>
            <ElTableColumn label="任务名称" min-width="150" align="center">
              <template #default="{ row }">
                <ElInput
                  v-model="row.name"
                  :disabled="isReadonly"
                  placeholder="请输入任务名称"
                  size="small"
                />
              </template>
            </ElTableColumn>
            <ElTableColumn label="升级条件" width="132" align="center">
              <template #default="{ row }">
                <div class="task-condition">
                  <ElInputNumber
                    v-model="row.num"
                    :disabled="isReadonly"
                    :min="0"
                    :controls="false"
                    size="small"
                    class="task-condition__num"
                  />
                  <span class="task-condition__unit">{{ row.unit }}</span>
                </div>
              </template>
            </ElTableColumn>
            <ElTableColumn label="任务描述" min-width="180" align="center">
              <template #default="{ row }">
                <ElInput
                  v-model="row.info"
                  :disabled="isReadonly"
                  type="textarea"
                  :rows="2"
                  placeholder="任务描述"
                  size="small"
                  resize="none"
                />
              </template>
            </ElTableColumn>
          </ElTable>
        </ElFormItem>
        <ElFormItem label="一级返佣(上浮比例)：" required>
          <div class="rate-field">
            <div class="rate-field__input">
              <ElInputNumber
                v-model="form.extension_one"
                :disabled="isReadonly"
                :min="0"
                :max="1000"
                :precision="2"
                :controls="false"
                class="rate-field__num"
              />
              <span class="rate-field__unit">%</span>
            </div>
            <div class="rate-field__hint">{{ oneRateHint }}</div>
          </div>
        </ElFormItem>
        <ElFormItem label="二级返佣(上浮比例)：" required>
          <div class="rate-field">
            <div class="rate-field__input">
              <ElInputNumber
                v-model="form.extension_two"
                :disabled="isReadonly"
                :min="0"
                :max="1000"
                :precision="2"
                :controls="false"
                class="rate-field__num"
              />
              <span class="rate-field__unit">%</span>
            </div>
            <div class="rate-field__hint">{{ twoRateHint }}</div>
          </div>
        </ElFormItem>
      </ElForm>
    </LevelDrawer>
  </Page>
</template>

<style scoped>
.level-list-icon {
  display: block;
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: contain;
  border: 1px solid hsl(var(--border));
}

.level-form {
  padding: 4px 8px 12px;
}

.level-form__rank {
  width: 160px;
}

.icon-picker {
  display: grid;
  grid-template-columns: repeat(3, 48px);
  gap: 8px;
  align-items: start;
  width: max-content;
  max-width: 100%;
}

.icon-picker__item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  padding: 4px;
  cursor: pointer;
  background: #fff;
  border: 1px dashed transparent;
  border-radius: 6px;
  box-sizing: border-box;
}

.icon-picker__item img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.icon-picker__item.is-active,
.icon-picker__custom.is-active :deep(.image-field__tile) {
  border-color: #c0c4cc;
  border-style: dashed;
  background: #fafafa;
}

.icon-picker__item:disabled {
  cursor: default;
  opacity: 0.85;
}

.icon-picker__custom {
  display: flex;
  flex-direction: column;
  gap: 2px;
  align-items: center;
  width: 48px;
}

.icon-picker__custom :deep(.image-field) {
  gap: 0;
}

.icon-picker__custom :deep(.image-field__tile) {
  border-radius: 6px;
}

.icon-picker__custom-label {
  font-size: 11px;
  line-height: 1.2;
  color: #909399;
  white-space: nowrap;
}

.task-table {
  width: 100%;
}

.task-table :deep(.el-table__cell) {
  padding: 8px 6px;
  vertical-align: middle;
}

.task-table__label {
  display: inline-block;
  padding-left: 4px;
  color: #303133;
  line-height: 1.4;
  white-space: normal;
}

.task-condition {
  display: inline-flex;
  gap: 6px;
  align-items: center;
  justify-content: center;
}

.task-condition__num {
  width: 72px;
}

.task-condition__unit {
  min-width: 1em;
  color: #606266;
}

.rate-field {
  width: 100%;
}

.rate-field__input {
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.rate-field__num {
  width: 160px;
}

.rate-field__unit {
  color: #606266;
}

.rate-field__hint {
  margin-top: 8px;
  max-width: 720px;
  font-size: 12px;
  line-height: 1.55;
  color: #909399;
}
</style>
