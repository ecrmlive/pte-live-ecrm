<script setup lang="ts">
import type { LiveWxProductForm, LiveWxProductItem } from '#/api/core/plus-live-wx-product';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm, z } from '#/adapter/form';
import { editLiveWxProductApi } from '#/api/core/plus-live-wx-product';

import ImageField from '#/components/shop/image-field.vue';

defineOptions({ name: 'LiveWxProductEditDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveWxProductItem | null;
}>();

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const productImage = ref('');

const nameRule = z
  .string()
  .min(1, { message: '请输入商品名称' })
  .min(3, { message: '长度在 3 到 17 个字符' })
  .max(17, { message: '长度在 3 到 17 个字符' });

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    fieldName: 'name',
    label: '商品名称',
    rules: nameRule,
  },
  {
    component: markRaw(ImageField),
    componentProps: {
      hint: '建议尺寸300*300,大小不超过2M',
    },
    fieldName: 'cover_img',
    label: '商品封面图',
    rules: 'required',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '一口价', value: 1 },
        { label: '价格区间', value: 2 },
        { label: '折扣价', value: 3 },
      ],
    },
    fieldName: 'price_type',
    label: '价格类型',
  },
  {
    component: 'Input',
    dependencies: {
      componentProps: (values) => ({
        placeholder:
          Number(values.price_type) === 1
            ? '请输入价格'
            : Number(values.price_type) === 2
              ? '请输入最低价格'
              : '请输入原价',
      }),
      show: (values) => [1, 2, 3].includes(Number(values.price_type)),
      triggerFields: ['price_type'],
    },
    fieldName: 'price',
    label: '价格',
    rules: 'required',
  },
  {
    component: 'Input',
    dependencies: {
      componentProps: (values) => ({
        placeholder: Number(values.price_type) === 2 ? '请输入最高价格' : '请输入现价',
      }),
      show: (values) => [2, 3].includes(Number(values.price_type)),
      triggerFields: ['price_type'],
    },
    fieldName: 'price2',
    label: '价格',
    rules: z.string().min(1, { message: '请输入价格' }),
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { size: 'small' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      if (!values.cover_img) {
        ElMessage.warning('请上传商品封面图');
        return;
      }
      submitting.value = true;
      try {
        await editLiveWxProductApi({
          cover_img: String(values.cover_img ?? ''),
          name: String(values.name ?? ''),
          price: String(values.price ?? ''),
          price2: String(values.price2 ?? ''),
          price_type: Number(values.price_type ?? 1),
          product_id: String(values.product_id ?? ''),
          shop_supplier_id: Number(values.shop_supplier_id ?? 0),
          wx_product_id: Number(values.wx_product_id ?? 0),
        });
        ElMessage.success('编辑成功');
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

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(
  () => [open.value, props.row] as const,
  ([visible, row]) => {
    if (visible) {
      modalApi.open();
    } else {
      modalApi.close();
    }
    if (!visible || !row) return;
    productImage.value = row.cover_img ?? '';
    void formApi.setValues({
      cover_img: row.cover_img ?? '',
      name: row.name ?? '',
      price: row.price ?? '',
      price2: row.price2 ?? '',
      price_type: Number(row.price_type ?? 1),
      product_id: row.product_id ?? '',
      shop_supplier_id: Number(row.shop_supplier_id ?? 0),
      wx_product_id: row.wx_product_id,
    });
  },
  { immediate: true },
);

async function submit() {
  await formApi.validateAndSubmitForm();
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[520px]"
    title="编辑商品"
  >
    <img
      v-if="productImage"
      :src="productImage"
      alt=""
      class="mb-3 h-[120px] w-[120px] object-cover"
    />
    <Form />
    <template #footer>
      <ElButton size="small" @click="open = false">取消</ElButton>
      <ElButton :loading="submitting" size="small" type="primary" @click="submit">
        立即提交
      </ElButton>
    </template>
  </Modal>
</template>
