<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus';

import FileApi from '#/api/core/file';

import type { PictureCategory } from './types';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  category: null | PictureCategory;
  fileType: string;
}>();

const emit = defineEmits<{
  success: [];
}>();

const formRef = ref<InstanceType<typeof ElForm>>();
const submitting = ref(false);
const title = ref('新增分类');

const form = reactive({
  category_id: 0,
  name: '',
});

function resetForm() {
  if (props.category?.category_id != null) {
    form.name = props.category.name;
    form.category_id = props.category.category_id;
    title.value = '编辑分类';
  } else {
    form.name = '';
    form.category_id = 0;
    title.value = '新增分类';
  }
}

const [Modal, modalApi] = useVbenDrawer({
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
      modalApi.open();
      return;
    }
    modalApi.close();
  },
  { immediate: true },
);

watch(
  () => props.category?.category_id,
  () => {
    if (open.value) {
      resetForm();
    }
  },
);

async function handleSubmit() {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    submitting.value = true;
    modalApi.setState({ confirmLoading: true });
    try {
      if (props.category?.category_id != null) {
        const res = await FileApi.editCategory({
          category_id: form.category_id,
          name: form.name,
        });
        if ((res as { code?: number }).code === 1 || res) {
          ElMessage.success('修改成功');
          open.value = false;
          emit('success');
        }
      } else {
        const res = await FileApi.addCategory({
          groupType: props.fileType,
          name: form.name,
        });
        if ((res as { code?: number }).code === 1 || res) {
          ElMessage.success('新增成功');
          open.value = false;
          emit('success');
        }
      }
    } catch {
      ElMessage.error('操作失败');
    } finally {
      submitting.value = false;
      modalApi.setState({ confirmLoading: false });
    }
  });
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :confirm-loading="submitting"
    :destroy-on-close="true"
    :title="title"
    class="w-[420px]"
  >
    <ElForm ref="formRef" :model="form" label-width="88px" size="small">
      <ElFormItem
        :rules="[{ message: '名字不能为空', required: true }]"
        label="分类名称"
        prop="name"
      >
        <ElInput v-model="form.name" autocomplete="off" />
      </ElFormItem>
    </ElForm>
  </Modal>
</template>
