<script lang="ts" setup>
import { Page } from '@vben/common-ui';

defineOptions({ name: 'MerchantListPage' });

withDefaults(
  defineProps<{
    /** 弹窗内嵌（不用 Page 壳） */
    embedded?: boolean;
    loading?: boolean;
  }>(),
  {
    embedded: false,
    loading: false,
  },
);
</script>

<template>
  <Page v-if="!embedded" content-class="merchant-list-page">
    <div v-loading="loading" class="merchant-list-page__body">
      <slot />
    </div>
  </Page>

  <div
    v-else
    v-loading="loading"
    class="merchant-list-page merchant-list-page--embedded"
  >
    <slot />
  </div>
</template>
