<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm, z } from '#/adapter/form';
import { editShopPasswordApi } from '#/api/core/passport';
import { validateBEndPassword } from '#/utils/b-end-password';

defineOptions({ name: 'ShopUpdatePasswordDialog' });

const emit = defineEmits<{ success: [] }>();

const open = defineModel<boolean>('open', { default: false });

const submitting = ref(false);

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      autocomplete: 'off',
      class: 'w-full',
      placeholder: '请输入原始密码',
      showPasswordOnClick: true,
      type: 'password',
    },
    fieldName: 'oldpass',
    label: '原始密码',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: {
      autocomplete: 'off',
      class: 'w-full',
      placeholder: '请输入新密码',
      showPasswordOnClick: true,
      type: 'password',
    },
    fieldName: 'password',
    label: '新密码',
    rules: z.string().superRefine((value, ctx) => {
      const msg = validateBEndPassword(value);
      if (msg) {
        ctx.addIssue({ code: 'custom', message: msg });
      }
    }),
  },
  {
    component: 'Input',
    componentProps: {
      autocomplete: 'off',
      class: 'w-full',
      placeholder: '请再次输入新密码',
      showPasswordOnClick: true,
      type: 'password',
    },
    dependencies: {
      rules(values) {
        return z
          .string()
          .min(1, '请确认新密码')
          .refine((val) => val === values.password, {
            message: '两次密码不相同',
          });
      },
      triggerFields: ['password', 'confirmPass'],
    },
    fieldName: 'confirmPass',
    label: '确认新密码',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      const password = String(values.password ?? '');
      const confirmPass = String(values.confirmPass ?? '');
      if (password !== confirmPass) {
        ElMessage.error('两次密码不相同');
        return;
      }

      submitting.value = true;
      try {
        await editShopPasswordApi({
          confirmPass,
          oldpass: String(values.oldpass ?? ''),
          password,
        });
        ElMessage.success('修改成功');
        open.value = false;
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

function resetForm() {
  void formApi.resetForm();
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (!isOpen) {
      resetForm();
      return;
    }
    resetForm();
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

function closeModal() {
  open.value = false;
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :destroy-on-close="true"
    class="w-[480px]"
    title="修改密码"
  >
    <Form />
    <template #footer>
      <ElButton @click="closeModal">取 消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确 定</ElButton>
    </template>
  </Modal>
</template>
