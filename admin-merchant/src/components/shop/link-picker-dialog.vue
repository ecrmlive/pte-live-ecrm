<script setup lang="ts">
import type { ShopLinkPickerItem } from '#/api/core/shop-link';
import type { ShopLinkValue } from '#/types/shop-link';

import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { ref, watch } from 'vue';

import LinkArticlePanel from './link-picker/panels/link-article-panel.vue';
import LinkDiyPagePanel from './link-picker/panels/link-diy-page-panel.vue';
import LinkH5Panel from './link-picker/panels/link-h5-panel.vue';
import LinkMarketingPanel from './link-picker/panels/link-marketing-panel.vue';
import LinkMenuPanel from './link-picker/panels/link-menu-panel.vue';
import LinkMpPanel from './link-picker/panels/link-mp-panel.vue';
import LinkPagesPanel from './link-picker/panels/link-pages-panel.vue';
import LinkProductPanel from './link-picker/panels/link-product-panel.vue';
import { finalizeShopLinkValue } from './link-picker/finalize-link';

defineOptions({ name: 'ShopLinkPickerDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  linkData?: ShopLinkValue | null;
}>();

const emit = defineEmits<{
  confirm: [ShopLinkValue | undefined];
}>();

const activeTab = ref('pages');
const activeData = ref<ShopLinkValue | null>(null);

function onPanelChange(item: ShopLinkPickerItem) {
  activeData.value = item;
}

function resetting() {
  activeData.value = { name: '请选择', type: '', url: '' };
  confirm(true);
}

function confirm(ok: boolean) {
  if (!ok) {
    open.value = false;
    emit('confirm', undefined);
    return;
  }
  const finalized = finalizeShopLinkValue(activeData.value);
  open.value = false;
  emit('confirm', finalized ?? undefined);
}

const [Modal, modalApi] = useVbenModal({
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
    class="shop-link-picker-dialog w-[800px]"
    title="超链接设置"
  >
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="页面" name="pages">
        <LinkPagesPanel v-if="activeTab === 'pages'" @change="onPanelChange" />
      </el-tab-pane>
      <el-tab-pane label="营销" name="market">
        <LinkMarketingPanel v-if="activeTab === 'market'" @change="onPanelChange" />
      </el-tab-pane>
      <el-tab-pane label="商品" name="product">
        <LinkProductPanel v-if="activeTab === 'product'" @change="onPanelChange" />
      </el-tab-pane>
      <el-tab-pane label="文章" name="Article">
        <LinkArticlePanel v-if="activeTab === 'Article'" @change="onPanelChange" />
      </el-tab-pane>
      <el-tab-pane label="小程序" name="SmallProgram">
        <LinkMpPanel
          v-if="activeTab === 'SmallProgram'"
          :link-data="linkData"
          @change="onPanelChange"
        />
      </el-tab-pane>
      <el-tab-pane label="H5" name="H5">
        <LinkH5Panel v-if="activeTab === 'H5'" :link-data="linkData" @change="onPanelChange" />
      </el-tab-pane>
      <el-tab-pane label="自定义" name="diypage">
        <LinkDiyPagePanel v-if="activeTab === 'diypage'" @change="onPanelChange" />
      </el-tab-pane>
      <el-tab-pane label="个人中心菜单" name="menu">
        <LinkMenuPanel v-if="activeTab === 'menu'" @change="onPanelChange" />
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <div class="flex items-start justify-between gap-4">
        <div class="min-w-0 flex-1 text-left text-sm">
          <template v-if="activeData?.url">
            <p class="truncate">
              <span>当前链接：</span>
              <span class="text-gray-400">{{ activeData.type }}</span>
              <span class="px-2 text-gray-300">/</span>
              <span class="text-primary">{{ activeData.name }}</span>
            </p>
            <p class="truncate text-[10px] text-gray-400">{{ activeData.url }}</p>
          </template>
          <span v-else class="text-gray-400">暂无</span>
        </div>
        <div class="flex shrink-0 gap-2">
          <ElButton size="small" @click="confirm(false)">取 消</ElButton>
          <ElButton size="small" type="warning" @click="resetting">清空</ElButton>
          <ElButton size="small" type="primary" @click="confirm(true)">确 定</ElButton>
        </div>
      </div>
    </template>
  </Modal>
</template>

<style scoped>
.shop-link-picker-dialog :deep(.marketing-box .el-tabs__item) {
  font-size: 12px;
}
</style>
