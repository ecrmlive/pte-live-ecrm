<template>
  <a-card :bordered="false" class="page-card">
    <a-alert
      type="info"
      show-icon
      message="会员价与店铺券"
      description="关闭叠加时：结算用了 SVIP 价则不可用店铺券（平台券仍可用）。秒杀行本身跳过会员价。"
      style="margin-bottom: 16px"
    />
    <a-form layout="vertical" style="max-width: 420px">
      <a-form-item label="SVIP 价可叠加店铺券">
        <a-switch
          v-model:checked="mergeOn"
          checked-children="可叠加"
          un-checked-children="互斥"
          :loading="loading"
          :disabled="!canUpdate"
          @change="onToggle"
        />
      </a-form-item>
    </a-form>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchSvipConfig, updateSvipConfig } from '@/api/svip';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canUpdate = computed(() => auth.hasPerm('svip/update'));

const loading = ref(false);
const mergeOn = ref(false);

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchSvipConfig();
    mergeOn.value = data.svip_coupon_merge === 1;
  } finally {
    loading.value = false;
  }
}

async function onToggle(checked: boolean) {
  if (!canUpdate.value) {
    mergeOn.value = !checked;
    return;
  }
  loading.value = true;
  try {
    await updateSvipConfig(checked ? 1 : 0);
    message.success(checked ? '已允许叠加' : '已设为互斥');
  } catch {
    mergeOn.value = !checked;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
</style>
