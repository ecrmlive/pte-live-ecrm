<script lang="ts" setup>
import { computed } from 'vue';

import { Icon as IconifyIcon } from '@iconify/vue';

import { normalizePlatformMenuIcon } from '#/constants/platform-lucide-icons';

const props = defineProps<{
  icon?: string;
  size?: number;
}>();

const resolved = computed(() => normalizePlatformMenuIcon(props.icon));
</script>

<template>
  <IconifyIcon
    v-if="resolved"
    :height="size ?? 18"
    :icon="resolved"
    :style="{ verticalAlign: 'middle' }"
    :width="size ?? 18"
    class="text-foreground/80"
  />
  <span v-else class="text-muted-foreground">—</span>
</template>
