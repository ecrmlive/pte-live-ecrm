<script setup lang="ts">
import type { ShopLinkPickerItem } from '#/api/core/shop-link';
import type { ShopLinkValue } from '#/types/shop-link';
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, reactive } from 'vue';

import { useVbenForm } from '#/adapter/form';

const props = defineProps<{
  linkData?: ShopLinkValue | null;
}>();

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

function parseCombinedString(str: string) {
  const regex = /([^=,]+)=([^,]*)(?:,|$)/g;
  const params: Record<string, string> = {};
  let match: RegExpExecArray | null;
  while ((match = regex.exec(str)) !== null) {
    params[match[1]!] = match[2] ?? '';
  }
  return params;
}

function emitValue(values: Record<string, unknown>) {
  const targetAppId = String(values.targetAppId ?? '');
  const targetGhId = String(values.targetGhId ?? '');
  const path = String(values.path ?? '');
  const url = `targetAppId=${targetAppId},targetGhId=${targetGhId},path=${path}`;
  const name = `appId:${targetAppId},原始Id:${targetGhId},路径:${path}`;
  emit('change', { name, type: 'mp', url });
}

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    fieldName: 'targetAppId',
    label: 'appId',
  },
  {
    component: 'Input',
    fieldName: 'targetGhId',
    label: '原始Id',
  },
  {
    component: 'Input',
    fieldName: 'path',
    label: '页面路径',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      labelWidth: 80,
    },
    handleValuesChange: (values) => {
      emitValue(values);
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

onMounted(() => {
  let targetAppId = '';
  let targetGhId = '';
  let path = '';
  if (props.linkData?.linkeType === 'mp') {
    const raw = String(props.linkData.linkUrl ?? props.linkData.url ?? '');
    const params = parseCombinedString(raw);
    targetAppId = params.targetAppId ?? '';
    targetGhId = params.targetGhId ?? '';
    path = params.path ?? '';
  }
  void formApi.setValues({ path, targetAppId, targetGhId });
  emitValue({ path, targetAppId, targetGhId });
});
</script>

<template>
  <Form />
  <div class="mt-2 space-y-1 text-xs leading-6 text-gray-500">
    <p>小程序内跳转请填写 appId；APP 内请填写原始 Id，H5 暂不支持。</p>
    <p>同小程序最多链接 10 个其他小程序；不可填写自身 appId。</p>
    <p>路径留空默认跳转小程序首页。</p>
  </div>
</template>
