<script setup lang="ts">
import { onMounted } from 'vue';

import { Page } from '@vben/common-ui';

import { dismissElPoppers } from '#/utils/el-popper-cleanup';

defineOptions({ name: 'NativeListPage' });

onMounted(() => {
  dismissElPoppers();
});
</script>

<template>
  <Page content-class="native-list-page">
    <!-- Hub Tab 壳；内嵌 panel 列表用 Page + useVbenVxeGrid -->
    <div
      class="list-panel flex flex-col"
      :class="{ 'list-panel--fill': Boolean($slots.tabs) }"
      data-native-list-layout="standard"
    >
      <div v-if="$slots.tabs" class="common-seach-wrap flex-shrink-0">
        <slot name="tabs" />
      </div>
      <slot />
    </div>
  </Page>
</template>
