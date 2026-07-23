<script setup lang="ts">
import type { LiveAnchorForm, LiveAnchorListItem } from '#/api/core/live';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm, z } from '#/adapter/form';
import { updateLiveAnchorApi } from '#/api/core/live';

import ImageField from '#/components/shop/image-field.vue';

defineOptions({ name: 'LiveAnchorEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveAnchorListItem | null;
}>();

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const anchorId = ref(0);

const phoneRule = z.string().superRefine((value, ctx) => {
  if (!/^1\d{10}$/.test(String(value ?? '').trim())) {
    ctx.addIssue({ code: 'custom', message: '请输入11位有效手机号' });
  }
});

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { disabled: true },
    description: '登录账号与手机号一致',
    fieldName: 'account',
    label: '登录账号',
  },
  {
    component: 'Input',
    componentProps: { maxlength: 11 },
    fieldName: 'phone',
    label: '手机号',
    rules: phoneRule,
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: '不修改请留空',
      showPassword: true,
      type: 'password',
    },
    fieldName: 'password',
    label: '新密码',
  },
  {
    component: 'Input',
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
      if (!anchorId.value) return;
      submitting.value = true;
      try {
        const payload: LiveAnchorForm & { anchor_id: number } = {
          anchor_id: anchorId.value,
          avatar: String(values.avatar ?? ''),
          intro: String(values.intro ?? ''),
          nick_name: String(values.nick_name ?? ''),
          phone: String(values.phone ?? ''),
          sort: Number(values.sort ?? 0),
          status: Number(values.status ?? 1),
          wechat: String(values.wechat ?? ''),
        };
        if (values.password) {
          payload.password = String(values.password);
        }
        const res = await updateLiveAnchorApi(payload);
        ElMessage.success(res.msg || '更新成功');
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
  },
});

function fillForm(row: LiveAnchorListItem) {
  anchorId.value = row.anchor_id;
  void formApi.setValues({
    account: row.account,
    avatar: row.avatar || '',
    intro: row.intro || '',
    nick_name: row.nick_name,
    password: '',
    phone: row.phone,
    sort: row.sort ?? 0,
    status: row.status ?? 1,
    wechat: row.wechat || '',
  });
}

watch(
  () => [open.value, props.row] as const,
  ([visible, row]) => {
    if (visible) {
      modalApi.open();
      if (row) fillForm(row);
      return;
    }
    modalApi.close();
  },
);

async function submit() {
  await formApi.validateAndSubmitForm();
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[640px]"
    title="编辑主播"
  >
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
