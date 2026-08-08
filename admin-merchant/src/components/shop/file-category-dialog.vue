<script setup lang="ts">
import type { ShopFileGroupItem } from '#/api/core/file';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import {
  addShopFileGroupApi,
  editShopFileGroupApi,
} from '#/api/core/file';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  category?: ShopFileGroupItem | null;
  fileType: 'image' | 'video';
}>();

const emit = defineEmits<{
  success: [];
}>();

const saving = ref(false);
const groupId = ref<number | undefined>();

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { maxlength: 32 },
    fieldName: 'categoryname',
    label: '分类名称',
    rules: 'required',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 88,
    },
    handleSubmit: async (values) => {
      const name = String(values.categoryname ?? '').trim();
      if (!name) {
        ElMessage.warning('请输入分类名称');
        return;
      }
      saving.value = true;
      try {
        if (groupId.value != null) {
          await editShopFileGroupApi({
            group_id: groupId.value,
            group_name: name,
          });
          ElMessage.success('修改成功');
        } else {
          await addShopFileGroupApi({
            group_name: name,
            group_type: props.fileType,
          });
          ElMessage.success('新增成功');
        }
        open.value = false;
        emit('success');
      } finally {
        saving.value = false;
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (!isOpen) return;
    groupId.value =
      props.category?.group_id != null ? props.category.group_id : undefined;
    void formApi.setValues({
      categoryname: props.category?.group_name ?? '',
    });
  },
});

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
    :title="groupId != null ? '编辑分类' : '新增分类'"
    :z-index="2200"
    class="w-[420px]"
  >
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="saving" type="primary" @click="submit">提交</ElButton>
    </template>
  </Modal>
</template>
