<script setup lang="ts">
import type { MerchantPaymentChannel, MerchantPaymentChannelCode } from '#/api/core/payment';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElAlert, ElButton, ElCard, ElForm, ElFormItem, ElInput, ElMessage, ElSkeleton, ElSwitch, ElTag } from 'element-plus';
import { computed, onMounted, reactive, ref } from 'vue';

import { getMerchantPaymentChannelsApi, updateMerchantPaymentChannelApi } from '#/api/core/payment';


const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

type PaymentForm = Record<string, string>;
const loading = ref(true);
const saving = ref(false);
const editing = ref<MerchantPaymentChannelCode>('wechat');
const channels = ref<MerchantPaymentChannel[]>([]);
const form = reactive<PaymentForm>({ enabled: 'false', app_id: '', h5_app_id: '', h5_site_url: '', mch_id: '', api_v3_key: '', serial_no: '', private_key: '', public_key: '', notify_url: '' });
const paymentItems = computed(() => [
  { channel: 'wechat' as const, name: '微信支付', description: '使用本店自己的微信商户号、APIv3 密钥和商户私钥。' },
  { channel: 'alipay' as const, name: '支付宝支付', description: '使用本店自己的支付宝 AppID、应用私钥和支付宝公钥。' },
]);
function item(channel: MerchantPaymentChannelCode) { return channels.value.find((entry) => entry.channel === channel); }
function enabled(channel: MerchantPaymentChannelCode) { return item(channel)?.enabled ?? false; }
function configured(channel: MerchantPaymentChannelCode) { return item(channel)?.configured ?? false; }
async function load() { loading.value = true; try { channels.value = (await getMerchantPaymentChannelsApi()).list ?? []; } finally { loading.value = false; } }
function openConfig(channel: MerchantPaymentChannelCode) {
  editing.value = channel;
  Object.assign(form, { enabled: enabled(channel) ? 'true' : 'false', app_id: '', h5_app_id: '', h5_site_url: '', mch_id: '', api_v3_key: '', serial_no: '', private_key: '', public_key: '', notify_url: '' });
  formDrawerApi.setState({ title: channel === 'wechat' ? '配置本店微信支付' : '配置本店支付宝支付' }).open();
}
async function save() {
  if (form.enabled === 'true') {
    const required = editing.value === 'wechat' ? ['app_id', 'mch_id', 'api_v3_key', 'serial_no', 'private_key', 'notify_url'] : ['app_id', 'private_key', 'public_key', 'notify_url'];
    if (required.some((key) => !form[key]?.trim())) { ElMessage.warning('请填写该支付渠道所需的商户参数'); return; }
  }
  saving.value = true;
  formDrawerApi.lock();
  try { await updateMerchantPaymentChannelApi(editing.value, { ...form }); await load(); formDrawerApi.close(); ElMessage.success('本店支付配置已加密保存'); } finally { saving.value = false; formDrawerApi.unlock(); }
}
onMounted(load);
</script>

<template>
  <Page auto-content-height>
    <ElAlert class="mb-4" type="warning" :closable="false" title="敏感参数仅服务端加密保存" description="私钥、APIv3 密钥不会回显；再次修改已启用渠道时，请重新填写完整商户参数。" />
    <ElSkeleton v-if="loading" :rows="6" animated />
    <div v-else class="grid max-w-4xl gap-4 md:grid-cols-2">
      <ElCard v-for="entry in paymentItems" :key="entry.channel" shadow="never">
        <template #header><div class="flex items-center justify-between gap-3"><span class="font-medium">{{ entry.name }}</span><ElTag :type="enabled(entry.channel) ? 'success' : 'info'">{{ enabled(entry.channel) ? '本店已启用' : '本店未启用' }}</ElTag></div></template>
        <p class="mb-4 text-sm leading-6 text-gray-500">{{ entry.description }}</p>
        <p class="mb-4 text-xs text-gray-400">{{ configured(entry.channel) ? '商户参数已加密保存' : '尚未配置商户参数' }}</p>
        <ElButton type="primary" plain @click="openConfig(entry.channel)">配置本店{{ entry.name }}</ElButton>
      </ElCard>
    </div>
    <FormDrawer>
      <ElForm label-position="top"><ElFormItem label="启用本店渠道"><ElSwitch v-model="form.enabled" active-value="true" inactive-value="false" /></ElFormItem><ElFormItem label="AppID" required><ElInput v-model="form.app_id" /></ElFormItem><template v-if="editing === 'wechat'"><ElFormItem label="微信 H5 支付 AppID"><ElInput v-model="form.h5_app_id" placeholder="用于 H5 MWEB 支付；不与小程序 AppID 混用" /></ElFormItem><ElFormItem label="微信 H5 授权站点（HTTPS）"><ElInput v-model="form.h5_site_url" placeholder="https://mall.example.com" /></ElFormItem><ElFormItem label="微信商户号" required><ElInput v-model="form.mch_id" /></ElFormItem><ElFormItem label="APIv3 密钥" required><ElInput v-model="form.api_v3_key" type="password" show-password /></ElFormItem><ElFormItem label="商户证书序列号" required><ElInput v-model="form.serial_no" /></ElFormItem></template><ElFormItem :label="editing === 'wechat' ? '微信商户私钥' : '支付宝应用私钥'" required><ElInput v-model="form.private_key" type="textarea" :rows="4" /></ElFormItem><ElFormItem v-if="editing === 'alipay'" label="支付宝公钥" required><ElInput v-model="form.public_key" type="textarea" :rows="4" /></ElFormItem><ElFormItem label="支付回调地址" required><ElInput v-model="form.notify_url" placeholder="https://业务域名/api/callback/v1/pay/..." /></ElFormItem></ElForm>
    </FormDrawer>
  </Page>
</template>
