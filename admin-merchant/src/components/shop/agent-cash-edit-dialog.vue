<script setup lang="ts">
import type { AgentCashItem } from '#/api/core/plus-agent';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { submitAgentCashApi } from '#/api/core/plus-agent';

defineOptions({ name: 'AgentCashEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  cash?: AgentCashItem;
}>();

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const rejectReason = ref('');
const isRejectedView = ref(false);
const cashId = ref(0);

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '审核通过', value: '20' },
        { label: '驳回', value: '30' },
      ],
    },
    fieldName: 'apply_status',
    label: '审核状态',
  },
  {
    component: 'Input',
    componentProps: { class: 'w-full', type: 'textarea', rows: 3 },
    dependencies: {
      show: (values) => values.apply_status === '30',
      triggerFields: ['apply_status'],
    },
    fieldName: 'reject_reason',
    label: '驳回原因',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      submitting.value = true;
      try {
        await submitAgentCashApi({
          apply_status: Number(values.apply_status ?? 20),
          id: cashId.value,
          reject_reason: String(values.reject_reason ?? ''),
        });
        ElMessage.success('修改成功');
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

function fillFromCash(cash: AgentCashItem) {
  const status = cash.apply_status?.value ?? 0;
  isRejectedView.value = status === 30;
  rejectReason.value = cash.reject_reason ?? '';
  cashId.value = cash.id;
  if (!isRejectedView.value) {
    void formApi.setValues({
      apply_status: '20',
      reject_reason: '',
    });
  }
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen && props.cash) {
      fillFromCash(props.cash);
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    if (props.cash) fillFromCash(props.cash);
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
    :title="isRejectedView ? '驳回原因' : '提现审核'"
    class="w-[560px]"
  >
    <Form v-if="!isRejectedView" />
    <p v-else class="reject-reason">{{ rejectReason || '-' }}</p>

    <template #footer>
      <template v-if="!isRejectedView">
        <ElButton @click="open = false">取消</ElButton>
        <ElButton :loading="submitting" type="primary" @click="submit">确定</ElButton>
      </template>
      <ElButton v-else type="primary" @click="open = false">关闭</ElButton>
    </template>
  </Modal>
</template>

<style scoped>
.reject-reason {
  margin: 0;
  line-height: 1.6;
  color: hsl(var(--foreground));
  white-space: pre-wrap;
}
</style>
