<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';
import type { ProductSpecFormData } from '#/utils/product-spec';

import { useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElInput,
  ElLoading,
  ElMessage,
  ElUpload,
} from 'element-plus';
import { computed, defineComponent, h, markRaw, reactive, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { importProductVirtualApi } from '#/api/core/product';

defineOptions({ name: 'ProductCarmineDialog' });

const open = defineModel<boolean>('open', { default: true });

const props = defineProps<{
  formData: ProductSpecFormData;
  scene?: string;
}>();

const emit = defineEmits<{
  close: [payload: { params: ProductSpecFormData; type: 'error' | 'success' }];
}>();

const params = reactive<ProductSpecFormData>({
  card_type: 10,
  virtualInfo: [],
});

const VirtualCardListField = defineComponent({
  name: 'VirtualCardListField',
  props: {
    modelValue: { default: () => [], type: Array as () => Array<{ card_no: string; card_pwd: string }> },
  },
  emits: ['update:modelValue'],
  setup(fieldProps, { emit }) {
    function updateList(list: Array<{ card_no: string; card_pwd: string }>) {
      emit('update:modelValue', list);
    }
    function removeList(index: number) {
      const list = [...(fieldProps.modelValue ?? [])];
      list.splice(index, 1);
      updateList(list);
    }
    function addList() {
      updateList([...(fieldProps.modelValue ?? []), { card_no: '', card_pwd: '' }]);
    }
    async function uploadExcel(options: { file: File }) {
      const loading = ElLoading.service({
        background: 'rgba(0, 0, 0, 0.7)',
        lock: true,
        text: '正在处理,请等待',
      });
      try {
        const res = await importProductVirtualApi(options.file);
        updateList([...(fieldProps.modelValue ?? []), ...(res.list ?? [])]);
        ElMessage.warning('导入成功');
      } catch {
        ElMessage.warning('本次处理失败');
      } finally {
        loading.close();
      }
    }
    function onBeforeUpload(file: File) {
      const ext = file.name.slice(file.name.lastIndexOf('.') + 1);
      const isExcel = ext === 'xlsx';
      const isLt10M = file.size / 1024 / 1024 < 10;
      if (!isExcel) {
        ElMessage.error('上传文件只能是 excel 格式!');
      }
      if (!isLt10M) {
        ElMessage.error('上传文件大小不能超过 10MB!');
      }
      return isExcel && isLt10M;
    }
    function downCard() {
      let baseUrl = import.meta.env.VITE_BASIC_URL as string | undefined;
      if (!baseUrl) {
        baseUrl = `${window.location.protocol}//${window.location.host}`;
      }
      window.location.href = `${baseUrl}/card.xlsx`;
    }
    return () =>
      h('div', [
        h(
          'div',
          { class: 'virtual-info' },
          (fieldProps.modelValue ?? []).map((item, index) =>
            h('div', { class: 'virtual-info__row', key: index }, [
              h('div', { class: 'virtual-info__field' }, [
                h('span', { class: 'virtual-info__label' }, `卡号${index + 1}：`),
                h(ElInput, {
                  class: 'max-w460',
                  modelValue: item.card_no,
                  'onUpdate:modelValue': (v: string) => {
                    item.card_no = v;
                  },
                }),
              ]),
              h('div', { class: 'virtual-info__field' }, [
                h('span', { class: 'virtual-info__label' }, `密码${index + 1}：`),
                h(ElInput, {
                  class: 'max-w460',
                  modelValue: item.card_pwd,
                  'onUpdate:modelValue': (v: string) => {
                    item.card_pwd = v;
                  },
                }),
              ]),
              h(ElButton, { link: true, type: 'primary', onClick: () => removeList(index) }, () => '删除'),
            ]),
          ),
        ),
        h('div', { class: 'virtual-info__actions' }, [
          h(ElButton, { size: 'small', type: 'primary', onClick: addList }, () => '添加卡密'),
          h(
            ElUpload,
            {
              accept: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
              action: '',
              beforeUpload: onBeforeUpload,
              httpRequest: uploadExcel,
              showFileList: false,
            },
            () => h(ElButton, { size: 'small', type: 'success' }, () => '导入卡密'),
          ),
          h(ElButton, { size: 'small', onClick: downCard }, () => '下载卡密模板'),
        ]),
      ]);
  },
});

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'RadioGroup',
    componentProps: {
      disabled: props.scene === 'edit',
      options: [
        { label: '固定卡密', value: 10 },
        { label: '一次性卡密', value: 20 },
      ],
      onChange: (value: number) => {
        if (
          value === 20 &&
          (!params.virtualInfo || params.virtualInfo === '')
        ) {
          params.virtualInfo = [{ card_no: '', card_pwd: '' }];
          void formApi.setFieldValue('virtualInfo', params.virtualInfo);
        }
      },
    },
    fieldName: 'card_type',
    label: '卡密类型',
  },
  {
    component: 'Textarea',
    componentProps: { placeholder: '填写卡密信息' },
    dependencies: {
      show: (values) => Number(values.card_type) === 10,
      triggerFields: ['card_type'],
    },
    fieldName: 'card_info',
    label: '卡密信息',
  },
  {
    component: 'Input',
    componentProps: { placeholder: '填写卡密库存' },
    dependencies: {
      show: (values) => Number(values.card_type) === 10,
      triggerFields: ['card_type'],
    },
    fieldName: 'stock_num',
    label: '库存',
    suffix: '件',
  },
  {
    component: markRaw(VirtualCardListField),
    dependencies: {
      show: (values) => Number(values.card_type) === 20,
      triggerFields: ['card_type'],
    },
    fieldName: 'virtualInfo',
    formItemClass: 'col-span-full max-w-none',
    hideLabel: true,
    label: '',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'max-w460' },
      labelWidth: 100,
    },
    handleValuesChange(values, fieldsChanged) {
      for (const key of fieldsChanged) {
        if (key in params) {
          (params as Record<string, unknown>)[key] = values[key];
        }
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
  }),
);

watch(
  () => props.formData,
  (value) => {
    Object.assign(params, value);
    if (!params.card_type) {
      params.card_type = 10;
    }
    if (
      params.card_type === 20 &&
      (!params.virtualInfo || params.virtualInfo === '')
    ) {
      params.virtualInfo = [{ card_no: '', card_pwd: '' }];
    }
    void formApi.setValues({ ...params });
  },
  { deep: true, immediate: true },
);

function confirm() {
  if (params.card_type === 10) {
    if (!params.card_info || !params.stock_num) {
      ElMessage.error('请输入卡密内容和件数');
      return;
    }
  }
  if (params.card_type === 20 && Array.isArray(params.virtualInfo)) {
    if (params.virtualInfo.some((item) => !item.card_no || !item.card_pwd)) {
      ElMessage.error('请输入卡号和密码');
      return;
    }
    params.stock_num = params.virtualInfo.length;
  }
  emit('close', { type: 'success', params: { ...params } });
  open.value = false;
}

function cancel() {
  emit('close', { type: 'error', params: { ...params } });
  open.value = false;
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) modalApi.open();
  else modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[600px]"
    title="卡密管理"
  >
    <div class="product-carmine-dialog">
      <Form />
    </div>

    <template #footer>
      <ElButton size="small" @click="cancel">取 消</ElButton>
      <ElButton size="small" type="primary" @click="confirm">确 定</ElButton>
    </template>
  </Modal>
</template>

<style scoped lang="scss">
.max-w460 {
  width: 100%;
  max-width: 460px;
}

.virtual-info {
  max-height: 40vh;
  margin-bottom: 16px;
  overflow-y: auto;
}

.virtual-info__row {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 16px;
}

.virtual-info__field {
  display: flex;
  gap: 8px;
  align-items: center;
}

.virtual-info__label {
  white-space: nowrap;
}

.virtual-info__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}
</style>
