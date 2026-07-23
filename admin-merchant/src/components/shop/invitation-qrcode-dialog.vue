<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { getInvitationQrcodeApi } from '#/api/core/plus-invitation';

defineOptions({ name: 'InvitationQrcodeDialog' });

const props = defineProps<{
  invitationGiftId: number;
  open: boolean;
}>();

const emit = defineEmits<{
  'update:open': [value: boolean];
}>();

const visible = ref(false);
const loading = ref(false);
const imageUrl = ref('');

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '微信小程序', value: 'wx' },
        { label: '公众号，H5网页', value: 'mp' },
      ],
    },
    fieldName: 'source',
    label: '下载类型',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { size: 'small' },
      labelWidth: 120,
    },
    handleValuesChange: (values) => {
      const source = values.source as 'mp' | 'wx' | undefined;
      if (source) {
        void loadQrcode(source);
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    visible.value = isOpen;
    if (!isOpen) {
      emit('update:open', false);
    }
  },
});

watch(
  () => props.open,
  (value) => {
    visible.value = value;
    if (value) {
      void formApi.setValues({ source: 'wx' });
      void loadQrcode('wx');
      modalApi.open();
    } else {
      modalApi.close();
    }
  },
  { immediate: true },
);

async function loadQrcode(source: 'mp' | 'wx') {
  loading.value = true;
  imageUrl.value = '';
  try {
    const res = await getInvitationQrcodeApi({
      id: props.invitationGiftId,
      source,
    });
    imageUrl.value = res.image || '';
  } catch {
    ElMessage.error('二维码加载失败');
  } finally {
    loading.value = false;
  }
}

function closeDialog() {
  visible.value = false;
  emit('update:open', false);
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[640px]"
    title="推广码"
  >
    <div v-loading="loading">
      <Form />
      <div v-if="imageUrl" class="ml-[120px] flex flex-col gap-2">
        <img :src="imageUrl" alt="推广二维码" class="h-[120px] w-[120px] object-contain" />
        <a :href="imageUrl" download rel="external nofollow">下载二维码</a>
      </div>
    </div>
    <template #footer>
      <ElButton @click="closeDialog">关闭</ElButton>
    </template>
  </Modal>
</template>
