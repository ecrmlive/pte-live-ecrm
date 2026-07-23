<script setup lang="ts">
import type { ProductListItem } from '#/api/core/product';
import type { LiveWxProductForm } from '#/api/core/plus-live-wx-product';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, defineComponent, h, markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm, z } from '#/adapter/form';
import { addLiveWxProductApi } from '#/api/core/plus-live-wx-product';

import ImageField from '#/components/shop/image-field.vue';
import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';

defineOptions({ name: 'LiveWxProductAddDialog' });

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  success: [];
}>();

const submitting = ref(false);
const productPickerOpen = ref(false);
const productName = ref('');
const productImage = ref('');

const nameRule = z
  .string()
  .min(1, { message: '请输入商品名称' })
  .min(3, { message: '长度在 3 到 17 个字符' })
  .max(17, { message: '长度在 3 到 17 个字符' });

const ProductPickField = defineComponent({
  name: 'ProductPickField',
  props: {
    modelValue: {
      default: '',
      type: [String, Number],
    },
  },
  emits: ['update:modelValue'],
  setup() {
    return () =>
      h('div', { class: 'flex flex-col gap-2' }, [
        h(
          'button',
          {
            class: 'inline-flex h-8 items-center rounded border px-3 text-sm hover:bg-muted',
            type: 'button',
            onClick: () => {
              productPickerOpen.value = true;
            },
          },
          '选择商品',
        ),
        productName.value ? h('span', productName.value) : null,
        productImage.value
          ? h('img', {
              alt: '',
              class: 'h-[120px] w-[120px] object-cover',
              src: productImage.value,
            })
          : null,
      ]);
  },
});

function onProductPicked(row: ProductListItem) {
  productName.value = row.product_name;
  productImage.value = row.image?.file_path ?? '';
  void formApi.setValues({
    cover_img: productImage.value || undefined,
    name: row.product_name,
    product_id: row.product_id,
    shop_supplier_id: Number(row.shop_supplier_id ?? 0),
  });
}

const schema = computed((): VbenFormSchema[] => [
  {
    component: markRaw(ProductPickField),
    fieldName: 'product_id',
    label: '选择商品',
    rules: 'selectRequired',
  },
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
    defaultValue: 1,
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
      if (!values.product_id) {
        ElMessage.warning('请绑定商品');
        return;
      }
      if (!values.cover_img) {
        ElMessage.warning('请上传商品封面图');
        return;
      }
      submitting.value = true;
      try {
        await addLiveWxProductApi({
          cover_img: String(values.cover_img ?? ''),
          name: String(values.name ?? ''),
          price: String(values.price ?? ''),
          price2: String(values.price2 ?? ''),
          price_type: Number(values.price_type ?? 1),
          product_id: values.product_id as number | string,
          shop_supplier_id: Number(values.shop_supplier_id ?? 0),
        });
        ElMessage.success('创建成功');
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

function resetForm() {
  productName.value = '';
  productImage.value = '';
  void formApi.resetForm();
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    resetForm();
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
    class="w-[520px]"
    title="添加商品"
  >
    <Form />
    <template #footer>
      <ElButton size="small" @click="open = false">取消</ElButton>
      <ElButton size="small" @click="resetForm">重置</ElButton>
      <ElButton :loading="submitting" size="small" type="primary" @click="submit">
        立即提交
      </ElButton>
    </template>

    <ProductPickerDialog v-model:open="productPickerOpen" @select="onProductPicked" />
  </Modal>
</template>
