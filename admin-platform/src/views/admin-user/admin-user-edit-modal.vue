<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import {
  ElButton,
  ElCheckbox,
  ElDialog,
  ElForm,
  ElFormItem,
  ElMessage,
  ElOption,
  ElSelect,
} from 'element-plus';

import PlatformAdminUserApi from '#/api/core/platform-admin-user';
import type { PlatformAdminUserRow } from '#/api/core/platform-admin-user';

const props = defineProps<{
  open: boolean;
  roleOptions: { role_id: number; role_name: string }[];
  user?: PlatformAdminUserRow;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const loading = ref(false);
const formData = reactive({
  is_super: false,
  role_ids: [] as number[],
});

watch(
  () => [props.open, props.user?.admin_user_id] as const,
  async ([open]) => {
    if (!open || !props.user) return;
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
  },
  { immediate: true },
);

function handleClose() {
  emit('update:open', false);
}

async function handleSubmit() {
  if (!props.user) return;
  loading.value = true;
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
      emit('update:open', false);
      emit('success');
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :model-value="open"
    title="分配角色"
    width="520px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
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
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">保存</ElButton>
    </template>
  </ElDialog>
</template>
