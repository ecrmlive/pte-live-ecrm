<script lang="ts" setup>
import { reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus';

import { editPlatformPasswordApi } from '#/api/core/user';
import { validateBEndPassword } from '#/utils/b-end-password';

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);

const form = reactive({
  oldpass: '',
  pass: '',
  checkPass: '',
});

const validateOldPass = (_: unknown, value: string, callback: (err?: Error) => void) => {
  if (!String(value ?? '').trim()) {
    callback(new Error('请输入原密码'));
    return;
  }
  callback();
};

const validatePass = (_: unknown, value: string, callback: (err?: Error) => void) => {
  const msg = validateBEndPassword(value);
  if (msg) {
    callback(new Error(msg));
    return;
  }
  if (form.checkPass) {
    formRef.value?.validateField('checkPass');
  }
  callback();
};

const validatePass2 = (_: unknown, value: string, callback: (err?: Error) => void) => {
  if (!value) {
    callback(new Error('请再次输入密码'));
  } else if (value !== form.pass) {
    callback(new Error('两次输入密码不一致!'));
  } else {
    callback();
  }
};

const rules = {
  oldpass: [{ required: true, trigger: 'blur', validator: validateOldPass }],
  pass: [{ required: true, trigger: 'blur', validator: validatePass }],
  checkPass: [{ required: true, trigger: 'blur', validator: validatePass2 }],
};

async function submitForm() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    loading.value = true;
    try {
      const res = await editPlatformPasswordApi({ ...form }, true);
      form.oldpass = '';
      form.pass = '';
      form.checkPass = '';
      ElMessage.success(res.msg || '修改成功');
    } finally {
      loading.value = false;
    }
  });
}
</script>

<template>
  <Page auto-content-height>
    <ElForm
      ref="formRef"
      :model="form"
      :rules="rules"
      class="max-w-xl"
      label-width="160px"
    >
      <ElFormItem label="原密码" prop="oldpass">
        <ElInput
          v-model="form.oldpass"
          autocomplete="off"
          class="max-w-[460px]"
          type="password"
        />
      </ElFormItem>
      <ElFormItem label="新密码" prop="pass">
        <ElInput v-model="form.pass" autocomplete="off" class="max-w-[460px]" type="password" />
      </ElFormItem>
      <ElFormItem label="确认新密码" prop="checkPass">
        <ElInput
          v-model="form.checkPass"
          autocomplete="off"
          class="max-w-[460px]"
          type="password"
        />
      </ElFormItem>
      <ElFormItem>
        <ElButton :loading="loading" type="primary" @click="submitForm">保存</ElButton>
      </ElFormItem>
    </ElForm>
  </Page>
</template>
