<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { useVbenModal } from '@vben/common-ui';
import { DevicePreviewFrame } from '@pte-live/diy';
import { ElButton, ElLink, ElSkeleton } from 'element-plus';

import {
  getPlatformProductEditApi,
  type PlatformProductEditDetail,
} from '#/api/core/platform-catalog';
import { getPlatformShopConfigApi } from '#/api/core/platform-mall-setting';

import ProductDetailDiyPreview, {
  type ProductDetailDiyPreviewData,
} from './ProductDetailDiyPreview.vue';

const props = defineProps<{
  productId?: number;
  productTitle?: string;
}>();

const emit = defineEmits<{
  closed: [];
}>();

const router = useRouter();
const siteUrl = ref('');
const loadingConfig = ref(false);
const loadingProduct = ref(false);
const detail = ref<PlatformProductEditDetail>();

/** 对齐 app-uni H5 hash 路由：pages/goods/detail?id= */
function buildGoodsDetailURL(base: string, productId: number) {
  const trimmed = base.replace(/\/+$/, '');
  return `${trimmed}/#/pages/goods/detail?id=${productId}`;
}

const previewURL = computed(() => {
  const id = Number(props.productId || 0);
  const base = siteUrl.value.trim();
  if (!id || !base) return '';
  return buildGoodsDetailURL(base, id);
});

const hasSiteURL = computed(() => Boolean(siteUrl.value.trim()));

const displayTitle = computed(
  () => detail.value?.title || props.productTitle || '',
);

const previewProduct = computed<ProductDetailDiyPreviewData | null>(() => {
  if (detail.value) return detail.value;
  if (!props.productTitle) return null;
  return { title: props.productTitle };
});

const [Modal, modalApi] = useVbenModal({
  title: '商品预览',
  class: 'w-[520px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  destroyOnClose: true,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      detail.value = undefined;
      emit('closed');
      return;
    }
    void Promise.all([loadSiteURL(), loadProduct()]);
  },
});

async function loadSiteURL() {
  loadingConfig.value = true;
  try {
    const data = await getPlatformShopConfigApi();
    siteUrl.value = String(data.config?.site_url || '').trim();
  } catch {
    siteUrl.value = '';
  } finally {
    loadingConfig.value = false;
  }
}

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
  modalApi.open();
}

function close() {
  modalApi.close();
}

function goShopSetting() {
  close();
  void router.push('/setting/shop');
}

function openInNewTab() {
  if (!previewURL.value) return;
  window.open(previewURL.value, '_blank', 'noopener,noreferrer');
}

defineExpose({ open, close });
</script>

<template>
  <Modal>
    <div class="product-preview-modal">
      <p v-if="displayTitle" class="product-preview-modal__name">
        {{ displayTitle }}
      </p>

      <div class="product-preview-modal__frame">
        <DevicePreviewFrame
          :show-device-switcher="true"
          :hide-nav="true"
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

      <div class="product-preview-modal__tips">
        <p class="product-preview-modal__warn">
          若页面未加载出，请前往系统配置页面填写网站域名
        </p>
        <ElButton type="primary" link @click="goShopSetting">点击前往</ElButton>
      </div>

      <div
        v-if="!loadingConfig && hasSiteURL"
        class="product-preview-modal__footer"
      >
        <ElLink type="primary" :underline="false" @click="openInNewTab">
          新窗口打开 H5
        </ElLink>
        <span class="product-preview-modal__url" :title="previewURL">
          {{ previewURL }}
        </span>
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.product-preview-modal {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: center;
  min-height: 420px;
}

.product-preview-modal__name {
  align-self: stretch;
  margin: 0;
  color: hsl(var(--foreground));
  font-size: 14px;
  font-weight: 600;
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

.product-preview-modal__tips {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
  width: 100%;
  padding: 4px 8px 0;
  text-align: center;
}

.product-preview-modal__warn {
  margin: 0;
  color: #f56c6c;
  font-size: 13px;
  line-height: 20px;
}

.product-preview-modal__footer {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
  width: 100%;
}

.product-preview-modal__url {
  max-width: 100%;
  overflow: hidden;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
