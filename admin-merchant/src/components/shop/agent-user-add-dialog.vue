<script setup lang="ts">
import type { MemberListItem } from '#/api/core/member';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { Plus } from '@element-plus/icons-vue';
import { ElButton, ElMessage } from 'element-plus';
import { computed, defineComponent, h, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { addAgentUserApi } from '#/api/core/plus-agent';
import MemberPickerDialog from '#/components/shop/member-picker-dialog.vue';

defineOptions({ name: 'AgentUserAddDialog' });

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  success: [];
}>();

const memberPickerOpen = ref(false);
const submitting = ref(false);
const selectedMember = ref<MemberListItem | null>(null);

const AgentMemberField = defineComponent({
  name: 'AgentMemberField',
  props: {
    modelValue: {
      default: '',
      type: [String, Number],
    },
  },
  emits: ['picked', 'update:modelValue'],
  setup() {
    return () =>
      h('div', { class: 'flex flex-col gap-2' }, [
        h(
          ElButton,
          {
            icon: Plus,
            type: 'primary',
            onClick: () => {
              memberPickerOpen.value = true;
            },
          },
          () => '选择会员',
        ),
        selectedMember.value
          ? h('div', { class: 'flex items-center gap-2' }, [
              selectedMember.value.avatarUrl
                ? h('img', {
                    alt: '',
                    class: 'size-20 rounded object-cover',
                    src: selectedMember.value.avatarUrl,
                  })
                : null,
              h(
                'span',
                `${selectedMember.value.nickName}(ID:${selectedMember.value.user_id})`,
              ),
            ])
          : null,
      ]);
  },
});

function resetForm() {
  selectedMember.value = null;
  void formApi.resetForm();
}

function onMemberPicked(member: MemberListItem) {
  selectedMember.value = member;
  void formApi.setValues({
    mobile: member.mobile ?? '',
    real_name: member.nickName ?? '',
    user_id: member.user_id,
  });
}

const schema = computed((): VbenFormSchema[] => [
  {
    component: markRaw(AgentMemberField),
    fieldName: 'user_id',
    label: '选择用户',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    componentProps: { autocomplete: 'off' },
    fieldName: 'real_name',
    label: '姓名',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { autocomplete: 'off' },
    fieldName: 'mobile',
    label: '手机号',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { autocomplete: 'off', type: 'number' },
    description: '如果没有上级则设置为 0',
    defaultValue: 0,
    fieldName: 'referee_id',
    label: '推荐人ID',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      submitting.value = true;
      try {
        await addAgentUserApi({
          mobile: String(values.mobile ?? ''),
          real_name: String(values.real_name ?? ''),
          referee_id: Number(values.referee_id ?? 0),
          user_id: values.user_id as number | string,
        });
        ElMessage.success('新增成功');
        open.value = false;
        resetForm();
        emit('success');
      } finally {
        submitting.value = false;
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      resetForm();
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    return;
  }
  modalApi.close();
});

async function submit() {
  await formApi.validateAndSubmitForm();
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :destroy-on-close="true"
    class="w-[600px]"
    title="新增"
  >
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>

  <MemberPickerDialog v-model:open="memberPickerOpen" @select="onMemberPicked" />
</template>
