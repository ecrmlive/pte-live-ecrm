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

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      placeholder: 'qxkejiwl.top/path',
    },
    fieldName: 'name',
    label: 'https://',
  },
  {
    component: 'Divider',
    fieldName: '_tips',
    formItemClass: 'col-span-full',
    hideLabel: true,
    renderComponentContent: () => ({
      default: () => null,
    }),
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      labelWidth: 80,
    },
    handleValuesChange: (values) => {
      const name = String(values.name ?? '');
      emit('change', {
        name,
        type: 'h5',
        url: `https://${name}`,
      });
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

onMounted(() => {
  let name = '';
  if (props.linkData?.linkeType === 'h5') {
    name = String(props.linkData.name ?? props.linkData.url ?? '').replace(/^https:\/\//, '');
  }
  void formApi.setValues({ name });
  emit('change', {
    name,
    type: 'h5',
    url: `https://${name}`,
  });
});
</script>

<template>
  <Form />
  <div class="mt-2 space-y-1 text-xs text-gray-500">
    <p>注意一：链接的网页必须支持 SSL，且以 https:// 开头</p>
    <p>注意二：链接的域名必须是备案域名</p>
    <p>
      注意三：域名须加入微信小程序业务域名，
      <a
        class="text-primary"
        href="https://developers.weixin.qq.com/miniprogram/dev/framework/ability/domain.html"
        rel="noopener noreferrer"
        target="_blank"
      >
        查看详情
      </a>
    </p>
  </div>
</template>
