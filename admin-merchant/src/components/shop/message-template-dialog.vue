<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type {
  MessageChannelType,
  MessageFieldItem,
  MessageListItem,
  MessageTemplateSettings,
} from '#/api/core/message-setting';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElInput, ElMessage } from 'element-plus';
import { computed, nextTick, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getMessageFieldListApi,
  saveMessageSettingsApi,
} from '#/api/core/message-setting';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  message: Pick<MessageListItem, 'message_id' | 'message_name'> | null;
  messageType: MessageChannelType;
}>();

const emit = defineEmits<{ success: [] }>();

const CHANNEL_META: Record<
  MessageChannelType,
  { placeholder: string; tips: string[]; title: string; varTip?: string }
> = {
  mt: {
    placeholder: '请填写申请的公众号模板消息 id',
    tips: [
      '公众号模板消息里有的字段才勾选，如果没有请勿勾选。',
      '模板变量替换成公众号模板里的字段。',
    ],
    title: '设置公众号模板消息',
    varTip: '变量填写、例如 {{thing4.DATA}}，只需要填写 thing4。',
  },
  mp: {
    placeholder: '请填写申请的公众号订阅通知 id',
    tips: [
      '公众号订阅通知里有的字段才勾选，如果没有请勿勾选。',
      '模板变量替换成公众号模板里的字段。',
    ],
    title: '设置公众号订阅通知',
    varTip: '变量填写、例如 {{thing4.DATA}}，只需要填写 thing4。',
  },
  wx: {
    placeholder: '请填写申请的微信小程序订阅模板消息 id',
    tips: [
      '微信小程序订阅模板里有的字段才勾选，如果没有请勿勾选。',
      '模板变量替换成微信小程序订阅模板里的字段。',
    ],
    title: '设置小程序消息订阅模板',
    varTip: '变量填写、例如 {{thing4.DATA}}，只需要填写 thing4。',
  },
  sms: {
    placeholder: '请填写申请的短信模板 code',
    tips: ['短信模板里有的字段才勾选，如果没有请勿勾选。', '模板变量替换成短信模板里的字段。'],
    title: '设置短信模板',
  },
};

const loading = ref(false);
const saving = ref(false);
const fieldList = ref<MessageFieldItem[]>([]);
const selectedRows = ref<MessageFieldItem[]>([]);

const dialogTitle = computed(() => {
  const channelMeta = CHANNEL_META[props.messageType];
  const name = props.message?.message_name ?? '';
  return name ? `${channelMeta.title}（${name}）` : channelMeta.title;
});

const meta = computed(() => CHANNEL_META[props.messageType]);

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      class: 'w-full',
      placeholder: meta.value.placeholder,
    },
    fieldName: 'templateId',
    label: '模板id',
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 88,
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

const gridOptions = reactive<VxeGridProps<MessageFieldItem>>({
  checkboxConfig: {
    highlight: true,
    reserve: true,
  },
  columns: [
    { type: 'checkbox', width: 55 },
    { field: 'field_name', minWidth: 120, title: '字段名称' },
    {
      field: 'field_new_ename',
      minWidth: 140,
      slots: { default: 'field_new_ename' },
      title: '模板变量名',
    },
    {
      field: 'filed_new_value',
      minWidth: 160,
      slots: { default: 'filed_new_value' },
      title: '模板内容',
    },
  ],
  minHeight: 280,
  proxyConfig: {
    ajax: {
      query: async () => ({
        items: fieldList.value,
        total: fieldList.value.length,
      }),
    },
  },
  rowConfig: {
    keyField: 'message_field_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [FieldGrid, fieldGridApi] = useVbenVxeGrid({ gridOptions });

function normalizeSettings(raw: MessageTemplateSettings | null | unknown) {
  if (!raw || typeof raw !== 'object') {
    return null;
  }
  const settings = raw as MessageTemplateSettings;
  if (!settings.var_data || typeof settings.var_data !== 'object') {
    return settings.template_id ? settings : null;
  }
  return settings;
}

function initChecked(settings: MessageTemplateSettings | null) {
  if (!settings?.var_data) {
    return;
  }
  const grid = fieldGridApi.grid;
  if (!grid) return;
  Object.keys(settings.var_data).forEach((key) => {
    const saved = settings.var_data[key];
    const field = fieldList.value.find((item) => item.field_ename === key);
    if (!field) {
      return;
    }
    field.field_new_ename = saved?.field_name ?? field.field_new_ename;
    field.filed_new_value = saved?.filed_value ?? field.filed_new_value;
    grid.setCheckboxRow(field, true);
  });
  selectedRows.value = (grid.getCheckboxRecords?.() ?? []) as MessageFieldItem[];
}

function onCheckboxChange() {
  selectedRows.value = (fieldGridApi.grid?.getCheckboxRecords?.() ??
    []) as MessageFieldItem[];
}

async function loadMeta() {
  if (!props.message) {
    return;
  }
  loading.value = true;
  fieldList.value = [];
  selectedRows.value = [];
  try {
    const res = await getMessageFieldListApi(props.message.message_id, props.messageType);
    fieldList.value = (res.list ?? []).map((field) => ({
      ...field,
      field_new_ename: field.field_ename,
      filed_new_value: field.filed_value,
    }));
    const settings = normalizeSettings(res.settings);
    await formApi.setValues({
      templateId: settings?.template_id ?? '',
    });
    await nextTick();
    await fieldGridApi.reload();
    await nextTick();
    initChecked(settings);
  } finally {
    loading.value = false;
  }
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void loadMeta();
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    return;
  }
  modalApi.close();
});

async function saveTemplate() {
  if (!props.message || saving.value) {
    return;
  }
  const values = await formApi.getValues();
  saving.value = true;
  try {
    await saveMessageSettingsApi({
      fieldList: selectedRows.value,
      message_id: props.message.message_id,
      message_type: props.messageType,
      template_id: String(values.templateId ?? '').trim(),
    });
    ElMessage.success('保存成功');
    open.value = false;
    emit('success');
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :title="dialogTitle"
    class="w-[640px] max-w-[95vw]"
  >
    <div v-loading="loading">
      <div class="mb-3 rounded-lg bg-[#909399] p-4 text-sm text-white">
        <p>配置说明：</p>
        <p v-for="(tip, index) in meta.tips" :key="index">{{ index + 1 }}、{{ tip }}</p>
        <p v-if="meta.varTip">3、{{ meta.varTip }}</p>
      </div>
      <Form />
      <FieldGrid
        class="mt-3"
        @checkbox-all="onCheckboxChange"
        @checkbox-change="onCheckboxChange"
      >
        <template #field_new_ename="{ row }">
          <ElInput v-model.trim="row.field_new_ename" size="small" />
        </template>
        <template #filed_new_value="{ row }">
          <ElInput
            v-model.trim="row.filed_new_value"
            :disabled="row.is_var === 1"
            size="small"
          />
        </template>
      </FieldGrid>
    </div>
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="saving" type="primary" @click="saveTemplate">确定</ElButton>
    </template>
  </Modal>
</template>
