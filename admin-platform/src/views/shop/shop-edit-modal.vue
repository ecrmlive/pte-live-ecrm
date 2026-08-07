<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElCheckbox,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
} from 'element-plus';

import ShopApi from '#/api/core/shop';
import { hasSpace, isAllSpace, replaceSpace } from '#/utils/form-text';
import { validateBEndPassword } from '#/utils/b-end-password';

import type { ShopRow } from './types';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  shop?: ShopRow;
}>();

const emit = defineEmits<{
  success: [];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);

const form = reactive<ShopRow>({
  app_id: 0,
  app_name: '',
  user_name: '',
  is_recycle: true,
  weixin_service: false,
  expire_time_text: '',
  no_expire: false,
  password: '',
  password_confirm: '',
});

const rules = {
  app_name: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!value) return callback(new Error('请输入商城名称'));
        if (isAllSpace(value)) return callback(new Error('不能全是空格'));
        callback();
      },
    },
  ],
  user_name: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!replaceSpace(value)) return callback(new Error('商家账户名'));
        if (hasSpace(value)) return callback(new Error('不能包含空格'));
        callback();
      },
    },
  ],
  password: [
    {
      trigger: 'change',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!value) return callback();
        const msg = validateBEndPassword(value);
        if (msg) return callback(new Error(msg));
        callback();
      },
    },
  ],
  password_confirm: [
    {
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!form.password) return callback();
        if (value !== form.password) return callback(new Error('确认密码不一致'));
        callback();
      },
    },
  ],
  expire_time_text: [
    {
      required: true,
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (form.no_expire) return callback();
        if (!value) return callback(new Error('请输入过期时间'));
        callback();
      },
    },
  ],
};

function syncForm() {
  if (!props.shop) return;
  Object.assign(form, props.shop);
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      syncForm();
    }
  },
  onConfirm: () => {
    void handleSubmit();
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      drawerApi.setState({ title: '编辑小程序商城' }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(
  () => props.shop?.app_id,
  () => {
    if (open.value) {
      syncForm();
    }
  },
);

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
    drawerApi.setState({ confirmLoading: true });
    try {
      await ShopApi.editShop({ ...form }, true);
      ElMessage.success('恭喜你，修改成功');
      open.value = false;
      emit('success');
    } finally {
      loading.value = false;
      drawerApi.setState({ confirmLoading: false });
    }
  });
}
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :confirm-loading="loading"
    :destroy-on-close="true"
    title="编辑小程序商城"
  >
    <div style="height: 0; overflow: hidden">
      <input type="password" />
    </div>
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="132px">
      <ElFormItem label="商城名称" prop="app_name">
        <ElInput v-model="form.app_name" autocomplete="off" placeholder="请输入商城名称" />
      </ElFormItem>
      <ElFormItem label="过期时间" prop="expire_time_text">
        <ElDatePicker
          v-model="form.expire_time_text"
          :disabled="form.no_expire"
          placeholder="过期时间"
          type="date"
          value-format="YYYY-MM-DD"
        />
        <ElCheckbox v-model="form.no_expire" class="pl-4">永不过期</ElCheckbox>
      </ElFormItem>
      <ElFormItem label="商家账户名" prop="user_name">
        <ElInput v-model="form.user_name" autocomplete="off" placeholder="请输入商家账户名" />
        <p class="text-xs text-[#909399]">注：商家后台用户名</p>
      </ElFormItem>
      <ElFormItem label="商家账户密码" prop="password">
        <ElInput
          v-model="form.password"
          autocomplete="off"
          placeholder="请输入密码"
          type="password"
        />
        <p class="text-xs text-[#909399]">注：商家后台用户密码</p>
      </ElFormItem>
      <ElFormItem label="确认密码" prop="password_confirm">
        <ElInput
          v-model="form.password_confirm"
          autocomplete="off"
          placeholder="请输入确认密码"
          type="password"
        />
      </ElFormItem>
      <ElFormItem label="微信服务商支付">
        <ElCheckbox v-model="form.weixin_service">开启</ElCheckbox>
      </ElFormItem>
    </ElForm>
  </Drawer>
</template>
