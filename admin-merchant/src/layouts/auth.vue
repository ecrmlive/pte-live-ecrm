<script lang="ts" setup>
import { computed, onMounted } from 'vue';

import { AuthPageLayout } from '@vben/layouts';
import { preferences } from '@vben/preferences';

import { $t } from '#/locales';
import { applyLoginPageBranding } from '#/utils/shop-display-name';

const logo = computed(() => preferences.logo.source);
const logoDark = computed(() => preferences.logo.sourceDark);

onMounted(() => {
  applyLoginPageBranding(preferences.logo.source);
});
</script>

<template>
  <AuthPageLayout
    :app-name="$t('authentication.pageTitle')"
    :logo="logo"
    :logo-dark="logoDark"
    :page-description="$t('authentication.pageDesc')"
    :page-title="$t('authentication.pageTitle')"
  >
    <!-- 自定义工具栏 -->
    <!-- <template #toolbar></template> -->
  </AuthPageLayout>
</template>
