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

import MessageApi from '#/api/core/message';

import type { MessageFormModel } from './types';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  form?: MessageFormModel;
}>();

const emit = defineEmits<{
  success: [];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);

const model = reactive<MessageFormModel>({
  message_name: '',
  message_ename: '',
  message_to: 10,
  message_type: 10,
  sort: 100,
  remark: '',
});

function syncForm() {
  if (!props.form) return;
  Object.assign(model, props.form);
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
      drawerApi.setState({ title: '编辑消息' }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(
  () => props.form?.message_id,
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
      await MessageApi.editMessage({ ...model }, true);
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
    title="编辑消息"
  >
    <ElForm ref="formRef" :model="model" label-width="132px">
      <ElFormItem
        :rules="[{ message: ' ', required: true }]"
        label="消息名称"
        prop="message_name"
      >
        <ElInput v-model="model.message_name" autocomplete="off" placeholder="请输入消息名称" />
      </ElFormItem>
      <ElFormItem
        :rules="[{ message: ' ', required: true }]"
        label="名称(英文唯一)"
        prop="message_ename"
      >
        <ElInput v-model="model.message_ename" autocomplete="off" placeholder="请输入消息英文名称" />
      </ElFormItem>
      <ElFormItem label="通知对象">
        <ElSelect v-model="model.message_to" placeholder="请选择通知对象">
          <ElOption :value="10" label="会员" />
          <ElOption :value="20" label="商家" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="消息类别">
        <ElSelect v-model="model.message_type" placeholder="请选择消息类别">
          <ElOption :value="10" label="订单" />
          <ElOption :value="20" label="分销" />
          <ElOption :value="30" label="通知" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="排序">
        <ElInput v-model="model.sort" placeholder="请输入数字" />
      </ElFormItem>
      <ElFormItem label="备注">
        <ElInput v-model="model.remark" autocomplete="off" placeholder="请输入备注" />
      </ElFormItem>
    </ElForm>
  </Drawer>
</template>
