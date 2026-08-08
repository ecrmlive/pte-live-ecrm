<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElCascader,
  ElCheckbox,
  ElCheckboxGroup,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createBusinessZone,
  createBusinessZoneAgent,
  deleteBusinessZone,
  fetchBusinessZone,
  fetchBusinessZoneAgents,
  fetchBusinessZones,
  fetchMerchantCategories,
  fetchMerchantTypes,
  fetchPlatformMerchants,
  fetchPlatformRoles,
  updateBusinessZone,
  updateBusinessZoneStatus,
  type BusinessZoneMerchantBrief,
  type BusinessZoneRow,
  type BusinessZoneSaveInput,
  type MerchantCategoryRow,
  type MerchantTypeRow,
  type PlatformMerchantRow,
  type PlatformRoleRow,
} from '#/api/core/ecrm';
import {
  listPlatformCategoriesApi,
  type PlatformCategory,
} from '#/api/core/platform-catalog';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit' | 'view';

const GOODS_TYPE_OPTIONS = [
  { label: '普通商品', value: 0 },
  { label: '虚拟', value: 1 },
  { label: '云盘', value: 2 },
  { label: '卡密', value: 3 },
  { label: '预约', value: 4 },
  { label: '年/次卡', value: 5 },
] as const;

type PlatformCategoryCascaderOption = {
  label: string;
  value: number;
  children?: PlatformCategoryCascaderOption[];
};

const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const activeTab = ref('basic');
const categories = ref<MerchantCategoryRow[]>([]);
const types = ref<MerchantTypeRow[]>([]);
const roles = ref<PlatformRoleRow[]>([]);
const agentOptions = ref<{ label: string; value: number; phone?: string }[]>(
  [],
);
const platformCategoryTree = ref<PlatformCategoryCascaderOption[]>([]);
const linkedMerchants = ref<BusinessZoneMerchantBrief[]>([]);

const form = reactive({
  name: '',
  business_store_category: undefined as number | undefined,
  business_store_type: undefined as number | undefined,
  platform_category_ids: [] as number[],
  goods_type: [] as number[],
  circle_agent_id: undefined as number | undefined,
  role_id: undefined as number | undefined,
  sort: 0,
  status: 1,
});

const adminForm = reactive({
  name: '',
  phone: '',
  account: '',
  password: '000000',
  uid: 0,
});

const isReadonly = computed(() => drawerMode.value === 'view');
const drawerTitle = computed(() => {
  if (drawerMode.value === 'view') return '商户详情';
  if (drawerMode.value === 'edit') return '编辑商户';
  return '新增商户';
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入商户名称' },
    fieldName: 'name',
    label: '商户名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '全部',
    },
    fieldName: 'circle_agent_id',
    label: '商户管理员',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '全部', value: '' },
        { label: '开启', value: 1 },
        { label: '关闭', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '商户状态',
  },
]);

const gridOptions: VxeGridProps<BusinessZoneRow> = {
  columns: [
    { field: 'circle_id', title: '商户ID', width: 88 },
    { field: 'name', minWidth: 140, title: '商户名称' },
    {
      field: 'circle_agent',
      formatter: ({ row }) => row.circle_agent?.name || '—',
      minWidth: 110,
      title: '商户管理员',
    },
    {
      field: 'phone',
      formatter: ({ row }) => row.circle_agent?.phone || '—',
      minWidth: 120,
      title: '手机号码',
    },
    {
      field: 'merchant_count',
      formatter: ({ cellValue }) => String(cellValue ?? 0),
      title: '关联店铺',
      width: 96,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '商户状态',
      width: 110,
    },
    { field: 'sort', title: '排序', width: 72 },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const statusRaw = formValues?.status;
        const result = await fetchBusinessZones({
          page: page.currentPage,
          limit: page.pageSize,
          type: 1,
          name: String(formValues?.name ?? '').trim() || undefined,
          circle_agent_id: formValues?.circle_agent_id
            ? Number(formValues.circle_agent_id)
            : undefined,
          status:
            statusRaw === 0 || statusRaw === 1 || statusRaw === '0' || statusRaw === '1'
              ? Number(statusRaw)
              : undefined,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '提交',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (isReadonly.value) {
      formDrawerApi.close();
      return;
    }
    await save();
  },
});

const [AdminDrawer, adminDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => saveAdmin(),
});

/* ---------- 选择店铺 ---------- */
const pickerFilter = reactive({ keyword: '' });
const pickerLoading = ref(false);
const pickerRows = ref<PlatformMerchantRow[]>([]);
const pickerTotal = ref(0);
const pickerPage = ref(1);
const pickerLimit = ref(10);
const pickerSelected = ref<BusinessZoneMerchantBrief[]>([]);
const pickerTableRef = ref<{
  clearSelection: () => void;
  toggleRowSelection: (row: PlatformMerchantRow, selected?: boolean) => void;
}>();

const [StorePickerModal, storePickerApi] = useVbenModal({
  title: '选择店铺',
  class: 'w-[920px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  onConfirm: () => {
    const map = new Map(
      linkedMerchants.value.map((item) => [item.mer_id, item]),
    );
    for (const item of pickerSelected.value) {
      map.set(item.mer_id, item);
    }
    linkedMerchants.value = [...map.values()];
    storePickerApi.close();
  },
});

function toPlatformCategoryCascaderOptions(
  rows: PlatformCategory[] = [],
): PlatformCategoryCascaderOption[] {
  const out: PlatformCategoryCascaderOption[] = [];
  for (const row of rows || []) {
    const value = Number(row.store_category_id);
    const label = String(row.cate_name || '').trim();
    if (!Number.isFinite(value) || value <= 0 || !label) continue;
    const children = toPlatformCategoryCascaderOptions(row.children || []);
    if (Number(row.is_show) === 0 && !children.length) continue;
    const option: PlatformCategoryCascaderOption = { label, value };
    if (children.length) option.children = children;
    if (Number(row.is_show) !== 0) out.push(option);
    else out.push(...children);
  }
  return out;
}

async function loadOptions() {
  const [cateRes, typeRes, roleRes, agentRes, categoryRes] = await Promise.all([
    fetchMerchantCategories(),
    fetchMerchantTypes(),
    fetchPlatformRoles({ page: 1, limit: 100 }),
    fetchBusinessZoneAgents({
      page: 1,
      limit: 100,
      type: 1,
      status: 1,
    }),
    listPlatformCategoriesApi(),
  ]);
  categories.value = cateRes.list || [];
  types.value = typeRes.list || [];
  roles.value = (roleRes.list || []).filter((row) => row.status === 1);
  agentOptions.value = (agentRes.list || []).map((row) => ({
    label: row.name,
    value: row.circle_agent_id,
    phone: row.phone,
  }));
  const categoryList = Array.isArray(categoryRes?.list)
    ? categoryRes.list
    : Array.isArray(categoryRes)
      ? (categoryRes as PlatformCategory[])
      : [];
  platformCategoryTree.value = toPlatformCategoryCascaderOptions(categoryList);
  gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'circle_agent_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: agentOptions.value,
        placeholder: '全部',
      },
    },
  ]);
}

function resetForm() {
  Object.assign(form, {
    name: '',
    business_store_category: undefined,
    business_store_type: undefined,
    platform_category_ids: [],
    goods_type: [],
    circle_agent_id: undefined,
    role_id: undefined,
    sort: 0,
    status: 1,
  });
  linkedMerchants.value = [];
  activeTab.value = 'basic';
}

function openCreate() {
  editingId.value = 0;
  drawerMode.value = 'create';
  resetForm();
  formDrawerApi.setState({ title: drawerTitle.value, showConfirmButton: true }).open();
}

async function openEdit(row: BusinessZoneRow, mode: DrawerMode = 'edit') {
  editingId.value = row.circle_id;
  drawerMode.value = mode;
  resetForm();
  const detail = await fetchBusinessZone(row.circle_id);
  Object.assign(form, {
    name: detail.name || '',
    business_store_category: detail.business_store_category || undefined,
    business_store_type: detail.business_store_type || undefined,
    platform_category_ids: [...(detail.platform_category_ids || [])],
    goods_type: [...(detail.goods_type || [])],
    circle_agent_id: detail.circle_agent_id || undefined,
    role_id: detail.role_id || undefined,
    sort: detail.sort || 0,
    status: detail.status ?? 1,
  });
  linkedMerchants.value = [...(detail.merchant || [])];
  if (
    detail.circle_agent_id &&
    !agentOptions.value.some((item) => item.value === detail.circle_agent_id)
  ) {
    agentOptions.value = [
      {
        label: detail.circle_agent?.name || `管理员#${detail.circle_agent_id}`,
        value: detail.circle_agent_id,
        phone: detail.circle_agent?.phone,
      },
      ...agentOptions.value,
    ];
  }
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: mode !== 'view',
    })
    .open();
}

function openView(row: BusinessZoneRow) {
  void openEdit(row, 'view');
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入商户名称');
    activeTab.value = 'basic';
    return;
  }
  if (!form.business_store_category) {
    ElMessage.warning('请选择店铺分类');
    activeTab.value = 'basic';
    return;
  }
  if (!form.business_store_type) {
    ElMessage.warning('请选择店铺类型');
    activeTab.value = 'basic';
    return;
  }
  if (!form.circle_agent_id) {
    ElMessage.warning('请选择商户管理员');
    activeTab.value = 'basic';
    return;
  }
  if (!form.role_id) {
    ElMessage.warning('请选择身份权限');
    activeTab.value = 'basic';
    return;
  }
  const payload: BusinessZoneSaveInput = {
    pid: 0,
    type: 1,
    name: form.name.trim(),
    business_store_category: form.business_store_category,
    business_store_type: form.business_store_type,
    platform_category_ids: [...form.platform_category_ids],
    goods_type: [...form.goods_type],
    circle_agent_id: form.circle_agent_id,
    role_id: form.role_id,
    sort: form.sort || 0,
    status: form.status,
    merchant_ids: linkedMerchants.value.map((item) => item.mer_id),
    commission_type: 0,
    commission_rate: 0,
  };
  formDrawerApi.lock();
  try {
    if (editingId.value) {
      await updateBusinessZone(editingId.value, payload);
    } else {
      await createBusinessZone(payload);
    }
    ElMessage.success('保存成功');
    formDrawerApi.close();
    gridApi.reload();
    await loadOptions();
  } finally {
    formDrawerApi.unlock();
  }
}

async function toggleStatus(row: BusinessZoneRow, value: number | string | boolean) {
  const next = Number(value);
  const prev = row.status;
  row.status = next;
  try {
    await updateBusinessZoneStatus(row.circle_id, next);
    ElMessage.success(next === 1 ? '已开启' : '已关闭');
  } catch {
    row.status = prev;
  }
}

async function remove(row: BusinessZoneRow) {
  try {
    await ElMessageBox.confirm(`确定删除商户「${row.name}」吗？`, '提示', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消',
    });
    await deleteBusinessZone(row.circle_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

function openAdminDrawer() {
  Object.assign(adminForm, {
    name: '',
    phone: '',
    account: '',
    password: '000000',
    uid: 0,
  });
  adminDrawerApi.setState({ title: '添加商户管理员' }).open();
}

watch(
  () => adminForm.phone,
  (phone, prev) => {
    if (!adminForm.account || adminForm.account === prev) {
      adminForm.account = phone;
    }
  },
);

async function saveAdmin() {
  if (!adminForm.name.trim()) {
    ElMessage.warning('请输入管理姓名');
    return;
  }
  if (!adminForm.phone.trim()) {
    ElMessage.warning('请输入手机号码');
    return;
  }
  if (!adminForm.account.trim()) {
    ElMessage.warning('请输入登录账号');
    return;
  }
  if (!adminForm.password.trim()) {
    ElMessage.warning('请输入登录密码');
    return;
  }
  adminDrawerApi.lock();
  try {
    const created = await createBusinessZoneAgent({
      type: 1,
      name: adminForm.name.trim(),
      phone: adminForm.phone.trim(),
      account: adminForm.account.trim(),
      password: adminForm.password,
      uid: adminForm.uid || 0,
      qualification: '',
      remark: '',
      payment_method: 0,
      payment_name: '',
      business_name: '',
    });
    await loadOptions();
    form.circle_agent_id = created.circle_agent_id;
    ElMessage.success('管理员已添加');
    adminDrawerApi.close();
  } finally {
    adminDrawerApi.unlock();
  }
}

function removeLinked(merId: number) {
  linkedMerchants.value = linkedMerchants.value.filter(
    (item) => item.mer_id !== merId,
  );
}

async function loadPicker() {
  pickerLoading.value = true;
  try {
    const result = await fetchPlatformMerchants({
      page: pickerPage.value,
      limit: pickerLimit.value,
      keyword: pickerFilter.keyword.trim() || undefined,
      status: 1,
    });
    pickerRows.value = result.list || [];
    pickerTotal.value = result.total || 0;
    await nextTick();
    const selectedIds = new Set(pickerSelected.value.map((item) => item.mer_id));
    for (const row of pickerRows.value) {
      if (selectedIds.has(row.mer_id)) {
        pickerTableRef.value?.toggleRowSelection(row, true);
      }
    }
  } finally {
    pickerLoading.value = false;
  }
}

function openStorePicker() {
  pickerFilter.keyword = '';
  pickerPage.value = 1;
  pickerSelected.value = [...linkedMerchants.value];
  storePickerApi.open();
  void loadPicker();
}

function onPickerSelectionChange(rows: PlatformMerchantRow[]) {
  pickerSelected.value = rows.map((row) => ({
    mer_id: row.mer_id,
    mer_name: row.mer_name,
    real_name: row.real_name,
    mer_phone: row.mer_phone,
  }));
}

onMounted(() => {
  void loadOptions();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增商户
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status"
          :active-value="1"
          :inactive-value="0"
          active-text="开启"
          inactive-text="关闭"
          inline-prompt
          @change="(val) => toggleStatus(row, val)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openView(row)">详情</ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElTabs v-model="activeTab">
        <ElTabPane label="基础信息" name="basic">
          <ElForm label-width="110px" class="pt-2">
            <ElFormItem label="商户名称" required>
              <ElInput
                v-model="form.name"
                :disabled="isReadonly"
                maxlength="20"
                show-word-limit
                placeholder="请输入商户名称"
                class="max-w-[460px]"
              />
            </ElFormItem>
            <ElFormItem label="店铺分类" required>
              <ElSelect
                v-model="form.business_store_category"
                :disabled="isReadonly"
                clearable
                class="max-w-[460px] w-full"
                placeholder="请选择店铺分类"
              >
                <ElOption
                  v-for="item in categories"
                  :key="item.merchant_category_id"
                  :label="item.category_name"
                  :value="item.merchant_category_id"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="店铺类型" required>
              <ElSelect
                v-model="form.business_store_type"
                :disabled="isReadonly"
                clearable
                class="max-w-[460px] w-full"
                placeholder="请选择店铺类型"
              >
                <ElOption
                  v-for="item in types"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="平台分类" required>
              <ElCascader
                v-model="form.platform_category_ids"
                :disabled="isReadonly"
                :options="platformCategoryTree"
                :props="{
                  multiple: true,
                  checkStrictly: false,
                  emitPath: false,
                  value: 'value',
                  label: 'label',
                  children: 'children',
                }"
                clearable
                filterable
                collapse-tags
                collapse-tags-tooltip
                class="max-w-[460px] w-full"
                placeholder="请选择平台分类"
              />
              <div class="field-tip">
                注：平台分类为空时，该商户下的店铺可以使用全部分类
              </div>
            </ElFormItem>
            <ElFormItem label="商品类型" required>
              <ElCheckboxGroup v-model="form.goods_type" :disabled="isReadonly">
                <ElCheckbox
                  v-for="item in GOODS_TYPE_OPTIONS"
                  :key="item.value"
                  :label="item.value"
                >
                  {{ item.label }}
                </ElCheckbox>
              </ElCheckboxGroup>
              <div class="field-tip">
                注：商品类型为空时，该商户下的店铺可以使用全部类型
              </div>
            </ElFormItem>
            <ElFormItem label="商户管理员" required>
              <div class="admin-row">
                <ElButton
                  v-if="!isReadonly"
                  type="primary"
                  plain
                  @click="openAdminDrawer"
                >
                  + 添加管理员
                </ElButton>
                <ElSelect
                  v-model="form.circle_agent_id"
                  :disabled="isReadonly"
                  filterable
                  clearable
                  class="admin-select"
                  placeholder="请输入选择商户管理员"
                >
                  <ElOption
                    v-for="item in agentOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </ElSelect>
              </div>
            </ElFormItem>
            <ElFormItem label="身份权限" required>
              <ElSelect
                v-model="form.role_id"
                :disabled="isReadonly"
                filterable
                clearable
                class="max-w-[460px] w-full"
                placeholder="请输入选择身份权限"
              >
                <ElOption
                  v-for="item in roles"
                  :key="item.role_id"
                  :label="item.role_name"
                  :value="item.role_id"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="排序">
              <ElInputNumber
                v-model="form.sort"
                :disabled="isReadonly"
                :min="0"
                :max="9999"
                :precision="0"
              />
              <div class="field-tip">数字越大越靠前</div>
            </ElFormItem>
            <ElFormItem label="开启状态">
              <ElSwitch
                v-model="form.status"
                :disabled="isReadonly"
                :active-value="1"
                :inactive-value="0"
                active-text="开启"
                inactive-text="关闭"
              />
              <div class="field-tip">关闭后，该商户禁止登录</div>
            </ElFormItem>
          </ElForm>
        </ElTabPane>
        <ElTabPane label="关联店铺" name="stores">
          <div class="stores-pane">
            <ElButton
              v-if="!isReadonly"
              type="primary"
              @click="openStorePicker"
            >
              选择店铺
            </ElButton>
            <ElTable :data="linkedMerchants" class="mt-3" border>
              <ElTableColumn label="店铺ID" prop="mer_id" width="90" />
              <ElTableColumn label="店铺名称" prop="mer_name" min-width="160" />
              <ElTableColumn label="联系人" prop="real_name" width="120" />
              <ElTableColumn label="联系电话" prop="mer_phone" width="140" />
              <ElTableColumn
                v-if="!isReadonly"
                label="操作"
                width="90"
                fixed="right"
              >
                <template #default="{ row }">
                  <ElButton link type="danger" @click="removeLinked(row.mer_id)">
                    删除
                  </ElButton>
                </template>
              </ElTableColumn>
            </ElTable>
          </div>
        </ElTabPane>
      </ElTabs>
    </FormDrawer>

    <AdminDrawer>
      <ElForm label-width="110px" class="pt-2">
        <ElFormItem label="管理姓名" required>
          <ElInput v-model="adminForm.name" placeholder="请输入管理员姓名" />
        </ElFormItem>
        <ElFormItem label="手机号码" required>
          <ElInput v-model="adminForm.phone" placeholder="请输入手机号码" />
          <div class="field-tip">
            手机号码为商户管理的登录账号，登录密码默认000000
          </div>
        </ElFormItem>
        <ElFormItem label="登录账号" required>
          <ElInput v-model="adminForm.account" placeholder="请输入登录账号" />
        </ElFormItem>
        <ElFormItem label="登录密码" required>
          <ElInput
            v-model="adminForm.password"
            type="password"
            show-password
            placeholder="请输入登录密码"
          />
        </ElFormItem>
        <ElFormItem label="关联用户">
          <ElInputNumber v-model="adminForm.uid" :min="0" />
          <div class="field-tip">可选；填写 C 端用户 UID，0 表示不关联</div>
        </ElFormItem>
      </ElForm>
    </AdminDrawer>

    <StorePickerModal>
      <div class="picker-filter">
        <ElForm inline @submit.prevent>
          <ElFormItem label="关键字">
            <ElInput
              v-model="pickerFilter.keyword"
              clearable
              class="picker-field"
              placeholder="店铺名称 / 联系人"
            />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="loadPicker">搜索</ElButton>
            <ElButton
              @click="
                () => {
                  pickerFilter.keyword = '';
                  loadPicker();
                }
              "
            >
              重置
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>
      <ElTable
        ref="pickerTableRef"
        v-loading="pickerLoading"
        :data="pickerRows"
        row-key="mer_id"
        border
        @selection-change="onPickerSelectionChange"
      >
        <ElTableColumn type="selection" width="48" />
        <ElTableColumn label="店铺ID" prop="mer_id" width="90" />
        <ElTableColumn label="店铺名称" prop="mer_name" min-width="160" />
        <ElTableColumn label="联系人" prop="real_name" width="120" />
        <ElTableColumn label="联系电话" prop="mer_phone" width="140" />
      </ElTable>
      <div class="picker-total">共 {{ pickerTotal }} 条</div>
    </StorePickerModal>
  </Page>
</template>

<style scoped>
.field-tip {
  width: 100%;
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.4;
  color: #909399;
}

.admin-row {
  display: flex;
  gap: 10px;
  align-items: center;
  width: 100%;
  max-width: 560px;
}

.admin-select {
  flex: 1;
  min-width: 0;
}

.stores-pane {
  padding-top: 4px;
}

.picker-filter {
  margin-bottom: 12px;
}

.picker-field {
  width: 220px;
}

.picker-total {
  margin-top: 10px;
  font-size: 13px;
  color: #909399;
}
</style>
