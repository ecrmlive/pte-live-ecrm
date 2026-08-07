<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import {
  ElCheckbox,
  ElForm,
  ElFormItem,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import PlatformAdminUserApi from '#/api/core/platform-admin-user';
import type { PlatformAdminUserRow } from '#/api/core/platform-admin-user';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  roleOptions: { role_id: number; role_name: string }[];
  user?: PlatformAdminUserRow;
}>();

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const formData = reactive({
  is_super: false,
  role_ids: [] as number[],
});

async function loadForm() {
  if (!props.user) return;
  loading.value = true;
  try {
    const res = await PlatformAdminUserApi.userEditInfo(
      props.user.admin_user_id,
      true,
    );
    const data = res.data as {
      is_super?: number;
      role_ids?: number[];
      user?: PlatformAdminUserRow;
    };
    formData.is_super = Number(data.user?.is_super ?? props.user.is_super) === 1;
    formData.role_ids = data.role_ids ?? props.user.role_ids ?? [];
  } finally {
    loading.value = false;
  }
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void loadForm();
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
      drawerApi.setState({ title: '分配角色' }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(
  () => props.user?.admin_user_id,
  () => {
    if (open.value) {
      void loadForm();
    }
  },
);

async function handleSubmit() {
  if (!props.user) return;
  loading.value = true;
  drawerApi.setState({ confirmLoading: true });
  try {
    const res = await PlatformAdminUserApi.userEdit(
      props.user.admin_user_id,
      {
        role_ids: formData.role_ids,
        is_super: formData.is_super ? 1 : 0,
      },
      true,
    );
    if (res.code === 1) {
      ElMessage.success(res.msg || '保存成功');
      open.value = false;
      emit('success');
    }
  } finally {
    loading.value = false;
    drawerApi.setState({ confirmLoading: false });
  }
}
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :confirm-loading="loading"
    :destroy-on-close="true"
    title="分配角色"
  >
    <ElForm v-loading="loading" label-width="100px">
      <ElFormItem label="用户名">
        <span>{{ user?.user_name }}</span>
      </ElFormItem>
      <ElFormItem label="角色">
        <ElSelect v-model="formData.role_ids" class="w-full" multiple>
          <ElOption
            v-for="item in roleOptions"
            :key="item.role_id"
            :label="item.role_name"
            :value="item.role_id"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem v-if="user && user.admin_user_id !== 10001" label="超管">
        <ElCheckbox v-model="formData.is_super">设为超级管理员（全权限）</ElCheckbox>
      </ElFormItem>
    </ElForm>
  </Drawer>
</template>
