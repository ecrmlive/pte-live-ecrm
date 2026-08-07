<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import { computed, onMounted, ref, watch } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { VbenTiptap, VbenTiptapPreview } from '@vben/plugins/tiptap';
import { ElButton, ElMessage, ElTabPane, ElTabs } from 'element-plus';

import { uploadAttachmentApi } from '#/api/core/attachment';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformAgreementApi,
  savePlatformAgreementApi,
} from '#/api/core/platform-content';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

/** 对齐 CRMEB 店铺设置 / 说明提示：仅两个 Tab */
const TABS = [
  { key: 'sys_merchant_type', label: '店铺类型说明', name: 'type' },
  { key: 'sys_merchant_category', label: '店铺分类说明', name: 'category' },
] as const;

const activeName = ref<(typeof TABS)[number]['name']>('type');
const content = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const previewOpen = ref(false);

const currentTab = computed(
  () => TABS.find((item) => item.name === activeName.value) ?? TABS[0],
);

const imageUpload: ImageUploadOptions = {
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const row = await uploadAttachmentApi(file);
    return row.attachment_src;
  },
  onUploadError: () => {
    ElMessage.error('图片上传失败');
  },
};

const [PreviewDrawer, previewDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  footer: false,
  onOpenChange(isOpen) {
    previewOpen.value = isOpen;
  },
});

async function load() {
  loading.value = true;
  try {
    const row = await getPlatformAgreementApi(currentTab.value.key);
    content.value = row.content || '';
  } finally {
    loading.value = false;
  }
}

async function save() {
  const html = content.value.trim();
  if (!html || html === '<p></p>') {
    ElMessage.warning('说明正文不能为空');
    return;
  }
  saving.value = true;
  try {
    const saved = await savePlatformAgreementApi(currentTab.value.key, html);
    content.value = saved.content;
    ElMessage.success('已提交保存');
  } finally {
    saving.value = false;
  }
}

function openPreview() {
  previewDrawerApi.setState({ title: currentTab.value.label }).open();
}

watch(activeName, () => {
  void load();
});

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value =
    codes.includes('setting.agreement.manage') ||
    codes.includes('store.description');
  await load();
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <template #tabs>
        <ElTabs v-model="activeName">
          <ElTabPane
            v-for="tab in TABS"
            :key="tab.name"
            :label="tab.label"
            :name="tab.name"
          />
        </ElTabs>
      </template>

      <VbenTiptap
        v-model="content"
        :editable="canManage"
        :image-upload="imageUpload"
        :max-height="560"
        :min-height="420"
        :previewable="false"
        placeholder="请输入说明正文…"
      />

      <template #actions>
        <ElButton @click="openPreview">预览</ElButton>
        <ElButton
          v-if="canManage"
          :loading="saving"
          type="primary"
          @click="save"
        >
          提交
        </ElButton>
      </template>
    </SettingsTabLayout>

    <PreviewDrawer class="w-[720px]">
      <VbenTiptapPreview
        v-if="previewOpen"
        :content="content"
        :min-height="240"
      />
    </PreviewDrawer>
  </Page>
</template>
