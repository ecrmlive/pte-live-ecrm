<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { DevicePreviewFrame } from '@pte-live/diy';
import { ElSkeleton } from 'element-plus';

import {
  getPlatformProductEditApi,
  type PlatformProductEditDetail,
} from '#/api/core/platform-catalog';

import ProductDetailDiyPreview, {
  type ProductDetailDiyPreviewData,
} from './ProductDetailDiyPreview.vue';

const props = withDefaults(
  defineProps<{
    /** 弹窗标题，默认「商品预览」 */
    modalTitle?: string;
    productId?: number;
    productTitle?: string;
    /** 覆盖机框内展示售价（如秒杀价） */
    displayPrice?: number;
    /** 覆盖划线价 */
    displayOtPrice?: number;
  }>(),
  {
    modalTitle: '商品预览',
  },
);

const emit = defineEmits<{
  closed: [];
}>();

const loadingProduct = ref(false);
const detail = ref<PlatformProductEditDetail>();

const previewProduct = computed<ProductDetailDiyPreviewData | null>(() => {
  const base = detail.value
    ? { ...detail.value }
    : props.productTitle
      ? ({ title: props.productTitle } as ProductDetailDiyPreviewData)
      : null;
  if (!base) return null;
  if (props.displayPrice !== undefined && props.displayPrice !== null) {
    base.price = Number(props.displayPrice);
  }
  if (props.displayOtPrice !== undefined && props.displayOtPrice !== null) {
    base.ot_price = Number(props.displayOtPrice);
  }
  return base;
});

const [Modal, modalApi] = useVbenModal({
  title: props.modalTitle,
  class: 'w-[520px] max-w-[96vw]',
  footer: false,
  destroyOnClose: true,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      detail.value = undefined;
      emit('closed');
      return;
    }
    modalApi.setState({ title: props.modalTitle || '商品预览' });
    void loadProduct();
  },
});

async function loadProduct() {
  const id = Number(props.productId || 0);
  if (!id) {
    detail.value = undefined;
    return;
  }
  loadingProduct.value = true;
  try {
    detail.value = await getPlatformProductEditApi(id);
  } catch {
    detail.value = undefined;
  } finally {
    loadingProduct.value = false;
  }
}

function open() {
  modalApi.setState({ title: props.modalTitle || '商品预览' }).open();
}

function close() {
  modalApi.close();
}

defineExpose({ open, close });
</script>

<template>
  <Modal>
    <div class="product-preview-modal">
      <div class="product-preview-modal__frame">
        <DevicePreviewFrame
          :show-device-switcher="true"
          :show-back="true"
          :show-expand="true"
          :hide-nav="false"
          :side-gutter="0"
          content-bg="#f5f5f5"
          title="商品详情"
        >
          <ElSkeleton
            :loading="loadingProduct && !previewProduct"
            animated
            :rows="8"
            class="product-preview-modal__skeleton"
          >
            <template #default>
              <ProductDetailDiyPreview :product="previewProduct" />
            </template>
          </ElSkeleton>
        </DevicePreviewFrame>
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.product-preview-modal {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 420px;
}

.product-preview-modal__frame {
  display: flex;
  justify-content: center;
  width: 100%;
  max-height: min(78vh, 780px);
  overflow: auto;
}

.product-preview-modal__skeleton {
  min-height: 480px;
  padding: 12px;
}
</style>
