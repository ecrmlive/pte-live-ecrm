<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import PlatformShopUserApi from '#/api/core/platform-shop-user';
import type { ShopUserRoleOption } from '#/api/core/platform-shop-user';
import { validateBEndPassword } from '#/utils/b-end-password';

const props = defineProps<{
  appId: number;
  open: boolean;
  shopUserId?: number;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);
const roleOptions = ref<ShopUserRoleOption[]>([]);

const form = reactive({
  user_name: '',
  real_name: '',
  password: '',
  confirm_password: '',
  role_id: [] as number[],
});

const rules = {
  user_name: [{ message: '请输入用户名', required: true, trigger: 'blur' }],
  real_name: [{ message: '请输入姓名', required: true, trigger: 'blur' }],
  password: [
    {
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!value) return callback();
        const msg = validateBEndPassword(value);
        if (msg) return callback(new Error(msg));
        callback();
      },
    },
  ],
  confirm_password: [
    {
      trigger: 'blur',
      validator: (_: unknown, value: string, callback: (err?: Error) => void) => {
        if (!form.password) return callback();
        if (value !== form.password) return callback(new Error('确认密码不一致'));
        callback();
      },
    },
  ],
  role_id: [{ message: '请选择所属角色', required: true, trigger: 'change', type: 'array' }],
};

watch(
  () => props.open,
  async (open) => {
    if (!open || !props.shopUserId || props.appId <= 0) return;
    loading.value = true;
    try {
      const res = await PlatformShopUserApi.editInfo(props.appId, props.shopUserId, true);
      if (res.code !== 1) return;
      const info = res.data?.info;
      roleOptions.value = res.data?.roleList ?? [];
      form.user_name = info?.user_name ?? '';
      form.real_name = info?.real_name ?? '';
      form.password = '';
      form.confirm_password = '';
      form.role_id = res.data?.role_arr ?? [];
    } finally {
      loading.value = false;
    }
  },
);

function handleClose() {
  emit('update:open', false);
}

async function handleSubmit() {
  if (!formRef.value || !props.shopUserId || props.appId <= 0) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
    try {
      const res = await PlatformShopUserApi.edit(
        {
          app_id: props.appId,
          shop_user_id: props.shopUserId,
          user_name: form.user_name,
          real_name: form.real_name,
          password: form.password || undefined,
          confirm_password: form.confirm_password || undefined,
          role_id: form.role_id,
        },
        true,
      );
      if (res.code === 1) {
        ElMessage.success(res.msg || '更新成功');
        emit('update:open', false);
        emit('success');
      }
    } finally {
      loading.value = false;
    }
  });
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :model-value="open"
    title="编辑商城管理员"
    width="560px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
  >
    <ElForm ref="formRef" v-loading="loading" :model="form" :rules="rules" label-width="100px">
      <ElFormItem label="用户名" prop="user_name">
        <ElInput v-model="form.user_name" autocomplete="off" />
      </ElFormItem>
      <ElFormItem label="姓名" prop="real_name">
        <ElInput v-model="form.real_name" />
      </ElFormItem>
      <ElFormItem label="所属角色" prop="role_id">
        <ElSelect v-model="form.role_id" class="w-full" multiple placeholder="请选择角色">
          <ElOption
            v-for="item in roleOptions"
            :key="item.role_id"
            :label="item.role_name_h1 || item.role_name"
            :value="item.role_id"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="新密码" prop="password">
        <ElInput
          v-model="form.password"
          autocomplete="new-password"
          placeholder="留空表示不修改"
          show-password
          type="password"
        />
      </ElFormItem>
      <ElFormItem label="确认密码" prop="confirm_password">
        <ElInput
          v-model="form.confirm_password"
          autocomplete="new-password"
          show-password
          type="password"
        />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">确定</ElButton>
    </template>
  </ElDialog>
</template>
