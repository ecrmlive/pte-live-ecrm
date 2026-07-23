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

const props = defineProps<{ open: boolean }>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const loading = ref(false);

const form = reactive({
  message_name: '',
  message_ename: '',
  message_to: 10,
  message_type: 10,
  sort: 100,
  remark: '',
  status: 0,
});

watch(
  () => props.open,
  (open) => {
    if (!open) return;
    Object.assign(form, {
      message_name: '',
      message_ename: '',
      message_to: 10,
      message_type: 10,
      sort: 100,
      remark: '',
      status: 0,
    });
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
      const res = await MessageApi.addMessage({ ...form }, true);
      if (res.code === 1) {
        ElMessage.success('恭喜你，添加成功');
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
    :close-on-press-escape="false"
    :model-value="open"
    title="添加消息"
    width="480px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
  >
    <ElForm ref="formRef" :model="form" label-width="132px">
      <ElFormItem
        :rules="[{ message: ' ', required: true }]"
        label="消息名称"
        prop="message_name"
      >
        <ElInput v-model="form.message_name" autocomplete="off" placeholder="请输入消息名称" />
      </ElFormItem>
      <ElFormItem
        :rules="[{ message: ' ', required: true }]"
        label="名称(英文唯一)"
        prop="message_ename"
      >
        <ElInput v-model="form.message_ename" autocomplete="off" placeholder="请输入消息英文名称" />
      </ElFormItem>
      <ElFormItem label="通知对象">
        <ElSelect v-model="form.message_to" placeholder="请选择通知对象">
          <ElOption :value="10" label="会员" />
          <ElOption :value="20" label="商家" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="消息类别">
        <ElSelect v-model="form.message_type" placeholder="请选择消息类别">
          <ElOption :value="10" label="订单" />
          <ElOption :value="20" label="分销" />
          <ElOption :value="30" label="通知" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="排序">
        <ElInput v-model="form.sort" placeholder="请输入数字" />
      </ElFormItem>
      <ElFormItem label="备注">
        <ElInput v-model="form.remark" autocomplete="off" placeholder="请输入备注" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">确定</ElButton>
    </template>
  </ElDialog>
</template>
