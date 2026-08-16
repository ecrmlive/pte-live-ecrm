<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { ArrowRight, Picture } from '@element-plus/icons-vue';
import { ElIcon, ElSwitch } from 'element-plus';
import { computed, ref } from 'vue';

import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';
import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import {
  diyColor,
  diyRadioGroup,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyInputField from './shared/diy-input-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsHotspot' });

type DiyRecord = Record<string, any>;
type HotModuleKey = 'category' | 'groupbuy' | 'newcomer' | 'seckill';

interface HotProduct {
  image: string;
  name: string;
  price: string;
  productId?: number;
}

interface HotModule {
  enabled: boolean;
  image: string;
  linkUrl: string;
  linkeType?: string;
  products?: HotProduct[];
  subtitle: string;
  title: string;
}

type HotModules = Record<HotModuleKey, HotModule>;

const hotModuleKeys: HotModuleKey[] = ['seckill', 'groupbuy', 'category', 'newcomer'];

const hotModuleMeta: Record<HotModuleKey, { hint: string; label: string; type: 'icon' | 'products' }> = {
  category: { hint: '图标与跳转链接独立配置', label: '热门分类', type: 'icon' },
  groupbuy: { hint: '选择 2 个团购商品', label: '热销团购', type: 'products' },
  newcomer: { hint: '图标与跳转链接独立配置', label: '新人专享', type: 'icon' },
  seckill: { hint: '选择 2 个秒杀商品', label: '限时秒杀', type: 'products' },
};

const { PrimaryButton } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: DiyRecord & {
    params?: DiyRecord;
    style?: DiyRecord;
  };
}>();

const editor = useDiyEditor();
const styleType = ref<'content' | 'style'>('content');
const productPickerOpen = ref(false);
const productTarget = ref<{ index: number; key: HotModuleKey } | null>(null);
const isLinkset = ref(false);
const linkTarget = ref<HotModuleKey | null>(null);
const linkData = ref<DiyRecord | null>(null);

function defaultProduct(slot: number): HotProduct {
  return {
    image: '',
    name: `商品${slot}`,
    price: '',
  };
}

function defaultHotspotModules(): HotModules {
  return {
    category: {
      enabled: true,
      image: '',
      linkUrl: '',
      subtitle: '全屋｜厨卫｜保洁',
      title: '日常保洁',
    },
    groupbuy: {
      enabled: true,
      image: '',
      linkUrl: '',
      products: [defaultProduct(1), defaultProduct(2)],
      subtitle: '超低价 随时退',
      title: '热销团购',
    },
    newcomer: {
      enabled: true,
      image: '',
      linkUrl: '',
      subtitle: '下单咨询即领取45元券',
      title: '新人专享',
    },
    seckill: {
      enabled: true,
      image: '',
      linkUrl: '',
      products: [defaultProduct(1), defaultProduct(2)],
      subtitle: '4折秒杀',
      title: '限时秒杀',
    },
  };
}

function ensureHotspotDefaults(item: DiyRecord) {
  const params = (item.params ??= {}) as DiyRecord;
  const style = (item.style ??= {}) as DiyRecord;
  const defaults = defaultHotspotModules();
  const savedModules = (params.modules ?? {}) as Partial<HotModules>;
  const modules = {} as HotModules;

  hotModuleKeys.forEach((key) => {
    const fallback = defaults[key];
    const saved = (savedModules[key] ?? {}) as Partial<HotModule>;
    modules[key] = { ...fallback, ...saved };
    if (fallback.products) {
      const savedProducts = Array.isArray(saved.products) ? saved.products : [];
      modules[key].products = fallback.products.map((product, index) => ({
        ...product,
        ...(savedProducts[index] ?? {}),
      }));
    }
  });
  params.modules = modules;

  if (style.float === undefined) style.float = 18;
  if (!style.background) style.background = 'rgba(255, 255, 255, 0)';
  if (!style.bgcolor) style.bgcolor = '#ffffff';
  if (style.paddingTop === undefined) style.paddingTop = 9;
  if (style.paddingBottom === undefined) style.paddingBottom = 0;
  if (style.paddingLeft === undefined) style.paddingLeft = 10;
  if (style.paddingRight === undefined) style.paddingRight = style.paddingLeft;
  if (style.marginTop === undefined) style.marginTop = 0;
  if (!style.radiusMode) style.radiusMode = 'all';
  if (style.cardRadius === undefined) style.cardRadius = 0;
  if (style.topLeftRadio === undefined) style.topLeftRadio = 0;
  if (style.topRightRadio === undefined) style.topRightRadio = 0;
  if (style.bottomLeftRadio === undefined) style.bottomLeftRadio = 0;
  if (style.bottomRightRadio === undefined) style.bottomRightRadio = 0;
  if (!style.cardShadow) style.cardShadow = 'off';
}

const modules = computed(() => {
  ensureHotspotDefaults(props.curItem);
  return props.curItem.params!.modules as HotModules;
});

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('内容设置', '四个模块分别独立配置；商品、图标和跳转页面均通过选择器关联。'),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style ?? {};
  return [
    diySection('卡片样式'),
    diySlider('style.float', '组件上浮：', { max: 48, min: 0 }),
    diyColor('style.background', '底部背景：', 'rgba(255, 255, 255, 0)', '透明'),
    diyColor('style.bgcolor', '组件背景：', '#ffffff', '透明'),
    diySlider('style.paddingTop', '上边距：', { max: 48, min: 0 }),
    diySlider('style.paddingBottom', '下边距：', { max: 48, min: 0 }),
    diySlider('style.paddingLeft', '左右边距：', { max: 48, min: 0 }),
    diySlider('style.marginTop', '页面上间距：', { max: 96, min: 0 }),
    diyRadioGroup(
      'style.radiusMode',
      '背景圆角：',
      [
        { label: '全部', value: 'all' },
        { label: '分别设置', value: 'individual' },
      ],
      true,
    ),
    ...(style.radiusMode === 'individual'
      ? [
          diySlider('style.topLeftRadio', '左上圆角：', { max: 48, min: 0 }),
          diySlider('style.topRightRadio', '右上圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomLeftRadio', '左下圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomRightRadio', '右下圆角：', { max: 48, min: 0 }),
        ]
      : [diySlider('style.cardRadius', '圆角值：', { max: 48, min: 0 })]),
    diyRadioGroup(
      'style.cardShadow',
      '开启阴影：',
      [
        { label: '关闭', value: 'off' },
        { label: '开启', value: 'on' },
      ],
    ),
  ];
});

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, {
  onInit: ensureHotspotDefaults,
});
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit: ensureHotspotDefaults,
});

function openProductPicker(key: HotModuleKey, index: number) {
  productTarget.value = { index, key };
  productPickerOpen.value = true;
}

function onProductPicked(product: PlatformProduct) {
  const target = productTarget.value;
  if (!target) return;
  const productSlot = modules.value[target.key].products?.[target.index];
  if (!productSlot) return;
  Object.assign(productSlot, {
    image: String(product.image ?? ''),
    name: String(product.store_name ?? '未命名商品'),
    price: String(product.price ?? ''),
    productId: Number(product.product_id),
  });
  productTarget.value = null;
}

function clearProduct(key: HotModuleKey, index: number) {
  modules.value[key].products?.splice(index, 1, defaultProduct(index + 1));
}

function selectModuleImage(key: HotModuleKey) {
  editor.onEditorSelectImage(modules.value[key], 'image');
}

function changeLink(key: HotModuleKey) {
  linkTarget.value = key;
  linkData.value = modules.value[key];
  isLinkset.value = true;
}

function closeLinkset(e: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  const key = linkTarget.value;
  if (!e || !key) return;
  Object.assign(modules.value[key], {
    linkUrl: e.url ?? '',
    linkeType: e.type,
  });
}
</script>

<template>
  <div>
    <div class="common-form common-form-new">
      <span>{{ curItem.name }}</span>
      <div class="diy-changes">
        <div class="diy-change" :class="{ active: styleType === 'content' }" @click="styleType = 'content'">
          内容
        </div>
        <div class="diy-change" :class="{ active: styleType === 'style' }" @click="styleType = 'style'">
          样式
        </div>
      </div>
    </div>

    <div v-show="styleType === 'content'">
      <ContentForm />
      <div class="hotspot-module-list">
        <section v-for="key in hotModuleKeys" :key="key" class="hotspot-module-card">
          <header class="hotspot-module-header">
            <div>
              <h4>{{ hotModuleMeta[key].label }}</h4>
              <p>{{ hotModuleMeta[key].hint }}</p>
            </div>
            <ElSwitch v-model="modules[key].enabled" aria-label="启用模块" />
          </header>

          <div class="hotspot-field-grid">
            <label class="hotspot-input-field">
              <span>主标题</span>
              <DiyInputField v-model="modules[key].title" :maxlength="8" />
            </label>
            <label class="hotspot-input-field">
              <span>副标题</span>
              <DiyInputField v-model="modules[key].subtitle" :maxlength="20" />
            </label>
          </div>

          <template v-if="hotModuleMeta[key].type === 'products'">
            <div class="hotspot-product-list">
              <article
                v-for="(product, productIndex) in modules[key].products"
                :key="productIndex"
                class="hotspot-product-slot"
              >
                <div class="hotspot-product-cover" @click="openProductPicker(key, productIndex)">
                  <img v-if="product.image" v-img-url="product.image" :alt="product.name" />
                  <div v-else class="hotspot-product-empty">
                    <ElIcon><Picture /></ElIcon>
                    <span>选择商品</span>
                  </div>
                </div>
                <div class="hotspot-product-copy">
                  <span>{{ product.name || '未选择商品' }}</span>
                  <strong v-if="product.price">¥{{ product.price }}</strong>
                </div>
                <div class="hotspot-product-actions">
                  <component :is="PrimaryButton" link size="small" @click="openProductPicker(key, productIndex)">
                    {{ product.productId ? '更换' : '选择商品' }}
                  </component>
                  <component
                    v-if="product.productId"
                    :is="PrimaryButton"
                    link
                    size="small"
                    @click="clearProduct(key, productIndex)"
                  >
                    清除
                  </component>
                </div>
              </article>
            </div>
          </template>

          <template v-else>
            <div class="hotspot-icon-config">
              <button class="hotspot-image-selector" type="button" @click="selectModuleImage(key)">
                <img v-if="modules[key].image" v-img-url="modules[key].image" :alt="modules[key].title" />
                <template v-else>
                  <ElIcon><Picture /></ElIcon>
                  <span>选择图标</span>
                </template>
              </button>
              <div class="hotspot-link-field">
                <span>跳转链接</span>
                <DiyInputField v-model="modules[key].linkUrl" placeholder="选择页面链接" readonly>
                  <template #suffix>
                    <ElIcon color="#333" size="16px" @click="changeLink(key)">
                      <ArrowRight />
                    </ElIcon>
                  </template>
                </DiyInputField>
              </div>
            </div>
          </template>
        </section>
      </div>
    </div>

    <div v-show="styleType === 'style'">
      <StyleForm />
    </div>

    <ProductPickerDialog v-model:open="productPickerOpen" @select="onProductPicked" />
    <DiyLinkPickerDialog
      v-if="isLinkset"
      :is_linkset="isLinkset"
      :link-data="linkData"
      @close-dialog="closeLinkset"
    >
      选择链接
    </DiyLinkPickerDialog>
  </div>
</template>

<style lang="scss" scoped>
.hotspot-module-list {
  display: grid;
  gap: 12px;
  padding: 0 16px 16px;
}

.hotspot-module-card {
  border: 1px solid #e6eaf2;
  border-radius: 10px;
  background: #fff;
  padding: 14px;
}

.hotspot-module-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;

  h4 {
    margin: 0;
    color: #1f2937;
    font-size: 14px;
    font-weight: 600;
  }

  p {
    margin: 4px 0 0;
    color: #94a3b8;
    font-size: 12px;
  }
}

.hotspot-field-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.hotspot-input-field {
  display: grid;
  gap: 6px;
  min-width: 0;
  color: #667085;
  font-size: 12px;
}

.hotspot-product-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-top: 12px;
}

.hotspot-product-slot {
  min-width: 0;
  border: 1px solid #edf0f5;
  border-radius: 8px;
  padding: 8px;
}

.hotspot-product-cover {
  display: flex;
  height: 88px;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 6px;
  background: #f7f9fc;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.hotspot-product-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: #98a2b3;
  font-size: 12px;

  .el-icon {
    font-size: 20px;
  }
}

.hotspot-product-copy {
  display: flex;
  min-height: 38px;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  margin-top: 7px;
  color: #475467;
  font-size: 12px;

  span {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  strong {
    flex: none;
    color: #ef4444;
    font-weight: 600;
  }
}

.hotspot-product-actions {
  display: flex;
  gap: 8px;
}

.hotspot-icon-config {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
}

.hotspot-image-selector {
  display: flex;
  width: 76px;
  height: 76px;
  flex: none;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  overflow: hidden;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
  color: #98a2b3;
  background: #f8fafc;
  font-size: 12px;
  cursor: pointer;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .el-icon {
    font-size: 20px;
  }
}

.hotspot-link-field {
  min-width: 0;
  flex: 1;

  > span {
    display: block;
    margin-bottom: 6px;
    color: #667085;
    font-size: 12px;
  }
}
</style>
