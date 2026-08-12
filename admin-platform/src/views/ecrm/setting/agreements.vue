<script setup lang="ts">
import type { ImageUploadOptions } from "@vben/plugins/tiptap";

import { computed, onMounted, ref } from "vue";

import { Page } from "@vben/common-ui";
import { DevicePreviewFrame } from "@pte-live/diy";
import { VbenTiptap, VbenTiptapPreview } from "@vben/plugins/tiptap";
import { ElButton, ElMessage } from "element-plus";

import { uploadAttachmentApi } from "#/api/core/attachment";
import { getAccessCodesApi } from "#/api/core/auth";
import {
  listPlatformAgreementsApi,
  savePlatformAgreementApi,
  type PlatformAgreement,
} from "#/api/core/platform-content";

const loading = ref(false);
const saving = ref(false);
const rows = ref<PlatformAgreement[]>([]);
const activeKey = ref("");
const content = ref("");
const savedContent = ref("");
const canManage = ref(false);

const current = computed(
  () => rows.value.find((item) => item.key === activeKey.value) ?? rows.value[0],
);
const isDirty = computed(() => content.value.trim() !== savedContent.value.trim());

const imageUpload: ImageUploadOptions = {
  accept: "image/jpeg,image/png,image/gif,image/webp",
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const row = await uploadAttachmentApi(file);
    return row.attachment_src;
  },
  onUploadError: () => {
    ElMessage.error("图片上传失败");
  },
};

function select(row: PlatformAgreement) {
  activeKey.value = row.key;
  content.value = row.content || "";
  savedContent.value = row.content || "";
}

function reset() {
  content.value = savedContent.value;
}

async function load() {
  loading.value = true;
  try {
    const data = await listPlatformAgreementsApi();
    rows.value = data.list || [];
    if (rows.value.length > 0) {
      select(rows.value.find((item) => item.key === activeKey.value) ?? rows.value[0]);
    }
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!current.value) return;
  const html = content.value.trim();
  if (!html || html === "<p></p>") {
    ElMessage.warning("协议正文不能为空");
    return;
  }

  saving.value = true;
  try {
    const saved = await savePlatformAgreementApi(current.value.key, html);
    const index = rows.value.findIndex((item) => item.key === saved.key);
    if (index >= 0) rows.value[index] = saved;
    content.value = saved.content;
    savedContent.value = saved.content;
    ElMessage.success("保存成功");
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes("setting.agreement.manage");
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <div v-loading="loading" class="agreement-settings">
      <aside class="agreement-settings__nav" aria-label="协议类型">
        <button
          v-for="row in rows"
          :key="row.key"
          :class="[
            'agreement-settings__nav-item',
            { 'agreement-settings__nav-item--active': current?.key === row.key },
          ]"
          type="button"
          @click="select(row)"
        >
          {{ row.label }}
        </button>
      </aside>

      <section class="agreement-settings__preview" aria-label="移动端效果预览">
        <DevicePreviewFrame
          device="ios"
          :show-back="true"
          :show-device-switcher="false"
          :show-expand="false"
          :side-gutter="0"
          content-bg="#ffffff"
          :title="current?.label || '协议内容'"
        >
          <article class="agreement-mobile-content">
            <h1>{{ current?.label || "协议内容" }}</h1>
            <VbenTiptapPreview
              :content="content"
              :min-height="0"
              class="agreement-mobile-content__body"
            />
          </article>
        </DevicePreviewFrame>
      </section>

      <section class="agreement-settings__editor">
        <header class="agreement-settings__editor-heading">
          <span>{{ current?.label || "协议设置" }}</span>
        </header>

        <div v-if="current" class="agreement-settings__editor-body">
          <VbenTiptap
            v-model="content"
            :editable="canManage"
            :image-upload="imageUpload"
            :max-height="680"
            :min-height="460"
            :previewable="false"
            placeholder="请输入协议正文…"
          />
        </div>
        <div v-else class="agreement-settings__empty">暂无可维护的协议</div>

        <footer class="agreement-settings__actions">
          <ElButton :disabled="!canManage || !isDirty" @click="reset"> 重置 </ElButton>
          <ElButton
            :disabled="!canManage || !isDirty"
            :loading="saving"
            type="primary"
            @click="save"
          >
            保存
          </ElButton>
        </footer>
      </section>
    </div>
  </Page>
</template>

<style scoped>
.agreement-settings {
  display: grid;
  grid-template-columns: 232px minmax(460px, 520px) minmax(0, 1fr);
  min-height: max(680px, calc(100vh - 168px));
  overflow: hidden;
  background: var(--vben-bg-color-overlay, #fff);
}

.agreement-settings__nav {
  overflow-y: auto;
  border-right: 1px solid var(--vben-border-color, #e5e7eb);
  padding: 24px 10px;
}

.agreement-settings__nav-item {
  display: flex;
  width: 100%;
  min-height: 52px;
  align-items: center;
  justify-content: center;
  border: 0;
  border-right: 3px solid transparent;
  background: transparent;
  color: var(--vben-text-color, #4b5563);
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
  text-align: center;
}

.agreement-settings__nav-item:hover {
  background: var(--vben-bg-color-hover, #f5f7fa);
}

.agreement-settings__nav-item--active {
  border-right-color: #1677ff;
  background: #edf4ff;
  color: #1677ff;
}

.agreement-settings__preview {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: center;
  overflow: auto;
  border-right: 1px solid var(--vben-border-color, #e5e7eb);
  background: #f4f6f9;
  padding: 28px 20px;
}

.agreement-mobile-content {
  box-sizing: border-box;
  min-height: 100%;
  padding: 24px 20px 40px;
  color: #1f2937;
  background: #fff;
}

.agreement-mobile-content h1 {
  margin: 0 0 22px;
  color: #111827;
  font-size: 22px;
  font-weight: 700;
  line-height: 1.35;
  text-align: center;
}

.agreement-mobile-content__body {
  color: #374151;
  font-size: 14px;
  line-height: 1.85;
}

.agreement-mobile-content__body :deep(*) {
  max-width: 100%;
}

.agreement-mobile-content__body :deep(h1),
.agreement-mobile-content__body :deep(h2),
.agreement-mobile-content__body :deep(h3) {
  margin: 20px 0 10px;
  color: #111827;
  font-size: 17px;
  font-weight: 700;
  line-height: 1.5;
}

.agreement-mobile-content__body :deep(p) {
  margin: 0 0 12px;
}

.agreement-settings__editor {
  display: flex;
  min-width: 0;
  flex-direction: column;
  overflow: hidden;
}

.agreement-settings__editor-heading {
  display: flex;
  min-height: 84px;
  align-items: center;
  border-bottom: 1px solid var(--vben-border-color, #e5e7eb);
  padding: 0 32px;
}

.agreement-settings__editor-heading span {
  border-left: 4px solid #1677ff;
  padding-left: 14px;
  color: #303133;
  font-size: 25px;
  font-weight: 700;
}

.agreement-settings__editor-body {
  min-height: 0;
  flex: 1;
  overflow: auto;
  padding: 24px 32px;
}

.agreement-settings__editor-body :deep(.vben-tiptap) {
  min-width: 620px;
}

.agreement-settings__empty {
  display: grid;
  min-height: 260px;
  flex: 1;
  place-items: center;
  color: #909399;
}

.agreement-settings__actions {
  display: flex;
  min-height: 76px;
  align-items: center;
  justify-content: center;
  gap: 12px;
  border-top: 1px solid var(--vben-border-color, #e5e7eb);
  background: #fff;
}

@media (max-width: 1200px) {
  .agreement-settings {
    grid-template-columns: 200px minmax(420px, 450px) minmax(0, 1fr);
  }

  .agreement-settings__preview {
    padding: 24px 14px;
  }

  .agreement-settings__editor-heading,
  .agreement-settings__editor-body {
    padding-left: 22px;
    padding-right: 22px;
  }
}

@media (max-width: 960px) {
  .agreement-settings {
    grid-template-columns: 190px minmax(0, 1fr);
  }

  .agreement-settings__preview {
    display: none;
  }
}
</style>
