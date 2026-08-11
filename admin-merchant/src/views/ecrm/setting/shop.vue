<script setup lang="ts">
import type { MerchantShopProfileInput } from '#/api/core/shop-setting';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElForm, ElFormItem, ElInput, ElMessage, ElSkeleton, ElTag } from 'element-plus';
import { onMounted, reactive, ref } from 'vue';

import {
  getMerchantShopProfileApi,
  updateMerchantShopProfileApi,
} from '#/api/core/shop-setting';

const loading = ref(true);
const saving = ref(false);
const shopID = ref(0);
const status = ref(0);
const form = reactive<MerchantShopProfileInput>({
  mer_address: '',
  mer_info: '',
  mer_name: '',
  mer_phone: '',
  real_name: '',
});

async function load() {
  loading.value = true;
  try {
    const data = await getMerchantShopProfileApi();
    shopID.value = data.mer_id;
    status.value = data.status;
    form.mer_name = data.mer_name ?? '';
    form.real_name = data.real_name ?? '';
    form.mer_phone = data.mer_phone ?? '';
    form.mer_address = data.mer_address ?? '';
    form.mer_info = data.mer_info ?? '';
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!form.mer_name.trim()) {
    ElMessage.warning('请填写店铺名称');
    return;
  }
  saving.value = true;
  try {
    const data = await updateMerchantShopProfileApi({ ...form });
    status.value = data.status;
    ElMessage.success('店铺资料已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page title="店铺设置" description="仅维护当前商户店铺资料；保存操作受 shop/update 按钮权限控制。">
    <ElSkeleton v-if="loading" :rows="8" animated />
    <template v-else>
      <ElAlert class="mb-4" type="info" :closable="false" title="店铺名称、联系人、联系电话、地址和简介会展示在商户侧业务场景中。" />
      <ElForm label-position="top" class="max-w-3xl">
        <div class="mb-5 flex items-center gap-3 text-sm">
          <span>店铺 ID：{{ shopID }}</span>
          <ElTag :type="status === 1 ? 'success' : 'warning'">{{ status === 1 ? '正常营业' : '已停用' }}</ElTag>
        </div>
        <ElFormItem label="店铺名称" required><ElInput v-model="form.mer_name" maxlength="64" show-word-limit /></ElFormItem>
        <ElFormItem label="联系人"><ElInput v-model="form.real_name" maxlength="32" /></ElFormItem>
        <ElFormItem label="联系电话"><ElInput v-model="form.mer_phone" maxlength="32" /></ElFormItem>
        <ElFormItem label="店铺地址"><ElInput v-model="form.mer_address" maxlength="255" /></ElFormItem>
        <ElFormItem label="店铺简介"><ElInput v-model="form.mer_info" type="textarea" :rows="5" maxlength="1000" show-word-limit /></ElFormItem>
        <ElButton type="primary" :loading="saving" @click="save">保存</ElButton>
      </ElForm>
    </template>
  </Page>
</template>
