<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, nextTick, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElTabPane,
  ElTabs,
  ElTree,
} from 'element-plus';
import { Grid as GridIcon, Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteMerchantType,
  fetchMerchantStoreMenus,
  fetchMerchantType,
  fetchMerchantTypes,
  saveMerchantType,
  type MerchantStoreMenuRow,
  type MerchantTypeRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit' | 'view';

type MenuTreeNode = {
  code: string;
  menu_id: number;
  menu_name: string;
  children?: MenuTreeNode[];
};

const canManage = ref(false);
const drawerMode = ref<DrawerMode>('create');
const activeTab = ref<'basic' | 'auth'>('basic');
const editing = ref<MerchantTypeRow>();
const menuTreeData = ref<MenuTreeNode[]>([]);
const detailAuthTree = ref<MenuTreeNode[]>([]);
const menuTreeRef = ref<InstanceType<typeof ElTree>>();
const detailTreeRef = ref<InstanceType<typeof ElTree>>();
const treeExpanded = ref(false);

const form = reactive({
  name: '',
  type_info: '',
  is_margin: false,
  margin: 0,
  description: '',
  remark: '',
  menu_codes: [] as string[],
});

const isReadonly = computed(() => drawerMode.value === 'view');
const headerName = computed(
  () => form.name.trim() || editing.value?.name || '未命名类型',
);
const marginSummary = computed(() =>
  form.is_margin ? `${Number(form.margin || 0).toFixed(2)}元` : '无',
);
const depositLabel = computed(() =>
  form.is_margin ? `${Number(form.margin || 0).toFixed(2)}元` : '无',
);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入店铺类型名称' },
    fieldName: 'keyword',
    label: '店铺类型名称',
  },
]);

const gridOptions: VxeGridProps<MerchantTypeRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 72 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      slots: { default: 'name' },
      title: '店铺类型名称',
    },
    {
      field: 'store_count',
      formatter: ({ cellValue }) => String(Number(cellValue || 0)),
      title: '店铺数量',
      width: 96,
    },
    {
      field: 'margin',
      formatter: ({ row }) => formatDeposit(row),
      title: '店铺保证金',
      width: 120,
    },
    {
      className: 'col--remark',
      field: 'type_info',
      formatter: ({ cellValue }) => {
        const text = String(cellValue ?? '').trim();
        return text || '—';
      },
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '店铺类型要求',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      showOverflow: false,
      title: '创建时间',
      width: 168,
    },
    platformListActionColumn({ width: 120 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim() || undefined;
        const list =
          (await fetchMerchantTypes({ keyword })).list || [];
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
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

const [TypeDrawer, typeDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (drawerMode.value === 'view') return;
    await save();
  },
  onOpenChange(isOpen) {
    if (!isOpen) {
      activeTab.value = 'basic';
      treeExpanded.value = false;
    }
  },
});

function formatDeposit(row: Pick<MerchantTypeRow, 'is_margin' | 'margin'>) {
  return row.is_margin ? `${Number(row.margin || 0).toFixed(2)}元` : '无';
}

function displayOrDash(value?: string | number | null) {
  if (value === 0) return '0';
  const text = String(value ?? '').trim();
  return text || '无';
}

function buildMenuTree(rows: MerchantStoreMenuRow[]): MenuTreeNode[] {
  const map = new Map<number, MenuTreeNode>();
  for (const row of rows) {
    if (!row.code) continue;
    map.set(row.menu_id, {
      code: row.code,
      menu_id: row.menu_id,
      menu_name: row.menu_name || row.code,
      children: [],
    });
  }
  const roots: MenuTreeNode[] = [];
  for (const row of rows) {
    const node = map.get(row.menu_id);
    if (!node) continue;
    const parent = row.pid ? map.get(row.pid) : undefined;
    if (parent) {
      parent.children = parent.children || [];
      parent.children.push(node);
    } else {
      roots.push(node);
    }
  }
  const prune = (nodes: MenuTreeNode[]) => {
    for (const node of nodes) {
      if (node.children?.length) prune(node.children);
      else delete node.children;
    }
  };
  prune(roots);
  return roots;
}

function filterAssignedTree(
  nodes: MenuTreeNode[],
  codes: Set<string>,
): MenuTreeNode[] {
  const result: MenuTreeNode[] = [];
  for (const node of nodes) {
    const children = node.children
      ? filterAssignedTree(node.children, codes)
      : [];
    if (codes.has(node.code) || children.length) {
      result.push({
        ...node,
        children: children.length ? children : undefined,
      });
    }
  }
  return result;
}

async function ensureMenus() {
  if (menuTreeData.value.length) return;
  const rows = await fetchMerchantStoreMenus();
  menuTreeData.value = buildMenuTree(rows);
}

function collectCheckedCodes() {
  const checked = (menuTreeRef.value?.getCheckedKeys(false) || []) as string[];
  const half = (menuTreeRef.value?.getHalfCheckedKeys() || []) as string[];
  return Array.from(new Set([...checked, ...half].filter(Boolean)));
}

function applyCheckedKeys(codes: string[]) {
  nextTick(() => {
    menuTreeRef.value?.setCheckedKeys(codes, false);
    setTreeExpand(false);
  });
}

function collectCodes(nodes: MenuTreeNode[]): string[] {
  const codes: string[] = [];
  const walk = (list: MenuTreeNode[]) => {
    for (const node of list) {
      codes.push(node.code);
      if (node.children?.length) walk(node.children);
    }
  };
  walk(nodes);
  return codes;
}

function setTreeExpand(expand: boolean) {
  treeExpanded.value = expand;
  const tree = isReadonly.value ? detailTreeRef.value : menuTreeRef.value;
  const data = isReadonly.value ? detailAuthTree.value : menuTreeData.value;
  nextTick(() => {
    for (const code of collectCodes(data)) {
      const node = tree?.getNode?.(code);
      if (node) node.expanded = expand;
    }
  });
}

function toggleTreeExpand() {
  setTreeExpand(!treeExpanded.value);
}

function resetForm(row?: MerchantTypeRow) {
  editing.value = row;
  Object.assign(form, {
    name: row?.name || '',
    type_info: row?.type_info || '',
    is_margin: row?.is_margin === 1,
    margin: Number(row?.margin || 0),
    description: row?.description || '',
    remark: row?.remark || '',
    menu_codes: [...(row?.menu_codes || [])],
  });
  detailAuthTree.value = filterAssignedTree(
    menuTreeData.value,
    new Set(form.menu_codes),
  );
}

async function openCreate() {
  drawerMode.value = 'create';
  activeTab.value = 'basic';
  await ensureMenus();
  resetForm();
  typeDrawerApi
    .setState({
      title: '添加店铺类型',
      showConfirmButton: true,
      confirmText: '提交',
      cancelText: '取消',
    })
    .open();
  applyCheckedKeys([]);
}

async function openEdit(row: MerchantTypeRow) {
  drawerMode.value = 'edit';
  activeTab.value = 'basic';
  await ensureMenus();
  typeDrawerApi
    .setState({
      loading: true,
      title: '编辑店铺类型',
      showConfirmButton: true,
      confirmText: '完成',
      cancelText: '取消',
    })
    .open();
  try {
    const detail = await fetchMerchantType(row.id);
    resetForm(detail);
    applyCheckedKeys(detail.menu_codes || []);
  } finally {
    typeDrawerApi.setState({ loading: false });
  }
}

async function openDetail(row: MerchantTypeRow) {
  drawerMode.value = 'view';
  activeTab.value = 'basic';
  await ensureMenus();
  typeDrawerApi
    .setState({
      loading: true,
      title: '店铺类型详情',
      showConfirmButton: false,
      cancelText: '关闭',
    })
    .open();
  try {
    const detail = await fetchMerchantType(row.id);
    resetForm(detail);
    nextTick(() => setTreeExpand(false));
  } finally {
    typeDrawerApi.setState({ loading: false });
  }
}

async function switchToEdit() {
  const typeId = editing.value?.id;
  if (!typeId || !canManage.value) return;
  drawerMode.value = 'edit';
  typeDrawerApi.setState({
    title: '编辑店铺类型',
    showConfirmButton: true,
    confirmText: '完成',
    cancelText: '取消',
  });
  // 详情→编辑：保留并回填同一 id，避免 checkbox 树重挂载后丢状态
  try {
    typeDrawerApi.setState({ loading: true });
    const detail = await fetchMerchantType(typeId);
    resetForm(detail);
    applyCheckedKeys(detail.menu_codes || []);
  } finally {
    typeDrawerApi.setState({ loading: false });
  }
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写店铺类型名称');
    activeTab.value = 'basic';
    return;
  }
  if (form.is_margin && Number(form.margin) <= 0) {
    ElMessage.warning('启用保证金时金额必须大于 0');
    activeTab.value = 'basic';
    return;
  }
  const menu_codes = collectCheckedCodes();
  typeDrawerApi.lock();
  try {
    await saveMerchantType(editing.value?.id, {
      name: form.name.trim(),
      type_info: form.type_info.trim(),
      is_margin: form.is_margin,
      margin: form.is_margin ? Number(form.margin) : 0,
      description: form.description.trim(),
      remark: form.remark.trim(),
      status: true,
      menu_codes,
    });
    typeDrawerApi.close();
    ElMessage.success('店铺类型已保存');
    gridApi.reload();
  } finally {
    typeDrawerApi.unlock();
  }
}

async function remove(row: MerchantTypeRow) {
  try {
    const countHint =
      Number(row.store_count || 0) > 0
        ? `当前有 ${row.store_count} 个店铺使用该类型，`
        : '';
    await confirm({
      content: `${countHint}删除“${row.name}”将移除其店铺菜单授权，是否继续？`,
      icon: 'warning',
      title: '删除店铺类型',
    });
    await deleteMerchantType(row.id);
    ElMessage.success('店铺类型已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value = codes.includes('merchant.type.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="types-toolbar">
          <ElAlert
            class="types-tip"
            type="warning"
            show-icon
            :closable="false"
            :title="
              '用于定义店铺经营模式，例如“旗舰店、自营店、工厂店”，不同的类型可设置不同保证金，店铺类型也会在移动端商家店铺的首页展示'
            "
          />
          <div class="types-toolbar__actions">
            <ElButton
              v-if="canManage"
              :icon="Plus"
              type="primary"
              @click="openCreate"
            >
              添加店铺类型
            </ElButton>
          </div>
        </div>
      </template>
      <template #name="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">
          {{ row.name || '—' }}
        </ElButton>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          type="danger"
          @click="remove(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <TypeDrawer>
      <div class="type-drawer">
        <div
          v-if="drawerMode !== 'create'"
          class="type-drawer__header"
        >
          <div class="type-drawer__brand">
            <div class="type-drawer__avatar">
              <GridIcon />
            </div>
            <div class="type-drawer__titles">
              <div class="type-drawer__name">{{ headerName }}</div>
              <div class="type-drawer__sub">保证金：{{ marginSummary }}</div>
            </div>
          </div>
          <div v-if="drawerMode === 'view' && canManage" class="type-drawer__actions">
            <ElButton type="primary" @click="switchToEdit">编辑</ElButton>
          </div>
        </div>

        <ElTabs v-model="activeTab" class="type-drawer__tabs">
          <ElTabPane label="基本信息" name="basic">
            <div class="type-section">
              <div class="type-section__title">店铺类型详情</div>

              <template v-if="isReadonly">
                <div class="type-desc-grid">
                  <div class="type-desc">
                    <span class="label">店铺类型名称</span>
                    <span class="value">{{ displayOrDash(form.name) }}</span>
                  </div>
                  <div class="type-desc">
                    <span class="label">店铺保证金</span>
                    <span class="value">{{ depositLabel }}</span>
                  </div>
                  <div class="type-desc">
                    <span class="label">店铺类型要求</span>
                    <span class="value">{{ displayOrDash(form.type_info) }}</span>
                  </div>
                  <div class="type-desc">
                    <span class="label">已有店铺数量</span>
                    <span class="value">
                      {{
                        editing?.store_count
                          ? editing.store_count
                          : '无'
                      }}
                    </span>
                  </div>
                  <div class="type-desc">
                    <span class="label">其他说明</span>
                    <span class="value">
                      {{ displayOrDash(form.description) }}
                    </span>
                  </div>
                  <div class="type-desc">
                    <span class="label">创建时间</span>
                    <span class="value">
                      {{
                        editing?.created_at
                          ? formatShanghaiDateTime(editing.created_at)
                          : '—'
                      }}
                    </span>
                  </div>
                  <div class="type-desc">
                    <span class="label">备注</span>
                    <span class="value">{{ displayOrDash(form.remark) }}</span>
                  </div>
                  <div class="type-desc">
                    <span class="label">最新修改时间</span>
                    <span class="value">
                      {{
                        editing?.updated_at
                          ? formatShanghaiDateTime(editing.updated_at)
                          : '—'
                      }}
                    </span>
                  </div>
                </div>
              </template>

              <ElForm v-else label-width="120px" class="type-form">
                <div class="type-form-grid">
                  <ElFormItem label="店铺类型名称" required>
                    <ElInput
                      v-model="form.name"
                      maxlength="128"
                      placeholder="请填写店铺类型名称"
                    />
                  </ElFormItem>
                  <ElFormItem label="店铺保证金">
                    <div class="flex flex-wrap items-center gap-3">
                      <ElRadioGroup v-model="form.is_margin">
                        <ElRadio :value="false">无</ElRadio>
                        <ElRadio :value="true">有</ElRadio>
                      </ElRadioGroup>
                      <template v-if="form.is_margin">
                        <ElInputNumber
                          v-model="form.margin"
                          :min="0.01"
                          :precision="2"
                          :step="1"
                          controls-position="right"
                        />
                        <span>元</span>
                      </template>
                    </div>
                  </ElFormItem>
                  <ElFormItem label="店铺类型要求">
                    <ElInput
                      v-model="form.type_info"
                      :rows="4"
                      maxlength="500"
                      placeholder="请填写店铺类型要求"
                      type="textarea"
                    />
                  </ElFormItem>
                  <ElFormItem label="其他说明">
                    <ElInput
                      v-model="form.description"
                      :rows="4"
                      maxlength="65535"
                      placeholder="请填写其他说明"
                      type="textarea"
                    />
                  </ElFormItem>
                  <ElFormItem label="备注">
                    <ElInput
                      v-model="form.remark"
                      :rows="3"
                      maxlength="500"
                      placeholder="请填写备注"
                      type="textarea"
                    />
                  </ElFormItem>
                </div>
              </ElForm>
            </div>
          </ElTabPane>

          <ElTabPane label="权限信息" name="auth">
            <p v-if="!isReadonly" class="type-auth-tip">
              如果提示权限不足，请先去掉任意一个权限提交，然后重新补充权限再次提交即可。
            </p>
            <div class="type-section">
              <div class="type-section__title">权限管理</div>
              <div class="type-auth-toolbar">
                <span class="label">
                  {{ isReadonly ? '权限范围：' : '权限：' }}
                </span>
                <ElButton link type="primary" @click="toggleTreeExpand">
                  展开/收起
                </ElButton>
              </div>
              <div class="type-auth-tree">
                <ElTree
                  v-if="isReadonly"
                  ref="detailTreeRef"
                  :data="detailAuthTree"
                  :default-expand-all="false"
                  :props="{ label: 'menu_name', children: 'children' }"
                  node-key="code"
                />
                <ElTree
                  v-else
                  ref="menuTreeRef"
                  :data="menuTreeData"
                  :default-expand-all="false"
                  :props="{ label: 'menu_name', children: 'children' }"
                  check-on-click-node
                  node-key="code"
                  show-checkbox
                />
              </div>
            </div>
          </ElTabPane>
        </ElTabs>
      </div>
    </TypeDrawer>
  </Page>
</template>

<style scoped>
.types-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.types-tip {
  width: 100%;
}

.types-toolbar__actions {
  display: flex;
  justify-content: flex-start;
}

.type-drawer {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 100%;
}

.type-drawer__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid hsl(var(--border));
}

.type-drawer__brand {
  display: flex;
  gap: 12px;
  align-items: center;
  min-width: 0;
}

.type-drawer__avatar {
  display: flex;
  width: 48px;
  height: 48px;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: hsl(var(--primary) / 12%);
  color: hsl(var(--primary));
  font-size: 22px;
}

.type-drawer__name {
  font-size: 18px;
  font-weight: 600;
  line-height: 28px;
}

.type-drawer__sub {
  margin-top: 2px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.type-drawer__actions {
  flex-shrink: 0;
}

.type-drawer__tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}

.type-section {
  margin-bottom: 8px;
}

.type-section__title {
  position: relative;
  margin-bottom: 14px;
  padding-left: 10px;
  font-size: 15px;
  font-weight: 600;
  line-height: 24px;
}

.type-section__title::before {
  position: absolute;
  top: 4px;
  left: 0;
  width: 3px;
  height: 16px;
  border-radius: 2px;
  background: hsl(var(--primary));
  content: '';
}

.type-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 20px;
}

.type-desc-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px 32px;
}

.type-desc {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.type-desc .label {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 18px;
}

.type-desc .value {
  font-size: 14px;
  line-height: 22px;
  word-break: break-word;
}

.type-auth-tip {
  margin: 0 0 12px;
  color: hsl(var(--destructive));
  font-size: 13px;
  line-height: 20px;
}

.type-auth-toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.type-auth-toolbar .label {
  color: hsl(var(--foreground));
  font-size: 14px;
}

.type-auth-tree {
  min-height: 240px;
  max-height: min(56vh, 520px);
  overflow: auto;
  padding: 8px 4px;
}

@media (max-width: 768px) {
  .type-form-grid,
  .type-desc-grid {
    grid-template-columns: 1fr;
  }
}
</style>
