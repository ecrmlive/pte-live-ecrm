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

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  success: [];
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

function resetForm() {
  Object.assign(form, {
    message_name: '',
    message_ename: '',
    message_to: 10,
    message_type: 10,
    sort: 100,
    remark: '',
    status: 0,
  });
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
      drawerApi.setState({ title: '添加消息' }).open();
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
      const res = await MessageApi.addMessage({ ...form }, true);
      if (res.code === 1) {
        ElMessage.success('恭喜你，添加成功');
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
    title="添加消息"
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
  </Drawer>
</template>
