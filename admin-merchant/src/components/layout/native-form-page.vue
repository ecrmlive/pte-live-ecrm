<script setup lang="ts">
import { Page } from '@vben/common-ui';

defineOptions({ name: 'NativeFormPage' });
</script>

<template>
  <Page auto-content-height content-class="native-form-page">
    <slot />
  </Page>
</template>
