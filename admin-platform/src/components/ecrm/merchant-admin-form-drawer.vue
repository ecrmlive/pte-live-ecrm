<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElAvatar,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import {
  createBusinessZoneAgent,
  fetchBusinessZoneAgent,
  fetchBusinessZones,
  updateBusinessZoneAgent,
  type BusinessZoneAgentRow,
  type BusinessZoneRow,
} from '#/api/core/ecrm';
import UserPickerModal, {
  type PickedPlatformUser,
} from '#/components/ecrm/user-picker-modal.vue';

const props = withDefaults(
  defineProps<{
    /** 是否展示「负责商户」多选（商户管理员列表需要；商户表单内新增时可关） */
    showResponsibleMerchants?: boolean;
  }>(),
  { showResponsibleMerchants: true },
);

const emit = defineEmits<{
  saved: [row: BusinessZoneAgentRow];
}>();

const editingId = ref(0);
const zoneOptions = ref<Array<{ label: string; value: number }>>([]);
const linkedUser = ref<PickedPlatformUser | null>(null);
const userPickerOpen = ref(false);

const form = reactive({
  name: '',
  phone: '',
  account: '',
  password: '000000',
  uid: 0,
  circle_ids: [] as number[],
});

const isEdit = computed(() => editingId.value > 0);
const drawerTitle = computed(() =>
  isEdit.value ? '编辑商户管理员' : '新增商户管理员',
);

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  Object.assign(form, {
    name: '',
    phone: '',
    account: '',
    password: '000000',
    uid: 0,
    circle_ids: [],
  });
  linkedUser.value = null;
}

async function loadZoneOptions() {
  if (!props.showResponsibleMerchants) return;
  try {
    const result = await fetchBusinessZones({
      page: 1,
      limit: 200,
      type: 1,
      status: 1,
    });
    zoneOptions.value = (result.list || []).map((row: BusinessZoneRow) => ({
      label: row.name,
      value: row.circle_id,
    }));
  } catch {
    zoneOptions.value = [];
  }
}

async function openCreate() {
  editingId.value = 0;
  resetForm();
  await loadZoneOptions();
  formDrawerApi.setState({ title: drawerTitle.value }).open();
}

async function openEdit(row: BusinessZoneAgentRow) {
  editingId.value = row.circle_agent_id;
  resetForm();
  await loadZoneOptions();
  formDrawerApi.setState({ title: drawerTitle.value }).open();
  formDrawerApi.lock();
  try {
    const detail = await fetchBusinessZoneAgent(row.circle_agent_id);
    Object.assign(form, {
      name: detail.name || '',
      phone: detail.phone || '',
      account: detail.account || detail.admin?.account || '',
      password: '',
      uid: detail.uid || 0,
      circle_ids: [...(detail.circle_ids || [])],
    });
    if (detail.uid) {
      linkedUser.value = {
        id: detail.uid,
        nickname: detail.nickname || `用户#${detail.uid}`,
        avatar_url: detail.avatar_url || '',
        mobile: '',
      };
    }
    // 编辑回填可能含已关闭商户，补进选项以免只显示 ID。
    for (const circle of detail.circles || []) {
      if (!zoneOptions.value.some((item) => item.value === circle.circle_id)) {
        zoneOptions.value.push({
          label: circle.name,
          value: circle.circle_id,
        });
      }
    }
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
}

function clearLinkedUser() {
  linkedUser.value = null;
  form.uid = 0;
}

watch(
  () => form.phone,
  (phone, prev) => {
    if (isEdit.value) return;
    if (!form.account || form.account === prev) {
      form.account = phone;
    }
  },
);

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入管理姓名');
    return;
  }
  if (!form.phone.trim()) {
    ElMessage.warning('请输入手机号码');
    return;
  }
  if (!isEdit.value && !form.account.trim()) {
    ElMessage.warning('请输入登录账号');
    return;
  }
  if (!isEdit.value && !form.password.trim()) {
    ElMessage.warning('请输入登录密码');
    return;
  }
  if (isEdit.value && form.password && form.password.trim().length < 6) {
    ElMessage.warning('登录密码至少 6 位，留空则不修改');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      type: 1 as const,
      name: form.name.trim(),
      phone: form.phone.trim(),
      uid: form.uid || 0,
      qualification: '',
      remark: '',
      payment_method: 0,
      payment_name: '',
      business_name: '',
      ...(props.showResponsibleMerchants
        ? { circle_ids: [...form.circle_ids] }
        : {}),
    };
    let saved: BusinessZoneAgentRow;
    if (isEdit.value) {
      saved = await updateBusinessZoneAgent(editingId.value, {
        ...payload,
        ...(form.password.trim() ? { password: form.password } : {}),
      });
    } else {
      saved = await createBusinessZoneAgent({
        ...payload,
        account: form.account.trim(),
        password: form.password,
      });
    }
    ElMessage.success(isEdit.value ? '管理员已更新' : '管理员已新增');
    formDrawerApi.close();
    emit('saved', saved);
  } finally {
    formDrawerApi.unlock();
  }
}

onMounted(() => {
  void loadZoneOptions();
});

defineExpose({ openCreate, openEdit });
</script>

<template>
  <FormDrawer :title="drawerTitle">
    <ElForm label-width="110px" class="pt-2">
      <ElFormItem label="管理姓名" required>
        <ElInput
          v-model="form.name"
          maxlength="64"
          placeholder="请输入管理员姓名"
          class="max-w-[460px]"
        />
      </ElFormItem>
      <ElFormItem label="手机号码" required>
        <ElInput
          v-model="form.phone"
          placeholder="请输入手机号码"
          class="max-w-[460px]"
        />
        <div class="field-tip">
          手机号码为商户管理的登录账号，登录密码默认000000
        </div>
      </ElFormItem>
      <ElFormItem label="登录账号" required>
        <ElInput
          v-model="form.account"
          :disabled="isEdit"
          placeholder="请输入登录账号"
          class="max-w-[460px]"
        />
      </ElFormItem>
      <ElFormItem :label="isEdit ? '登录密码' : '登录密码'" :required="!isEdit">
        <ElInput
          v-model="form.password"
          type="password"
          show-password
          autocomplete="new-password"
          :placeholder="isEdit ? '留空则不修改' : '请输入登录密码'"
          class="max-w-[460px]"
        />
      </ElFormItem>
      <ElFormItem label="关联用户">
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
              <div class="linked-user-id">UID {{ linkedUser.id }}</div>
            </div>
            <ElButton link type="danger" @click="clearLinkedUser">
              清空
            </ElButton>
          </div>
          <ElButton type="primary" plain @click="openUserPicker">
            选择用户
          </ElButton>
        </div>
        <div class="field-tip">可选；选择 C 端用户后提交时关联其 UID</div>
      </ElFormItem>
      <ElFormItem v-if="showResponsibleMerchants" label="负责商户">
        <ElSelect
          v-model="form.circle_ids"
          multiple
          filterable
          clearable
          collapse-tags
          collapse-tags-tooltip
          placeholder="请选择负责商户（可多选）"
          class="max-w-[460px] w-full"
        >
          <ElOption
            v-for="item in zoneOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </ElSelect>
      </ElFormItem>
    </ElForm>
  </FormDrawer>

  <UserPickerModal v-model:open="userPickerOpen" @select="onUserPicked" />
</template>

<style scoped>
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
</style>
