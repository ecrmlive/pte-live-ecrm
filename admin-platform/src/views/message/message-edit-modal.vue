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

import MessageApi from '#/api/core/message';

import type { MessageFormModel } from './types';

const props = defineProps<{
  form?: MessageFormModel;
  open: boolean;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
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

watch(
  () => [props.open, props.form?.message_id] as const,
  ([open]) => {
    if (!open || !props.form) return;
    Object.assign(model, props.form);
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
      await MessageApi.editMessage({ ...model }, true);
      ElMessage.success('恭喜你，修改成功');
      emit('update:open', false);
      emit('success');
    } finally {
      loading.value = false;
    }
  });
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :model-value="open"
    title="编辑消息"
    width="480px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
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
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">确定</ElButton>
    </template>
  </ElDialog>
</template>
