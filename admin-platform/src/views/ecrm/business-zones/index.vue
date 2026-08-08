<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import QRCode from 'qrcode';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createBusinessZone,
  createBusinessZoneRegionAgent,
  deleteBusinessZone,
  fetchBusinessZone,
  fetchBusinessZoneAgentOptions,
  fetchBusinessZoneInvite,
  fetchBusinessZoneOptions,
  fetchBusinessZones,
  fetchPlatformRoles,
  updateBusinessZone,
  updateBusinessZoneStatus,
  type BusinessZoneAgentOption,
  type BusinessZoneMerchantBrief,
  type BusinessZoneOptionNode,
  type BusinessZoneRow,
  type PlatformRoleRow,
} from '#/api/core/ecrm';
import ImageField from '#/components/shop/image-field.vue';
import StorePickerModal, {
  type PickedStore,
} from '#/components/ecrm/store-picker-modal.vue';
import UserPickerModal, {
  type PickedPlatformUser,
} from '#/components/ecrm/user-picker-modal.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit';
type ZoneTreeRow = BusinessZoneRow & {
  hasChild?: boolean;
  children?: ZoneTreeRow[];
};

/** 区域树一次拉全量同级节点（后端 limit 上限 9999） */
const ZONE_LIST_LIMIT = 9999;

const drawerMode = ref<DrawerMode>('create');
const editingID = ref<number>();
const activeTab = ref('basic');
const parentLocked = ref(false);
const agentOptions = ref<BusinessZoneAgentOption[]>([]);
const roleOptions = ref<PlatformRoleRow[]>([]);
const parentOptions = ref<Array<{ label: string; value: number }>>([
  { label: '顶级区域', value: 0 },
]);
const relatedStores = ref<PickedStore[]>([]);
const storePickerOpen = ref(false);
const userPickerOpen = ref(false);
const pickedUser = ref<PickedPlatformUser | null>(null);
const inviteImageUrl = ref('');
const inviteH5Url = ref('');
const inviteZoneName = ref('');
const statusUpdating = ref<Record<number, boolean>>({});

const form = reactive({
  pid: 0,
  name: '',
  circle_agent_id: undefined as number | undefined,
  role_id: undefined as number | undefined,
  commission_type: 0,
  commission_rate: 0,
  sort: 0,
  status: 1,
  remark: '',
});

const agentForm = reactive({
  name: '',
  phone: '',
  qualification: '',
  uid: 0,
  remark: '',
});

const isReadonly = computed(() => false);
const formDrawerTitle = computed(() => {
  if (drawerMode.value === 'edit') return '编辑区域';
  return parentLocked.value && form.pid > 0 ? '新增下级区域' : '新增区域';
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入区域名称' },
    fieldName: 'name',
    label: '区域名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '请选择区域代理',
    },
    fieldName: 'circle_agent_id',
    label: '区域代理',
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
    label: '开启状态',
  },
]);

function flattenZoneOptions(
  nodes: BusinessZoneOptionNode[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.label}`, value: node.value },
    ...flattenZoneOptions(node.children || [], `${prefix}— `),
  ]);
}

function toTreeRow(row: BusinessZoneRow): ZoneTreeRow {
  return {
    ...row,
    hasChild: Boolean(row.has_child) && Number(row.level) < 2,
  };
}

async function fetchZoneChildren(pid: number) {
  const result = await fetchBusinessZones({
    page: 1,
    limit: ZONE_LIST_LIMIT,
    pid,
    type: 0,
  });
  return (result.list || []).map(toTreeRow);
}

function commissionText(row: BusinessZoneRow) {
  if (row.commission_type === 1) {
    return `${Number(row.commission_rate || 0).toFixed(2)}%`;
  }
  return '默认设置';
}

function agentName(row: BusinessZoneRow) {
  return row.circle_agent?.name || '—';
}

function agentPhone(row: BusinessZoneRow) {
  return row.circle_agent?.phone || '—';
}

const gridOptions: VxeGridProps<ZoneTreeRow> = {
  columns: [
    {
      field: 'name',
      minWidth: 180,
      showOverflow: false,
      title: '区域名称',
      treeNode: true,
    },
    {
      field: 'circle_agent',
      formatter: ({ row }) => agentName(row),
      minWidth: 120,
      title: '区域代理',
    },
    {
      field: 'phone',
      formatter: ({ row }) => agentPhone(row),
      minWidth: 130,
      title: '手机号码',
    },
    {
      field: 'commission_rate',
      formatter: ({ row }) => commissionText(row),
      minWidth: 110,
      title: '代理提成',
    },
    {
      field: 'merchant_count',
      formatter: ({ cellValue }) => String(cellValue ?? 0),
      minWidth: 100,
      title: '关联店铺',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '开启状态',
      width: 100,
    },
    { field: 'sort', title: '排序', width: 80 },
    platformListActionColumn({ width: 300 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const statusRaw = formValues?.status;
        const agentRaw = formValues?.circle_agent_id;
        const name = String(formValues?.name ?? '').trim() || undefined;
        const hasFilter =
          Boolean(name) ||
          statusRaw === 0 ||
          statusRaw === 1 ||
          Boolean(agentRaw);
        // 无筛选：只拉一级区域；有筛选：按条件全量匹配（不限 pid）。
        const result = await fetchBusinessZones({
          page: 1,
          limit: ZONE_LIST_LIMIT,
          name,
          type: 0,
          pid: hasFilter ? undefined : 0,
          circle_agent_id:
            agentRaw === 0 || agentRaw
              ? Number(agentRaw) || undefined
              : undefined,
          status:
            statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
        });
        const items = (result.list || []).map(toTreeRow);
        return { items, total: items.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_id' },
  treeConfig: {
    childrenField: 'children',
    expandAll: false,
    hasChildField: 'hasChild',
    lazy: true,
    loadMethod: async ({ row }) => fetchZoneChildren(row.circle_id),
    transform: false,
  },
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
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (isReadonly.value) {
      formDrawerApi.close();
      return;
    }
    await saveZone();
  },
});

const [AgentDrawer, agentDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  title: '新增区域代理',
  onConfirm: async () => saveAgent(),
});

const [InviteModal, inviteModalApi] = useVbenModal({
  title: '邀请入驻',
  class: 'w-[480px]',
  footer: false,
});

async function loadAgentOptions(keyword = '') {
  const result = await fetchBusinessZoneAgentOptions(0);
  let list = result.list || [];
  const kw = keyword.trim().toLowerCase();
  if (kw) {
    list = list.filter(
      (item) =>
        item.name.toLowerCase().includes(kw) ||
        String(item.phone || '').includes(kw),
    );
  }
  agentOptions.value = list;
  const options = list.map((item) => ({
    label: item.phone ? `${item.name}（${item.phone}）` : item.name,
    value: item.circle_agent_id,
  }));
  gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'circle_agent_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options,
        placeholder: '请选择区域代理',
      },
    },
  ]);
}

async function loadRoleOptions(keyword = '') {
  const result = await fetchPlatformRoles({ page: 1, limit: 100 });
  let list = (result.list || []).filter(
    (item) => item.status === 1 && (item.is_agent === 1 || item.code.includes('region')),
  );
  if (!list.length) {
    list = (result.list || []).filter((item) => item.status === 1);
  }
  const kw = keyword.trim().toLowerCase();
  if (kw) {
    list = list.filter((item) => item.role_name.toLowerCase().includes(kw));
  }
  roleOptions.value = list;
}

async function loadParentOptions() {
  const result = await fetchBusinessZoneOptions(0);
  parentOptions.value = [
    { label: '顶级区域', value: 0 },
    ...flattenZoneOptions(result.list || []),
  ];
}

function resetForm(pid = 0) {
  Object.assign(form, {
    pid,
    name: '',
    circle_agent_id: undefined,
    role_id: undefined,
    commission_type: 0,
    commission_rate: 0,
    sort: 0,
    status: 1,
    remark: '',
  });
  relatedStores.value = [];
  activeTab.value = 'basic';
}

function resetAgentForm() {
  Object.assign(agentForm, {
    name: '',
    phone: '',
    qualification: '',
    uid: 0,
    remark: '',
  });
  pickedUser.value = null;
}

async function openCreate(pid = 0) {
  drawerMode.value = 'create';
  editingID.value = undefined;
  parentLocked.value = pid > 0;
  resetForm(pid);
  await Promise.all([loadAgentOptions(), loadRoleOptions(), loadParentOptions()]);
  formDrawerApi
    .setState({
      title: formDrawerTitle.value,
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
}

async function openEdit(row: BusinessZoneRow) {
  drawerMode.value = 'edit';
  editingID.value = row.circle_id;
  parentLocked.value = true;
  await Promise.all([loadAgentOptions(), loadRoleOptions(), loadParentOptions()]);
  const detail = await fetchBusinessZone(row.circle_id);
  Object.assign(form, {
    pid: detail.pid,
    name: detail.name,
    circle_agent_id: detail.circle_agent_id || undefined,
    role_id: detail.role_id || undefined,
    commission_type: detail.commission_type,
    commission_rate: detail.commission_rate,
    sort: detail.sort,
    status: detail.status,
    remark: detail.remark || '',
  });
  relatedStores.value = (detail.merchant || []).map((item: BusinessZoneMerchantBrief) => ({
    mer_id: item.mer_id,
    mer_name: item.mer_name,
    real_name: item.real_name || '',
    mer_phone: item.mer_phone || '',
  }));
  activeTab.value = 'basic';
  formDrawerApi
    .setState({
      title: '编辑区域',
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
}

async function saveZone() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写区域名称');
    activeTab.value = 'basic';
    return;
  }
  if (!form.circle_agent_id) {
    ElMessage.warning('请选择区域代理');
    activeTab.value = 'basic';
    return;
  }
  if (!form.role_id) {
    ElMessage.warning('请选择身份权限');
    activeTab.value = 'basic';
    return;
  }
  if (form.commission_type === 1 && (form.commission_rate < 0 || form.commission_rate > 100)) {
    ElMessage.warning('提成比例需在 0～100 之间');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      pid: form.pid || 0,
      name: form.name.trim(),
      circle_agent_id: Number(form.circle_agent_id),
      role_id: Number(form.role_id),
      commission_type: form.commission_type,
      commission_rate: form.commission_type === 1 ? Number(form.commission_rate || 0) : 0,
      sort: form.sort || 0,
      status: form.status,
      remark: form.remark,
      type: 0,
      merchant_ids: relatedStores.value.map((item) => item.mer_id),
      business_store_category: 0,
      business_store_type: 0,
    };
    if (editingID.value) {
      await updateBusinessZone(editingID.value, payload);
    } else {
      await createBusinessZone(payload);
    }
    formDrawerApi.close();
    ElMessage.success('保存成功');
    gridApi.reload();
    void loadParentOptions();
  } finally {
    formDrawerApi.unlock();
  }
}

function openAgentCreate() {
  if (isReadonly.value) return;
  resetAgentForm();
  agentDrawerApi.open();
}

async function saveAgent() {
  if (!agentForm.name.trim() || !agentForm.phone.trim()) {
    ElMessage.warning('代理名称和联系电话必填');
    return;
  }
  if (!agentForm.uid || agentForm.uid <= 0) {
    ElMessage.warning('请选择区域代理关联用户');
    return;
  }
  agentDrawerApi.lock();
  try {
    const created = await createBusinessZoneRegionAgent({
      name: agentForm.name.trim(),
      phone: agentForm.phone.trim(),
      qualification: agentForm.qualification,
      uid: agentForm.uid,
      remark: agentForm.remark,
      type: 0,
      auto_approve: true,
    });
    await loadAgentOptions();
    form.circle_agent_id = created.circle_agent_id;
    agentDrawerApi.close();
    ElMessage.success('代理人已新增');
  } finally {
    agentDrawerApi.unlock();
  }
}

function onUserPicked(user: PickedPlatformUser) {
  pickedUser.value = user;
  agentForm.uid = user.id;
  if (!agentForm.phone && user.mobile) {
    agentForm.phone = user.mobile;
  }
}

function onStoresPicked(stores: PickedStore[]) {
  relatedStores.value = stores;
}

function removeRelatedStore(merId: number) {
  if (isReadonly.value) return;
  relatedStores.value = relatedStores.value.filter((item) => item.mer_id !== merId);
}

async function onStatusChange(row: BusinessZoneRow, value: number | string | boolean) {
  const next = Number(value) === 1 ? 1 : 0;
  const prev = row.status;
  if (next === prev) return;
  statusUpdating.value[row.circle_id] = true;
  row.status = next;
  try {
    await updateBusinessZoneStatus(row.circle_id, next);
    ElMessage.success(next === 1 ? '已开启' : '已关闭');
  } catch {
    row.status = prev;
  } finally {
    statusUpdating.value[row.circle_id] = false;
  }
}

async function openInvite(row: BusinessZoneRow) {
  inviteZoneName.value = row.name;
  inviteImageUrl.value = '';
  inviteH5Url.value = '';
  inviteModalApi.open();
  try {
    const payload = await fetchBusinessZoneInvite(row.circle_id);
    inviteH5Url.value = payload.h5_url;
    inviteImageUrl.value = await QRCode.toDataURL(payload.h5_url, {
      margin: 1,
      width: 200,
    });
  } catch {
    ElMessage.error('邀请二维码加载失败');
  }
}

function downloadInviteQR() {
  if (!inviteImageUrl.value) {
    ElMessage.warning('二维码尚未生成');
    return;
  }
  const link = document.createElement('a');
  link.href = inviteImageUrl.value;
  link.download = `邀请入驻-${inviteZoneName.value || '区域'}.png`;
  link.click();
}

async function remove(row: BusinessZoneRow) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”后不可恢复，是否继续？`,
      '删除区域',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' },
    );
    await deleteBusinessZone(row.circle_id);
    ElMessage.success('已删除');
    gridApi.reload();
    void loadParentOptions();
  } catch {
    /* 取消 */
  }
}

onMounted(() => {
  void loadAgentOptions();
  void loadRoleOptions();
  void loadParentOptions();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate(0)">
          新增区域
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status"
          :active-value="1"
          :inactive-value="0"
          :loading="!!statusUpdating[row.circle_id]"
          @change="(val) => onStatusChange(row, val)"
        />
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="Number(row.level) < 2"
          link
          type="primary"
          @click="openCreate(row.circle_id)"
        >
          新增下级
        </ElButton>
        <ElButton link type="primary" @click="openInvite(row)">邀请入驻</ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer :title="formDrawerTitle">
      <ElTabs v-model="activeTab" class="zone-drawer__tabs">
        <ElTabPane label="基础信息" name="basic">
          <ElForm label-width="110px" :disabled="isReadonly">
            <ElFormItem label="上级区域" required>
              <ElSelect
                v-model="form.pid"
                class="w-full"
                filterable
                :disabled="isReadonly || parentLocked || drawerMode === 'edit'"
                placeholder="请选择上级区域"
              >
                <ElOption
                  v-for="item in parentOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                  :disabled="editingID === item.value"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="区域名称" required>
              <ElInput
                v-model="form.name"
                maxlength="20"
                show-word-limit
                placeholder="请输入区域名称"
              />
            </ElFormItem>
            <ElFormItem label="区域代理" required>
              <div class="zone-agent-row">
                <ElButton
                  v-if="!isReadonly"
                  :icon="Plus"
                  @click="openAgentCreate"
                >
                  新增代理人
                </ElButton>
                <ElSelect
                  v-model="form.circle_agent_id"
                  class="zone-agent-row__select"
                  filterable
                  clearable
                  placeholder="请输入选择代理人姓名"
                >
                  <ElOption
                    v-for="item in agentOptions"
                    :key="item.circle_agent_id"
                    :label="item.phone ? `${item.name}（${item.phone}）` : item.name"
                    :value="item.circle_agent_id"
                  />
                </ElSelect>
              </div>
              <div class="field-help">
                区域代理可以查看/管理该区域下的店铺数据。
              </div>
            </ElFormItem>
            <ElFormItem label="身份权限" required>
              <ElSelect
                v-model="form.role_id"
                class="w-full"
                filterable
                clearable
                placeholder="请输入选择身份权限"
              >
                <ElOption
                  v-for="item in roleOptions"
                  :key="item.role_id"
                  :label="item.role_name"
                  :value="item.role_id"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="代理提成">
              <ElRadioGroup v-model="form.commission_type">
                <ElRadio :value="0">默认设置</ElRadio>
                <ElRadio :value="1">单独设置</ElRadio>
              </ElRadioGroup>
            </ElFormItem>
            <ElFormItem v-if="form.commission_type === 1" label="提成比例">
              <div class="zone-rate-row">
                <ElInputNumber
                  v-model="form.commission_rate"
                  :min="0"
                  :max="100"
                  :precision="2"
                  :step="0.01"
                />
                <span class="zone-rate-row__unit">%</span>
              </div>
              <div class="field-help">
                下级区域订单提成 = 平台手续费 ×（本级比例 − 下级比例）
              </div>
              <div class="field-help">
                区域店铺订单提成 = 平台手续费 × 本级比例
              </div>
            </ElFormItem>
            <ElFormItem label="排序">
              <ElInputNumber v-model="form.sort" :min="0" :max="9999" :precision="0" />
              <div class="field-help">数字越大越靠前</div>
            </ElFormItem>
            <ElFormItem label="开启状态">
              <ElSwitch
                v-model="form.status"
                :active-value="1"
                :inactive-value="0"
                active-text="开启"
                inactive-text="关闭"
              />
              <div class="field-help">关闭后，该区域禁止登录</div>
            </ElFormItem>
          </ElForm>
        </ElTabPane>
        <ElTabPane label="关联店铺" name="stores">
          <div class="zone-stores">
            <div v-if="!isReadonly" class="zone-stores__toolbar">
              <ElButton type="primary" @click="storePickerOpen = true">
                选择店铺
              </ElButton>
              <span class="field-help">已选 {{ relatedStores.length }} 家店铺</span>
            </div>
            <ElTable :data="relatedStores" size="small" empty-text="暂无关联店铺">
              <ElTableColumn prop="mer_id" label="店铺ID" width="90" />
              <ElTableColumn prop="mer_name" label="店铺名称" min-width="160" />
              <ElTableColumn prop="real_name" label="联系人" min-width="100" />
              <ElTableColumn prop="mer_phone" label="联系电话" min-width="120" />
              <ElTableColumn v-if="!isReadonly" label="操作" width="80" fixed="right">
                <template #default="{ row }">
                  <ElButton link type="danger" @click="removeRelatedStore(row.mer_id)">
                    移除
                  </ElButton>
                </template>
              </ElTableColumn>
            </ElTable>
          </div>
        </ElTabPane>
      </ElTabs>
    </FormDrawer>

    <AgentDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="代理名称" required>
          <ElInput v-model="agentForm.name" maxlength="64" placeholder="请输入代理名称" />
        </ElFormItem>
        <ElFormItem label="联系电话" required>
          <ElInput v-model="agentForm.phone" maxlength="16" placeholder="请输入联系电话" />
        </ElFormItem>
        <ElFormItem label="身份资质">
          <ImageField
            v-model="agentForm.qualification"
            button-text="从素材库选择"
            :preview-size="96"
            default-library="system"
          />
        </ElFormItem>
        <ElFormItem label="区域代理" required>
          <div class="zone-user-pick">
            <ElButton @click="userPickerOpen = true">选择用户</ElButton>
            <ElTag v-if="pickedUser" class="ml-2" type="info" closable @close="pickedUser = null; agentForm.uid = 0">
              {{ pickedUser.nickname }}（ID:{{ pickedUser.id }}）
              <template v-if="pickedUser.mobile"> · {{ pickedUser.mobile }}</template>
            </ElTag>
            <span v-else class="field-help ml-2">必选；选择 C 端用户并绑定其 UID</span>
          </div>
        </ElFormItem>
        <ElFormItem label="说明">
          <ElInput
            v-model="agentForm.remark"
            type="textarea"
            :rows="3"
            maxlength="255"
            show-word-limit
            placeholder="请输入说明"
          />
        </ElFormItem>
      </ElForm>
    </AgentDrawer>

    <InviteModal>
      <div class="zone-invite">
        <div class="zone-invite__label">H5邀请入驻</div>
        <div v-if="inviteImageUrl" class="zone-invite__qr">
          <img :src="inviteImageUrl" alt="H5邀请入驻二维码" />
        </div>
        <div v-else class="zone-invite__placeholder">二维码生成中…</div>
        <p v-if="inviteH5Url" class="zone-invite__url">{{ inviteH5Url }}</p>
        <ElButton type="primary" class="mt-3" @click="downloadInviteQR">
          下载二维码
        </ElButton>
      </div>
    </InviteModal>

    <StorePickerModal v-model:open="storePickerOpen" :selected="relatedStores" @confirm="onStoresPicked" />
    <UserPickerModal v-model:open="userPickerOpen" @select="onUserPicked" />
  </Page>
</template>

<style scoped>
.zone-drawer__tabs {
  min-height: 420px;
}

.zone-agent-row {
  display: flex;
  gap: 10px;
  align-items: center;
  width: 100%;
}

.zone-agent-row__select {
  flex: 1;
  min-width: 0;
}

.zone-rate-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.zone-rate-row__unit {
  color: hsl(var(--foreground));
}

.field-help {
  margin-top: 4px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 18px;
}

.zone-stores__toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
}

.zone-user-pick {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.zone-invite {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px 20px;
  text-align: center;
}

.zone-invite__label {
  margin-bottom: 16px;
  font-size: 14px;
  color: hsl(var(--foreground));
}

.zone-invite__qr img {
  display: block;
  width: 200px;
  height: 200px;
  background: #fff;
}

.zone-invite__placeholder {
  width: 200px;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: hsl(var(--muted-foreground));
  border: 1px dashed hsl(var(--border));
}

.zone-invite__url {
  margin-top: 12px;
  max-width: 100%;
  word-break: break-all;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}
</style>
