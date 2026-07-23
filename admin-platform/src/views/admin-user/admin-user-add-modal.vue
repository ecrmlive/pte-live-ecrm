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

import PlatformAdminUserApi from '#/api/core/platform-admin-user';

const props = defineProps<{
  open: boolean;
  roleOptions: { role_id: number; role_name: string }[];
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);

const form = reactive({
  password: '',
  role_ids: [] as number[],
  user_name: '',
});

watch(
  () => props.open,
  (open) => {
    if (!open) return;
    form.user_name = '';
    form.password = '';
    form.role_ids = [];
  },
);

function handleClose() {
  emit('update:open', false);
}

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
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
        ElMessage.success(res.msg || '添加成功');
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
    title="新增平台账号"
    width="520px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
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
            v-for="item in roleOptions"
            :key="item.role_id"
            :label="item.role_name"
            :value="item.role_id"
          />
        </ElSelect>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">保存</ElButton>
    </template>
  </ElDialog>
</template>
