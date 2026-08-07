<script setup lang="ts">
import type { ShopAdminUserForm, ShopAdminUserRoleOption } from '#/api/core/shop-auth';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { editShopAdminUserApi, getShopAdminUserEditInfoApi } from '#/api/core/shop-auth';

defineOptions({ name: 'ShopAdminUserEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  shopUserId?: number;
}>();

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const submitting = ref(false);
const roleList = ref<ShopAdminUserRoleOption[]>([]);

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
      options: roleList.value.map((item) => ({
        label: item.role_name_h1,
        value: item.role_id,
      })),
    },
    fieldName: 'access_id',
    label: '所属角色',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    componentProps: {
      class: 'w-full',
      placeholder: '不修改请留空',
      type: 'password',
    },
    fieldName: 'password',
    label: '登录密码',
  },
  {
    component: 'Input',
    componentProps: {
      class: 'w-full',
      placeholder: '请输入确认密码',
      type: 'password',
    },
    dependencies: {
      show: (values) => Boolean(values.password),
      triggerFields: ['password'],
    },
    fieldName: 'confirm_password',
    label: '确认密码',
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
          access_id: (values.access_id as Array<number | string>) ?? [],
          confirm_password: String(values.confirm_password ?? ''),
          password: String(values.password ?? ''),
          real_name: String(values.real_name ?? ''),
          shop_user_id: Number(values.shop_user_id ?? 0),
          user_name: String(values.user_name ?? ''),
        };
        if (!payload.password) {
          delete payload.password;
          delete payload.confirm_password;
        }
        await editShopAdminUserApi(payload);
        ElMessage.success('恭喜你，修改成功');
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

async function loadData() {
  if (!props.shopUserId) return;
  loading.value = true;
  try {
    const res = await getShopAdminUserEditInfoApi(props.shopUserId);
    roleList.value = res.roleList ?? [];
    formApi.updateSchema(schema.value);
    await formApi.setValues({
      access_id: res.role_arr ?? [],
      confirm_password: '',
      password: '',
      real_name: res.info?.real_name ?? '',
      shop_user_id: props.shopUserId,
      user_name: res.info?.user_name ?? '',
    });
  } finally {
    loading.value = false;
  }
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen && props.shopUserId) {
      void loadData();
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
    title="修改管理员"
    class="w-[600px]"
  >
    <div v-loading="loading">
      <Form />
    </div>
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
