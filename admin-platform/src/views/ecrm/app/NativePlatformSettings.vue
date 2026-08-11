<script setup lang="ts">
import type {
  CloudConfigGroup,
  NativePlatform,
} from '#/api/core/cloud-config';

import { computed, onMounted, reactive, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTabs,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getMobileAppConfigApi,
  getPushConfigApi,
  saveMobileAppConfigApi,
  savePushConfigApi,
} from '#/api/core/cloud-config';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

type ConfigKind = 'mobile' | 'push';

interface ConfigSection {
  canManage: boolean;
  group?: CloudConfigGroup;
  hint: string;
  kind: ConfigKind;
  title: string;
}

const TABS: Array<{ label: string; value: NativePlatform }> = [
  { label: 'iOS', value: 'ios' },
  { label: 'Android', value: 'android' },
  { label: 'HarmonyOS', value: 'harmony' },
];

const activePlatform = ref<NativePlatform>('ios');
const loading = ref(false);
const canManageMobile = ref(false);
const canManagePush = ref(false);
const groups = reactive<Partial<Record<ConfigKind, CloudConfigGroup>>>({});
const forms = reactive<Record<ConfigKind, Record<string, string>>>({
  mobile: {},
  push: {},
});
const secretConfigured = reactive<Record<ConfigKind, Record<string, boolean>>>({
  mobile: {},
  push: {},
});
const saving = reactive<Record<ConfigKind, boolean>>({
  mobile: false,
  push: false,
});

const sections = computed<ConfigSection[]>(() => [
  {
    kind: 'mobile',
    title: '基础与发布',
    hint: '应用标识用于各原生端打包、下载更新、跳转与服务端配置对齐。',
    canManage: canManageMobile.value,
    group: groups.mobile,
  },
  {
    kind: 'push',
    title: '推送配置',
    hint: '友盟及 APNs 密钥均加密保存，保存后不会回显明文；密钥字段留空即保留原值。',
    canManage: canManagePush.value,
    group: groups.push,
  },
]);

function isEnabled(value: string | undefined) {
  return value === '1' || value === 'true';
}

function setSwitch(kind: ConfigKind, key: string, enabled: boolean) {
  forms[kind][key] = enabled ? 'true' : 'false';
}

function applyGroup(kind: ConfigKind, data: CloudConfigGroup) {
  groups[kind] = data;
  const form = forms[kind];
  const configured = secretConfigured[kind];
  for (const key of Object.keys(form)) delete form[key];
  for (const key of Object.keys(configured)) delete configured[key];
  for (const field of data.fields) {
    const value = data.values[field.key] || '';
    form[field.key] = field.secret ? '' : value;
    configured[field.key] = field.secret && Boolean(value);
  }
}

async function load(kind: ConfigKind) {
  const data = kind === 'mobile'
    ? await getMobileAppConfigApi(activePlatform.value)
    : await getPushConfigApi(activePlatform.value);
  applyGroup(kind, data);
}

async function loadAll() {
  loading.value = true;
  try {
    const tasks: Promise<void>[] = [];
    if (canManageMobile.value) tasks.push(load('mobile'));
    else delete groups.mobile;
    if (canManagePush.value) tasks.push(load('push'));
    else delete groups.push;
    await Promise.all(tasks);
  } finally {
    loading.value = false;
  }
}

async function save(kind: ConfigKind) {
  const current = groups[kind];
  if (!current) return;
  const values: Record<string, string> = {};
  for (const field of current.fields) {
    const value = (forms[kind][field.key] || '').trim();
    if (
      field.required &&
      !value &&
      !(field.secret && secretConfigured[kind][field.key])
    ) {
      ElMessage.warning(`请填写${field.label}`);
      return;
    }
    values[field.key] = value;
  }
  saving[kind] = true;
  try {
    const data = kind === 'mobile'
      ? await saveMobileAppConfigApi(activePlatform.value, values)
      : await savePushConfigApi(activePlatform.value, values);
    applyGroup(kind, data);
    ElMessage.success(`${kind === 'mobile' ? '应用' : '推送'}配置已保存`);
  } finally {
    saving[kind] = false;
  }
}

watch(activePlatform, () => {
  void loadAll();
});

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManageMobile.value = codes.includes('app.mobile.manage');
  canManagePush.value = codes.includes('app.push.manage');
  await loadAll();
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <template #tabs>
        <ElTabs v-model="activePlatform">
          <ElTabPane
            v-for="tab in TABS"
            :key="tab.value"
            :label="tab.label"
            :name="tab.value"
          />
        </ElTabs>
      </template>

      <template v-for="(section, index) in sections" :key="section.kind">
        <ElDivider v-if="index > 0" />
        <section class="native-platform-settings__section">
          <div class="native-platform-settings__heading">
            <h2>{{ section.title }}</h2>
          </div>
          <ElAlert
            :title="section.hint"
            type="warning"
            :closable="false"
            class="mb-5"
          />

          <ElForm
            v-if="section.group"
            label-width="168px"
            class="max-w-3xl"
            :disabled="!section.canManage"
            @submit.prevent
          >
            <ElFormItem
              v-for="field in section.group.fields"
              :key="field.key"
              :label="field.label"
              :required="field.required"
            >
              <ElSwitch
                v-if="field.input_type === 'switch' || field.key === 'enabled'"
                :model-value="isEnabled(forms[section.kind][field.key])"
                @update:model-value="setSwitch(section.kind, field.key, $event)"
              />
              <ElSelect
                v-else-if="field.input_type === 'select'"
                v-model="forms[section.kind][field.key]"
                class="w-full"
                :placeholder="field.hint || `请选择${field.label}`"
              >
                <ElOption
                  v-for="option in field.options || []"
                  :key="option"
                  :label="option.toUpperCase()"
                  :value="option"
                />
              </ElSelect>
              <ElInput
                v-else
                v-model="forms[section.kind][field.key]"
                :type="
                  field.key === 'apns_p8_key' || field.input_type === 'textarea'
                    ? 'textarea'
                    : field.secret
                      ? 'password'
                      : field.input_type === 'number'
                        ? 'number'
                        : 'text'
                "
                :rows="field.key === 'apns_p8_key' || field.input_type === 'textarea' ? 5 : undefined"
                :show-password="field.secret && field.key !== 'apns_p8_key'"
                :placeholder="
                  field.secret && secretConfigured[section.kind][field.key]
                    ? `已配置；如需更换请输入新的${field.label}`
                    : field.hint || `请输入${field.label}`
                "
              />
              <p
                v-if="field.secret && secretConfigured[section.kind][field.key]"
                class="native-platform-settings__hint"
              >
                已配置，留空保存即可保留原值。
              </p>
            </ElFormItem>
          </ElForm>

          <div v-else class="native-platform-settings__forbidden">
            你没有维护此部分配置的按钮权限。
          </div>

          <div class="native-platform-settings__actions">
            <ElButton :disabled="!section.canManage" @click="load(section.kind)">
              重置
            </ElButton>
            <ElButton
              v-if="section.canManage"
              type="primary"
              :loading="saving[section.kind]"
              @click="save(section.kind)"
            >
              保存
            </ElButton>
          </div>
        </section>
      </template>
    </SettingsTabLayout>
  </Page>
</template>

<style scoped>
.native-platform-settings__section + .native-platform-settings__section {
  margin-top: 24px;
}

.native-platform-settings__heading {
  margin-bottom: 12px;
}

.native-platform-settings__heading h2 {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

.native-platform-settings__actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  max-width: 768px;
  margin-top: 8px;
}

.native-platform-settings__forbidden {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  line-height: 32px;
}

.native-platform-settings__hint {
  width: 100%;
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}
</style>
