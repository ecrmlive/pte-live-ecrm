<script setup lang="ts">
import type { LiveAnchorForm } from '#/api/core/live';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm, z } from '#/adapter/form';
import { createLiveAnchorApi } from '#/api/core/live';

import ImageField from '#/components/shop/image-field.vue';

defineOptions({ name: 'LiveAnchorAddDialog' });

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);

const phoneRule = z.string().superRefine((value, ctx) => {
  if (!/^1\d{10}$/.test(String(value ?? '').trim())) {
    ctx.addIssue({ code: 'custom', message: '请输入11位有效手机号' });
  }
});

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      autocomplete: 'off',
      maxlength: 11,
      name: 'anchor_phone',
      placeholder: '作为登录账号，全平台唯一',
      type: 'tel',
    },
    description: '登录账号与手机号相同，无需单独填写',
    fieldName: 'phone',
    label: '手机号',
    rules: phoneRule,
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: '初始密码',
      showPassword: true,
      type: 'password',
    },
    fieldName: 'password',
    label: '登录密码',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { placeholder: '展示昵称' },
    fieldName: 'nick_name',
    label: '昵称',
    rules: 'required',
  },
  {
    component: markRaw(ImageField),
    componentProps: { hint: '建议上传正方形头像图片' },
    fieldName: 'avatar',
    label: '头像',
  },
  {
    component: 'Input',
    fieldName: 'wechat',
    label: '微信号',
  },
  {
    component: 'Textarea',
    componentProps: { rows: 3 },
    fieldName: 'intro',
    label: '简介',
  },
  {
    component: 'InputNumber',
    componentProps: { max: 9999, min: 0 },
    defaultValue: 0,
    fieldName: 'sort',
    label: '排序',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '启用', value: 1 },
        { label: '禁用', value: 0 },
      ],
    },
    defaultValue: 1,
    fieldName: 'status',
    label: '状态',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      submitting.value = true;
      try {
        const phone = String(values.phone ?? '').trim();
        const res = await createLiveAnchorApi({
          account: phone,
          avatar: String(values.avatar ?? ''),
          intro: String(values.intro ?? ''),
          nick_name: String(values.nick_name ?? ''),
          password: String(values.password ?? ''),
          phone,
          sort: Number(values.sort ?? 0),
          status: Number(values.status ?? 1),
          wechat: String(values.wechat ?? ''),
        });
        ElMessage.success(res.msg || '添加成功');
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

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void resetAddForm();
    }
  },
});

async function resetAddForm() {
  await formApi.setValues({
    avatar: '',
    intro: '',
    nick_name: '',
    password: '',
    phone: '',
    sort: 0,
    status: 1,
    wechat: '',
  });
  await formApi.resetValidate();
}

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
    class="w-[640px]"
    title="添加主播"
  >
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
