<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElMessage,
  ElTabPane,
  ElTabs,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformMerchantApplyConfigApi,
  savePlatformMerchantApplyConfigApi,
  type MerchantApplyFormField,
  type PlatformMerchantApplyConfig,
} from '#/api/core/platform-mall-setting';
import ImageField from '#/components/shop/image-field.vue';

import FormDiyEditor from './apply-form-diy/FormDiyEditor.vue';
import type { ApplyFormField } from './apply-form-diy/types';

type TabName = 'page' | 'form';

const activeTab = ref<TabName>('page');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformMerchantApplyConfig>({
  background_image: '',
  form_fields: [],
});

function normalizeLoadedFields(
  list: MerchantApplyFormField[],
): ApplyFormField[] {
  return list.map((field) => ({
    ...field,
    content_type: field.content_type || 'text',
    max_upload:
      field.type === 'image' ? field.max_upload || 8 : field.max_upload,
    city_level:
      field.type === 'city'
        ? field.city_level || 'province_city_district'
        : field.city_level,
    default_visible: field.default_visible,
    default_mode: field.default_mode,
    specify_value: field.specify_value || '',
  }));
}

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformMerchantApplyConfigApi();
    form.value = {
      background_image: data.config.background_image || '',
      form_fields: normalizeLoadedFields(
        Array.isArray(data.config.form_fields) ? data.config.form_fields : [],
      ),
    };
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!canManage.value) return;
  const titles = new Set<string>();
  for (const field of form.value.form_fields) {
    const title = field.title.trim();
    if (!title) {
      ElMessage.warning('自定义字段标题不能为空');
      activeTab.value = 'form';
      return;
    }
    if (titles.has(title)) {
      ElMessage.error(`存在重复标题：${title}`);
      activeTab.value = 'form';
      return;
    }
    titles.add(title);
  }
  saving.value = true;
  try {
    form.value = await savePlatformMerchantApplyConfigApi({
      background_image: form.value.background_image.trim(),
      form_fields: form.value.form_fields.map((field) => ({
        ...field,
        title: field.title.trim(),
        placeholder: (field.placeholder || '').trim(),
        default_value: (field.default_value || '').trim(),
        content_type: field.content_type || 'text',
        options:
          field.type === 'radio' ||
          field.type === 'checkbox' ||
          field.type === 'select'
            ? field.options?.filter(Boolean)
            : undefined,
        max_upload: field.type === 'image' ? field.max_upload || 8 : undefined,
        city_level: field.type === 'city' ? field.city_level : undefined,
        default_visible:
          field.type === 'date' ||
          field.type === 'daterange' ||
          field.type === 'time' ||
          field.type === 'timerange'
            ? field.default_visible || 'show'
            : undefined,
        default_mode:
          field.type === 'date' ||
          field.type === 'daterange' ||
          field.type === 'time' ||
          field.type === 'timerange'
            ? field.default_mode || 'current'
            : undefined,
        specify_value:
          field.type === 'date' ||
          field.type === 'daterange' ||
          field.type === 'time' ||
          field.type === 'timerange'
            ? (field.specify_value || '').trim()
            : undefined,
      })),
    });
    ElMessage.success('商户设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [codes] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = codes.includes('merchant.mgmt.settings.manage');
});
</script>

<template>
  <Page auto-content-height content-class="!bg-card !p-0">
    <div v-loading="loading" class="apply-setting">
      <div class="apply-setting__tabs">
        <ElTabs v-model="activeTab">
          <ElTabPane label="入驻页面配置" name="page" />
          <ElTabPane label="入驻表单设置" name="form" />
        </ElTabs>
      </div>

      <div v-show="activeTab === 'page'" class="apply-setting__page">
        <ElForm
          label-width="100px"
          :disabled="!canManage"
          class="apply-setting__page-form"
        >
          <ElFormItem label="背景图">
            <div>
              <ImageField
                v-model="form.background_image"
                :disabled="!canManage"
                default-library="system"
                :preview-size="52"
                :show-button="false"
              />
              <div class="tips">点击缩略图从素材库选择入驻页背景</div>
            </div>
          </ElFormItem>
          <ElFormItem v-if="canManage">
            <ElButton :loading="saving" type="primary" @click="save">
              保存
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>

      <div v-show="activeTab === 'form'" class="apply-setting__form">
        <FormDiyEditor
          v-model:fields="form.form_fields"
          :disabled="!canManage"
        />
        <div class="apply-setting__footer">
          <ElButton
            v-if="canManage"
            :loading="saving"
            type="primary"
            @click="save"
          >
            保存
          </ElButton>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped lang="scss">
.apply-setting {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: hsl(var(--card));
}

.apply-setting__tabs {
  flex-shrink: 0;
  padding: 0 16px;
  border-bottom: 1px solid hsl(var(--border));
}

.apply-setting__tabs :deep(.el-tabs__header) {
  margin: 0;
}

.apply-setting__tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.apply-setting__tabs :deep(.el-tabs__item) {
  height: 44px;
  line-height: 44px;
}

.apply-setting__page {
  flex: 1;
  min-height: 0;
  padding: 24px 32px 40px;
  overflow: auto;
}

.apply-setting__page-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.apply-setting__page-form :deep(.el-form-item__label) {
  align-items: center;
  font-weight: 400;
}

.tips {
  margin-top: 6px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
}

.apply-setting__form {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
}

.apply-setting__footer {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  height: 60px;
  background: #fff;
  box-shadow: 0 -1px 4px rgb(0 0 0 / 10%);
}
</style>
