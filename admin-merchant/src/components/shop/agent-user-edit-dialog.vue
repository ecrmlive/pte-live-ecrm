<script setup lang="ts">
import type { AgentGradeOption, AgentUserItem } from '#/api/core/plus-agent';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { matchFormSelectValue } from '#/utils/form-select';
import {
  editAgentUserApi,
  getAgentRefereeCheckApi,
} from '#/api/core/plus-agent';

defineOptions({ name: 'AgentUserEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  gradeList?: AgentGradeOption[];
  user?: AgentUserItem;
}>();

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const oldRefereeId = ref(0);
const avatarUrl = ref('');
const refereeNickName = ref('-');

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { disabled: true },
    fieldName: 'nickName',
    label: '昵称',
  },
  {
    component: 'Input',
    componentProps: { disabled: true },
    fieldName: 'real_name',
    label: '姓名',
  },
  {
    component: 'Input',
    fieldName: 'mobile',
    label: '手机号',
  },
  {
    component: 'Select',
    componentProps: {
      options: (props.gradeList ?? []).map((item) => ({
        label: item.name,
        value: item.grade_id,
      })),
      placeholder: '-请选择等级-',
    },
    fieldName: 'grade_id',
    label: '等级',
  },
  {
    component: 'Input',
    componentProps: { type: 'number' },
    description: '如果没有上级则设置为 0',
    fieldName: 'referee_id',
    label: '推荐人ID',
  },
  {
    component: 'Input',
    componentProps: { disabled: true },
    fieldName: 'refereeNickName',
    label: '推荐人昵称',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      const params = {
        grade_id: Number(values.grade_id ?? 0),
        mobile: String(values.mobile ?? ''),
        real_name: String(values.real_name ?? ''),
        user_id: Number(values.user_id ?? 0),
      };
      const refereeId = Number(values.referee_id ?? 0);
      if (refereeId !== oldRefereeId.value) {
        const res = await getAgentRefereeCheckApi(refereeId);
        if (refereeId > 0 && !res.model) {
          ElMessage.error('上级用户不存在');
          return;
        }
        if (params.user_id === refereeId) {
          ElMessage.error('不能绑定自己');
          return;
        }
        const msg = res.model
          ? `确定要修改上级到${res.model.real_name}？`
          : '确定要修改上级为空？';
        await ElMessageBox.confirm(msg, '提示', { type: 'warning' });
        await saveUser({ ...params, referee_id: refereeId });
        return;
      }
      await saveUser(params);
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

async function saveUser(params: {
  grade_id?: number;
  mobile?: string;
  real_name?: string;
  referee_id?: number | string;
  user_id: number;
}) {
  submitting.value = true;
  try {
    await editAgentUserApi(params);
    ElMessage.success('修改成功');
    open.value = false;
    emit('success');
  } finally {
    submitting.value = false;
  }
}

function fillForm(user: AgentUserItem) {
  avatarUrl.value = user.avatarUrl ?? '';
  refereeNickName.value = user.referee?.nickName || '-';
  oldRefereeId.value = user.referee_id;
  const gradeOptions = (props.gradeList ?? []).map((item) => ({
    label: item.name,
    value: item.grade_id,
  }));
  void formApi.setValues({
    grade_id: matchFormSelectValue(user.grade_id, gradeOptions),
    mobile: user.mobile,
    nickName: user.nickName ?? '',
    real_name: user.real_name,
    referee_id: user.referee_id,
    refereeNickName: refereeNickName.value,
    user_id: user.user_id,
  });
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen && props.user) {
      fillForm(props.user);
    }
  },
});

watch(
  () => [open.value, props.user] as const,
  ([visible, user]) => {
    if (visible) {
      modalApi.open();
      if (user) fillForm(user);
      return;
    }
    modalApi.close();
  },
);

async function submit() {
  await formApi.validateAndSubmitForm();
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :destroy-on-close="true"
    class="w-[640px]"
    title="编辑分销商"
  >
    <div v-if="avatarUrl" class="agent-user-edit__avatar">
      <img :alt="user?.nickName || 'avatar'" :src="avatarUrl" />
    </div>
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>

<style scoped>
.agent-user-edit__avatar {
  margin: 0 0 16px 100px;
}

.agent-user-edit__avatar img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid hsl(var(--border));
}
</style>
