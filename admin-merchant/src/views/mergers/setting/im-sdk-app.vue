<script setup lang="ts">
import type { MerchantIMSDKApp, MerchantIMSDKAppInput } from '#/api/core/im-sdk-app';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElDialog, ElForm, ElFormItem, ElInput, ElMessage, ElOption, ElPopconfirm, ElSelect, ElTable, ElTableColumn, ElTag } from 'element-plus';
import { onMounted, reactive, ref } from 'vue';

import {
  activateMerchantIMSDKAppApi,
  createMerchantIMSDKAppApi,
  disableMerchantIMSDKAppApi,
  getMerchantIMSDKAppsApi,
  updateMerchantIMSDKAppApi,
} from '#/api/core/im-sdk-app';

const rows = ref<MerchantIMSDKApp[]>([]);
const loading = ref(false);
const saving = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const form = reactive<MerchantIMSDKAppInput>({
  sdk_app_id: '', name: '', status: 'disabled', api_public_url: '', ws_public_url: '', pte_profile_id: '',
});

async function load() {
  loading.value = true;
  try { rows.value = (await getMerchantIMSDKAppsApi()).list ?? []; } finally { loading.value = false; }
}
function openCreate() {
  editingID.value = undefined;
  Object.assign(form, { sdk_app_id: '', name: '', status: 'disabled', api_public_url: '', ws_public_url: '', pte_profile_id: '' });
  dialogOpen.value = true;
}
function openEdit(row: any) {
  editingID.value = row.id;
  Object.assign(form, { sdk_app_id: row.sdk_app_id, name: row.name, status: row.status, api_public_url: row.api_public_url, ws_public_url: row.ws_public_url, pte_profile_id: row.pte_profile_id });
  dialogOpen.value = true;
}
async function save() {
  if (!form.sdk_app_id.trim() || !form.name.trim()) { ElMessage.warning('请填写 SDK AppId 和名称'); return; }
  saving.value = true;
  try {
    if (editingID.value) await updateMerchantIMSDKAppApi(editingID.value, { ...form });
    else await createMerchantIMSDKAppApi({ ...form });
    dialogOpen.value = false;
    await load();
    ElMessage.success('IM SDK AppId 已保存');
  } finally { saving.value = false; }
}
async function activate(row: any) {
  await activateMerchantIMSDKAppApi(row.id); await load(); ElMessage.success('已切换当前 IM SDK AppId');
}
async function disable(row: any) {
  await disableMerchantIMSDKAppApi(row.id); await load(); ElMessage.success('IM SDK AppId 已停用');
}
onMounted(load);
</script>

<template>
  <Page title="IM SDK AppId" description="一个商户可维护多个 pte-live-im SDK AppId，但同一时刻只允许一个启用项。">
    <ElAlert class="mb-4" type="warning" :closable="false" title="SDK AppId 不等于店铺 X-AppId" description="店铺 X-AppId 用于七禧业务隔离；本页只管理 pte-live-im 映射。集成令牌和 UserSig 永不在本页面保存或显示。" />
    <div class="mb-4"><ElButton type="primary" @click="openCreate">新增 IM SDK AppId</ElButton></div>
    <ElTable v-loading="loading" :data="rows" border>
      <ElTableColumn prop="name" label="名称" min-width="140" />
      <ElTableColumn prop="sdk_app_id" label="SDK AppId" min-width="180" />
      <ElTableColumn label="状态" width="150"><template #default="{ row }"><ElTag :type="row.is_active ? 'success' : row.status === 'enabled' ? 'warning' : 'info'">{{ row.is_active ? '当前启用' : row.status === 'enabled' ? '可启用' : '已停用' }}</ElTag></template></ElTableColumn>
      <ElTableColumn prop="api_public_url" label="IM API 公网地址" min-width="220" show-overflow-tooltip />
      <ElTableColumn prop="ws_public_url" label="WebSocket 公网地址" min-width="220" show-overflow-tooltip />
      <ElTableColumn label="操作" width="220" fixed="right"><template #default="{ row }"><ElButton link type="primary" @click="openEdit(row)">编辑</ElButton><ElButton v-if="!row.is_active && row.status === 'enabled'" link type="success" @click="activate(row)">设为当前</ElButton><ElPopconfirm title="停用后不能签发新的 IM 凭证，确认继续？" @confirm="disable(row)"><template #reference><ElButton v-if="row.status === 'enabled'" link type="danger">停用</ElButton></template></ElPopconfirm></template></ElTableColumn>
    </ElTable>
    <ElDialog v-model="dialogOpen" :title="editingID ? '编辑 IM SDK AppId' : '新增 IM SDK AppId'" width="640px" destroy-on-close>
      <ElForm label-position="top"><ElFormItem label="pte-live-im SDK AppId" required><ElInput v-model="form.sdk_app_id" :disabled="Boolean(editingID)" /></ElFormItem><ElFormItem label="名称" required><ElInput v-model="form.name" /></ElFormItem><ElFormItem label="状态"><ElSelect v-model="form.status"><ElOption label="已停用" value="disabled" /><ElOption label="可启用" value="enabled" /></ElSelect></ElFormItem><ElFormItem label="IM API 公网地址"><ElInput v-model="form.api_public_url" placeholder="https://im.example.com" /></ElFormItem><ElFormItem label="WebSocket 公网地址"><ElInput v-model="form.ws_public_url" placeholder="wss://im.example.com/ws" /></ElFormItem><ElFormItem label="pte 配置引用"><ElInput v-model="form.pte_profile_id" placeholder="仅引用，不填写集成 Token" /></ElFormItem></ElForm>
      <template #footer><ElButton @click="dialogOpen = false">取消</ElButton><ElButton type="primary" :loading="saving" @click="save">保存</ElButton></template>
    </ElDialog>
  </Page>
</template>
