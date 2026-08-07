<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, h, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import {
  editAgentProductApi,
  getAgentProductEditMetaApi,
} from '#/api/core/plus-agent';

defineOptions({ name: 'AgentProductDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  initialIsAgent?: number;
  productId?: number;
  productName?: string;
}>();

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const submitting = ref(false);
const gradeLevel = ref(2);
const unit = ref('%');

function changeMoneyType(value: number) {
  unit.value = value === 10 ? '%' : '元';
}

const schema = computed((): VbenFormSchema[] => {
  const fields: VbenFormSchema[] = [
    {
      component: 'Input',
      componentProps: { disabled: true },
      defaultValue: props.productName ?? '',
      fieldName: 'productName',
      label: '商品名称',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '参与', value: 1 },
          { label: '不参与', value: 0 },
        ],
      },
      fieldName: 'is_agent',
      label: '是否参与推广',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '开启', value: 1 },
          { label: '关闭', value: 0 },
        ],
      },
      dependencies: {
        show: (values) => Number(values.is_agent) === 1,
        triggerFields: ['is_agent'],
      },
      fieldName: 'is_ind_agent',
      label: '是否开启单独分销',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '百分比', value: 10 },
          { label: '固定金额', value: 20 },
        ],
      },
      dependencies: {
        show: (values) =>
          Number(values.is_agent) === 1 && Number(values.is_ind_agent) === 1,
        triggerFields: ['is_agent', 'is_ind_agent'],
      },
      fieldName: 'agent_money_type',
      label: '分销佣金类型',
    },
    {
      component: 'Input',
      componentProps: { class: 'w-full', min: 0, type: 'number' },
      dependencies: {
        show: (values) =>
          Number(values.is_agent) === 1 && Number(values.is_ind_agent) === 1,
        triggerFields: ['is_agent', 'is_ind_agent'],
      },
      fieldName: 'first_money',
      label: '单独分销设置',
      renderComponentContent: () => ({
        append: () => h('span', unit.value),
        prepend: () => h('span', '一级佣金：'),
      }),
    },
  ];

  if (gradeLevel.value >= 2) {
    fields.push({
      component: 'Input',
      componentProps: { class: 'w-full', min: 0, type: 'number' },
      dependencies: {
        show: (values) =>
          Number(values.is_agent) === 1 && Number(values.is_ind_agent) === 1,
        triggerFields: ['is_agent', 'is_ind_agent'],
      },
      fieldName: 'second_money',
      label: ' ',
      renderComponentContent: () => ({
        append: () => h('span', unit.value),
        prepend: () => h('span', '二级佣金：'),
      }),
    });
  }

  if (gradeLevel.value >= 3) {
    fields.push({
      component: 'Input',
      componentProps: { class: 'w-full', min: 0, type: 'number' },
      dependencies: {
        show: (values) =>
          Number(values.is_agent) === 1 && Number(values.is_ind_agent) === 1,
        triggerFields: ['is_agent', 'is_ind_agent'],
      },
      fieldName: 'third_money',
      label: ' ',
      renderComponentContent: () => ({
        append: () => h('span', unit.value),
        prepend: () => h('span', '三级佣金：'),
      }),
    });
  }

  return fields;
});

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 140,
    },
    handleSubmit: async (values) => {
      if (!props.productId) return;
      submitting.value = true;
      try {
        await editAgentProductApi({
          agent_money_type: Number(values.agent_money_type ?? 10),
          first_money: String(values.first_money ?? ''),
          is_agent: Number(values.is_agent ?? 0),
          is_ind_agent: Number(values.is_ind_agent ?? 0),
          product_id: props.productId,
          second_money: String(values.second_money ?? ''),
          third_money: String(values.third_money ?? ''),
        });
        ElMessage.success('保存成功');
        open.value = false;
        emit('success');
      } finally {
        submitting.value = false;
      }
    },
    handleValuesChange: (values, fieldsChanged) => {
      if (fieldsChanged.includes('agent_money_type')) {
        changeMoneyType(Number(values.agent_money_type ?? 10));
      }
      if (fieldsChanged.includes('productName') && props.productName) {
        void formApi.setFieldValue('productName', props.productName);
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

async function loadFormData() {
  if (!props.productId) return;
  loading.value = true;
  try {
    let isAgent = Number(props.initialIsAgent ?? 0);
    const res = await getAgentProductEditMetaApi(props.productId);
    gradeLevel.value = Number(res.basicSetting?.level ?? 2);
    const agentProduct = res.agent_product;
    if (agentProduct) {
      isAgent = Number(agentProduct.is_agent ?? 0);
      void formApi.setValues({
        agent_money_type: agentProduct.agent_money_type ?? 10,
        first_money: agentProduct.first_money ?? '',
        is_agent: isAgent,
        is_ind_agent: agentProduct.is_ind_agent ?? 0,
        productName: props.productName ?? '',
        second_money: agentProduct.second_money ?? '',
        third_money: agentProduct.third_money ?? '',
      });
      changeMoneyType(Number(agentProduct.agent_money_type ?? 10));
    } else {
      void formApi.setValues({
        agent_money_type: 10,
        first_money: '',
        is_agent: isAgent,
        is_ind_agent: 0,
        productName: props.productName ?? '',
        second_money: '',
        third_money: '',
      });
      changeMoneyType(10);
    }
  } finally {
    loading.value = false;
  }
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void loadFormData();
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
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :destroy-on-close="true"
    class="w-[640px]"
    title="商品规则"
  >
    <Form v-loading="loading" />

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
    </template>
  </Modal>
</template>
