<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';

import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSwitch,
  ElTabPane,
  ElTabs,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getNotificationConfigApi,
  listNotificationConfigsApi,
  type NotificationAudience,
  type NotificationConfig,
  type NotificationConfigSaveInput,
  saveNotificationConfigApi,
  syncNotificationTemplatesApi,
} from '#/api/core/platform-notification';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';

type NotificationForm = NotificationConfigSaveInput & {
  notice_type: string;
  scene: string;
};

const audience = ref<NotificationAudience>('member');
const editing = ref<NotificationConfig>();
const form = reactive<NotificationForm>({
  mini_program_enabled: 0,
  mini_program_text: '',
  notice_type: '',
  scene: '',
  sms_enabled: 0,
  sms_text: '',
  wechat_enabled: 0,
  wechat_text: '',
});

const gridOptions: VxeGridProps<NotificationConfig> = {
  columns: [
    { field: 'notification_id', title: 'ID', width: 88 },
    { field: 'notice_type', minWidth: 210, title: '通知类型' },
    { field: 'scene', minWidth: 260, title: '通知场景说明' },
    {
      align: 'center',
      field: 'wechat_enabled',
      slots: { default: 'wechat' },
      title: '公众号模板',
      width: 150,
    },
    {
      align: 'center',
      field: 'mini_program_enabled',
      slots: { default: 'miniProgram' },
      title: '小程序订阅',
      width: 150,
    },
    {
      align: 'center',
      field: 'sms_enabled',
      slots: { default: 'sms' },
      title: '发送短信',
      width: 140,
    },
    platformListActionColumn({ width: 116 }),
  ],
  emptyText: '暂无通知配置',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const result = await listNotificationConfigsApi(audience.value, {
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'notification_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });
const [ConfigDrawer, configDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[min(96vw,1000px)]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: save,
});

function setForm(row: NotificationConfig) {
  Object.assign(form, {
    mini_program_enabled: row.mini_program_enabled,
    mini_program_text: row.mini_program_text || '',
    notice_type: row.notice_type,
    scene: row.scene,
    sms_enabled: row.sms_enabled,
    sms_text: row.sms_text || '',
    wechat_enabled: row.wechat_enabled,
    wechat_text: row.wechat_text || '',
  });
}

async function openConfig(row: NotificationConfig) {
  const latest = await getNotificationConfigApi(row.notification_id);
  editing.value = latest;
  setForm(latest);
  configDrawerApi.setState({ title: '设置通知' }).open();
}

async function save() {
  if (!editing.value) return;
  const payload: NotificationConfigSaveInput = {
    mini_program_enabled: form.mini_program_enabled,
    mini_program_text: form.mini_program_text.trim(),
    sms_enabled: form.sms_enabled,
    sms_text: form.sms_text.trim(),
    wechat_enabled: form.wechat_enabled,
    wechat_text: form.wechat_text.trim(),
  };
  if (
    (payload.wechat_enabled && !payload.wechat_text) ||
    (payload.mini_program_enabled && !payload.mini_program_text) ||
    (payload.sms_enabled && !payload.sms_text)
  ) {
    ElMessage.warning('已开启的发送渠道必须填写固定文本');
    return;
  }
  configDrawerApi.lock();
  try {
    await saveNotificationConfigApi(editing.value.notification_id, payload);
    ElMessage.success('通知配置已保存');
    configDrawerApi.close();
    await gridApi.reload();
  } finally {
    configDrawerApi.unlock();
  }
}

async function onSync(channel: 'mini_program' | 'wechat') {
  try {
    await syncNotificationTemplatesApi(audience.value, channel);
  } catch {
    ElMessage.info('当前仅维护默认发送行为与固定文本，未配置外部凭据，未执行同步。');
  }
}

async function changeAudience(value: number | string) {
  audience.value = value as NotificationAudience;
  await gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElTabs :model-value="audience" class="notification-tabs" @update:model-value="changeAudience">
          <ElTabPane label="通知会员" name="member" />
          <ElTabPane label="通知店铺" name="store" />
        </ElTabs>
        <div class="mb-4 flex flex-wrap gap-3">
          <ElButton type="primary" @click="onSync('mini_program')">同步小程序订阅消息</ElButton>
          <ElButton type="success" @click="onSync('wechat')">同步公众号模板消息</ElButton>
        </div>
        <ElAlert
          class="mb-4"
          description="为每个业务场景设置默认发送渠道和固定文本；外部模板同步仅在完成应用凭据配置后执行。"
          :closable="false"
          title="通知发送规则"
          type="info"
        />
      </template>

      <template #wechat="{ row }">
        <ElSwitch :model-value="row.wechat_enabled === 1" aria-label="公众号模板已启用" disabled />
      </template>
      <template #miniProgram="{ row }">
        <ElSwitch :model-value="row.mini_program_enabled === 1" aria-label="小程序订阅已启用" disabled />
      </template>
      <template #sms="{ row }">
        <ElSwitch :model-value="row.sms_enabled === 1" aria-label="短信通知已启用" disabled />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openConfig(row)">设置</ElButton>
      </template>
    </Grid>

    <ConfigDrawer>
      <ElForm label-width="130px">
        <ElFormItem label="通知类型">
          <ElInput :model-value="form.notice_type" disabled />
        </ElFormItem>
        <ElFormItem label="通知场景说明">
          <ElInput :model-value="form.scene" disabled />
        </ElFormItem>
        <ElFormItem label="公众号模板">
          <ElSwitch v-model="form.wechat_enabled" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="公众号固定文本">
          <ElInput v-model="form.wechat_text" :maxlength="500" :rows="3" show-word-limit type="textarea" />
        </ElFormItem>
        <ElFormItem label="小程序订阅">
          <ElSwitch v-model="form.mini_program_enabled" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="小程序固定文本">
          <ElInput v-model="form.mini_program_text" :maxlength="500" :rows="3" show-word-limit type="textarea" />
        </ElFormItem>
        <ElFormItem label="发送短信">
          <ElSwitch v-model="form.sms_enabled" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="短信固定文本">
          <ElInput v-model="form.sms_text" :maxlength="500" :rows="3" show-word-limit type="textarea" />
        </ElFormItem>
      </ElForm>
    </ConfigDrawer>
  </Page>
</template>

<style scoped>
.notification-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}
</style>
