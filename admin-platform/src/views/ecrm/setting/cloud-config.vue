<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Page } from '@vben/common-ui';

import {
  fetchCloudConfigs,
  saveCloudConfig,
  type CloudConfigGroup,
} from '#/api/core/cloud-config';

const groups = ref<CloudConfigGroup[]>([]);
const active = ref('payment');
const loading = ref(false);
const saving = ref(false);

async function load() {
  loading.value = true;
  try {
    groups.value = await fetchCloudConfigs();
    if (!groups.value.some((item) => item.group_key === active.value)) {
      active.value = groups.value[0]?.group_key || '';
    }
  } catch {
    ElMessage.error('加载云服务配置失败');
  } finally {
    loading.value = false;
  }
}

async function save(group: CloudConfigGroup) {
  saving.value = true;
  try {
    const updated = await saveCloudConfig(group.group_key, group.values);
    const index = groups.value.findIndex((item) => item.group_key === updated.group_key);
    if (index >= 0) groups.value[index] = updated;
    ElMessage.success('已加密保存');
  } catch {
    ElMessage.error('保存失败，请检查字段与权限');
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page title="云服务配置" description="支付、腾讯云、直播与回调配置均加密存库；密钥保存后仅显示掩码。">
    <ElAlert
      class="mb-4"
      title="密钥字段不会从接口返回明文。保留掩码并保存时不会覆盖已有密钥。"
      type="warning"
      :closable="false"
    />
    <ElSkeleton v-if="loading" :rows="8" animated />
    <ElTabs v-else v-model="active" tab-position="left" class="min-h-120">
      <ElTabPane v-for="group in groups" :key="group.group_key" :name="group.group_key" :label="group.label">
        <ElForm label-position="top" class="max-w-4xl">
          <ElRow :gutter="16">
            <ElCol v-for="field in group.fields" :key="field.key" :span="12">
              <ElFormItem :label="field.label" :required="field.required">
                <ElInput
                  v-model="group.values[field.key]"
                  :placeholder="field.hint || (field.secret ? '输入新值以替换已有密钥' : `请输入${field.label}`)"
                  :show-password="field.secret"
                  :type="field.secret ? 'password' : 'text'"
                />
              </ElFormItem>
            </ElCol>
          </ElRow>
          <ElButton type="primary" :loading="saving" @click="save(group)">保存配置</ElButton>
        </ElForm>
      </ElTabPane>
    </ElTabs>
  </Page>
</template>
