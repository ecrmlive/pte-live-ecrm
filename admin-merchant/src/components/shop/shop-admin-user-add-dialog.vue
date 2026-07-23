<script setup lang="ts">
import type { ShopAdminUserForm, ShopAdminUserRoleOption } from '#/api/core/shop-auth';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { addShopAdminUserApi } from '#/api/core/shop-auth';

defineOptions({ name: 'ShopAdminUserAddDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  roleList?: ShopAdminUserRoleOption[];
}>();

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { class: 'w-full', placeholder: '请输入用户名' },
    fieldName: 'user_name',
    label: '用户名',
    rules: 'required',
  },
  {
    component: 'Select',
    componentProps: {
      class: 'w-full',
      multiple: true,
      options: (props.roleList ?? []).map((item) => ({
        label: item.role_name_h1,
        value: item.role_id,
      })),
    },
    fieldName: 'role_id',
    label: '所属角色',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    componentProps: {
      class: 'w-full',
      placeholder: '请输入登录密码',
      type: 'password',
    },
    fieldName: 'password',
    label: '登录密码',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: {
      class: 'w-full',
      placeholder: '请输入确认密码',
      type: 'password',
    },
    fieldName: 'confirm_password',
    label: '确认密码',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { class: 'w-full', placeholder: '请输入姓名' },
    fieldName: 'real_name',
    label: '姓名',
    rules: 'required',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 96,
    },
    handleSubmit: async (values) => {
      submitting.value = true;
      try {
        const payload: ShopAdminUserForm = {
          confirm_password: String(values.confirm_password ?? ''),
          password: String(values.password ?? ''),
          real_name: String(values.real_name ?? ''),
          role_id: (values.role_id as Array<number | string>) ?? [],
          user_name: String(values.user_name ?? ''),
        };
        await addShopAdminUserApi(payload);
        ElMessage.success('恭喜你，添加成功');
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

const [Modal, modalApi] = useVbenModal({
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
    :destroy-on-close="true"
    title="添加管理员"
    class="w-[600px]"
  >
    <Form />
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
