<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { ElButton, ElEmpty, ElMessage } from 'element-plus';

import PlatformListQueryForm from '#/components/platform-list/PlatformListQueryForm.vue';
import { Page } from '@vben/common-ui';
import TrafficPanel from '#/views/shop/traffic-panel.vue';

const route = useRoute();
const router = useRouter();

const fromShop = computed(() => route.query.from === 'shop');
const activeAppId = ref('');

const queryFormOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '如 10001' },
      defaultValue: '',
      fieldName: 'app_id',
      formItemClass: 'pb-0',
      label: '商城 ID',
    },
  ],
};

function goBack() {
  router.push(fromShop.value ? '/shop/Index' : '/shop/traffic/Index');
}

function handleSubmit(values: Record<string, unknown>) {
  const id = Number(values.app_id);
  if (!id) {
    ElMessage.warning('请输入商城 ID');
    return;
  }
  activeAppId.value = String(id);
}

function handleReset() {
  activeAppId.value = '';
}

onMounted(() => {
  const q = route.query.app_id;
  if (q) {
    activeAppId.value = String(q);
  }
});
</script>

<template>
  <Page>
    <PlatformListQueryForm
      :options="{
        ...queryFormOptions,
        ...(route.query.app_id
          ? { initialValues: { app_id: String(route.query.app_id) } }
          : {}),
      }"
      class="mb-3"
      @reset="handleReset"
      @submit="handleSubmit"
    />

    <div class="mb-3">
      <ElButton size="small" @click="goBack">
        ← {{ fromShop ? '返回商城列表' : '返回商城流量' }}
      </ElButton>
    </div>

    <TrafficPanel v-if="activeAppId" :key="activeAppId" :app-id="activeAppId" />

    <ElEmpty v-else description="请输入商城 ID 查询">
      <ElButton type="primary" @click="router.push('/shop/traffic/Index')">
        查看全部商城流量
      </ElButton>
    </ElEmpty>
  </Page>
</template>
