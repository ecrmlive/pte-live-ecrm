<script setup lang="ts">
import { computed } from 'vue';

import { VbenTiptapPreview } from '@vben/plugins/tiptap';
import {
  ElEmpty,
  ElForm,
  ElFormItem,
  ElImage,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
} from 'element-plus';

import type { PlatformProductEditDetail, PlatformProductEditSKU } from '#/api/core/platform-catalog';
import type { PlatformPresell } from '#/api/core/platform-presell';
import ImageField from '#/components/shop/image-field.vue';
import ImagesField from '#/components/shop/images-field.vue';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const props = withDefaults(
  defineProps<{
    product?: PlatformProductEditDetail;
    productMissing?: boolean;
    productMissingTip?: string;
    presell?: PlatformPresell;
    showPresellTab?: boolean;
  }>(),
  {
    productMissing: false,
    productMissingTip: '关联商品编辑信息暂不可用，仍可查看/编辑预售活动信息。',
    showPresellTab: false,
  },
);

const activeTab = defineModel<string>({ default: 'basic' });

const coverImage = computed(
  () => props.product?.image || props.presell?.image || '',
);
const sliderImages = computed(() => props.product?.slider_image || []);

const skuSpecKeys = computed(() => {
  const keys = new Set<string>();
  for (const sku of props.product?.skus || []) {
    Object.keys(sku.spec || {}).forEach((k) => keys.add(k));
  }
  return [...keys];
});

const deliveryWayText = computed(() => {
  const ways = props.product?.delivery_way || [];
  if (!ways.length) return '—';
  return ways
    .map((w) =>
      Number(w) === 1 ? '到店自提' : Number(w) === 2 ? '快递配送' : String(w),
    )
    .join(' / ');
});

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '—';
  return String(v);
}

function skuSpecValue(sku: PlatformProductEditSKU, key: string) {
  return sku.spec?.[key] || '—';
}
</script>

<template>
  <div class="presell-product-tabs">
    <p v-if="productMissing" class="presell-product-tabs__tip">
      {{ productMissingTip }}
    </p>

    <ElTabs v-model="activeTab" class="presell-product-tabs__tabs">
      <ElTabPane label="商品信息" name="basic">
        <div class="presell-product-tabs__section-title">商品信息</div>
        <ElForm
          label-position="left"
          label-width="108px"
          class="presell-product-tabs__form"
        >
          <ElFormItem label="商品名称">
            <span>{{
              dash(product?.title || product?.store_name || presell?.store_name)
            }}</span>
          </ElFormItem>
          <div class="presell-product-tabs__field-grid">
            <ElFormItem label="平台分类">
              <span>{{ dash(product?.cate_name) }}</span>
            </ElFormItem>
            <ElFormItem label="品牌">
              <span>{{ dash(product?.brand_name) }}</span>
            </ElFormItem>
            <ElFormItem label="商品单位">
              <span>{{ dash(product?.unit_name) }}</span>
            </ElFormItem>
            <ElFormItem label="商品关键字">
              <span>{{ dash(product?.keyword) }}</span>
            </ElFormItem>
            <ElFormItem label="运费说明">
              <span>{{ deliveryWayText }}</span>
            </ElFormItem>
            <ElFormItem label="运费模板">
              <span>—</span>
            </ElFormItem>
            <ElFormItem label="商品分类">
              <span>{{ dash(product?.mer_cate_name) }}</span>
            </ElFormItem>
          </div>
          <ElFormItem label="商品简介">
            <span>{{ dash(product?.store_info || presell?.store_info) }}</span>
          </ElFormItem>
          <ElFormItem label="商品封面图">
            <ImageField
              v-if="coverImage"
              :model-value="coverImage"
              disabled
              :preview-size="72"
            />
            <span v-else>—</span>
          </ElFormItem>
          <ElFormItem label="商品轮播图">
            <ImagesField
              v-if="sliderImages.length"
              :model-value="sliderImages"
              disabled
              :preview-size="72"
              :limit="Math.max(sliderImages.length, 1)"
            />
            <div v-else-if="coverImage" class="presell-product-tabs__sliders">
              <ElImage
                class="presell-product-tabs__thumb"
                :src="resolveCosMediaUrl(coverImage)"
                fit="cover"
                :preview-src-list="[resolveCosMediaUrl(coverImage)]"
              />
            </div>
            <span v-else>—</span>
          </ElFormItem>
        </ElForm>
      </ElTabPane>

      <ElTabPane label="商品详情" name="content">
        <div class="presell-product-tabs__section-title">商品详情</div>
        <div class="presell-product-tabs__content">
          <VbenTiptapPreview
            v-if="product?.content"
            :content="product.content"
            :min-height="240"
          />
          <ElEmpty
            v-else
            :description="productMissing ? '商品详情不可用' : '暂无详情'"
            :image-size="72"
          />
        </div>
      </ElTabPane>

      <ElTabPane label="其他设置" name="other">
        <div class="presell-product-tabs__section-title">其他设置</div>
        <ElForm label-position="left" label-width="108px">
          <ElFormItem label="商品排序">
            <span>{{ dash(product?.rank) }}</span>
          </ElFormItem>
          <ElFormItem label="支持退款">
            <span>{{
              product
                ? Number(product.refund_switch) === 1
                  ? '是'
                  : '否'
                : '—'
            }}</span>
          </ElFormItem>
          <ElFormItem label="最少购买件数">
            <span>{{
              product
                ? Number(product.once_min_count || 0) > 0
                  ? String(product.once_min_count)
                  : '不限购'
                : '—'
            }}</span>
          </ElFormItem>
        </ElForm>
      </ElTabPane>

      <ElTabPane label="商品规格" name="sku">
        <div class="presell-product-tabs__section-title">商品规格</div>
        <ElForm label-position="left" label-width="108px">
          <ElFormItem label="多规格说明">
            <span>{{
              product
                ? Number(product.spec_type) === 1
                  ? '多规格'
                  : '单规格'
                : '—'
            }}</span>
          </ElFormItem>
          <ElFormItem label="佣金设置">
            <span>{{ dash(product?.commission_text) || '默认' }}</span>
          </ElFormItem>
        </ElForm>

        <template v-if="product?.skus?.length">
          <div class="presell-product-tabs__section-title">规格表</div>
          <ElTable :data="product.skus" border empty-text="暂无规格">
            <ElTableColumn
              v-for="key in skuSpecKeys"
              :key="key"
              :label="key"
              min-width="90"
            >
              <template #default="{ row }">
                {{ skuSpecValue(row, key) }}
              </template>
            </ElTableColumn>
            <ElTableColumn label="图片" width="72" align="center">
              <template #default="{ row }">
                <ElImage
                  v-if="row.image"
                  class="presell-product-tabs__sku-img"
                  :src="resolveCosMediaUrl(row.image)"
                  fit="cover"
                />
                <span v-else>—</span>
              </template>
            </ElTableColumn>
            <ElTableColumn label="售价" width="96" align="center">
              <template #default="{ row }">
                {{ Number(row.price || 0).toFixed(2) }}
              </template>
            </ElTableColumn>
            <ElTableColumn label="划线价" width="96" align="center">
              <template #default="{ row }">
                {{ Number(row.ot_price || 0).toFixed(2) }}
              </template>
            </ElTableColumn>
            <ElTableColumn prop="stock" label="库存" width="80" align="center" />
            <ElTableColumn label="一级返佣" width="96" align="center">
              <template #default="{ row }">
                {{ Number(row.extension_one || 0).toFixed(2) }}
              </template>
            </ElTableColumn>
          </ElTable>
        </template>
        <div v-else class="presell-product-tabs__sku-fallback">
          <ElEmpty
            :description="
              productMissing
                ? '关联商品不存在或不可读，规格已降级为预售快照（无 SKU 表）'
                : '暂无 SKU，已按单规格预售快照展示'
            "
            :image-size="72"
          />
          <div v-if="presell && !productMissing" class="presell-product-tabs__snapshot">
            <div>预售价：¥{{ Number(presell.price || 0).toFixed(2) }}</div>
            <div>
              库存 / 限量：{{ presell.stock ?? 0 }} /
              {{ presell.stock_count ?? 0 }}
            </div>
          </div>
        </div>
      </ElTabPane>

      <ElTabPane v-if="showPresellTab" label="预售信息" name="presell">
        <div class="presell-product-tabs__section-title">预售信息</div>
        <slot name="presell" />
      </ElTabPane>
    </ElTabs>
  </div>
</template>

<style scoped>
.presell-product-tabs__tip {
  margin: 0 0 12px;
  padding: 8px 12px;
  border-radius: 6px;
  background: hsl(var(--warning) / 0.12);
  color: hsl(var(--foreground));
  font-size: 13px;
  line-height: 20px;
}

.presell-product-tabs__tabs {
  flex: 1;
  min-height: 320px;
  padding-bottom: 12px;
}

.presell-product-tabs__section-title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  color: hsl(var(--foreground));
  font-size: 14px;
  font-weight: 600;
}

.presell-product-tabs__section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 2px;
  background: hsl(var(--primary));
}

.presell-product-tabs__field-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 12px;
}

.presell-product-tabs__field-grid :deep(.el-form-item) {
  margin-bottom: 12px;
}

.presell-product-tabs__sliders {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.presell-product-tabs__thumb,
.presell-product-tabs__sku-img {
  width: 64px;
  height: 64px;
  border-radius: 4px;
}

.presell-product-tabs__content {
  min-height: 240px;
  padding: 8px 0;
}

.presell-product-tabs__sku-fallback {
  margin-top: 8px;
}

.presell-product-tabs__snapshot {
  margin-top: -8px;
  color: hsl(var(--muted-foreground));
  font-size: 13px;
  line-height: 22px;
  text-align: center;
}

@media (max-width: 960px) {
  .presell-product-tabs__field-grid {
    grid-template-columns: 1fr;
  }
}
</style>
