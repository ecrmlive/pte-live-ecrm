<script setup lang="ts">
import type { AgentGradeItem } from '#/api/core/plus-agent';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import {
  editAgentGradeApi,
  getAgentGradeEditMetaApi,
} from '#/api/core/plus-agent';
import ImageField from '#/components/shop/image-field.vue';

import TabbarColorField from '#/views/native/page/tabbar/tabbar-color-field.vue';

defineOptions({ name: 'AgentGradeEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  grade?: AgentGradeItem;
}>();

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const gradeLevel = ref(2);
const isDefault = ref(0);

const schema = computed((): VbenFormSchema[] => {
  const fields: VbenFormSchema[] = [
    {
      component: 'Input',
      componentProps: { autocomplete: 'off' },
      fieldName: 'name',
      label: '等级名称',
      rules: 'required',
    },
  ];

  if (isDefault.value === 0) {
    fields.push(
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
          label: '选择图片',
        },
        fieldName: 'image',
        label: '背景图片',
        rules: 'required',
      },
      {
        component: markRaw(TabbarColorField),
        componentProps: { defaultColor: '#333333' },
        fieldName: 'font_color',
        label: '文字颜色',
      },
    );

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
      fieldName: 'condition_type',
      label: '升级条件',
    });
  }

  return fields;
});

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 120,
    },
    handleSubmit: async (values) => {
      if (!props.grade) return;
      submitting.value = true;
      try {
        await editAgentGradeApi({
          ...props.grade,
          condition_type: String(values.condition_type ?? props.grade.condition_type),
          first_percent: Number(values.first_percent ?? props.grade.first_percent),
          font_color: String(values.font_color ?? props.grade.font_color),
          image: String(values.image ?? props.grade.image),
          name: String(values.name ?? props.grade.name),
          second_percent: Number(values.second_percent ?? props.grade.second_percent),
          weight: Number(values.weight ?? props.grade.weight),
        });
        ElMessage.success('等级修改成功');
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

function fillForm(grade: AgentGradeItem) {
  isDefault.value = grade.is_default ?? 0;
  void formApi.setValues({ ...grade });
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen && props.grade) {
      fillForm(props.grade);
    }
  },
});

watch(
  () => [open.value, props.grade] as const,
  ([visible, grade]) => {
    if (visible) {
      modalApi.open();
      if (grade) fillForm(grade);
      return;
    }
    modalApi.close();
  },
);

async function loadMeta() {
  const res = await getAgentGradeEditMetaApi();
  gradeLevel.value = Number(res.basicSetting?.level ?? 2);
}

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
    title="编辑等级"
  >
    <Form />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
