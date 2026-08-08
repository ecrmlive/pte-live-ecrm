<script setup lang="ts">
import type { PropType } from 'vue';

import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, defineComponent, h, markRaw, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElAvatar,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createBusinessZoneAgent,
  fetchBusinessZoneAgent,
  fetchBusinessZoneAgents,
  resetBusinessZoneAgentPassword,
  revokeBusinessZoneAgent,
  updateBusinessZoneAgent,
  type BusinessZoneAgentRow,
} from '#/api/core/ecrm';
import UserPickerModal, {
  type PickedPlatformUser,
} from '#/components/ecrm/user-picker-modal.vue';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit';
type UserSearchField = 'uid' | 'user_phone' | 'nickname';
type UserSearchValue = { field: UserSearchField; keyword: string };

const QUALIFICATION_MAX = 5;

const USER_SEARCH_OPTIONS: Array<{ label: string; value: UserSearchField }> = [
  { label: 'UID', value: 'uid' },
  { label: '手机号', value: 'user_phone' },
  { label: '用户昵称', value: 'nickname' },
];

function userSearchPlaceholder(field: UserSearchField) {
  if (field === 'uid') return '请输入UID';
  if (field === 'user_phone') return '请输入手机号';
  if (field === 'nickname') return '请输入用户昵称';
  return '请输入用户信息';
}

/** 用户搜索：左侧类型 Select + 右侧关键词 Input，绑定同一表单项。 */
const UserSearchComposite = defineComponent({
  name: 'UserSearchComposite',
  props: {
    modelValue: {
      type: Object as PropType<UserSearchValue>,
      default: () => ({ field: 'uid' as UserSearchField, keyword: '' }),
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    function patch(partial: Partial<UserSearchValue>) {
      emit('update:modelValue', {
        field: props.modelValue?.field || 'uid',
        keyword: props.modelValue?.keyword || '',
        ...partial,
      });
    }
    return () => {
      const field = (props.modelValue?.field || 'uid') as UserSearchField;
      const keyword = props.modelValue?.keyword || '';
      return h('div', { class: 'agent-user-search' }, [
        h(
          ElSelect,
          {
            modelValue: field,
            'onUpdate:modelValue': (value: UserSearchField) =>
              patch({ field: value }),
            class: 'agent-user-search__type',
          },
          () =>
            USER_SEARCH_OPTIONS.map((opt) =>
              h(ElOption, {
                key: opt.value,
                label: opt.label,
                value: opt.value,
              }),
            ),
        ),
        h(ElInput, {
          modelValue: keyword,
          'onUpdate:modelValue': (value: string) => patch({ keyword: value }),
          clearable: true,
          class: 'agent-user-search__keyword',
          placeholder: userSearchPlaceholder(field),
        }),
      ]);
    };
  },
});

const drawerMode = ref<DrawerMode>('create');
const editingID = ref(0);
const resetTarget = ref<BusinessZoneAgentRow>();
const linkedUser = ref<PickedPlatformUser | null>(null);
const userPickerOpen = ref(false);

const form = reactive({
  name: '',
  phone: '',
  qualificationImages: [''] as string[],
  uid: 0,
  remark: '',
});

const passwordReset = reactive({
  password: '',
  confirmPassword: '',
  reason: '',
});

const drawerTitle = computed(() =>
  drawerMode.value === 'edit' ? '编辑区域代理' : '新增区域代理',
);

function normalizeQualificationSlots(urls: string[]): string[] {
  const filled = urls.map((item) => item.trim()).filter(Boolean);
  if (!filled.length) return [''];
  if (filled.length >= QUALIFICATION_MAX) {
    return filled.slice(0, QUALIFICATION_MAX);
  }
  return [...filled, ''];
}

function parseQualificationImages(raw?: string): string[] {
  const text = String(raw ?? '').trim();
  if (!text) return [''];
  try {
    const parsed = JSON.parse(text) as unknown;
    if (Array.isArray(parsed)) {
      return normalizeQualificationSlots(
        parsed.map((item) => String(item ?? '').trim()),
      );
    }
  } catch {
    /* 兼容历史纯文本/单 URL */
  }
  if (text.startsWith('http') || text.startsWith('/')) {
    return normalizeQualificationSlots([text]);
  }
  return [''];
}

function serializeQualificationImages(images: string[]): string {
  const urls = images.map((item) => item.trim()).filter(Boolean);
  return urls.length ? JSON.stringify(urls) : '';
}

function setQualificationImage(index: number, value?: string | null) {
  form.qualificationImages[index] = String(value ?? '').trim();
  form.qualificationImages = normalizeQualificationSlots(
    form.qualificationImages,
  );
}

function addQualificationSlot() {
  if (form.qualificationImages.length >= QUALIFICATION_MAX) return;
  const last = form.qualificationImages[form.qualificationImages.length - 1];
  if (!String(last ?? '').trim()) return;
  form.qualificationImages.push('');
}

/** 负责区域：名称(提成%)，顿号连接。 */
function formatResponsibleRegions(row: BusinessZoneAgentRow) {
  const parts = (row.circles || [])
    .filter((item) => Number(item?.type ?? 0) === 0)
    .map((item) => {
      const name = String(item?.name ?? '').trim();
      if (!name) return '';
      const rate = Number(item?.commission_rate ?? 0);
      return `${name}(${Number.isFinite(rate) ? rate : 0}%)`;
    })
    .filter(Boolean);
  if (parts.length) return parts.join('、');
  // 无 type=0 时回退全部绑定（兼容旧数据）
  const fallback = (row.circles || [])
    .map((item) => {
      const name = String(item?.name ?? '').trim();
      if (!name) return '';
      const rate = Number(item?.commission_rate ?? 0);
      return `${name}(${Number.isFinite(rate) ? rate : 0}%)`;
    })
    .filter(Boolean);
  return fallback.length ? fallback.join('、') : '—';
}

function formatUserInfo(row: BusinessZoneAgentRow) {
  if (!row.uid) return '—';
  return `${row.nickname || '用户'} | ${row.uid}`;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入代理名称',
    },
    fieldName: 'name',
    label: '代理名称',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入联系电话',
    },
    fieldName: 'phone',
    label: '联系电话',
  },
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '创建时间',
    componentProps: {
      ...LIST_DATE_RANGE_FIELD.componentProps,
      startPlaceholder: '开始日期',
      endPlaceholder: '结束日期',
    },
  },
  {
    component: markRaw(UserSearchComposite),
    defaultValue: { field: 'uid', keyword: '' } satisfies UserSearchValue,
    fieldName: 'user_search',
    label: '用户搜索',
  },
]);

const gridOptions: VxeGridProps<BusinessZoneAgentRow> = {
  columns: [
    { field: 'name', minWidth: 120, title: '代理名称' },
    { field: 'phone', title: '联系电话', width: 140 },
    {
      field: 'uid',
      formatter: ({ row }) => formatUserInfo(row),
      minWidth: 160,
      title: '用户信息',
    },
    {
      field: 'circles',
      minWidth: 200,
      slots: { default: 'circles' },
      title: '负责区域',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ minWidth: 200, title: '代理操作' }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const userSearch = (formValues?.user_search ||
          {}) as Partial<UserSearchValue>;
        const userField = (userSearch.field || 'uid') as UserSearchField;
        const userKeyword = String(userSearch.keyword ?? '').trim();
        // 仅下发当前用户搜索类型对应参数，不附带其它类型键。
        const params: Parameters<typeof fetchBusinessZoneAgents>[0] = {
          page: page.currentPage,
          limit: page.pageSize,
          type: 0,
          name: String(formValues?.name ?? '').trim() || undefined,
          phone: String(formValues?.phone ?? '').trim() || undefined,
          date_from: range[0] ? String(range[0]) : undefined,
          date_to: range[1] ? String(range[1]) : undefined,
        };
        if (userKeyword) {
          if (userField === 'uid') {
            const uid = Number(userKeyword);
            if (Number.isFinite(uid) && uid > 0) params.uid = uid;
          } else if (userField === 'user_phone') {
            params.user_phone = userKeyword;
          } else if (userField === 'nickname') {
            params.nickname = userKeyword;
          }
        }
        const result = await fetchBusinessZoneAgents(params);
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_agent_id' },
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
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const [ResetDrawer, resetDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确认重置',
  cancelText: '取消',
  placement: 'right',
  title: '重置区域代理后台密码',
  onConfirm: async () => submitPasswordReset(),
});

function resetForm() {
  Object.assign(form, {
    name: '',
    phone: '',
    qualificationImages: [''],
    uid: 0,
    remark: '',
  });
  linkedUser.value = null;
}

function fillForm(row: BusinessZoneAgentRow) {
  Object.assign(form, {
    name: row.name || '',
    phone: row.phone || '',
    qualificationImages: parseQualificationImages(row.qualification),
    uid: row.uid || 0,
    remark: row.remark || '',
  });
  if (row.uid) {
    linkedUser.value = {
      id: row.uid,
      nickname: row.nickname || `用户#${row.uid}`,
      avatar_url: row.avatar_url || '',
      mobile: row.phone || '',
    };
  } else {
    linkedUser.value = null;
  }
}

function openCreate() {
  drawerMode.value = 'create';
  editingID.value = 0;
  resetForm();
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
}

async function openEdit(row: BusinessZoneAgentRow) {
  drawerMode.value = 'edit';
  editingID.value = row.circle_agent_id;
  resetForm();
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
  formDrawerApi.lock();
  try {
    const detail = await fetchBusinessZoneAgent(row.circle_agent_id);
    fillForm(detail);
  } finally {
    formDrawerApi.unlock();
  }
}

function openUserPicker() {
  userPickerOpen.value = true;
}

function onUserPicked(user: PickedPlatformUser) {
  linkedUser.value = user;
  form.uid = user.id;
  if (!form.phone.trim() && user.mobile) {
    form.phone = user.mobile;
  }
}

function clearLinkedUser() {
  linkedUser.value = null;
  form.uid = 0;
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入代理名称');
    return;
  }
  if (!form.phone.trim()) {
    ElMessage.warning('请输入联系电话');
    return;
  }
  if (!form.uid || form.uid <= 0) {
    ElMessage.warning('请选择区域代理关联用户');
    return;
  }
  formDrawerApi.lock();
  const isEdit = drawerMode.value === 'edit' && editingID.value > 0;
  try {
    const payload = {
      type: 0 as const,
      name: form.name.trim(),
      phone: form.phone.trim(),
      uid: form.uid,
      qualification: serializeQualificationImages(form.qualificationImages),
      remark: form.remark.trim(),
      auto_approve: true,
    };
    if (isEdit) {
      await updateBusinessZoneAgent(editingID.value, payload);
    } else {
      await createBusinessZoneAgent(payload);
    }
    formDrawerApi.close();
    ElMessage.success(isEdit ? '已保存' : '已新增区域代理');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

function openPasswordReset(row: BusinessZoneAgentRow) {
  resetTarget.value = row;
  Object.assign(passwordReset, {
    password: '',
    confirmPassword: '',
    reason: '',
  });
  resetDrawerApi.open();
}

async function submitPasswordReset() {
  const reason = passwordReset.reason.trim();
  const target = resetTarget.value;
  if (
    !target ||
    passwordReset.password.length < 12 ||
    passwordReset.password.length > 72 ||
    passwordReset.password !== passwordReset.confirmPassword ||
    reason.length < 2 ||
    reason.length > 500
  ) {
    ElMessage.warning(
      '请填写两次一致的 12 至 72 位新密码和 2 至 500 字的重置原因',
    );
    return;
  }
  resetDrawerApi.lock();
  try {
    await resetBusinessZoneAgentPassword(target.circle_agent_id, {
      password: passwordReset.password,
      reason,
      idempotency_key: `agent-password-${target.circle_agent_id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('后台密码已重置，该代理旧后台会话已失效');
    resetDrawerApi.close();
  } finally {
    resetDrawerApi.unlock();
  }
}

async function remove(row: BusinessZoneAgentRow) {
  try {
    await ElMessageBox.confirm(
      `确定删除代理「${row.name}」吗？删除为资格撤销，不会硬删历史审核与结算事实；已关联区域或仍有佣金余额时无法撤销。`,
      '提示',
      {
        type: 'warning',
        confirmButtonText: '删除',
        cancelButtonText: '取消',
      },
    );
    await revokeBusinessZoneAgent(row.circle_agent_id, {
      reason: `删除区域代理 ${row.name}`,
      idempotency_key: `agent-delete-${row.circle_agent_id}-${Date.now()}`,
    });
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

/** TODO: 待对接区域代理导出 API（当前无独立 export 接口） */
function exportList() {
  ElMessage.info('导出功能待对接，暂不可用');
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增
        </ElButton>
        <ElButton @click="exportList">导出</ElButton>
      </template>
      <template #circles="{ row }">
        {{ formatResponsibleRegions(row) }}
      </template>
      <template #action="{ row }">
        <ElButton link type="warning" @click="openPasswordReset(row)">
          重置密码
        </ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer :title="drawerTitle">
      <ElForm label-width="110px">
        <ElFormItem label="代理名称" required>
          <ElInput
            v-model="form.name"
            clearable
            maxlength="64"
            placeholder="请输入代理名称"
          />
        </ElFormItem>
        <ElFormItem label="联系电话" required>
          <ElInput
            v-model="form.phone"
            clearable
            maxlength="16"
            placeholder="请输入联系电话"
          />
        </ElFormItem>
        <ElFormItem label="身份资质">
          <div class="agent-image-slots">
            <ImageField
              v-for="(_, index) in form.qualificationImages"
              :key="`qual-${index}`"
              :model-value="form.qualificationImages[index]"
              :preview-size="72"
              :show-button="false"
              default-library="system"
              @update:model-value="(value) => setQualificationImage(index, value)"
            />
            <ElButton
              v-if="
                form.qualificationImages.length < QUALIFICATION_MAX &&
                !!String(
                  form.qualificationImages[
                    form.qualificationImages.length - 1
                  ] || '',
                ).trim()
              "
              plain
              type="primary"
              @click="addQualificationSlot"
            >
              新增
            </ElButton>
          </div>
        </ElFormItem>
        <ElFormItem label="区域代理" required>
          <div class="linked-user-row">
            <div v-if="linkedUser" class="linked-user-summary">
              <ElAvatar
                v-if="linkedUser.avatar_url"
                :size="32"
                :src="linkedUser.avatar_url"
              />
              <ElAvatar v-else :size="32">
                {{ (linkedUser.nickname || '?').slice(0, 1) }}
              </ElAvatar>
              <div class="linked-user-meta">
                <div class="linked-user-name">{{ linkedUser.nickname }}</div>
                <div class="linked-user-id">
                  UID {{ linkedUser.id
                  }}<template v-if="linkedUser.mobile">
                    · {{ linkedUser.mobile }}</template
                  >
                </div>
              </div>
              <ElButton link type="danger" @click="clearLinkedUser">
                清空
              </ElButton>
            </div>
            <ElButton type="primary" plain @click="openUserPicker">
              选择用户
            </ElButton>
          </div>
          <div class="field-tip">必选；选择 C 端用户并绑定其 UID</div>
        </ElFormItem>
        <ElFormItem label="说明">
          <ElInput
            v-model="form.remark"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
            placeholder="请输入说明"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <ResetDrawer>
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="仅对已关联且启用的统一后台区域账号生效；不回显、不记录或回传密码，提交后旧后台会话立即失效。"
      />
      <ElForm label-width="108px" @submit.prevent="submitPasswordReset">
        <ElFormItem label="代理">
          <ElInput :model-value="resetTarget?.name || ''" disabled />
        </ElFormItem>
        <ElFormItem label="新密码" required>
          <ElInput
            v-model="passwordReset.password"
            autocomplete="new-password"
            show-password
            type="password"
          />
        </ElFormItem>
        <ElFormItem label="确认新密码" required>
          <ElInput
            v-model="passwordReset.confirmPassword"
            autocomplete="new-password"
            show-password
            type="password"
          />
        </ElFormItem>
        <ElFormItem label="重置原因" required>
          <ElInput
            v-model="passwordReset.reason"
            maxlength="500"
            show-word-limit
            type="textarea"
            :rows="3"
          />
        </ElFormItem>
      </ElForm>
    </ResetDrawer>

    <UserPickerModal v-model:open="userPickerOpen" @select="onUserPicked" />
  </Page>
</template>

<style scoped>
.agent-image-slots {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.field-tip {
  width: 100%;
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.4;
  color: #909399;
}

.linked-user-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  width: 100%;
}

.linked-user-summary {
  display: flex;
  gap: 10px;
  align-items: center;
  min-width: 0;
  padding: 6px 10px;
  background: var(--el-fill-color-light);
  border-radius: 8px;
}

.linked-user-meta {
  min-width: 0;
}

.linked-user-name {
  font-size: 14px;
  line-height: 1.3;
  color: var(--el-text-color-primary);
}

.linked-user-id {
  font-size: 12px;
  line-height: 1.3;
  color: var(--el-text-color-secondary);
}

:deep(.agent-user-search) {
  display: flex;
  width: 100%;
  gap: 8px;
  align-items: center;
}

:deep(.agent-user-search__type) {
  width: 112px;
  flex-shrink: 0;
}

:deep(.agent-user-search__keyword) {
  flex: 1;
  min-width: 0;
}
</style>
