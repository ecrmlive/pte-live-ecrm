<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import {
  addAgentGradeApi,
  getAgentGradeAddMetaApi,
} from '#/api/core/plus-agent';
import ImageField from '#/components/shop/image-field.vue';

import TabbarColorField from '#/views/native/page/tabbar/tabbar-color-field.vue';

defineOptions({ name: 'AgentGradeAddDialog' });

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const gradeLevel = ref(2);

const schema = computed((): VbenFormSchema[] => {
  const fields: VbenFormSchema[] = [
    {
      component: 'Input',
      componentProps: { placeholder: '请输入等级名称' },
      fieldName: 'name',
      label: '等级名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入等级权重', type: 'number' },
      description: '权重越大，等级越高',
      fieldName: 'weight',
      label: '等级权重',
      rules: 'required',
    },
    {
      component: markRaw(ImageField),
      componentProps: {
        hint: '建议图片上传尺寸为 654px×300px',
        previewSize: 96,
      },
      fieldName: 'image',
      label: '背景图片',
      rules: 'required',
    },
    {
      component: markRaw(TabbarColorField),
      componentProps: { defaultColor: '#333333' },
      defaultValue: '#333333',
      fieldName: 'font_color',
      label: '文字颜色',
    },
  ];

  if (gradeLevel.value >= 1) {
    fields.push({
      component: 'Input',
      componentProps: { placeholder: '请输入上调比例', type: 'number' },
      description:
        '在默认分佣比例上上调，如果原来是 5%，这里上调 3%，那么一级拿 8% 佣金',
      fieldName: 'first_percent',
      label: '一级上调',
      rules: 'required',
    });
  }

  if (gradeLevel.value >= 2) {
    fields.push({
      component: 'Input',
      componentProps: { placeholder: '请输入上调比例', type: 'number' },
      fieldName: 'second_percent',
      label: '二级上调',
      rules: 'required',
    });
  }

  fields.push({
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '满足所有任务条件', value: 'and' },
        { label: '满足任意任务条件', value: 'or' },
      ],
    },
    defaultValue: 'and',
    fieldName: 'condition_type',
    label: '升级条件',
  });

  return fields;
});

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 120,
    },
    handleSubmit: async (values) => {
      submitting.value = true;
      try {
        await addAgentGradeApi({
          condition_type: String(values.condition_type ?? 'and'),
          first_percent: Number(values.first_percent ?? 0),
          font_color: String(values.font_color ?? '#333333'),
          image: String(values.image ?? ''),
          name: String(values.name ?? ''),
          second_percent: Number(values.second_percent ?? 0),
          weight: Number(values.weight ?? 100),
        });
        ElMessage.success('新增成功');
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

async function loadMeta() {
  const res = await getAgentGradeAddMetaApi();
  gradeLevel.value = Number(res.basicSetting?.level ?? 2);
}

async function resetForm() {
  await loadMeta();
  void formApi.resetForm();
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void resetForm();
    }
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

void loadMeta();
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :destroy-on-close="true"
    class="w-[640px]"
    title="新增等级"
  >
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
