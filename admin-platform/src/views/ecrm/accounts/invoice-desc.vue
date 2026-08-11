<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { VbenTiptap, VbenTiptapPreview } from '@vben/plugins/tiptap';
import { ElButton, ElMessage } from 'element-plus';

import { uploadAttachmentApi } from '#/api/core/attachment';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformAgreementApi,
  savePlatformAgreementApi,
} from '#/api/core/platform-content';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

/** 对齐 CRMEB CacheRepository::RECEIPT_AGREE */
const AGREE_KEY = 'sys_receipt_agree';
const PAGE_TITLE = '发票说明';

const content = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const previewOpen = ref(false);

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
    const row = await getPlatformAgreementApi(AGREE_KEY);
    content.value = row.content || '';
  } finally {
    loading.value = false;
  }
}

async function save() {
  const html = content.value.trim();
  if (!html || html === '<p></p>') {
    ElMessage.warning('发票说明正文不能为空');
    return;
  }
  saving.value = true;
  try {
    const saved = await savePlatformAgreementApi(AGREE_KEY, html);
    content.value = saved.content;
    ElMessage.success('已提交保存');
  } finally {
    saving.value = false;
  }
}

function openPreview() {
  previewDrawerApi.setState({ title: PAGE_TITLE }).open();
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value =
    codes.includes('accounts.invoice.desc.manage') ||
    codes.includes('setting.agreement.manage');
  await load();
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <div class="invoice-desc">
        <div class="invoice-desc__title">{{ PAGE_TITLE }}</div>
        <VbenTiptap
          v-model="content"
          :editable="canManage"
          :image-upload="imageUpload"
          :max-height="560"
          :min-height="420"
          :previewable="false"
          placeholder="请输入发票说明…"
        />
      </div>

      <template #actions>
        <ElButton type="primary" @click="openPreview">预览</ElButton>
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

<style scoped>
.invoice-desc__title {
  margin-bottom: 16px;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
  color: hsl(var(--foreground));
  text-align: center;
}
</style>
