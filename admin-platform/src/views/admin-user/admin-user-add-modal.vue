<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import {
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import PlatformAdminUserApi from '#/api/core/platform-admin-user';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  roleOptions: { role_id: number; role_name: string }[];
}>();

const emit = defineEmits<{
  success: [];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);

const form = reactive({
  password: '',
  role_ids: [] as number[],
  user_name: '',
});

function resetForm() {
  form.user_name = '';
  form.password = '';
  form.role_ids = [];
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      resetForm();
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
      drawerApi.setState({ title: '新增平台账号' }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
    drawerApi.setState({ confirmLoading: true });
    try {
      const res = await PlatformAdminUserApi.userAdd(
        {
          password: form.password,
          role_ids: form.role_ids,
          user_name: form.user_name,
        },
        true,
      );
      if (res.code === 1) {
        ElMessage.success(res.msg || '新增成功');
        open.value = false;
        emit('success');
      }
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
    title="新增平台账号"
  >
    <ElForm ref="formRef" :model="form" label-width="100px">
      <ElFormItem
        :rules="[{ message: '请输入用户名', required: true }]"
        label="用户名"
        prop="user_name"
      >
        <ElInput v-model="form.user_name" autocomplete="off" placeholder="登录用户名" />
      </ElFormItem>
      <ElFormItem
        :rules="[{ message: '请输入密码', required: true }]"
        label="密码"
        prop="password"
      >
        <ElInput
          v-model="form.password"
          autocomplete="new-password"
          placeholder="登录密码"
          show-password
          type="password"
        />
      </ElFormItem>
      <ElFormItem label="角色">
        <ElSelect v-model="form.role_ids" class="w-full" multiple placeholder="可选">
          <ElOption
            v-for="item in props.roleOptions"
            :key="item.role_id"
            :label="item.role_name"
            :value="item.role_id"
          />
        </ElSelect>
      </ElFormItem>
    </ElForm>
  </Drawer>
</template>
